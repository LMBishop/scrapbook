package html

import (
	"net/url"

	. "github.com/LMBishop/scrapbook/web/skeleton"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func AuthenticateSitePage(err, redirect, siteName string) Node {
	return Page("Authenticate",
		H1(Text("A password is required to visit this site")),

		If(err != "", AlertError(err)),

		Form(
			Action("/authenticate?redirect="+url.QueryEscape(redirect)),
			Method("post"),

			FieldSet(
				Legend(Text("Authentication")),
				Label(
					For("password"),
					Text("Password"),
				),
				Input(
					ID("password"),
					Name("password"),
					Type("password"),
				),
				Span(
					Class("form-help"),
					Text("Enter the password to continue."),
				),
			),

			Div(
				Class("control-group group-right"),
				Input(
					Type("submit"),
					Value("Submit"),
				),
			),
		),
	)
}
