package main

import (
	"errors"
	"math/rand"
	"time"
	"strings"

	"github.com/garyburd/redigo/redis"
	. "github.com/xyproto/browserspeak"
	"github.com/xyproto/web"
)

const (
	ONLY_LOGIN      = "100"
	ONLY_LOGOUT     = "010"
	ONLY_REGISTER   = "001"
	EXCEPT_LOGIN    = "011"
	EXCEPT_LOGOUT   = "101"
	EXCEPT_REGISTER = "110"
	NOTHING         = "000"

	MINIMUM_CONFIRMATION_CODE_LENGTH = 20
	USERNAME_ALLOWED_LETTERS = "abcdefghijklmnopqrstuvwxyzæøåABCDEFGHIJKLMNOPQRSTUVWXYZÆØÅ_0123456789"
)

type UserState struct {
	// see: http://redis.io/topics/data-types
	users      *RedisHashMap    // Hash map of users, with several different fields per user ("loggedin", "confirmed", "email" etc)
	usernames  *RedisSet        // A list of all usernames, for easy enumeration
	unconfirmed *RedisSet       // A list of unconfirmed usernames, for easy enumeration
	connection redis.Conn
}

func InitUserSystem(connection redis.Conn) *UserState {

	// For the secure cookies
	// This must happen before the random seeding, or 
	// else people will have to log in again after every server restart
	web.Config.CookieSecret = RandomCookieFriendlyString(30)

	rand.Seed(time.Now().UnixNano())

	// For the database
	state := new(UserState)
	state.users = NewRedisHashMap(connection, "users")
	state.usernames = NewRedisSet(connection, "usernames")
	state.unconfirmed = NewRedisSet(connection, "unconfirmed")
	state.connection = connection

	return state
}

// TODO: Rethink this. Use templates for Login/Logout button?
// Generate "1" or "0" values for showing the login, logout or register menus,
// depending on the cookie status and UserState
func GenerateShowLoginLogoutRegister(state *UserState) SimpleContextHandle {
	return func(ctx *web.Context) string {
		if username := GetBrowserUsername(ctx); username != "" {
			//print("USERNAME", username)
			// Has a username stored in the browser
			if state.LoggedIn(username) {
				// Ok, logged in to the system + login cookie in the browser
				// Only present the "Logout" menu
				return ONLY_LOGOUT
			} else {
				// Has a login cookie, but is not logged in.
				// Keep the browser cookie (could be tempting to remove it)
				// Present only the "Login" menu
				//return "100"
				// Present both "Login" and "Register", just in case it's a new user
				// in the same browser.
				return EXCEPT_LOGOUT
			}
		} else {
			// Does not have a username stored in the browser
			// Present the "Register" and "Login" menu
			return EXCEPT_LOGOUT
		}
		// Everything went wrong, should never reach this point
		return NOTHING
	}
}

// TODO: Don't return false if there is an error, the user may exist
func (state *UserState) HasUser(username string) bool {
	val, err := state.usernames.Has(username)
	if err != nil {
		return false
	}
	return val
}

// Creates a user without doing ANY checks
func AddUserUnchecked(state *UserState, username, password, email string) {
	// Add the user
	state.usernames.Add(username)

	// Add password and email
	state.users.Set(username, "password", password)
	state.users.Set(username, "email", email)

	// Addditional fields
	state.users.Set(username, "loggedin", "false")
	state.users.Set(username, "confirmed", "false")
}

func IsConfirmed(state *UserState, username string) bool {
	if !state.HasUser(username) {
		return false
	}
	confirmed, err := state.users.Get(username, "confirmed")
	if err != nil {
		return false
	}
	return TruthValue(confirmed)
}

func IsLoggedIn(state *UserState, username string) bool {
	if !state.HasUser(username) {
		return false
	}
	loggedin, err := state.users.Get(username, "loggedin")
	if err != nil {
		return false
	}
	return TruthValue(loggedin)
}

