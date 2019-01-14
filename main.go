package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const interval = 5 * time.Second

type Data struct {
	Quotes []struct {
		Bid    string `json:"bid"`
		Ask    string `json:"ask"`
		Symbol string `json:"currencyPairCode"`
	} `json:"quotes"`
}

func main() {
	start := time.Now()
	wait := start.Truncate(interval).Add(interval).Sub(start)
	time.Sleep(wait)
	ticker := time.Tick(interval)
	for now := range ticker {
		fmt.Println(now.String())
		if err := fetch(); err != nil {
			panic(err)
		}
	}
}

func fetch() error {
	res, err := http.Get("https://www.gaitameonline.com/rateaj/getrate")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var data Data
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}
	for _, d := range data.Quotes {
		fmt.Printf("Symbol:\t%v\n", d.Symbol)
		fmt.Printf("Bid:\t%v\n", d.Bid)
		fmt.Printf("Ask:\t%v\n", d.Ask)
		fmt.Println("==============")
	}
	return nil
}
