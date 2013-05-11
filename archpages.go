package main

import (
	"strconv"
	"time"

	"github.com/hoisie/web"
	. "github.com/xyproto/browserspeak"
	. "github.com/xyproto/genericsite"
	. "github.com/xyproto/siteengines"
)

// The default settings for Arch Linux content pages
func ArchBaseCP(state *UserState) *ContentPage {
	cp := DefaultCP(state)
	cp.Title = "Arch Linux"
	cp.Subtitle = "no"

	y := time.Now().Year()

	// TODO: Use templates for the footer, for more accurate measurment of the time made to generate the page
	cp.FooterText = "Alexander Rødseth, " + strconv.Itoa(y)

	cp.Url = "/" // Is replaced when the contentpage is published

	cp.ColorScheme = NewArchColorScheme()

	// Behind the text
	//cp.BgImageURL = "/img/nasty_fabric.png"
	//cp.BgImageURL = "/img/cloth_alike.png"
	//cp.BgImageURL = "/img/strange_bullseyes.png"
	cp.BgImageURL = "/img/rough_diagonal.png"
	cp.StretchBackground = false

	// Behind the menu
	//cp.BackgroundTextureURL = "/img/bg2.png"
	//cp.BackgroundTextureURL = "/img/simple_dashed.png"
	//cp.BackgroundTextureURL = "/img/grey.png"
	//cp.BackgroundTextureURL = "/img/pw_maze_black.png"
	//cp.BackgroundTextureURL = "/img/black_twill.png"
	cp.BackgroundTextureURL = "/img/dark_wood.png"
	//cp.BackgroundTextureURL = "/img/hixs_pattern_evolution.png"
	//ps_neutral.png"

	return cp
}

// Returns a ArchBaseCP with the contentTitle set
func ArchBaseTitleCP(contentTitle string, userState *UserState) *ContentPage {
	cp := ArchBaseCP(userState)
	cp.ContentTitle = contentTitle
	return cp
}

func OverviewCP(userState *UserState, url string) *ContentPage {
	cp := ArchBaseCP(userState)
	cp.ContentTitle = "Arch Linux Norway"
	//cp.ContentHTML = `This site is currently under construction.<br />Visit the <a href="https://bbs.archlinux.org/viewtopic.php?id=4998">Arch Linux Forum</a> in the meantime.<br /><br /><i>- Alexander Rødseth &lt;rodseth / gmail&gt;</i>`
	cp.ContentHTML = `<img src="https://si0.twimg.com/profile_images/71498752/arch-linux-logo.png">`
	cp.Url = url
	return cp
}

func MirrorsCP(userState *UserState, url string) *ContentPage {
	cp := ArchBaseCP(userState)
	cp.ContentTitle = "Mirrors"
	cp.ContentHTML = "List over Norwegian Arch Linux mirrors:"
	cp.Url = url
	return cp
}

func CountCP(userState *UserState, url string) *ContentPage {
	apc := ArchBaseCP(userState)
	apc.ContentTitle = "Counting"
	apc.ContentHTML = "1 2 3"
	apc.Url = url
	return apc
}

func BobCP(userState *UserState, url string) *ContentPage {
	apc := ArchBaseCP(userState)
	apc.ContentTitle = "Bob"
	if userState.HasUser("bob") {
		apc.ContentHTML = "has bob, l "
	} else {
		apc.ContentHTML = "no bob, l "
	}
	if userState.IsLoggedIn("bob") {
		apc.ContentHTML += "yes"
	} else {
		apc.ContentHTML += "no"
	}
	apc.Url = url
	return apc
}

func JQueryCP(userState *UserState, url string) *ContentPage {
	apc := ArchBaseCP(userState)
	apc.ContentTitle = "JQuery"

	apc.ContentHTML = "<button id=clickme>bob</button><br />"
	apc.ContentHTML += "<div id=status>status</div>"

	//apc.contentJS = OnClick("#clickme", GetTest())
	//apc.contentJS += OnClick("#clickme", SetText("#clickme", "ost"))
	//apc.contentJS += OnClick("#clickme", SetTextFromURL("#clickme", "http://archlinux.no/status/bob"))
	//apc.contentJS += OnClick("#clickme", GetTest())

	apc.ContentJS += Load("#status", "/status/elg")
	apc.ContentJS += OnClick("#clickme", Load("#status", "/status/bob"))
	apc.ContentJS += SetText("#menuJQuery", "Heppa")

	apc.Url = url

	return apc
}

