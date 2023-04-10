package config

import (
	"os"
)

var (
	FileSecretToken = os.Getenv("FILESECRET_TOKEN")
	TIME_ZONE		= os.Getenv("TIME_ZONE")
)
