package Eskhata

import (
	"fmt"
	"parser_of_credits/structs"
	"parser_of_credits/utils"
	"regexp"
	"strings"
	"unicode"

	"github.com/gocolly/colly"
)

// Colly constructor for prepare to parsing
func constructor() *colly.Collector {
	return colly.NewCollector(
		colly.AllowedDomains("www.eskhata.com", "eskhata.com"),
	)
}

// Getting (from site - parsing) credit by url, anchor and prompts v1. The function returns the struct of Credit type
func getEskhataCredit(url, anchor string, prompts map[string]string) (credit structs.Credit) {
	coli := constructor()
	credit.Url = url + anchor
	for key, value := range prompts {
		if !strings.Contains(value, anchor) {
			prompts[key] = fmt.Sprintf("%s %s", anchor, value)
		}
	}
	// getting title
	coli.OnHTML(prompts["Title"], func(h *colly.HTMLElement) {
		credit.Title = h.Text
	})

	// getting description
	coli.OnHTML(prompts["Description"], func(h *colly.HTMLElement) {
		var description string
		h.ForEachWithBreak("p", func(i int, e *colly.HTMLElement) bool {
			if ok, err := regexp.MatchString("Цель кредита:", e.Text); err == nil && ok {
				return false
			}

			description = description + strings.TrimSpace(strings.ReplaceAll(e.Text, "\n", ""))
			return true
		})

		credit.Description = utils.TernarOp(description != "" && description != credit.Description, description, credit.Description)
	})

	// getting purpose
	coli.OnHTML(prompts["Purpose-Amount-LoanPeriod"], func(h *colly.HTMLElement) {
		h.ForEachWithBreak("p", func(i int, e *colly.HTMLElement) bool {
			if ok, err := regexp.MatchString("Наши условия:", e.Text); err == nil && ok {
				return false
			} else if ok, err := regexp.MatchString("Цель кредита:", e.Text); err == nil && ok {
				credit.Purpose = strings.TrimSpace(e.DOM.Next().Text())
			}

			return true
		})

		var amount, period string
		h.ForEachWithBreak("p", func(i int, e *colly.HTMLElement) bool {
			if ok, err := regexp.MatchString("Процентная ставка:", e.Text); err == nil && ok {
				return false
			} else if ok, err := regexp.MatchString("Наши условия:", e.Text); err == nil && ok {
				pattern := regexp.MustCompile(`до \d.*`)
				creditAmount := e.DOM.Next()
				creditPeriod := creditAmount.Next()
				amount = pattern.FindString(creditAmount.Text())
				period = pattern.FindString(creditPeriod.Text())

				if h.Index == 0 {
					credit.Amount = append(credit.Amount, amount[:len(amount)-2])
				}
				credit.LoanPeriod = period
			}

			return true
		})
	})

	// getting interest rate & need documents
	coli.OnHTML(prompts["InterestRate-NeedDocuments"], func(h *colly.HTMLElement) {
		if ok, err := regexp.MatchString("%", h.Text); err == nil && ok {
			h.Text = strings.ReplaceAll(h.Text, "\n", "")
			content := strings.Split(h.Text, "      ")

			for _, val := range content {
				credit.InterestRate = append(credit.InterestRate, strings.TrimSpace(val))
			}
		}
	})

	// getting collateral
	coli.OnHTML(prompts["Collateral"], func(h *colly.HTMLElement) {
		h.ForEachWithBreak("p", func(i int, e *colly.HTMLElement) bool {
			if ok, err := regexp.MatchString("Необходимые документы:", e.Text); err == nil && ok {
				return false
			} else if ok, err := regexp.MatchString("Залоговое обеспечение:", e.Text); err == nil && ok {
				creditCollateral1 := e.DOM.Next()
				creditCollateral2 := creditCollateral1.Next()

				if h.Index == 0 {
					credit.Collateral = append(credit.Collateral, creditCollateral1.Text(), creditCollateral2.Text())
				}
			}

			return true
		})
	})

	// getting need documents
	coli.OnHTML(prompts["InterestRate-NeedDocuments"], func(h *colly.HTMLElement) {
		if h.Index == 1 {
			h.Text = strings.ReplaceAll(h.Text, "\n", "")
			h.Text = strings.Trim(h.Text, " ")

			for _, val := range strings.Split(h.Text, "    ") {
				credit.NeedDocuments = append(credit.NeedDocuments, strings.TrimLeft(val, " "))
			}
		}
	})

	coli.Visit(url)

	return
}

