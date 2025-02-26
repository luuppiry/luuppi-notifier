package types

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/bwmarrin/discordgo"
)

type NotificationType int

const (
	News NotificationType = iota
	Event
)

type Notification struct {
	NotificationType NotificationType
	ID               string
	Published        time.Time
	StartTime        time.Time
	EndTime          time.Time
	CreateTime       time.Time
	UpdateTime       time.Time
	SlugEn           string
	SlugFi           string
	TitleEn          string
	TitleFi          string
	DescriptionEn    string
	DescriptionFi    string
	FulltextEn       string
	FulltextFi       string
	LocationEn       string
	LocationFi       string
}

var eventRssTemplateFi = template.Must(template.New("descFi").Parse("{{.TitleFi}}\nSijainti: {{.LocationFi}}\nTapahtuma alkaa: {{.StartTime}}\n{{.DescriptionFi}}"))
var eventRssTemplateEn = template.Must(template.New("descEn").Parse("{{.TitleEn}}\nLocation: {{.LocationEn}}\nEvent starts at: {{.StartTime}}\n{{.DescriptionEn}}"))

var newsRssTemplateFi = template.Must(template.New("descFi").Parse("{{.DescriptionFi}}"))
var newsRssTemplateEn = template.Must(template.New("descEn").Parse("{{.DescriptionEn}}"))

func (f *Notification) Rss_format(basePath string, locale string) (Rss_item, error) {
	if f.NotificationType == Event {
		xi := Rss_item{}
		desc := bytes.Buffer{}
		link := ""
		title := ""
		if locale == "fi" {
			link = fmt.Sprintf("%s%s", basePath, f.SlugFi)
			eventRssTemplateFi.Execute(&desc, f)
			title = f.TitleFi
		} else if locale == "en" {
			link = fmt.Sprintf("%s%s", basePath, f.SlugEn)
			eventRssTemplateEn.Execute(&desc, f)
			title = f.TitleEn

		} else {
			link = fmt.Sprintf("%s%s", basePath, f.SlugFi)
			eventRssTemplateFi.Execute(&desc, f)
			desc.Write([]byte("\n----------------\n"))
			eventRssTemplateEn.Execute(&desc, f)
			title = fmt.Sprintf("%s / %s", f.TitleFi, f.TitleEn)

		}

		xi.Title = title
		xi.Link = link
		xi.Description = desc.String()
		xi.Guid = link
		xi.PubDate = f.Published.String()
		return xi, nil
	}
	if f.NotificationType == News {
		xi := Rss_item{}
		desc := bytes.Buffer{}
		link := ""
		title := ""
		if locale == "fi" {
			link = fmt.Sprintf("%s%s", basePath, f.SlugFi)
			newsRssTemplateFi.Execute(&desc, f)
			title = f.TitleFi
		} else if locale == "en" {
			link = fmt.Sprintf("%s%s", basePath, f.SlugEn)
			newsRssTemplateEn.Execute(&desc, f)
			title = f.TitleEn

		} else {
			link = fmt.Sprintf("%s%s", basePath, f.SlugFi)
			newsRssTemplateFi.Execute(&desc, f)
			desc.Write([]byte("\n----------------\n"))
			newsRssTemplateEn.Execute(&desc, f)
			title = fmt.Sprintf("%s / %s", f.TitleFi, f.TitleEn)

		}
		xi.Title = title
		xi.Link = link
		xi.Description = desc.String()
		xi.Guid = link
		xi.PubDate = f.Published.String()
		return xi, nil
	}
	return Rss_item{}, fmt.Errorf("Unknown notification type")

}
func (f *Notification) Discord_format(basePath string, locale string) (Discord_message, error) {
	link := ""
	title := ""
	content := ""
	if locale == "fi" {
		link = fmt.Sprintf("%s%s", basePath, f.SlugFi)
		title = f.TitleFi
		content = f.DescriptionFi
	} else if locale == "en" {
		link = fmt.Sprintf("%s%s", basePath, f.SlugEn)
		title = f.TitleEn
		content = f.DescriptionEn

	} else {
		link = fmt.Sprintf("%s%s", basePath, f.SlugFi)
		title = fmt.Sprintf("%s / %s", f.TitleFi, f.TitleEn)
		content = fmt.Sprintf("%s\n----------\n%s", f.DescriptionFi, f.DescriptionEn)
	}
	msg := &discordgo.MessageEmbed{
		URL:         link,
		Type:        discordgo.EmbedTypeArticle,
		Title:       title,
		Description: content,
	}
	return Discord_message{
		Id:        f.ID,
		Content:   msg,
		Published: &f.Published,
	}, nil
}
