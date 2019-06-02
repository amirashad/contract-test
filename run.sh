# /bin/sh

go run examples/basic-test/main.go &

sleep 1s
go run main.go test.go env.go -t examples/basic-test/health.json
go run main.go test.go env.go -t examples/basic-test/get-user.json
go run main.go test.go env.go -t examples/basic-test/get-users.json
go run main.go test.go env.go -t examples/basic-test/get-states.json

curl -s http://localhost:8080/shutdown