package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const APIVersion = "v1"

type Req struct {
	Endpoint string
	Method   string
	Payload  map[string]any
	ExpCode  int
}

func callApi(req Req, c http.Client, delay int, url string) error {
	var body io.Reader

	if req.Payload != nil {
		b, err := json.Marshal(req.Payload)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(b)
	}

	fullURL := fmt.Sprintf("%s/%s/%s", url, APIVersion, req.Endpoint)

	r, err := http.NewRequest(strings.ToUpper(req.Method), fullURL, body)
	if err != nil {
		return err
	}

	time.Sleep(time.Duration(delay) * time.Second)
	resp, err := c.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != req.ExpCode {
		return fmt.Errorf(
			"unexpected status code: %d expected: %d\nfor request %s\nplayload: %v",
			resp.StatusCode,
			req.ExpCode,
			req.Endpoint,
			req.Payload,
		)
	}

	log.Printf("✓ created %s", req.Endpoint)
	return nil
}

func main() {
	log.SetFlags(0)
	delay := flag.Int("d", 0, "enable a delay between requests (seconds)")
	url := flag.String("u", "http://localhost:3000", "api url (e.g., http://localhost:3000)")
	flag.Parse()

	if err := Run(*delay, *url); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func Run(delay int, url string) error {
	c := http.Client{
		Timeout: 3 * time.Second,
	}

	for _, request := range requestList {
		err := callApi(request, c, delay, url)
		if err != nil {
			return err
		}

	}
	return nil
}
