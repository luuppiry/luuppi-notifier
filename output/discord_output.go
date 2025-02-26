package output

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/luuppiry/luuppi-rss-service/types"
)

type Discord_output struct {
	token           string
	guild           string
	channel         string
	basePath        string
	contentBasePath string
	locale          string
	connection      *discordgo.Session
	cache           map[string]cache_data
}

type luuppi_discord_header_protocol struct {
	version string
	id      string
	hash    string
}

type cache_data struct {
	msg_id string
	hash   string
}

func (h *luuppi_discord_header_protocol) encode() string {
	return fmt.Sprintf("||%s;%s;%s||", h.version, h.id, h.hash)
}
func decode(h string) (*luuppi_discord_header_protocol, error) {
	h = strings.Trim(h, "|")
	parts := strings.Split(h, ";")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Discord message header is not correctly formatted")
	}
	return &luuppi_discord_header_protocol{
		version: parts[0],
		id:      parts[1],
		hash:    parts[2],
	}, nil

}
func (d *Discord_output) Initialize() error {
	dc, err := discordgo.New("Bot " + d.token)
	if err != nil {
		return fmt.Errorf("Failed creating discord session %w", err)
	}
	d.connection = dc
	dc.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	err = dc.Open()
	if err != nil {
		return fmt.Errorf("Failed opening discord connection %w", err)
	}
	msgs, err := dc.ChannelMessages(d.channel, 50, "", "", "")
	if err != nil {
		return fmt.Errorf("Failed to fetch previous messages from discord channel %w", err)
	}
	for _, msg := range msgs {
		if len(msg.Embeds) == 1 {
			before, _, found := strings.Cut(msg.Embeds[0].Description, "\n")
			if !found {
				//Expect that this is just someone elses message and skip
				continue
			}
			h, err := decode(before)
			if err != nil || h.version != "v1" {
				//Expect that this is just someone elses message and skip
				continue
			}
			d.cache[h.id] = cache_data{msg_id: msg.ID, hash: h.hash}
		}
	}
	return nil
}

func (o *Discord_output) Update(data []types.Notification) error {
	formatted := []types.Discord_message{}
	for _, d := range data {
		f, err := d.Discord_format(o.basePath, o.locale)
		if err != nil {
			log.Printf("Formatting discord message failed %v", d)
			continue
		}
		formatted = append(formatted, f)
	}
	slices.SortFunc(formatted, func(a, b types.Discord_message) int { return a.Published.Compare(*b.Published) })
	//create new posts if any
	for _, c := range formatted {
		if _, ok := o.cache[c.Id]; ok {
			continue
		}
		post := c.Content
		hasher := sha1.New()
		hasher.Write([]byte(post.Description))
		header := luuppi_discord_header_protocol{
			id:      c.Id,
			version: "v1",
			hash:    base64.URLEncoding.EncodeToString(hasher.Sum(nil)),
		}
		post.Description = header.encode() + "\n" + post.Description
		msg, err := o.connection.ChannelMessageSendEmbed(o.channel, post)
		if err != nil {
			return err
		}
		o.cache[c.Id] = cache_data{msg_id: msg.ID, hash: header.hash}
	}
	//update posts
	for _, c := range formatted {
		cached, ok := o.cache[c.Id]
		if !ok {
			//we are only updating here, if not found do nothing
			continue
		}
		post := c.Content
		hasher := sha1.New()
		hasher.Write([]byte(post.Description))
		hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
		if hash == cached.hash {
			continue
		}
		header := luuppi_discord_header_protocol{
			id:      c.Id,
			version: "v1",
			hash:    hash,
		}
		post.Description = header.encode() + "\n" + post.Description
		msg, err := o.connection.ChannelMessageEditEmbed(o.channel, cached.msg_id, post)
		if err != nil {
			return err
		}
		o.cache[c.Id] = cache_data{msg_id: msg.ID, hash: header.hash}
	}

	return nil
}

func NewDiscordOutput(conf map[string]string) *Discord_output {
	return &Discord_output{
		token:           conf["token"],
		guild:           conf["server"],
		channel:         conf["channel"],
		basePath:        conf["basePath"],
		contentBasePath: conf["contentBasePath"],
		locale:          conf["locale"],
		cache:           map[string]cache_data{},
	}
}
