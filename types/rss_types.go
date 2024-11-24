package types

type Rss struct {
	Channel Rss_channel `xml:"channel"`
	Version string      `xml:"version,attr"`
	XMLName struct{}    `xml:"rss"`
}

type Rss_channel struct {
	Title       string     `xml:"title"`
	Link        string     `xml:"link"`
	Description string     `xml:"description"`
	Items       []Rss_item `xml:"item"`
}

type Rss_item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Guid        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
}
