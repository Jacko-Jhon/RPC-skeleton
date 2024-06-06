package Client
type serviceUrls struct {
	Status  bool     `json:"status"`
	Info    string   `json:"info"`
	Name    string   `json:"name"`
	Ips     []string `json:"ips"`
	Ports   []int    `json:"ports"`
	Factors []int32  `json:"factors"`
	Args    []string `json:"args"`
	Ret     []string `json:"ret"`
}

type urls struct {
	s       int32
	ips     []string
	ports   []int
	factors []int32
}


type MulArgs struct {
	A	float32	`json:"A"`
	B	float32	`json:"B"`
	Res	float32	`json:"Res"`
}
type SubArgs struct {
	A	float32	`json:"A"`
	B	float32	`json:"B"`
	Res	float32	`json:"Res"`
}
type UnregisterArgs struct {
	Name	string	`json:"Name"`
	Id	string	`json:"Id"`
	Status	bool	`json:"Status"`
	Info	string	`json:"Info"`
}
type DelayTestArgs struct {
	T	int64	`json:"T"`
	Res	int64	`json:"Res"`
}
type GetUrlArgs struct {
	Get	bool	`json:"Get"`
	Url	string	`json:"Url"`
}
type AddArgs struct {
	A	float32	`json:"A"`
	B	float32	`json:"B"`
	Res	float32	`json:"Res"`
}
type DivArgs struct {
	A	float32	`json:"A"`
	B	float32	`json:"B"`
	Res	float32	`json:"Res"`
}
type MergeSortArgs struct {
	Vector	[]float32	`json:"Vector"`
	Res	[]float32	`json:"Res"`
}
type QSortArgs struct {
	Vector	[]float32	`json:"Vector"`
	Res	[]float32	`json:"Res"`
}
type RegisterArgs struct {
	Name	string	`json:"Name"`
	Id	string	`json:"Id"`
	Status	bool	`json:"Status"`
	Info	string	`json:"Info"`
}