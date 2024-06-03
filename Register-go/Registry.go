package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Registry struct {
	nameToId map[string]*[]string
	services map[string]*ServiceInfo
}

var GlobalRegistry = &Registry{
	nameToId: make(map[string]*[]string),
	services: make(map[string]*ServiceInfo),
}

func IdGenerate() string {
	key := time.Now().Format("20060102150405.0000000") + strconv.Itoa(rand.Intn(100000))
	data := []byte(key)
	md := md5.Sum(data)
	return hex.EncodeToString(md[:])
}

// dump 将服务列表写入文件
func (rg *Registry) dump(path string) {
	fp := filepath.Clean(path)
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
	services := make([]ServiceInfo, 0, len(rg.services))
	for _, info := range rg.services {
		services = append(services, *info)
	}
	jsonData, err := json.Marshal(services)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write(jsonData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("dump services.json")
}

// load 从文件中读取服务列表
func (rg *Registry) load(path string) {
	fp := filepath.Clean(path)
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
	var services []ServiceInfo
	err = json.Unmarshal(data, &services)
	if err != nil {
		log.Fatal(err)
	}
	for _, info := range services {
		rg.services[info.Id] = &info
		if _, ok := rg.nameToId[info.Name]; ok {
			*rg.nameToId[info.Name] = append(*rg.nameToId[info.Name], info.Id)
		} else {
			rg.nameToId[info.Name] = &[]string{info.Id}
		}
	}
	fmt.Println("load services.json")
}
