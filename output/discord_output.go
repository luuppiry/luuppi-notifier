package output

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/luuppiry/luuppi-rss-service/types"
)

type Discord_output struct {
	token           string
	guild           string
	channel         string
	basePath        string
	contentBasePath string
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
		before, _, found := strings.Cut(msg.Content, "\n")
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

	return nil
}

func (d *Discord_output) Update(data []Formattable) error {
	formatted := []types.Discord_message{}
	for _, d := range data {
		f, err := d.Discord_format()
		if err != nil {
			log.Printf("Formatting discord message failed %v", d)
			continue
		}
		formatted = append(formatted, f)
	}
	correlated := correlate(formatted)
	sort.Slice(correlated, func(i, j int) bool { return correlated[i][0].Published.Compare(*correlated[j][0].Published) > 0 })
	//create new posts if any
	for _, c := range correlated {
		id, post := combine(c, d.basePath, d.contentBasePath)
		if _, ok := d.cache[id]; ok {
			break
		}
		hasher := sha1.New()
		hasher.Write([]byte(post.Description))
		header := luuppi_discord_header_protocol{
			id:      id,
			version: "v1",
			hash:    base64.URLEncoding.EncodeToString(hasher.Sum(nil)),
		}
		post.Description = header.encode() + "\n" + post.Description
		msg, err := d.connection.ChannelMessageSendEmbed(d.channel, post)
		if err != nil {
			return err
		}
		d.cache[id] = cache_data{msg_id: msg.ID, hash: header.hash}
	}
	//update posts
	for _, c := range correlated {
		id, post := combine(c, d.basePath, d.contentBasePath)
		cached, ok := d.cache[id]
		if !ok {
			//we are only updating here, if not found do nothing
			continue
		}
		hasher := sha1.New()
		hasher.Write([]byte(post.Description))
		hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
		if hash == cached.hash {
			continue
		}
		header := luuppi_discord_header_protocol{
			id:      id,
			version: "v1",
			hash:    hash,
		}
		post.Description = header.encode() + "\n" + post.Description
		msg, err := d.connection.ChannelMessageEditEmbed(d.channel, cached.msg_id, post)
		if err != nil {
			return err
		}
		d.cache[id] = cache_data{msg_id: msg.ID, hash: header.hash}
	}

	return nil
}
func correlate(data []types.Discord_message) [][]types.Discord_message {
	ret := [][]types.Discord_message{}
	cutoff := len(data)
	for i := 0; i < cutoff; i++ {
		n := data[i]
		c := closure(*n.Published)
		same := []types.Discord_message{}
		same = append(same, n)
		start := i + 1
		for true {
			matching := slices.IndexFunc(data[start:cutoff], c)
			if matching == -1 {
				break
			}
			same = append(same, data[matching+start])
			data[matching+start], data[cutoff-1] = data[cutoff-1], data[matching+start]
			cutoff--
		}
		ret = append(ret, same)
	}
	return ret
}

func closure(a time.Time) func(types.Discord_message) bool {
	return func(nn types.Discord_message) bool {
		return isCloseInTime(a, *nn.Published)
	}
}

func isCloseInTime(a, b time.Time) bool {
	return a.Truncate(30 * time.Minute).Equal(b.Truncate(30 * time.Minute))
}

func combine(data []types.Discord_message, basePath string, contentBasePath string) (string, *discordgo.MessageEmbed) {
	sort.Slice(data, func(i, j int) bool { return data[i].Id < data[j].Id })
	ids := []string{}
	contents := []string{}
	titles := []string{}
	slug := ""
	banner := ""
	ind := slices.IndexFunc(data, func(e types.Discord_message) bool { return e.Locale == "fi" })
	if ind != -1 {
		slug = data[ind].Id
		banner = data[ind].Image
	}
	for _, d := range data {
		ids = append(ids, d.Id)
		titles = append(titles, d.Title)
		withLink := d.Content + "\n" + basePath + "/" + d.Locale + "/news/" + slug
		contents = append(contents, withLink)
	}
	id := strings.Join(ids, "")
	title := strings.Join(titles, "/")
	content := strings.Join(contents, "\n-----\n")
	return id, &discordgo.MessageEmbed{
		URL:         basePath + "/fi" + "/news/" + slug,
		Type:        discordgo.EmbedTypeArticle,
		Description: content,
		Title:       title,
		Image:       &discordgo.MessageEmbedImage{URL: contentBasePath + banner},
	}
}

func NewDiscordOutput(conf map[string]string) *Discord_output {
	return &Discord_output{
		token:           conf["token"],
		guild:           conf["server"],
		channel:         conf["channel"],
		basePath:        conf["basePath"],
		contentBasePath: conf["contentBasePath"],
		cache:           map[string]cache_data{},
	}
}
