package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tdewolff/parse/html"
	"github.com/tdewolff/parse/js"
)

func processMarkdownActions(input string) string {
	var builder strings.Builder
	lines := strings.Split(input, "\n")
	inAction := false
	for i := range lines {
		line := lines[i]
		if len(line) > 0 {
			if line[0] == '|' {
				if !inAction {
					builder.WriteString("\n<section class=\"highlighted\"><i class=\"fas fa-wrench\"></i>\n<div>\n")
					inAction = true
				}

				if len(line) == 1 || line[1] != ' ' {
					builder.WriteString(line[1:])
				} else {
					builder.WriteString(line[2:])
				}
			} else {
				if inAction {
					builder.WriteString("</div>\n</section>\n\n")
					inAction = false
				}
				builder.WriteString(line)
			}
		} else {
			if inAction {
				builder.WriteString("</div>\n</section>\n\n")
				inAction = false
			}
		}
		builder.WriteRune('\n')
	}
	return builder.String()
}

type page struct {
	index                int
	name, title, content string
	url, prevURL         string
	linkedfile           string
}

func (p *page) write(nextURL string) {
	fName := fmt.Sprintf("%02d-%s.md", p.index, p.name)
	f, err := os.Create(fName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	writer := bufio.NewWriter(f)
	fmt.Fprintf(writer, "---\nlayout: tutorialpage\nweight: %d\ntitle: %s\npermalink: %s\n",
		p.index, p.title, p.url)
	if len(p.linkedfile) > 0 {
		fmt.Fprintf(writer, "linkedfile: %s\n", p.linkedfile)
	}
	fmt.Fprintf(writer, "previous: %s\nnext: %s\n---\n",
		p.prevURL, nextURL)
	writer.WriteString(p.content)
	writer.Flush()
}

type pageWriter struct {
	pending page
}

func (pw *pageWriter) put(name string, title string, content string, linkedfile string) {
	var url string
	if name == "index" {
		url = "/plugins/tutorial/"
	} else {
		url = "/plugins/tutorial/" + name
	}
	if pw.pending.index != 0 {
		pw.pending.write(url)
	}
	pw.pending = page{index: pw.pending.index + 1,
		name: name, title: title, content: processMarkdownActions(content),
		linkedfile: linkedfile,
		url:        url, prevURL: pw.pending.url}
}

func (pw *pageWriter) processMdFile(title string, name string) {
	content, err := ioutil.ReadFile("../plugin-tutorial/" + name + ".md")
	if err != nil {
		panic(err)
	}
	pw.put(name, title, string(content), "")
}

func (pw *pageWriter) processGoFile(path string) {
	fileset := token.NewFileSet()
	contents, err := ioutil.ReadFile("../../plugin-tutorial/" + path)
	if err != nil {
		panic(err)
	}
	file, err := parser.ParseFile(fileset, path, contents, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	var builder strings.Builder
	curOffset := 0
	postponedEnd := 0
	title := ""

	for i := range file.Comments {
		group := file.Comments[i]
		if group.List[0].Text[:2] == "//" {
			continue
		}
		gOffset := fileset.Position(group.Pos()).Offset
		if curOffset == 0 {
			postponedEnd = gOffset
		} else if gOffset > curOffset {
			builder.WriteString("\n---\n\n```go\n")
			if postponedEnd != 0 {
				builder.WriteString(strings.TrimSpace(string(contents[0:postponedEnd])))
				builder.WriteString("\n\n")
				postponedEnd = 0
			}
			builder.WriteString(strings.TrimSpace(string(contents[curOffset:gOffset])))
			builder.WriteString("\n```\n\n")
		}
		text := group.Text()
		index := strings.IndexByte(text, '\n')
		if strings.HasPrefix(text, "title: ") {
			title = text[7:index]
			text = text[index+1:]
		}
		indent := 0
		for text[indent] == '\t' {
			indent++
		}
		lines := strings.Split(text, "\n")
		for i := range lines {
			if len(lines[i]) > indent {
				builder.WriteString(lines[i][indent:len(lines[i])])
			}
			builder.WriteByte('\n')
		}
		curOffset = fileset.Position(group.End()).Offset
	}
	filename := filepath.Base(path)
	pw.put(strings.TrimSuffix(filename, filepath.Ext(filename)),
		title, builder.String(), path)
}

func (pw *pageWriter) processHTMLFile(path string) {
	contents, err := ioutil.ReadFile("../../plugin-tutorial/" + path)
	if err != nil {
		panic(err)
	}
	lex := html.NewLexer(bytes.NewReader(contents))

	var builder strings.Builder
	title := ""
	mark := 0
	curOffset := 0

lexloop:
	for {
		tt, data := lex.Next()
		switch tt {
		case html.ErrorToken:
			break lexloop
		case html.CommentToken:
			break
		default:
			mark += len(data)
			continue
		}
		if mark != curOffset {
			builder.WriteString("\n---\n\n```html\n")
			builder.WriteString(strings.TrimSpace(string(contents[curOffset:mark])))
			builder.WriteString("\n```\n\n")
		}
		text := string(data)
		if strings.HasPrefix(text, "<!--title: ") {
			index := strings.IndexByte(text, '\n')
			title = text[11:index]
			text = text[index+1 : len(text)-4]
		} else {
			text = text[4 : len(text)-4]
		}
		builder.WriteString(text)
		mark += len(data)
		curOffset = mark
	}
	filename := filepath.Base(path)
	pw.put(strings.TrimSuffix(filename, filepath.Ext(filename)),
		title, builder.String(), path)
}

func (pw *pageWriter) processJSFile(path string) {
	contents, err := ioutil.ReadFile("../../plugin-tutorial/" + path)
	if err != nil {
		panic(err)
	}
	lex := js.NewLexer(bytes.NewReader(contents))

	var builder strings.Builder
	title := ""
	mark := 0
	curOffset := 0

lexloop:
	for {
		tt, data := lex.Next()
		switch tt {
		case js.ErrorToken:
			break lexloop
		case js.MultiLineCommentToken:
			break
		default:
			mark += len(data)
			continue
		}
		if mark != curOffset {
			builder.WriteString("\n---\n\n```javascript\n")
			builder.WriteString(strings.TrimSpace(string(contents[curOffset:mark])))
			builder.WriteString("\n```\n\n")
		}
		text := string(data)
		if strings.HasPrefix(text, "/*title: ") {
			index := strings.IndexByte(text, '\n')
			title = text[9:index]
			text = text[index+1 : len(text)-2]
		} else {
			text = text[2 : len(text)-2]
		}
		builder.WriteString(text)
		mark += len(data)
		curOffset = mark
	}
	filename := filepath.Base(path)
	pw.put(strings.TrimSuffix(filename, filepath.Ext(filename)),
		title, builder.String(), path)
}

func (pw *pageWriter) Close() {
	if pw.pending.index != 0 {
		pw.pending.write("")
		pw.pending.index = 0
	}
}

func main() {
	var pw pageWriter
	defer pw.Close()
	pw.processMdFile("Plugin Tutorial", "index")
	pw.processMdFile("Introduction", "introduction")
	pw.processGoFile("calendar/universitydate.go")
	pw.processGoFile("calendar/state.go")
	pw.processGoFile("calendar/renderer.go")
	pw.processGoFile("calendar/descriptor.go")
	pw.processHTMLFile("web/html/templates.html")
	pw.processJSFile("web/js/controllers.js")
	pw.processGoFile("plugin.go")
	pw.processMdFile("Building the Plugin", "building")
}
