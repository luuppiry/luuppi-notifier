package fetchers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/luuppiry/luuppi-rss-service/types"
	"github.com/luuppiry/luuppi-rss-service/utils"
)

type Content_piece any
type Strapiv4NewsFetcher struct {
	urls []string
}

func (f *Strapiv4NewsFetcher) Fetch() ([]types.Notification, error) {
	ret := []types.Notification{}
	for _, url := range f.urls {
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		sv4 := Strapiv4News{}
		err = json.Unmarshal(data, &sv4)
		if err != nil {
			return nil, err
		}
		news := sv4.NewsToNotifications()
		for _, n := range news {
			ret = append(ret, n)
		}
	}
	return ret, nil
}
func NewStrapiv4NewsFetcher(conf map[string]string) *Strapiv4NewsFetcher {
	return &Strapiv4NewsFetcher{
		urls: strings.Split(conf["urls"], ","),
	}
}

type Localizations struct {
	Data []L12_content
}

type L12_content struct {
	Id         int
	Attributes L12_attributes
}

type L12_attributes struct {
	Title       string
	AuthorName  string
	Description string
	Content     []Content_piece
	Slug        string
	Category    string
	AuthorTitle string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Locale      string
	PublishedAt time.Time
}
type Strapiv4News struct {
	Data []Content_data
	Meta Content_meta
}

type Content_meta any

type Content_data struct {
	Id         int
	Attributes StrapiNewsAttributes
}

type Banner struct {
	Url string
}

type StrapiNewsAttributes struct {
	Title         string
	AuthorName    string
	Description   string
	Content       []Content_piece
	Slug          string
	Category      string
	AuthorTitle   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Locale        string
	PublishedAt   time.Time
	Banner        Banner
	Localizations Localizations
}

func (s *Strapiv4News) NewsToNotifications() []types.Notification {
	out := []types.Notification{}
	for _, d := range s.Data {
		out = append(out, *NewsToNotification(d))
	}
	return out
}

func NewsToNotification(c Content_data) *types.Notification {
	if len(c.Attributes.Localizations.Data) == 1 {
		return &types.Notification{
			NotificationType: types.News,
			TitleFi:          c.Attributes.Title,
			TitleEn:          c.Attributes.Localizations.Data[0].Attributes.Title,
			DescriptionFi:    c.Attributes.Description,
			DescriptionEn:    c.Attributes.Localizations.Data[0].Attributes.Description,
			SlugFi:           fmt.Sprintf("/fi/news/%s", c.Attributes.Slug),
			SlugEn:           fmt.Sprintf("/en/news/%s", c.Attributes.Slug),
			Published:        c.Attributes.PublishedAt,
			CreateTime:       c.Attributes.CreatedAt,
			UpdateTime:       c.Attributes.UpdatedAt,
			ID:               c.Attributes.Slug,
		}

	} else {
		return &types.Notification{
			NotificationType: types.News,
			TitleFi:          c.Attributes.Title,
			DescriptionFi:    c.Attributes.Description,
			FulltextFi:       "",
			SlugFi:           fmt.Sprintf("/fi/news/%s", c.Attributes.Slug),
			Published:        c.Attributes.PublishedAt,
			CreateTime:       c.Attributes.CreatedAt,
			UpdateTime:       c.Attributes.UpdatedAt,
			ID:               c.Attributes.Slug,
		}
	}
}

type Strapiv4EventsFetcher struct {
	urls []string
}

func (f *Strapiv4EventsFetcher) Fetch() ([]types.Notification, error) {
	ret := []types.Notification{}
	for _, url := range f.urls {
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		sv4 := Strapiv4Events{}
		err = json.Unmarshal(data, &sv4)
		if err != nil {
			return nil, err
		}
		events := sv4.EventsToNotifications()
		for _, e := range events {
			ret = append(ret, e)
		}
	}
	return ret, nil
}
func NewStrapiv4eventsFetcher(conf map[string]string) *Strapiv4EventsFetcher {
	return &Strapiv4EventsFetcher{
		urls: strings.Split(conf["urls"], ","),
	}
}

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

func (s *Strapiv4Events) EventsToNotifications() []types.Notification {
	out := []types.Notification{}
	for _, d := range s.Data {
		out = append(out, *eventToNotification(d))
	}
	return out
}

func eventToNotification(c Strapiv4EventContents) *types.Notification {
	return &types.Notification{
		NotificationType: types.Event,
		TitleFi:          c.Attributes.NameFi,
		TitleEn:          c.Attributes.NameEn,
		DescriptionFi:    utils.ParseRTFJson(c.Attributes.DescriptionFi),
		DescriptionEn:    utils.ParseRTFJson(c.Attributes.DescriptionEn),
		ID:               c.Attributes.Id,
		SlugFi:           fmt.Sprintf("/fi/events/%d", c.Id),
		SlugEn:           fmt.Sprintf("/en/events/%d", c.Id),
		Published:        c.Attributes.PublishedAt,
		StartTime:        c.Attributes.StartDate,
		CreateTime:       c.Attributes.CreatedAt,
		UpdateTime:       c.Attributes.UpdatedAt,
		EndTime:          c.Attributes.EndDate,
		LocationFi:       c.Attributes.LocationFi,
		LocationEn:       c.Attributes.LocationEn,
	}
}
