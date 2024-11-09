package types

import "time"

type Strapiv4News struct {
	Data []Content_data
	Meta Content_meta
}

type Content_meta any

type Content_data struct {
	Id         int
	Attributes StrapiNewsAttributes
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
}
