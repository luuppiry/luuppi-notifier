package output

import "github.com/luuppiry/luuppi-rss-service/types"

type Formattable interface {
	Rss_format(basePath string) (types.Rss_item, error)
	Discord_format() (types.Discord_message, error)
}
