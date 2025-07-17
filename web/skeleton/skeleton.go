package skeleton

import (
	_ "embed"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

//go:embed style.css
var styles string

func Page(title string, children ...Node) Node {
	return HTML5(HTML5Props{
		Title:    title,
		Language: "en",
		Head: []Node{
			StyleEl(Raw(styles)),
		},
		Body: []Node{
			Div(Class("container"),
				Group(children),
				footer(),
			),
		},
	})
}

func footer() Node {
	return Footer(
		Hr(),
		Text("scrapbook"),
	)
}

func NavButton(label string, dest string) Node {
	return A(
		Class("button"),
		Href(dest),
		Text(label),
	)
}

func Alert(label string, class string) Node {
	return Div(
		Class("alert "+class),
		Text(label),
	)
}

func AlertError(label string) Node {
	return Alert(label, "error")
}

func AlertSuccess(label string) Node {
	return Alert(label, "success")
}
