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
		Args:       []string{"A float32", "B float32"},
		Ret:        []string{"Res float32"},
		Timeout:    1,
		MaxProcess: 1,
		run:        add,
	},
	{
		Name:       "sub",
		Args:       []string{"A float32", "B float32"},
		Ret:        []string{"Res float32"},
		Timeout:    1,
		MaxProcess: 1,
		run:        sub,
	},
	{
		Name:       "mul",
		Args:       []string{"A float32", "B float32"},
		Ret:        []string{"Res float32"},
		Timeout:    1,
		MaxProcess: 1,
		run:        mul,
	},
	{
		Name:       "div",
		Args:       []string{"A float32", "B float32"},
		Ret:        []string{"Res float32"},
		Timeout:    1,
		MaxProcess: 1,
		run:        div,
	},
	{
		Name:       "MergeSort",
		Args:       []string{"Vector []float32"},
		Ret:        []string{"Res []float32"},
		Timeout:    1000,
		MaxProcess: 1,
		run:        mergeSort,
	},
	{
		Name:       "QSort",
		Args:       []string{"Vector []float32"},
		Ret:        []string{"Res []float32"},
		Timeout:    1000,
		MaxProcess: 1,
		run:        qSort,
	},
	{
		Name:       "DelayTest",
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
