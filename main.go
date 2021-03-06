package main

import (
	"log"
	"os"
	"sync"

	"github.com/innovate-technologies/WHMCS-currency-update/whmcs"
)

const dryRun = false

var mainCurrency = "GBP"
var currencies map[string]whmcs.Currency
var products map[int64]whmcs.Product
var bannedProducts = map[int64]bool{47: true}

var api whmcs.API

func main() {
	var err error
	api = whmcs.New(os.Getenv("WHMCS_USERNAME"), os.Getenv("WHMCS_MD5"), os.Getenv("WHMCS_ACCESSKEY"), os.Getenv("WHMCS_URL"))

	currencies, err = api.GetCurrencies()
	if err != nil {
		panic(err)
	}
	for _, currency := range currencies {
		log.Printf("Rate %s: %f \n", currency.Code, currency.Rate)
	}

	products, err = api.GetAllProducts()
	if err != nil {
		panic(err)
	}

	goOverClientProducts()
}

func goOverClientProducts() {
	needsMore := true
	count := 0
	step := 100
	for needsMore {
		products, err := api.GetClientsProducts(count, step)
		if err != nil {
			panic(err)
		}

		var wg sync.WaitGroup
		for _, product := range products {
			wg.Add(1)
			go updatePrice(product, &wg)
		}
		wg.Wait()

		count += step
		log.Printf("Batch processed, total: %d \n", count)

		if len(products) == 0 {
			needsMore = false
		}
	}
}

func updatePrice(product whmcs.ClientProduct, wg *sync.WaitGroup) {
	defer wg.Done()

	if product.Billingcycle == "Free Account" || product.Recurringamount <= 0.0 {
		return // free product
	}

	if _, banned := bannedProducts[product.PID]; banned {
		return // we should never update this product
	}

	currency, err := api.GetClientCurrency(product.ClientID)
	if err != nil {
		log.Println(err)
		return
	}

	if currency == mainCurrency {
		//log.Printf("[%d] pays in GBP \n", product.ID)
		//return // no update needed here
	}

	parent := products[product.PID]
	var newPrice float64
	switch product.Billingcycle {
	case "Monthly":
		newPrice = parent.Pricing[mainCurrency].Monthly * currencies[currency].Rate
		break
	case "Quarterly":
		newPrice = parent.Pricing[mainCurrency].Quarterly * currencies[currency].Rate
		break
	case "Semi-Annually":
		newPrice = parent.Pricing[mainCurrency].Semiannually * currencies[currency].Rate
		break
	case "Annually":
		newPrice = parent.Pricing[mainCurrency].Annually * currencies[currency].Rate
		break
	case "Biennially":
		newPrice = parent.Pricing[mainCurrency].Biennially * currencies[currency].Rate
		break
	case "Triennially":
		newPrice = parent.Pricing[mainCurrency].Triennially * currencies[currency].Rate
		break
	}
	log.Printf("[%d] update price from %f to %f \n", product.ID, product.Recurringamount, newPrice)

	if newPrice == 0 {
		return // something went wrong
	}

	if !dryRun {
		err := api.UpdatePrice(product.ID, newPrice)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("[%d] updated price \n", product.ID)
	}
}
