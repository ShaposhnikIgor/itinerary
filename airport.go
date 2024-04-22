package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// AirportLookup represents a data structure for airport lookup.
type AirportLookup struct {
	Name         string // Airport name
	IsoCountry   string // ISO country code
	Municipality string // Municipality/City
	IcaoCode     string // ICAO code
	IataCode     string // IATA code
	Coordinates  string // Coordinates
}

// TextFormattingSettings represents the formatting settings for text.
type TextFormattingSettings struct {
	Color         string // ANSI color code
	Bold          bool   // Bold text indicator
	Italic        bool   // Italic text indicator
	Underline     bool   // Underline text indicator
	Strikethrough bool   // Strikethrough text indicator
}

func main() {
	// Command-line arguments parsing
	showHelp := flag.Bool("h", false, "Display usage information")
	toOutput := flag.Bool("o", false, "Write output to file")

	flag.Usage = func() {
		fmt.Println("itinerary usage:")
		fmt.Printf("go run . ./input.txt ./output.txt ./airport-lookup.csv\n")
	}
	flag.Parse()

	// Checking for the correct number of arguments
	if (len(os.Args) != 4 && !*toOutput) || (*showHelp && !*toOutput) {
		flag.Usage()
		os.Exit(0)
	}

	inputPath := ""
	outputPath := ""
	lookupPath := ""
	settingsPath := "./user_settings.txt"

	if *toOutput {
		inputPath = "input.txt"
		lookupPath = "airport-lookup.csv"
	} else {
		inputPath = os.Args[1]
		outputPath = os.Args[2]
		lookupPath = os.Args[3]
	}

	// Start the main processing
	startGeneretion()

	// Decoding the expected digital signature from hex string to byte slice
	expectedSignatureHex := "d5c5c7c0aaebd7c8975aad8be7b192ab96fd859708060610006e1a0347504f53"
	expectedSignature, err := hex.DecodeString(expectedSignatureHex)
	if err != nil {
		log.Fatal("Failed to decode expected signature:", err)
	}

	// Specify the full path to the signature.txt file
	signatureFilePath := "signature.txt"
	// Reading the digital signature from the file
	signatureBytes, err := os.ReadFile(signatureFilePath)
	if err != nil {
		fmt.Println("Failed to read signature file:", err)
		os.Exit(0)
	}
	// Converting the read digital signature from string to byte slice
	signatureData, err := hex.DecodeString(string(signatureBytes))
	if err != nil {
		fmt.Println("Failed to decode signature data:", err)
		os.Exit(0)
	}

	if bytes.Equal(signatureData, expectedSignature) {
		fmt.Println("Digital signatures match!")

		// Removing the signature.txt file after key comparison
		err := os.Remove(signatureFilePath)
		if err != nil {
			fmt.Println("Failed to remove signature file:", err)
		} else {
			// Executing the processing pipeline if keys match
			executeProcessingPipeline(inputPath, outputPath, lookupPath, settingsPath, *toOutput)
		}
	} else {
		fmt.Println("Digital signatures do not match!")
		// Removing the signature.txt file in case of mismatched keys
		err := os.Remove(signatureFilePath)
		if err != nil {
			fmt.Println("Failed to remove signature file:", err)
		}
		// Exiting the program
		os.Exit(0)
	}
}

