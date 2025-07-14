package html

import (
	"fmt"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func DeletePage(success, err string, siteName string) Node {
	return page("Delete "+siteName,
		H1(Text("Delete "+siteName)),

		If(success != "", Group{
			alertSuccess(success),
			Div(
				Class("control-group group-right"),
				navButton("OK", "/"),
			),
		}),

		If(success == "", Group{
			If(err != "", alertError(err)),

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
