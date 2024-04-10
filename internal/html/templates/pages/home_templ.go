// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.648
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "github.com/heyztb/lists-backend/internal/html/templates/shared"

func Home(title string) templ.Component {
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
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!--\n// v0 by Vercel.\n// https://v0.dev/t/GGQFhyAjdhj\n--> <div class=\"flex flex-col min-h-[100dvh]\"><header class=\"bg-background shadow px-4 lg:px-6 h-14 flex items-center fixed top-0 left-0 right-0\"><a class=\"flex items-center justify-center\" href=\"/\"><svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"h-6 w-6\"><path d=\"m8 3 4 8 5-5 5 15H2L8 3z\"></path></svg> <span class=\"sr-only\">Lists</span></a><nav class=\"ml-auto flex gap-4 sm:gap-6\"><a class=\"text-sm font-medium hover:underline underline-offset-4\" href=\"#features\">Features</a> <a class=\"text-sm font-medium hover:underline underline-offset-4\" href=\"/about\">About</a><div class=\"border border-foreground border-l\"></div><a class=\"text-sm font-medium hover:underline underline-offset-4\" href=\"/login\">Log in</a> <a class=\"text-sm font-medium hover:underline underline-offset-4\" href=\"/register\">Sign up\t</a></nav></header><main class=\"flex-1\"><section class=\"w-full py-12 md:py-24 lg:py-32\"><div class=\"container px-4 md:px-6\"><div class=\"flex flex-col items-center space-y-4 text-center mt-6 lg:mt-0\"><div class=\"space-y-2\"><h1 class=\"text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl max-w-2/3\">Streamline your productivity, without sacrificing your privacy.</h1><p class=\"mx-auto max-w-[700px] text-gray-500 md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed dark:text-gray-400\">Over-engineered with your security and privacy in mind.</p></div><div class=\"space-x-4\"><a class=\"inline-flex h-9 items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-gray-50 shadow\" href=\"/register\">Get Started</a> <a class=\"inline-flex h-9 items-center justify-center rounded-md border border-gray-200 bg-white px-4 py-2 text-sm font-medium shadow-sm transition-colors hover:bg-gray-100 hover:text-gray-900 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-gray-950 disabled:pointer-events-none disabled:opacity-50 dark:border-gray-800 dark:bg-gray-950 dark:hover:bg-gray-800 dark:hover:text-gray-50 dark:focus-visible:ring-gray-300\" href=\"#features\">Learn more</a></div></div></div></section><section class=\"w-full py-12 md:py-24 lg:py-32 border-t\" id=\"features\"><div class=\"container grid items-center gap-4 px-4 text-center md:px-6 lg:gap-10\"><div class=\"space-y-3\"><h2 class=\"text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl\">Features</h2><p class=\"mx-auto max-w-[600px] text-gray-500 md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed dark:text-gray-400\">Everything you expect, carefully designed to respect your privacy.</p></div><div class=\"mx-auto grid max-w-5xl items-start gap-8 sm:grid-cols-2 lg:gap-12\"><div class=\"grid gap-1\"><h3 class=\"text-lg font-bold\">Task Organization</h3><p class=\"text-sm text-gray-500 dark:text-gray-400\">Keep your tasks organized with categories, tags, and priorities.</p></div><div class=\"grid gap-1\"><h3 class=\"text-lg font-bold\">Due Date Reminders</h3><p class=\"text-sm text-gray-500 dark:text-gray-400\">Never miss a deadline with automatic due date reminders.</p></div><div class=\"grid gap-1\"><h3 class=\"text-lg font-bold\">Collaboration</h3><p class=\"text-sm text-gray-500 dark:text-gray-400\">Share your tasks with your team and collaborate in real time.</p></div><div class=\"grid gap-1\"><h3 class=\"text-lg font-bold\">Private</h3><p class=\"text-sm text-gray-500 dark:text-gray-400\">You own your data. Lists is designed to never store or transmit data in plaintext.</p></div><div class=\"grid gap-1 col-span-2\"><h3 class=\"text-lg font-bold\">Open source. Forever.</h3><p class=\"text-sm text-gray-500 dark:text-gray-400\">Lists will always be open source. Feel free to <a class=\"underline underline-offset-4\" href=\"https://github.com/heyztb/lists\">take a look</a>.</p></div></div></div></section><section class=\"w-full py-12 md:py-24 lg:py-32 bg-gray-100 dark:bg-gray-800\"><div class=\"container grid items-center gap-4 px-4 text-center md:px-6 lg:gap-10\"><div class=\"space-y-3\"><h2 class=\"text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl\">A modern take on the age old classic.</h2><p class=\"mx-auto max-w-[600px] text-gray-500 md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed dark:text-gray-400\">Focus on planning your next big thing, let us worry about the rest.</p></div><div class=\"mx-auto w-full max-w-sm space-y-2\"><form class=\"flex space-x-2\"><input type=\"email\" class=\"flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 max-w-lg flex-1\" placeholder=\"Enter your email\"> <button class=\"inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2\" type=\"submit\">Sign Up</button></form><p class=\"text-xs text-gray-500 dark:text-gray-400\">Sign up to get notified when we launch. <a class=\"underline underline-offset-2\" href=\"#\" rel=\"ugc\">Terms &amp; Conditions</a></p></div></div></section></main><footer class=\"flex flex-col gap-2 sm:flex-row py-6 w-full shrink-0 items-center px-4 md:px-6 border-t\"><p class=\"text-xs text-gray-500 dark:text-gray-400\">© 2024 Zachary Blake. All rights reserved.</p><nav class=\"sm:ml-auto flex gap-4 sm:gap-6\"><a class=\"text-xs hover:underline underline-offset-4\" href=\"/tos\">Terms of Service</a> <a class=\"text-xs hover:underline underline-offset-4\" href=\"/privacy\">Privacy</a></nav></footer></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !templ_7745c5c3_IsBuffer {
				_, templ_7745c5c3_Err = io.Copy(templ_7745c5c3_W, templ_7745c5c3_Buffer)
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = shared.Page("Home").Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}