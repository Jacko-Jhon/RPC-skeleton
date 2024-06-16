server:
	@echo "Server is building"
	cd ./Server-go && go build -o ../server

server-cp:
	@echo "Server is copying"
	rm -f ./Server-go/MyFunctions.go
	cp ./test/MyFunctions.go ./Server-go/MyFunctions.go
	cd ./Server-go && go build -o ../server

client:
	@echo "Client is building"
	cd ./Client-go && go build -o ../client

client-test:
	@echo "Client is testing"
	cp -r ./Client ./clientTest
	cd ./clientTest && go build -o ../test

register:
	@echo "Register is building"
	cd ./Register-go && go build -o ../register