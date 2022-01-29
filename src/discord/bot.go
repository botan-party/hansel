package discord

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session *discordgo.Session
}

func NewBot() (*Bot, error) {
	session, err := discordgo.New("Bot " + os.Getenv("BOT_ID"))
	if err != nil {
		return &Bot{}, nil
	}

	return &Bot{
		session: session,
	}, nil
}

func (b *Bot) Start() error {
	b.session.AddHandler(b.receive)

	err := b.session.Open()
	if err != nil {
		log.Println("Failed : Start Bot")
		return err
	}

	log.Println("Succeeded : Start Bot")
	return nil
}

func (b *Bot) Stop() error {
	return b.session.Close()
}
