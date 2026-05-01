package html

import (
	. "github.com/LMBishop/scrapbook/web/skeleton"
	. "maragu.dev/gomponents"
)

func ErrorPage(err string) Node {
	return Page("Error",
		AlertError(err),

		NavButton("Home", "/"),
	)
}
