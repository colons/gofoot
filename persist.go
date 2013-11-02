package main

import (
	"github.com/cznic/kv"
	"fmt"
	"os"
	"os/signal"
)

var DB *kv.DB

func Remember(key string) (value string) {
	fmt.Printf("Get %s\n", key)
	val, err := DB.Get(nil, []byte(key))
	if err != nil {
		fmt.Printf("Failed DB.Get(); %s", err)
	}
	return string(val)
}

func Persist(key, value string) {
	fmt.Printf("Persist %s = %s\n", key, value)
	DB.BeginTransaction()
	if err := DB.Set([]byte(key), []byte(value)); err != nil {
		fmt.Printf("Failed DB.Set(); %s", err)
	}
	DB.Commit()
}

func InitPersist(network string) {
	dbPath := os.Getenv("HOME") + "/.gofoot/" + network + ".db"
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Println("No database. Making one...")

		if DB, err = kv.Create(dbPath, &kv.Options{}); err != nil {
			fmt.Printf("Could not create db %s:\n", err)
			os.Exit(1)
		}
	} else {
		if DB, err = kv.Open(dbPath, &kv.Options{}); err != nil {
			fmt.Printf("Could not open db %s:\n", err)
			os.Exit(1)
		}
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		for sig := range c {
			if err := DB.Close(); err != nil {
				fmt.Printf("Could not close db %s:\n", err)
			}
			fmt.Println(sig)
			os.Exit(0)
		}
	}()
}
