package formatters

import (
	"encoding/xml"
	"errors"
	"fmt"
	"slices"

	"github.com/luuppiry/luuppi-rss-service/types"
)

type RssNewsFormatter struct {
	ChannelTitle       string
	BasePath           string
	ChannelDescription string
}

func (f *RssNewsFormatter) Format(n any) ([]byte, error) {
	news := n.([]types.NormalizedNews)
	slices.SortFunc(news, func(a, b types.NormalizedNews) int {
		return a.Published.Compare(b.Published)
	})
	xo := types.Rss{}
	xo.Channel = types.Channel{Title: f.ChannelTitle, Link: f.BasePath, Description: f.ChannelDescription}
	for _, n := range news {
		xi := types.Item{}
		xi.Title = n.Title
		xi.Link = fmt.Sprintf("%s/%s", f.BasePath, n.Slug)
		xi.Description = n.Description
		xo.Channel.Items = append(xo.Channel.Items, xi)
	}
	out, err := xml.MarshalIndent(xo, "", "  ")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to marshal xml: %s in: %s", err, f))
	}
	out = []byte(xml.Header + string(out))
	return out, nil

}
func NewRssNewsFormatter(conf map[string]string) *RssNewsFormatter {
	title := conf["title"]
	basePath := conf["basePath"]
	description := conf["description"]
	return &RssNewsFormatter{
		ChannelTitle:       title,
		BasePath:           basePath,
		ChannelDescription: description,
	}
}
