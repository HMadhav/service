

run:
	go run app/services/sales-api/main.go

help:
	go run app/services/sales-api/main.go -h

version:
	go run app/services/sales-api/main.go -v

tidy:
	go mod tidy
	go mod vendor