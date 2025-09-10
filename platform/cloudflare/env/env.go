//go:build js && wasm

package env

import (
	"os"

	jsconv "github.com/Darckfast/workers-go/internal/conv"
	"github.com/Darckfast/workers-go/platform/cloudflare/lifecycle"
)

func LoadEnvs() error {
	if !lifecycle.Env.Truthy() {
		return nil
	}

	envs, err := jsconv.JSValueToMap(lifecycle.Env)

	if err != nil {
		return err
	}

	for key := range envs {
		envValue, ok := envs[key].(string)

		if ok {
			err = os.Setenv(key, envValue)

			if err != nil {
				return err
			}
		}
	}

	return nil
}
