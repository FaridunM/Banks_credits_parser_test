package Humo

import (
	"parser_of_credits/structs"
	"parser_of_credits/utils"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

// Colly constructor for prepare to parsing
func constructor() *colly.Collector {
	return colly.NewCollector(
		colly.AllowedDomains("www.humo.tj", "humo.tj"),
	)
}

// Getting (from site - parsing) credit by credit type, url and prompts. The function returns the struct of Credit type
func getHumoCredit(creditType, url string, prompts map[string]string) (credit structs.Credit) {
	coli := constructor()

	coli.OnHTML(prompts["Title"], func(h *colly.HTMLElement) {
		credit.Title = strings.TrimSpace(h.Text)
	})

	if strings.Contains("agricultural, livestock", creditType) {
		coli.OnHTML(prompts["Description"], func(h *colly.HTMLElement) {
			if ok, err := regexp.MatchString("(на выращивание сельскохозяйственных культур|разведение домашнего скота).*", h.Text); err == nil && ok {
				pattern := regexp.MustCompile("(Позволяет|Кредит на разведение домашнего скота) [а-яА-Я, ]*")
				description := pattern.FindString(strings.Replace(h.Text, ".", "", -1))
				credit.Description = strings.TrimSpace(description)
			}
		})
	} else {
		counter := 0
		coli.OnHTML(prompts["Description"], func(h *colly.HTMLElement) {
			counter++
			credit.Description = utils.TernarOp(counter == 1, strings.TrimSpace(h.Text), credit.Description)
		})
	}

	coli.OnHTML(prompts["Currency"], func(h *colly.HTMLElement) {
		if ok, err := regexp.MatchString("(Валюта кредита):([а-яА-Я ]*)", h.Text); err == nil && ok {
			pattern := regexp.MustCompile("(Валюта кредита):([а-яА-Я ]*)")
			credit.Currency = pattern.FindAllStringSubmatch(h.Text, -1)[0][2]
		}
	})

	coli.OnHTML(prompts["LandingMethodology"], func(h *colly.HTMLElement) {
		if ok, err := regexp.MatchString("(Методология кредитования):([а-я ]*)", h.Text); err == nil && ok {
			pattern := regexp.MustCompile("(Методология кредитования):([а-я ]*)")
			credit.LandingMethodology = strings.TrimSpace(pattern.FindAllStringSubmatch(h.Text, -1)[0][2])
		}
	})

	if strings.Contains("agricultural, livestock", creditType) {
		coli.OnHTML(prompts["Collateral"], func(h *colly.HTMLElement) {
			if ok, err := regexp.MatchString("[а-я ]:[а-я ]", h.DOM.Children().Text()); err == nil && ok {
				credit.Collateral = strings.Split(strings.Replace(h.DOM.Children().Text(), ".", "", -1), ";")
			}
		})
	} else {
		coli.OnHTML(prompts["Collateral"], func(h *colly.HTMLElement) {
			if ok, err := regexp.MatchString("(Обеспечение|Гарантия|Залоговое обеспечение): ([а-я ]*)", h.Text); err == nil && ok {
				pattern := regexp.MustCompile("(Обеспечение|Гарантия|Залоговое обеспечение): ([а-я ]*)")
				credit.Collateral = append(credit.Collateral, pattern.FindAllStringSubmatch(h.Text, -1)[0][2])
			}
		})
	}

	if strings.Contains("education, agricultural, livestock", creditType) {
		coli.OnHTML(prompts["Amount"], func(h *colly.HTMLElement) {
			if ok, err := regexp.MatchString("[0-9]{3,}", strings.Replace(h.Text, " ", "", -1)); err == nil && ok {
				replaceElement := utils.TernarOp(creditType == "education", ".", ". ")
				credit.Amount = strings.Split(strings.Replace(h.DOM.Children().Text(), replaceElement, "", -1), ";")
			}
		})
	} else {
		coli.OnHTML(prompts["Amount"], func(h *colly.HTMLElement) {
			if ok, err := regexp.MatchString("[0-9]{3,}", strings.Replace(h.Text, " ", "", -1)); err == nil && ok {
				pattern := utils.TernarOp(creditType == "startBusiness", "(Сумма):(.*)", "(Сумма кредита|Сумма):([а-яА-Я0-9, ]*)")
				regex := regexp.MustCompile(pattern)
				tempStr := utils.LoopLeftTrim(regex.FindAllStringSubmatch(strings.ReplaceAll(h.Text, ".", ""), -1)[0][2], " ")

				credit.Amount = append(credit.Amount, tempStr)
			}
		})
	}

	coli.OnHTML(prompts["LoanPeriod"], func(h *colly.HTMLElement) {
		if ok, err := regexp.MatchString("Продолжительность займа|Продолжительность|Срок кредита.*", h.Text); err == nil && ok {
			pattern := regexp.MustCompile("(Продолжительность займа|Продолжительность|Срок кредита): ([а-я0-9, ]*)")
			credit.LoanPeriod = pattern.FindAllStringSubmatch(strings.ReplaceAll(h.Text, ".", ""), -1)[0][2]
		}
	})

	coli.OnHTML(prompts["GracePeriod"], func(h *colly.HTMLElement) {
		if ok, err := regexp.MatchString("Льготный период.*", h.Text); err == nil && ok {
			pattern := regexp.MustCompile("Льготный период:(.*)")
			credit.GracePeriod = strings.TrimSpace(pattern.FindAllStringSubmatch(strings.ReplaceAll(h.Text, ".", ""), -1)[0][1])
		}
	})

	coli.OnHTML(prompts["InterestRate"], func(h *colly.HTMLElement) {
		if ok, err := regexp.MatchString("Процентная ставка: ", h.Text); err == nil && ok {
			pattern := regexp.MustCompile("(Процентная ставка):([а-яА-Я0-9%, ]*)")
			credit.InterestRate = append(credit.InterestRate, strings.TrimSpace(pattern.FindStringSubmatch(h.Text)[2]))
		}
	})

	coli.OnHTML(prompts["InitialPayment"], func(h *colly.HTMLElement) {
		if ok, err := regexp.MatchString("Первоначальный взнос: ", h.Text); err == nil && ok {
			pattern := regexp.MustCompile("(Первоначальный взнос):([а-яА-Я0-9%, ]*)")
			credit.InitialPayment = strings.TrimSpace(pattern.FindStringSubmatch(h.Text)[2])
		}
	})

	switch creditType {
	case "mortgage":
		coli.OnHTML(prompts["LoanTerms"], func(h *colly.HTMLElement) {
			if ok, err := regexp.MatchString("Кредит выдается на покупку жилья.*", h.Text); err == nil && ok {
				terms := strings.ReplaceAll(h.Text, "\n\t", "")
				terms = utils.LoopRigthTrim(terms, ".")
				credit.LoanTerms = strings.Split(terms, ".")
				for k, v := range credit.LoanTerms {
					credit.LoanTerms[k] = strings.TrimLeft(v, " ")
				}
			}
		})
	case "agricultural", "livestock":
		allTerms := []string{}
		coli.OnHTML(prompts["LoanTerms"], func(h *colly.HTMLElement) {
			if ok, err := regexp.MatchString("(Физические лица старше 20 лет|Количество заемщиков группы от 2 до 10 человек)", h.Text); err == nil && ok {
				h.Text = strings.ReplaceAll(h.Text, ".", "")
				allTerms = strings.Split(strings.ReplaceAll(h.Text, "\n\t", ""), ";")
				credit.LoanTerms = append(credit.LoanTerms, allTerms...)
			}
		})
	default:
		coli.OnHTML(prompts["LoanTerms"], func(h *colly.HTMLElement) {
			if ok, err := regexp.MatchString("Физические лица старше 20 лет", h.Text); err == nil && ok {
				h.Text = strings.ReplaceAll(h.Text, ".", "")
				credit.LoanTerms = strings.Split(strings.ReplaceAll(h.Text, "\n\t", ""), ";")
			}
		})
	}

	coli.OnHTML(prompts["NeedDocuments"], func(h *colly.HTMLElement) {
		if ok, err := regexp.MatchString("Кредитная заявка", h.Text); err == nil && ok {
			h.Text = strings.ReplaceAll(h.Text, ".", "")
			credit.NeedDocuments = strings.Split(strings.ReplaceAll(h.Text, "\n\t", ""), ";")
		}
	})
	credit.Url = url

	coli.Visit(url)
	return
}

// Getting (from site - parsing) credit Orzu which some different from others credit by parsing. Function gets in argument url and return struct of Credit type
func getOrzuCredit(url string) (credit structs.Credit) {
	coli := func() *colly.Collector {
		return colly.NewCollector(
			colly.AllowedDomains("orzu.humo.tj"),
		)
	}()

	credit.Title = "Орзу"
	credit.Description = "Кредит наличными до 50 000 сомони. На любые цели с картой «Орзу»"
	credit.InterestRate = []string{"15д - 2%", "1м - 4%", "2м - 6%", "3м - 8%", "6м - 13%", "9м - 18%", "12м - 23%", "18м - 34%"}
	credit.Url = url

	coli.OnHTML(".Information_title__tlejv", func(h *colly.HTMLElement) {
		if ok, err := regexp.MatchString("[0-9]{2,3} сомони", h.Text); err == nil && ok {
			credit.Amount = append(credit.Amount, h.Text)
		}

		if ok, err := regexp.MatchString("до [0-9] лет", h.Text); err == nil && ok {
			credit.LoanPeriod = strings.TrimSpace(h.Text)
		}
	})

	coli.Visit(url)

	return credit
}