func CorrectPassword(state *UserState, username, password string) bool {
	hashedPassword, err := state.users.Get(username, "password")
	if err != nil {
		return false
	}
	if hashedPassword == HashPassword(password) {
		return true
	}
	return false
}

// Goes through all the secrets of all the unconfirmed users
// and checks if this secret already is in use
func AlreadyHasSecret(state *UserState, secret string) bool {
	unconfirmedUsernames, err := state.usernames.GetAll()
	if err != nil {
		return false
	}
	for _, aUsername := range unconfirmedUsernames {
		aSecret, err := state.users.Get(aUsername, "secret")
		if err != nil {
			// TODO: Inconsistent user, log this
			continue
		}
		if secret == aSecret {
			// Found it
			return true
		}
	}
	return false
}

// Create a user by adding the username to the list of usernames
func GenerateConfirmUser(state *UserState) WebHandle {
	return func(ctx *web.Context, val string) string {
		secret := val

		unconfirmedUsernames, err := state.usernames.GetAll()
		if err != nil {
			return MessageOKurl("Confirmation", "All users are confirmed already.", "/register")
		}

		// TODO: Only generate unique secrets

		// Find the username by looking up the secret on unconfirmed users
		username := ""
		for _, aUsername := range unconfirmedUsernames {
			aSecret, err := state.users.Get(aUsername, "secret")
			if err != nil {
				// TODO: Inconsistent user! Log this.
				continue
			}
			if secret == aSecret {
				// Found the right user
				username = aUsername
				break
			}
		}

		// Check that the user is there
		if username == "" {
			// Say "no longer" because we don't care about people that just try random confirmation links
			return MessageOKurl("Confirmation", "The confirmation link is no longer valid.", "/register")
		}
		if !state.HasUser(username) {
			return MessageOKurl("Confirmation", "The user you wish to confirm does not exist anymore.", "/register")
		}

		// Remove from the list of unconfirmed usernames
		state.unconfirmed.Del(username)
		// Remove the secret from the user
		state.users.Del(username, "secret")

		// Mark user as confirmed
		state.users.Set(username, "confirmed", "true")

		return MessageOKurl("Confirmation", "Thank you " + username + ", you can now log in.", "/login")
	}
}

// Create a user by adding the username to the list of usernames
func GenerateRemoveUser(state *UserState) WebHandle {
	return func(ctx *web.Context, val string) string {
		if val == "" {
			return "Can't remove blank user"
		}
		if !state.HasUser(val) {
			return "user " + val + " doesn't exists, could not remove"
		}

		// Remove the user
		state.usernames.Del(val)

		// Remove additional data as well
		state.users.Del(val, "loggedin")

		return "OK, user " + val + " removed"
	}
}

// Log in a user by changing the loggedin value
func GenerateLoginUser(state *UserState) WebHandle {
	return func(ctx *web.Context, val string) string {
		// Fetch password from ctx
		password, found := ctx.Params["password"]
		if !found {
			return MessageOKback("Login", "Can't log in without a password.")
		}
		username := val
		if username == "" {
			return MessageOKback("Login", "Can't log in with a blank username.")
		}
		if !state.HasUser(username) {
			return MessageOKback("Login", "User " + username + " does not exist, could not log in.")
		}

		if !IsConfirmed(state, username) {
			return MessageOKback("Login", "The email for " + username + " has not been confirmed, check your email and follow the link.")
		}

		// TODO: Hash password, check with hash from database

		if !CorrectPassword(state, username, password) {
			return MessageOKback("Login", "Wrong password.")
		}

		// Log in the user by changing the database and setting a secure cookie
		state.users.Set(username, "loggedin", "true")
		state.SetBrowserUsername(ctx, username)

		// TODO: Use a welcoming messageOK where the user can see when he/she last logged in and from which host

		// TODO: Then redirect to the page the user was at before logging in
		ctx.SetHeader("Refresh", "0; url=/", true)

		return ""
	}
}

