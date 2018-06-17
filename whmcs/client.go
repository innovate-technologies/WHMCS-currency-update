package whmcs

import (
	"strings"

	resty "gopkg.in/resty.v0"
)

// API allows to call the WHMCS API
type API struct {
	username  string
	password  string
	accesskey string
	url       string
}

// New returns a new WHMCSAPI for an WHMCS install
func New(username, password, accesskey, url string) API {
	return API{
		username:  username,
		password:  password,
		accesskey: accesskey,
		url:       strings.Trim(url, "/"),
	}
}

func (a *API) createRequest(action string, data map[string]string) (*resty.Response, error) {
	data["username"] = a.username
	data["password"] = a.password
	data["accesskey"] = a.accesskey
	data["action"] = action
	data["responsetype"] = "json"
	return resty.R().SetFormData(data).Post(a.url + "/includes/api.php")
}
