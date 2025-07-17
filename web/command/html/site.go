package html

import (
	"fmt"

	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/site"
	. "github.com/LMBishop/scrapbook/web/skeleton"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func SitePage(mainConfig *config.MainConfig, site *site.Site) Node {
	versions, err := site.GetAllVersions()
	currentVersion, _ := site.GetCurrentVersion()

	return Page("Site "+site.Name,
		H1(Text("Site "+site.Name)),

		If(site.EvaluateSiteStatus() != "live", AlertError(site.EvaluateSiteStatusReason())),

		FieldSet(
			Legend(Text("Site actions")),

			Div(
				Class("control-group"),

				NavButton("Change host", "host"),
				NavButton("Set flags", "flags"),
				NavButton("Delete site", "delete"),
			),
		),

		H2(Text("Version history")),

		If(len(versions) == 0, Span(Class("span"), Alert("There are no versions to display", ""))),
		If(err != nil, Span(Class("span"), AlertError(fmt.Errorf("Cannot show site versions: %w", err).Error()))),
		If(len(versions) > 0 && err == nil, Group{
			Div(
				Class("table versions-table"),
				Group{
					Span(
						Class("header date"),
						Text("Date"),
					),
					Span(
						Class("header actions"),
						Text("Actions"),
					),
				},

				Map(versions, func(version string) Node {
					return Group{
						Span(
							Class("date"),
							Span(Text(version)),
							If(currentVersion == version, Span(Class("current"), Text("current"))),
						),
						Span(
							Class("actions"),
							If(currentVersion != version, NavButton("Set current", fmt.Sprintf("/site/%s/", site.Name))),
							NavButton("Details", fmt.Sprintf("version/%s/", version)),
						),
					}
				}),
			),
		}),
		Div(
			Class("control-group group-right"),

			NavButton("Upload new version", "upload"),
		),

		H2(Text("Information")),

		P(Text("API endpoint for new versions: "), Code(Text(fmt.Sprintf("POST https://%s/api/site/%s/upload", mainConfig.Command.Host, site.Name)))),

		P(Text("Data directory on system: "), Code(Text(site.Path))),

		Br(),

		NavButton("Go back", "/"),
	)
}
