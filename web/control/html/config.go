package html

import (
	"fmt"

	. "github.com/LMBishop/scrapbook/web/skeleton"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func ConfigPage(success, err, siteName, configValue string) Node {
	return Page("Change configuration for "+siteName,
		H1(Text("Change configuration for "+siteName)),

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

				P(Text("This configuration file defines how scrapbook handles this site. To change the behaviour of the internal web server, set site flags instead. This configuration file is saved at the root of the site data directory.")),

				FieldSet(
					Legend(Text("Configuration")),
					Textarea(
						ID("config"),
						Name("config"),
						Rows("15"),
						Text(configValue),
					),
					Span(
						Class("form-help"),
						Text("See below for a configuration reference."),
					),
				),

				Div(
					Class("control-group group-gap margin-bottom"),
					NavButton("Go back", fmt.Sprintf("/site/%s/", siteName)),
					Input(
						Type("submit"),
						Value("Submit"),
					),
				),

				FieldSet(
					Legend(Text("Reference")),
					Text("Yadda yadda"),
				),
			),
		}),
	)
}
