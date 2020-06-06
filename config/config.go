package config

import (
	"golang.org/x/xerrors"
	"google.golang.org/api/option"
	"log"
	"os"
)

var (
	ProjectID        string
	ProjectNo        string
	GCPClientOptions []option.ClientOption
)

func Configure() {
	var err error
	ProjectID, err = MustGetEnv("PROJECT_ID")
	if err != nil {
		log.Fatal(err)
	}
	ProjectNo, err = MustGetEnv("PROJECT_NO")
	if err != nil {
		log.Fatal(err)
	}

	isRelase := os.Getenv("RELEASE") == "release"
	if !isRelase {
		GCPClientOptions = []option.ClientOption{
			option.WithCredentialsJSON([]byte(os.Getenv("SERVICE_ACCOUNT_JSON"))),
		}
	}

}

func MustGetEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", xerrors.Errorf("Not exists %s", key)
	}
	return value, nil
}
