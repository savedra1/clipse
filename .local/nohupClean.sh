ps aux | grep main.go
pkill -f "go run main.go"

nohup go run main.go --start >/dev/null 2>&1 &