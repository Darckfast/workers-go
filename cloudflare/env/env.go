//go:build js && wasm

package env

import (
	"os"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	jsconv "github.com/Darckfast/workers-go/internal/conv"
)

func LoadEnvs() {
	if jsclass.Env.IsUndefined() {
		return
	}

	envs := jsconv.JSValueToMap(jsclass.Env)

	for key := range envs {
		envValue, ok := envs[key].(string)

		if ok {
			os.Setenv(key, envValue)
		}
	}
}
