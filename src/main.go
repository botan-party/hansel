package main

import (
	"hansel/discord"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	bot, err := discord.NewBot()
	if err != nil {
		log.Fatalln(err)
		return
	}

	err = bot.Start()
	if err != nil {
		log.Fatalln(err)
		return
	}

	defer bot.Stop()

	// 終了を待機
	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		os.Interrupt,
		os.Kill,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	select {
	case <-signalChan:
		log.Println("bye")
		return
	}
}
