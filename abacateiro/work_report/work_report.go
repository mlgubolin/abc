package workreport

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

const maxBytes = 50 << 20 // 50MB

type Topic struct {
	Title string
	Text  string
}

func (t Topic) String() string {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, "=====================================")
	fmt.Fprintln(&buf, t.Title)
	fmt.Fprintln(&buf, "=====================================")
	fmt.Fprintln(&buf)
	fmt.Fprintln(&buf, t.Text)

	return buf.String()
}

func ExtractText(zr *zip.Reader) (text string, topics []Topic, err error) {
	mainDoc, err := getMainDocument(zr)
	if err != nil {
		return "", nil, err
	}

	text, topics, err = extractTexts(mainDoc)

	return
}

func extractTexts(f *zip.File) (fullText string, topics []Topic, err error) {
	r, err := f.Open()
	if err != nil {
		return "", nil, err
	}
	defer r.Close()

	dec := xml.NewDecoder(io.LimitReader(r, maxBytes))
	dec.Strict = true
	inTopic := false

	var curr Topic

	for {
		t, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", nil, err
		}

		switch v := t.(type) {
		case xml.CharData:
			fullText += string(v)
			if inTopic && len(curr.Title) > 0 {
				curr.Text += string(v)
			}

		case xml.EndElement:
			if v.Name.Local == "p" && inTopic && len(curr.Title) == 0 {
				startTitle := strings.LastIndexByte(fullText, '\n') + 1
				if startTitle > 0 {
					curr.Title = strings.TrimSpace(fullText[startTitle:])
				}
				// TODO: o que fazer se o titulo nao for identificado?
			}

		case xml.StartElement:
			if v.Name.Local == "pStyle" && len(v.Attr) > 0 {
				for _, attr := range v.Attr {
					switch attr.Value {
					case "Ttulo3":
						// transição de um tópico para outro
						if inTopic && len(curr.Title) > 0 {
							topics = append(topics, curr)
						}
						inTopic = true
						curr = Topic{}
					case "Ttulo1", "Ttulo2":
						// transição para um outro título sem ser tópico
						if len(curr.Title) > 0 {
							topics = append(topics, curr)
						}
						curr = Topic{}
						inTopic = false
					}
				}
			}

			if (v.Name.Local == "br") || (v.Name.Local == "p") || (v.Name.Local == "tab") {
				fullText += "\n"
				if inTopic && len(curr.Title) > 0 {
					curr.Text += "\n"
				}
			} else if (v.Name.Local == "instrText") || (v.Name.Local == "script") {
				// pula todo o corpo entre tags <instrText></instrText> e <script></script>
				if err := skipXML(dec); err != nil {
					return "", nil, err
				}
			}
		}
	}

	if inTopic && len(curr.Title) > 0 {
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
	depth := 1
	for {
		t, err := dec.Token()
		if err != nil {
			return err
		}

		switch t.(type) {
		case xml.StartElement:
			depth++
		case xml.EndElement:
			depth--
		}

		if depth == 0 {
			break
		}
	}
	return nil
}
