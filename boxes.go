package main

// Various elements of a webpage

import (
	"strings"

	"github.com/xyproto/browserspeak"
)

func AddTopBox(page *browserspeak.Page, title, subtitle, searchURL, searchButtonText, backgroundTextureURL string) (*browserspeak.Tag, error) {
	body, err := page.GetTag("body")
	if err != nil {
		return nil, err
	}

	div := body.AddNewTag("div")
	div.AddAttr("id", "topbox")
	div.AddStyle("display", "block")
	div.AddStyle("width", "100%")
	div.AddStyle("margin", "0")
	div.AddStyle("padding", "0 0 1em 0")
	div.AddStyle("position", "absolute")
	div.AddStyle("top", "0")
	div.AddStyle("left", "0")
	div.AddStyle("background-color", NICEGRAY)
	div.AddStyle("position", "fixed")

	titlebox := AddTitleBox(div, title, subtitle)
	titlebox.AddAttr("id", "titlebox")
	titlebox.AddStyle("margin", "0 0 0 0")
	// Padding-top + height should be 5em, padding decides the position
	titlebox.AddStyle("padding", "1.8em 0 0 2.8em")
	titlebox.AddStyle("height", "3.2em")
	titlebox.AddStyle("width", "100%")
	titlebox.AddStyle("position", "fixed")
	titlebox.AddStyle("background-color", NICEGRAY) // gray, could be a gradient
	titlebox.AddStyle("background", "url('"+backgroundTextureURL+"')")

	searchbox := AddSearchBox(div, searchURL, searchButtonText)
	searchbox.AddAttr("id", "searchbox")
	searchbox.AddStyle("position", "relative")
	searchbox.AddStyle("margin-top", "2em")
	searchbox.AddStyle("margin-right", "2%")

	return div, nil
}

func AddFooter(page *browserspeak.Page, footerText, footerTextColor, footerColor string) (*browserspeak.Tag, error) {
	body, err := page.GetTag("body")
	if err != nil {
		return nil, err
	}
	div := body.AddNewTag("div")
	div.AddAttr("id", "notice")
	div.AddStyle("position", "fixed")
	div.AddStyle("bottom", "0")
	div.AddStyle("left", "0")
	div.AddStyle("width", "100%")
	div.AddStyle("display", "block")
	div.AddStyle("padding", "0")
	div.AddStyle("margin", "0")
	div.AddStyle("background-color", footerColor)
	div.AddStyle("font-size", "0.6em")
	div.AddStyle("text-align", "right")
	div.AddStyle("box-shadow", "1px -2px 3px rgba(0,0,0, .5)")

	innerdiv := div.AddNewTag("div")
	innerdiv.AddAttr("id", "innernotice")
	innerdiv.AddStyle("padding", "0 2em 0 0")
	innerdiv.AddStyle("margin", "0")
	innerdiv.AddStyle("color", footerTextColor)
	innerdiv.AddContent(footerText)

	return div, nil
}

func AddContent(page *browserspeak.Page, contentTitle, contentHTML string) (*browserspeak.Tag, error) {
	body, err := page.GetTag("body")
	if err != nil {
		return nil, err
	}

	div := body.AddNewTag("div")
	div.AddAttr("id", "content")
	div.AddStyle("z-index", "-1")
	div.AddStyle("color", "black") // content headline color
	div.AddStyle("min-height", "80%")
	div.AddStyle("min-width", "60%")
	div.AddStyle("float", "left")
	div.AddStyle("position", "relative")
	div.AddStyle("margin-left", "5%")
	div.AddStyle("margin-top", "10em")
	div.AddStyle("margin-right", "5em")
	div.AddStyle("padding-left", "4em")
	div.AddStyle("padding-right", "5em")
	div.AddStyle("padding-top", "1em")
	div.AddStyle("padding-bottom", "2em")
	div.AddStyle("background-color", "rgba(255,255,255,0.92)") // light gray
	div.RoundedBox()

	h2 := div.AddNewTag("h2")
	h2.AddAttr("id", "textheader")
	h2.AddContent(contentTitle)
	h2.CustomSansSerif("Armata")

	p := div.AddNewTag("p")
	p.AddAttr("id", "textparagraph")
	p.AddStyle("margin-top", "0.5em")
	p.CustomSansSerif("Junge")
	p.AddStyle("font-size", "1.0em")
	p.AddStyle("color", "black") // content text color
	p.AddContent(contentHTML)

	return div, nil
}

