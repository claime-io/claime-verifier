package main

import (
	"claime-verifier/lib/functions/lib/common/log"
	"claime-verifier/lib/functions/lib/guild"
	"claime-verifier/lib/functions/lib/infrastructure/ssm"
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	token string
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	token, err := ssm.New().DiscordBotToken(context.Background())
	if err != nil {
		log.Error("error get bot token", err)
		return
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Error("error creating Discord session,", err)
		return
	}

	dg.AddHandler(guild.GuildMemberAdd)

	dg.Identify.Intents = discordgo.IntentsAll

	err = dg.Open()
	if err != nil {
		log.Error("error opening connection,", err)
		return
	}

	log.Info("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}
