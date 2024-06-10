package ScanSanAPI

import (
	"encoding/json"
	"fmt"
	l "log/slog"
	"net/url"
	"strings"
)

type AreaSummaryResponse struct {
	AreaCode     string `json:"area_code"`
	AreaCodeType string `json:"area_code_type"`
	ResponseTime string `json:"response_time"`
	Data         []struct {
		TotalProperties               float32   `json:"total_properties"`
		TotalPropertiesSoldInLast5Yrs float32   `json:"total_properties_sold_in_last_5yrs"`
		SoldPriceRangeInLast5Yrs      []float32 `json:"sold_price_range_in_last_5yrs"`
		CurrentRentListings           float32   `json:"current_rent_listings"`
		CurrentRentListingsPcmRange   []float32 `json:"current_rent_listings_pcm_range"`
		CurrentSaleListings           float32   `json:"current_sale_listings"`
		CurrentRentListingsPriceRange []float32 `json:"current_rent_listings_price_range"`
	} `json:"data"`
}

func (r *Request) AreaSummary(AreaCode string) (string, error) {

	// Make sure the Code is clean
	d := strings.ReplaceAll(AreaCode, " ", "")
	d = strings.TrimSpace(d)

	out := ""

	ScanSanURL := fmt.Sprintf("%s/area_codes/%s/summary", r.Server, url.QueryEscape(d))
	body, _, err := webRequest(ScanSanURL)

	if r.Display {
		fmt.Printf("RAW RESPONSE \n%s\n-- \n", body)
	}

	if r.DEBUG {
		l.With("Body: ", string(body)).Debug("ScanSanAPI Response")
		l.With("ScanSanURL: ", ScanSanURL).Debug("ScanSanURL")
	}

	if err != nil {
		return out, err
	}

	AS := AreaSummaryResponse{}
	err = json.Unmarshal(body, &AS)
	if err != nil {
		l.With("error", err, "body", string(body)).Error("Error unmarshalling ML Response")
		return out, err
	}

	out = "### Area Summary for " + AS.AreaCode + "\n"
	out = out + "- Area Code Type: " + AS.AreaCodeType + "\n"

	for _, v := range AS.Data {
		// out = out + fmt.Sprintf("- Total Properties: %d\n", int(v.TotalProperties))
		out = out + fmt.Sprintf("- Total Properties Sold in Last 5 Years: %d\n", int(v.TotalPropertiesSoldInLast5Yrs))
		out = out + fmt.Sprintf("- Sold Price Range in Last 5 Years: £%d - %d\n", int(v.SoldPriceRangeInLast5Yrs[0]), int(v.SoldPriceRangeInLast5Yrs[1]))
		out = out + fmt.Sprintf("- Current Rent Listings: %d\n", int(v.CurrentRentListings))
		out = out + fmt.Sprintf("- Current Rent Listings PCM Range: £%d - %d\n", int(v.CurrentRentListingsPcmRange[0]), int(v.CurrentRentListingsPcmRange[1]))
		out = out + fmt.Sprintf("- Current Sale Listings: %d\n", int(v.CurrentSaleListings))
		out = out + fmt.Sprintf("- Current Sale Listings Price Range: £%d - %d\n", int(v.CurrentRentListingsPriceRange[0]), int(v.CurrentRentListingsPriceRange[1]))
	}

	if r.DEBUG {
		fmt.Printf("--  \n%s\n-- \n", out)
	}
	return out, nil
}
