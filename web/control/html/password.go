package html

import (
	"fmt"

	. "github.com/LMBishop/scrapbook/web/skeleton"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func PasswordPage(success, err, siteName, passwordValue string) Node {
	return Page("Change password for "+siteName,
		H1(Text("Change password for "+siteName)),

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

				FieldSet(
					Legend(Text("Password")),
					Input(
						ID("password"),
						Name("password"),
						Value(passwordValue),
					),
					Span(
						Class("form-help"),
						Text("The password visitors must enter to be able to view the site."),
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
