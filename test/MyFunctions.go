package main

type Runnable func([]byte) ([]byte, error)

type Function struct {
	Name string
	Id   string // if you want to register by name, you should set id that already exists
	// else you can set id to "" or don't set id
	Args       []string
	Ret        []string
	Timeout    int
	MaxProcess int32
	run        Runnable
}

var MyFunctions = []Function{
	{
		Name:       "add",
		Id:         "6fc6da58961ad9230b946a4853767d46",
		Args:       []string{"A float32", "B float32"},
		Ret:        []string{"Res float32"},
		Timeout:    1,
		MaxProcess: 1,
		run:        add,
	},
	{
		Name:       "sub",
		Id:         "99eb104e4bbbe798addf57244d453322",
		Args:       []string{"A float32", "B float32"},
		Ret:        []string{"Res float32"},
		Timeout:    1,
		MaxProcess: 1,
		run:        sub,
	},
	{
		Name:       "mul",
		Id:         "d47a84634898d3714c559627da06d36a",
		Args:       []string{"A float32", "B float32"},
		Ret:        []string{"Res float32"},
		Timeout:    1,
		MaxProcess: 1,
		run:        mul,
	},
	{
		Name:       "div",
		Id:         "31e616f5540db0310ef8c5581b019a4b",
		Args:       []string{"A float32", "B float32"},
		Ret:        []string{"Res float32"},
		Timeout:    1,
		MaxProcess: 1,
		run:        div,
	},
	{
		Name:       "MergeSort",
		Id:         "23a3cccce11d2d9330ebf50d0799d533",
		Args:       []string{"Vector []float32"},
		Ret:        []string{"Res []float32"},
		Timeout:    1000,
		MaxProcess: 1,
		run:        mergeSort,
	},
	{
		Name:       "QSort",
		Id:         "0d722e613e1e280d3c83bb4c542db121",
		Args:       []string{"Vector []float32"},
		Ret:        []string{"Res []float32"},
		Timeout:    1000,
		MaxProcess: 1,
		run:        qSort,
	},
	{
		Name:       "DelayTest",
		Id:         "bb764f7bf83affae61b7bdd78821b1c5",
		Args:       []string{"T int64"},
		Ret:        []string{"Res int64"},
		Timeout:    1,
		MaxProcess: 1,
		run:        delayTest,
	}, {
		Name:       "register",
		Args:       []string{"Name string", "Id string"},
		Ret:        []string{"Status bool", "Info string"},
		Timeout:    0,
		MaxProcess: 1,
		run:        register,
	},
	{
		Name:       "unregister",
		Args:       []string{"Name string", "Id string"},
		Ret:        []string{"Status bool", "Info string"},
		Timeout:    0,
		MaxProcess: 1,
		run:        unregister,
	}, {
		Name:       "GetUrl",
		Args:       []string{"Get bool"},
		Ret:        []string{"Url string"},
		Timeout:    0,
		MaxProcess: 1,
		run:        getUrl,
	},
}
