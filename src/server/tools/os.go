package tools

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/former"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

func WriteResponsesToJSONFile(root_dir string, resp former.FormResponse) error {
	storage_dir := root_dir + "/data/" + resp.FormName + "/" + resp.ResponderID + "/"

	json_resp := ConvertFormResponseToJSONFormResponse(root_dir, resp)

	json_bytes, err := json.MarshalIndent(json_resp, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(storage_dir+"responses.json", json_bytes, 0644)
	return err
}

func ConvertFormResponseToJSONFormResponse(root_dir string, resp former.FormResponse) former.JSONFormResponse {
	json_resp := former.JSONFormResponse{}
	json_resp.FormName = resp.FormName
	json_resp.RelationalID = resp.RelationalID
	json_resp.ResponderID = resp.ResponderID
	json_resp.Responses = resp.Responses
	json_resp.FilePaths = make(map[string]string)
	storage_dir := root_dir + "/data/" + resp.FormName + "/" + resp.ResponderID + "/"
	for k, v := range resp.FileObjects {
		json_resp.FilePaths[k] = storage_dir + "files/" + v.Header.Filename
	}

	return json_resp
}

func WriteFilesFromMultipart(root_dir string, response_struct former.FormResponse) []error {
	storage_dir := root_dir + "/data/" + response_struct.FormName + "/" + response_struct.ResponderID + "/files/"
	var err_list []error = []error{}
	for field_name, file_object := range response_struct.FileObjects {
		fname := field_name + "-" + file_object.Header.Filename
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

// borrowing from https://gist.github.com/mimoo/25fc9716e0f1353791f5908f94d6e726
func CreateDownloadableForGivenForm(initialization_folder string, form_name string) error {
	form_dir := initialization_folder + "/data/" + form_name + "/"
	file_path := initialization_folder + "/data/" + form_name + "/downloadable.tar.gz"

	os.Remove(file_path)

	var buf bytes.Buffer
	zr := gzip.NewWriter(&buf)
	tw := tar.NewWriter(zr)

	// walk through every file in the folder
	err := filepath.Walk(form_dir, func(file string, fi os.FileInfo, err error) error {
		// generate tar header
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		// must provide real name
		// (see https://golang.org/src/archive/tar/common.go?#L626)
		header.Name = filepath.ToSlash(file)

		// write header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		// if not a dir, write file content
		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// produce tar
	if err := tw.Close(); err != nil {
		return err
	}
	// produce gzip
	if err := zr.Close(); err != nil {
		return err
	}

	compressed_file, err := os.OpenFile(file_path, os.O_CREATE|os.O_RDWR, os.FileMode(0644))
	if err != nil {
		return err
	}
	if _, err := io.Copy(compressed_file, &buf); err != nil {
		return err
	}
	return nil
}

func WriteJSONReadmeToDir(filename string, field_map map[string]string) error {
	nice_marshal, err := json.MarshalIndent(field_map, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, nice_marshal, 0644)
}

func WriteCSVToDir(filename string, csv_data [][]string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(f)
	return writer.WriteAll(csv_data)
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
