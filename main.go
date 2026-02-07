package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var BaseDelay = 5

type Sender interface {
	Send(string) error
}

func main() {
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
	var delay int = 1
	var fails int = 0
	for {
		err := checker.Check()
		if err != nil {
			fails++
			delay = min((delay + delay), 60)
			log.Printf("Something went wrong %s: new Delay %d\n", checker.GetUrl(), delay)
			if fails >= 3 {
				d.Send(fmt.Sprintf("%-10s\n%d tries\n%-10s\n%s", time.Now().Format(time.RFC3339), fails, checker.GetUrl(), err.Error()))
			}
			time.Sleep(time.Duration(BaseDelay+delay) * time.Minute)
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
