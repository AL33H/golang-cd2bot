package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
)

var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.Contains(m.Content, "!mr") {
		ambiente, url, projeto, numeroMr, branch, atividadeId, formatModification, descricao := extractInfoMessage(m)
		embeded := generateEmbed(m, url, ambiente, numeroMr, projeto, formatModification, atividadeId, descricao, branch)
		s.ChannelMessageSendEmbed(m.ChannelID, embeded)
	}

	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

func extractInfoMessage(m *discordgo.MessageCreate) (string, string, string, string, string, string, string, string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Comando inválido!")
		}
	}()

	splitedContent := strings.Split(m.Content, " ")
	urlMr := strings.Split(splitedContent[2], "/")

	ambiente := splitedContent[1]

	url := splitedContent[2]
	projeto := urlMr[5]
	numeroMr := urlMr[8]

	urlBranch := strings.Split(splitedContent[3], "/")
	branch := splitedContent[3]
	atividadeId := urlBranch[1]
	formatModification := urlBranch[0]
	descricao := urlBranch[2]

	return ambiente, url, projeto, numeroMr, branch, atividadeId, formatModification, descricao
}

func generateEmbed(m *discordgo.MessageCreate, url string, ambiente string, numeroMr string, projeto string, formatModification string, atividadeId string, descricao string, branch string) *discordgo.MessageEmbed {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Comando inválido!")
		}
	}()

	embeded := embed.NewGenericEmbed("Solicitação de MR ", m.Author.Username+" ~ "+time.Now().Format("01/02/2006"))
	embeded.Color = 0x00ff00
	embeded.URL = url
	embeded.Fields = append(embeded.Fields,
		&discordgo.MessageEmbedField{Name: "AMBIENTE", Value: ambiente, Inline: true},
		&discordgo.MessageEmbedField{Name: "MR", Value: numeroMr, Inline: true},
		&discordgo.MessageEmbedField{Name: "PROJETO", Value: projeto, Inline: true},
		&discordgo.MessageEmbedField{Name: "TIPO", Value: formatModification, Inline: true},
		&discordgo.MessageEmbedField{Name: "ATIVIDADE", Value: atividadeId, Inline: true},
		&discordgo.MessageEmbedField{Name: "DESCRICAO", Value: descricao, Inline: true},
		&discordgo.MessageEmbedField{Name: "BRANCH", Value: branch, Inline: true},
		&discordgo.MessageEmbedField{Name: "URL", Value: url, Inline: false})
	embeded.Image = &discordgo.MessageEmbedImage{URL: "https://media-exp1.licdn.com/dms/image/C4D0BAQGXoqfzgYQxuw/company-logo_200_200/0/1647390907682?e=1676505600&v=beta&t=DSlD-Bm6O3VSBwSsxeVziWYLu2GUddKnHPmqckZO1s0",
		Height: 50, Width: 100}
	embeded.Footer = &discordgo.MessageEmbedFooter{Text: "Pessoas fazem tudo!"}
	return embeded
}
