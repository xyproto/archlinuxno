package main

import (
	"github.com/hoisie/web"
	"github.com/xyproto/browserspeak"
	"github.com/xyproto/genericsite"
	"github.com/xyproto/instapage"
	"github.com/xyproto/siteengines"
)

const JQUERY_VERSION = "2.0.0"

func hello(val string) string {
	return instapage.Message("root page", "hello: "+val)
}

func helloHandle(ctx *web.Context, name string) string {
	return "Hello, " + name
}

func notFound2(ctx *web.Context, val string) {
	ctx.ResponseWriter.WriteHeader(404)
	ctx.ResponseWriter.Write([]byte(browserspeak.NotFound(ctx, val)))
}

func ServeEngines(userState *genericsite.UserState, mainMenuEntries genericsite.MenuEntries) {
	// The user engine
	userEngine := siteengines.NewUserEngine(userState)
	userEngine.ServePages("archlinux.no")

	// The admin engine
	adminEngine := siteengines.NewAdminEngine(userState)
	adminEngine.ServePages(ArchBaseCP, mainMenuEntries)

	// TODO: Move this one to roboticoverlords instead
	// The dynamic IP webpage (returns an *IPState)
	ipEngine := siteengines.NewIPEngine(userState)
	ipEngine.ServePages()

	// The chat system (see also the menu entry in ArchBaseCP)
	chatEngine := siteengines.NewChatEngine(userState)
	chatEngine.ServePages(ArchBaseCP, mainMenuEntries)

	// Wiki engine
	wikiEngine := siteengines.NewWikiEngine(userState)
	wikiEngine.ServePages(ArchBaseCP, mainMenuEntries)

	// Blog engine
	//blogEngine := NewBlogEngine(userState)
	//blogEngine.ServePages(ArchBaseCP, mainMenuEntries)
}

// TODO: One database per site
func main() {

	// UserState with a Redis Connection Pool
	userState := genericsite.NewUserState()
	defer userState.Close()

	// The archlinux.no webpage,
	mainMenuEntries := ServeArchlinuxNo(userState, "/js/jquery-"+JQUERY_VERSION+".min.js")

	ServeEngines(userState, mainMenuEntries)

	// Compilation errors, vim-compatible filename
	web.Get("/error", browserspeak.GenerateErrorHandle("errors.err"))
	web.Get("/errors", browserspeak.GenerateErrorHandle("errors.err"))

	// Various .php and .asp urls that showed up in the log
	genericsite.ServeForFun()

	// TODO: Incorporate this check into web.go, to only return
	// stuff in the header when the HEAD method is requested:
	// if ctx.Request.Method == "HEAD" { return }
	// See also: curl -I

	// Serve on port 3000 for the Nginx instance to use
	web.Run("0.0.0.0:3000")
}
