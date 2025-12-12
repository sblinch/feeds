package feeds

import (
	"fmt"
	"html"
	"io"
	"strings"
	"time"
)

// HTML is used to convert a generic Feed to HTML.
type HTML struct {
	*Feed
}

// ToHTML encodes f into a HTML string. Returns an error if marshalling fails.
func (f *HTML) ToHTML() (string, error) {
	b := strings.Builder{}
	if err := f.WriteHTML(&b); err == nil {
		return b.String(), nil
	} else {
		return "", err
	}
}

const (
	htmlIndent  = "    "
	htmlNewline = "\n"
)

type htmlWriter struct {
	w      io.Writer
	err    error
	indent int
}

func (w *htmlWriter) Err() error {
	return w.err
}

func (w *htmlWriter) String(s string) {
	if w.err != nil {
		return
	}
	_, w.err = io.WriteString(w.w, s)
}

func (w *htmlWriter) SafeString(s string) {
	w.String(html.EscapeString(s))
}

func (w *htmlWriter) Indent() {
	for i := 0; i < w.indent; i++ {
		if w.err != nil {
			return
		}
		_, w.err = io.WriteString(w.w, htmlIndent)
	}
}

func (w *htmlWriter) Line(s string) {
	w.Indent()
	w.String(s)
	if w.err == nil {
		_, w.err = io.WriteString(w.w, htmlNewline)
	}
}

func (w *htmlWriter) printf(format string, args ...interface{}) {
	for i, arg := range args {
		if s, ok := arg.(string); ok {
			args[i] = html.EscapeString(s)
		}
	}
	_, w.err = fmt.Fprintf(w.w, format, args...)
}

func (w *htmlWriter) Printf(format string, args ...interface{}) {
	w.Indent()
	if w.err == nil {
		w.printf(format, args...)
	}
	if w.err == nil {
		_, w.err = io.WriteString(w.w, htmlNewline)
	}
}

func (w *htmlWriter) WrapTag(tag string, f func(), attrPairs ...string) {
	w.Indent()
	w.OpenTag(tag, attrPairs...)
	w.String(htmlNewline)

	w.indent++
	f()
	w.indent--

	w.Indent()
	w.CloseTag(tag)
	w.String(htmlNewline)
}

func (w *htmlWriter) MaybeWrapTag(tag string, wrap bool, f func(), attrPairs ...string) {
	if wrap {
		w.WrapTag(tag, f, attrPairs...)
	} else {
		f()
	}
}

func (w *htmlWriter) OpenTag(name string, attrPairs ...string) {
	if w.err != nil {
		return
	}
	_, w.err = io.WriteString(w.w, "<")
	if w.err == nil {
		_, w.err = io.WriteString(w.w, name)
	}
	if len(attrPairs) > 0 {
		key := ""
		for i, v := range attrPairs {
			if i%2 == 0 {
				key = v
			} else if v != "" {
				if w.err == nil {
					_, w.err = io.WriteString(w.w, " ")
				}
				if w.err == nil {
					_, w.err = io.WriteString(w.w, key)
				}
				if w.err == nil {
					_, w.err = io.WriteString(w.w, `="`)
				}
				if w.err == nil {
					_, w.err = io.WriteString(w.w, html.EscapeString(v))
				}
				if w.err == nil {
					_, w.err = io.WriteString(w.w, `"`)
				}
			}
		}
	}
	if w.err == nil {
		_, w.err = io.WriteString(w.w, ">")
	}
}

func (w *htmlWriter) CloseTag(name string) {
	if w.err != nil {
		return
	}
	_, w.err = io.WriteString(w.w, "</")
	if w.err == nil {
		_, w.err = io.WriteString(w.w, name)
	}
	if w.err == nil {
		_, w.err = io.WriteString(w.w, ">")
	}
}

func (w *htmlWriter) StandaloneTag(name string, attrPairs ...string) {
	w.Indent()
	w.OpenTag(name, attrPairs...)
	if w.err == nil {
		_, w.err = io.WriteString(w.w, htmlNewline)
	}
}

func (w *htmlWriter) Tag(name string, value string, attrPairs ...string) {
	w.Indent()
	w.OpenTag(name, attrPairs...)
	if w.err == nil {
		_, w.err = io.WriteString(w.w, html.EscapeString(value))
	}
	w.CloseTag(name)
	if w.err == nil {
		_, w.err = io.WriteString(w.w, htmlNewline)
	}
}

func firstOf[T comparable](v ...T) T {
	var zero T
	for _, x := range v {
		if x != zero {
			return x
		}
	}
	return zero
}

