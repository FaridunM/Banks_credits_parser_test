package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"parser_of_credits/structs"
	"strings"

	"github.com/spf13/viper"
)

// if the condition is true, then return a, else return b
func TernarOp[T any](condition bool, a, b T) T {
	if condition {
		return a
	}

	return b
}

// removing last spaces via recursion
func DeleteLastSpace(s string) string {
	if len(s) == 0 {
		return ""
	}

	if s[len(s)-1] == ' ' {
		s = DeleteLastSpace(s[:len(s)-1])
	}

	return s
}

// Remove cutSymbol from left side of string
func LoopLeftTrim(str string, cutSymbol string) (newStr string) {
	newStr = strings.TrimLeft(str, cutSymbol)

	if string(newStr[0]) == cutSymbol {
		newStr = LoopLeftTrim(newStr, cutSymbol)
	}

	return newStr
}

// Remove cutSymbol from rigth side of string
func LoopRigthTrim(str string, cutSymbol string) (newStr string) {
	newStr = strings.TrimLeft(str, cutSymbol)

	if string(newStr[len(newStr)-1]) == cutSymbol {
		newStr = LoopRigthTrim(newStr, cutSymbol)
	}

	return newStr
}

func GetCreditUrls() (creditsUrl structs.CreditsUrl, err error) {
	viper.SetConfigFile("credits.json")
	if err = viper.ReadInConfig(); err != nil {
		log.Println("error in read in json-file:", err)
		return structs.CreditsUrl{}, err
	}

	if err = viper.Unmarshal(&creditsUrl); err != nil {
		log.Println("error in unmarshaling:", err)
		return structs.CreditsUrl{}, err
	}

	return
}

func ModifyCreditText(credit structs.Credit) structs.Credit {
	for key, value := range credit.Collateral {
		value = strings.TrimSpace(value)
		credit.Collateral[key] = value
	}

	for key, value := range credit.Amount {
		value = strings.TrimSpace(value)
		credit.Amount[key] = value
	}

	for key, value := range credit.LoanTerms {
		value = strings.TrimSpace(value)
		credit.LoanTerms[key] = value
	}

	for key, value := range credit.NeedDocuments {
		value = strings.TrimSpace(value)
		credit.NeedDocuments[key] = value
	}

	for key, value := range credit.InterestRate {
		value = strings.TrimSpace(value)
		credit.NeedDocuments[key] = value
	}

	return credit
}

func PanicCatcher() {
	if err := recover(); err != nil {
		log.Println("panic msg -", err)
	}
}

func WriteJson(w http.ResponseWriter, structure any) error {
	defer PanicCatcher()
	w.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(structure)
	if err != nil {
		return err
	}

	fmt.Fprint(w, string(body))

	return nil
}
