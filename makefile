
run-admin:
	go run app/tooling/sales-admin/main.go | go run app/tooling/logfmt/main.go
	
run:
	go run app/services/sales-api/main.go | go run app/tooling/logfmt/main.go

help:
	go run app/services/sales-api/main.go -h

version:
	go run app/services/sales-api/main.go -v

tidy:
	go mod tidy
	go mod vendor