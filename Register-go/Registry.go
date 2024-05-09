package Register_go

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Registry struct {
	nameToId           map[string][]string
	services           map[string]ServiceInfo
	clientsWithService map[string][]string
	clients            map[string]int64
}

var GlobalRegistry = &Registry{
	nameToId:           make(map[string][]string),
	services:           make(map[string]ServiceInfo),
	clientsWithService: make(map[string][]string),
}

// IdGenerate @cs string("client" or "server")
func IdGenerate() string {
	key := time.Now().Format("20060102150405.0000000") + strconv.Itoa(rand.Intn(100000))
	data := []byte(key)
	md := md5.Sum(data)
	return hex.EncodeToString(md[:])
}

func (rg Registry) dump() {
	services := make([]ServiceInfo, 0, len(rg.services))
	for _, info := range rg.services {
		services = append(services, info)
	}
	jsonData, err := json.Marshal(services)
	if err != nil {
		log.Fatal(err)
	}
	f, err1 := os.Create("services.json")
	if err1 != nil {
		log.Fatal(err1)
	}
	defer f.Close()
	_, err = f.Write(jsonData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("dump services.json")
}

func (rg Registry) load() {
	f, err := os.Open("services.json")
	if err != nil {
		log.Fatal(err)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	var services []ServiceInfo
	err = json.Unmarshal(data, &services)
	if err != nil {
		log.Fatal(err)
	}
	for _, info := range services {
		rg.services[info.Id] = info
	}
	fmt.Println("load services.json")
}
