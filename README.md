Google App Engine Guestbook for Golang with Bootstrap
=====================================================

This project is just a first attempt at using google app engine with go.  The goal was to use templates and static assets to get a sense of how to lay out an app.  This app basically follows the go app engine tutorial, except it adds a authOnly handler to force users to log in to access the guestbook.

To get started:

* [Read about setting up GAE dev environment in go](https://developers.google.com/appengine/docs/go/gettingstarted/devenvironment)
* git clone git@github.com:nickdufresne/gae-guestbook-bootstrap-go.git
* ensure goapp bin is in your path
* goapp serve gae-guestbook-bootstrap-go
* navigate to [localhost:8080](http://localhost:8080/)
