package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"flag"

	_ "github.com/joho/godotenv/autoload"
	"github.com/swexbe/zettleIT/api"
)

const timeFormat = "2006-01-02"

// IO
func getDateInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Payout Date: ")
	date, err := reader.ReadString('\n')

	if err != nil {
		return "", errors.New("Unable to get input")
	}

	return strings.TrimSuffix(date, "\n"), nil
}

func getDateString(date string) (string, string, error) {

	endDate, err := time.Parse(timeFormat, date)

	if err != nil {
		return "", "", fmt.Errorf("Invalid Date: %s", date)
	}

	startDate := endDate.AddDate(0, 0, -14)

	sdString := startDate.Format(timeFormat)
	edString := endDate.Format(timeFormat)

	return sdString, edString, nil

}

func handleFatalErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// mapUUIDToSeller maps the
func mapUUIDToSeller(purchases []api.Purchase, verbose bool) map[string]string {

	purchasesMap := make(map[string]string)

	for _, v := range purchases {
		if verbose {
			fmt.Printf("Purchase UUID: %s mapped to seller: %s\n", v.Payments[0].UUID, v.UserDisplayName)
		}
		purchasesMap[v.Payments[0].UUID] = v.UserDisplayName
	}

	return purchasesMap

}

func countTransactions(transactions []api.Transaction, purchasesMap map[string]string, verbose bool) (map[string]int, int) {

	amountSold := make(map[string]int)

	numPayouts := 0
	total := 0

	for _, v := range transactions {

		if v.Type == "PAYOUT" {
			numPayouts++
		}

		if numPayouts == 0 {
			continue
		}

		if numPayouts >= 2 {
			break
		}

		seller := purchasesMap[v.UUID]
		if seller == "" {
			total += v.Amount
			continue
		}
		if verbose {
			fmt.Printf("Seller %s sold item for %d öre\n", seller, v.Amount)
		}
		amountSold[seller] += v.Amount

	}

	return amountSold, -total
}

func printDistribution(amountSold map[string]int, total int) {
	fmt.Print("\nDISTRIBUTION OF PAYMENTS: \n")

	fmt.Printf("Total : %.2f kr \n", float64(total)/100)

	for key, value := range amountSold {

		amount := float64(value) / 100

		fmt.Printf("%s : %.2f kr \n", key, amount)
	}
}

func main() {

	verbose := flag.Bool("v", false, "")
	dateP := flag.String("d", "", "")
	flag.Parse()

	if *verbose {
		fmt.Printf("***********VERBOSE MODE***********\n")
		fmt.Printf("Date input is: %s\n", *dateP)
	}

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	token, err := api.GetAuthkey(username, password)

	var date string

	if *dateP != "" {
		date = *dateP
	} else {
		var err error
		date, err = getDateInput()
		handleFatalErr(err)
	}

	startDate, endDate, err := getDateString(date)
	handleFatalErr(err)

	if *verbose {
		fmt.Printf("StartDate: %s EndDate: %s\n", startDate, endDate)
	}

	transactions, err := api.GetTransactions(startDate, endDate, token)
	handleFatalErr(err)
	purchases, err := api.GetPurchases(startDate, endDate, token)
	handleFatalErr(err)

	if *verbose {
		fmt.Printf("Number of transaction %d \nNumber of purchases %d\n", len(transactions), len(purchases))
	}

	purchasesMap := mapUUIDToSeller(purchases, *verbose)
	amountSold, total := countTransactions(transactions, purchasesMap, *verbose)

	printDistribution(amountSold, total)
}
