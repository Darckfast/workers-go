//go:build js && wasm

package httpr2

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
	mathr "math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/Darckfast/workers-go/cloudflare/r2"
)

const BUFFER_SIZE int = 4096
const IV_SIZE int = 16

func EncryptFileCTR(in io.Reader, out *bytes.Buffer) (int64, error) {
	block, err := aes.NewCipher([]byte(os.Getenv("ENC_KEY")))
	if err != nil {
		return 0, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return 0, err
	}

	var size int64

	n, err := out.Write(iv)

	if err != nil {
		return 0, err
	}

	size += int64(n)

	stream := cipher.NewCTR(block, iv)
	h := hmac.New(sha256.New, []byte(os.Getenv("HMAC_KEY")))
	h.Write(iv)

	buf := make([]byte, BUFFER_SIZE)

	for {
		n, readErr := in.Read(buf)
		if n > 0 {
			encrypted := make([]byte, n)
			stream.XORKeyStream(encrypted, buf[:n])
			h.Write(encrypted)
			n, err := out.Write(encrypted)

			if err != nil {
				return 0, err
			}

			size += int64(n)
		}

		if readErr == io.EOF {
			break
		}

		if readErr != nil {
			return 0, readErr
		}
	}

	n, err = out.Write(h.Sum(nil))
	if err != nil {
		return 0, err
	}

	size += int64(n)
	return size, nil
}

func DecryptFileCTR(in io.Reader, out io.Writer, size int64) error {
	dataSize := size - int64(IV_SIZE)
	block, err := aes.NewCipher([]byte(os.Getenv("ENC_KEY")))

	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(in, iv); err != nil {
		return err
	}

	stream := cipher.NewCTR(block, iv)
	h := hmac.New(sha256.New, []byte(os.Getenv("HMAC_KEY")))
	h.Write(iv)

	buf := make([]byte, BUFFER_SIZE)
	readBytes := int64(sha256.Size)
	var totalRead int64

	for readBytes < dataSize {
		toRead := BUFFER_SIZE

		if dataSize-readBytes < int64(BUFFER_SIZE) {
			toRead = int(dataSize - readBytes)
		}

		n, err := io.ReadFull(in, buf[:toRead])

		totalRead += int64(n)

		if err != nil {
			return err
		}

		h.Write(buf[:n])
		decrypted := make([]byte, n)
		stream.XORKeyStream(decrypted, buf[:n])
		if _, err := out.Write(decrypted); err != nil {
			return err
		}

		readBytes += int64(n)
	}

	expectedHMAC := make([]byte, sha256.Size)

	if _, err := io.ReadFull(in, expectedHMAC); err != nil {
		return err
	}

	if !hmac.Equal(h.Sum(nil), expectedHMAC) {
		return errors.New("HMAC verification failed — file may be corrupted or tampered with")
	}

	return nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[mathr.Intn(len(letterRunes))]
	}
	return string(b)
}

var POST_ENC_R2 = func(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	encreader, encwriter := io.Pipe()
	buf := bytes.Buffer{}

	defer r.Body.Close()
	size, err := EncryptFileCTR(r.Body, &buf)

	go func() {
		defer encwriter.Close()
		io.Copy(encwriter, &buf)
	}()

	bucket, _ := r2.NewBucket("TEST_BUCKET")
	itemKey := RandStringRunes(32)

	o, err := bucket.Put(itemKey, encreader, size, &r2.PutOptions{
		CustomMetadata: map[string]string{
			"time":     time.Now().String(),
			"go":       runtime.Version(),
			"os":       runtime.GOOS,
			"arch":     runtime.GOARCH,
			"cpus":     strconv.Itoa(runtime.NumCPU()),
			"duration": strconv.Itoa(int(time.Since(start).Milliseconds())),
		},
		HTTPMetadata: r.Header,
	})

	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": err.Error(),
		})

		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"data": o,
	})
}

var GET_ENC_R2 = func(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	if key == "" {
		w.WriteHeader(400)
		return
	}
	bucket, _ := r2.NewBucket("TEST_BUCKET")
	o, err := bucket.Get(key, nil)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("content-type", "image/png")
	err = DecryptFileCTR(o.Body, w, int64(o.Size))

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))

		return
	}

}
