package out

import (
	"log"
	"net/http"
	"net/url"
)

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
