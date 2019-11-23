package main

import (
	"log"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/swexbe/zettleIT/api"
)

const timeFormat = "2006-01-02"

func main() {

	token := api.GetAuthkey()

	d := "2019-11-17"

	endDate, err := time.Parse(timeFormat, d)
	if err != nil {
		log.Fatalln(err)
	}

	startDate := endDate.AddDate(0, 0, -14)

	transactions := api.GetTransactions(startDate.Format(timeFormat), endDate.Format(timeFormat), token)

	var transactionsMap map[string]api.Transaction

	for _, v := range transactions {
		transactionsMap[v.UUID] = v
	}
}