// Getting (from site - parsing) credit by url, anchor and prompts v2. The function returns the struct of Credit type
func getEskhataCredit2(url, anchor string, prompts map[string]string) (credit structs.Credit) {
	coli := constructor()
	credit.Url = url + anchor
	for key, value := range prompts {
		prompts[key] = fmt.Sprintf("%s %s", anchor, value)
	}

	// getting title
	coli.OnHTML(prompts["Title"], func(h *colly.HTMLElement) {
		credit.Title = h.Text
	})

	// getting description
	coli.OnHTML(prompts["Description"], func(h *colly.HTMLElement) {
		var description string
		h.ForEachWithBreak("p", func(i int, e *colly.HTMLElement) bool {
			if ok, err := regexp.MatchString("Цель кредита:", e.Text); err == nil && ok {
				return false
			}

			e.Text = strings.ReplaceAll(e.Text, ".", "")
			description = description + strings.TrimSpace(e.Text)
			return true
		})

		credit.Description = utils.TernarOp(description != "" && description != credit.Description, description, credit.Description)
	})

	// getting purpose
	coli.OnHTML(prompts["Purpose"], func(h *colly.HTMLElement) {
		h.ForEachWithBreak("p", func(i int, e *colly.HTMLElement) bool {
			if ok, err := regexp.MatchString("Наши условия:", e.Text); err == nil && ok {
				return false
			} else if ok, err := regexp.MatchString("Цель кредита:", e.Text); err == nil && ok {
				domElement := strings.ReplaceAll(e.DOM.Next().Text(), ".", "")
				credit.Purpose = strings.TrimSpace(domElement)
			}

			return true
		})
	})

	// getting amount
	coli.OnHTML(prompts["Amount"], func(h *colly.HTMLElement) {
		if ok, err := regexp.MatchString("[С|C]умма кредита: [а-я0-9 ].*", h.Text); err == nil && ok {
			pattern := regexp.MustCompile("([С|C]умма кредита): ([а-я0-9 ].*)")
			amount := pattern.FindAllStringSubmatch(h.Text, -1)[0][2]
			credit.Amount = append(credit.Amount, removeAllBesideCDPrSPlS(amount))
		}
	})

	// getting currency
	coli.OnHTML(prompts["Currency"], func(h *colly.HTMLElement) {
		if ok, err := regexp.MatchString("Валюта кредита: [а-я].*", h.Text); err == nil && ok {
			pattern := regexp.MustCompile("(Валюта кредита): ([а-я].*)")
			credit.Currency = pattern.FindAllStringSubmatch(h.Text, -1)[0][2]
		}
	})

	// getting loan period
	coli.OnHTML(prompts["LoanPeriod"], func(h *colly.HTMLElement) {
		if ok, err := regexp.MatchString("Срок кредита: [а-я0-9 ].*", h.Text); err == nil && ok {
			pattern := regexp.MustCompile("(Срок кредита): ([а-я0-9 ].*)")
			credit.LoanPeriod = pattern.FindAllStringSubmatch(h.Text, -1)[0][2]
		}
	})

	// getting interest rate
	coli.OnHTML(prompts["InterestRate"], func(h *colly.HTMLElement) {
		h.ForEachWithBreak("li", func(i int, e *colly.HTMLElement) bool {
			if ok, err := regexp.MatchString("Процентн(ые|ая) ставк(и|а): [а-я0-9%].*|\\+[а-я0-9%].*", e.Text); err == nil && ok {
				e.Text = e.Text[strings.Index(e.Text, ":")+1:]
				e.Text = strings.Trim(e.Text, " ")
				if unicode.IsSpace(rune(e.Text[len(e.Text)-1])) {
					credit.InterestRate = append(credit.InterestRate, removeAllBesideCDPrSPlS(e.Text))
				}
			}

			return true
		})
	})

	// getting collateral
	coli.OnHTML(prompts["Collateral"], func(h *colly.HTMLElement) {
		h.ForEachWithBreak("li", func(i int, e *colly.HTMLElement) bool {
			if ok, err := regexp.MatchString("Залоговое обеспечение: [а-я0-9%].+", e.Text); err == nil && ok {
				e.Text = e.Text[strings.Index(e.Text, ":")+1:]
				credit.Collateral = append(credit.Collateral, strings.Trim(e.Text, " "))
			}

			return true
		})
	})

	// getting need documents
	coli.OnHTML(prompts["NeedDocuments"], func(h *colly.HTMLElement) {
		if h.Index == 1 {
			h.Text = strings.ReplaceAll(h.Text, "\n", "")
			h.Text = strings.Trim(h.Text, " ")
			credit.NeedDocuments = append(credit.NeedDocuments, strings.Split(h.Text, "   ")...)
		}
	})

	coli.Visit(url)

	return
}

