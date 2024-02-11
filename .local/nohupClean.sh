ps aux | grep main.go
# kill <pid>

nohup go run main.go --start >/dev/null 2>&1 &