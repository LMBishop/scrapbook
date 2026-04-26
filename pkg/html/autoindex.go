package html

import (
	"time"

	. "github.com/LMBishop/scrapbook/web/skeleton"
	"github.com/dustin/go-humanize"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type File struct {
	Name  string
	IsDir bool
	Size  int64
	Mtime time.Time
}

func IndexPage(dir string, err bool, files []File) Node {
	return Page("Index of "+dir,
		H1(Text("Index of "+dir)),

		Div(
			Class("table files-table"),
			Group{
				Span(
					Class("header name"),
					Text("Name"),
				),
				Span(
					Class("header size"),
					Text("Size"),
				),
				Span(
					Class("header last-modified"),
					Text("Last modified"),
				),
			},

			If(files != nil, Map(files, func(file File) Node {
				var fileName string
				if file.IsDir {
					fileName = file.Name + "/"
				} else {
					fileName = file.Name
				}
				return Group{
					A(
						Class("name"),
						Href(fileName),
						Text(fileName),
					),
					Span(
						Class("size"),
						If(file.IsDir, Text("--")),
						If(!file.IsDir, Text(humanize.IBytes(uint64(file.Size)))),
					),
					Span(
						Class("last-modified"),
						If(file.IsDir, Text("--")),
						If(!file.IsDir, Text(file.Mtime.Format("2006-01-02 15:04:05"))),
					),
				}
			})),
		),

		If(err, AlertError("Failed to list directory")),
	)
}
