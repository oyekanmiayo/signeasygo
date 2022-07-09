package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	signeasy "signeasygo/v2"
)

func main() {
	accessToken := flag.String("token", "", "Access Token")
	fmt.Println("Args: ", os.Args)
	flag.Parse()

	client := signeasy.NewClient(http.DefaultClient, *accessToken)
	response, httpResp, err := client.Originals.ListOriginals()
	fmt.Println("\n Resp: ", response, "\n httpResp: ", httpResp, "\n Error:", err)
}