func validEnclosure(encl *Enclosure) bool {
	return encl != nil && encl.Url != ""
}
func validAuthor(author *Author) bool {
	return author != nil && (author.Name != "" || author.Email != "")
}
func validLink(link *Link) bool {
	return link != nil && link.Href != ""
}
func validImage(image *Image) bool {
	return image != nil && image.Url != ""
}

func authorName(author *Author, combine bool) string {
	if author == nil {
		return ""
	}

	if author.Email != "" && author.Name != "" {
		if combine {
			return fmt.Sprint(author.Name, " (", author.Email, ")")
		} else {
			return author.Name
		}
	} else {
		return firstOf(author.Name, author.Email)
	}
}

func (f *HTML) WriteHTML(w io.Writer) error {
	sw := htmlWriter{w: w}
	sw.Line("<!doctype html>")
	sw.Line("<html>")
	sw.WrapTag("head", func() {
		if f.Title != "" {
			sw.Tag("title", f.Title)
		}

		if validLink(f.Link) {
			sw.StandaloneTag("link", "rel", firstOf(f.Link.Rel, "author"), "href", f.Link.Href)
		}

		if validAuthor(f.Author) {
			sw.StandaloneTag("meta", "name", "author", "value", authorName(f.Author, true))
		}

		if f.Description != "" {
			sw.StandaloneTag("meta", "name", "description", "value", f.Description)
		}
	})

	sw.WrapTag("body", func() {
		if validImage(f.Image) {
			sw.WrapTag("p", func() {
				sw.MaybeWrapTag(
					"a",
					f.Image.Link != "",
					func() {
						sw.StandaloneTag("img", "src", f.Image.Url, "title", f.Image.Title)
					},
					"href", f.Image.Link,
				)
			})
		}

		if f.Title != "" {
			sw.Tag("h1", f.Title)
		}
		if f.Subtitle != "" {
			sw.Tag("h2", f.Subtitle)
		}

		sw.WrapTag("ul", func() {
			for _, item := range f.Items {
				sw.WrapTag("li", func() {
					sw.WrapTag("p", func() {
						if item.Id != "" {
							sw.Tag("a", "name", item.Id)
						}

						itemTime := anyTimeFormat(time.RFC1123, item.Updated, item.Created)

						title := firstOf(item.Title, item.Id, itemTime)
						if title == itemTime {
							itemTime = ""
						}

						link := "#"
						if validLink(item.Link) {
							link = item.Link.Href
						}

						sw.Tag("a", title, "href", link)

						if itemTime != "" {
							sw.StandaloneTag("br")
							sw.Tag("small", itemTime)
						}
					})
					if validEnclosure(item.Enclosure) {
						sw.WrapTag("p", func() {
							sw.StandaloneTag("img", "src", item.Enclosure.Url)
						})
					}

					if item.Description != "" {
						itemHasContent := item.Content != ""
						descriptionHasPTag := strings.HasPrefix(item.Description, "<p>")

						sw.MaybeWrapTag("p", !descriptionHasPTag, func() {
							sw.MaybeWrapTag("em", itemHasContent, func() {
								// item.Description is intentionally not escaped as it seems intended to contain HTML
								sw.Line(item.Description)
							})
						})
					}

					if item.Content != "" {
						contentHasPTag := strings.HasPrefix(item.Content, "<p>")
						sw.MaybeWrapTag("p", !contentHasPTag, func() {
							// item.Content is intentionally not escaped as it seems intended to contain HTML
							sw.Line(item.Content)
						})
					}

					hasAuthor := validAuthor(item.Author)
					hasSource := validLink(item.Source)
					if hasAuthor || hasSource {
						sw.WrapTag("p", func() {
							sw.WrapTag("cite", func() {
								if hasAuthor {
									author := authorName(item.Author, false)
									if item.Author.Email != "" {
										sw.Tag("a", author, "href", item.Author.Email)
									} else {
										sw.Printf("%s", author)
									}
								}
								if hasSource {
									// avoid using words to eliminate the need for i18n
									sw.Tag("a", " (↗)", "href", item.Source.Href)
								}
								sw.StandaloneTag("br")
							})
						})
					}

				})
			}
		})

		sw.WrapTag("p", func() {
			if f.Copyright != "" {
				sw.Printf("%s", f.Copyright)
			}
			if t := anyTimeFormat(time.RFC1123, f.Updated, f.Created); t != "" {
				sw.StandaloneTag("br")
				sw.Tag("small", t)
			}
		})

	})
	sw.Line("</html>")

	return sw.Err()
}
