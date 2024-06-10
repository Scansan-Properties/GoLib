package main

import (
    "fmt"
    
    _ "github.com/joho/godotenv/autoload"
    
    "github.com/SpaceTent/GoLib/ScanSanAPI"
)

func main() {
    
    r := ScanSanAPI.NewRequest()
    err := error(nil)
    // r.DEBUG = true
    
    summary := ""
    summary, err = r.AreaSummary("E14")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    fmt.Printf("summary:\n%+v\n", summary)
    
    rentdemand := ""
    rentdemand, err = r.AreaSummary("E14")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    fmt.Printf("rentdemand:\n%+v\n", rentdemand)
    
    salesdemand := ""
    salesdemand, err = r.SalesDemand("E14")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    fmt.Printf("salesdemand:\n%+v\n", salesdemand)
    
    AllDetails, Districts, err := r.AreaSearch("Barking")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Details:\n%+v\n", AllDetails)
    fmt.Printf("Districts:\n%+v\n", Districts)
    
}
