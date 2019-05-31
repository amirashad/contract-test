# /bin/sh

go run examples/basic-test/main.go &

sleep 1s
echo "============ Test 1 =============="
go run main.go test.go env.go -t examples/basic-test/health.json
echo "============ Test 2 =============="
go run main.go test.go env.go -t examples/basic-test/get-user.json

curl -s http://localhost:8080/shutdown