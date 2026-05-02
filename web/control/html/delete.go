package html

import (
	"fmt"

	. "github.com/LMBishop/scrapbook/web/skeleton"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func DeletePage(success, err, what string) Node {
	return Page("Delete "+what,
		H1(Text("Delete "+what)),

		If(success != "", Group{
			AlertSuccess(success),
			Div(
				Class("control-group group-right"),
				NavButton("OK", "../.."),
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
							Text(fmt.Sprintf("Really delete %s?", what)),
						),
					),
					Span(
						Class("form-help"),
						Text("Check the box to confirm deletion. All data on disk will be deleted. This action is irreversible."),
					),
				),

				Div(
					Class("control-group group-gap"),
					NavButton("Go back", "."),
					Input(
						Type("submit"),
						Value("Submit"),
					),
				),
			),
		}),
	)
}
