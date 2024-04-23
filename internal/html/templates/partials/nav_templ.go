// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.648
package partials

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func Nav() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<header class=\"bg-background shadow px-4 lg:px-6 h-14 flex items-center fixed top-0 left-0 right-0\"><a class=\"flex items-center justify-center\" href=\"/\"><svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"h-6 w-6\"><path d=\"m8 3 4 8 5-5 5 15H2L8 3z\"></path></svg> <span class=\"sr-only\">Lists</span></a><nav class=\"ml-auto flex gap-4 sm:gap-6\"><a class=\"text-sm font-medium hover:underline underline-offset-4\" href=\"#features\">Features</a> <a class=\"text-sm font-medium hover:underline underline-offset-4\" href=\"/about\">About</a><div class=\"border border-foreground border-l\"></div><a class=\"text-sm font-medium hover:underline underline-offset-4\" href=\"/login\">Log in</a> <a class=\"text-sm font-medium hover:underline underline-offset-4\" href=\"/register\">Sign up\t</a></nav></header>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
