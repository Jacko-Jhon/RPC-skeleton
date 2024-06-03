package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
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
	var err error
	RegisterSocket, err = net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(RegisterIP),
		Port: RegisterPort,
	})
	if err != nil {
		panic(err)
	}
	defer func(RegisterSocket *net.UDPConn) {
		err := RegisterSocket.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(RegisterSocket)
	buf := DailRegistry("1")
	var res ServiceList
	err = json.Unmarshal(buf, &res)
	if err != nil {
		panic(err)
	}
	if res.Status {
		fmt.Println("Online Services:")
		for i, name := range res.List {
			fmt.Println(strconv.Itoa(i+1)+",", name)
		}
	} else {
		fmt.Println(res.Info)
	}
}

func GenerateCode() {
	path := "./Client"
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}
	file1, err := os.Create(path + "/error.go")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file1)
	_, err = file1.WriteString(error_tpl)
	if err != nil {
		panic(err)
	}
	file2, err := os.Create(path + "/util.go")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file2)
	_, err = file2.WriteString(util_tpl)
	if err != nil {
		panic(err)
	}
	file3, err := os.Create(path + "/struct.go")
	if err != nil {
		panic(err)
	}
	str := generateStruct()
	tpl := struct_tpl + strings.Join(str, "\n")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file3)
	_, err = file3.WriteString(tpl)
	if err != nil {
		panic(err)
	}
	file4, err := os.Create(path + "/client.go")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file4)
	nl, caseL := generateCall()
	tpl2 := strings.Replace(client_tpl, "$case$", nl, -1)
	tpl3 := strings.Replace(call_tpl, "$case$", caseL, -1)
	tpl2 = tpl2 + tpl3
	_, err = file4.WriteString(tpl2)
	if err != nil {
		panic(err)
	}
}

func GetServices(sv string) {
	var err error
	RegisterSocket, err = net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(RegisterIP),
		Port: RegisterPort,
	})
	if err != nil {
		panic(err)
	}
	defer func(RegisterSocket *net.UDPConn) {
		err := RegisterSocket.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(RegisterSocket)
	split := strings.Split(sv, " ")
	for _, s := range split {
		if _, ok := MyServices[s]; ok {
			fmt.Println("Service", s, "already got")
			continue
		}
		buf := DailRegistry("0", s)
		var res ServiceUrls
		err = json.Unmarshal(buf, &res)
		if res.Status {
			fmt.Println("Service", res.Name, "already got")
			MyServices[sv] = ServiceBrief{
				Name: res.Name,
				Args: res.Args,
				Ret:  res.Ret,
			}
		} else {
			fmt.Println(res.Info)
		}
	}
	dump()
}

func UnGetServices(ug string) {
	split := strings.Split(ug, " ")
	for _, s := range split {
		if s == "-all" {
			var input string
			fmt.Println("Are you sure to delete all services?(y/n)")
			_, err := fmt.Scan(&input)
			if err != nil {
				return
			}
			if input != "y" {
				continue
			}
			for k := range MyServices {
				delete(MyServices, k)
			}
			break
		}
		if _, ok := MyServices[s]; ok {
			delete(MyServices, s)
			fmt.Println("Service", s, "already deleted")
		} else {
			fmt.Println("Service", s, "not found")
		}
	}
	dump()
}

func generateStruct() []string {
	var structs []string
	for name, service := range MyServices {
		title := strings.ToUpper(name[0:1]) + name[1:] + "Args"
		var data []string
		for _, arg := range service.Args {
			spl := strings.Split(arg, " ")
			a := strings.ToUpper(spl[0][0:1]) + spl[0][1:]
			js := "`json:\"" + spl[0] + "\"`"
			m := "\t" + a + "\t" + spl[1] + "\t" + js
			data = append(data, m)
		}
		for _, ret := range service.Ret {
			spl := strings.Split(ret, " ")
			a := strings.ToUpper(spl[0][0:1]) + spl[0][1:]
			js := "`json:\"" + spl[0] + "\"`"
			m := "\t" + a + "\t" + spl[1] + "\t" + js
			data = append(data, m)
		}
		mergeD := strings.Join(data, "\n")
		mergeA := "type " + title + " struct {\n" + mergeD + "\n}"
		structs = append(structs, mergeA)
	}
	return structs
}

func generateCall() (string, string) {
	nl := make([]string, 0)
	caseL := make([]string, 0)
	for name := range MyServices {
		title := strings.ToUpper(name[0:1]) + name[1:] + "Args"
		n := "\"" + name + "\""
		nl = append(nl, n)
		fn := "\tcase " + n + ":\n" +
			"\t\targs := ArgsType.(*" + title + ")\n" +
			"\t\tbuff, err := json.Marshal(args)\n\t\tfatalError(err)\n\t\tcall(&buff, ServiceName)\n" +
			"\t\tvar res " + title + "\n" +
			"\t\terr = json.Unmarshal(buff, &res)\n\t\tfatalError(err)\n\t\tArgsType = res"
		caseL = append(caseL, fn)
	}
	res1 := strings.Join(nl, ", ")
	res2 := strings.Join(caseL, "\n")
	return res1, res2
}
