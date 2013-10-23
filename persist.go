package main

import (
	"github.com/cznic/kv"
	"fmt"
	"os"
)

var DB *kv.DB

func Remember(key string) (value string) {
	theKey := Config.ourNetwork + ":" + key
	fmt.Printf("Get %s\n", theKey)
	val, err := DB.Get(nil, []byte(theKey))
	if err != nil {
		fmt.Printf("Failed DB.Get(); %s", err)
	}
	return string(val)
}

func Persist(key, value string) {
	theKey := Config.ourNetwork + ":" + key
	fmt.Printf("Persist %s = %s\n", theKey, value)
	if err := DB.Set([]byte(theKey), []byte(value)); err != nil {
		fmt.Printf("Failed DB.Set(); %s", err)
	}
}

func InitPersist() {
	dbPath := os.Getenv("HOME") + "/.gofoot/db"
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Println("No database. Making one...")

		if DB, err = kv.Create(dbPath, &kv.Options{}); err != nil {
			fmt.Printf("Could not create db %s:\n", err)
		}
	} else {
		if DB, err = kv.Open(dbPath, &kv.Options{}); err != nil {
			fmt.Printf("Could not open db %s:\n", err)
		}
	}
}
