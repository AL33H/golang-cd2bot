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
	dg.AddHandler(interaction)
	dg.AddHandler(ready)
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions | discordgo.IntentGuildEmojis

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

func interaction(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Comando inválido! interacttion")
		}
	}()

	msg, _ := s.ChannelMessage(m.ChannelID, m.MessageID)

	if len(msg.Embeds) == 0 || m.UserID == "1041368492655530065" {
		return
	}

	if m.Emoji.Name == "✅" {
		msg.Embeds[0].Color = 0x40cf1d
		s.ChannelMessageEditEmbed(m.ChannelID, m.MessageID, msg.Embeds[0])
	}

	if m.Emoji.Name == "❌" {
		msg.Embeds[0].Color = 0xd61a1a
		s.ChannelMessageEditEmbed(m.ChannelID, m.MessageID, msg.Embeds[0])
	}
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateGameStatus(1, "Pessoas fazem tudo!")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer func() {
		if r := recover(); r != nil {
			userChat, err := s.UserChannelCreate(m.Author.ID)
			if err == nil {
				s.ChannelMessageSend(userChat.ID, "A mensagem \""+m.Content+"\" fere algum padrão encontrado na implementação. Por favor, necessário ajuste para continuar!")
			}
		}
	}()

	if strings.Contains(m.Content, "!mrassai") {
		msg := extractInfoMessageAssai(m)
		embeded := generateEmbedAssai(m, msg)
		messageCreate, _ := s.ChannelMessageSendEmbed(chatAssai, embeded)
		fmt.Println(&messageCreate)
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

	if m.Content == "!help assai" {
		s.ChannelMessageSend(m.ChannelID, "O Padrão de solicitação é \"!mrassai ambiente linkMr nomebranch\"")
	}

	if m.Content == "!help hippo" {
		s.ChannelMessageSend(m.ChannelID, "O Padrão de solicitação é \"!mrahippo ambiente linkMr nomebranch\"")
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
		projeto:   urlMr[4],
		branch:    splitedContent[3],
		descricao: urlBranch[1],
	}

	return msg
}

func generateEmbedAssai(m *discordgo.MessageCreate, msg mensagem) *discordgo.MessageEmbed {

	embeded := embed.NewGenericEmbed("Solicitação de MR ", m.Author.Username+" ~ "+time.Now().Format("02/01/2006  15:04"))
	embeded.Color = 0xe6cd53
	embeded.URL = msg.url
	embeded.Fields = append(embeded.Fields,
		&discordgo.MessageEmbedField{Name: "AMBIENTE", Value: msg.ambiente, Inline: false},
		&discordgo.MessageEmbedField{Name: "MR", Value: msg.numeroMr, Inline: false},
		&discordgo.MessageEmbedField{Name: "PROJETO", Value: msg.projeto, Inline: false},
		&discordgo.MessageEmbedField{Name: "TIPO", Value: msg.formatModification, Inline: false},
		&discordgo.MessageEmbedField{Name: "ATIVIDADE", Value: msg.atividadeId, Inline: false},
		&discordgo.MessageEmbedField{Name: "DESCRICAO", Value: msg.descricao, Inline: false},
		&discordgo.MessageEmbedField{Name: "BRANCH", Value: msg.branch, Inline: false},
		&discordgo.MessageEmbedField{Name: "URL", Value: msg.url, Inline: false})
	embeded.Image = &discordgo.MessageEmbedImage{URL: "https://media-exp1.licdn.com/dms/image/C4D0BAQGXoqfzgYQxuw/company-logo_200_200/0/1647390907682?e=1676505600&v=beta&t=DSlD-Bm6O3VSBwSsxeVziWYLu2GUddKnHPmqckZO1s0",
		Height: 50, Width: 100}
	embeded.Footer = &discordgo.MessageEmbedFooter{Text: "Pessoas fazem tudo!"}
	return embeded
}

func generateEmbedHippo(m *discordgo.MessageCreate, msg mensagem) *discordgo.MessageEmbed {

	embeded := embed.NewGenericEmbed("Solicitação de MR ", m.Author.Username+" ~ "+time.Now().Format("02/01/2006  15:04"))
	embeded.Color = 0xe6cd53
	embeded.URL = msg.url
	embeded.Fields = append(embeded.Fields,
		&discordgo.MessageEmbedField{Name: "AMBIENTE", Value: msg.ambiente, Inline: false},
		&discordgo.MessageEmbedField{Name: "PROJETO", Value: msg.projeto, Inline: false},
		&discordgo.MessageEmbedField{Name: "DESCRICAO", Value: msg.descricao, Inline: false},
		&discordgo.MessageEmbedField{Name: "BRANCH", Value: msg.branch, Inline: false},
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
