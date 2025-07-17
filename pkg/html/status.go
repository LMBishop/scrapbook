package html

import (
	"fmt"

	. "github.com/LMBishop/scrapbook/web/skeleton"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func NotFoundUrlPage(url, host string) Node {
	return Page("Page not found",
		H1(Text("Page not found")),

		P(Text(fmt.Sprintf("The URL %s could not be found on site %s", url, host))),
	)
}

func NotFoundSitePage(host string) Node {
	return Page("Site not found",
		H1(Text("Site not found")),

		P(Text(fmt.Sprintf("The site %s is unknown", host))),
	)
}

func ForbiddenDisabledPage(host string) Node {
	return Page("Forbidden",
		H1(Text("Forbidden")),

		P(Text(fmt.Sprintf("Site %s is disabled", host))),
	)
}