func HashPassword(password string) string {
	// TODO: Implement actual hashing, with salt
	return "abc123" + password + "abc123"
}

// Register a new user
func GenerateRegisterUser(state *UserState) WebHandle {
	return func(ctx *web.Context, val string) string {

		// Password checks
		password1, found := ctx.Params["password1"]
		if password1 == "" || !found {
			return MessageOKback("Register", "Can't register without a password.")
		}
		password2, found := ctx.Params["password2"]
		if password2 == "" || !found {
			return MessageOKback("Register", "Please confirm the password by typing it in twice.")
		}
		if password1 != password2 {
			return MessageOKback("Register", "The password and confirmation password must be equal.")
		}

		// Email checks
		email, found := ctx.Params["email"]
		if !found {
			return MessageOKback("Register", "Can't register without an email address.")
		}
		// must have @ and ., but no " "
		if !strings.Contains(email, "@") || !strings.Contains(email, ".") || strings.Contains(email, " ") {
			return MessageOKback("Register", "Please use a valid email address.")
		}

		// Username checks
		username := val
		if username == "" {
			return MessageOKback("Register", "Can't register without a username.")
		}
		if state.HasUser(username) {
			return MessageOKback("Register", "That user already exists, try another username.")
		}
		// Only some letters are allowed
		NEXT: for _, letter := range username {
			for _, allowedLetter := range USERNAME_ALLOWED_LETTERS {
				if letter == allowedLetter {
					continue NEXT
				}
			}
			return MessageOKback("Register", "Only a-å, A-Å, 0-9 and _ are allowed in usernames.")
		}
		if username == password1 {
			return MessageOKback("Register", "Username and password must be different, try another password.")
		}

		// Register the user
		password := HashPassword(password1)
		AddUserUnchecked(state, username, password, email)
		state.users.Set(username, "confirmed", "false")

		// The confirmation code must be a minimum of 8 letters long
		length := MINIMUM_CONFIRMATION_CODE_LENGTH
		secretConfirmationCode := RandomHumanFriendlyString(length)
		for AlreadyHasSecret(state, secretConfirmationCode) {
			// Increase the length of the secret random string every time there is a collision
			length++
			secretConfirmationCode = RandomHumanFriendlyString(length)
			if length > 100 {
				// Something is seriously wrong if this happens
				// TODO: Log this and sysexit
			}
		}

		// Send confirmation email
		ConfirmationEmail("archlinux.no", "https://archlinux.no/confirm/" + secretConfirmationCode, username, email)

		// Register the need to be confirmed
		state.unconfirmed.Add(username)
		state.users.Set(username, "secret", secretConfirmationCode)

		// Redirect
		//ctx.SetHeader("Refresh", "0; url=/login", true)

		return MessageOKurl("Registration complete", "Thanks for registering, the confirmation e-mail has been sent.", "/login")
	}
}

// Log out a user by changing the loggedin value
func GenerateLogoutCurrentUser(state *UserState) SimpleContextHandle {
	return func(ctx *web.Context) string {
		username := GetBrowserUsername(ctx)
		if username == "" {
			return MessageOKback("Logout", "No user to log out")
		}
		if !state.HasUser(username) {
			return MessageOKback("Logout", "user " + username + " does not exist, could not log out")
		}

		// TODO: Check if the user is logged in already

		// Log out the user by changing the database, the cookie can stay
		state.users.Set(username, "loggedin", "false")

		//return "OK, user " + username + " logged out"

		// TODO: Redirect to the page the user was at before logging out
		//ctx.SetHeader("Refresh", "0; url=/", true)

		//return ""

		return MessageOKurl("Logout", username + " is now logged out. Hope to see you soon!", "/login")
	}
}

func GenerateGetAllUsernames(state *UserState) SimpleWebHandle {
	return func(val string) string {
		s := ""
		usernames, err := state.usernames.GetAll()
		if err == nil {
			for _, val := range usernames {
				s += val + "<br />"
			}
		}
		return MessageOKback("Usernames", s)
	}
}

