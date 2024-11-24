package fetchers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/luuppiry/luuppi-rss-service/output"
	"github.com/luuppiry/luuppi-rss-service/types"
)

type Strapiv4NewsFetcher struct {
	urls []string
}

func (f *Strapiv4NewsFetcher) Fetch() ([]output.Formattable, error) {
	ret := []output.Formattable{}
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
		sv4 := types.Strapiv4News{}
		err = json.Unmarshal(data, &sv4)
		if err != nil {
			return nil, err
		}
		news := sv4.MapToNormalizedNews()
		for _, n := range news {
			ret = append(ret, &n)
		}
	}
	return ret, nil
}
func NewStrapiv4NewsFetcher(conf map[string]string) *Strapiv4NewsFetcher {
	return &Strapiv4NewsFetcher{
		urls: strings.Split(conf["urls"], ","),
	}
}

type Strapiv4EventsFetcher struct {
	urls []string
}

func (f *Strapiv4EventsFetcher) Fetch() ([]output.Formattable, error) {
	ret := []output.Formattable{}
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
		sv4 := types.Strapiv4Events{}
		err = json.Unmarshal(data, &sv4)
		if err != nil {
			return nil, err
		}
		events := sv4.MapToNormalizedEvents()
		for _, e := range events {
			ret = append(ret, &e)
		}
	}
	return ret, nil
}
func NewStrapiv4eventsFetcher(conf map[string]string) *Strapiv4EventsFetcher {
	return &Strapiv4EventsFetcher{
		urls: strings.Split(conf["urls"], ","),
	}
}
