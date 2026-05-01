package html

import (
	. "github.com/LMBishop/scrapbook/web/skeleton"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type CreatePageForm struct {
	Name string
	Host string
}

func CreatePage(success, err string, formValues CreatePageForm) Node {
	return Page("Create site",
		H1(Text("Create site")),

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
						Text("The fully qualified domain name for which this site is to be served on. If this site is not to be served by scrapbook, leave blank."),
					),
				),

				Div(
					Class("control-group group-right"),
					NavButton("Go back", "/"),
					Input(
						Type("submit"),
						Value("Submit"),
					),
				),
			),
		}),
	)
}
