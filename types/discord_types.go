package types

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type Discord_message struct {
	Id        string
	Published *time.Time
	Content   *discordgo.MessageEmbed
}
