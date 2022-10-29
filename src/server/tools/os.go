package tools

import (
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/former"
	"path/filepath"
	"strings"
  "errors"
	"bufio"
	"fmt"
	"io"
	"os"
)

func LogError(storage_dir string, message string) {
	err_handler, err := os.OpenFile(storage_dir+"errors.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer err_handler.Close()
	_, err = err_handler.WriteString("File write fail: " + message + "\n")
	if err != nil {
		panic(err)
	}
}

func CheckSafeDirectoryName(dir string) bool {
	bad := strings.Contains(dir, "/")
	if bad {
		return false
	}
	bad = strings.Contains(dir, "\\")
	if bad {
		return false
	}
	return true
}


func WriteFilesFromMultipart(root_dir string, response_struct former.FormResponse) []error {
	storage_dir := root_dir + "/data/" + response_struct.FormName + "/" + response_struct.ResponderID + "/files/"
	var err_list []error = []error{}
	for _, file_object := range response_struct.FileObjects {
		fname := file_object.Header.Filename
		if len(fname) > 255 {
			ext := filepath.Ext(fname)
			base := filepath.Base(fname)
			fname = base[:255 - len(ext)] + ext
		}
		if strings.Contains(fname, "/") {
			LogError(storage_dir, storage_dir+fname)
			err_list = append(err_list, errors.New("File "+fname+" contained illegal characters"))
			continue
		}
		handler, err := os.OpenFile(storage_dir+fname, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			LogError(storage_dir, storage_dir+fname)
			err_list = append(err_list, err)
			continue
		}
		defer handler.Close()
		_, err = io.Copy(handler, file_object.File)
		if err != nil {
			LogError(storage_dir, storage_dir+fname)
			err_list = append(err_list, err)
			continue
		}
	}
	return err_list
}

func ReadLine(item string) string {
	fmt.Print("  (" + item + ") --> ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	// convert CRLF to LF
	text = strings.Replace(text, "\n", "", -1)
	return text
}

func ArrayReverse[T any](arr []T) []T {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}
