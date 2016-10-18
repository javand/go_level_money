package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/shopspring/decimal" //makes dealing with money a bit easier
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	flag.StringVar(&Configuration.UID, "uid", Configuration.UID, "Level User Id")
	flag.StringVar(&Configuration.AuthToken, "auth_token", Configuration.AuthToken, "Leven Auth Token")
	flag.StringVar(&Configuration.APIToken, "api_token", Configuration.APIToken, "Level API Token")
	flag.StringVar(&Configuration.ConfigFile, "config_file", Configuration.ConfigFile, "location of your configuration file")
	flag.BoolVar(&Configuration.HelpRequested, "help", Configuration.HelpRequested, "display help information")
	flag.StringVar(&Configuration.LogDir, "log_dir", Configuration.LogDir, "location logs should be stored")
	flag.BoolVar(&Configuration.IgnoreDonuts, "ignore-donuts", Configuration.IgnoreDonuts, "oooooh Donuts")

	LoadConfigurationFile()

	ParseCLIArgs()

	file, err := os.OpenFile(fmt.Sprintf("%s/%s.log", Configuration.LogDir, "go_level_money"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", ":", err)
	}

	log.SetOutput(file)
}

func main() {

	fmt.Println("Level Money For Life!!\n")

	at := GetAllTransactions()
	BuildMonthlyTransactionReport(at)
}

func GetAllTransactions() []Transaction {

	uid64, err := strconv.ParseInt(Configuration.UID, 10, 64)

	type Args struct {
		Uid      int64  `json:"uid"`
		Token    string `json:"token"`
		ApiToken string `json:"api-token"`
	}

	type ArgHolder struct {
		Val Args `json:"args"`
	}

	a := Args{Uid: uid64, Token: Configuration.AuthToken, ApiToken: Configuration.APIToken}
	ah := ArgHolder{Val: a}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(ah)
	req, err := http.NewRequest("POST", ConfigurationFileMap.BaseURL+"/get-all-transactions", b)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var ts TransactionsResponse
	json.NewDecoder(resp.Body).Decode(&ts)
	return ts.TransactionList
}

func IsADonut(possibleDonut string) (newval bool) {
	if strings.Contains(possibleDonut, "Donuts") {
		newval = true
	} else if strings.Contains(possibleDonut, "DUNKIN") {
		newval = true
	}
	return
}

func BuildMonthlyTransactionReport(transactions []Transaction) {

	report := make(map[string]MonthlyReport)
	layout := "2006-01-02T15:04:05.000Z"
	for _, t := range transactions {
		if Configuration.IgnoreDonuts {
			if IsADonut(t.Merchant) {
				continue
			}
		}
		t1, _ := time.Parse(layout, t.TransactionTime)
		mr, ok := report[t1.Format("2006-01")]
		if !ok {
			if t.Amount < 0 {
				t.Amount = -t.Amount
				mr = MonthlyReport{Spent: decimal.New(t.Amount, 0), Income: decimal.New(0, 0)}
			} else {
				mr = MonthlyReport{Spent: decimal.New(0, 0), Income: decimal.New(t.Amount, 0)}
			}
			report[t1.Format("2006-01")] = mr
		} else {
			if t.Amount < 0 {
				t.Amount = -t.Amount
				mr.Spent = mr.Spent.Add(decimal.New(t.Amount, 0))
				report[t1.Format("2006-01")] = mr
			} else {
				mr.Income = mr.Income.Add(decimal.New(t.Amount, 0))
				report[t1.Format("2006-01")] = mr
			}
		}
	}

	var divider = decimal.New(10000, 0)
	var totalspend = decimal.New(0, 0)
	var totalincome = decimal.New(0, 0)

	for k, v := range report {
		totalspend = totalspend.Add(v.Spent)
		totalincome = totalincome.Add(v.Income)
		v.Spent = v.Spent.Div(divider)
		v.Income = v.Income.Div(divider)
		report[k] = v
	}

	b, err := json.Marshal(report)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

	n := strconv.Itoa(len(report))
	totalMonthsReported, err := decimal.NewFromString(n)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(report) > 0 {
		asr := MonthlyReport{totalspend.Div(totalMonthsReported).Div(divider).Round(2), totalincome.Div(totalMonthsReported).Div(divider).Round(2)}

		b1, err := json.Marshal(asr)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("average:", string(b1))
	}
}
