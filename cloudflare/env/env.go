package cloudflare

import (
	"os"

	jsutil "github.com/syumai/workers/internal/utils"
)

func init() {
	if jsutil.RuntimeEnv.IsUndefined() {
		return
	}

	envList := jsutil.ObjectClass.Call("entries", jsutil.RuntimeEnv)
	for i := range envList.Length() {
		envPair := envList.Index(i)
		envKey := envPair.Index(0).String()
		envValue := envPair.Index(1)

		if envValue.Type().String() == "string" {
			os.Setenv(envKey, envValue.String())
		}
	}
}
