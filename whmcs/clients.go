package whmcs

import (
	"encoding/json"
	"errors"
	"fmt"
)

// GetClientCurrency gets back the currency of a client in WHMCS
func (a *API) GetClientCurrency(id int64) (string, error) {
	resp, err := a.createRequest("GetClientsDetails", map[string]string{
		"clientid": fmt.Sprintf("%d", id),
	})
	if err != nil {
		return "", err
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return "", err
	}

	if data["result"] != "success" {
		return "", errors.New(string(resp.Body()))
	}

	return (data["client"].(map[string]interface{}))["currency_code"].(string), nil
}
