package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	//t := time.Now().UTC().String()
	T := strings.Split(time.Now().UTC().String(), "+0000")[0]

	fmt.Println(T)
}
