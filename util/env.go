package util

import "os"

func GetEnv(env string, defaultVal string) string {
	res := os.Getenv(env)
	if len(res) == 0 {
		return defaultVal
	}
	return res
}