func TextCP(userState *UserState, url string) *ContentPage {
	apc := ArchBaseCP(userState)
	apc.ContentTitle = "YOLO narwhal"
	apc.ContentHTML = `<p>Locavore Austin fanny pack pickled.  Marfa hoodie pitchfork american apparel, flexitarian YOLO pickled keytar twee cred craft beer seitan authentic raw denim kogi.  Selvage mixtape blog, pickled cosby sweater williamsburg skateboard brooklyn lo-fi twee.  Blue bottle echo park kale chips, selvage fap skateboard swag chambray tousled.  Street art etsy four loko fap, iphone carles cliche banh mi fashion axe PBR authentic leggings.  Narwhal mumblecore street art tumblr.  Messenger bag vice art party, next level aesthetic church-key tumblr direct trade  typewriter street art.</p><p>Messenger bag blue bottle VHS before they sold out.  Artisan pickled swag, VHS meggings jean shorts blog tonx salvia cosby sweater mumblecore aesthetic literally narwhal.  Brunch tofu gluten-free disrupt blog occupy.  Austin bicycle rights sartorial narwhal, butcher trust fund cred.  Neutra kale chips letterpress literally, williamsburg kogi brunch bicycle rights.  Williamsburg craft beer brunch quinoa, forage YOLO swag put a bird on it four loko mixtape banksy.  Tumblr semiotics yr fixie.</p><p>Iphone banksy wolf squid wayfarers, VHS photo booth banh mi fap.  Tonx flexitarian vinyl scenester terry richardson squid synth deep v.  VHS tousled godard, cardigan american apparel lo-fi flannel.  Vice church-key cliche, hashtag banh mi direct trade  skateboard.  Sriracha meh pitchfork, wayfarers helvetica leggings try-hard viral YOLO lo-fi fingerstache synth ennui next level ugh.  Wayfarers organic american apparel fingerstache craft beer bicycle rights, beard keffiyeh banksy four loko butcher hashtag mumblecore banjo wes anderson.  Williamsburg next level deep v pickled typewriter kogi.</p><p>Meggings gastropub flexitarian, before they sold out DIY wes anderson cred authentic artisan dreamcatcher aesthetic ennui food truck.  Fanny pack selvage synth vegan pug.  YOLO shoreditch pitchfork, letterpress whatever put a bird on it truffaut mumblecore flannel terry richardson irony cray master cleanse ethnic gluten-free.  Fap banksy blog pickled meh ethnic food truck +1, vice leggings retro quinoa.  Small batch vice pop-up mustache.  +1 ethnic echo park semiotics letterpress raw denim.  Keytar photo booth wes anderson, freegan before they sold out skateboard seitan brooklyn.</p><p>Wes anderson high life banksy messenger bag art party plaid disrupt tattooed, next level swag viral raw denim.  Cliche meggings terry richardson cray.  Next level 3 wolf moon retro marfa.  Pork belly authentic banjo, iphone lomo williamsburg letterpress cosby sweater Austin typewriter quinoa skateboard hoodie.  Plaid kale chips godard farm-to-table.  Fashion axe mixtape freegan, pop-up chambray ugh etsy YOLO jean shorts dreamcatcher meggings.  Banh mi letterpress tousled, skateboard stumptown high life vegan fap typewriter shoreditch 8-bit lo-fi master cleanse selfies bespoke.</p>`
	apc.Url = url
	return apc
}

func HelloCP(userState *UserState, url string) *ContentPage {
	apc := ArchBaseCP(userState)
	apc.ContentTitle = "This is it"
	apc.Url = url
	return apc
}

// This is where the possibilities for the menu are listed
// TODO: Find a way to support many languages
// TODO: Refactor this ugly function
func Cps2MenuEntries(cps []ContentPage) MenuEntries {
	links := []string{
		"Overview:/",
		"Login:/login",
		"Register:/register",
		"Logout:/logout",
		"Admin:/admin",
		"Chat:/chat",
		"Wiki:/wiki",
		"Text:/text",
	}
	//for _, cp := range cps {
	//	text_and_url := cp.ContentTitle + ":" + cp.Url
	//	links = append(links, text_and_url)
	//}
	return Links2menuEntries(links)
}

// Routing for the archlinux.no webpage
// Admin, search and user management is already provided
func ServeArchlinuxNo(userState *UserState, jquerypath string) MenuEntries {
	cps := []ContentPage{
		*OverviewCP(userState, "/"),
		*TextCP(userState, "/text"),
		*JQueryCP(userState, "/jquery"),
		*BobCP(userState, "/bob"),
		*CountCP(userState, "/counting"),
		*MirrorsCP(userState, "/mirrors"),
		*HelloCP(userState, "/feedback"),
		*LoginCP(ArchBaseCP, userState, "/login"),
		*RegisterCP(ArchBaseCP, userState, "/register"),
	}

	menuEntries := Cps2MenuEntries(cps)

	// template content generator
	tvgf := DynamicMenuFactoryGenerator(menuEntries)

	// TODO: Simplify ServeSearchPages to take fewer parameters
	ServeSearchPages(ArchBaseCP, userState, cps, ArchBaseCP(userState).ColorScheme, tvgf(userState))
	ServeSite(ArchBaseCP, userState, cps, tvgf, jquerypath)

	// "dynamic" pages
	// Makes helloSF handle the content for /hello/(.*) urls, but wrapped in a BaseCP with the title "Hello"
	web.Get("/hello/(.*)", ArchBaseTitleCP("Hello", userState).WrapWebHandle(helloHandle, tvgf(userState)))

	return menuEntries
}

func NewArchColorScheme() *ColorScheme {
	var cs ColorScheme
	cs.Darkgray = "#202020"
	cs.Nicecolor = "#5080D0"   // nice blue
	cs.Menu_link = "#c0c0c0"   // light gray
	cs.Menu_hover = "#efefe0"  // light gray, somewhat yellow
	cs.Menu_active = "#ffffff" // white
	cs.Default_background = "#000030"
	cs.TitleText = "#e0e0e0" // The first word of the title text
	return &cs
}
