package html

import (
	"fmt"

	. "github.com/LMBishop/scrapbook/web/skeleton"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func UploadPage(success, err string, siteName string) Node {
	return Page("Upload new version to "+siteName,
		H1(Text("Upload new version to "+siteName)),

		If(success != "", Group{
			AlertSuccess(success),
			Div(
				Class("control-group group-right"),
				NavButton("OK", fmt.Sprintf("/site/%s/", siteName)),
			),
		}),

		If(success == "", Group{
			If(err != "", AlertError(err)),

			Form(
				Method("post"),
				EncType("multipart/form-data"),

				FieldSet(
					Legend(Text("Upload")),
					Input(
						ID("upload"),
						Name("upload"),
						Type("file"),
					),
					Span(
						Class("form-help"),
						Text("A zip file with the contents of the site."),
					),
				),

				Div(
					Class("control-group group-right"),
					NavButton("Go back", fmt.Sprintf("/site/%s/", siteName)),
					Input(
						Type("submit"),
						Value("Submit"),
					),
				),
			),
		}),
	)
}
