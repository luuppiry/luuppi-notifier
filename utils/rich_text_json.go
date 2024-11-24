package utils

type RTF []RTF_piece

type RTF_piece struct {
	Type     string
	Children []RTF_piece
	Url      string
	Text     string
}

func ParseRTFJson(data RTF) string {
	return ""
}
