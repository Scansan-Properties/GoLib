package ScanSanAPI

import (
	"encoding/json"
	"fmt"
	l "log/slog"
	"net/url"
	"strings"
)

type SalesDemandResponse struct {
	AreaCode     string  `json:"area_code"`
	AreaCodeType string  `json:"area_code_type"`
	TargetMonth  string  `json:"target_month"`
	TargetYear   float64 `json:"target_year"`
	ResponseTime string  `json:"response_time"`
	Data         struct {
		SaleDemand []struct {
			TotalPropertiesForSale float64 `json:"total_properties_for_sale"`
			AverageTransactionsPcm float64 `json:"average_transactions_pcm"`
			MonthsOfInventory      float64 `json:"months_of_inventory"`
			TurnoverPercentagePcm  float64 `json:"turnover_percentage_pcm"`
			DaysOnMarket           float64 `json:"days_on_market"`
			AveragePrice           float64 `json:"average_price"`
			MedianPrice            float64 `json:"median_price"`
			Currency               string  `json:"currency"`
			Rating                 string  `json:"rating"`
		} `json:"sale_demand"`
		AdditionalSaleData struct {
			PriceCategory []struct {
				MinPriceRange     float64 `json:"min_price_range"`
				MaxPriceRange     float64 `json:"max_price_range"`
				PropertiesForSale float64 `json:"properties_for_sale"`
				AverageOtmDays    float64 `json:"average_otm_days"`
				Currency          string  `json:"currency"`
			} `json:"price_category"`
			BedroomsCategory []struct {
				NoOfBedrooms      float64 `json:"no_of_bedrooms"`
				PropertiesForSale float64 `json:"properties_for_sale"`
				AveragePrice      float64 `json:"average_price"`
				MedianPrice       float64 `json:"median_price"`
				AverageOtmDays    float64 `json:"average_otm_days"`
				Currency          string  `json:"currency"`
			} `json:"bedrooms_category"`
			TypeCategory []struct {
				PropertyType      string  `json:"property_type"`
				PropertiesForSale float64 `json:"properties_for_sale"`
				AveragePrice      float64 `json:"average_price"`
				MedianPrice       float64 `json:"median_price"`
				AverageOtmDays    float64 `json:"average_otm_days"`
				Currency          string  `json:"currency"`
			} `json:"type_category"`
			OtmCategory []struct {
				MinOtmMonths      float64 `json:"min_otm_months"`
				MaxOtmMonths      float64 `json:"max_otm_months"`
				PropertiesForSale float64 `json:"properties_for_sale"`
			} `json:"otm_category"`
		} `json:"additional_sale_data"`
	} `json:"data"`
}

func (r *Request) SalesDemand(AreaCode string) (string, error) {

	// Make sure the Code is clean
	d := strings.ReplaceAll(AreaCode, " ", "")
	d = strings.TrimSpace(d)

	out := ""

	ScanSanURL := fmt.Sprintf("%s/area_codes/%s/sale/demand?additional_rental_data=true", r.Server, url.QueryEscape(d))
	body, _, err := webRequest(ScanSanURL)

	if r.DEBUG {
		l.With("Body: ", string(body)).Debug("ScanSanAPI Response")
		l.With("ScanSanURL: ", ScanSanURL).Debug("ScanSanURL")
	}

	if err != nil {
		return out, err
	}

	SD := SalesDemandResponse{}
	err = json.Unmarshal(body, &SD)
	if err != nil {
		l.With("error", err, "body", string(body)).Error("Error unmarshalling ML Response")
		return out, err
	}

	out = "### Area Sales Demand for " + SD.AreaCode + "\n"
	out = out + "- Area Code Type: " + SD.AreaCodeType + "\n"

	for _, v := range SD.Data.SaleDemand {

		out = out + fmt.Sprintf("- Total Properties for Sale: %d\n", int(v.TotalPropertiesForSale))
		out = out + fmt.Sprintf("- Average Transactions PCM: %d\n", int(v.AverageTransactionsPcm))
		out = out + fmt.Sprintf("- Months of Inventory: %d\n", int(v.MonthsOfInventory))
		out = out + fmt.Sprintf("- Turnover Percentage PCM: %d\n", int(v.TurnoverPercentagePcm))
		out = out + fmt.Sprintf("- Days on Market: %d\n", int(v.DaysOnMarket))
		out = out + fmt.Sprintf("- Average Price: £%d\n", int(v.AveragePrice))
		out = out + fmt.Sprintf("- Median Price: £%d\n", int(v.MedianPrice))
		out = out + fmt.Sprintf("- Currency: %s\n", v.Currency)
		out = out + fmt.Sprintf("- Rating: %s\n", v.Rating)

	}
	out = out + "\n"

	if len(SD.Data.AdditionalSaleData.PriceCategory) > 0 {
		out = out + "### Additional Sales Data\n"
		for _, v := range SD.Data.AdditionalSaleData.PriceCategory {
			out = out + fmt.Sprintf("- Price Range: £%d - £%d\n", int(v.MinPriceRange), int(v.MaxPriceRange))
			out = out + fmt.Sprintf("- Properties for Sale: %d\n", int(v.PropertiesForSale))
			out = out + fmt.Sprintf("- Average OTM Days: %d\n", int(v.AverageOtmDays))

		}
		out = out + "\n"
	}

	if len(SD.Data.AdditionalSaleData.BedroomsCategory) > 0 {
		out = out + "### Additional Sales Data - Number of Bedrooms\n"
		for _, v := range SD.Data.AdditionalSaleData.BedroomsCategory {
			out = out + fmt.Sprintf("- Bedrooms: %d\n", int(v.NoOfBedrooms))
			out = out + fmt.Sprintf("- Properties for Sale: %d\n", int(v.PropertiesForSale))
			out = out + fmt.Sprintf("- Average Price: £%d\n", int(v.AveragePrice))
			out = out + fmt.Sprintf("- Median Price: £%d\n", int(v.MedianPrice))
			out = out + fmt.Sprintf("- Average OTM Days: %d\n", int(v.AverageOtmDays))
		}
		out = out + "\n"
	}

	if len(SD.Data.AdditionalSaleData.TypeCategory) > 0 {
		out = out + "### Additional Sales Data - Property Type\n"
		for _, v := range SD.Data.AdditionalSaleData.TypeCategory {
			out = out + fmt.Sprintf("- Property Type: %s\n", v.PropertyType)
			out = out + fmt.Sprintf("- Properties for Sale: %d\n", int(v.PropertiesForSale))
			out = out + fmt.Sprintf("- Average Price: £%d\n", int(v.AveragePrice))
			out = out + fmt.Sprintf("- Median Price: £%d\n", int(v.MedianPrice))
			out = out + fmt.Sprintf("- Average OTM Days: %d\n", int(v.AverageOtmDays))
		}
		out = out + "\n"
	}

	if len(SD.Data.AdditionalSaleData.OtmCategory) > 0 {
		out = out + "### Additional Sales Data - On the Market\n"
		for _, v := range SD.Data.AdditionalSaleData.OtmCategory {
			out = out + fmt.Sprintf("- Min Number of Months on the Market: %d\n", int(v.MinOtmMonths))
			out = out + fmt.Sprintf("- Max Number of Months on the Market: %d\n", int(v.MaxOtmMonths))
			out = out + fmt.Sprintf("- Properties for Sale: %d\n", int(v.PropertiesForSale))
		}
	}

	if r.DEBUG {
		fmt.Printf("--  \n%s\n-- \n", out)
	}

	return out, err
}
