package googlesheets

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/logger"
	"go.uber.org/zap"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var spreadsheets sync.Map

var sheetsService *sheets.Service
var driveService *drive.Service

const SheetAccessType = "anyone"

var Error = errors.New("unable to create and edit google sheet, please, try later")
var CsvError = errors.New("unable to convert google sheet to csv")

type Spreadsheet struct {
	OwnerEmail string
	Obj        *sheets.Spreadsheet
}

//func getSpreadsheet(email string) *Spreadsheet {
//	client, ok := spreadsheets.Load(email)
//	if !ok {
//		return nil
//	}
//	return client.(*Spreadsheet)
//}

func setSpreadsheet(email string, client *Spreadsheet) {
	spreadsheets.Store(email, client)
}

func initGoogleServices() error {
	// Load the Google API credentials from your JSON file.
	credsPath, err := filepath.Abs("./credentials.json")
	if err != nil {
		logger.Error(fmt.Errorf("unable to get absolute path of credentials file: %w", err))
		return Error
	}
	creds, err := os.ReadFile(credsPath)
	if err != nil {
		logger.Error(fmt.Errorf("unable to read credentials file: %w", err))
		return Error
	}

	config, err := google.JWTConfigFromJSON(creds, sheets.SpreadsheetsScope, drive.DriveFileScope)
	if err != nil {
		logger.Error(fmt.Errorf("unable to create JWT config: %w", err))
		return Error
	}

	client := config.Client(context.Background())

	sheetsService, err = sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		logger.Error(fmt.Errorf("unable to create Google Sheets service: %w", err))
		return Error
	}

	driveService, err = drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		logger.Error(fmt.Errorf("unable to create Google Drive service: %w", err))
		return Error
	}

	return nil
}

func (s *Spreadsheet) getNthSheetName(n int) (string, error) {
	ctx := context.Background()
	spreadsheet, err := sheetsService.Spreadsheets.Get(s.Obj.SpreadsheetId).Context(ctx).Do()
	if err != nil {
		logger.Error(fmt.Errorf("unable to retrieve spreadsheet metadata: %w", err))
		return "", Error
	}

	if len(spreadsheet.Sheets) == 0 {
		logger.Error(fmt.Errorf("no sheets found in the spreadsheet: %w", err))
		return "", Error
	}

	if n == 0 || n > len(spreadsheet.Sheets) {
		n = 1
	}
	nthSheetName := spreadsheet.Sheets[n-1].Properties.Title
	return nthSheetName, nil
}

func (s *Spreadsheet) WriteHeaders(headers []interface{}) error {
	ctx := context.Background()

	var vr sheets.ValueRange

	vr.Values = append(vr.Values, headers)
	sheetName, err := s.getNthSheetName(1)
	if err != nil {
		return err
	}

	cellFrom := "A1"
	cellTo := string(rune('A'+len(headers))) + "1"

	_, err = sheetsService.Spreadsheets.Values.Update(s.Obj.SpreadsheetId,
		fmt.Sprintf("%s!%s:%s", sheetName, cellFrom, cellTo), &vr,
	).ValueInputOption("RAW").Context(ctx).Do()
	if err != nil {
		logger.Error(fmt.Errorf("unable to write headers to sheet: %w", err))
		return Error
	}

	logger.Logger().Info("Headers written successfully!",
		zap.String("spreadsheet ID", s.Obj.SpreadsheetId),
		zap.String("owner email", s.OwnerEmail))

	return nil
}

