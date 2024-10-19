package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

type Product struct {
	Name  string
	Price int
	Image string
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	body, err := loadPage("index")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, applyTemplate(body))
}

func initDB() {
	dsn := "root:root@tcp(localhost:3306)/check24"

	db, err := sql.Open("mysql", dsn)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Products (
	 	Name VARCHAR(255) PRIMARY KEY NOT NULL,
	 	Price INT NOT NULL,
	 	Image VARCHAR(255) NOT NULL
	 )`)

	if err != nil {
		panic(err)
	}
	file, err := os.Open("dump.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var products []Product
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&products)
	if err != nil {
		panic(err)
	}

	for _, product := range products {
		_, err := db.Exec(
			"INSERT INTO Products (Name, Price, Image) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE Price=VALUES(Price), Image=VALUES(Image)", product.Name, product.Price, product.Image)
		if err != nil {
			panic(err)
		}
	}
}

func queryProducts(q string) []Product {

	dsn := "root:root@tcp(localhost:3306)/check24"
	// /run/mysqld/mysqld.sock
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	rows, _ := db.Query("SELECT Name, Price, Image FROM Products WHERE Name LIKE ? " + "%" + q + "%")

	var products []Product

	// Iterate through the result set
	for rows.Next() {
		var product Product

		// Scan the result into the Product struct
		rows.Scan(&product.Name, &product.Price, &product.Image)
		fmt.Println(product.Name)

		// Append the product to the slice
		products = append(products, product)
	}

	return products
}

func getSearch(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("q")
	fmt.Printf("got /search request with query %s", searchQuery)
	body, err := loadPage("search")
	fmt.Println("a")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("b")

	products := []Product{}

	if searchQuery != "" {
		products = queryProducts(searchQuery)
	}
	fmt.Println("c")

	fmt.Println("d")

	io.WriteString(w, applySearchTemplate(body, products))
	fmt.Println("e")

}

func loadPage(name string) (string, error) {
	filename := name + ".html"
	body, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(body), nil
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
	return result.String()
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
	fmt.Println("asd")

	initDB()
	fmt.Println("asd")
	err := http.ListenAndServe(":3333", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)

	}
}
