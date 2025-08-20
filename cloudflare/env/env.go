//go:build js && wasm

package env

import (
	"os"

	jsclass "github.com/Darckfast/workers-go/internal/class"
	jsconv "github.com/Darckfast/workers-go/internal/conv"
)

func LoadEnvs() error {
	if !jsclass.Env.Truthy() {
		return nil
	}

	envs, err := jsconv.JSValueToMap(jsclass.Env)

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
