package html

import (
	"fmt"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func HostPage(success, err, siteName, hostValue string) Node {
	return page("Change host for "+siteName,
		H1(Text("Change host for "+siteName)),

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

				FieldSet(
					Legend(Text("Host")),
					Input(
						ID("host"),
						Name("host"),
						Value(hostValue),
					),
					Span(
						Class("form-help"),
						Text("The fully qualified domain name for which this site is to be served on. If this site is not to be served by scrapbook, leave blank."),
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
