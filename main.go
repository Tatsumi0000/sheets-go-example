package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

func main() {
	secret, err := ioutil.ReadFile("secret.json")
	if err != nil {
		log.Fatal(err)
	}

	conf, err := google.JWTConfigFromJSON(secret, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(context.Background())
	srv, err := sheets.New(client)
	if err != nil {
		log.Fatal(err)
	}

	spreadsheetID := os.Getenv("SHEET_ID")

	// 参照
	readRange := "Sheet1!A:B"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
		return
	}
	for _, row := range resp.Values {
		fmt.Printf("%s, %s\n", row[0], row[1])
	}

	// 更新
	var vr sheets.ValueRange
	now := time.Now().Format("2006/01/02 15:04:05")
	for i := 0; i < len(resp.Values); i++ {
		vr.Values = append(vr.Values, []interface{}{now})
	}
	updateRange := "Sheet1!B1:B"
	if _, err = srv.Spreadsheets.Values.Update(spreadsheetID, updateRange, &vr).ValueInputOption("RAW").Do(); err != nil {
		log.Fatal(err)
	}
}
