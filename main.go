package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type matched struct {
	Id             int
	Gender         string
	Height         int
	Age            int
	Ask_gender     string
	Ask_height_up  int
	Ask_height_low int
	Ask_age_up     int
	Ask_age_low    int
}

type people struct {
	Person []matched
}

var pe = people{}
var tmpl = template.Must(template.ParseFiles("index.html"))

func searchhandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		pe.Person = nil
		r.ParseForm()

		gender := r.Form.Get("gender")

		heightStr := r.Form.Get("height")
		height, err := strconv.Atoi(heightStr)
		if err != nil {
			http.Error(w, "Invalid input, please enter a valid number", http.StatusBadRequest)
			return
		}

		ageStr := r.Form.Get("age")
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			http.Error(w, "Invalid input, please enter a valid number", http.StatusBadRequest)
			return
		}

		ask_gender := r.Form.Get("ask_gender")

		ask_height_upStr := r.Form.Get("ask_height_up")
		ask_height_up, err := strconv.Atoi(ask_height_upStr)
		if err != nil {
			http.Error(w, "Invalid input, please enter a valid number", http.StatusBadRequest)
			return
		}

		ask_height_lowStr := r.Form.Get("ask_height_low")
		ask_height_low, err := strconv.Atoi(ask_height_lowStr)
		if err != nil {
			http.Error(w, "Invalid input, please enter a valid number", http.StatusBadRequest)
			return
		}

		ask_age_upStr := r.Form.Get("ask_age_up")
		ask_age_up, err := strconv.Atoi(ask_age_upStr)
		if err != nil {
			http.Error(w, "Invalid input, please enter a valid number", http.StatusBadRequest)
			return
		}

		ask_age_lowStr := r.Form.Get("ask_age_low")
		ask_age_low, err := strconv.Atoi(ask_age_lowStr)
		if err != nil {
			http.Error(w, "Invalid input, please enter a valid number", http.StatusBadRequest)
			return
		}

		db, err := sql.Open("mysql", "user:1234@tcp(localhost:3306)/kkbox?charset=utf8")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		fmt.Println("Executing query with parameters:", ask_gender, gender, ask_height_up, ask_height_low, height, height, ask_age_up, ask_age_low, age, age)

		rows, err := db.Query("select * from matched_table where "+
			"ask_gender = ? and gender = ? "+
			"and ask_height_up >= ? and ask_height_low <= ? "+
			"and height <= ? and height >= ? "+
			"and ask_age_up >= ? and ask_age_low <= ? "+
			"and age <= ? and age >= ? ",
			ask_gender, gender, ask_height_up, ask_height_low, height, height, ask_age_up, ask_age_low,
			age, age)
		if err != nil {
			fmt.Println("Error executing SQL query:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			ma := matched{}
			err := rows.Scan(&ma.Id, &ma.Gender, &ma.Height, &ma.Age, &ma.Ask_gender, &ma.Ask_height_up, &ma.Ask_height_low, &ma.Ask_age_up, &ma.Ask_age_low)
			if err != nil {
				fmt.Println("Error scanning database row:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			pe.Person = append(pe.Person, ma)
			fmt.Println("Matched Person:", ma)
		}

		maxID, err := getMaxID(db)
		if err != nil {
			fmt.Println("Error getting max ID:", err)
			return
		}

		fmt.Println("Max ID:", maxID)
		insert, err := db.Exec("INSERT INTO matched_table VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
			maxID+1, gender, height, age, ask_gender, ask_height_up, ask_height_low, ask_age_up, ask_age_low)
		if err != nil {
			http.Error(w, "Error insert data", http.StatusInternalServerError)
			return
		}
		fmt.Print(insert)

		delete, err := db.Exec("DELETE FROM matched_table WHERE id = ?", maxID+1)
		if err != nil {
			http.Error(w, "Error delete data", http.StatusInternalServerError)
			return
		}
		fmt.Print(delete)
	}

	err := template.Must(template.ParseFiles("result.html")).Execute(w, pe)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
	}
}

func getMaxID(db *sql.DB) (int, error) {
	var maxID int
	err := db.QueryRow("SELECT MAX(id) FROM matched_table").Scan(&maxID)
	if err != nil {
		return 0, err
	}
	return maxID, nil
}

func handle(w http.ResponseWriter, r *http.Request) {
	err := tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println("Error rendering template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()

		IDStr := r.Form.Get("ID")
		ID, err := strconv.Atoi(IDStr)
		if err != nil {
			http.Error(w, "Invalid input for ID", http.StatusBadRequest)
			return
		}

		gender := r.Form.Get("subject")

		heightStr := r.Form.Get("height")
		height, err := strconv.Atoi(heightStr)
		if err != nil {
			http.Error(w, "Invalid input for height", http.StatusBadRequest)
			return
		}

		ageStr := r.Form.Get("age")
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			http.Error(w, "Invalid input for age", http.StatusBadRequest)
			return
		}

		db, err := sql.Open("mysql", "user:1234@tcp(localhost:3306)/kkbox?charset=utf8")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		insert, err := db.Exec("INSERT INTO matched_table VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
			ID, gender, height, age, "", 0, 0, 0, 0)
		if err != nil {
			http.Error(w, "Error insert data", http.StatusInternalServerError)
			return
		}
		fmt.Print(insert)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func main() {
	http.HandleFunc("/", handle)
	http.HandleFunc("/search", searchhandle)
	http.HandleFunc("/add", addHandler)
	http.ListenAndServe(":8080", nil)
}
