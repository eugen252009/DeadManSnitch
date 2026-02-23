package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	BaseDelay       = 5
	ReportCount int = 3
)

type Sender interface {
	Send(string) error
}

func main() {
	count := os.Getenv("REPORTCOUNT")
	if n, err := strconv.Atoi(count); err != nil {
		log.Printf("%s is not a Number!, using %d", count, ReportCount)
	} else {
		ReportCount = n
	}
	hosts := strings.Split(os.Getenv("HOSTS"), ",")
	if len(hosts) == 0 {
		panic("No Hosts Found!\nSpecify Env HOSTS")
	}
	d := DiscordResult{URL: os.Getenv("DISCORD_WEBHOOK")}
	for _, url := range hosts {
		switch true {
		case strings.HasPrefix(url, "tcp://"):
			go Heartbeat(&d, &TCP{Url: url[6:]})
		case strings.HasPrefix(url, "https://"):
			go Heartbeat(&d, &HTTP{Url: url})
		}
	}
	select {}
}

type Checker interface {
	Check() error
	GetUrl() string
}

func Heartbeat(d Sender, checker Checker) {
	delay := 1
	fails := 0
	for {
		err := checker.Check()
		if err != nil {
			msg := err.Error()
			if errors.Is(err, context.DeadlineExceeded) {
				msg = "Error: Could not reach Host"
			}
			fails++
			log.Printf("Something went wrong %s: tries %d\n", checker.GetUrl(), fails)
			if fails >= ReportCount {
				d.Send(fmt.Sprintf("%-10s\n%-10.0d tries\n%-10s\n%s", time.Now().Format(time.RFC3339), fails, checker.GetUrl(), msg))
			}
			time.Sleep(time.Duration(BaseDelay+delay) * time.Minute)
			delay = min((delay + delay), 60)
			continue
		}
		delay = 1
		time.Sleep(time.Duration(BaseDelay+delay) * time.Minute)
	}
}

type (
	HTTP struct {
		Url string
	}
	TCP struct {
		Url string
	}
)

func (t *TCP) Check() error {
	conn, err := net.DialTimeout("tcp", t.Url, 5*time.Second)
	if err != nil {
		return fmt.Errorf("%s TCP Port geschlossen", t.Url)
	}
	defer conn.Close()
	return nil
}

func (t *TCP) GetUrl() string {
	return t.Url
}

func (h *HTTP) GetUrl() string {
	return h.Url
}

func (h *HTTP) Check() error {
	Client := &http.Client{
		Timeout: 5 * time.Second,
	}
	request, err := Client.Head(h.Url)
	if err != nil {
		return err
	}
	request.Body.Close()
	if request.StatusCode >= 300 {
		return fmt.Errorf("%-25s %-20s %s", time.Now().Format(time.RFC3339), h.Url, request.Status)
	}
	return nil
}

type DiscordResult struct {
	URL string
}

func (d *DiscordResult) Send(content string) error {
	data := url.Values{}
	data.Set("content", content)
	resp, err := http.PostForm(d.URL, data)
	if err != nil {
		log.Println("Webhook down!")
		log.Printf("%s\n", content)
		return err
	}
	defer resp.Body.Close()
	return nil
}
