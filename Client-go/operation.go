package main

import (
	"fmt"
	"strconv"
)

func ListMyServices() {
	fmt.Println("My Services:")
	i := 1
	for name := range MyServices {
		fmt.Println(strconv.Itoa(i)+",", name)
		i++
	}
}

func ListOnlineServices() {

}

func GenerateCode() {

}

func GetServices(sv string) {

}

func UnGetServices(ug string) {

}
