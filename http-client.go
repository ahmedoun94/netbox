package netbox

import (
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

//HTTP request received by a server or to be sent by a client.
//To make a request with custom headers, use NewRequest and DefaultClient.Do
func Client(reqcurl string, Token string) []byte {
	req, err := http.NewRequest("GET", reqcurl, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Token "+Token)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)

	if err != nil {

		log.Fatal("Unable to reach Netbox")
	}

	//checking the validity of the token
	if resp.StatusCode == 403 {
		log.Fatal("Netbox token not valid")

	} else if resp.StatusCode != 403 {
		log.Info("Response received from netbox")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}
