package Register_go

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type Register struct {
	nameToId map[string][]string
	services map[string]ServiceInfo
	clients  map[string]int
}

var GlobalRegister = &Register{
	nameToId: make(map[string][]string),
	services: make(map[string]ServiceInfo),
	clients:  make(map[string]int),
}

func IdGenerate() string {
	key := strconv.Itoa(len(GlobalRegister.services) + 100000000)
	data := []byte(key)
	md := md5.New()
	md.Write(data)
	return hex.EncodeToString(md.Sum(nil))
}

func (rg Register) dump() {
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

func (rg Register) load() {
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
