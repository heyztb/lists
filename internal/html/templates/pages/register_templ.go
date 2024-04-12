// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.648
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "github.com/heyztb/lists-backend/internal/html/templates/shared"
import "github.com/heyztb/lists-backend/internal/html/templates/partials"

func submitRegistration() templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_submitRegistration_1d2f`,
		Function: `function __templ_submitRegistration_1d2f(){event.preventDefault()
  const form = new FormData(event.target)
  const data = Object.fromEntries(form)
  console.log(data)
  try {
    registerUser(data['email'], data['password']).then(success => {
      console.log(success)
    })
  } catch (e) {
    console.log(e)
  }
}`,
		Call:       templ.SafeScript(`__templ_submitRegistration_1d2f`),
		CallInline: templ.SafeScriptInline(`__templ_submitRegistration_1d2f`),
	}
}

func Register() templ.Component {
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
		templ_7745c5c3_Var2 := templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
			if !templ_7745c5c3_IsBuffer {
				templ_7745c5c3_Buffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script type=\"module\">\n\t\timport { registerUser } from '/assets/srp.min.js'\n\t\t\n\t\tasync function onSubmit(event) {\n\t\t  event.preventDefault()\n      const form = new FormData(event.target)\n      const data = Object.fromEntries(form)\n      console.log(data)\n      try {\n        const success = await registerUser(data['email'], data['password'])\n        if (success) {\n          document.getElementById('response').innerHTML = 'User registered, redirecting to login in 3 seconds..'\n          setTimeout(() => {\n            window.location.pathname = '/login'\n          }, 3000)\n        }\n      } catch (e) {\n        console.log(e)\n          document.getElementById('response').innerHTML = e \n          document.getElementById('password').value = ''\n          document.getElementById('confirm-password').value = ''\n      }\n\t\t}\n\n\t\tconst form = document.getElementById('registration')\n\t\tform.addEventListener('submit', onSubmit)\n\t\t</script> <div class=\"w-full lg:grid lg:grid-cols-2 min-h-screen\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = partials.NavAuth(true).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"hidden bg-muted lg:block\"><h1 class=\"text-4xl font-bold lg:text-5xl\">Lists has been created with your privacy and data security in mind.</h1></div><div class=\"flex items-center justify-center py-12\"><div class=\"mx-auto grid w-[350px] gap-6\"><div class=\"grid gap-2 text-center\"><h1 class=\"text-3xl font-bold\">Create an account</h1><p class=\"text-balance text-muted-foreground\">Get started with Lists today</p></div><form class=\"grid gap-4\" id=\"registration\" method=\"POST\"><div class=\"grid gap-2\"><label class=\"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70\" for=\"email\">Email</label> <input class=\"flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50\" name=\"email\" id=\"email\" placeholder=\"m@example.com\" required type=\"email\"></div><div class=\"grid gap-2\"><div class=\"flex items-center\"><label class=\"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70\" for=\"password\">Password</label></div><input class=\"flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50\" name=\"password\" id=\"password\" required type=\"password\"></div><div class=\"grid gap-2\"><div class=\"flex items-center\"><label class=\"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70\" for=\"confirm-password\">Confirm password</label></div><input class=\"flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50\" name=\"confirm-password\" id=\"confirm-password\" required type=\"password\"></div><div class=\"flex flex-col gap-2\"><button class=\"inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2 flex-1\" type=\"submit\">Create account</button></div></form><div id=\"response\" class=\"text-md\"></div></div></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !templ_7745c5c3_IsBuffer {
				_, templ_7745c5c3_Err = io.Copy(templ_7745c5c3_W, templ_7745c5c3_Buffer)
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = shared.Page("Register").Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
