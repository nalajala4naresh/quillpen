package main

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed index.html
var Templates embed.FS

type Mpage struct{

	Countries []string
	Data []Post

}

type Post struct {
    Title string
	Content string

}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", Index)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()

}

func Index(resp http.ResponseWriter, req *http.Request) {

	t := template.Must(template.ParseFS(Templates,"index.html"))
	var data []Post
    opost := Post{
		Title: "Coming Soon",
		Content: "We are working on building the next generation blog for the locals",
	}
	data = append(data,opost)

	page_data := Mpage{ Countries: []string{"India","USA","Canada","UK"}, Data: data}
	t.Execute(resp, page_data)
}

