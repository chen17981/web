package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	pricesStr string
	prices    map[string]float64
)

func handle(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "form.html")

	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		webinput := strings.TrimSpace(r.FormValue("prices"))
		if len(webinput) != 0 {
			fmt.Fprintf(w, "Product Prices = %s\n", webinput)
			names, nums, ok := validatePricelist(webinput)
			if !ok {
				fmt.Fprintf(w, "%s is not valid\n", webinput)
				return
			}

			//prices is a map which records the pair of product and its price
			prices = setPrices(names, nums)
		}

		//fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		webinput = strings.TrimSpace(r.FormValue("name"))
		fmt.Fprintf(w, "Shopping Items = %s\n", webinput)

		items, ok := validateShoppingItems(webinput, prices)
		if !ok {
			fmt.Fprintf(w, "%s is not valid\n", webinput)
			return
		}

		calculate(items, prices, true, w)

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {

	pricesStr = "CH1,3.11,AP1,6.00,CF1,11.23,MK1,4.75,OM1,3.69"

	names, nums, ok := validatePricelist(pricesStr)
	if !ok {
		return
	}

	//prices is a map which records the pair of product and its price
	prices = setPrices(names, nums)

	http.HandleFunc("/", handle)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func validatePricelist(str string) ([]string, []float64, bool) {

	nns := strings.Split(str, ",")

	if len(nns)%2 != 0 {
		fmt.Println("ERROR: the input string is not valid, it should be a list of pairs of product and price.")
		return nil, nil, false
	}

	names := make([]string, len(nns)/2)
	nums := make([]float64, len(nns)/2)

	for i, j := 0, 0; i < len(nns); i, j = i+2, j+1 {

		_, err := strconv.ParseFloat(nns[i], 32)
		if err == nil {
			fmt.Printf("ERROR: <%s> is not valid product name\n", nns[i])
			return nil, nil, false

		}

		v, err := strconv.ParseFloat(nns[i+1], 64)
		if err != nil {
			fmt.Printf("ERROR: <%s> is not a float number\n", nns[i+1])
			return nil, nil, false
		}

		names[j] = strings.TrimSpace(nns[i])
		nums[j] = float64(v)
	}

	return names, nums, true
}

func setPrices(names []string, nums []float64) map[string]float64 {

	res := make(map[string]float64)
	for i := 0; i < len(names); i++ {
		res[names[i]] = nums[i]
	}

	return res
}

func validateShoppingItems(str string, prices map[string]float64) ([]string, bool) {

	if len(str) == 0 {
		fmt.Printf("ERROR: the shopping list is empty\n")
		return nil, false
	}

	items := strings.Split(str, ",")

	res := make([]string, len(items))

	for i, v := range items {
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}

		if _, ok := prices[v]; !ok {
			fmt.Printf("ERROR: ***%s*** did not existed in the product list\n", v)
			return nil, false
		}
		res[i] = v
	}
	return res, true
}

func validateUserInput(str string, prices map[string]float64) (string, bool) {

	if _, ok := prices[str]; !ok {
		fmt.Printf("ERROR: ***%s*** was not a valid shopping item\n", str)
		return "", false
	}
	return str, true
}

func printNormal(ok bool, str string, v float64, w http.ResponseWriter) {

	if !ok {
		return
	}

	fmt.Fprintf(w, "%s\t\t\t\t%8.2f\n", str, v)
}

func printDiscount(ok bool, code string, v float64, w http.ResponseWriter) {
	if !ok {
		return
	}

	fmt.Fprintf(w, "\t\t%s\t\t%8.2f\n", code, v)
}

func printHead(ok bool, w http.ResponseWriter) {

	if !ok {
		return
	}

	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "Item\t\t\t\t   Price\n")
	fmt.Fprintf(w, "----\t\t\t\t   -----\n")
}

func printEnd(ok bool, v float64, w http.ResponseWriter) {

	if !ok {
		return
	}

	fmt.Fprintf(w, "----------------------------------------\n")
	fmt.Fprintf(w, "%40.2f\n", v)

}

func calculate(strs []string, prices map[string]float64, isPrint bool, w http.ResponseWriter) float64 {

	cnts_map, onOff_map := setDiscount(strs, prices)

	res := float64(0.0)

	printHead(isPrint, w)

	for _, v := range strs {

		printNormal(isPrint, v, prices[v], w)
		res += prices[v]

		discount := getDiscount(v, cnts_map, onOff_map, prices, isPrint, w)
		res += discount
	}

	printEnd(isPrint, res, w)

	return res
}

//Return two maps: cnts_map, onOff_map which will be used to calculate the discount value.
func setDiscount(products []string, prices map[string]float64) (map[string]int, map[string]bool) {

	cnts_map := make(map[string]int)
	onOff_map := make(map[string]bool)
	for k, _ := range prices {
		cnts_map[k] = 0
		onOff_map[k] = false
	}

	for _, s := range products {
		cnts_map[s] += 1
	}

	//CHMK policy
	if cnts_map["CH1"] > 0 {
		//onOff_ch1 = true
		onOff_map["CH1"] = true
	}

	return cnts_map, onOff_map
}

//Calculate the most recent discount value according to different discount policy.
func getDiscount(str string, cnts_map map[string]int, onOff_map map[string]bool, prices map[string]float64, isPrint bool, w http.ResponseWriter) float64 {

	var res float64

	switch str {
	case "CF1":
		{
			if onOff_map["CF1"] {
				printDiscount(isPrint, "BOGO", -prices["CF1"], w)
				res = -prices["CF1"]
			}
			onOff_map["CF1"] = !onOff_map["CF1"]
			return res
		}
	case "MK1":
		{
			if onOff_map["CH1"] {
				printDiscount(isPrint, "CHMK", -prices["MK1"], w)
				res = -prices["MK1"]
				onOff_map["CH1"] = false
			}
			return res
		}
	case "AP1":
		{

			//APOM policy, prefer to use APOM discount, as it will give the largest discount than APPL policy
			if cnts_map["OM1"] > 0 {
				diff := -prices["AP1"] / 2
				printDiscount(isPrint, "APOM", diff, w)
				res = diff
				cnts_map["OM1"] -= 1

				//return here, as one product only can have one discount
				return res
			}

			//APPL policy
			if cnts_map["AP1"] >= 3 {
				res = 4.50 - prices["AP1"]
				printDiscount(isPrint, "APPL", res, w)
			}
			return res
		}
	default:
		return res

	}
}
