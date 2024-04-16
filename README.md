# Usage

## Installation
To use this program, ensure you have Go installed on your system. If not, 
you can download and install it from the official Go website.


## Running the Program
To run the program, execute the following command in your terminal:

``` 
go run . ./input.txt ./output.txt ./airport-lookup.csv
```
Replace input.txt, output.txt, and airport-lookup.csv with the paths to your input text file, output text file,
 and airport lookup CSV file respectively.

## Command-line Options
The program supports the following command-line options:


-h: Display usage information.
```
go run . -h
```

-o: Write output to the console.
```
go run . -o
```


## Text Processing
### Airport Code Conversion*

The program converts airport codes (e.g., IATA and ICAO codes) found in the text data to their corresponding airport names and municipalities.
 It uses an airport lookup CSV file to perform this conversion.

### ISO Date and Time Conversion

The program converts ISO-formatted date and time strings in the text data to a more human-readable format.
 It supports both 12-hour and 24-hour time formats and adjusts the time zone offset accordingly.

### Text Formatting

The program allows users to define custom text formatting settings in a settings file (user_settings.txt). 
It supports formatting options such as color, bold, italic, underline, and strikethrough.
 These settings are applied to the processed text output.


## Error Handling

The program handles various error scenarios, including:

- Missing or incorrect command-line arguments.
- Missing or malformed input files (e.g., input text file, airport lookup CSV file).
- Errors during file operations (e.g., reading, writing, removing files).
- Errors in parsing date and time strings.


## Logging

The program logs errors and informational messages to a log file (app.log). 
It records details such as file names, line numbers, and error descriptions to aid in troubleshooting.


## Customization

Users can customize the program's behavior by modifying the following components:

- Input, output and Airport lookup CSV file paths.
- Text formatting settings in the user_settings.txt file.


## Bonuses

### Error logging

Error and informational message logging to the log file (app.log)
To ensure error and informational message tracking, the application uses the app.log log file. 
This file stores messages about any events that may be useful for debugging or monitoring the application's operation.

### Dynamic Airport Lookup

The program works with non-standard column orders when searching for airports.
 This means it can correctly handle CSV files with airport data where columns may be in any order.
 
### There are three files named airport-lookup.csv:

- **airport-lookup1.csv** :  standard file without any modifications.
- **airport-lookup2.csv** : with a non-standard column order (for testing the bonus feature).
- **airport-lookup3.csv** : has a **'name'** column with no data and a **'name1'** column that contains data (for testing errors).
To apply one of the options, you need to remove the digit.

### Text Formatting

The application allows text formatting according to provided settings. 
It uses arguments specified in the user_settings.txt settings file. 
Users can set formatting parameters such as text color, boldness, italics, underline,
 and strikethrough for various types of text, such as airport names and dates.

In the user_settings.txt file, you can customize text formatting parameters for various elements:

- Date Formatting
- Time Formatting
- Positive Offset Formatting
- Negative Offset Formatting
- Airport Formatting
- City Formatting

You have the following text formatting options available:

- Color: Choose any color by specifying its ANSI code (e.g., 31 for red).
- Bold: Disabled/Enabled
- Italic: Disabled/Enabled
- Underline: Disabled/Enabled
- Strikethrough: Disabled/Enabled
 
You can mix and match text styles in any combination that suits your preferences.



### Error Handling

The application provides error handling in several cases:

- If the app.log log file cannot be opened or created, the program exits with a fatal error.
- If input or airport lookup files are not found, the program outputs the corresponding message and exits with an error code.
- If the airport data structure in the CSV file does not match the expected structure, the program outputs an error message and exits.
- During text processing, if errors occur while parsing dates/times, they will be logged in the log file.

### City Name Conversion from Airport Codes
The application provides functionality to replace airport codes in the text with their corresponding city names.
 For example, the code #JFK will be replaced with the city name corresponding to this airport. Similarly, 
the code *#JFK will be replaced with the airport name.

### Digital Signature Verification
The program checks the digital signature of the input data against an expected signature.
