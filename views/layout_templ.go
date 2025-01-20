// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.819
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

type contents struct {
	body templ.Component
}

func layout(c contents) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!doctype html><html><head><link rel=\"stylesheet\" href=\"/static/css/tailwind.css\"><meta charset=\"utf-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"><title>Project harmony</title></head><body class=\"bg-secondary-50 dark:bg-stone-800\"><div class=\"min-h-screen flex flex-col\"><!--\n\t\t\t\t<header class=\"bg-primary-600 text-zinc-50 p-4 font-bold flex\">\n        --><header class=\"p-4 flex items-center\"><a href=\"/\" aria-title=\"Go to home\"><svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 200 200\" width=\"2rem\" height=\"2rem\"><!-- Background Square with Rounded Corners --><rect x=\"10\" y=\"10\" width=\"180\" height=\"180\" rx=\"20\" ry=\"20\" fill=\"#e6e6e6\" stroke=\"#999999\" stroke-width=\"4\"></rect></svg></a> <a href=\"/auth/login\" class=\"bg-ctabase-300 hover:bg-ctabase-300 bg-opacity-20 focus:ring-4\nfocus:ring-blue-300 font-medium text-sm rounded-2xl px-8 py-[0.375rem]\ndark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none\nborder-cta border\ndark:focus:ring-blue-800 ml-auto\">Login</a></header><div id=\"body-root\" class=\"flex-grow flex items-stretch flex-col\" hx-swap-oob=\"true\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = c.body.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "</div><!--\n\t\t\t\t<footer\n\t\t\t\t\tclass=\"bg-primary-600 text-white p-4 border-t\n      border-primary-800\"\n\t\t\t\t>Footer content</footer>\n        --></div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
