package main

import (
	"encoding/json"
)

type Runnable func([]byte) ([]byte, error)

type Function struct {
	Name string
	Id   string // if you want to register by name, you should set id that already exists
	// else you can set id to "" or don't set id
	Args []string
	Ret  []string
	run  Runnable
}

var MyFunctions = []Function{
	{
		Name: "add",
		Args: []string{"A int", "B int"},
		Ret:  []string{"Res int"},
		run:  add,
	},
	{
		Name: "sub",
		Args: []string{"A float32", "B float32"},
		Ret:  []string{"Res float32"},
		run:  sub,
	},
}

type AddType struct {
	A   int `json:"A"`
	B   int `json:"b"`
	Res int `json:"res"`
}

func add(args []byte) ([]byte, error) {
	var argType AddType
	err := json.Unmarshal(args, &argType)
	if err != nil {
		return nil, err
	}

	argType.Res = argType.A + argType.B

	js, _ := json.Marshal(argType)
	return js, nil
}

type SubType struct {
	A   float32 `json:"a"`
	B   float32 `json:"b"`
	Res float32 `json:"res"`
}

func sub(args []byte) ([]byte, error) {
	var argType SubType
	err := json.Unmarshal(args, &argType)
	if err != nil {
		return nil, err
	}

	argType.Res = argType.A - argType.B

	js, _ := json.Marshal(argType)
	return js, nil
}
