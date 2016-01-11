package main

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"encoding/json"
)

type Person struct {
	Id   int `json:"id"`
	Name string `json:"name"`
	Age  int `json:"age"`
}

func (p Person) String() string {
	return fmt.Sprintf("[id=%v, name=%v, age=%v]", p.Id, p.Name, p.Age)
}

func (p Person) Smt() int {
	return p.Age * 2
}

func initDatabse() (*sql.DB, error) {
	db, err := sql.Open("postgres", "user=yourcode dbname=yourcode sslmode=disable")
	if err != nil {
		return nil, err
	}
	return db, err
}

func findPerson(minAge int, db *sql.DB) (*Person, error) {
	var person Person
	err := db.QueryRow("SELECT id, name, age FROM people WHERE age > $1", minAge).Scan(&person.Id, &person.Name, &person.Age)
	if err != nil {
		return nil, err
	}
	return &person, err
}

func listPeopleHandler(w http.ResponseWriter, r *http.Request) {
	people := []Person{Person{Id: 1, Name: "Michal", Age: 26}, Person{Id: 1, Name: "Michal", Age: 26}}
	json, err := json.Marshal(people)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func listPeople(db *sql.DB) ([]Person, error) {
	people := []Person{}
	rows, err := db.Query("SELECT id, name, age FROM people")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var person Person
		err = rows.Scan(&person.Id, &person.Name, &person.Age)
		if err != nil {
			return nil, err
		}
		people = append(people, person)
	}
	return people, err
}

func listDbPeopleHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	people, err := listPeople(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json, err := json.Marshal(people)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func makeHandlerHelper(db *sql.DB, fn func(http.ResponseWriter, *http.Request, *sql.DB)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, db)
	}
}

func makeHandlerClosure(db *sql.DB) func(func(http.ResponseWriter, *http.Request, *sql.DB)) http.HandlerFunc {
	return func(fn func(http.ResponseWriter, *http.Request, *sql.DB)) http.HandlerFunc {
		return makeHandlerHelper(db, fn)
	}
}

func startServer(db *sql.DB) {
	log.Println("Start listening on port :8080")
	makeHandler := makeHandlerClosure(db)
	http.HandleFunc("/people", makeHandler(listDbPeopleHandler))
	http.HandleFunc("/peopleStatic", listPeopleHandler)
	http.ListenAndServe(":8080", nil)
}

func main() {
	var db *sql.DB
	var err error

	db, err = initDatabse()
	if err != nil {
		log.Fatal(err)
	}

	var person *Person
	person, err = findPerson(26, db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(person)
	fmt.Println(person.Smt())

	var people []Person
	people, err = listPeople(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(people)

	startServer(db)
}