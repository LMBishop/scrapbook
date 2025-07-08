package html

import (
	"fmt"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func UploadPage(success, err string, siteName string) Node {
	return page("Upload new version to "+siteName,
		H1(Text("Upload new version to "+siteName)),

		If(success != "", Group{
			alertSuccess(success),
			Div(
				Class("control-group group-right"),
				navButton("OK", fmt.Sprintf("/site/%s/", siteName)),
			),
		}),

		If(success == "", Group{
			If(err != "", alertError(err)),

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
					navButton("Go back", fmt.Sprintf("/site/%s/", siteName)),
					Input(
						Type("submit"),
						Value("Submit"),
					),
				),
			),
		}),
	)
}
