package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func dump() {
	fp := "./ClientConfig.json"
	f, err1 := os.Create(fp)
	if err1 != nil {
		log.Fatal(err1)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)
	jsonData, err := json.Marshal(MyServices)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write(jsonData)
	if err != nil {
		log.Fatal(err)
	}
}

func load() {
	fp := "./ClientConfig.json"
	f, err := os.Open(fp)
	if err != nil {
		return
	}
	data, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)
	err = json.Unmarshal(data, &MyServices)
	if err != nil {
		log.Fatal(err)
	}
}
