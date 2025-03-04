// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.819
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func AuthLogin(redirectUrl string, data LoginFormData) templ.Component {
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
		templ_7745c5c3_Err = layout(contents{body: login_body(redirectUrl, data)}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

type LoginFormData struct {
	Email              string
	EmailMissing       bool
	Password           string
	PasswordMissing    bool
	InvalidCredentials bool
}

func boolToString(b bool) string {
	if b {
		return "true"
	} else {
		return "false"
	}
}

func login_body(redirectUrl string, formData LoginFormData) templ.Component {
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
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div class=\"w-full m-auto flex flex-col items-center\"><div class=\"my-4 font-bold text-2xl dark:text-white\">Project Harmony</div><div class=\"bg-white rounded-lg shadow-md border md:mt-0 w-full sm:max-w-xl\n  xl:p-0 dark:bg-gray-800 dark:border-gray-700\"><main class=\"p-6 space-y-4 md:space-y-6 sm:p-8\"><h1 class=\"text-center\ntext-xl font-bold leading-tight tracking-tight text-gray-900 md:text-4xl dark:text-white\n        \">Login</h1><form id=\"login-form\" class=\"space-y-4 md:space-y-6\" hx-post=\"/auth/login\" hx-swap=\"innerHTML\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = LoginForm(redirectUrl, formData).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "</form></main></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func invalid(v bool) (res templ.Attributes) {
	return templ.Attributes{
		"aria-invalid": boolToString(v),
	}
}

func LoginForm(redirectUrl string, formData LoginFormData) templ.Component {
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
		templ_7745c5c3_Var3 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var3 == nil {
			templ_7745c5c3_Var3 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "<input type=\"hidden\" name=\"redirectUrl\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(redirectUrl)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/server/views/auth_login.templ`, Line: 61, Col: 21}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = field(fieldOptions{
			inputOptions: inputOptions{
				id:              "email",
				name:            "email",
				inputType:       "text",
				required:        true,
				validationError: "Email is required",
				autofocus:       true,
				value:           formData.Email,
				invalid:         formData.EmailMissing,
				attributes:      invalid(formData.EmailMissing),
			},
			label: "Email"}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = field(fieldOptions{
			inputOptions: inputOptions{
				id:              "password",
				name:            "password",
				inputType:       "password",
				validationError: "Password is required",
				required:        true,
				placeholder:     "••••••••",
				value:           formData.Password,
				invalid:         formData.PasswordMissing,
				attributes:      invalid(formData.PasswordMissing),
			},
			label: "Password",
		}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "<div class=\"flex items-center justify-between\"><div class=\"flex items-start\"><!--\n\t\t\t\t\t\t<div class=\"flex items-center h-5\">\n\t\t\t\t\t\t\t<input\n\t\t\t\t\t\t\t\tclass=\"w-4 h-4 border border-gray-300 rounded bg-gray-50 focus:ring-3 focus:ring-primary-300 dark:bg-gray-700 dark:border-gray-600 dark:focus:ring-primary-600 dark:ring-offset-gray-800\"\n\t\t\t\t\t\t\t\ttype=\"checkbox\"\n\t\t\t\t\t\t\t\tid=\"remember\"\n\t\t\t\t\t\t\t/>\n\t\t\t\t\t\t</div>\n\t\t\t\t\t\t<div class=\"ml-3 text-sm\">\n\t\t\t\t\t\t\t<label\n\t\t\t\t\t\t\t\tclass=\"block text-sm font-medium  text-gray-500 dark:text-gray-300\"\n\t\t\t\t\t\t\t\tfor=\"remember\"\n\t\t\t\t\t\t\t>Remember me</label>\n\t\t\t\t\t\t</div>\n            --></div><!--\n\t\t\t\t\t<a href=\"#\" class=\"text-sm font-medium text-primary-600 hover:underline dark:text-primary-500\">Forgot password?</a>\n          --></div><button id=\"submit-login-form-button\" type=\"submit\" class=\"w-full text-white bg-cta hover:bg-ctabase-900 focus:ring-4\n    focus:outline-none focus:ring-primary-300 font-medium rounded-lg text-sm\n    px-5 py-2.5 text-center dark:bg-primary-600 dark:hover:bg-primary-700\n    dark:focus:ring-primary-800\">Sign in</button> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if formData.InvalidCredentials {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "<div id=\"alert-div\" role=\"alert\" aria-live=\"assertive\" class=\"text-red-700\">Email or password did not match</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, "<!--\n    <p class=\"text-sm font-light text-gray-500 dark:text-gray-400\">\n        Don’t have an account yet? <a href=\"#\" class=\"font-medium text-primary-600 hover:underline dark:text-primary-500\">Sign up</a>\n    </p>\n    -->")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
