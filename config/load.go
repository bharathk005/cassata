package config

import (
	"fmt"
	"log"
	"os"
)

type Secret string

const (
	JWT_SECRET_KEY Secret = "JWT_SECRET_KEY"
	DATABASE_DSN   Secret = "DATABASE_DSN"
)

var requiredSecrets = []Secret{
	JWT_SECRET_KEY,
	DATABASE_DSN,
}

var Env = make(map[Secret]string)

func LoadConfig() {
	for _, secret := range requiredSecrets {
		if os.Getenv(string(secret)) == "" {
			log.Fatalf("Required environment variable %s is not set", secret)
		}
		Env[secret] = os.Getenv(string(secret))
	}
	fmt.Println("All required environment variables are set")
}