// executeProcessingPipeline runs the processing pipeline.
func executeProcessingPipeline(inputPath, outputPath, lookupPath, settingsPath string, toOutput bool) {
	// Create a log file
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer logFile.Close()
	// Set log output to file
	log.SetOutput(logFile)

	// Check for the existence of the settings file
	settingsFile, err := os.Open(settingsPath)
	if err != nil {
		fmt.Println("Settings not found")
		_, file, line, _ := runtime.Caller(0)
		log.Printf(file, line, "Settings not found")
		os.Exit(0)
	}
	defer settingsFile.Close()

	// Check for the existence of the input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		fmt.Println("Input not found")
		_, file, line, _ := runtime.Caller(0)
		log.Printf(file, line, "Input not found")
		os.Exit(0)
	}
	defer inputFile.Close()

	// Check for the existence of the airport lookup file
	lookupFile, err := os.Open(lookupPath)
	if err != nil {
		fmt.Println("Airport lookup not found")
		_, file, line, _ := runtime.Caller(0)
		log.Printf(file, line, "Airport lookup not found")
		os.Exit(0)
	}
	defer lookupFile.Close()

	// Read airport data from CSV
	airportLookup, err := readAirportLookup(lookupFile)
	if err != nil {
		fmt.Println("Airport lookup malformed")
		_, file, line, _ := runtime.Caller(0)
		log.Printf(file, line, "Airport lookup malformed")
		os.Exit(0)
	}

	if toOutput {
		stdOutputText := processInput(inputFile, airportLookup, true)
		fmt.Println(stdOutputText)
	} else {
		outputText := processInput(inputFile, airportLookup, false)
		// Write processed text to output file
		outputFile, err := os.Create(outputPath)
		if err != nil {
			fmt.Println("Error creating output file:", err)
			_, file, line, _ := runtime.Caller(0)
			log.Printf(file, line, "Error creating output file:", err)
			os.Exit(0)
		}
		defer outputFile.Close()

		outputWriter := bufio.NewWriter(outputFile)
		_, err = outputWriter.WriteString(outputText)
		if err != nil {
			fmt.Println("Error writing to output file:", err)
			_, file, line, _ := runtime.Caller(0)
			log.Printf(file, line, "Error writing to output file:", err)
			os.Exit(0)
		}

		outputWriter.Flush()

		fmt.Println("Processing complete.")
	}
}

// readAirportLookup reads airport data from CSV and returns it as a slice of AirportLookup structures.
func readAirportLookup(file io.Reader) ([]AirportLookup, error) {
	reader := csv.NewReader(file)

	// Skip header row
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	// Expected headers
	expectedHeaders := []string{"name", "iso_country", "municipality", "icao_code", "iata_code", "coordinates"}
	var mismatchedHeaders []string

	// Check for all expected headers and determine their indexes
	headerIndices := make(map[string]int)
	for i, col := range header {
		headerIndices[col] = i
	}
	for _, expectedHeader := range expectedHeaders {
		_, ok := headerIndices[expectedHeader]
		if !ok {
			mismatchedHeaders = append(mismatchedHeaders, expectedHeader) // Add mismatched headers to the slice
		}
	}

	if len(mismatchedHeaders) > 0 {
		// Combine all mismatched headers into a single string
		mismatchedHeaderStr := strings.Join(mismatchedHeaders, ", ")
		_, file, line, _ := runtime.Caller(0)
		log.Printf("%s:%d: airport lookup malformed: missing or incorrect headers: %s", file, line, mismatchedHeaderStr)
		return nil, fmt.Errorf("airport lookup malformed: missing or incorrect headers: '%s'", mismatchedHeaderStr)
	}

	// Define column indexes
	nameIdx := headerIndices["name"]
	isoCountryIdx := headerIndices["iso_country"]
	municipalityIdx := headerIndices["municipality"]
	icaoCodeIdx := headerIndices["icao_code"]
	iataCodeIdx := headerIndices["iata_code"]
	coordinatesIdx := headerIndices["coordinates"]

	var airportLookup []AirportLookup

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Check for data in specified indexes
		if record[nameIdx] == "" || record[isoCountryIdx] == "" || record[municipalityIdx] == "" || record[icaoCodeIdx] == "" || record[iataCodeIdx] == "" || record[coordinatesIdx] == "" {
			_, file, line, _ := runtime.Caller(0)
			log.Printf("%s:%d: airport lookup malformed: missing data in one or more columns", file, line)
			return nil, fmt.Errorf("airport lookup malformed: missing data in one or more columns")
		}

		// Check for all fields in the record
		if len(record) != len(header) {
			_, file, line, _ := runtime.Caller(0)
			log.Printf("%s:%d: airport lookup malformed: missing or incomplete fields", file, line)
			return nil, fmt.Errorf("airport lookup malformed: missing or incomplete fields")
		}

		airport := AirportLookup{
			Name:         record[nameIdx],
			IsoCountry:   record[isoCountryIdx],
			Municipality: record[municipalityIdx],
			IcaoCode:     record[icaoCodeIdx],
			IataCode:     record[iataCodeIdx],
			Coordinates:  record[coordinatesIdx],
		}

		airportLookup = append(airportLookup, airport)
	}

	return airportLookup, nil
}

