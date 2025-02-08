package api

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"os"
	"strings"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func AddUserToSheet(spreadsheetId, userID, username string) error {
	currentTime := time.Now()
	monthName := currentTime.Month().String()
	sheetName := monthName

	creds, err := os.ReadFile("api/credentials.json")
	if err != nil {
		return fmt.Errorf("Error reading credentials: %v", err)
	}

	// Creating JWT-based config
	config, err := google.JWTConfigFromJSON(creds, sheets.SpreadsheetsScope)
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
	monthName := currentTime.Month().String()
	sheetName := monthName

	creds, err := os.ReadFile("api/credentials.json")
	if err != nil {
		return fmt.Errorf("error reading credentials: %v", err)
	}

	config, err := google.JWTConfigFromJSON(creds, sheets.SpreadsheetsScope)
	if err != nil {
		return fmt.Errorf("error loading JWT config: %v", err)
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
