package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
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
	if err := fetch(); err != nil {
		panic(err)
	}
}

func fetch() error {
	res, err := http.Get("https://www.gaitameonline.com/rateaj/getrate")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return err
	}
	defer c.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	counter := counter(c)
	counter = addCount(c, counter)
	key := name(counter)
	fmt.Println(key)
	c.Do("SET", key, body)

	b, _ := redis.Bytes(c.Do("GET", key))

	var data Data
	if err := json.Unmarshal(b, &data); err != nil {
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

func counter(c redis.Conn) int {
	// 初期値0
	count, _ := redis.Int(c.Do("GET", "count"))
	return count
}

func addCount(c redis.Conn, count int) int {
	count++
	c.Do("SET", "count", count)
	return count
}

func name(count int) string {
	return "ticker_" + strconv.Itoa(count)
}
