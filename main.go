package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

// BaseHolder keeps API data
type BaseHolder struct {
	rates map[string]float64
}

// InitHolder with a given Base currency
// return new Holder object filled with the actual API data
func InitHolder(baseCur string) (*BaseHolder, error) {
	resp, err := http.Get("https://api.exchangeratesapi.io/latest?base=" + baseCur)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	if v, ok := data["error"]; ok {
		return nil, errors.New(v.(interface{}).(string))
	}

	if _, ok := data["rates"]; !ok {
		return nil, errors.New("Failed to get rates from API data")
	}
	rates := data["rates"].(interface{}).(map[string]interface{})

	var holder BaseHolder
	holder.rates = make(map[string]float64)
	for k, v := range rates {
		holder.rates[k] = v.(interface{}).(float64)
	}

	return &holder, nil
}

// getRate return given currency rate against base currency
func (h BaseHolder) getRate(curr string) (float64, error) {
	if v, ok := h.rates[curr]; ok {
		return v, nil
	}
	return 0, errors.New("Target currency " + curr + " is not found")
}

// Formal parameters check
func parseCmd() (amount float64, srcCurr, dstCurr string, err error) {
	if len(os.Args) != 4 {
		err = errors.New("Usage: " + os.Args[0] + " <amount:float> <src_symbol:string> <dst_symbol:string>")
		return
	}

	f, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		return
	}

	amount = f
	srcCurr = os.Args[2]
	dstCurr = os.Args[3]
	return
}

func main() {
	// I prefer to not use external packages (like `go-flags`), whereis possible
	amount, srcCurr, dstCurr, err := parseCmd()
	if err != nil {
		fmt.Println(err)
		return
	}

	holder, err := InitHolder(srcCurr)
	if err != nil {
		fmt.Println(err)
		return
	}

	m, err := holder.getRate(dstCurr)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(amount, srcCurr, "=", amount*m, dstCurr)
}