// parseSettingsFromFile reads settings from a file and returns a map of TextFormattingSettings with their respective identifiers
func parseSettingsFromFile(filename string) (map[string]TextFormattingSettings, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	settingsMap := make(map[string]TextFormattingSettings)
	currentIdentifier := ""

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentIdentifier = strings.TrimSuffix(strings.TrimPrefix(line, "["), "]")
		} else {
			parts := strings.Split(line, "=")
			if len(parts) != 2 {
				continue // Skip invalid lines
			}
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			settings := settingsMap[currentIdentifier]
			switch key {
			case "Color":
				settings.Color = value
			case "Bold":
				settings.Bold = (value == "true")
			case "Italic":
				settings.Italic = (value == "true")
			case "Underline":
				settings.Underline = (value == "true")
			case "Strikethrough":
				settings.Strikethrough = (value == "true")
			}
			settingsMap[currentIdentifier] = settings
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return settingsMap, nil
}

// processInput обрабатывает текст из входного файла, преобразуя аэропортные коды и даты/время.
func processInput(input io.Reader, airportLookup []AirportLookup, toStdOutput bool) string {
	var outputText strings.Builder
	scanner := bufio.NewScanner(input)

	// Initialize variables to track consecutive blank lines
	consecutiveBlankLines := 0

	// Get caller information
	_, file, line, _ := runtime.Caller(0)

	for scanner.Scan() {
		line := scanner.Text()

		// Преобразование аэропортных кодов
		line = convertAirportCodes(line, airportLookup, toStdOutput)

		// Преобразование ISO дат и времени
		line = convertISODateTime(line, toStdOutput)

		// Trim leading and trailing whitespace
		line = strings.TrimSpace(line)

		// Replace line-break characters (\v, \f, \r) and their text representations with \n
		line = strings.ReplaceAll(line, "\v", "\n")
		line = strings.ReplaceAll(line, "\\v", "\n")
		line = strings.ReplaceAll(line, "\f", "\n")
		line = strings.ReplaceAll(line, "\\f", "\n")
		line = strings.ReplaceAll(line, "\r", "\n")
		line = strings.ReplaceAll(line, "\\r", "\n")

		// Check for consecutive blank lines
		if line == "" {
			consecutiveBlankLines++
			if consecutiveBlankLines > 1 {
				continue // Пропустить запись, если более одной последовательной пустой строки
			}
		} else {
			consecutiveBlankLines = 0 // Сбросить счетчик последовательных пустых строк
		}

		outputText.WriteString(cleanNonASCIIChars(line) + "\n")
	}

	// Получаем итоговый текст после обработки ввода
	processedInput := outputText.String()

	// Log the completion of processing
	log.Printf("%s:%d: processing of input completed", file, line)

	// Сокращаем пустые строки до одной
	finalOutput := reduceEmptyLines(processedInput)

	return finalOutput
}

func reduceEmptyLines(input string) string {
	// Разбиваем текст на строки
	lines := strings.Split(input, "\n")

	// Инициализируем переменную для хранения результата
	var result strings.Builder

	// Инициализируем флаг, который показывает, что предыдущая строка была пустой
	prevEmpty := false

	// Проходим по каждой строке
	for _, line := range lines {
		// Если строка не пустая, записываем ее
		if line != "" {
			// Если предыдущая строка была пустой, а текущая не пустая, добавляем пустую строку
			if prevEmpty {
				result.WriteString("\n")
			}
			result.WriteString(line + "\n")
			// Сбрасываем флаг пустой строки
			prevEmpty = false
		} else {
			// Если текущая строка пустая, устанавливаем флаг
			prevEmpty = true
		}
	}

	if prevEmpty && len(lines) > 1 {
		resultString := result.String()
		resultString = strings.TrimSuffix(resultString, "\n")
		return resultString
	}

	return result.String()
}

// cleanNonASCIIChars removes non-ASCII characters from a string
func cleanNonASCIIChars(line string) string {
	var cleanedLine strings.Builder
	for _, char := range line {
		if char < 128 { // Check if the character is within ASCII range
			cleanedLine.WriteRune(char)
		}
	}
	return cleanedLine.String()
}

