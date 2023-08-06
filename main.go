package main

import (
	"log"
	"net/http"

	handler "parser_of_credits/api"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter().PathPrefix("/api").Subrouter()
	router.Use(mux.CORSMethodMiddleware(router))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello. For use this service using route '/getGreditBy' with parameters bank and creditType"))
	}).Methods("GET")
	router.HandleFunc("/getGreditBy", handler.GetGreditBy).Methods("GET")

	// starting server
	log.Println("Starting server... on ", "0.0.0.0:8080")
	err := http.ListenAndServe("0.0.0.0:8080", router)
	if err != nil {
		log.Fatal(err)
	}

	// -------------------------------------------------------------------

	// Getting all credits from Humo
	// humoCredits, err := Humo.GetHumoCredits()
	// if err != nil {
	// 	log.Println("err -", err)
	// }
	// fmt.Println("humo credits -", humoCredits)

	// Getting credit by title from Humo
	// res, err := Humo.GetHumoCreditBy("orzu")
	// if err != nil {
	// 	log.Println("err -", err)
	// }
	// fmt.Println("credit for orzu -", res)

	// Getting Eskhata credits
	// eskhataCredits, err := Eskhata.GetEskhataCredits()
	// if err != nil {
	// 	log.Println("err -", err)
	// }
	// fmt.Println("eskhata credits -", eskhataCredits)

	// expressEskhataCredits, err := Eskhata.GetEskhataCreditBy("express")
	// if err != nil {
	// 	log.Println("err -", err)
	// }
	// fmt.Println("express eskhata credit -", expressEskhataCredits)
}
