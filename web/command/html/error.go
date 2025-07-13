package html

import (
	. "maragu.dev/gomponents"
)

func ErrorPage(err string) Node {
	return page("Error",
		alertError(err),

		navButton("Home", "/"),
	)
}
