package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type logonMessage struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Domain   string `json:"domain,omitempty"`
}

type logonResponse struct {
	Scope        string `json:"scope"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	CSPAuthToken string `json:"cspAuthToken"`
}

func main() {

	userPtr := flag.String("user", "administrator", "User Name")
	passwordPtr := flag.String("password", "", "Password")
	domainPtr := flag.String("domain", "", "Domain to log into")
	serverPtr := flag.String("server", "api.mgmt.cloud.vmware.com", "FQDN/Hostname of vRA server")
	tokenTypePtr := flag.String("type", "refresh", "Token Type [refresh|access]")

	flag.Parse()

	if !strings.EqualFold("refresh", *tokenTypePtr) && !strings.EqualFold("access", *tokenTypePtr) {
		log.Fatal("Invalid token type specifed. Not 'refresh' or 'access'")
	}

	logonJSON, err := json.Marshal(logonMessage{*userPtr, *passwordPtr, *domainPtr})

	if err != nil {
		fmt.Print("Error occurred")
	}

	postURL := fmt.Sprintf("https://%s/csp/gateway/am/api/login", *serverPtr)

	if strings.EqualFold("refresh", *tokenTypePtr) {
		postURL += "?access_token"
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Post(postURL, "application/json", bytes.NewBuffer(logonJSON))
	if err != nil {
		log.Fatal(err)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	if resp.StatusCode != 200 {
		fmt.Println("Non-200 Status Code received", resp.StatusCode, string(respBody[:]))
		log.Fatal("Non-200 Status Code received")
	}

	var logonResp logonResponse
	json.Unmarshal(respBody, &logonResp)
	if strings.EqualFold("refresh", *tokenTypePtr) {
		fmt.Print(logonResp.RefreshToken)
	} else if strings.EqualFold("access", *tokenTypePtr) {
		fmt.Print(logonResp.CSPAuthToken)
	}

}
