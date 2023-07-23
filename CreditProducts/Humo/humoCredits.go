package Humo

import (
	"fmt"
	"log"
	"parser_of_credits/structs"
	"parser_of_credits/utils"
)

var (
	consumerAndBusinessPrompts      = map[string]string{"Title": "#mainh1", "Description": "#mainBody p", "Collateral": "#mainBody p", "Amount": "#mainBody p", "LoanPeriod": "#mainBody p", "LoanTerms": "#mainBody ul", "NeedDocuments": "#mainBody ul"}
	mortgagePrompts                 = map[string]string{"Title": "#mainh1", "Description": "#mainBody p", "Collateral": "#mainBody p", "Amount": "#mainBody p", "InterestRate": "#mainBody p", "LoanPeriod": "#mainBody p", "InitialPayment": "#mainBody p", "LoanTerms": "#mainBody p", "NeedDocuments": "#mainBody ul"}
	educationPrompts                = map[string]string{"Title": "#mainh1", "Currency": "#mainBody p", "Amount": "#mainBody ul", "LoanPeriod": "#mainBody", "GracePeriod": "#mainBody", "LoanTerms": "#mainBody ul", "NeedDocuments": "#mainBody ul"}
	agriculturalAndLivestockPrompts = map[string]string{"Title": "#mainh1", "Description": "#mainBody", "LandingMethodology": "#mainBody", "Collateral": "#mainBody ul", "Amount": "#mainBody ul", "LoanPeriod": "#mainBody p", "LoanTerms": "#mainBody ul", "NeedDocuments": "#mainBody ul"}
)

// Getting url by credit type
func getCreditUrlBy(creditType string) (url string, err error) {
	creditUrls, err := utils.GetCreditUrls()
	if err != nil {
		log.Println("Error in getting credit urls")
	}

	switch creditType {
	case "consumer":
		url = creditUrls.Humo.Consumer
	case "startBusiness":
		url = creditUrls.Humo.StartBusiness
	case "mortgage":
		url = creditUrls.Humo.Mortgage
	case "education":
		url = creditUrls.Humo.Education
	case "agricultural":
		url = creditUrls.Humo.Agricultural
	case "livestock":
		url = creditUrls.Humo.Livestock
	case "orzu":
		url = creditUrls.Humo.Orzu
	default:
		url = ""
		err = fmt.Errorf("creditType %s not found", creditType)
	}

	return
}

// Getting all humo credits
func GetHumoCredits() (humoCredits structs.HumoCreditProducts, err error) {
	defer utils.PanicCatcher()

	consumerUrl, err := getCreditUrlBy("consumer")
	if err != nil {
		log.Println("err in getting consumer url:", err)
	}
	humoCredits.Consumer = utils.ModifyCreditText(getHumoCredit("consumer", consumerUrl, consumerAndBusinessPrompts))

	url, err := getCreditUrlBy("startBusiness")
	if err != nil {
		log.Println("err in getting start business url:", err)
	}
	humoCredits.StartBusiness = utils.ModifyCreditText(getHumoCredit("startBusiness", url, consumerAndBusinessPrompts))

	mortgageUrl, err := getCreditUrlBy("mortgage")
	if err != nil {
		log.Println("err in getting mortgage url:", err)
	}
	humoCredits.Mortgage = utils.ModifyCreditText(getHumoCredit("mortgage", mortgageUrl, mortgagePrompts))

	educationUrl, err := getCreditUrlBy("education")
	if err != nil {
		log.Println("err in getting education url:", err)
	}
	humoCredits.Education = utils.ModifyCreditText(getHumoCredit("education", educationUrl, educationPrompts))

	agriculturalUrl, err := getCreditUrlBy("agricultural")
	if err != nil {
		log.Println("err in getting agricultural url:", err)
	}
	humoCredits.Agricultural = utils.ModifyCreditText(getHumoCredit("agricultural", agriculturalUrl, agriculturalAndLivestockPrompts))

	livestockUrl, err := getCreditUrlBy("livestock")
	if err != nil {
		log.Println("err in getting livestock url:", err)
	}
	humoCredits.Livestock = utils.ModifyCreditText(getHumoCredit("livestock", livestockUrl, agriculturalAndLivestockPrompts))

	orzuUrl, err := getCreditUrlBy("orzu")
	if err != nil {
		log.Println("err in getting orzu url:", err)
	}
	humoCredits.Orzu = utils.ModifyCreditText(getOrzuCredit(orzuUrl))

	return
}

// Obtaining a specific credit
func GetHumoCreditBy(title string) (humoCredits structs.Credit, err error) {
	defer utils.PanicCatcher()

	url, err := getCreditUrlBy(title)
	if err != nil {
		log.Println(err)
	}

	switch title {
	case "consumer", "startBusiness":
		humoCredits = getHumoCredit(title, url, consumerAndBusinessPrompts)
	case "mortgage":
		humoCredits = getHumoCredit(title, url, mortgagePrompts)
	case "education":
		humoCredits = getHumoCredit(title, url, educationPrompts)
	case "agricultural", "livestock":
		humoCredits = getHumoCredit(title, url, agriculturalAndLivestockPrompts)
	case "orzu":
		humoCredits = getOrzuCredit(url)
	}

	return utils.ModifyCreditText(humoCredits), nil
}
