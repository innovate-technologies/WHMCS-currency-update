package whmcs

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type Product struct {
	// incomplete, for now we don't need everything
	PID     int64
	GID     int64
	Name    string
	Pricing map[string]ProductPrice
}

type ProductPrice struct {
	Monthly      float64
	Quarterly    float64
	Semiannually float64
	Annually     float64
	Biennially   float64
	Triennially  float64
}

type ClientProduct struct {
	ID              int64
	PID             int64
	ClientID        int64
	Billingcycle    string
	Recurringamount float64
}

// GetAllProducts gets back all products in WHMCS
func (a *API) GetAllProducts() (map[int64]Product, error) {
	resp, err := a.createRequest("GetProducts", map[string]string{})
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return nil, err
	}

	if data["result"] != "success" {
		return nil, errors.New(string(resp.Body()))
	}

	products := map[int64]Product{}
	for _, c := range (data["products"].(map[string]interface{}))["product"].([]interface{}) {
		new := Product{}
		d := c.(map[string]interface{})

		new.PID = int64(d["pid"].(float64))
		new.GID = int64(d["gid"].(float64))
		new.Name = d["name"].(string)
		new.Pricing = map[string]ProductPrice{}

		pricing := d["pricing"].(map[string]interface{})
		for code, price := range pricing {
			priceMap := price.(map[string]interface{})

			monthly, _ := strconv.ParseFloat(priceMap["monthly"].(string), 64)
			quarterly, _ := strconv.ParseFloat(priceMap["quarterly"].(string), 64)
			semiannually, _ := strconv.ParseFloat(priceMap["semiannually"].(string), 64)
			annually, _ := strconv.ParseFloat(priceMap["annually"].(string), 64)
			biennially, _ := strconv.ParseFloat(priceMap["biennially"].(string), 64)
			triennially, _ := strconv.ParseFloat(priceMap["triennially"].(string), 64)

			new.Pricing[code] = ProductPrice{
				Monthly:      monthly,
				Quarterly:    quarterly,
				Semiannually: semiannually,
				Annually:     annually,
				Biennially:   biennially,
				Triennially:  triennially,
			}
		}

		products[new.PID] = new
	}

	return products, nil
}

func (a *API) GetClientsProducts(start, limit int) ([]ClientProduct, error) {
	resp, err := a.createRequest("GetClientsProducts", map[string]string{
		"limitstart": fmt.Sprintf("%d", start),
		"limitnum":   fmt.Sprintf("%d", limit),
	})
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return nil, err
	}

	if data["result"] != "success" {
		return nil, errors.New(string(resp.Body()))
	}

	products := []ClientProduct{}

	if data["products"] == nil {
		return products, nil
	}

	for _, c := range (data["products"].(map[string]interface{}))["product"].([]interface{}) {
		new := ClientProduct{}
		d := c.(map[string]interface{})

		new.ID = int64(d["id"].(float64))
		new.PID = int64(d["pid"].(float64))
		new.ClientID = int64(d["clientid"].(float64))
		new.Billingcycle = d["billingcycle"].(string)
		new.Recurringamount, _ = strconv.ParseFloat(d["recurringamount"].(string), 64)

		products = append(products, new)
	}

	return products, nil
}

func (a *API) UpdatePrice(id int64, amount float64) error {
	// UpdateClientProduct
	resp, err := a.createRequest("UpdateClientProduct", map[string]string{
		"serviceid":       fmt.Sprintf("%d", id),
		"recurringamount": fmt.Sprintf("%.2f", amount),
	})
	if err != nil {
		return err
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return err
	}

	if data["result"] != "success" {
		return errors.New(string(resp.Body()))
	}

	return nil
}
