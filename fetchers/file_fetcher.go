package fetchers

import (
	"encoding/json"
	"io"
	"os"

	"github.com/luuppiry/luuppi-rss-service/types"
)

type FileNewsFetcher struct {
	files []string
}

func (f *FileNewsFetcher) Fetch() ([]types.Notification, error) {
	ret := []types.Notification{}
	for _, file := range f.files {
		res, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer res.Close()
		data, err := io.ReadAll(res)
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
func NewFileNewsFetcher(files []string) *FileNewsFetcher {
	return &FileNewsFetcher{files: files}
}

type FileEventsFetcher struct {
	files []string
}

func (f *FileEventsFetcher) Fetch() ([]types.Notification, error) {
	ret := []types.Notification{}
	for _, file := range f.files {
		res, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer res.Close()
		data, err := io.ReadAll(res)
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
func NewFileEventsFetcher(files []string) *FileEventsFetcher {
	return &FileEventsFetcher{files: files}
}
