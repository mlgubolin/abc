package topic

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

const (
	MaxWorkReportTopicTitleLen = 300
	ErrorClosingZipFile        = "error closing zip file: %s"
	ErrorClosingFileReader     = "error closing file reader: %s"
)
const maxBytes = 50 << 20 // 50MB

type Topic struct {
	Title, Text string
}

func (t Topic) String() string {
	return fmt.Sprintf("=====================================\n%s\n=====================================\n\n%s\n", t.Title, t.Text)
}

func ExtractText(zr *zip.Reader) (string, []Topic, error) {
	mainDoc, err := getMainDocument(zr)
	if err != nil {
		return "", nil, err
	}
	return extractTexts(mainDoc)
}

func extractTexts(f *zip.File) (fullText string, topics []Topic, err error) {
	r, err := f.Open()
	if err != nil {
		return "", nil, err
	}
	defer func(r io.ReadCloser) {
		err := r.Close()
		if err != nil {
			_ = fmt.Errorf(ErrorClosingFileReader, err)
		}
	}(r)

	dec := xml.NewDecoder(io.LimitReader(r, maxBytes))
	inTopic := false
	var curr Topic

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", nil, err
		}

		switch v := tok.(type) {
		case xml.CharData:
			text := string(v)
			fullText += text
			if inTopic && curr.Title != "" {
				curr.Text += text
			}
		case xml.EndElement:
			if v.Name.Local == "p" && inTopic && curr.Title == "" {
				curr.Title = extractTitle(fullText)
			}
		case xml.StartElement:
			handleStartElement(v, &inTopic, &curr, &topics)
			if isLineBreak(v.Name.Local) {
				fullText += "\n"
				if inTopic && curr.Title != "" {
					curr.Text += "\n"
				}
			} else if v.Name.Local == "instrText" || v.Name.Local == "script" {
				if err := skipXML(dec); err != nil {
					return "", nil, err
				}
			}
		}
	}

	if inTopic && curr.Title != "" {
		topics = append(topics, curr)
	}
	return
}

func getMainDocument(zr *zip.Reader) (*zip.File, error) {
	for _, zf := range zr.File {
		if zf.Name == "word/document.xml" {
			return zf, nil
		}
	}
	return nil, fmt.Errorf("arquivo word/document.xml não encontrado")
}

func skipXML(dec *xml.Decoder) error {
	for depth := 1; depth > 0; {
		tok, err := dec.Token()
		if err != nil {
			return err
		}
		switch tok.(type) {
		case xml.StartElement:
			depth++
		case xml.EndElement:
			depth--
		}
	}
	return nil
}

func extractTitle(fullText string) string {
	start := strings.LastIndexByte(fullText, '\n') + 1
	if start > 0 {
		return strings.TrimSpace(fullText[start:])
	}
	return ""
}

func handleStartElement(el xml.StartElement, inTopic *bool, curr *Topic, topics *[]Topic) {
	if el.Name.Local == "pStyle" {
		for _, attr := range el.Attr {
			switch attr.Value {
			case "Ttulo3":
				if *inTopic && curr.Title != "" {
					*topics = append(*topics, *curr)
				}
				*inTopic = true
				*curr = Topic{}
			case "Ttulo1", "Ttulo2":
				if curr.Title != "" {
					*topics = append(*topics, *curr)
				}
				*curr = Topic{}
				*inTopic = false
			}
		}
	}
}

func isLineBreak(tag string) bool {
	return tag == "br" || tag == "p" || tag == "tab"
}
