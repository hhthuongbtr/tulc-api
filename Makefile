build:
	GOOS=linux GOARCH=amd64 go build -o deployment/tulc-api main.go
	docker build -t tulc-api deployment
clean:
	rm -f deployment/tulc-api
run:
	docker run -p 8088:8088 -itd -e GIN_MODE=release tulc-api