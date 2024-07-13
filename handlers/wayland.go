package handlers

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/savedra1/clipse/config"
	"github.com/savedra1/clipse/utils"
)

func StoreWLData() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		utils.LogERROR(fmt.Sprintf("failed to read stdin: %s", err))
		return
	}

	dataType = "text" // defined in defaultHandler
	img, format, err := image.Decode(bytes.NewReader(input))
	if err == nil {
		dataType = "image"
	}

	switch dataType {
	case "text":
		inputStr := string(input)
		if inputStr == "" {
			return
		}
		if err := config.AddClipboardItem(inputStr, "null"); err != nil {
			utils.LogERROR(fmt.Sprintf("failed to add new item `( %s )` | %s", input, err))
		}

	case "image":
		/*
			When saving image data from the stdin using wl-paste --watch,
			the byte size if different to when the image data is copied
			from the saved file with wl-copy -t image/path.
			This means to maintain consistency with identifying duplicates
			we need to save a temporary image file, then update the file name
			using the file size as the identifier, so any duplicates can be
			auto-removed during the AddClipboardItem call in the same way
			as non-wayland specific data.
		*/
		fileName := fmt.Sprintf("%s.%s", utils.GetTimeStamp(), format)
		filePath := filepath.Join(config.ClipseConfig.TempDirPath, fileName)

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
		default:
			// default to PNG if format not recognized
			err = png.Encode(out, img)
		}

		if err != nil {
			utils.LogERROR(fmt.Sprintf("failed to encode img data: %s", err))
			return
		}

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			utils.LogERROR(fmt.Sprintf("failed to read new image file: %s", err))
			return
		}

		fileSize := fileInfo.Size()
		updatedFileName := fmt.Sprintf(
			"%s-%s.%s",
			strconv.Itoa(int(fileSize)),
			strings.Split(fileName, ".")[0],
			format,
		)
		updatedFilePath := filepath.Join(config.ClipseConfig.TempDirPath, updatedFileName)

		if err = os.Rename(filePath, updatedFilePath); err != nil {
			utils.LogERROR(fmt.Sprintf("failed to rename new image file: %s", err))
			return
		}

		itemTitle := fmt.Sprintf("%s %s", imgIcon, updatedFileName)

		if err := config.AddClipboardItem(itemTitle, updatedFilePath); err != nil {
			utils.LogERROR(fmt.Sprintf("failed to save image | %s", err))
		}

	}
}
