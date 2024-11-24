package formatters

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"slices"
	"text/template"

	"github.com/luuppiry/luuppi-rss-service/types"
)

type RssEventsFormatter struct {
	ChannelTitle       string
	BasePath           string
	ChannelDescription string
}

var descTmpl = template.Must(template.New("desc").Parse("{{.TitleFi}}\nSijainti: {{.LocationFi}}\nTapahtuma alkaa: {{.StartTime}}\n{{.DescriptionFi}}\n------\n{{.TitleEn}}\nLocation: {{.LocationEn}}\nEvent starts at: {{.StartTime}}\n{{.DescriptionEn}}"))

func (f *RssEventsFormatter) Format(n any) ([]byte, error) {
	news := n.([]types.NormalizedEvents)
	slices.SortFunc(news, func(a, b types.NormalizedEvents) int {
		return a.Published.Compare(b.Published)
	})
	xo := types.Rss{}
	xo.Version = "2.0"
	xo.Channel = types.Channel{Title: f.ChannelTitle, Link: f.BasePath, Description: f.ChannelDescription}
	for _, n := range news {
		xi := types.Item{}
		desc := bytes.Buffer{}
		descTmpl.Execute(&desc, n)
		link := fmt.Sprintf("%s%s", f.BasePath, n.Id)
		xi.Title = fmt.Sprintf("%s / %s ", n.TitleFi, n.TitleEn)
		xi.Link = link
		xi.Description = desc.String()
		xi.Guid = link
		xi.PubDate = n.StartTime.String()
		xo.Channel.Items = append(xo.Channel.Items, xi)
	}
	out, err := xml.MarshalIndent(xo, "", "  ")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to marshal xml: %s in: %s", err, f))
	}
	out = []byte(xml.Header + string(out))
	return out, nil

}
func NewRssEventsFormatter(conf map[string]string) *RssEventsFormatter {
	title := conf["title"]
	basePath := conf["basePath"]
	description := conf["description"]
	return &RssEventsFormatter{
		ChannelTitle:       title,
		BasePath:           basePath,
		ChannelDescription: description,
	}
}
