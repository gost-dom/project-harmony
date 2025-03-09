package gviews

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Index() Node {
	return Layout(
		Div(Class("container mx-auto py-4"),
			H1(Class("text-4xl font-bold text-center py-4"),
				Text("Project Harmony")),
			P(Text("This is just a demo")),
		),
	)
}
