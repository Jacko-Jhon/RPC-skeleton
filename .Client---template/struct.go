package Client_tpl

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
