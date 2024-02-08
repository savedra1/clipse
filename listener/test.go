package main

import (
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
)

func min() {
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	interrupt := make(chan os.Signal, 1)

	fmt.Println("Press CTRL+/ to trigger success message.")

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyCtrlSlash && char != 100 {
			break
		}
	}
	<-interrupt
	fmt.Println("Exiting...")
}
