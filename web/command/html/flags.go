package html

import (
	"fmt"

	"github.com/LMBishop/scrapbook/pkg/config"
	. "github.com/LMBishop/scrapbook/web/skeleton"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func FlagsPage(success, err string, siteName string, flags config.SiteFlag) Node {
	return Page("Set flags for "+siteName,
		H1(Text("Set flags for "+siteName)),

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

				P(Text("These flags affect the behaviour of scrapbook's internal web server. They will have no effect if you are serving the site using a different web server.")),

				FieldSet(
					Legend(Text("Flags")),
					Span(
						Input(
							ID("disable"),
							Name("disable"),
							Type("checkbox"),
							If(flags&config.FlagDisable != 0, Checked()),
						),
						Label(
							For("disable"),
							Text("Disable"),
						),
					),
					Span(
						Class("form-help"),
						Text("Disallow access to this site."),
					),

					Span(
						Input(
							ID("tls"),
							Name("tls"),
							Type("checkbox"),
							Disabled(),
							If(flags&config.FlagTLS != 0, Checked()),
						),
						Label(
							For("tls"),
							Text("TLS"),
						),
					),
					Span(
						Class("form-help"),
						Text("Serve this site on the HTTPS socket."),
					),

					Span(
						Input(
							ID("index"),
							Name("index"),
							Type("checkbox"),
							If(flags&config.FlagIndex != 0, Checked()),
						),
						Label(
							For("index"),
							Text("Automatic index"),
						),
					),
					Span(
						Class("form-help"),
						Text("Generate index.html files on the fly if they do not exist."),
					),

					Span(
						Input(
							ID("password"),
							Name("password"),
							Type("checkbox"),
							Disabled(),
							If(flags&config.FlagPassword != 0, Checked()),
						),
						Label(
							For("password"),
							Text("Password protect"),
						),
					),
					Span(
						Class("form-help"),
						Text("Require visitors to enter a password to view the site."),
					),

					Span(
						Input(
							ID("readonly"),
							Name("readonly"),
							Type("checkbox"),
							If(flags&config.FlagReadOnly != 0, Checked()),
						),
						Label(
							For("readonly"),
							Text("Read only"),
						),
					),
					Span(
						Class("form-help"),
						Text("Disallow new site revisions or modification."),
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
