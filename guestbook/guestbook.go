package guestbook

import (
	"html/template"
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

type Greeting struct {
	Author  string
	Content string
	Date    time.Time
}

type GreetingPage struct {
	User       *user.User
	SignOutURL string
	Greetings  []Greeting
}

func authOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)
		u := user.Current(c)
		if u == nil {
			url, _ := user.LoginURL(c, "/")
			if err := signinTemplate.Execute(w, url); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		next(w, r)
	}
}

// guestbookKey returns the key used for all guestbook entries.
func guestbookKey(c appengine.Context) *datastore.Key {
	// The string "default_guestbook" here could be varied to have multiple guestbooks.
	return datastore.NewKey(c, "Guestbook", "default_guestbook", 0, nil)
}

func root(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	if u == nil {
		panic("User should not be nil")
	}
	// Ancestor queries, as shown here, are strongly consistent with the High
	// Replication Datastore. Queries that span entity groups are eventually
	// consistent. If we omitted the .Ancestor from this query there would be
	// a slight chance that Greeting that had just been written would not
	// show up in a query.
	q := datastore.NewQuery("Greeting").Ancestor(guestbookKey(c)).Order("-Date").Limit(10)
	greetings := make([]Greeting, 0, 10)
	if _, err := q.GetAll(c, &greetings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url, _ := user.LogoutURL(c, "/")
	//fmt.Fprintf(w, `Welcome, %s! (<a href="%s">sign out</a>)`, u, url)

	page := GreetingPage{
		User:       u,
		SignOutURL: url,
		Greetings:  greetings,
	}

	if err := guestbookTemplate.Execute(w, page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func sign(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	g := Greeting{
		Content: r.FormValue("content"),
		Date:    time.Now(),
	}
	if u := user.Current(c); u != nil {
		g.Author = u.String()
	}
	// We set the same parent key on every Greeting entity to ensure each Greeting
	// is in the same entity group. Queries across the single entity group
	// will be consistent. However, the write rate to a single entity group
	// should be limited to ~1/second.
	key := datastore.NewIncompleteKey(c, "Greeting", guestbookKey(c))
	_, err := datastore.Put(c, key, &g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

var guestbookTemplate *template.Template
var signinTemplate *template.Template

func init() {
	http.HandleFunc("/", authOnly(root))
	http.HandleFunc("/sign", authOnly(sign))

	guestbookTemplate = template.Must(template.ParseFiles("tmpl/index.tmpl"))
	signinTemplate = template.Must(template.ParseFiles("tmpl/signin.tmpl"))
}
