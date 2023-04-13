package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sport/config"
)

func main() {

	mode := os.Args[1]

	key, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	req, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatal(err)
	}

	switch mode {
	case "e":
		nonce := make([]byte, gcm.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			log.Fatal(err)
		}

		etext := gcm.Seal(nil, nonce, req, nil)

		fmt.Printf("%s.%s",
			base64.StdEncoding.EncodeToString(nonce),
			base64.StdEncoding.EncodeToString(etext))
	case "d":
		d, err := config.Decrypt(string(req), key)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(d)
	default:
		log.Fatal("invalid mode")
	}
}