// convertAirportCodes replaces airport codes with their names in the text string.
func convertAirportCodes(text string, airportLookup []AirportLookup, toStdOutput bool) string {
	uniqueAirErrors := make(map[string]bool)
	var errAirList []string

	settingsPath := "./user_settings.txt"

	// Make a settings map from user settings
	settingsMap, err := parseSettingsFromFile(settingsPath)
	if err != nil {
		fmt.Println("Error:", err)
	}

	re := regexp.MustCompile(`#([A-Z]{3})|##([A-Z]{4})|\*\#([A-Z]{3})|\*\##([A-Z]{4})`)
	processedText := re.ReplaceAllStringFunc(text, func(code string) string {

		for _, airport := range airportLookup {
			if code == "#"+airport.IataCode || code == "##"+airport.IcaoCode {
				if toStdOutput {
					return formatText(airport.Name, settingsMap["Airport"]) // Format the text for standard output
				} else {
					return airport.Name // Return the airport name as it is
				}
			}
			if code == "*#"+airport.IataCode || code == "*##"+airport.IcaoCode {
				if toStdOutput {
					return formatText(airport.Municipality, settingsMap["City"]) // Format the text for standard output
				} else {
					return airport.Municipality // Return the airport name as it is
				}
			}
		}

		errStr := code
		if !uniqueAirErrors[errStr] {
			uniqueAirErrors[errStr] = true
			errAirList = append(errAirList, errStr)
		}
		return code // Return the airport code unchanged if not found
	})
	if len(errAirList) > 0 {
		_, file, line, _ := runtime.Caller(0)
		errAirStrings := strings.Join(errAirList, ", ")
		log.Printf("%s:%d: airport code not found: %s", file, line, errAirStrings)
	}
	return processedText
}

// convertISODateTime converts ISO date and time strings in the text to specified format.
func convertISODateTime(line string, toStdOutput bool) string {
	var errStringsList []string
	uniqueErrors := make(map[string]bool)

	re := regexp.MustCompile(`(D|T12|T24)\(([^)]+)\)`)
	matches := re.FindAllStringSubmatch(line, -1)

	settingsPath := "./user_settings.txt"

	// Make a settings map from user settings
	settingsMap, err := parseSettingsFromFile(settingsPath)
	if err != nil {
		fmt.Println("Error:", err)
	}

	for _, match := range matches {
		prefix := match[1]
		timeStr := match[2]

		// Replace minus symbol (−) with regular hyphen (-)
		timeStr = strings.ReplaceAll(timeStr, "−", "-")

		t, err := time.Parse("2006-01-02T15:04Z07:00", timeStr)
		if err != nil {
			errStr := timeStr + ": " + err.Error()
			if !uniqueErrors[errStr] {
				uniqueErrors[errStr] = true
				errStringsList = append(errStringsList, errStr)
			}
			continue
		}

		var formattedTime string
		var offset string
		var formattedTimeStd string

		if t.Location().String() == "UTC" {
			offset = "(-07:00)" // Set the offset according to your timezone
		} else {
			offset = t.Format("(-07:00)")
		}

		formatedOfset := formatText(offset, settingsMap["OffsetPos"])
		if strings.HasPrefix(offset, "(-") {
			formatedOfset = formatText(offset, settingsMap["OffsetNeg"])
		}

		switch prefix {
		case "D":
			formattedTime = t.Format("02 Jan 2006")
			formattedTimeStd = formatText(formattedTime, settingsMap["Date"])
		case "T12":
			formattedTime = t.Format("03:04PM ") + offset
			formattedTimeStd = formatText(formattedTime, settingsMap["Time"]) + formatedOfset
		case "T24":
			formattedTime = t.Format("15:04 ") + offset
			formattedTimeStd = formatText(formattedTime, settingsMap["Time"]) + formatedOfset
		}

		if toStdOutput {
			line = strings.Replace(line, match[0], formattedTimeStd, 1)
		} else {
			line = strings.Replace(line, match[0], formattedTime, 1)
		}

	}

	if len(errStringsList) > 0 {
		errStrings := strings.Join(errStringsList, ", ")
		_, file, lineNum, _ := runtime.Caller(0)
		log.Printf("%s:%d: failed to parse time string: %s", file, lineNum, errStrings)
	}

	return line
}

// formatText formats a line of text based on provided settings
func formatText(line string, settings TextFormattingSettings) string {
	formatString := ""

	if settings.Color != "" {
		formatString += "\033[" + settings.Color + "m"
	}

	if settings.Bold {
		formatString += "\033[1m"
	}

	if settings.Italic {
		formatString += "\033[3m"
	}

	if settings.Underline {
		formatString += "\033[4m"
	}

	if settings.Strikethrough {
		formatString += "\033[9m"
	}

	// Adding reset ANSI code at the end to reset all formatting
	formatString += line + "\033[0m"

	return formatString
}
