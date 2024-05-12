package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type EtfInfo struct {
	Title              string
	Replication        string
	Earnings           string
	TotalExpenseRatio  string
	TrackingDifference string
	FundSize           string
}

func main() {

	// Retrieving isin number for user input
	// An International Securities Identification Number, or ISIN, is a unique twelve-digit code that is assigned to every security issuance in the world. This number is used to facilitate the trading, clearing, and settlement of securities transactions, especially across borders.
	var isin string
	fmt.Print("Enter ISIN Number: ")
	fmt.Scanln(&isin)

	// Sample ISINs Avalible from TrackingDifferences.com
	// Vanguard S&P 500 ETF: IE00B3XXRP09
	// More US ETFs: https://www.trackingdifferences.com/ETF/Region/North%20America

	etfInfo := EtfInfo{}

	c := colly.NewCollector(colly.AllowedDomains("www.trackingdifferences.com", "trackingdifferences.com"))

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US;q=0.9")
		fmt.Printf("Visiting %s\n", r.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Printf("Error while scraping: %s\n", e.Error())
	})

	c.OnHTML("h1.page-title", func(h *colly.HTMLElement) {
		etfInfo.Title = h.Text
	})

	c.OnHTML("div.descfloat p.desc", func(h *colly.HTMLElement) {
		selection := h.DOM

		childNodes := selection.Children().Nodes
		if len(childNodes) == 3 {
			description := cleanDesc(selection.Find("span.desctitle").Text())
			value := selection.FindNodes(childNodes[2]).Text()

			switch description {
			case "Replication":
				etfInfo.Replication = value
				break
			case "TER":
				etfInfo.TotalExpenseRatio = value
				break
			case "TD":
				etfInfo.TrackingDifference = value
				break
			case "Earnings":
				etfInfo.Earnings = value
				break
			case "Fund size":
				etfInfo.FundSize = value
				break
			}
		}
	})

	c.OnScraped(func(r *colly.Response) {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "\t")
		enc.Encode(etfInfo)
	})

	c.Visit(scrapeUrl(isin))

}

func cleanDesc(s string) string {
	return strings.TrimSpace(s)
}

func scrapeUrl(isin string) string {
	return "https://www.trackingdifferences.com/ETF/ISIN/" + isin
}