package api

import (
	"PowerBook/utils"
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func getSheetname(time time.Time) string {
	year := time.Year()
	month := time.Month()
	count := 33

	switch year {
	case 2025:
		count += 12
	case 2026:
		count += 24
	case 2027:
		count += 36
	}

	switch month {
	case 1:
		count += 1
	case 2:
		count += 2
	case 3:
		count += 3
	case 4:
		count += 4
	case 5:
		count += 5
	case 6:
		count += 6
	case 7:
		count += 7
	case 8:
		count += 8
	case 9:
		count += 9
	case 10:
		count += 10
	case 11:
		count += 11
	case 12:
		count += 12
	}

	return "Круг " + strconv.Itoa(count)
}

func AddUserToSheet(spreadsheetId, userID, username string) error {
	currentTime := time.Now()
	sheetName := getSheetname(currentTime)

	utils.LoadConfig()
	creds := os.Getenv("GOOGLE_CREDENTIALS")
	if creds == "" {
		return fmt.Errorf("Error: GOOGLE_CREDENTIALS environment variable not set")
	}
	credsBytes := []byte(creds)

	// Creating JWT-based config
	config, err := google.JWTConfigFromJSON(credsBytes, sheets.SpreadsheetsScope)
	if err != nil {
		return fmt.Errorf("Error loading JWT config: %v", err)
	}

	// Creating HTTP client with JWT credentials
	client := config.Client(context.Background())

	// Creating Sheets service
	service, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("Error connecting to Sheets API: %v", err)
	}

	// Prepare data for appending
	values := [][]interface{}{
		{userID, username},
	}

	// Define the range (Sheet name and columns)
	appendRange := fmt.Sprintf("%s!A:B", sheetName)
	valueRange := &sheets.ValueRange{
		Values: values,
	}

	// Append the data to the sheet
	_, err = service.Spreadsheets.Values.Append(spreadsheetId, appendRange, valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		return fmt.Errorf("Error appending data to sheet: %v", err)
	}

	fmt.Println("User added successfully!")
	return nil
}

//////

func AddReadingMinutes(spreadsheetId, userID string, minutes int) error {
	currentTime := time.Now()
	sheetName := getSheetname(currentTime)

	utils.LoadConfig()
	creds := os.Getenv("GOOGLE_CREDENTIALS")
	if creds == "" {
		return fmt.Errorf("Error: GOOGLE_CREDENTIALS environment variable not set")
	}
	credsBytes := []byte(creds)

	// Creating JWT-based config
	config, err := google.JWTConfigFromJSON(credsBytes, sheets.SpreadsheetsScope)
	if err != nil {
		return fmt.Errorf("Error loading JWT config: %v", err)
	}

	client := config.Client(context.Background())
	service, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("error connecting to Sheets API: %v", err)
	}

	// Read the entire sheet to ensure no data is missed
	readRange := fmt.Sprintf("%s!A1:ZZ10000", sheetName)
	resp, err := service.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		return fmt.Errorf("error reading sheet data: %v", err)
	}

	// Find the current date column
	currentDate := time.Now().Format("2.01") // e.g., 06-08-2024
	dateColumnIndex := -1

	if len(resp.Values) > 0 {
		for colIdx, colValue := range resp.Values[0] { // Header row
			if strings.TrimSpace(fmt.Sprintf("%v", colValue)) == currentDate {
				dateColumnIndex = colIdx
				break
			}
		}
	}

	if dateColumnIndex == -1 {
		return fmt.Errorf("date %s not found in the sheet header", currentDate)
	}

	// Find the user row by userID in the first column
	userRowIndex := -1
	for rowIdx, row := range resp.Values {
		fmt.Printf("Row %d data: %v\n", rowIdx+1, row)

		if len(row) > 0 {
			sheetUserID := fmt.Sprintf("%v", row[0])
			fmt.Printf("Checking row %d: sheetUserID='%s' against input userID='%s'\n", rowIdx+1, sheetUserID, userID)
			if sheetUserID == strings.TrimSpace(userID) {
				userRowIndex = rowIdx
				break
			}
		}
	}

	if userRowIndex == -1 {
		return fmt.Errorf("userID %s not found in the sheet", userID)
	}

	// Convert column index to letter (e.g., 1 -> B)
	colLetter := getColumnLetter(dateColumnIndex + 1)
	cellRef := fmt.Sprintf("%s!%s%d", sheetName, colLetter, userRowIndex+1)

	// Prepare data to update
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{{minutes}},
	}

	// Update the cell
	_, err = service.Spreadsheets.Values.Update(spreadsheetId, cellRef, valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		return fmt.Errorf("error updating reading minutes: %v", err)
	}

	fmt.Printf("Reading minutes (%d) added for userID %s on %s.\n", minutes, userID, currentDate)
	return nil
}

func getColumnLetter(index int) string {
	letters := ""
	for index > 0 {
		index--
		letters = string('A'+(index%26)) + letters
		index /= 26
	}
	return letters
}
