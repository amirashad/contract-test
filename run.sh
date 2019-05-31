# /bin/sh

go run examples/basic-test/main.go &
sleep 2s
go run main.go test.go env.go
curl http://localhost:8080/shutdown