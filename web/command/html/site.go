package html

import (
	"fmt"

	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/site"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func SitePage(mainConfig *config.MainConfig, site *site.Site) Node {
	versions, err := site.GetAllVersions()
	currentVersion, _ := site.GetCurrentVersion()

	return page("Site "+site.Name,
		H1(Text("Site "+site.Name)),

		If(site.EvaluateSiteStatus() != "live", alertError(site.EvaluateSiteStatusReason())),

		FieldSet(
			Legend(Text("Site actions")),

			Div(
				Class("control-group"),

				navButton("Upload new version", "upload"),
				navButton("Set flags", "flags"),
				navButton("Delete site", "delete"),
			),
		),

		H2(Text("Version history")),

		If(len(versions) == 0, Span(Class("span"), alert("There are no versions to display", ""))),
		If(err != nil, Span(Class("span"), alertError(fmt.Errorf("Cannot show site versions: %w", err).Error()))),
		If(len(versions) > 0 && err == nil, Group{
			Div(
				Class("versions-table"),
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
							If(currentVersion != version, navButton("Set current", fmt.Sprintf("/site/%s/", site.Name))),
							navButton("Details", fmt.Sprintf("version/%s/", version)),
						),
					}
				}),
			),
		}),

		H2(Text("Information")),

		P(Text("API endpoint for new versions: "), Code(Text(fmt.Sprintf("POST https://%s/api/site/%s/upload", mainConfig.Command.Host, site.Name)))),

		P(Text("Data directory on system: "), Code(Text(site.Path))),

		navButton("Go back", "/"),
	)
}
