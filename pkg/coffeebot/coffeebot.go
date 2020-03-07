package coffeebot

import (
	"coffeeBot/internal/slackbot"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func Run() {
	cfg, err := readConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGABRT, syscall.SIGTERM)

	r := mux.NewRouter()

	if cfg.Slack.AccessToken != "" {
		log.Println("Starting Slack bot.")
		slack, err := slackbot.New(&cfg.Slack)
		if err != nil {
			log.Fatalln("Failed to start Slack bot:", err)
		}
		r.HandleFunc("/slack", func(w http.ResponseWriter, r *http.Request) {
			slack.HandleCommand(w, r)
		})
	}

	if cfg.Discord.AccessToken != "" {
		log.Println("Discord is not supported (yet).")
	}

	srv := http.Server{
		Addr:              ":" + strconv.Itoa(cfg.Net.Port),
		Handler:           r,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 0,
		WriteTimeout:      15 * time.Second,
	}

	go func() {
		defer close(sig)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	log.Println("Listening for HTTP on port", cfg.Net.Port)

	<-sig
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}
