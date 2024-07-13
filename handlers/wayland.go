package handlers

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/savedra1/clipse/config"
	"github.com/savedra1/clipse/utils"
)

func StoreWLData() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		utils.LogERROR(fmt.Sprintf("failed to read stdin: %s", err))
		return
	}

	dataType = "text"
	img, format, err := image.Decode(bytes.NewReader(input))
	if err == nil {
		dataType = "image"
	}

	switch dataType {
	case "text":
		if string(input) == "" {
			return
		}
		if err := config.AddClipboardItem(string(input), "null"); err != nil {
			utils.LogERROR(fmt.Sprintf("failed to add new item `( %s )` | %s", input, err))
		}
	case "image":
		fileName := fmt.Sprintf("%s-%s.%s", strconv.Itoa(len(input)), utils.GetTimeStamp(), format)
		itemTitle := fmt.Sprintf("%s %s", imgIcon, fileName)
		filePath := filepath.Join(config.ClipseConfig.TempDirPath, fileName)

		if err := config.AddClipboardItem(itemTitle, filePath); err != nil {
			utils.LogERROR(fmt.Sprintf("failed to save image | %s", err))
		}

		out, err := os.Create(filePath)
		if err != nil {
			utils.LogERROR(fmt.Sprintf("failed to store img file: %s", err))
			return
		}

		defer out.Close()

		switch format {
		case "png":
			err = png.Encode(out, img)
		case "jpeg", "jpg":
			err = jpeg.Encode(out, img, nil)
		case "gif":
			err = gif.Encode(out, img, nil)
		default:
			// If format is not recognized, default to PNG
			err = png.Encode(out, img)
		}

		if err != nil {
			utils.LogERROR(fmt.Sprintf("failed to encode img data: %s", err))
			return
		}
	}
}
