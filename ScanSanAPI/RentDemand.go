package ScanSanAPI

import (
	"encoding/json"
	"fmt"
	l "log/slog"
	"net/url"
	"strings"
)

type RentDemandResponse struct {
	AreaCode     string `json:"area_code"`
	AreaCodeType string `json:"area_code_type"`
	TargetMonth  string `json:"target_month"`
	TargetYear   int    `json:"target_year"`
	ResponseTime string `json:"response_time"`
	Data         struct {
		RentalDemand []struct {
			TotalPropertiesForRent float64 `json:"total_properties_for_rent"`
			AverageTransactionsPcm float64 `json:"average_transactions_pcm"`
			MonthsOfInventory      float64 `json:"months_of_inventory"`
			TurnoverPercentagePcm  float64 `json:"turnover_percentage_pcm"`
			DaysOnMarket           float64 `json:"days_on_market"`
			AverageRentPcm         float64 `json:"average_rent_pcm"`
			MedianRentPcm          float64 `json:"median_rent_pcm"`
			Currency               string  `json:"currency"`
			Rating                 string  `json:"rating"`
		} `json:"rental_demand"`
		AdditionalRentalData struct {
			PriceCategory []struct {
				PriceCategoryPcm  string  `json:"price_category_pcm"`
				PropertiesForRent float64 `json:"properties_for_rent"`
			} `json:"price_category"`
			BedroomsCategory []struct {
				NoOfBedrooms      float64 `json:"no_of_bedrooms"`
				PropertiesForRent float64 `json:"properties_for_rent"`
				AverageRentPcm    float64 `json:"average_rent_pcm"`
				MedianRentPcm     float64 `json:"median_rent_pcm"`
				Currency          string  `json:"currency"`
			} `json:"bedrooms_category"`
			TypeCategory []struct {
				Type              string  `json:"type"`
				PropertiesForRent float64 `json:"properties_for_rent"`
				AverageRentPcm    float64 `json:"average_rent_pcm"`
				MedianRentPcm     float64 `json:"median_rent_pcm"`
			} `json:"type_category"`
		} `json:"additional_rental_data"`
	} `json:"data"`
}

func (r *Request) RentDemand(AreaCode string) (string, error) {

	// Make sure the Code is clean
	d := strings.ReplaceAll(AreaCode, " ", "")
	d = strings.TrimSpace(d)

	out := ""

	ScanSanURL := fmt.Sprintf("%s/area_codes/%s/rent/demand?additional_rental_data=true", r.Server, url.QueryEscape(d))
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

	RD := RentDemandResponse{}
	err = json.Unmarshal(body, &RD)
	if err != nil {
		l.With("error", err, "body", string(body)).Error("Error unmarshalling ML Response")
		return out, err
	}

	out = "### Area Rental Demand for " + RD.AreaCode + "\n"
	out = out + "- Area Code Type: " + RD.AreaCodeType + "\n"

	for _, v := range RD.Data.RentalDemand {

		out = out + fmt.Sprintf("- Total Properties for Rent: %d\n", int(v.TotalPropertiesForRent))
		out = out + fmt.Sprintf("- Average Transactions PCM: £%d\n", int(v.AverageTransactionsPcm))
		out = out + fmt.Sprintf("- Months of Inventory: %d\n", int(v.MonthsOfInventory))
		out = out + fmt.Sprintf("- Turnover Percentage PCM: %d\n", int(v.TurnoverPercentagePcm))
		out = out + fmt.Sprintf("- Days on Market: %d\n", int(v.DaysOnMarket))
		out = out + fmt.Sprintf("- Average Rent PCM: £%d\n", int(v.AverageRentPcm))
		out = out + fmt.Sprintf("- Median Rent PCM: £%d\n", int(v.MedianRentPcm))
		out = out + fmt.Sprintf("- Currency: %s\n", v.Currency)
		out = out + fmt.Sprintf("- Rating: %s\n", v.Rating)

	}
	out = out + "\n"

	out = out + "### Additional Rental by Price Category " + RD.AreaCode + "\n"
	for _, v := range RD.Data.AdditionalRentalData.PriceCategory {
		out = out + fmt.Sprintf("- Price Category PCM: %s\n", v.PriceCategoryPcm)
		out = out + fmt.Sprintf("- Properties for Rent: %d\n", int(v.PropertiesForRent))
	}
	out = out + "\n"

	out = "### Additional Rental by Number of Bedrooms " + RD.AreaCode + "\n"
	for _, v := range RD.Data.AdditionalRentalData.BedroomsCategory {
		out = out + fmt.Sprintf("- Number of Bedrooms: %d\n", int(v.NoOfBedrooms))
		out = out + fmt.Sprintf("- Properties for Rent: %d\n", int(v.PropertiesForRent))
		out = out + fmt.Sprintf("- Average Rent PCM: £%d\n", int(v.AverageRentPcm))
		out = out + fmt.Sprintf("- Median Rent PCM: £%d\n", int(v.MedianRentPcm))
	}
	out = out + "\n"

	out = "### Additional Rental by Property Type " + RD.AreaCode + "\n"
	for _, v := range RD.Data.AdditionalRentalData.TypeCategory {
		out = out + fmt.Sprintf("- Type: %s\n", v.Type)
		out = out + fmt.Sprintf("- Properties for Rent: %d\n", int(v.PropertiesForRent))
		out = out + fmt.Sprintf("- Average Rent PCM: £%d\n", int(v.AverageRentPcm))
		out = out + fmt.Sprintf("- Median Rent PCM: £%d\n", int(v.MedianRentPcm))
	}

	if r.DEBUG {
		fmt.Printf("--  \n%s\n-- \n", out)
	}

	return out, nil
}
