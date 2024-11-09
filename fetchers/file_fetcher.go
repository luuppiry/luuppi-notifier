package fetchers

type FileFetcher struct{}

func (f *FileFetcher) Fetch(string) []byte {
	return []byte{}
}
func NewFileFetcher() *FileFetcher {
	return &FileFetcher{}
}
