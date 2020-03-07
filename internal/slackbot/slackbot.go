package slackbot

import (
	"coffeeBot/internal/api"
	"encoding/json"
	"fmt"
	"github.com/ajg/form"
	"github.com/google/uuid"
	"github.com/nlopes/slack"
	"log"
	"net/http"
)

func New(cfg *Config) (api.Interface, error) {
	slackApi := slack.New(cfg.AccessToken)
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
	p, err := decodePayload(r)
	if err != nil {
		log.Println("[Slack] Failed to decode payload.", err)
		return
	}

	log.Println("[Slack] Handling", p.Command, p.Text)
	switch p.Command {
	case "/echo":
		if _, err := w.Write([]byte(p.Text)); err != nil {
			log.Println("[Slack] Failed to respond to", p.Command)
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	case "/cof":
		fallthrough
	case "/covfefe":
		fallthrough
	case "/coffee":
		handleCoffeeCommand(w, &p)
	default:
		log.Println("[Slack] Unrecognized command:", p.Command)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func decodePayload(r *http.Request) (p payload, err error) {
	d := form.NewDecoder(r.Body)
	d.IgnoreUnknownKeys(true)
	err = d.Decode(&p)
	return
}

func handleCoffeeCommand(w http.ResponseWriter, p *payload) {
	w.Header().Set("Content-Type", "application/json")

	milkType := api.MilkTypeDairy

	id := uuid.New()

	heading := slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf(":coffee: *@%s is making coffee in %d minutes.*", p.UserName, 5), false, false)

	spaces := slack.NewTextBlockObject(slack.PlainTextType, fmt.Sprintf("There is space for %d more people.", 2), false, false)

	joinButton := slack.NewButtonBlockElement("join", id.String(), slack.NewTextBlockObject(slack.PlainTextType, milkType, true, false))

	joiners := []slack.MixedElement{
		slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/profile_1.png", "Michael Scott"),
	}

	blockMessage := slack.NewBlockMessage(
		slack.NewSectionBlock(heading, nil, nil),
		slack.NewDividerBlock(),
		slack.NewSectionBlock(spaces, nil, slack.NewAccessory(joinButton)),
		slack.NewContextBlock(id.String(), joiners...),
	)

	if err := json.NewEncoder(w).Encode(blockMessage); err != nil {
		log.Println("[Slack] Failed to create block response.", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
