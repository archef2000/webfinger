package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type User struct {
	Subject string `json:"subject"`
	Links   []Link `json:"links"`
}

type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

func WebFingerHandler(links []Link) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resource := r.URL.Query().Get("resource")

		if resource == "" {
			http.Error(w, "Missing 'resource' parameter", http.StatusBadRequest)
			return
		}
		fmt.Println("WebFinger got accessed with resource: ", resource)
		fmt.Println(r.URL.Query().Get("rel"))
		user := User{
			Subject: resource,
			Links:   links,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(user); err != nil {
			fmt.Println("Error encoding JSON: " + err.Error())
			http.Error(w, "Error encoding JSON: see server logs", http.StatusInternalServerError)
			return
		}
	}
}

func getLinks() []Link {
	pattern := regexp.MustCompile("^LINKS_\\d_REL$")
	var links []Link
	envs := os.Environ()
	for _, env := range envs {
		pair := strings.SplitN(env, "=", 2)
		key := pair[0]
		if pattern.MatchString(key) {
			parts := strings.Split(key, "_")
			index := parts[1]
			hrefEnv := "LINKS_" + index + "_HREF"
			href := os.Getenv(hrefEnv)
			if href == "" {
				fmt.Printf("Found value for 'LINKS_%s_REL' but '%s' is not set.\n", index, hrefEnv)
				os.Exit(1)
			}
			value := pair[1]
			link := Link{Rel: value, Href: href}
			links = append(links, link)
		}
	}
	return links
}

func main() {
	router := mux.NewRouter()

	links := getLinks()
	router.HandleFunc("/.well-known/webfinger", WebFingerHandler(links)).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	serverAddr := fmt.Sprintf(":%s", port)
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	fmt.Printf("Server is running on http://localhost:%s\n", port)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
