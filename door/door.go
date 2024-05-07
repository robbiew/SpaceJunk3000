// BBS door utilities for handling ANSI art, user input, and dropfile data.

package door

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/eiannone/keyboard"
	"golang.org/x/text/encoding/charmap"
)

const (
	Esc = "\u001B["
	Osc = "\u001B]"
	Bel = "\u0007"
)

// Common fonts, supported by SyncTerm
const (
	Mosoul          = Esc + "0;38 D"
	Potnoodle       = Esc + "0;37 D"
	Microknight     = Esc + "0;41 D"
	Microknightplus = Esc + "0;39 D"
	Topaz           = Esc + "0;42 D"
	Topazplus       = Esc + "0;40 D"
	Ibm             = Esc + "0;0 D"
	Ibmthin         = Esc + "0;26 D"
)

// Common ANSI escapes sequences. This is not a complete list.
const (
	CursorBackward = Esc + "D"
	CursorPrevLine = Esc + "F"
	CursorLeft     = Esc + "G"
	CursorTop      = Esc + "d"
	CursorTopLeft  = Esc + "H"

	CursorBlinkEnable  = Esc + "?12h"
	CursorBlinkDisable = Esc + "?12I"

	ScrollUp   = Esc + "S"
	ScrollDown = Esc + "T"

	TextInsertChar = Esc + "@"
	TextDeleteChar = Esc + "P"
	TextEraseChar  = Esc + "X"
	TextInsertLine = Esc + "L"
	TextDeleteLine = Esc + "M"

	EraseRight  = Esc + "K"
	EraseLeft   = Esc + "1K"
	EraseLine   = Esc + "2K"
	EraseDown   = Esc + "J"
	EraseUp     = Esc + "1J"
	EraseScreen = Esc + "2J"

	Black     = Esc + "30m"
	Red       = Esc + "31m"
	Green     = Esc + "32m"
	Yellow    = Esc + "33m"
	Blue      = Esc + "34m"
	Magenta   = Esc + "35m"
	Cyan      = Esc + "36m"
	White     = Esc + "37m"
	BlackHi   = Esc + "30;1m"
	RedHi     = Esc + "31;1m"
	GreenHi   = Esc + "32;1m"
	YellowHi  = Esc + "33;1m"
	BlueHi    = Esc + "34;1m"
	MagentaHi = Esc + "35;1m"
	CyanHi    = Esc + "36;1m"
	WhiteHi   = Esc + "37;1m"

	BgBlack     = Esc + "40m"
	BgRed       = Esc + "41m"
	BgGreen     = Esc + "42m"
	BgYellow    = Esc + "43m"
	BgBlue      = Esc + "44m"
	BgMagenta   = Esc + "45m"
	BgCyan      = Esc + "46m"
	BgWhite     = Esc + "47m"
	BgBlackHi   = Esc + "40;1m"
	BgRedHi     = Esc + "41;1m"
	BgGreenHi   = Esc + "42;1m"
	BgYellowHi  = Esc + "43;1m"
	BgBlueHi    = Esc + "44;1m"
	BgMagentaHi = Esc + "45;1m"
	BgCyanHi    = Esc + "46;1m"
	BgWhiteHi   = Esc + "47;1m"

	Reset = Esc + "0m"
)

// Prompt the user and get their choice
func PromptYesNo(question string) (string, error) {
	// Print the prompt
	fmt.Printf("%s (yes/no)", question)

	// Open keyboard listener
	err := keyboard.Open()
	if err != nil {
		return "", err
	}
	defer keyboard.Close()

	// Listen for single key press
	char, _, err := keyboard.GetSingleKey()
	if err != nil {
		return "", err
	}

	// Convert the pressed key to lowercase string
	choice := string(unicode.ToLower(rune(char)))

	// Check if the choice is valid
	if choice != "y" && choice != "n" {
		fmt.Println("Invalid choice. Please enter 'y' or 'n'.")
		return PromptYesNo(question)
	}

	return choice, nil
}

// Move cursor to X, Y location
func MoveCursor(x int, y int) {
	fmt.Printf(Esc+"%d;%df", y, x)
}

// Erase the screen
func ClearScreen() {
	fmt.Println(EraseScreen)
	MoveCursor(0, 0)
}

// Show the cursor.
func CursorShow() {
	fmt.Print(Esc + "?25h")
}

// Hide the cursor.
func CursorHide() {
	fmt.Print(Esc + "?25l")
}

// WaitForAnyKey waits for a user to press any key to continue.
func WaitForAnyKey() error {
	// Open the keyboard listener
	err := keyboard.Open()
	if err != nil {
		return err
	}
	defer keyboard.Close() // Ensure that the keyboard listener is closed when done

	// Wait for a single key press
	_, _, err = keyboard.GetSingleKey()
	if err != nil {
		return err
	}

	return nil
}

func DisplayAnsiFile(filePath string, localDisplay bool) {
	content, err := ReadAnsiFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file %s: %v", filePath, err)
	}
	ClearScreen()
	PrintAnsi(content, 0, localDisplay)
}

func ReadAnsiFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// Print ANSI art with a delay between lines
func PrintAnsi(artContent string, delay int, localDisplay bool) { // localDisplay as an argument for UTF-8 conversion
	noSauce := TrimStringFromSauce(artContent) // strip off the SAUCE metadata
	lines := strings.Split(noSauce, "\r\n")

	for i, line := range lines {
		if localDisplay {
			// Convert line from CP437 to UTF-8
			utf8Line, err := charmap.CodePage437.NewDecoder().String(line)
			if err != nil {
				fmt.Printf("Error converting to UTF-8: %v\n", err)
				continue
			}
			line = utf8Line
		}

		if i < len(lines)-1 && i != 24 { // Check for the 25th line (index 24)
			fmt.Println(line) // Print with a newline
		} else {
			fmt.Print(line) // Print without a newline (for the 25th line and the last line of the art)
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
}

// Print ANSI art at an X, Y location after removing SAUCE metadata
func PrintAnsiLoc(artfile string, x, y int) error {
	// Open the file
	file, err := os.Open(artfile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Move to the specified Y coordinate
	fmt.Printf("\033[%d;%df", y, x)

	// Read and print each line at the specified location after removing SAUCE metadata
	for scanner.Scan() {
		line := scanner.Text()
		line = TrimStringFromSauce(line)
		fmt.Printf("%s\n", line)
		y++
		fmt.Printf("\033[%d;%df", y, x) // Move to the next line at the specified X coordinate
	}

	// Check for any scanner errors
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// TrimStringFromSauce trims SAUCE metadata from a string.
func TrimStringFromSauce(s string) string {
	return trimMetadata(s, "COMNT", "SAUCE00")
}

// trimMetadata trims metadata based on delimiters.
func trimMetadata(s string, delimiters ...string) string {
	for _, delimiter := range delimiters {
		if idx := strings.Index(s, delimiter); idx != -1 {
			return trimLastChar(s[:idx])
		}
	}
	return s
}

// trimLastChar trims the last character from a string.
func trimLastChar(s string) string {
	if len(s) > 0 {
		_, size := utf8.DecodeLastRuneInString(s)
		return s[:len(s)-size]
	}
	return s
}

// Print text at an X, Y location
func PrintStringLoc(text string, x int, y int) {
	fmt.Fprintf(os.Stdout, Esc+strconv.Itoa(y)+";"+strconv.Itoa(x)+"f"+text)
}

// CenterAlignText center-aligns text while preserving ANSI escape sequences and supports foreground and background colors.
func CenterAlignText(text string, width int, foreground, background string) string {
	// Regular expression to find ANSI escape sequences
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]+m`)

	// Split the input text into segments of plain text and ANSI escape sequences
	segments := ansiRegex.Split(text, -1)
	escapeSequences := ansiRegex.FindAllString(text, -1)

	// Calculate the total length of the input text (including ANSI escape sequences)
	totalLength := 0
	for _, segment := range segments {
		totalLength += len(segment)
	}
	totalLength += len(strings.Join(escapeSequences, ""))

	// Calculate the number of spaces needed for center alignment
	spacesNeeded := (width - totalLength) / 2

	// Prepare the center-aligned text
	var alignedText strings.Builder
	for i := 0; i < len(segments); i++ {
		if i == 0 && foreground != "" {
			alignedText.WriteString(foreground)
		}
		if i == 0 && background != "" {
			alignedText.WriteString(background)
		}
		alignedText.WriteString(strings.Repeat(" ", spacesNeeded))
		alignedText.WriteString(segments[i])
		if i < len(escapeSequences) {
			alignedText.WriteString(strings.Repeat(" ", spacesNeeded))
			alignedText.WriteString(escapeSequences[i])
		}
	}

	return alignedText.String()
}

// RightAlignText right-aligns text while preserving ANSI escape sequences and supports foreground and background colors.
func RightAlignText(text string, width int, foreground, background string) string {
	// Regular expression to find ANSI escape sequences
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]+m`)

	// Split the input text into segments of plain text and ANSI escape sequences
	segments := ansiRegex.Split(text, -1)
	escapeSequences := ansiRegex.FindAllString(text, -1)

	// Calculate the total length of the input text (including ANSI escape sequences)
	totalLength := 0
	for _, segment := range segments {
		totalLength += len(segment)
	}
	totalLength += len(strings.Join(escapeSequences, ""))

	// Calculate the number of spaces needed for right alignment
	spacesNeeded := width - totalLength

	// Prepare the right-aligned text
	var alignedText strings.Builder
	for i := 0; i < len(segments); i++ {
		if i < len(escapeSequences) {
			alignedText.WriteString(escapeSequences[i])
		}
		if i == 0 && foreground != "" {
			alignedText.WriteString(foreground)
		}
		if i == 0 && background != "" {
			alignedText.WriteString(background)
		}
		alignedText.WriteString(strings.Repeat(" ", spacesNeeded))
		alignedText.WriteString(segments[i])
	}

	return alignedText.String()
}

func DropFileData(path string) (string, int, int, int, error) {
	// Append trailing slash to path if it doesn't exist
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	// Check if the file exists
	fileInfo, err := os.Stat(strings.ToLower(path + "door32.sys"))
	if err != nil {
		if os.IsNotExist(err) {
			return "", 0, 0, 0, errors.New("file does not exist")
		}
		return "", 0, 0, 0, err
	}

	// Check if the file is empty
	if fileInfo.Size() == 0 {
		return "", 0, 0, 0, errors.New("file is empty")
	}

	// Open the file
	file, err := os.Open(strings.ToLower(path + "door32.sys"))
	if err != nil {
		return "", 0, 0, 0, err
	}
	defer file.Close()

	// Read lines from the file
	scanner := bufio.NewScanner(file)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	// Extract drop file data
	dropAlias := text[6]
	dropTimeLeft := text[8]
	dropEmulation := text[9]
	nodeNum := text[10]

	// Convert timeLeft and emulation to integers
	timeInt, err := strconv.Atoi(dropTimeLeft)
	if err != nil {
		return "", 0, 0, 0, err
	}
	emuInt, err := strconv.Atoi(dropEmulation)
	if err != nil {
		return "", 0, 0, 0, err
	}
	nodeInt, err := strconv.Atoi(nodeNum)
	if err != nil {
		return "", 0, 0, 0, err
	}

	return dropAlias, timeInt, emuInt, nodeInt, nil
}
