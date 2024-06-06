server:
	@echo "Server is building"
	go build -o ./server ./Server-go

client:
	@echo "Client is building"
	go build -o ./client ./Client-go

register:
	@echo "Register is building"
	go build -o ./register ./Register-go