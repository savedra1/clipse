package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os/exec"
	"strconv"
	"time"

	"github.com/atotto/clipboard"
)

func main() {

	/*
		in := os.Stdin
		input, err := io.ReadAll(in)
		copyCmd := fmt.Sprintf("wl-copy --type image/png < %s", input)
		cmd := exec.Command("sh", "-c", copyCmd)
		err = cmd.Run()
		if err != nil {
			fmt.Println("error running copy cmd")
		}

		cmd = exec.Command("sh", "-c", "wl-paste --type image/png > new_img.png")
		err = cmd.Run()
		if err != nil {
			fmt.Println("error running paste cmd")
		}

		time.Sleep(1000 * time.Second)*/

	text, err := clipboard.ReadAll()

	fileType := FileType([]byte(text))
	fmt.Println(fileType)
	if fileType == "png" {
		fname := fmt.Sprintf("%s.png", strconv.Itoa(len(text)))
		fmt.Println(fname)

	}

	fmt.Println(len(text))
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println([]byte(text))
	time.Sleep(1000 * time.Second)
	err = exec.Command("sh", "-c", "wl-paste -t image/png > /home/michael/.config/clipboard_manager/tmp_files/image3.png").Run()
	if err != nil {
		fmt.Println(err)
	}

	//time.Sleep(1000 * time.Second)
	/*
	   cmdString := fmt.Sprintf("wl-paste --type image/png > %s", fileName)
	   cmd := os.exec("sh", "-c", cmdString)
	   err = cmd.Run()

	   	if err != nil {
	   		fmt.Println("failed to exec cmd.")
	   	}

	   time.Sleep(1000 * time.Second)
	*/
}

func FileType(data []byte) string {
	reader := bytes.NewReader(data)
	_, err := png.Decode(reader)
	if err == nil {
		return "png"
	}
	_, err = jpeg.Decode(reader)
	if err == nil {
		return "jpg"
	}

	return ""

}
