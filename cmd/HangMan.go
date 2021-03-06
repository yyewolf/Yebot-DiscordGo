package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"DiscGo.discordgo/config"
	"github.com/bwmarrin/discordgo"
)

type Resp struct {
	Success bool     `json:"success"`
	Data    []string `json:"words"`
}

type HMGame struct {
	CreatedAt time.Time
	Word      string
	Category  string
	GuildID   string
	ChannelID string
	Guesses   []string
	Lost      int
	finished  bool
	Message   *discordgo.Message
}

var AllHMGames []HMGame

func RemoveHMGame(s []HMGame, ChannelID string) []HMGame {
	i := 0
	exist := false
	for i = 0; i < len(s); i++ {
		if s[i].ChannelID == ChannelID {
			exist = true
			break
		}
	}
	if exist {
		s[i] = s[len(s)-1]
		return s[:len(s)-1]
	} else {
		return s
	}
}

func Try() (string, string, error) {
	HangManCat := []string{
		"1;Animals",
		"56;Common Animals",
		"4;Places",
		"65;Sports",
		"2;Food and Cooking",
		"63;Nature",
		"3;People",
		"52;Around the House",
		"53;Around the Office",
		"66;Travel",
		"95;Categories",
		"55;Colors",
		"57;Dog Breeds",
		"59;Feelings and Emotions",
		"58;English Litterature",
		"60;Food and Cooking",
		"61;Math",
		"62;Music",
		"64;Science",
		"54;Art",
	}
	k := rand.Intn(len(HangManCat))
	resp, err := http.Get("https://www.thegamegal.com/wordgenerator/generator.php?game=1&category=" + strings.Split(HangManCat[k], ";")[0])
	if err != nil {
		return "", "", err
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	res := Resp{}
	json.Unmarshal([]byte(body), &res)
	i := rand.Intn(len(res.Data))
	Word := res.Data[i]
	return Word, strings.Split(HangManCat[k], ";")[1], nil
}

func isWin(s HMGame) bool {
	win := true
	t := HMFields(s)
	if strings.Contains(t[2].Value, ":large_blue_diamond:") {
		win = false
	}
	return win
}

func HMFields(s HMGame) []*discordgo.MessageEmbedField {
	value := ""  // Field Word
	value2 := "" // Field Guesses
	for i := 0; i < len(s.Word); i++ {
		letter := false
		for k := 0; k < len(s.Guesses); k++ {
			if s.Guesses[k] == string([]byte(s.Word)[i]) {
				value += ":regional_indicator_" + string([]byte(s.Word)[i]) + ": "
				letter = true
			}

			if i == 0 && s.Guesses[k] != "" && !strings.Contains(value2, "**"+s.Guesses[k]+"** ") {
				value2 += "**" + s.Guesses[k] + "** "
			}
		}
		if !letter && string([]byte(s.Word)[i]) != " " && string([]byte(s.Word)[i]) != "-" {
			value += ":large_blue_diamond: "
		} else if !letter && string([]byte(s.Word)[i]) == " " {
			value += ":white_circle: "
		} else if !letter && string([]byte(s.Word)[i]) == "-" {
			value += ":white_large_square: "
		}
	}
	if value2 == "" {
		value2 = "No guesses yet."
	}
	FieldGuesses := &discordgo.MessageEmbedField{
		Name:   "Guesses :",
		Value:  value2,
		Inline: true,
	}
	FieldCat := &discordgo.MessageEmbedField{
		Name:   "Category :",
		Value:  "**" + s.Category + "**",
		Inline: true,
	}
	FieldWord := &discordgo.MessageEmbedField{
		Name:   "Word :",
		Value:  value,
		Inline: false,
	}
	Fields := []*discordgo.MessageEmbedField{FieldGuesses, FieldCat, FieldWord}
	return Fields
}

func HMPlay(s *discordgo.Session, m *discordgo.MessageCreate) {
	HangManImages := []string{
		"https://media.discordapp.net/attachments/364445463333830678/566352682827251733/0.jpg",
		"https://media.discordapp.net/attachments/364445463333830678/566352693451423765/1.jpg",
		"https://media.discordapp.net/attachments/364445463333830678/566352699315060737/2.jpg",
		"https://media.discordapp.net/attachments/364445463333830678/566352703366758450/3.jpg",
		"https://media.discordapp.net/attachments/364445463333830678/566352708534009856/4.jpg",
		"https://cdn.discordapp.com/attachments/364445463333830678/566352715358273556/5.jpg",
		"https://cdn.discordapp.com/attachments/364445463333830678/566352719460302852/6.jpg",
		"https://cdn.discordapp.com/attachments/364445463333830678/566352725483454476/7.jpg",
		"https://cdn.discordapp.com/attachments/364445463333830678/566352730202046482/8.jpg",
		"https://cdn.discordapp.com/attachments/364445463333830678/566352735838928958/9.jpg",
		"https://cdn.discordapp.com/attachments/364445463333830678/566352739412606996/10.jpg",
	}
	i := 0
	exist := false
	for i = 0; i < len(AllHMGames); i++ {
		if AllHMGames[i].ChannelID == m.ChannelID && !AllHMGames[i].finished {
			exist = true
			break
		}
	}
	if !exist {
		WordString, Category, err := Try()
		if err != nil {
			return
		}
		NewGame := HMGame{
			CreatedAt: time.Now(),
			ChannelID: m.ChannelID,
			GuildID:   m.GuildID,
			Word:      WordString,
			Category:  Category,
			Guesses:   []string{},
			Lost:      0,
			finished:  false,
		}

		AvatarURL := m.Author.AvatarURL("512")
		Fields := HMFields(NewGame)
		AllFields := []*discordgo.MessageEmbedField{Fields[0], Fields[1], Fields[2]}
		embed := &discordgo.MessageEmbed{
			Title:       "Hangman",
			Description: "You have to type `g!h` followed by a letter to submit a guess !",
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: HangManImages[0],
			},
			Footer: &discordgo.MessageEmbedFooter{
				IconURL: AvatarURL,
				Text:    m.Author.Username,
			},
			Fields: AllFields,
			Color:  0xFFDD00,
		}
		Msg, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			return
		}
		NewGame.Message = Msg

		AllHMGames = append(AllHMGames, NewGame)
	}
}

