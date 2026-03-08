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

func callApi(req Req, c http.Client, url string) (bool, error) {
	var body io.Reader

	if req.Payload != nil {
		b, err := json.Marshal(req.Payload)
		if err != nil {
			return false, err
		}
		body = bytes.NewBuffer(b)
	}

	fullURL := fmt.Sprintf("%s/%s/%s", url, APIVersion, req.Endpoint)

	r, err := http.NewRequest(strings.ToUpper(req.Method), fullURL, body)
	if err != nil {
		return false, err
	}

	resp, err := c.Do(r)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != req.ExpCode {
		switch resp.StatusCode {
		case 429:
			return true, nil
		default:
			return false, fmt.Errorf(
				"unexpected status code: %d expected: %d\nfor request %s\nplayload: %v\nresponse: %s",
				resp.StatusCode,
				req.ExpCode,
				req.Endpoint,
				req.Payload,
				body,
			)
		}

	}

	log.Printf("✓ created %s", req.Endpoint)
	return false, nil
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

	currentDelay := time.Duration(delay) * time.Second

	for _, request := range requestList {
		// Loop until a request fails or succeeds.
		// Only retry if we hit a ratelimiteed response code.
		for range 3 {
			time.Sleep(currentDelay)

			retry, err := callApi(request, c, url)
			if err != nil {
				return err
			}

			if retry {
				currentDelay += 1 * time.Second
				continue
			}
			break
		}
	}
	return nil
}
