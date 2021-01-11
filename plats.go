package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type PlatRecording struct {
	EntryNumber, Year, PartyFrom, PartyTo, Date string
}

type DocumentDetail struct {
	ReceivingParty, Description string
}

//GetAllPlatsRecordedSince will query the Utah County Recorder's office and will pull all
//subdivision plats (designation "S PLAT") that have been recorded since the date 'd' passed in
func GetAllPlatsRecordedSince(d time.Time) (map[string]PlatRecording, error) { 
	//todo add funcionality to schedule this job and save the last recorded time it ran so it picks up every one

	offset := 0
	doc, err := getSPlatHTMLSince(d, offset)
	if err != nil {
		return nil, fmt.Errorf("Error getting S PLAT HTML since %v. Error: %v", d, err)
	}

	plats := map[string]PlatRecording{}
	numFound, err := extractPlatRecordingsFromDocToMap(doc, plats)
	if err != nil {
		return nil, fmt.Errorf("Unable to extract plat recordings to map: %s", err.Error())
	}

	for numFound == 100 {
		offset += 100
		doc, err := getSPlatHTMLSince(d, offset)
		if err != nil {
			return nil, fmt.Errorf("Error getting S PLAT HTML since %v. Error: %v", d, err)
		}

		numFound, err = extractPlatRecordingsFromDocToMap(doc, plats)
		if err != nil {
			return nil, fmt.Errorf("Unable to extract plat recordings to map: %s", err.Error())
		}
	}

	return plats, nil
}

func extractPlatRecordingsFromDocToMap(doc *html.Node, platsMap map[string]PlatRecording) (numFound int, err error) {

	dataRows, err := htmlquery.QueryAll(doc, "//tr")
	if err != nil {
		return 0, fmt.Errorf("Error querying for all S PLAT data rows: %s", err.Error())
	}
	if len(dataRows) == 0 {
		return 0, fmt.Errorf("Unable to locate any S PLAT data rows")
	}

	for _, row := range dataRows {
		e := getEntryNumber(row)
		platsMap[e] = getPlatRecording(e, row)
	}

	return len(dataRows), nil
}

func getEntryNumber(row *html.Node) string {
	a := htmlquery.FindOne(row, "//a")
	return htmlquery.InnerText(a)
}

func getPlatRecording(entryNum string, row *html.Node) PlatRecording {
	y := htmlquery.FindOne(row, "//td[3]")
	d := htmlquery.FindOne(row, "//td[4]")
	f := htmlquery.FindOne(row, "//td[5]")
	t := htmlquery.FindOne(row, "//td[6]")

	return PlatRecording{
		EntryNumber: entryNum,
		Year: htmlquery.InnerText(y),
		Date: htmlquery.InnerText(d),
		PartyFrom: htmlquery.InnerText(f),
		PartyTo: htmlquery.InnerText(t),
	}
}

//getSPlatHTMLSince will grab the HTML page from the Utah County Recorder's website
//Offset is how pagination is handled. An offset of '0' will yield the first page of results. Pages are 100 entries in length at max.
func getSPlatHTMLSince(d time.Time, offset int) (*html.Node, error) {
	entryDate := d.Format("1/02/2006")
	entryDate = strings.Replace(entryDate, "/", "%2F", 3)
	doc, err := htmlquery.LoadURL(fmt.Sprintf("http://www.utahcounty.gov/LandRecords/DocKoi.asp?avKoi=S+PLAT&avEntryDate=%s&Submit=Search&offset=%d", entryDate, offset))
	if err != nil {
		return nil, fmt.Errorf("Error performing S PLAT search: %v", err.Error())
	}

	return doc, nil
}
