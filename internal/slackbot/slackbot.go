package slackbot

import (
	"coffeeBot/internal/api"
	"github.com/ajg/form"
	"github.com/nlopes/slack"
	"log"
	"net/http"
)

func New(cfg *Config) (api.Interface, error) {
	slackApi := slack.New(cfg.AccessToken, slack.OptionDebug(true))
	if r, err := slackApi.AuthTest(); err != nil {
		return nil, err
	} else {
		log.Println("[Slack] Authenticated for team", r.Team, "as user", r.User)
	}
	if err := slackApi.SetUserAsActive(); err != nil {
		return nil, err
	}
	return &bot{client: slackApi}, nil
}

func (b *bot) HandleCommand(w http.ResponseWriter, r *http.Request) {
	p, err := decodePayload(w, r)
	if err != nil {
		log.Println("[Slack] Failed to decode payload.", err)
		return
	}

	log.Println("[Slack] Handling", p.Command, p.Text)
	switch p.Command {
	case "/echo":
		if _, err := w.Write([]byte(p.Text)); err != nil {
			log.Println("Failed to respond to", p.Command)
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func decodePayload(w http.ResponseWriter, r *http.Request) (p payload, err error) {
	d := form.NewDecoder(r.Body)
	d.IgnoreUnknownKeys(true)
	err = d.Decode(&p)
	return
}
