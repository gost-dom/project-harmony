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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!doctype html><html><head><link rel=\"stylesheet\" href=\"/static/css/tailwind.css\"><meta charset=\"utf-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"><title>Meet the Locals</title></head><body class=\"bg-secondary-50 dark:bg-stone-800\"><div class=\"min-h-screen flex flex-col\"><header class=\"bg-primary-600 text-zinc-50 p-4 font-bold\"><svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 200 200\" width=\"200\" height=\"200\"><!-- Background Square with Rounded Corners --><rect x=\"10\" y=\"10\" width=\"180\" height=\"180\" rx=\"20\" ry=\"20\" fill=\"#f0f0f0\" stroke=\"#cccccc\" stroke-width=\"4\"></rect><!-- Left Person --><circle cx=\"60\" cy=\"80\" r=\"20\" fill=\"#f5c7a9\"></circle><!-- Head --><path d=\"M45 100 Q60 115, 75 100 Q70 135, 50 135 Q30 135, 45 100 Z\" fill=\"#d1eaff\"></path><!-- Body --><!-- Right Person --><circle cx=\"140\" cy=\"80\" r=\"20\" fill=\"#f5c7a9\"></circle><!-- Head --><path d=\"M125 100 Q140 115, 155 100 Q150 135, 130 135 Q110 135, 125 100 Z\" fill=\"#ffd1d1\"></path><!-- Body --><!-- Meeting Gesture --><line x1=\"75\" y1=\"100\" x2=\"125\" y2=\"100\" stroke=\"#cc8c6b\" stroke-width=\"6\"></line><!-- Arms --><circle cx=\"100\" cy=\"100\" r=\"8\" fill=\"#e0a889\"></circle><!-- Connection --></svg> Meet the Locals</header><div id=\"body-root\" class=\"flex-grow flex items-stretch flex-col\" hx-swap-oob=\"true\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = c.body.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "</div><!--\n\t\t\t\t<footer\n\t\t\t\t\tclass=\"bg-primary-600 text-white p-4 border-t\n      border-primary-800\"\n\t\t\t\t>Footer content</footer>\n        --></div></body></html><!--\n    .min-h-screen.flex.flex-col\n      header.bg-primary-600.text-zinc-50.p-4.font-bold\n        | Meet the Locals\n\n      div#body-root.flex-grow.flex.items-stretch.flex-col(hx-swap-oob='true')\n        block content\n\n      footer.bg-primary-600.text-white.p-4.border-t.border-primary-800\n        | Footer content\n        -->")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
