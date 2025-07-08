package html

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type CreatePageForm struct {
	Name string
	Host string
}

func CreatePage(success, err string, formValues CreatePageForm) Node {
	return page("Create site",
		H1(Text("Create site")),

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
				Action("/create"),
				Method("post"),

				FieldSet(
					Legend(Text("Site details")),
					Label(
						For("name"),
						Text("Name"),
					),
					Input(
						ID("name"),
						Name("name"),
						Value(formValues.Name),
					),
					Span(
						Class("form-help"),
						Text("The unique identifier for this site. This must be a valid directory name, and should be lower case with no spaces."),
					),

					Label(
						For("host"),
						Text("Host"),
					),
					Input(
						ID("host"),
						Name("host"),
						Value(formValues.Host),
					),
					Span(
						Class("form-help"),
						Text("The fully qualified domain name for which this site is to be served on."),
					),
				),

				Div(
					Class("control-group group-right"),
					navButton("Go back", "/"),
					Input(
						Type("submit"),
						Value("Submit"),
					),
				),
			),
		}),
	)
}
