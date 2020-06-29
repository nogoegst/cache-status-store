package main

import (
	"flag"
	"log"

	cachestatusstore "github.com/nogoegst/cache-status-store"
)

func run() error {
	var debug = flag.Bool("debug", false, "print progress")
	var url = flag.String("url", "https://1.1.1.1", "host")
	var password = flag.String("password", "", "unique password")
	var write = flag.String("write", "", "string to write")
	var readlength = flag.Int64("read", 0, "length of the string to read")
	flag.Parse()

	cache := cachestatusstore.NewCloudFlareCache(*url)
	storage := cachestatusstore.NewStorage(cache)
	storage.PrintDebugBits = *debug
	encryptedStorage := cachestatusstore.NewEncryptedStorage(storage)

	if *write != "" {
		if err := encryptedStorage.SetBytes([]byte(*password), []byte(*write)); err != nil {
			return err
		}
	} else {
		out, err := encryptedStorage.GetBytes([]byte(*password), *readlength)
		if err != nil {
			return err
		}
		log.Printf("read out the string back: %s, %x", out, out)

	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