func (s *Spreadsheet) WriteData(data [][]interface{}) error {
	ctx := context.Background()
	var vr sheets.ValueRange

	for _, row := range data {
		vr.Values = append(vr.Values, row)
	}

	sheetName, err := s.getNthSheetName(1)
	if err != nil {
		return err
	}

	letterFrom, numberFrom := "A", 2
	letterTo, numberTo := string(rune('A'+len(data[0]))), len(data)+numberFrom-1

	_, err = sheetsService.Spreadsheets.Values.Update(s.Obj.SpreadsheetId,
		fmt.Sprintf("%s!%s%d:%s%d", sheetName, letterFrom, numberFrom, letterTo, numberTo), &vr,
	).ValueInputOption("USER_ENTERED").Context(ctx).Do()
	if err != nil {
		logger.Error(fmt.Errorf("unable to write data to sheet: %w", err))
		return Error
	}

	logger.Logger().Info("Data written successfully!",
		zap.String("spreadsheet ID", s.Obj.SpreadsheetId),
		zap.String("owner email", s.OwnerEmail))

	numColumns := len(vr.Values)

	req := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				AutoResizeDimensions: &sheets.AutoResizeDimensionsRequest{
					Dimensions: &sheets.DimensionRange{
						SheetId:    0,
						Dimension:  "COLUMNS",
						StartIndex: 0,
						EndIndex:   int64(numColumns),
					},
				},
			},
		},
	}

	_, err = sheetsService.Spreadsheets.BatchUpdate(s.Obj.SpreadsheetId, req).Context(ctx).Do()
	if err != nil {
		logger.Error(fmt.Errorf("unable to resize columns: %w", err))
		return nil
	}

	logger.Logger().Info("Columns resized successfully!",
		zap.String("spreadsheet ID", s.Obj.SpreadsheetId),
		zap.String("owner email", s.OwnerEmail))

	return nil
}

func CreateSpreadsheet(title, ownerEmail string) (*Spreadsheet, error) {
	if sheetsService == nil || driveService == nil {
		err := initGoogleServices()
		if err != nil {
			return nil, err
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
		logger.Error(fmt.Errorf("unable to create spreadsheet: %w", err))
		return nil, Error
	}
	logger.Logger().Info("Spreadsheet created: "+resp.SpreadsheetUrl,
		zap.String("spreadsheet ID", resp.SpreadsheetId),
		zap.String("owner email", ownerEmail))

	spreadsheetObj := &Spreadsheet{
		OwnerEmail: ownerEmail,
		Obj:        resp,
	}

	setSpreadsheet(ownerEmail, spreadsheetObj)

	err = spreadsheetObj.shareSpreadsheet(resp.SpreadsheetId)
	if err != nil {
		return nil, err
	}

	return spreadsheetObj, nil
}

func (s *Spreadsheet) shareSpreadsheet(spreadsheetId string) error {
	ctx := context.Background()
	permission := &drive.Permission{
		Type: SheetAccessType,
		Role: "writer",
		//EmailAddress: s.OwnerEmail,
	}

	_, err := driveService.Permissions.Create(spreadsheetId, permission).Context(ctx).Do()
	if err != nil {
		var gErr *googleapi.Error
		if errors.As(err, &gErr) {
			for _, e := range gErr.Errors {
				if e.Reason == "invalid" && strings.Contains(e.Message, "email") {
					logger.Error(fmt.Errorf("invalid email address provided for sharing: %w", err))
					return errors.New("invalid email address provided for sharing: " + s.OwnerEmail)
				}
			}
		}
		logger.Error(fmt.Errorf("unable to share spreadsheet: %w", err))
		return Error
	}

	logger.Logger().Info("Spreadsheet shared successfully!",
		zap.String("spreadsheet ID", s.Obj.SpreadsheetId),
		zap.String("owner email", s.OwnerEmail))

	return nil
}

func (s *Spreadsheet) GetSpreadsheetURL() string {
	return s.Obj.SpreadsheetUrl
}

func (s *Spreadsheet) GetSpreadsheetCSV() ([]byte, error) {
	resp, err := driveService.Files.Export(s.Obj.SpreadsheetId, "text/csv").Download()
	if err != nil {
		logger.Error(fmt.Errorf("unable to export file: %w", err))
		return nil, CsvError
	}
	defer resp.Body.Close()

	buff := new(bytes.Buffer)

	_, err = io.Copy(buff, resp.Body)
	if err != nil {
		logger.Error(fmt.Errorf("unable to write to buffer: %w", err))
		return nil, CsvError
	}

	return buff.Bytes(), nil

}

func ToInterfaceSlice[T any](src []T) []interface{} {
	interfaceSlice := make([]interface{}, len(src))
	for i, h := range src {
		interfaceSlice[i] = h
	}
	return interfaceSlice
}
