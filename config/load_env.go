package config 

import (
	"os"
)

var (
	FileSecretToken = os.Getenv("FILESECRET_TOKEN")
)