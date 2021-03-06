# /bin/sh

go run examples/basic-test/main.go &

go run main.go test.go env.go wait.go -t examples/basic-test/health.json --wait-for-endpoint "http://localhost:8080/health"
go run main.go test.go env.go wait.go -t examples/basic-test/get-user.json
go run main.go test.go env.go wait.go -t examples/basic-test/get-users.json
go run main.go test.go env.go wait.go -t examples/basic-test/get-states.json
go run main.go test.go env.go wait.go -t examples/basic-test/create-user.json
go run main.go test.go env.go wait.go -t examples/basic-test/create-state.json

curl -s http://localhost:8080/shutdown