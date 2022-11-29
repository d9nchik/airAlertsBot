package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"airAlertsBot/pkg/core/bot"
	"airAlertsBot/pkg/core/puller"
)

func main() {
	token := os.Getenv("token")
	chatName := os.Getenv("chatName")

	b := bot.NewBot(token, chatName)
	p := puller.NewPuller(b)

	quitCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	p.Run(quitCtx)
}
