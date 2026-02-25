package integration

import (
	"fmt"
	"net/http"
	"time"
)

func (h *HTTP) GetURL() string {
	return h.URL
}

func (h *HTTP) Check() error {
	Client := &http.Client{
		Timeout: 5 * time.Second,
	}
	request, err := Client.Head(h.URL)
	if err != nil {
		return err
	}
	request.Body.Close()
	if request.StatusCode >= 300 {
		return fmt.Errorf("%-25s %-20s %s", time.Now().Format(time.RFC3339), h.URL, request.Status)
	}
	return nil
}
