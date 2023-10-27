package asciiart

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Function that takes arguments and banner and
// prints the corresponding ascii-art
func Execute() {
	args := os.Args[1:]
	if VerifyArguments(args) {
		output, err := findSubstringAfterEqualSign(args[0])
		if err != nil{
			fmt.Println(err)
			os.Exit(0)
		}
		if len(args[1]) > 0 && TextVerification(args[1]) {
			if lines := strings.Split(args[1], `\n`); EmptyLines(lines) {
				for i := 0; i < len(lines)-1; i++ {
					WriteFile(output, `\n`)
				}
			} else if len(args) == 2 {
				WriteFile(output, args[1])
			} else if len(args) == 3 {
				WriteFile(output, args[1], args[2])
			}
		} else if len(args[1]) > 0 {
			fmt.Println("Votre texte contient des caracteres non pris en charge.")
		}
	} else {
		fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]\nEX: go run . --output=<fileName.txt> something standard")
	}
}

// Function that gets all characters from
// the given ascii-art file
func GetAllChars(banner ...string) map[rune][]string {
	bannerFile := "standard"
	if len(banner) != 0 {
		bannerFile = banner[0]
	}

	var char []string
	var count rune = 32
	chars := make(map[rune][]string)
	lines := ReadFile(bannerFile)

	for i, val := range lines {
		if (i+1)%9 == 0 {
			chars[count] = char
			char = []string{}
			count++
		} else {
			char = append(char, val)
		}
	}
	return chars
}

// Read lines from the given banner
func ReadFile(bannerFile string) []string {
	s, err := os.ReadFile(bannerFile + ".txt")
	if err == nil {
		// Deletion of carriage ret ("\r") noticed inside "thinkertoy" file
		lines := strings.Split(strings.ReplaceAll(string(s), "\r", ""), "\n")[1:]
		return lines
	} else {
		fmt.Println("INVALID BANNER")
		os.Exit(0)
		return nil
	}
}



// Function that returns the ascii-art text corresponding
// to a given string
func GetChars(s string, banner ...string) [][]string {
	allChars := GetAllChars(banner...)
	var charsTab [][]string
	for _, val := range s {
		charsTab = append(charsTab, allChars[rune(val)])
	}
	return charsTab
}

// Function that receive a text and return
// it's ascii-art printable text
func TextToPrint(s string, banner ...string) string {
	lines := strings.Split(s, `\n`)
	text := ""
	for i, line := range lines {
		chars := GetChars(line, banner...)
		if len(chars) > 0 {
			for i := 0; i < len(chars[0]); i++ {
				for _, char := range chars {
					if len(char) > i {
						text += char[i]
					}
				}
				if i < len(chars[0])-1 {
					text += "\n"
				}
			}
		}
		if i < len(lines)-1 {
			text += "\n"
		}
	}
	return text
}

// Function that checks if all the lines are empty
func EmptyLines(lines []string) bool {
	for _, line := range lines {
		if line != "" {
			return false
		}
	}
	return true
}

// Function that prints the ascii-art for
// the given text
func WriteFile(filename, s string, banner ...string) {
	text := TextToPrint(s, banner...)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	if text == "\n"{
		_, err = file.WriteString(text)
		if err != nil {
			fmt.Println(err)
		}
	}else{
		_, err = file.WriteString(text+"\n")
		if err != nil {
			fmt.Println(err)
		}
	}
}

// Function that checks if the string contains
// only characters that are in the given file
func TextVerification(s string) bool {
	re := regexp.MustCompile(`[^[:ascii:]]`)
	return len(re.FindAllString(s, -1)) == 0
}

// Function that checks if the correct number of command line arguments is provided.
func VerifyArguments(args []string) bool {
	return len(args) == 2 || len(args) == 3
}

func findSubstringAfterEqualSign(input string) (string, error) {
	beforeEqual := regexp.MustCompile(`([^= ]+)`)
	afterEqual := regexp.MustCompile(`=(.*)`)
	
	// Find all matches
	before := beforeEqual.FindStringSubmatch(input)
	if before[0] != "--output"{
		fmt.Println("Usage: go run www. [OPTION] [STRING] [BANNER]\nEX: go run . --output=<fileName.txt> something standard")
		os.Exit(0)
	}
	after := afterEqual.FindStringSubmatch(input) //this is the file name
	if len(after) == 0 {
		return "", fmt.Errorf("Usage: go run . [OPTION] [STRING] [BANNER]\nEX: go run . --output=<fileName.txt> something standard")
	}

	return after[1], nil
}
