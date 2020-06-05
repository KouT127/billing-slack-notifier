package config

import (
	"cloud.google.com/go/compute/metadata"
	"log"
)

var (
	ProjectID string
	ProjectNo string
)

func Configure() {
	var err error
	ProjectID, err = metadata.ProjectID()
	if err != nil {
		log.Fatal(err)
	}
	ProjectNo, err = metadata.NumericProjectID()
	if err != nil {
		log.Fatal(err)
	}
}
