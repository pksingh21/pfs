package main

import (
	"fmt"
	"os"
	"os/user"
	"syscall"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type File struct {
	isDirectory bool
	isFile      bool
	Name        string
	Owner       string
	Group       string
}

func GetFilesAndDirectories(path string) []File {
	var filesAndDirectories []File
	fileInfos, err := os.ReadDir(path)
	if err != nil {
		return nil
	}

	for _, fileInfo := range fileInfos {
		if fileInfo.Name() == "." || fileInfo.Name() == ".." {
			continue
		}
		extraFileInfo, _ := os.Stat(path + "/" + fileInfo.Name())
		uuid := extraFileInfo.Sys().(*syscall.Stat_t).Uid
		uuid_string := fmt.Sprint(uuid)
		userName, _ := user.LookupId(uuid_string)
		groupName, _ := user.LookupGroupId(fmt.Sprint(extraFileInfo.Sys().(*syscall.Stat_t).Gid))
		file := File{
			isDirectory: fileInfo.IsDir(),
			isFile:      !fileInfo.IsDir(),
			Name:        fileInfo.Name(),
			Owner:       userName.Name,
			Group:       groupName.Name,
		}
		filesAndDirectories = append(filesAndDirectories, file)
	}

	return filesAndDirectories
}
func main() {
	filesAndDirectories := GetFilesAndDirectories("/home/pks")
	for _, file := range filesAndDirectories {
		fmt.Println(file.isDirectory, file.isFile, file.Name, file.Owner,file.Group)
	}
	a := app.New()
	w := a.NewWindow("Hello")

	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))

	w.ShowAndRun()
}
