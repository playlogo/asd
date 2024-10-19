package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"
	"time"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	body, err := loadPage("index")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, applyTemplate(body))
}

func getSearch(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("q")
	fmt.Printf("got /search request with query %s", searchQuery)
	body, err := loadPage("search")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	products := []Product{}

	if searchQuery != "" {
		/* !TODO fetch products from DB */
	}

	products = []Product{
		{
			Name:  "Test",
			Price: 4079,
			Image: "https://image-placeholder.com/images/actual-size/200x200.png",
		},
	}

	io.WriteString(w, applySearchTemplate(body, products))
}

func loadPage(name string) (string, error) {
	filename := name + ".html"
	body, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

type Product struct {
	Name  string
	Price int
	Image string
}

func applySearchTemplate(input string, products []Product) string {
	type Data struct {
		Products []Product
	}

	data := Data{products}

	tmpl, err := template.New("search").Parse(input)
	if err != nil {
		panic(err)
	}
	var result bytes.Buffer
	err = tmpl.Execute(&result, data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("server closed\n")
	return ""
}

func applyTemplate(input string) string {

	type Data struct {
		Time string
	}
	t := time.Now()

	sweaters := Data{t.Format("2006.01.02 - 15:04:05")}
	tmpl, err := template.New("root").Parse(input)
	if err != nil {
		panic(err)
	}
	var result bytes.Buffer
	err = tmpl.Execute(&result, sweaters)
	if err != nil {
		panic(err)
	}
	return result.String()
}

func main() {
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/search", getSearch)

	err := http.ListenAndServe(":3333", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}
