package main

import (
	"DiscordLogCollector/inputs"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gtuk/discordwebhook"
	"log"
	"net/http"
	"os"
)

func discordSender(id string, color string, token string, title *string, message *string, footer *string) error {
	embed := []discordwebhook.Embed{
		{
			Title:       title,
			Description: message,
			Color:       &color,
		},
	}
	if footer != nil {
		embed[0].Footer = &discordwebhook.Footer{
			Text: footer,
		}
	}

	return discordwebhook.SendMessage(fmt.Sprintf("https://discord.com/api/webhooks/%s/%s", id, token), discordwebhook.Message{
		Embeds: &embed,
	})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/stalwart/{id}/{token}", inputs.StalwartHandler(discordSender)).Methods(http.MethodPost)

	log.Println("Started server on port " + os.Getenv("PORT"))
	if http.ListenAndServe(":"+os.Getenv("PORT"), handlers.LoggingHandler(os.Stdout, handlers.ProxyHeaders(handlers.RecoveryHandler()(r)))) != nil {
		log.Fatal("Error starting server")
	}
}
