package types

import (
	"fmt"
	"time"
)

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
	Banner      Banner
}

func (s *Strapiv4News) MapToNormalizedNews() []NormalizedNews {
	out := []NormalizedNews{}
	for _, d := range s.Data {
		out = append(out, *mapToNormalizedNews(d))
	}
	return out
}

func mapToNormalizedNews(c Content_data) *NormalizedNews {
	return &NormalizedNews{
		Title:       c.Attributes.Title,
		Author:      c.Attributes.AuthorName,
		Description: c.Attributes.Description,
		FullText:    "",
		Slug:        c.Attributes.Slug,
		Published:   c.Attributes.PublishedAt,
		Locale:      c.Attributes.Locale,
		BannerURL:   c.Attributes.Banner.Url,
	}
}

type Content_piece any

type NormalizedNews struct {
	Title       string
	Author      string
	Description string
	FullText    string
	Slug        string
	Published   time.Time
	Locale      string
	BannerURL   string
}

func (f *NormalizedNews) Rss_format(basePath string) (Rss_item, error) {
	xi := Rss_item{}
	link := fmt.Sprintf("%s%s", basePath, f.Slug)
	xi.Title = f.Title
	xi.Link = link
	xi.Description = f.Description
	xi.Guid = link
	xi.PubDate = f.Published.String()
	return xi, nil
}

func (f *NormalizedNews) Discord_format() (Discord_message, error) {
	content := f.Description
	return Discord_message{
		Id:        f.Slug,
		Content:   content,
		Published: &f.Published,
		Locale:    f.Locale,
		Title:     f.Title,
		Image:     f.BannerURL,
	}, nil
}
