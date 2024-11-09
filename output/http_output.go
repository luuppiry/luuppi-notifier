package output

import "net/http"

type Http_output struct {
	slug string
	data []byte
}

func (h *Http_output) Initialize() error {
	http.Handle(h.slug, h)
	return nil
}

func (h *Http_output) Update(data []byte) error {
	h.data = data
	return nil
}

func (h *Http_output) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write(h.data)
}

func NewHttpOutput(conf map[string]string) *Http_output {
	return &Http_output{
		slug: conf["slug"],
	}
}
