package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const cep = "17018600"

type Response struct {
	Source   string
	Response string
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	url1 := "https://brasilapi.com.br/api/cep/v1/" + cep
	url2 := "http://viacep.com.br/ws/" + cep + "/json/"

	resChan := make(chan Response)
	errChan := make(chan error)
	go FetchData(ctx, url1, resChan, errChan)
	go FetchData(ctx, url2, resChan, errChan)

	select {
	case <-time.After(1 * time.Second):
		log.Println("Request timeout")
		return
	case err := <-errChan:
		fmt.Println("error received: " + err.Error())
		return
	case res := <-resChan:
		fmt.Println("Response received:")
		fmt.Println("Source: " + res.Source)
		fmt.Println("Content: " + res.Response)
	}
}

func FetchData(ctx context.Context, url string, resChan chan<- Response, errChan chan<- error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		errChan <- err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		errChan <- err
	}
	defer res.Body.Close()

	respValue, err := io.ReadAll(res.Body)
	if err != nil {
		errChan <- err
	}
	resChan <- Response{
		Source:   url,
		Response: string(respValue),
	}
}
