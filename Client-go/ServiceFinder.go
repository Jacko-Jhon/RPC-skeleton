package main

import (
	"flag"
	"fmt"
	"net"
)

var (
	MyServices     map[string]ServiceBrief
	RegisterSocket *net.UDPConn
	RegisterIP     string
	RegisterPort   int
	Online         bool
	GetList        bool
	Generate       bool
	Get            string
	UnGet          string
)

func main() {
	flag.StringVar(&RegisterIP, "i", "127.0.0.1", "Set register ip")
	flag.IntVar(&RegisterPort, "p", 8888, "Set register port")
	flag.BoolVar(&GetList, "l", false, "Get my services list")
	flag.StringVar(&Get, "get", "", "Get services")
	flag.StringVar(&UnGet, "un-get", "", "Delete form my services, \nif you want to delete all, input '-all'")
	flag.BoolVar(&Online, "o-l", false, "Get online service list")
	flag.BoolVar(&Generate, "gen", false, "Generate code for your project")
	flag.Parse()
	load()
	if GetList && !Online && !Generate && Get == "" && UnGet == "" {
		ListMyServices()
	} else if !GetList && Online && !Generate && Get == "" && UnGet == "" {
		ListOnlineServices()
	} else if !GetList && !Online && Generate && Get == "" && UnGet == "" {
		GenerateCode()
	} else if !GetList && !Online && !Generate && Get != "" && UnGet == "" {
		GetServices(Get)
	} else if !GetList && !Online && !Generate && Get == "" && UnGet != "" {
		UnGetServices(UnGet)
	} else {
		fmt.Println("Command Error!")
	}
}
