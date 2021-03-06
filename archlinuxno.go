// First attempt at a home grown webpage for archlinux.no
package main

import (
	"github.com/hoisie/web"
	"github.com/xyproto/genericsite"
	"github.com/xyproto/permissions2"
	"github.com/xyproto/pinterface"
	"github.com/xyproto/siteengines"
	"github.com/xyproto/webhandle"
)

const JQUERY_VERSION = "2.0.0"

func hello(val string) string {
	return webhandle.Message("root page", "hello: "+val)
}

func helloHandle(ctx *web.Context, name string) string {
	return "Hello, " + name
}

func notFound2(ctx *web.Context, val string) {
	ctx.ResponseWriter.WriteHeader(404)
	ctx.ResponseWriter.Write([]byte(webhandle.NotFound(ctx, val)))
}

func ServeEngines(userState pinterface.IUserState, mainMenuEntries genericsite.MenuEntries) error {
	// The user engine
	userEngine, err := siteengines.NewUserEngine(userState)
	if err != nil {
		return err
	}
	userEngine.ServePages("archlinux.no")

	// The admin engine
	adminEngine, err := siteengines.NewAdminEngine(userState)
	if err != nil {
		return err
	}
	adminEngine.ServePages(ArchBaseCP, mainMenuEntries)

	// TODO: Move this one to roboticoverlords instead
	// The dynamic IP webpage (returns an *IPState)
	ipEngine, err := siteengines.NewIPEngine(userState)
	if err != nil {
		return err
	}
	ipEngine.ServePages()

	// The chat system (see also the menu entry in ArchBaseCP)
	chatEngine, err := siteengines.NewChatEngine(userState)
	if err != nil {
		return err
	}
	chatEngine.ServePages(ArchBaseCP, mainMenuEntries)

	// Wiki engine
	wikiEngine, err := siteengines.NewWikiEngine(userState)
	if err != nil {
		return err
	}
	wikiEngine.ServePages(ArchBaseCP, mainMenuEntries)

	// Blog engine
	//blogEngine := NewBlogEngine(userState)
	//blogEngine.ServePages(ArchBaseCP, mainMenuEntries)

	return nil
}

// TODO: One database per site
func main() {

	// UserState with a Redis Connection Pool
	userState := permissions.NewUserState(0, true, ":6379")

	defer userState.Close()

	// The archlinux.no webpage,
	mainMenuEntries := ServeArchlinuxNo(userState, "/js/jquery-"+JQUERY_VERSION+".min.js")

	ServeEngines(userState, mainMenuEntries)

	// Compilation errors, vim-compatible filename
	web.Get("/error", webhandle.GenerateErrorHandle("errors.err"))
	web.Get("/errors", webhandle.GenerateErrorHandle("errors.err"))

	// Various .php and .asp urls that showed up in the log
	genericsite.ServeForFun()

	// TODO: Incorporate this check into web.go, to only return
	// stuff in the header when the HEAD method is requested:
	// if ctx.Request.Method == "HEAD" { return }
	// See also: curl -I

	// Serve on port 3009 for the Nginx instance to use
	web.Run("0.0.0.0:3009")
}
