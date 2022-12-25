package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// Usuario :)
type Usuario struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

// UsuarioHandler analisa o request e delega para função adequada
func UsuarioHandler(w http.ResponseWriter, r *http.Request) {
	sid := strings.TrimPrefix(r.URL.Path, "/usuarios/")
	id, _ := strconv.Atoi(sid)

	switch {
	case r.Method == "GET" && id > 0:
		usuarioPorID(w, r, id)
	case r.Method == "GET":
		usuarioTodos(w, r)
	case r.Method == "POST":
		usuarioInsert(w, r)
	case r.Method == "PUT":
		usuarioUpdate(w, r)
	case r.Method == "DELETE":
		usuarioDeleteProc(w, r)
		//usuarioDelete(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Desculpa... :(")
	}
}

func usuarioPorID(w http.ResponseWriter, r *http.Request, id int) {
	db, err := sql.Open("mysql", "root:123456@/cursogo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var u Usuario
	db.QueryRow("select id, nome from usuarios where id = ?", id).Scan(&u.ID, &u.Nome)

	json, _ := json.Marshal(u)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))
}

func usuarioTodos(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:123456@/cursogo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, _ := db.Query("select id, nome from usuarios")
	defer rows.Close()

	var usuarios []Usuario
	for rows.Next() {
		var usuario Usuario
		rows.Scan(&usuario.ID, &usuario.Nome)
		usuarios = append(usuarios, usuario)
	}

	json, _ := json.Marshal(usuarios)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))
}

func usuarioInsert(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:123456@/cursogo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	resBody, err := ioutil.ReadAll(r.Body)
	if err != nil {

		w.WriteHeader(401)
		fmt.Printf("client: could not read response body: %s\n", err)

	} else {
		//imprimindo response body no foramto json:{"Id":xxx, "Nome":"yyyy"}
		fmt.Printf("client: response body: %s\n", resBody)
		var u Usuario
		//convertendo json para struct
		err_json := json.Unmarshal(resBody, &u)
		if err_json != nil {

			w.WriteHeader(401)
			fmt.Printf("client: could not read json response body: %s\n", err_json)

		} else {

			fmt.Printf("client: response body Nome: %s\n", u.Nome)
			stmt, _ := db.Prepare("insert into usuarios(nome) values(?)")
			res, _ := stmt.Exec(u.Nome)

			id, _ := res.LastInsertId()

			//imprimindo id
			fmt.Println(id)
			//passando id para 32 bits
			u.ID = int(id)
			//convertendo para json
			json, _ := json.Marshal(u)

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, string(json))

		}

	}

}

func usuarioUpdate(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:123456@/cursogo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	resBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		log.Fatal(err)
	}
	//imprimindo response body no foramto json:{"Id":xxx, "Nome":"yyyy"}
	fmt.Printf("client: response body: %s\n", resBody)
	var u Usuario
	//convertendo json para struct
	err_json := json.Unmarshal(resBody, &u)
	if err_json != nil {
		w.WriteHeader(401)
		fmt.Printf("client: could not read json response body: %s\n", err_json)

	} else {

		var uQ Usuario
		db.QueryRow("select id, nome from usuarios where id = ?", u.ID).Scan(&uQ.ID, &uQ.Nome)
		if uQ.ID != u.ID {
			w.WriteHeader(401)
			fmt.Printf("client error update: doesn't exist user: 401\n")
		} else {

			fmt.Printf("client: response body Nome: %s\n", u.Nome)
			stmt, _ := db.Prepare("update usuarios set nome = ? where id = ?")
			_, error_upd := stmt.Exec(u.Nome, u.ID)
			if error_upd != nil {
				w.WriteHeader(401)
				fmt.Printf("client: error update: %s\n", error_upd)

			} else {
				//convertendo para json
				json, _ := json.Marshal(u)

				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, string(json))
			}
		}
	}
}

func usuarioDeleteProc(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:123456@/cursogo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	resBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		log.Fatal(err)
	}
	//imprimindo response body no foramto json:{"Id":xxx, "Nome":"yyyy"}
	fmt.Printf("client: response body: %s\n", resBody)
	var u Usuario
	//convertendo json para struct
	err_json := json.Unmarshal(resBody, &u)
	if err_json != nil {

		w.WriteHeader(401)
		fmt.Printf("client: could not read json response body: %s\n", err_json)

	} else {

		stmt, _ := db.Prepare("CALL DeleteUser(?)")
		res, _ := stmt.Exec(u.ID)
		linhas, _ := res.RowsAffected()

		if linhas > 0 {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(401)
		}
	}
}
