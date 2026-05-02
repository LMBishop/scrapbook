package html

import (
	. "github.com/LMBishop/scrapbook/web/skeleton"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func SuccessPage(msg, nav string) Node {
	return Page("Success",
		H1(Text("Success")),

		AlertSuccess(msg),

		Div(
			Class("control-group group-right"),
			NavButton("Continue", nav),
		),
	)
}
