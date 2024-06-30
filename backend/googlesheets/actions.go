package googlesheets

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var spreadsheetID string
var spreadsheetURL string

var sheetsService *sheets.Service
var driveService *drive.Service

func initGoogleServices() error {
	// Load the Google API credentials from your JSON file.
	credsPath, err := filepath.Abs("./credentials.json")
	if err != nil {
		return fmt.Errorf("unable to get absolute path of credentials file: %v", err)
	}
	creds, err := os.ReadFile(credsPath)
	if err != nil {
		return fmt.Errorf("unable to read credentials file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(creds, sheets.SpreadsheetsScope, drive.DriveFileScope)
	if err != nil {
		return fmt.Errorf("unable to create JWT config: %v", err)
	}

	client := config.Client(context.Background())

	sheetsService, err = sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to create Google Sheets service: %v", err)
	}

	driveService, err = drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to create Google Drive service: %v", err)
	}

	return nil
}

func ReadData() {
	ctx := context.Background()
	resp, err := sheetsService.Spreadsheets.Values.Get(spreadsheetID, "Лист1!A1:E12").Context(ctx).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
		return
	}

	for _, row := range resp.Values {
		fmt.Println(row)
	}
}

func getNthSheetName(n int) (string, error) {
	ctx := context.Background()
	spreadsheet, err := sheetsService.Spreadsheets.Get(spreadsheetID).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve spreadsheet metadata: %v", err)
	}

	if len(spreadsheet.Sheets) == 0 {
		return "", fmt.Errorf("no sheets found in the spreadsheet")
	}

	if n == 0 || n > len(spreadsheet.Sheets) {
		n = 1
	}
	nthSheetName := spreadsheet.Sheets[n-1].Properties.Title
	return nthSheetName, nil
}

func WriteHeaders(headers []interface{}) {
	ctx := context.Background()
	var vr sheets.ValueRange
	vr.Values = append(vr.Values, headers)
	sheetName, err := getNthSheetName(1)
	if err != nil {
		log.Println(err)
	}

	cellFrom := "A1"
	cellTo := string(rune('A'+len(headers))) + "1"

	_, err = sheetsService.Spreadsheets.Values.Update(spreadsheetID,
		fmt.Sprintf("%s!%s:%s", sheetName, cellFrom, cellTo), &vr,
	).ValueInputOption("RAW").Context(ctx).Do()
	if err != nil {
		log.Fatalf("Unable to write headers to sheet: %v", err)
	}
	fmt.Println("Headers written successfully!")
}

func WriteData(data [][]interface{}) {
	ctx := context.Background()
	var vr sheets.ValueRange

	for _, row := range data {
		vr.Values = append(vr.Values, row)
	}

	sheetName, err := getNthSheetName(1)
	if err != nil {
		log.Println(err)
	}

	letterFrom, numberFrom := "A", 2
	letterTo, numberTo := string(rune('A'+len(data[0]))), len(data)+numberFrom-1

	_, err = sheetsService.Spreadsheets.Values.Update(spreadsheetID,
		fmt.Sprintf("%s!%s%d:%s%d", sheetName, letterFrom, numberFrom, letterTo, numberTo), &vr,
	).ValueInputOption("RAW").Context(ctx).Do()

	if err != nil {
		log.Fatalf("Unable to write data to sheet: %v", err)
	}
	fmt.Println("Data written successfully!")
}

func CreateSpreadsheet(title string) (string, error) {
	if sheetsService == nil || driveService == nil {
		err := initGoogleServices()
		if err != nil {
			return "", err
		}
	}

	if strings.TrimSpace(title) == "" {
		title = "New sheet"
	}

	ctx := context.Background()
	spreadsheet := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: title,
		},
	}

	resp, err := sheetsService.Spreadsheets.Create(spreadsheet).Context(ctx).Do()
	if err != nil {
		log.Printf("Unable to create spreadsheet: %v", err)
		return "", err
	}
	fmt.Printf("Spreadsheet created: %s\n", resp.SpreadsheetUrl)

	shareSpreadsheet(resp.SpreadsheetId)

	spreadsheetID = resp.SpreadsheetId
	spreadsheetURL = resp.SpreadsheetUrl

	return resp.SpreadsheetId, nil
}

func shareSpreadsheet(spreadsheetId string) {
	ctx := context.Background()
	permission := &drive.Permission{
		Type: "anyone",
		Role: "writer", // Replace with your service account email
	}

	_, err := driveService.Permissions.Create(spreadsheetId, permission).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Unable to share spreadsheet: %v", err)
		return
	}
	fmt.Println("Spreadsheet shared successfully!")
}

func ToInterfaceSlice[T any](src []T) []interface{} {
	interfaceSlice := make([]interface{}, len(src))
	for i, h := range src {
		interfaceSlice[i] = h
	}
	return interfaceSlice
}

func GetSpreadsheetURL() string {
	return spreadsheetURL
}
