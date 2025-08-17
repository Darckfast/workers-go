package env

import (
	"os"

	jsconv "github.com/syumai/workers/internal/conv"
	jsutil "github.com/syumai/workers/internal/utils"
)

func LoadEnvs() {
	if jsutil.RuntimeEnv.IsUndefined() {
		return
	}

	envs := jsconv.JSValueToMap(jsutil.RuntimeEnv)

	for key := range envs {
		envValue, ok := envs[key].(string)

		if ok {
			os.Setenv(key, envValue)
		}
	}
}
