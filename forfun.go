package main

// Handle various request that just appeared in the log

import (
	"github.com/hoisie/web"
)

// Honeypot?
func ServeForFun() {
	bogus := []string{"/signup", "/wp-login.php", "/join.php", "/register.php", "/profile.php", "/user/register/", "/tools/quicklogin.one", "/sign_up.html", "/profile.php", "/ucp.php", "/account/register.php", "/join_form.php", "/tiki-register.php", "/YaBB.cgi/", "/YaBB.pl/", "/member/register", "/signup.php", "/blogs/load/recent", "/member/join.php"}
	for _, location := range bogus {
		web.Get(location, Hello)
	}

	bogusParam := []string{"/index.php", "/viewtopic.php"}
	for _, location := range bogusParam {
		web.Get(location, ParamExample)
	}
}
