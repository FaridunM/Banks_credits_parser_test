package handler

import (
	"log"
	"net/http"
	"parser_of_credits/CreditProducts/Eskhata"
	"parser_of_credits/CreditProducts/Humo"
	"parser_of_credits/utils"
)

func GetGreditBy(w http.ResponseWriter, r *http.Request) {
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
			utils.WriteJson(w, humoCredits)
			break
		}

		humoCredit, err := Humo.GetHumoCreditBy(creditType)
		if err != nil {
			log.Println("err -", err.Error())
		}
		utils.WriteJson(w, humoCredit)
	case "eskhata":
		if creditType == "" {
			eskhataCredits, err := Eskhata.GetEskhataCredits()
			if err != nil {
				log.Println("err -", err.Error())
			}
			utils.WriteJson(w, eskhataCredits)
			break
		}

		eskhataCredit, err := Eskhata.GetEskhataCreditBy(creditType)
		if err != nil {
			log.Println("err -", err.Error())
		}
		utils.WriteJson(w, eskhataCredit)
	default:
		w.Write([]byte("bank not found"))
	}
}
