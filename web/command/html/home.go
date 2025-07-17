package html

import (
	"fmt"

	"github.com/LMBishop/scrapbook/pkg/index"
	"github.com/LMBishop/scrapbook/pkg/site"
	. "github.com/LMBishop/scrapbook/web/skeleton"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func HomePage(siteIndex *index.SiteIndex) Node {
	return Page("All sites",
		H1(Text("All sites")),

		Div(
			Class("table sites-table"),
			Group{
				Span(
					Class("header name"),
					Text("Site"),
				),
				Span(
					Class("header status"),
					Text("Status"),
				),
				Span(
					Class("header flags"),
					Text("Flags"),
				),
				Span(
					Class("header actions"),
					Text("Actions"),
				),
			},

			Map(siteIndex.GetSites(), func(site *site.Site) Node {
				status := site.EvaluateSiteStatus()
				good := false
				if status == "live" {
					good = true
				}
				return Group{
					Span(
						Class("name"),
						Span(Text(site.Name)),
						If(site.SiteConfig.Host == "", Span(Text("no host"))),
						If(site.SiteConfig.Host != "", Span(Text(fmt.Sprintf("on %s", site.SiteConfig.Host)))),
					),
					Span(
						If(good, Class("status text-good")),
						If(!good, Class("status text-bad")),
						Text(status),
					),
					Span(
						Class("flags"),
						Text(site.ConvertFlagsToString()),
					),
					Span(
						Class("actions"),
						NavButton("Details", fmt.Sprintf("/site/%s/", site.Name)),
					),
				}
			}),
		),

		If(len(siteIndex.GetSites()) == 0, Alert("There are no sites to display", "")),

		Div(
			Class("control-group group-right"),

			NavButton("Create new", "/create"),
		),
	)
}
