package whmcs

import (
	"encoding/json"
	"errors"
	"strconv"
)

type Currency struct {
	ID   int64
	Code string
	Rate float64
}

// GetCurrencies gets back the currencies in WHMCS
func (a *API) GetCurrencies() (map[string]Currency, error) {
	resp, err := a.createRequest("GetCurrencies", map[string]string{})
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

	currencies := map[string]Currency{}
	for _, c := range (data["currencies"].(map[string]interface{}))["currency"].([]interface{}) {
		new := Currency{}
		d := c.(map[string]interface{})

		new.ID = int64(d["id"].(float64))
		new.Code = d["code"].(string)
		new.Rate, _ = strconv.ParseFloat(d["rate"].(string), 64)

		currencies[new.Code] = new
	}

	return currencies, nil
}
