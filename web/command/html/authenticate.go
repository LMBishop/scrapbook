package html

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func AuthenticatePage(err string) Node {
	return page("Authenticate",
		H1(Text("Welcome to scrapbook")),

		If(err != "", alertError(err)),

		Form(
			Action("/authenticate"),
			Method("post"),

			FieldSet(
				Legend(Text("Authentication")),
				Label(
					For("token"),
					Text("Secret key"),
				),
				Input(
					ID("token"),
					Name("token"),
					Type("password"),
				),
				Span(
					Class("form-help"),
					Text("Enter the secret key to continue."),
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
