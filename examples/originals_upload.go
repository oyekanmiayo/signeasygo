package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	signeasy "signeasygo/v3"
)

func main() {
	accessToken := flag.String("token", "", "Access Token")
	filePath := flag.String("file", "", "File Path")
	name := flag.String("name", "", "Name. Must have extension.")
	rename := flag.String("rename", "true", "Rename if exists")

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open(*filePath)
	defer file.Close()

	base := filepath.Base(*filePath)
	part1, errFile1 := writer.CreateFormFile("file", base)
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println(errFile1)
		return
	}
	_ = writer.WriteField("name", *name)
	_ = writer.WriteField("rename_if_exists", *rename)
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := signeasy.NewClient(http.DefaultClient, *accessToken)
	response, httpResp, err := client.Originals.ImportDocument(
		&signeasy.ImportDocumentBodyParams{
			Payload: payload, MultipartContentType: writer.FormDataContentType(),
		})
	fmt.Println("\n Resp: ", response, "\n httpResp: ", httpResp, "\n Error:", err)
}
