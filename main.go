package main

import (
	"embed"
	"fmt"
	"log"
	"os/exec"

	// "net/http"

	"github.com/webview/webview"
)

//go:embed app/*
var app embed.FS

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println(w, "You requested: %s\n", r.URL.Path)
// }

func main() {
	// http.HandleFunc("/", handler)

	// fmt.Println("go listening: 3030")
	// http.ListenAndServe(":3030", nil)

	serve := exec.Command("app/php", "-S", "localhost:3000", "-t", "app")
	err = serve.Start()
	if err != nil {
		log.Fatal(err)
	}
	defer serve.Process.Kill()

	view := webview.New(true)
	defer view.Destroy()
	view.SetTitle("Embed")
	view.SetSize(800, 600, webview.Hint(webview.HintNone))
	view.Navigate("http://localhost:3000")
	view.Run()
}