// Getting (from site - parsing) credit by url, anchor and prompts v3. The function returns the struct of Credit type
func getEskhataCredit3(url, anchor string, prompts map[string]string) (credit structs.Credit) {
	coli := constructor()
	credit.Url = url + anchor
	for key, value := range prompts {
		prompts[key] = fmt.Sprintf("%s %s", anchor, value)
	}

	// getting title
	coli.OnHTML(prompts["Title"], func(h *colly.HTMLElement) {
		credit.Title = h.Text
	})

	// getting description
	coli.OnHTML(prompts["Description"], func(h *colly.HTMLElement) {
		credit.Description += utils.TernarOp(h.Index < 2, h.Text, "")
	})

	// getting amount
	coli.OnHTML(prompts["Amount"], func(h *colly.HTMLElement) {
		if ok, err := regexp.MatchString("[С|C]умма кредита: [а-я0-9 ].*", h.Text); err == nil && ok {
			pattern := regexp.MustCompile("([С|C]умма кредита): ([а-я0-9 ].*)")
			amount := pattern.FindAllStringSubmatch(h.Text, -1)[0][2]
			credit.Amount = append(credit.Amount, removeAllBesideCDPrSPlS(amount))
		}
	})

	// getting currency
	coli.OnHTML(prompts["Currency"], func(h *colly.HTMLElement) {
		if ok, err := regexp.MatchString("Валюта кредита: [а-я].*", h.Text); err == nil && ok {
			pattern := regexp.MustCompile("(Валюта кредита): ([а-я].*)")
			credit.Currency = pattern.FindAllStringSubmatch(h.Text, -1)[0][2]
		}
	})

	// getting interest rate
	coli.OnHTML(prompts["InterestRate"], func(h *colly.HTMLElement) {
		h.ForEachWithBreak("li", func(i int, e *colly.HTMLElement) bool {
			if ok, err := regexp.MatchString("(Годовая процентная|Процентн(ые|ая)) ставк(и|а): [а-я0-9%].*|\\+[а-я0-9%].*", e.Text); err == nil && ok {
				pattern := regexp.MustCompile("(Годовая процентная ставка): ([а-я0-9%].*)")
				credit.InterestRate = append(credit.InterestRate, pattern.FindAllStringSubmatch(h.Text, -1)[0][2])
			}

			return true
		})
	})

	// getting loan period
	coli.OnHTML(prompts["LoanPeriod"], func(h *colly.HTMLElement) {
		if ok, err := regexp.MatchString("Срок кредита: [а-я0-9 ].*", h.Text); err == nil && ok {
			pattern := regexp.MustCompile("(Срок кредита): ([а-я0-9 ].*)")
			credit.LoanPeriod = pattern.FindAllStringSubmatch(h.Text, -1)[0][2]
		}
	})

	// getting grace period
	coli.OnHTML(prompts["GracePeriod"], func(h *colly.HTMLElement) {
		if ok, err := regexp.MatchString("Есть.+льготн(ый|ого) период.*", h.Text); err == nil && ok {
			pattern := regexp.MustCompile("Есть+[а-я ]+ льготн(ый|ого) период.*")
			credit.GracePeriod = removeAllBesideCDPrSPlS(pattern.FindString(h.Text))
		}
	})

	// getting need documents
	coli.OnHTML(prompts["NeedDocuments"], func(h *colly.HTMLElement) {
		if h.Index == 1 {
			h.Text = strings.ReplaceAll(h.Text, "\n", "")
			h.Text = strings.Trim(h.Text, " ")
			credit.NeedDocuments = append(credit.NeedDocuments, strings.Split(h.Text, "   ")...)
		}
	})

	coli.Visit(url)

	return
}

