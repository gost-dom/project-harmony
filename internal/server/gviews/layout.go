package gviews

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"

	hx "maragu.dev/gomponents-htmx"
)

func Layout(children ...Node) Node {
	return Doctype(HTML(
		Head(
			Link(Rel("stylesheet"), Href("/static/css/tailwind.css")),
			Script(Src("/static/js/htmx.js")),
			Meta(Charset("utf-8")),
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
			Title("Project Harmony"),
		),
		Body(Class("bg-secondary-50 dark:bg-stone-800"),
			Div(
				Class("min-h-screen flex flex-col"),
				Header(Class("p-4 flex items-center"),
					A(
						Href("/"),
						Aria("title", "Go to home"),
						Raw(logo),
					),
					A(
						hx.Boost("true"),
						Href("/auth/login"),
						Class(`bg-ctabase-300/20 hover:bg-ctabase-300 focus:ring-4
					focus:ring-blue-300 font-medium text-sm rounded-2xl px-8 py-[0.375rem]
					dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none
					border-cta border
					dark:focus:ring-blue-800 ml-auto`),
						Text("Login"),
					),
				),
				Group(children),
			),
		),
	))
}

// Background Square with Rounded Corners
const logo = `<svg
	xmlns="http://www.w3.org/2000/svg"
	viewBox="0 0 200 200"
	width="2rem"
	height="2rem"
><rect x="10" y="10" width="180" height="180" rx="20" ry="20" fill="#e6e6e6" stroke="#999999" stroke-width="4"></rect></svg>`
