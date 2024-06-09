package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"net"
	"strconv"
	"time"
)

func add(args []byte) ([]byte, error) {
	var argType AddArgs
	err := json.Unmarshal(args, &argType)
	if err != nil {
		return nil, err
	}

	ret := AddArgs{
		Res: argType.A + argType.B,
	}

	js, _ := json.Marshal(ret)
	return js, nil
}

func sub(args []byte) ([]byte, error) {
	var argType SubArgs
	err := json.Unmarshal(args, &argType)
	if err != nil {
		return nil, err
	}

	ret := SubArgs{
		Res: argType.A - argType.B,
	}

	js, _ := json.Marshal(ret)
	return js, nil
}

func mul(args []byte) ([]byte, error) {
	var argType MulArgs
	err := json.Unmarshal(args, &argType)
	if err != nil {
		return nil, err
	}

	ret := MulArgs{
		Res: argType.A * argType.B,
	}

	js, _ := json.Marshal(ret)
	return js, nil
}

func div(args []byte) ([]byte, error) {
	var argType DivArgs
	err := json.Unmarshal(args, &argType)
	if err != nil {
		return nil, err
	}

	if argType.B == 0 {
		return nil, err
	}
	ret := DivArgs{
		Res: argType.A / argType.B,
	}

	js, _ := json.Marshal(ret)
	return js, nil
}

func qSort(args []byte) ([]byte, error) {
	var argType SortArgs
	err := json.Unmarshal(args, &argType)
	if err != nil {
		return nil, err
	}

	_qSort(argType.Vector, 0, len(argType.Vector)-1)
	ret := SortArgs{
		Res: argType.Vector,
	}

	js, _ := json.Marshal(ret)
	return js, nil
}

func _qSort(floats []float32, left, right int) {
	if left >= right {
		return
	}
	p := floats[left]
	i := left
	j := right
	for i < j {
		for floats[j] > p {
			j--
		}
		for floats[i] < p {
			i++
		}
		swap(&floats[i], &floats[j])
	}
	_qSort(floats, left, i)
	_qSort(floats, i+1, right)
}

func swap(a, b *float32) {
	tmp := *a
	*a = *b
	*b = tmp
}

func mergeSort(args []byte) ([]byte, error) {
	var argType SortArgs
	err := json.Unmarshal(args, &argType)
	if err != nil {
		return nil, err
	}

	_mergeSort(argType.Vector, 0, len(argType.Vector)-1)
	ret := SortArgs{
		Res: argType.Vector,
	}

	js, _ := json.Marshal(ret)
	return js, nil
}

func _mergeSort(v []float32, left, right int) {
	if left <= right {
		return
	}
	mid := (left + right) / 2
	_mergeSort(v, left, mid)
	_mergeSort(v, mid+1, right)
	tmp := make([]float32, right-left+1)
	l := left
	r := mid + 1
	i := 0
	for ; l <= mid && r <= right; i++ {
		if v[l] < v[r] {
			tmp[i] = v[l]
			l++
		} else {
			tmp[i] = v[r]
			r++
		}
	}
	for ; l <= mid; l++ {
		tmp[i] = v[l]
		i++
	}
	for ; r <= right; r++ {
		tmp[i] = v[r]
		i++
	}
	copy(v[left:right+1], tmp)
}

func delayTest(args []byte) ([]byte, error) {
	var argType TestArgs
	err := json.Unmarshal(args, &argType)
	if err != nil {
		return nil, err
	}

	T2 := time.Now().UnixNano()
	ret := TestArgs{
		Res: (T2 - argType.T),
	}

	js, _ := json.Marshal(ret)
	return js, nil
}

func register(args []byte) ([]byte, error) {
	var argType IdentifyArgs
	err := json.Unmarshal(args, &argType)
	if err != nil {
		return nil, err
	}
	_, ok := IDMap.Load(argType.ID)
	if ok {
		ret := IdentifyArgs{
			Status: false,
			Info:   "name already exists",
		}
		js, _ := json.Marshal(ret)
		return js, nil
	} else {
		id := IdGenerate()
		IDMap.Store(argType.Name, id)
		ret := IdentifyArgs{
			Status: true,
			Info:   "success",
			Name:   argType.Name,
			ID:     id,
		}
		js, _ := json.Marshal(ret)
		return js, nil
	}
}

func IdGenerate() string {
	key := time.Now().Format("20060102150405.0000000") + strconv.Itoa(rand.Intn(100000))
	data := []byte(key)
	md := md5.Sum(data)
	return hex.EncodeToString(md[:])
}

func unregister(args []byte) ([]byte, error) {
	var argType IdentifyArgs
	err := json.Unmarshal(args, &argType)
	if err != nil {
		return nil, err
	}
	id, ok := IDMap.Load(argType.Name)
	if ok && id.(string) == argType.ID {
		IDMap.Delete(argType.Name)
		ret := IdentifyArgs{
			Status: true,
			Info:   "success",
		}
		js, _ := json.Marshal(ret)
		return js, nil
	} else {
		ret := IdentifyArgs{
			Status: false,
			Info:   "failed",
		}
		js, _ := json.Marshal(ret)
		return js, nil
	}
}

func getUrl(args []byte) ([]byte, error) {
	var argType UrlArgs
	err := json.Unmarshal(args, &argType)
	if err != nil {
		return nil, err
	}
	ip, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	ret := UrlArgs{
		Url: ip[0].String(),
	}
	js, _ := json.Marshal(ret)
	return js, nil
}
