package html

import (
	_ "embed"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

//go:embed style.css
var styles string

func page(title string, children ...Node) Node {
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

func navButton(label string, dest string) Node {
	return A(
		Href(dest),
		Text(label),
	)
}

func alert(label string, class string) Node {
	return Div(
		Class("alert "+class),
		Text(label),
	)
}

func alertError(label string) Node {
	return alert(label, "error")
}

func alertSuccess(label string) Node {
	return alert(label, "success")
}
