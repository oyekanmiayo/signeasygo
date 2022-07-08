package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	signeasy "signeasygo/v2"
)

func main() {
	flags := flag.NewFlagSet("flags", flag.ExitOnError)
	accessToken := flags.String("access-token", "", "Access Token")
	filePath := flags.String("file-path", "", "File Path")
	name := flags.String("name", "", "Name. Must have extension.")
	rename := flags.String("rename", "true", "Rename if exists")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("error: %v", err)
	}

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
	err = writer.Close()
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
