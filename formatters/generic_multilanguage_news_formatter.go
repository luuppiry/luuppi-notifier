package formatters

import (
	"encoding/xml"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/luuppiry/luuppi-rss-service/types"
)

type GenericMultiLanguageFormatter struct {
}

func (f *GenericMultiLanguageFormatter) Format(n any) ([]byte, error) {
	news := n.([]types.NormalizedNews)
	correlated := correlateNews(news)
	return out, nil

}

func correlateNews(data []types.NormalizedNews) [][]types.NormalizedNews {
	ret := [][]types.NormalizedNews{}
	cutoff := len(data)
	for i := 0; i < cutoff; i++ {
		n := data[i]
		c := closure(n.Published)
		same := []types.NormalizedNews{}
		for true {
			matching := slices.IndexFunc(data[i:cutoff], c)
			if matching == -1 {
				break
			}
			same = append(same, data[matching])
			data[matching], data[cutoff] = data[cutoff], data[matching]
			cutoff--
		}
		ret = append(ret, same)
	}
	return ret
}

func closure(a time.Time) func(types.NormalizedNews) bool {
	return func(nn types.NormalizedNews) bool {
		return isCloseInTime(a, nn.Published)
	}
}

func isCloseInTime(a, b time.Time) bool {
	return a.Truncate(time.Minute).Equal(b.Truncate(time.Minute))
}

func NewGenericMultiLanguageFormatter(conf map[string]string) *GenericMultiLanguageFormatter {
	return &GenericMultiLanguageFormatter{}
}
