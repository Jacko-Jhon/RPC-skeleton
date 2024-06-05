package main

import "sync"

var IDMap sync.Map

type AddArgs struct {
	A   float32 `json:"A"`
	B   float32 `json:"B"`
	Res float32 `json:"Res"`
}

type SubArgs struct {
	A   float32 `json:"A"`
	B   float32 `json:"B"`
	Res float32 `json:"Res"`
}

type MulArgs struct {
	A   float32 `json:"A"`
	B   float32 `json:"B"`
	Res float32 `json:"Res"`
}

type DivArgs struct {
	A   float32 `json:"A"`
	B   float32 `json:"B"`
	Res float32 `json:"Res"`
}

type SortArgs struct {
	Vector []float32 `json:"Vector"`
	Res    []float32 `json:"Res"`
}

type TestArgs struct {
	T   int64 `json:"T"`
	Res int64 `json:"Res"`
}

type IdentifyArgs struct {
	Status bool   `json:"Status"`
	Info   string `json:"Info"`
	Name   string `json:"Name"`
	ID     string `json:"ID"`
}

type UrlArgs struct {
	Get bool   `json:"Get"`
	Url string `json:"Url"`
}
