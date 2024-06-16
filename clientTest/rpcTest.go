package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"test/Client"
	"time"
)

var w sync.WaitGroup
var rd = rand.New(rand.NewSource(time.Now().UnixMicro()))
var test = 100000
var sleep = 10

func add_test() {
	t1 := time.Now().UnixMilli()
	for i := 0; i < test; i++ {
		a := rd.Float32() * 100
		b := rd.Float32() * 100
		arg := Client.AddArgs{A: a, B: b, Res: 0}
		Client.Call("add", &arg)
		if arg.Res != a+b {
			fmt.Println("error")
		}
		time.Sleep(time.Duration(sleep) * time.Millisecond)
	}
	t2 := time.Now().UnixMilli()
	w.Done()
	fmt.Println("add done, time:", t2-t1, "ms realtime:", t2-t1-int64(sleep*test), "ms")
}

func sub_test() {
	t1 := time.Now().UnixMilli()
	for i := 0; i < test; i++ {
		a := rd.Float32() * 100
		b := rd.Float32() * 100
		arg := Client.SubArgs{A: a, B: b, Res: 0}
		Client.Call("sub", &arg)
		if arg.Res != a-b {
			fmt.Println("error")
		}
		time.Sleep(time.Duration(sleep) * time.Millisecond)
	}
	t2 := time.Now().UnixMilli()
	w.Done()
	fmt.Println("sub done, time:", t2-t1, "ms realtime:", t2-t1-int64(sleep*test), "ms")
}

func mul_test() {
	t1 := time.Now().UnixMilli()
	for i := 0; i < test; i++ {
		a := rd.Float32() * 100
		b := rd.Float32() * 100
		arg := Client.MulArgs{A: a, B: b, Res: 0}
		Client.Call("mul", &arg)
		if arg.Res != a*b {
			fmt.Println("error")
		}
		time.Sleep(time.Duration(sleep) * time.Millisecond)
	}
	t2 := time.Now().UnixMilli()
	w.Done()
	fmt.Println("mul done, time:", t2-t1, "ms realtime:", t2-t1-int64(sleep*test), "ms")
}

func div_test() {
	t1 := time.Now().UnixMilli()
	for i := 0; i < test; i++ {
		a := rd.Float32() * 100
		b := rd.Float32() * 100
		for b == 0 {
			b = rd.Float32() * 100
		}
		arg := Client.DivArgs{A: a, B: b, Res: 0}
		Client.Call("div", &arg)
		if arg.Res != a/b {
			fmt.Println("error")
		}
		time.Sleep(time.Duration(sleep) * time.Millisecond)
	}
	t2 := time.Now().UnixMilli()
	w.Done()
	fmt.Println("div done, time:", t2-t1, "ms realtime:", t2-t1-int64(sleep*test), "ms")
}

func MergeSort_test() {
	v1 := make([]float32, 100)
	t1 := time.Now().UnixMilli()
	for i := 0; i < test/100; i++ {
		for j := 0; j < 100; j++ {
			v1[j] = rand.Float32() * 100
		}
		m := Client.MergeSortArgs{Vector: v1}
		Client.Call("MergeSort", &m)
		BubbleSort(v1)
		if !cmp100(v1, m.Res) {
			fmt.Println("error")
		}
		time.Sleep(time.Duration(sleep) * time.Millisecond)
	}
	t2 := time.Now().UnixMilli()
	w.Done()
	fmt.Println("mergesort done, time:", t2-t1, "ms realtime:", t2-t1-int64(sleep*test/100), "ms")
}

func QSort_test() {
	v1 := make([]float32, 100)
	t1 := time.Now().UnixMilli()
	for i := 0; i < test/100; i++ {
		for j := 0; j < 100; j++ {
			v1[j] = rand.Float32() * 100
		}
		q := Client.QSortArgs{Vector: v1}
		Client.Call("QSort", &q)
		BubbleSort(v1)
		if !cmp100(v1, q.Res) {
			fmt.Println("error")
		}
		time.Sleep(time.Duration(sleep) * time.Millisecond)
	}
	t2 := time.Now().UnixMilli()
	w.Done()
	fmt.Println("qsort done, time:", t2-t1, "ms realtime:", t2-t1-int64(sleep*test/100), "ms")
}

func BubbleSort(vector []float32) {
	for i := 0; i < len(vector); i++ {
		for j := 0; j < len(vector)-i-1; j++ {
			if vector[j] > vector[j+1] {
				vector[j], vector[j+1] = vector[j+1], vector[j]
			}
		}
	}
}

func cmp100(a []float32, b []float32) bool {
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func main() {
	var ip string
	var port int
	flag.StringVar(&ip, "ip", "", "server ip")
	flag.IntVar(&port, "port", 0, "server port")
	flag.Parse()
	Client.Init(ip, port)
	w.Add(6)
	go add_test()
	go sub_test()
	go mul_test()
	go div_test()
	go MergeSort_test()
	go QSort_test()
	w.Wait()
	fmt.Println("all done, test =", test, "sleep =", sleep)
	Client.Close()
}