func HM(s *discordgo.Session, m *discordgo.MessageCreate) {
	HangManImages := []string{
		"https://media.discordapp.net/attachments/364445463333830678/566352682827251733/0.jpg",
		"https://media.discordapp.net/attachments/364445463333830678/566352693451423765/1.jpg",
		"https://media.discordapp.net/attachments/364445463333830678/566352699315060737/2.jpg",
		"https://media.discordapp.net/attachments/364445463333830678/566352703366758450/3.jpg",
		"https://media.discordapp.net/attachments/364445463333830678/566352708534009856/4.jpg",
		"https://cdn.discordapp.com/attachments/364445463333830678/566352715358273556/5.jpg",
		"https://cdn.discordapp.com/attachments/364445463333830678/566352719460302852/6.jpg",
		"https://cdn.discordapp.com/attachments/364445463333830678/566352725483454476/7.jpg",
		"https://cdn.discordapp.com/attachments/364445463333830678/566352730202046482/8.jpg",
		"https://cdn.discordapp.com/attachments/364445463333830678/566352735838928958/9.jpg",
		"https://cdn.discordapp.com/attachments/364445463333830678/566352739412606996/10.jpg",
	}
	message := strings.Replace(m.Content, config.Prefix+"h ", "", 1)
	command := strings.Split(message, " ")[0]

	i := 0
	exist := false
	for i = 0; i < len(AllHMGames); i++ {
		if AllHMGames[i].ChannelID == m.ChannelID && !AllHMGames[i].finished {
			exist = true
			break
		}
	}
	if exist && len(command) < 2 && strings.Contains(AllHMGames[i].Word, command) {
		//SUCCESS
		err := s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			return
		}

		if !config.Contains(AllHMGames[i].Guesses, command) {
			AllHMGames[i].Guesses = append(AllHMGames[i].Guesses, command)
		}

		AvatarURL := m.Author.AvatarURL("512")
		Fields := HMFields(AllHMGames[i])
		AllFields := []*discordgo.MessageEmbedField{Fields[0], Fields[1], Fields[2]}
		embed := &discordgo.MessageEmbed{
			Title:       "Hangman",
			Description: "You have to type `g!h` followed by a letter to submit a guess !",
			Footer: &discordgo.MessageEmbedFooter{
				IconURL: AvatarURL,
				Text:    m.Author.Username,
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: HangManImages[AllHMGames[i].Lost],
			},
			Fields: AllFields,
			Color:  0xFFDD00,
		}

		if isWin(AllHMGames[i]) {
			embed.Description = "You found the word ! Congratulations"
			AllHMGames[i].finished = true
			RemoveHMGame(AllHMGames, m.ChannelID)
		}

		Edit := &discordgo.MessageEdit{
			ID:      AllHMGames[i].Message.ID,
			Channel: AllHMGames[i].Message.ChannelID,
			Embed:   embed,
		}
		s.ChannelMessageEditComplex(Edit)

	} else if exist && len(command) < 2 {
		//FAIL
		err := s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			return
		}
		if !config.Contains(AllHMGames[i].Guesses, command) {
			AllHMGames[i].Guesses = append(AllHMGames[i].Guesses, command)
		}
		AllHMGames[i].Lost += 1
		AvatarURL := m.Author.AvatarURL("512")

		Fields := HMFields(AllHMGames[i])
		AllFields := []*discordgo.MessageEmbedField{Fields[0], Fields[1], Fields[2]}
		embed := &discordgo.MessageEmbed{
			Title:       "Hangman",
			Description: "You have to type `g!h` followed by a letter to submit a guess !",
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: HangManImages[AllHMGames[i].Lost],
			},
			Footer: &discordgo.MessageEmbedFooter{
				IconURL: AvatarURL,
				Text:    m.Author.Username,
			},
			Fields: AllFields,
			Color:  0xFFDD00,
		}
		if AllHMGames[i].Lost == 10 {
			embed.Description = "You lost ! The word was : **" + AllHMGames[i].Word + "**"
			AllHMGames[i].finished = true
			RemoveHMGame(AllHMGames, m.ChannelID)
		}
		Edit := &discordgo.MessageEdit{
			ID:      AllHMGames[i].Message.ID,
			Channel: AllHMGames[i].Message.ChannelID,
			Embed:   embed,
		}
		s.ChannelMessageEditComplex(Edit)
	}
}
