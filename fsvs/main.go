package main

import (
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "dir/index.html")
}

func main() {
	http.Handle("/assets/", http.FileServer(http.Dir("dir1")))
	http.Handle("/bower_components/", http.FileServer(http.Dir("dir2")))
	http.HandleFunc("/", RootHandler)
	http.ListenAndServe(":8088", nil)
}
