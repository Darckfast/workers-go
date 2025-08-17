package env

import (
	"os"

	jsclass "github.com/syumai/workers/internal/class"
	jsconv "github.com/syumai/workers/internal/conv"
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
