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
	"github.com/savedra1/clipse/shell"
	"github.com/savedra1/clipse/utils"
)

func StoreWLData() {
	/* See `man wl-clipboard` for more information */
	if os.Getenv("CLIPBOARD_STATE") == "sensitive" {
		return
	}

	// Check if the clipboard content should be excluded based on source application
	activeWindow := utils.GetActiveWindowTitle()
	if utils.IsAppExcluded(activeWindow, config.ClipseConfig.ExcludedApps) {
		utils.LogINFO(fmt.Sprintf("Skipping clipboard content from excluded app: %s", activeWindow))
		return
	}

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		utils.LogERROR(fmt.Sprintf("failed to read stdin: %s", err))
		return
	}

	dt := Text
	if len(input) > 0 && input[0] == 0x89 && string(input[1:4]) == "PNG" {
		dt = PNG
	} else if len(input) > 10 && string(input[6:10]) == "JFIF" {
		dt = JPEG
	}

	switch dt {
	case Text:
		inputStr := string(input)
		if inputStr == "" {
			return
		}
		if err := config.AddClipboardItem(inputStr, "null"); err != nil {
			utils.LogERROR(fmt.Sprintf("failed to add new item `( %s )` | %s", input, err))
		}

	case PNG, JPEG:
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

		fileName := fmt.Sprintf("%s.%s", utils.GetTimeStamp(), "png")
		filePath := filepath.Join(config.ClipseConfig.TempDirPath, fileName)

		img, format, err := image.Decode(bytes.NewReader(input))

		if err != nil {
			/*
				if the image data cannot be decoded here it means this has
				not loaded properly from the wl-paste --watch api.
				the image is then created using `wl-paste -t image/png <path>`
			*/

			if err = shell.SaveImage(filePath, "wayland"); err != nil {
				utils.LogERROR(fmt.Sprintf("failed to save new image: %s", err))
				return
			}

			updatedFileName, updatedFilePath, err := renameImgFile(filePath, fileName, dt)
			if err != nil {
				utils.LogERROR(fmt.Sprintf("failed to rename new image file: %s", err))
				return
			}

			itemTitle := fmt.Sprintf("%s %s", imgIcon, updatedFileName)

			if err := config.AddClipboardItem(itemTitle, updatedFilePath); err != nil {
				utils.LogERROR(fmt.Sprintf("failed to save image | %s", err))
			}

			return
		}

		out, err := os.Create(filePath)
		if err != nil {
			utils.LogERROR(fmt.Sprintf("failed to create img file: %s", err))
			return
		}

		defer out.Close()

		switch format {
		case PNG:
			err = png.Encode(out, img)
		case JPEG, JPG:
			err = jpeg.Encode(out, img, nil)
		default:
			// default to PNG if format not recognized
			err = png.Encode(out, img)
		}

		if err != nil {
			utils.LogERROR(fmt.Sprintf("failed to encode img data: %s", err))
			return
		}

		updatedFileName, updatedFilePath, err := renameImgFile(filePath, fileName, dt)

		if err != nil {
			utils.LogERROR(fmt.Sprintf("failed to rename new image file: %s", err))
			return
		}

		itemTitle := fmt.Sprintf("%s %s", imgIcon, updatedFileName)

		if err := config.AddClipboardItem(itemTitle, updatedFilePath); err != nil {
			utils.LogERROR(fmt.Sprintf("failed to save image | %s", err))
		}
	}
}

func renameImgFile(filePath, fileName, dt string) (string, string, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", "", err
	}

	fileSize := fileInfo.Size()
	updatedFileName := fmt.Sprintf(
		"%s-%s.%s",
		strconv.Itoa(int(fileSize)),
		strings.Split(fileName, ".")[0],
		dt,
	)
	updatedFilePath := filepath.Join(config.ClipseConfig.TempDirPath, updatedFileName)

	if err = os.Rename(filePath, updatedFilePath); err != nil {
		return "", "", err
	}

	return updatedFileName, updatedFilePath, nil
}