// Add a search box to the page, actionURL is the url to use as a get action,
// buttonText is the text on the search button
func AddSearchBox(tag *browserspeak.Tag, actionURL, buttonText string) *browserspeak.Tag {

	div := tag.AddNewTag("div")
	div.AddAttr("id", "searchboxdiv")
	div.AddStyle("text-align", "right")
	div.AddStyle("display", "block")

	form := div.AddNewTag("form")
	form.AddAttr("id", "search")
	form.AddAttr("method", "get")
	form.AddAttr("action", actionURL)

	innerDiv := form.AddNewTag("div")
	innerDiv.AddAttr("id", "innerdiv")
	innerDiv.AddStyle("overflow", "hidden")
	innerDiv.AddStyle("padding-right", "0.5em")
	innerDiv.AddStyle("display", "inline-block")
	//innerDiv.AddStyle("background-color", "red")
	//innerDiv.AddStyle("display", "box")
	//innerDiv.AddStyle("box-align", "center")
	//innerDiv.AddStyle("display", "table-cell")
	//innerDiv.AddStyle("float", "left")

	inputText := innerDiv.AddNewTag("input")
	inputText.AddAttr("name", "q")
	inputText.AddAttr("size", "22")
	//inputText.AddStyle("position", "absolute")

	inputButton := form.AddNewTag("input")
	inputButton.AddStyle("margin-left", "0.4em")
	inputButton.AddStyle("float", "right")
	inputButton.AddAttr("type", "submit")
	inputButton.AddAttr("value", buttonText)
	inputButton.CustomSansSerif("Armata")
	//inputButton.AddStyle("vertical-align", "middle")
	//inputButton.AddStyle("top", "100px")
	//inputButton.AddStyle("position", "absolute")

	return div
}

func AddTitleBox(tag *browserspeak.Tag, title, subtitle string) *browserspeak.Tag {

	div := tag.AddNewTag("div")
	div.AddAttr("id", "titlebox")
	div.AddStyle("display", "block")
	div.AddStyle("position", "fixed")

	word1 := title
	word2 := ""
	if strings.Contains(title, " ") {
		word1 = strings.SplitN(title, " ", 2)[0]
		word2 = strings.SplitN(title, " ", 2)[1]
	}

	a := div.AddNewTag("a")
	a.AddAttr("id", "homelink")
	a.AddAttr("href", "/")
	a.AddStyle("text-decoration", "none")

	font0 := a.AddNewTag("font")
	font0.AddAttr("id", "whitetitle")
	font0.AddStyle("color", "white")
	font0.CustomSansSerif("Armata")
	font0.AddStyle("font-size", "2.0em")
	font0.AddStyle("font-weight", "bolder")
	font0.AddContent(word1)

	font1 := a.AddNewTag("font")
	font1.AddAttr("id", "bluetitle")
	font1.AddStyle("color", NICEBLUE)
	font1.CustomSansSerif("Armata")
	font1.AddStyle("font-size", "2.0em")
	font1.AddStyle("font-weight", "bold")
	font1.AddContent(word2)

	font2 := a.AddNewTag("font")
	font2.AddAttr("id", "graytitle")
	font2.AddStyle("font-size", "0.5em")
	font2.AddStyle("color", "#707070")
	font2.CustomSansSerif("Armata")
	font2.AddStyle("font-size", "1.25em")
	font2.AddStyle("font-weight", "normal")
	font2.AddContent(subtitle)

	return div
}

// Split a string at the colon into two strings
// If there's no colon, return the string and an empty string
func colonsplit(s string) (string, string) {
	if strings.Contains(s, ":") {
		sl := strings.SplitN(s, ":", 2)
		return sl[0], sl[1]
	}
	return s, ""
}

// Takes a page and a colon-separated string slice of text:url
func AddMenuBox(page *browserspeak.Page, links []string, darkBackgroundTexture string) (*browserspeak.Tag, error) {
	body, err := page.GetTag("body")
	if err != nil {
		return nil, err
	}

	div := body.AddNewTag("div")
	div.AddAttr("id", "menubox")
	div.AddStyle("display", "block")
	div.AddStyle("width", "100%")
	div.AddStyle("margin", "0")
	div.AddStyle("padding", "0.1em 0 0.2em 0")
	div.AddStyle("position", "absolute")
	div.AddStyle("top", "5em")
	div.AddStyle("left", "0")
	div.AddStyle("background-color", "#0c0c0c") // dark gray, fallback
	div.AddStyle("background", "url('"+darkBackgroundTexture+"')")
	div.AddStyle("position", "fixed")
	//div.AddStyle("-moz-box-shadow", "10px 10px 5px #606060")
	//div.AddStyle("-webkit-box-shadow", "10px 10px 5px #606060")
	div.AddStyle("box-shadow", "1px 3px 5px rgba(0,0,0, .8)")

	ul := div.AddNewTag("ul")
	ul.AddStyle("list-style-type", "none")
	ul.AddStyle("float", "left")
	ul.AddStyle("margin", "0")
	//ul.AddStyle("padding", "0")

	var a, li, sep *browserspeak.Tag
	var text, url string

	styleadded := false
	for i, text_url := range links {
		text, url = colonsplit(text_url)

		li = ul.AddNewTag("li")
		li.AddStyle("display", "inline")

		a = li.AddNewTag("a")
		a.AddAttr("id", "menulink")
		a.AddAttr("href", url)
		if !styleadded {
			a.AddStyle("font-weight", "bold")
			a.AddStyle("color", "#303030")
			a.AddStyle("text-decoration", "none")
			//a.AddStyle("padding", "8px 1.2em")
			//a.AddStyle("margin", "0")
			a.CustomSansSerif("Armata")
			//a.AddStyle("display", "block")
			//a.AddStyle("width", "60px")
			styleadded = true
		}
		a.AddContent(text)

		// For every element, but not the last one
		if i < (len(links) - 1) {
			// Insert a '|' character in a div
			sep = li.AddNewTag("div")
			sep.AddContent("|")
			sep.AddStyle("display", "inline")
			sep.AddStyle("color", "#a0a0a0")
		}
	}

	return div, nil
}
