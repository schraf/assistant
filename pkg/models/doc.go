package models

import "github.com/schraf/assistant/internal/content"

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
