package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/http2"
)

func main() {
	// Enable HTTP/2
	http2.ConfigureServer(nil, &http2.Server{})

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Request failed with status", resp.Status)
		return
	}

	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}
}