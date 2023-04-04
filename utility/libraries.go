package utility

import (
	"math/rand"
	"os"
	"io/ioutil"
)

func RandSecret() ([]byte, error){
	buf := make([]byte, 128)
	_, err := rand.Read(buf)
	if err != nil {
		return buf, err
	} else {
		return buf, nil
	}
}

func WriteSecretToFile(secret []byte, namefile string) error {
	file, err := os.OpenFile(namefile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer file.Close()
	_, err = file.Write(secret)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func ReadFile(filename string) ([]byte, error) {
    // Read the entire contents of the file into memory
    content, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    return content, nil
}