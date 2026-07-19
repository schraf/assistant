package models

import (
	"github.com/schraf/assistant/internal/content"
)

type DocumentSection struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"paragraphs"`
}

type Document struct {
	Title    string            `json:"title"`
	Author   string            `json:"author"`
	Sections []DocumentSection `json:"sections"`
}

func (d Document) Length() int {
	length := len(d.Title) + len(d.Author)

	for _, section := range d.Sections {
		length += len(section.Title)

		for _, paragraph := range section.Paragraphs {
			length += len(paragraph)
		}
	}

	return length
}

func (d *Document) AddSection(title string, body string) *DocumentSection {
	index := len(d.Sections)

	d.Sections = append(d.Sections, DocumentSection{
		Title:      title,
		Paragraphs: content.SplitParagraphs(body),
	})

	return &d.Sections[index]
}

func (d *Document) Clean() {
	d.Title = content.CleanText(d.Title)
	d.Author = content.CleanText(d.Author)

	for i := range d.Sections {
		d.Sections[i].Title = content.CleanText(d.Sections[i].Title)
		for j := range d.Sections[i].Paragraphs {
			d.Sections[i].Paragraphs[j] = content.CleanText(d.Sections[i].Paragraphs[j])
		}
	}
}

func DocumentSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"title": map[string]any{
				"type":        "string",
				"description": "The comprehensive title of the macro research report.",
			},
			"sections": map[string]any{
				"type":        "array",
				"description": "The structural body chapters detailing deep subtopic discoveries.",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"title": map[string]any{
							"type":        "string",
							"description": "The chapter or subtopic title header.",
						},
						"paragraphs": map[string]any{
							"type":        "array",
							"description": "The list of paragraphs for a single section",
							"items": map[string]any{
								"type":        "string",
								"description": "An exhaustive, highly detailed narrative paragraph containing at least 250-400 words of granular analysis, data points, and explicit facts. Do NOT summarize or condense the researcher's raw data.",
							},
						},
					},
					"required": []string{"title", "paragraphs"},
				},
			},
		},
		"required": []string{"title", "sections"},
	}
}
