package commands

import (
	"strconv"
	"strings"

	"DiscGo.discordgo/perm"
	"github.com/bwmarrin/discordgo"
)

func Kick(s *discordgo.Session, m *discordgo.MessageCreate) {
	member, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		return
	}
	isHe := Permissions.HasPermission(member, s, m.GuildID, Permissions.PERM_ADMINISTRATOR)
	if isHe {
		if len(m.Mentions) != 0 {
			s.GuildMemberDelete(m.GuildID, m.Mentions[0].ID)
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "You're not an administrator !")
	}
}

func Ban(s *discordgo.Session, m *discordgo.MessageCreate) {
	member, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		return
	}
	isHe := Permissions.HasPermission(member, s, m.GuildID, Permissions.PERM_ADMINISTRATOR)
	if isHe {
		if len(m.Mentions) != 0 {
			s.GuildBanCreate(m.GuildID, m.Mentions[0].ID, 1)
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "You're not an administrator !")
	}
}

func Clear(s *discordgo.Session, m *discordgo.MessageCreate) {
	member, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		return
	}
	isHe := Permissions.HasPermission(member, s, m.GuildID, Permissions.PERM_ADMINISTRATOR)
	if isHe {
		CommandSplit := strings.Split(m.Content, " ")
		number, err := strconv.Atoi(CommandSplit[1])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "I had some problem reading the number you entered.")
		} else if number < 100 {
			MessageSlice, err := s.ChannelMessages(m.ChannelID, number, "", "", "")
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "I had some problem deleting the messages.")
			} else {
				var Params []string
				for i := 0; i < len(MessageSlice); i++ {
					Params = append(Params, MessageSlice[i].ID)
				}
				s.ChannelMessagesBulkDelete(m.ChannelID, Params)
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "I cannot delete more than 100 messages.")
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "You're not an administrator !")
	}
}

func Ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Pong!")
}
