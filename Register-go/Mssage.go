package Register_go

type MessageToServer struct {
	status bool
	id     string
	info   string
}

type MessageToClient struct {
	status bool
	id     string
	info   string
	json   []byte
}
