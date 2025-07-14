package html

import (
	"fmt"

	"github.com/LMBishop/scrapbook/pkg/index"
	"github.com/LMBishop/scrapbook/pkg/site"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func HomePage(siteIndex *index.SiteIndex) Node {
	return page("All sites",
		H1(Text("All sites")),

		Div(
			Class("sites-table"),
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
				return Group{
					Span(
						Class("name"),
						Span(Text(site.Name)),
						If(site.SiteConfig.Host == "", Span(Text("no host"))),
						If(site.SiteConfig.Host != "", Span(Text(fmt.Sprintf("on %s", site.SiteConfig.Host)))),
					),
					Span(
						Class("status"),
						Text(site.EvaluateSiteStatus()),
					),
					Span(
						Class("flags"),
						Text(site.ConvertFlagsToString()),
					),
					Span(
						Class("actions"),
						navButton("Details", fmt.Sprintf("/site/%s/", site.Name)),
					),
				}
			}),
		),

		navButton("Create new", "/create"),
	)
}
