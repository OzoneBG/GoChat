package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/stretchr/gomniauth/common"

	"github.com/stretchr/objx"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
)

var env string

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	data := map[string]interface{}{
		"Host": r.Host,
	}

	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w, data)
}

func getGoogleProvider() common.Provider {
	clientID := "199014956970-cubjkv0gs3qd37kvl2k29t2rcd33tl1j.apps.googleusercontent.com"
	secret := "_a2tOls2gp-F8KrAjIMpEbLW"

	var callbackURL string
	if env == "production" {
		callbackURL = "http://gochatr.ddns.net/auth/callback/google"
	} else if env == "development" {
		callbackURL = "http://localhost:8080/auth/callback/google"
	}

	return google.New(clientID, secret, callbackURL)
}

func main() {
	var addr = flag.String("addr", ":80", "The addr of the application")
	envFlag := flag.String("env", "production", "Current environment")
	flag.Parse()

	env = *envFlag

	gomniauth.SetSecurityKey("myverysecretkey")
	gomniauth.WithProviders(
		getGoogleProvider(),
	)

	r := newRoom(UseFileSystemAvatar)
	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/avatars/",
		http.StripPrefix("/avatars/",
			http.FileServer(http.Dir("./avatars"))))

	go r.run()

	log.Println("Starting the web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
