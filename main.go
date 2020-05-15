package main

import (
    "fmt"
	"net/http"
	"html/template"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"encoding/json"
)
type User struct {
	Id 		int
    Name   string 
    Surname string
}

func logData(w http.ResponseWriter, r *http.Request){
	fmt.Println("w", w)
	fmt.Println("r", r)
	fmt.Println("*********************")
}

func dbConn() (db *sql.DB) {
    dbDriver 	:= "mysql"
	dbName 		:= "go"
    dbUser 		:= "root"
    dbPass 		:= "2121"


    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    if err != nil {
        panic(err.Error())
	}
	fmt.Println("db connection successfully")
	log.Println("db connect log")
    return db
}

func getData(w http.ResponseWriter, r *http.Request){
	logData(w, r)
	db := dbConn()
	search := r.FormValue("search")
	results, err := db.Query("SELECT * FROM users WHERE CONCAT_WS(' ', name, surname) LIKE ?", "%"+search+"%")

	user := User{}
	res := []User{}
    for results.Next() {
        err = results.Scan(&user.Id, &user.Name, &user.Surname)
        if err != nil {
            panic(err.Error())
		}
		res = append(res, user)
		log.Println(res)
	}

	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func index(w http.ResponseWriter, r *http.Request){
	db := dbConn()
	results, err := db.Query("SELECT * FROM users")

	user := User{}
	res := []User{}
    for results.Next() {
        err = results.Scan(&user.Id, &user.Name, &user.Surname)
        if err != nil {
            panic(err.Error())
		}
		res = append(res, user)
		log.Println(res)
	}

	
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, res)
}

func main() {
	fmt.Println("Server starting http://localhost:8080")
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
    http.HandleFunc("/", index)
    http.HandleFunc("/getData", getData)
    http.ListenAndServe(":8080", nil)
}