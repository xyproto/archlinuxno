package main

// Creates an anonymous function
func fn(source string) string {
	return "function() { " + source + " }"
}

func quote(src string) string {
	return "\"" + src + "\""
}

// Creates an event
func event(tagname, event, source string) string {
	return "$(" + quote(tagname) + ")." + event + "(" + fn(source) + ");"
}

// Call a method on a tag
func method(tagname, methodname, value string) string {
	return "$(" + quote(tagname) + ")." + methodname + "(" + value + ");"
}

// Call a method on a tag that takes a string as a parameter
func methodString(tagname, methodname, value string) string {
	return method(tagname, methodname, quote(value))
}

// Run code when the document is ready
func OnDocumentReady(source string) string {
	return "$(document).ready(" + fn(source) + ");"
}

// Display an intruding message
func Alert(msg string) string {
	return "alert(" + quote(msg) + ");"
}

// When a tag is clicked at
func OnClick(tagname, source string) string {
	return event(tagname, "click", source)
}

func SetText(tagname, text string) string {
	return methodString(tagname, "text", text)
}

func SetHTML(tagname, html string) string {
	return method(tagname, "html", html)
}

func SetValue(tagname, val string) string {
	return methodString(tagname, "val", val)
}

func SetRawValue(tagname, val string) string {
	return method(tagname, "val", val)
}

func Hide(tagname string) string {
	return "$(" + quote(tagname) + ").hide();"
}

func HideAnimated(tagname string) string {
	return "$(" + quote(tagname) + ").hide('fast');"
}

func Show(tagname string) string {
	return "$(" + quote(tagname) + ").show();"
}

func ShowAnimated(tagname string) string {
	return "$(" + quote(tagname) + ").show('fast');"
}

// Same as Show, but set display to inline instead of block
func ShowInline(tagname string) string {
	return "$(" + quote(tagname) + ").css('display', 'inline');"
}

// Same as Show, but set display to inline instead of block
func ShowInlineAnimated(tagname string) string {
	return ShowInline(tagname) + Hide(tagname) + ShowAnimated(tagname)
}

func Load(tagname, url string) string {
	return methodString(tagname, "load", url)
}

// Hide a tag if booleanURL doesn't return "1" (true)
func HideIfNot(booleanURL, tagname string) string {
	return "$.get(" + quote(booleanURL) + ", function(data) { if (data != \"1\") {" + Hide(tagname) + "}; });"
}

// Optimized function for login, logout and register
func HideIfNotLoginLogoutRegister(threeBooleanURL, logintag, logouttag, registertag string) string {
	src := "$.get(" + quote(threeBooleanURL) + ", function(data) {"
	// TODO: See what happens if data < 3 length
	src += "if (data[0] != \"1\") {" + Hide(logintag) + "};"
	src += "if (data[1] != \"1\") {" + Hide(logouttag) + "};"
	src += "if (data[2] != \"1\") {" + Hide(registertag) + "};"
	src += "});"
	return src
}

// Optimized function for login, logout and register
func ShowIfLoginLogoutRegister(threeBooleanURL, logintag, logouttag, registertag string) string {
	src := "$.get(" + quote(threeBooleanURL) + ", function(data) {"
	// TODO: See what happens if data < 3 length
	src += "if (data[0] == \"1\") {" + ShowInline(logintag) + "};"
	src += "if (data[1] == \"1\") {" + ShowInline(logouttag) + "};"
	src += "if (data[2] == \"1\") {" + ShowInline(registertag) + "};"
	src += "});"
	return src
}

// Returns html to run javascript once the document is ready
// Returns an empty string if there is no javascript to run.
func BodyJS(source string) string {
	if source != "" {
		return "<script type=\"text/javascript\">" + OnDocumentReady(source) + "</script>"
	}
	return ""
}
