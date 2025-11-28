package models

type DocumentSection struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"paragraphs"`
}

type Document struct {
	Title    string            `json:"title"`
	Author   string            `json:"author"`
	Sections []DocumentSection `json:"sections"`
}
