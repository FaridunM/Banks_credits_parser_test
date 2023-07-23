package Eskhata

import (
	"fmt"
	"log"
	"parser_of_credits/structs"
	"parser_of_credits/utils"
)

var (
	multiPurposeConsumerPrompts          = map[string]string{"Title": "h2", "Description": "div", "Purpose-Amount-LoanPeriod": "div", "InterestRate-NeedDocuments": "div ul", "Collateral": "div"}
	promptsForExpressHunarkhoiMardumiCar = map[string]string{"Title": "h2", "Description": "div", "Purpose": "div", "Amount": "div ul", "Currency": "div ul", "LoanPeriod": "div ul", "InterestRate": "div ul", "Collateral": "div ul", "NeedDocuments": "div ul"}
	imkonPrompts                         = map[string]string{"Title": "h2", "Description": "div p font", "Amount": "div ul", "Currency": "div ul", "InterestRate": "div ul", "LoanPeriod": "div ul", "GracePeriod": "div ul", "NeedDocuments": "div ul"}
	manziliPrompts                       = map[string]string{"Title": "h2", "Description": "div p font", "Purpose": "div", "AdditionalFiled": "div", "NeedDocuments": "div ul"}
)

// Getting url and anchor by credit type
func getCreditUrlBy(creditType string) (url, anchor string, err error) {
	creditUrls, err := utils.GetCreditUrls()
	if err != nil {
		log.Println("Error in getting credit urls")
	}
	url = creditUrls.Eskhata.CommonUrl

	switch creditType {
	case "multiPurposeConsumer":
		anchor = creditUrls.Eskhata.MultiPurposeConsumer
	case "express":
		anchor = creditUrls.Eskhata.Express
	case "hunarkhoiMardumi":
		anchor = creditUrls.Eskhata.HunarkhoiMardumi
	case "imkon":
		anchor = creditUrls.Eskhata.Imkon
	case "car":
		anchor = creditUrls.Eskhata.Car
	case "manzili":
		anchor = creditUrls.Eskhata.Manzili
	default:
		err = fmt.Errorf("creditType %s not found", creditType)
	}

	return
}

// Getting all eskhata credits
func GetEskhataCredits() (eskhataCredits structs.EskhataCreditProducts, err error) {
	defer utils.PanicCatcher()

	multiPurposeConsumerUrl, anchor, err := getCreditUrlBy("multiPurposeConsumer")
	if err != nil {
		log.Println("error:", err.Error())
	}
	eskhataCredits.MultiPurposeConsumer = getEskhataCredit(multiPurposeConsumerUrl, anchor, multiPurposeConsumerPrompts)

	expressUrl, anchor, err := getCreditUrlBy("express")
	if err != nil {
		log.Println("error:", err.Error())
	}
	eskhataCredits.Express = getEskhataCredit2(expressUrl, anchor, promptsForExpressHunarkhoiMardumiCar)

	hunarkhoiMardumiUrl, anchor, err := getCreditUrlBy("hunarkhoiMardumi")
	if err != nil {
		log.Println("error:", err.Error())
	}
	eskhataCredits.HunarkhoiMardumi = getEskhataCredit2(hunarkhoiMardumiUrl, anchor, promptsForExpressHunarkhoiMardumiCar)

	carUrl, anchor, err := getCreditUrlBy("car")
	if err != nil {
		log.Println("error:", err.Error())
	}
	eskhataCredits.Car = getEskhataCredit2(carUrl, anchor, promptsForExpressHunarkhoiMardumiCar)

	imkonUrl, anchor, err := getCreditUrlBy("imkon")
	if err != nil {
		log.Println("error:", err.Error())
	}
	eskhataCredits.Imkon = getEskhataCredit3(imkonUrl, anchor, imkonPrompts)

	manziliUrl, anchor, err := getCreditUrlBy("manzili")
	if err != nil {
		log.Println("error:", err.Error())
	}
	eskhataCredits.Manzili = getEskhataCredit4(manziliUrl, anchor, manziliPrompts)

	return
}

// Obtaining a specific credit
func GetEskhataCreditBy(title string) (eskhataCredits structs.Credit, err error) {
	defer utils.PanicCatcher()

	url, anchor, err := getCreditUrlBy(title)
	if err != nil {
		log.Println(err)
	}

	switch title {
	case "multiPurposeConsumer":
		eskhataCredits = getEskhataCredit(url, anchor, multiPurposeConsumerPrompts)
	case "express":
		eskhataCredits = getEskhataCredit2(url, anchor, promptsForExpressHunarkhoiMardumiCar)
	case "hunarkhoiMardumi":
		eskhataCredits = getEskhataCredit2(url, anchor, promptsForExpressHunarkhoiMardumiCar)
	case "car":
		eskhataCredits = getEskhataCredit2(url, anchor, promptsForExpressHunarkhoiMardumiCar)
	case "imkon":
		eskhataCredits = getEskhataCredit3(url, anchor, imkonPrompts)
	case "manzili":
		eskhataCredits = getEskhataCredit4(url, anchor, manziliPrompts)
	default:
		url = ""
		err = fmt.Errorf("credit %s not found", title)
	}

	return utils.ModifyCreditText(eskhataCredits), nil
}
