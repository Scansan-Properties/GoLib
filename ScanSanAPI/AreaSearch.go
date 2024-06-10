package ScanSanAPI

import (
	"encoding/json"
	"fmt"
	l "log/slog"
	"net/url"
	"strings"
)

type AreaSearchResponse struct {
	SearchQuery  string `json:"search_query"`
	SearchFound  string `json:"search_found"`
	ResponseTime string `json:"response_time"`
	Data         []struct {
		PostcodeDistrict string   `json:"postcode_district"`
		PostcodeCount    int      `json:"postcode_count"`
		PostcodeList     []any    `json:"postcode_list"`
		Boroughs         []string `json:"boroughs"`
		Wards            []string `json:"wards"`
		StreetCount      int      `json:"street_count"`
	} `json:"data"`
}

func (r *Request) AreaSearch(searchTerm string) (string, []string, error) {

	// Make sure the Code is clean
	searchTerm = strings.TrimSpace(searchTerm)

	out := ""
	districts := []string{}

	ScanSanURL := fmt.Sprintf("%s/area_codes/search?area_name=%s", r.Server, url.QueryEscape(searchTerm))
	body, code, err := webRequest(ScanSanURL)

	if r.Display {
		fmt.Printf("RAW RESPONSE \n%s\n-- \n", body)
	}

	if r.DEBUG {
		l.With("Body: ", string(body)).Debug("ScanSanAPI Response")
		l.With("ScanSanURL: ", ScanSanURL).Debug("ScanSanURL")
	}

	if err != nil {
		return out, districts, err
	}

	if code == 404 {
		return fmt.Sprintf("No area information found for search term '%s'\n", searchTerm), districts, nil
	}

	AS := AreaSearchResponse{}
	err = json.Unmarshal(body, &AS)
	if err != nil {
		l.With("error", err, "body", string(body)).Error("Error unmarshalling ML Response")
		return out, districts, err
	}

	out = "### Area Search for " + AS.SearchQuery + "\n"
	out = out + "- Type of Area found: " + AS.SearchFound + "\n"

	for _, v := range AS.Data {
		districts = append(districts, v.PostcodeDistrict)
		out = out + fmt.Sprintf("## Postcode District in this area: %s\n", v.PostcodeDistrict)
		out = out + fmt.Sprintf("- Number of Postcodes in %s: %d\n", v.PostcodeDistrict, v.PostcodeCount)
		out = out + fmt.Sprintf("- Number of Streets in %s: %d\n", v.PostcodeDistrict, v.StreetCount)
		out = out + fmt.Sprintf("- Boroughs: %s\n", strings.Join(v.Boroughs, ", "))
		out = out + fmt.Sprintf("- Wards: %s\n\n", strings.Join(v.Wards, ", "))
	}

	if r.DEBUG {
		fmt.Printf("--  \n%s\n-- \n", out)
	}
	return out, districts, nil
}