func GenerateStatus(state *UserState) SimpleWebHandle {
	return func(val string) string {
		username := val
		if username == "" {
			return MessageOKback("Status", "No username given")
		}
		if !state.HasUser(username) {
			return MessageOKback("Status", username + " does not exist")
		}
		loggedinStatus := "not logged in"
		if IsLoggedIn(state, username) {
			loggedinStatus = "logged in"
		}
		confirmStatus := "email has not been confirmed"
		if IsConfirmed(state, username) {
			confirmStatus = "email has been confirmed"
		}
		return MessageOKback("Status", username + " is " + loggedinStatus + " and " + confirmStatus)
	}
}

func GenerateStatusCurrentUser(state *UserState) SimpleContextHandle {
	return func(ctx *web.Context) string {
		username := GetBrowserUsername(ctx)
		if username == "" {
			return MessageOKback("Current user status", "No user logged in")
		}
		if !state.HasUser(username) {
			return MessageOKback("Current user status", username + " does not exist")
		}
		if !(state.LoggedIn(username)) {
			return MessageOKback("Current user status", "User " + username + " is not logged in")
		}
		return MessageOKback("Current user status", "User " + username + " is logged in")
	}
}

// Checks if the given username is logged in or not
func (state *UserState) LoggedIn(username string) bool {
	if !state.HasUser(username) {
		return false
	}
	status, err := state.users.Get(username, "loggedin")
	if err != nil {
		return false
	}
	return TruthValue(status)
}

func GenerateGetCookie(state *UserState) SimpleContextHandle {
	return func(ctx *web.Context) string {
		username := GetBrowserUsername(ctx)
		//username, _ := ctx.GetSecureCookie("user")
		return "Cookie: username = " + username // + " err: " + fmt.Sprintf("%v", exists) + " val: " + val
	}
}

// Gets the username that is stored in a cookie in the browser, if available
func GetBrowserUsername(ctx *web.Context) string {
	username, _ := ctx.GetSecureCookie("user")
	return username
}

func (state *UserState) SetBrowserUsername(ctx *web.Context, username string) error {
	if username == "" {
		return errors.New("Can't set cookie for empty username")
	}
	if !state.HasUser(username) {
		return errors.New("Can't store cookie for non-existsing user")
	}
	// Create a cookie that lasts for one hour,
	// this is the equivivalent of a session for a given username
	ctx.SetSecureCookiePath("user", username, 3600, "/")
	//"Cookie stored: user = " + username + "."
	return nil
}

// NB! Set the cookie at / for it to work in the paths underneath!
func GenerateSetCookie(state *UserState) WebHandle {
	return func(ctx *web.Context, val string) string {
		username := val
		if username == "" {
			return "Can't set cookie for empty username"
		}
		if !state.HasUser(username) {
			return "Can't store cookie for non-existsing user"
		}
		// Create a cookie that lasts for one hour,
		// this is the equivivalent of a session for a given username
		ctx.SetSecureCookiePath("user", username, 3600, "/")
		return "Cookie stored: user = " + username + "."
	}
}

// TODO: RESTful services?
func ServeUserSystem(connection redis.Conn) *UserState {
	state := InitUserSystem(connection)

	web.Post("/register/(.*)", GenerateRegisterUser(state))
	web.Post("/login/(.*)", GenerateLoginUser(state))
	web.Get("/logout", GenerateLogoutCurrentUser(state))
	web.Get("/confirm/(.*)", GenerateConfirmUser(state))

	// TODO: debug pages, comment out
	web.Get("/status", GenerateStatusCurrentUser(state))
	web.Get("/status/(.*)", GenerateStatus(state))
	web.Get("/remove/(.*)", GenerateRemoveUser(state))
	web.Get("/users/(.*)", GenerateGetAllUsernames(state))
	web.Get("/cookie/get", GenerateGetCookie(state))
	web.Get("/cookie/set/(.*)", GenerateSetCookie(state))

	return state
}
