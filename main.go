package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	i "github.com/eugen252009/deadmansnitch/bin/integration"
	"github.com/eugen252009/deadmansnitch/bin/out"
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
		log.Printf("REPORTCOUNT: %s is not a Number!, using %d", count, ReportCount)
	} else {
		ReportCount = n
	}
	hosts := strings.Split(os.Getenv("HOSTS"), ",")
	if len(hosts) == 0 {
		panic("No Hosts Found!\nSpecify Env HOSTS")
	}
	d := out.DiscordResult{URL: os.Getenv("DISCORD_WEBHOOK")}
	for _, url := range hosts {
		switch true {
		case strings.HasPrefix(url, "tcp://"):
			go Heartbeat(&d, &i.TCP{URL: url[6:]})
		case strings.HasPrefix(url, "https://"):
			go Heartbeat(&d, &i.HTTP{URL: url})
		}
	}
	select {}
}

func Heartbeat(d Sender, checker i.Checker) {
	delay := 0
	fails := 0
	for {
		err := checker.Check()
		if err != nil {
			fails++
			msg := err.Error()
			if errors.Is(err, context.DeadlineExceeded) {
				msg = "Error: Could not reach Host"
			}
			if fails >= ReportCount {
				err := d.Send(fmt.Sprintf("%-10s\n%-10.0d tries\n%-10s\n%s", time.Now().Format(time.RFC3339), fails, checker.GetURL(), msg))
				if err != nil {
					log.Printf("Webhook Error: %s\n", err)
					log.Printf("Something went wrong %s: tries %d\n", checker.GetURL(), fails)
				}
			}
			delay = min(delay+BaseDelay, 60)
			time.Sleep(time.Duration(delay) * time.Minute)
			continue
		}
		fails = 0
		delay = BaseDelay
		time.Sleep(time.Duration(delay) * time.Minute)
	}
}
