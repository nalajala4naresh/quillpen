package sessionManager

import "github.com/gorilla/sessions"

var SessionName = "quill"
var SessionUserId = "userid"
var SessionIsAuthenticated = "loggedin"
var Store = sessions.NewCookieStore([]byte("Quillpen!!"))


