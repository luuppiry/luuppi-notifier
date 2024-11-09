package fetchers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/luuppiry/luuppi-rss-service/types"
)

type Strapiv4NewsFetcher struct {
	url string
}

func (f *Strapiv4NewsFetcher) Fetch() (any, error) {
	res, err := http.Get(f.url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	sv4 := types.Strapiv4News{}
	err = json.Unmarshal(data, &sv4)
	if err != nil {
		return nil, err
	}
	return sv4.MapToNormalizedNews(), nil
}
func NewStrapiv4NewsFetcher(conf map[string]string) *Strapiv4NewsFetcher {
	return &Strapiv4NewsFetcher{
		url: conf["url"],
	}
}
