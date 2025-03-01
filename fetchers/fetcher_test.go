package fetchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewsParsing(t *testing.T) {
	fetcher := NewFileNewsFetcher([]string{"../testdata/news.json"})
	data, err := fetcher.Fetch()
	assert.Nil(t, err)
	assert.NotNil(t, data)
}
func TestEventsParsing(t *testing.T) {
	fetcher := NewFileEventsFetcher([]string{"../testdata/events.json"})
	data, err := fetcher.Fetch()
	assert.Nil(t, err)
	assert.NotNil(t, data)
}
