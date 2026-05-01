package html

import (
	"fmt"

	. "github.com/LMBishop/scrapbook/web/skeleton"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func DeletePage(success, err string, siteName string) Node {
	return Page("Delete "+siteName,
		H1(Text("Delete "+siteName)),

		If(success != "", Group{
			AlertSuccess(success),
			Div(
				Class("control-group group-right"),
				NavButton("OK", "/"),
			),
		}),

		If(success == "", Group{
			If(err != "", AlertError(err)),

			Form(
				Method("post"),

				FieldSet(
					Legend(Text("Delete")),
					Span(
						Input(
							ID("delete"),
							Name("delete"),
							Type("checkbox"),
						),
						Label(
							For("delete"),
							Text(fmt.Sprintf("Really delete site %s?", siteName)),
						),
					),
					Span(
						Class("form-help"),
						Text("Check the box to confirm deletion. Data on disk (including all site versions) will be deleted. This action is irreversible."),
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
