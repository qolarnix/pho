package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"io"
	// "net/http"

	"github.com/webview/webview"
)

//go:embed bin/php
var php []byte

//go:embed app/*
var app embed.FS

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println(w, "You requested: %s\n", r.URL.Path)
// }

func main() {
	// http.HandleFunc("/", handler)

	// fmt.Println("go listening: 3030")
	// http.ListenAndServe(":3030", nil)

	mountDir, err := os.MkdirTemp("", "app")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(mountDir)

	binPath := filepath.Join(mountDir, "php")
	err = os.WriteFile(binPath, php, 0700)
	if err != nil {
		log.Fatal(err)
	}

	err = fs.WalkDir(app, "app", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		file, err := app.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		stripPrefix := strings.TrimPrefix(path, "app/")
		fmt.Println(stripPrefix)

		writeFile, err := os.Create(filepath.Join(mountDir, stripPrefix))
		if err != nil {
			return err
		}
		defer writeFile.Close()

		_, err = io.Copy(writeFile, file)
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(binPath)
	fmt.Println(mountDir)

	serve := exec.Command(binPath, "-S", "localhost:3000", "-t", mountDir)
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