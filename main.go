package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"parser_of_credits/CreditProducts/Eskhata"
	"parser_of_credits/CreditProducts/Humo"
	"parser_of_credits/utils"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.Use(mux.CORSMethodMiddleware(router))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello. For use this service using route '/getGreditBy' with parameters bank and creditType"))
	}).Methods("GET")
	router.HandleFunc("/getGreditBy", getGreditBy).Methods("GET")

	// starting server
	log.Println("Starting server... on ", "0.0.0.0:7777")
	err := http.ListenAndServe("0.0.0.0:7777", router)
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

func getGreditBy(w http.ResponseWriter, r *http.Request) {
	bank := r.URL.Query().Get("bank")
	if bank == "" {
		w.Write([]byte("bank not found"))
		return
	}

	creditType := r.URL.Query().Get("creditType")

	switch bank {
	case "humo":
		if creditType == "" {
			humoCredits, err := Humo.GetHumoCredits()
			if err != nil {
				log.Println("err -", err.Error())
			}
			WriteJson(w, humoCredits)
			break
		}

		humoCredit, err := Humo.GetHumoCreditBy(creditType)
		if err != nil {
			log.Println("err -", err.Error())
		}
		WriteJson(w, humoCredit)
	case "eskhata":
		if creditType == "" {
			eskhataCredits, err := Eskhata.GetEskhataCredits()
			if err != nil {
				log.Println("err -", err.Error())
			}
			WriteJson(w, eskhataCredits)
			break
		}

		eskhataCredit, err := Eskhata.GetEskhataCreditBy(creditType)
		if err != nil {
			log.Println("err -", err.Error())
		}
		WriteJson(w, eskhataCredit)
	default:
		w.Write([]byte("bank not found"))
	}
}

func WriteJson(w http.ResponseWriter, structure any) error {
	defer utils.PanicCatcher()
	w.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(structure)
	if err != nil {
		return err
	}

	fmt.Fprint(w, string(body))

	return nil
}
