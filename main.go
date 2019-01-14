package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Data struct {
	Quotes []struct {
		Bid    string `json:"bid"`
		Ask    string `json:"ask"`
		Symbol string `json:"currencyPairCode"`
	} `json:"quotes"`
}

func main() {
	res, err := http.Get("https://www.gaitameonline.com/rateaj/getrate")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var data Data
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}
	for _, d := range data.Quotes {
		fmt.Printf("Symbol:\t%v\n", d.Symbol)
		fmt.Printf("Bid:\t%v\n", d.Bid)
		fmt.Printf("Ask:\t%v\n", d.Ask)
		fmt.Println("==============")
	}
}
