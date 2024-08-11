// main.go

package main

import (
    "encoding/json"
    "fmt"
    "html/template"
    "net/http"
    "os"

    "github.com/yorologo/GoPhrasesGenerator/image"
)

func imageHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("templates/index.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    tmpl.Execute(w, nil)
}

func loadMoreImagesHandler(w http.ResponseWriter, r *http.Request) {
    imagePaths, err := image.GenerateImages(10)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(imagePaths)
}

func main() {
    if _, err := os.Stat("img"); os.IsNotExist(err) {
        os.Mkdir("img", 0755)
    }

    http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))
    http.HandleFunc("/", imageHandler)
    http.HandleFunc("/load-more-images", loadMoreImagesHandler)

    fmt.Println("Server started at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
