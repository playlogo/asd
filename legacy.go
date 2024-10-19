/*package main

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
	body, err := loadPage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, applyTemplate(body))
}

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

func loadPage() (string, error) {
	filename := "index.html"
	body, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func applyTemplate(input string) string {

	type Data struct {
		Time string
	}
	t := time.Now()

	sweaters := Data{t.Format("2006.01.02 - 15:04:05")}
	tmpl, err := template.New("test").Parse(input)
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
	http.HandleFunc("/hello", getHello)

	err := http.ListenAndServe(":3333", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}
*/