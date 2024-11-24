package types

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/luuppiry/luuppi-rss-service/utils"
)

type Strapiv4Events struct {
	Data []Strapiv4EventContents
	Meta Content_meta
}

type Strapiv4EventContents struct {
	Id         int
	Attributes StrapiEventsAttributes
}

type StrapiEventsAttributes struct {
	NameFi        string
	NameEn        string
	LocationFi    string
	LocationEn    string
	DescriptionFi utils.RTF
	DescriptionEn utils.RTF
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	PublishedAt   time.Time
	StartDate     time.Time
	EndDate       time.Time
}

func (s *Strapiv4Events) MapToNormalizedEvents() []NormalizedEvents {
	out := []NormalizedEvents{}
	for _, d := range s.Data {
		out = append(out, *mapToNormalizedEvents(d))
	}
	return out
}

func mapToNormalizedEvents(c Strapiv4EventContents) *NormalizedEvents {
	return &NormalizedEvents{
		TitleFi:       c.Attributes.NameFi,
		TitleEn:       c.Attributes.NameEn,
		DescriptionFi: utils.ParseRTFJson(c.Attributes.DescriptionFi),
		DescriptionEn: utils.ParseRTFJson(c.Attributes.DescriptionEn),
		Id:            c.Attributes.Id,
		Published:     c.Attributes.PublishedAt,
		StartTime:     c.Attributes.StartDate,
		LocationFi:    c.Attributes.LocationFi,
		LocationEn:    c.Attributes.LocationEn,
	}
}

type NormalizedEvents struct {
	TitleFi       string
	TitleEn       string
	DescriptionFi string
	DescriptionEn string
	Id            string
	Published     time.Time
	StartTime     time.Time
	LocationFi    string
	LocationEn    string
}

var descTmpl = template.Must(template.New("desc").Parse("{{.TitleFi}}\nSijainti: {{.LocationFi}}\nTapahtuma alkaa: {{.StartTime}}\n{{.DescriptionFi}}\n------\n{{.TitleEn}}\nLocation: {{.LocationEn}}\nEvent starts at: {{.StartTime}}\n{{.DescriptionEn}}"))

func (f *NormalizedEvents) Rss_format(basePath string) (Rss_item, error) {
	xi := Rss_item{}
	desc := bytes.Buffer{}
	descTmpl.Execute(&desc, f)
	link := fmt.Sprintf("%s%s", basePath, f.Id)
	xi.Title = fmt.Sprintf("%s / %s ", f.TitleFi, f.TitleEn)
	xi.Link = link
	xi.Description = desc.String()
	xi.Guid = link
	xi.PubDate = f.StartTime.String()
	return xi, nil
}