// Getting (from site - parsing) credit by url, anchor and prompts v4. The function returns the struct of Credit type
func getEskhataCredit4(url, anchor string, prompts map[string]string) (credit structs.Credit) {
	coli := constructor()
	credit.Url = url + anchor
	for key, value := range prompts {
		prompts[key] = fmt.Sprintf("%s %s", anchor, value)
	}

	// getting title
	coli.OnHTML(prompts["Title"], func(h *colly.HTMLElement) {
		credit.Title = h.Text
	})

	// getting description
	coli.OnHTML(prompts["Description"], func(h *colly.HTMLElement) {
		credit.Description = utils.TernarOp(h.Index == 0, removeAllBesideCDPrSPlS(h.Text), credit.Description)
	})

	// getting purpose
	coli.OnHTML(prompts["Purpose"], func(h *colly.HTMLElement) {
		h.ForEachWithBreak("p", func(i int, e *colly.HTMLElement) bool {
			if ok, err := regexp.MatchString("Условия кредита:", e.Text); err == nil && ok {
				return false
			} else if ok, err := regexp.MatchString("Цель кредита:", e.Text); err == nil && ok {
				credit.Purpose = removeAllBesideCDPrSPlS(e.DOM.Next().Text())
			}

			return true
		})
	})

	// getting additional fields like conditions for repair and build or purchase house
	coli.OnHTML(prompts["AdditionalFiled"], func(h *colly.HTMLElement) {
		var conditionTitle1, conditionTitle2 string
		h.ForEachWithBreak("p", func(i int, e *colly.HTMLElement) bool {
			if ok, err := regexp.MatchString("(ремонт|строительство) жилья", e.Text); err == nil && ok {
				conditionTitle1 = removeAllBesideCDPrSPlS(e.Text)
			} else if ok, err := regexp.MatchString("покупк[а|у] жилья", e.Text); err == nil && ok {
				conditionTitle2 = removeAllBesideCDPrSPlS(e.Text)
				credit.AdditionalFields = map[string][]string{conditionTitle1: {}, conditionTitle2: {}}
				return false
			}
			return true
		})

		h.ForEachWithBreak("ul", func(i int, e *colly.HTMLElement) bool {
			if e.Index == 0 {
				e.Text = strings.ReplaceAll(e.Text, "\n", "")
				credit.AdditionalFields[conditionTitle1] = strings.Split(removeAllBesideCDPrSPlS(e.Text), "       ")
			} else if e.Index == 2 {
				e.Text = strings.ReplaceAll(e.Text, "\n", "")
				credit.AdditionalFields[conditionTitle2] = strings.Split(removeAllBesideCDPrSPlS(e.Text), "       ")
				return false
			}

			return true
		})
	})

	// getting need documents
	coli.OnHTML(prompts["NeedDocuments"], func(h *colly.HTMLElement) {
		if h.Index == 3 {
			h.Text = strings.ReplaceAll(h.Text, "\n", "")
			h.Text = strings.Trim(h.Text, " ")
			credit.NeedDocuments = append(credit.NeedDocuments, strings.Split(h.Text, "   ")...)
		}
	})

	coli.Visit(url)
	return
}

// Removing all beside cyrillic, digits, percent symbol and plus symbol
func removeAllBesideCDPrSPlS(s string) string {
	pattern := regexp.MustCompile(`[^А-Яа-я 0-9%\\+]+`)
	rgx := pattern.ReplaceAllString(s, "")
	return strings.Trim(rgx, " ")
}
