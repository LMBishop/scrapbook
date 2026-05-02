package html

import (
	"fmt"
	"time"

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

		func() Node {
			status, reason := site.EvaluateSiteStatus()
			if status != "live" {
				return AlertError(reason)
			} else {
				return Raw("")
			}
		}(),

		FieldSet(
			Legend(Text("Site actions")),

			Div(
				Class("control-group"),

				NavButton("Set flags", "flags"),
				NavButton("Change config", "config"),
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
						Text("Version"),
					),
					Span(
						Class("header date"),
						Text("Date"),
					),
					Span(
						Class("header actions"),
						Text("Actions"),
					),
				},

				Map(versions, func(version config.VersionMeta) Node {
					return Group{
						Span(
							Class("date"),
							Span(Text(version.Hash[:8])),
						),
						Span(
							Class("date"),
							Span(Text(time.Unix(int64(version.Created), 0).Format(time.DateTime))),
						),
						Span(
							Class("actions"),
							If(currentVersion == version.Hash, Span(Class("current"), Text("current"))),
							NavButton("Details", fmt.Sprintf("version/%s/", version.Hash)),
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

		P(Text("API endpoint for new versions: "), Code(Text(fmt.Sprintf("POST https://%s/api/site/%s/upload", mainConfig.Control.Host, site.Name)))),

		P(Text("Data directory on system: "), Code(Text(site.Path))),

		Br(),

		NavButton("Go back", "/"),
	)
}
