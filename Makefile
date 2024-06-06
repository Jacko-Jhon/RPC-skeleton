server:
	@echo "Server is building"
	cd ./Server-go && go build -o ../server

client:
	@echo "Client is building"
	cd ./Client-go && go build -o ../client

register:
	@echo "Register is building"
	cd ./Register-go && go build -o ../register