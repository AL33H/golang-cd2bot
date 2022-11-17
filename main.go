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

const chatAssai = "1042884138425913486"
const chatHippo = "1042883976773251202"

func main() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Comando inválido!")
		}
	}()

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Erro ao criar sessão do discord.", err)
		return
	}
	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("Erro ao abrir conexão com discord!", err)
		return
	}

	fmt.Println("Bot rodando com sucesso!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.Contains(m.Content, "!mrassai") {
		msg := extractInfoMessageAssai(m)
		embeded := generateEmbedAssai(m, msg)
		messageCreate, _ := s.ChannelMessageSendEmbed(chatAssai, embeded)

		s.MessageReactionAdd(messageCreate.ChannelID, messageCreate.ID, "✅")
		s.MessageReactionAdd(messageCreate.ChannelID, messageCreate.ID, "❌")
	}

	if strings.Contains(m.Content, "!mrhippo") {
		msg := extractInfoMessageHippo(m)
		embeded := generateEmbedHippo(m, msg)
		messageCreate, _ := s.ChannelMessageSendEmbed(chatHippo, embeded)

		s.MessageReactionAdd(messageCreate.ChannelID, messageCreate.ID, "✅")
		s.MessageReactionAdd(messageCreate.ChannelID, messageCreate.ID, "❌")
	}

	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

func extractInfoMessageAssai(m *discordgo.MessageCreate) mensagem {
	splitedContent := strings.Split(m.Content, " ")
	urlMr := strings.Split(splitedContent[2], "/")
	urlBranch := strings.Split(splitedContent[3], "/")

	msg := mensagem{
		ambiente:           splitedContent[1],
		url:                splitedContent[2],
		projeto:            urlMr[5],
		numeroMr:           urlMr[8],
		branch:             splitedContent[3],
		atividadeId:        urlBranch[1],
		formatModification: urlBranch[0],
		descricao:          urlBranch[2],
	}

	return msg
}

func extractInfoMessageHippo(m *discordgo.MessageCreate) mensagem {
	splitedContent := strings.Split(m.Content, " ")
	urlMr := strings.Split(splitedContent[2], "/")
	urlBranch := strings.Split(splitedContent[3], "/")

	msg := mensagem{
		ambiente:  splitedContent[1],
		url:       splitedContent[2],
		projeto:   urlMr[5],
		branch:    splitedContent[3],
		descricao: urlBranch[2],
	}

	return msg
}

func generateEmbedAssai(m *discordgo.MessageCreate, msg mensagem) *discordgo.MessageEmbed {
	embeded := embed.NewGenericEmbed("Solicitação de MR ", m.Author.Username+" ~ "+time.Now().Format(time.RFC822))
	embeded.Color = 0xe6cd53
	embeded.URL = msg.url
	embeded.Fields = append(embeded.Fields,
		&discordgo.MessageEmbedField{Name: "AMBIENTE", Value: msg.ambiente, Inline: true},
		&discordgo.MessageEmbedField{Name: "MR", Value: msg.numeroMr, Inline: true},
		&discordgo.MessageEmbedField{Name: "PROJETO", Value: msg.projeto, Inline: true},
		&discordgo.MessageEmbedField{Name: "TIPO", Value: msg.formatModification, Inline: true},
		&discordgo.MessageEmbedField{Name: "ATIVIDADE", Value: msg.atividadeId, Inline: true},
		&discordgo.MessageEmbedField{Name: "DESCRICAO", Value: msg.descricao, Inline: true},
		&discordgo.MessageEmbedField{Name: "BRANCH", Value: msg.branch, Inline: true},
		&discordgo.MessageEmbedField{Name: "URL", Value: msg.url, Inline: false})
	embeded.Image = &discordgo.MessageEmbedImage{URL: "https://media-exp1.licdn.com/dms/image/C4D0BAQGXoqfzgYQxuw/company-logo_200_200/0/1647390907682?e=1676505600&v=beta&t=DSlD-Bm6O3VSBwSsxeVziWYLu2GUddKnHPmqckZO1s0",
		Height: 50, Width: 100}
	embeded.Footer = &discordgo.MessageEmbedFooter{Text: "Pessoas fazem tudo!"}
	return embeded
}

func generateEmbedHippo(m *discordgo.MessageCreate, msg mensagem) *discordgo.MessageEmbed {

	embeded := embed.NewGenericEmbed("Solicitação de MR ", m.Author.Username+" ~ "+time.Now().Format(time.RFC822))
	embeded.Color = 0xe6cd53
	embeded.URL = msg.url
	embeded.Fields = append(embeded.Fields,
		&discordgo.MessageEmbedField{Name: "AMBIENTE", Value: msg.ambiente, Inline: true},
		&discordgo.MessageEmbedField{Name: "PROJETO", Value: msg.projeto, Inline: true},
		&discordgo.MessageEmbedField{Name: "DESCRICAO", Value: msg.descricao, Inline: true},
		&discordgo.MessageEmbedField{Name: "BRANCH", Value: msg.branch, Inline: true},
		&discordgo.MessageEmbedField{Name: "URL", Value: msg.url, Inline: false})
	embeded.Image = &discordgo.MessageEmbedImage{URL: "https://media-exp1.licdn.com/dms/image/C4D0BAQGXoqfzgYQxuw/company-logo_200_200/0/1647390907682?e=1676505600&v=beta&t=DSlD-Bm6O3VSBwSsxeVziWYLu2GUddKnHPmqckZO1s0",
		Height: 50, Width: 100}
	embeded.Footer = &discordgo.MessageEmbedFooter{Text: "Pessoas fazem tudo!"}
	return embeded
}

type mensagem struct {
	ambiente           string
	url                string
	projeto            string
	numeroMr           string
	branch             string
	atividadeId        string
	formatModification string
	descricao          string
}
