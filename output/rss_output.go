package output

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/luuppiry/luuppi-rss-service/types"
)

type Rss_output struct {
	slug               string
	channelTitle       string
	basePath           string
	channelDescription string
	data               []byte
	locale             string
}

func (h *Rss_output) Initialize() error {
	http.Handle(h.slug, h)
	return nil
}

func (h *Rss_output) Update(data []types.Notification) error {
	xo := types.Rss{}
	xo.Version = "2.0"
	xo.Channel = types.Rss_channel{Title: h.channelTitle, Link: h.basePath, Description: h.channelDescription}
	for _, d := range data {
		o, err := d.Rss_format(h.basePath, h.locale)
		if err != nil {
			return err
		}
		xo.Channel.Items = append(xo.Channel.Items, o)
	}
	out, err := xml.MarshalIndent(xo, "", "  ")
	if err != nil {
		return fmt.Errorf("Failed to marshal xml: %w in: %s", err, h)
	}
	h.data = []byte(xml.Header + string(out))

	return nil
}

func (h *Rss_output) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write(h.data)
}

func NewRssOutput(conf map[string]string) *Rss_output {
	return &Rss_output{
		slug:               conf["slug"],
		channelTitle:       conf["title"],
		basePath:           conf["basePath"],
		channelDescription: conf["description"],
		locale:             conf["locale"],
	}
}
