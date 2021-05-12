package dotenvUtil

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kingokeke/go-example-1/utils"
)

func GetValue(key string) string {
	dotEnvFile := ".env"

	if strings.HasPrefix(os.Getenv("APP_ENV"), "prod") {
		dotEnvFile = ".env.prod"
	}

	e := godotenv.Load(dotEnvFile)
	utils.CheckError(e)

	return os.Getenv(key)
}