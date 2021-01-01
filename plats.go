package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

//Want to get all plats recorded since the

type PlatRecording struct {
	Year, PartyFrom, PartyTo, Date string
}

type DocumentDetail struct {
	ReceivingParty, Description string
}

func GetAllPlatsRecordedSince(d time.Time) (map[string]PlatRecording, error) { 
	//todo add funcionality to schedule this job and save the last recorded time it ran so it picks up every one
	params := url.Values{
		"Submit": {"Search"},
		"avKoi": {"S+PLAT"},
		"avEntryDate": {"11-30-2020"},
	}
	url := url.URL{
		Host: "www.utahcounty.gov",
		Scheme: "http",
		Path: "LandRecords/DocKoi.asp",
		RawQuery: params.Encode(),
	}

	resp, err := http.Get(url.String())
	if err != nil {
		return nil, fmt.Errorf("Error with the request to get the webpage at URL: %s", url.String())
	} 
	if resp.StatusCode >= 300 {
		msg, err := getFailedResponseError(resp)
		if err != nil {
			return nil, fmt.Errorf("%s | Unable to get >300 response error message", err.Error())
		}
		return nil, errors.New(msg)
	}

	data, err := extractPlatDocumentDataFromHtml(resp)
	if err != nil {
		return nil, fmt.Errorf("Unable to extract plat document data from HTML. Error: %s", err.Error())
	}
	
	return data, nil
}

func getFailedResponseError(resp *http.Response) (string, error) {
	if resp.StatusCode < 300 {
		return "", fmt.Errorf("Cannot perform operation on request that didn't fail")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Unable to read response body")
	}

	return fmt.Sprintf("Error with request to URL: '%s'. Status code: '%d'. Body: '%s'", resp.Request.URL.String(), resp.StatusCode, body), nil
}

func extractPlatDocumentDataFromHtml(resp *http.Response) (map[string]PlatRecording, error) {
	
	defer resp.Body.Close()
	z := html.NewTokenizer(resp.Body)

    content := make(map[string]PlatRecording)

    // While have not hit the </html> tag
    for z.Token().Data != "html" {
        tt := z.Next()
		if tt != html.StartTagToken { continue }
		
		t := z.Token()
		if t.Data != "td" { continue }

		inner := z.Next()
		if inner != html.TextToken { continue } //only grab tokens inside tds that have text in it

		z.
		addDataToMap(content, inner, z)

    }
    // Print to check the slice's content
    fmt.Println(content)

	return nil, nil
}

func addDataToMap(bucket *map[string]PlatRecording, inner html.Token, z html.Tokenizer) {
	
	text := (string)(z.Text())
	text = strings.TrimSpace(text)

	if isEntryNumber, _ := regexp.MatchString(`\d{6}`, text); isEntryNumber {
		//will need to grab the year of the sibling td
		return
	}
	if isDate, _ := regexp.MatchString(`\d{1,2}\/\d{1,2}\/\d{4}`, text); isDate {
		
		return
	}


}