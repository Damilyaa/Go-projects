package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// adjustArticles корректирует артикли в тексте
func adjustArticles(txt string) string {
	words := strings.Fields(txt)
	var result []string
	exceptions := map[string]bool{
		"for": true, "and": true, "nor": true, "but": true,
		"or": true, "so": true, "yet": true,
	}
	// Adding exceptions for cases where 'an' is used even with consonant sounds.
	articleExceptionsU := map[string]bool{
		"universe": true, "unicode": true,
	}

	articleExceptionsH := map[string]bool{
		"hour": true, "heir": true, "honest": true,
	}

	for i := 0; i < len(words); i++ {
		word := words[i]
		if strings.ToLower(word) == "a" || strings.ToLower(word) == "an" {
			if i+1 < len(words) {
				nextWord := words[i+1]
				if exceptions[strings.ToLower(nextWord)] {
					result = append(result, word)
					continue
				}
				if articleExceptionsU[strings.ToLower(nextWord)] {
					if unicode.IsUpper(rune(word[0])) {
						result = append(result, "A")
					} else {
						result = append(result, "a")
					}
					continue
				}

				if articleExceptionsH[strings.ToLower(nextWord)] {
					if unicode.IsUpper(rune(word[0])) {
						result = append(result, "An")
					} else {
						result = append(result, "an")
					}
					continue
				}

				if isSingleCharacter(nextWord) {
					result = append(result, word)
					continue
				}
				if isVowel(rune(nextWord[0])) {
					if unicode.IsUpper(rune(word[0])) {
						if len(word) > 1 && unicode.IsUpper(rune(word[1])) {
							result = append(result, "AN")
						} else {
							result = append(result, "An")
						}
					} else {
						result = append(result, "an")
					}
				} else {
					if unicode.IsUpper(rune(word[0])) {
						result = append(result, "A")
					} else {
						result = append(result, "a")
					}
				}
			} else {
				result = append(result, word)
			}
		} else {
			result = append(result, word)
		}
	}
	return strings.Join(result, " ")
}

// changeCase изменяет регистр слов в тексте
func changeCase(inputText string) string {
	modsRegex := regexp.MustCompile(`\((up|low|cap)(,\s*\d+)?\)`)
	for {
		changed := false
		match := modsRegex.FindStringIndex(inputText)
		if match != nil {
			prefix := strings.TrimSpace(inputText[:match[0]])
			suffix := inputText[match[1]:]
			matchMod := inputText[match[0]:match[1]]
			mod, count := extractModAndCount(matchMod)
			if !containsLetterOrDigit(prefix) {
				inputText = prefix + suffix
				continue
			}
			modifiedPrefix := modifyPrefix(mod, count, prefix)
			inputText = modifiedPrefix + suffix
			changed = true
		}
		if !changed {
			break
		}
	}
	return strings.TrimSpace(inputText)
}

// formatPunctuation форматирует пунктуацию в тексте
func formatPunctuation(input string) string {
	re := regexp.MustCompile(`(.)(\s*)([.,!?;:]+)(.*?)`)
	punctuationRegex := regexp.MustCompile(`([.,!?;:])`)
	addSpace := punctuationRegex.ReplaceAllString(input, "$1 ")
	MatchText := re.ReplaceAllStringFunc(addSpace, func(match string) string {
		re2 := regexp.MustCompile(`(\s*)`)
		FinalStep := re2.ReplaceAllString(match, "")
		return FinalStep
	})
	return MatchText
}

// hexadecimalToDecimal преобразует шестнадцатеричные числа в десятичные
func hexadecimalToDecimal(input string) string {
	re := regexp.MustCompile(`\(\s*hex\s*\)`)
	for {
		char := re.FindStringIndex(input)
		if char == nil {
			break
		}
		before := input[:char[0]]
		after := input[char[1]:]
		before = strings.TrimRight(before, " ")
		if len(before) == 0 {
			input = after
			continue
		}
		lastSpace := strings.LastIndex(before, " ")
		var word, beforeWord string
		if lastSpace == -1 {
			word = before
			beforeWord = ""
		} else {
			word = before[lastSpace+1:]
			beforeWord = before[:lastSpace]
		}
		decimalNum, err := strconv.ParseInt(word, 16, 64)
		if err != nil {
			input = beforeWord + " " + word + after
			continue
		}
		newWord := fmt.Sprintf("%d", decimalNum)
		if beforeWord == "" {
			input = newWord + after
		} else {
			input = beforeWord + " " + newWord + after
		}
	}
	return strings.TrimSpace(input)
}

// binaryToDecimal преобразует двоичные числа в десятичные
func binaryToDecimal(input string) string {
	re := regexp.MustCompile(`\b([01]+)\s?\(bin\)`)
	return re.ReplaceAllStringFunc(input, func(match string) string {
		submatches := re.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}
		binValue := submatches[1]
		decimalValue, err := strconv.ParseInt(binValue, 2, 64)
		if err != nil {
			return match
		}
		return fmt.Sprintf("%d", decimalValue)
	})
}

// cleanWhitespace очищает лишние пробелы в тексте
func cleanWhitespace(input string) string {
	re := regexp.MustCompile(`\s+`)
	adjusted := re.ReplaceAllString(input, " ")
	return strings.TrimSpace(adjusted)
}

// formatSingleQuotes форматирует одинарные кавычки в тексте
func formatSingleQuotes(input string) string {
	runes := []rune(input)
	var output []rune
	i := 0
	for i < len(runes) {
		if runes[i] == '\'' {
			nextQuote := -1
			for j := i + 1; j < len(runes); j++ {
				if runes[j] == '\'' {
					nextQuote = j
					break
				}
			}
			if nextQuote != -1 {
				if len(output) > 0 && !unicode.IsSpace(output[len(output)-1]) {
					output = append(output, ' ')
				}
				output = append(output, '\'')
				content := strings.TrimSpace(string(runes[i+1 : nextQuote]))
				output = append(output, []rune(content)...)
				output = append(output, '\'')
				if nextQuote+1 < len(runes) {
					j := nextQuote + 1
					for j < len(runes) && unicode.IsSpace(runes[j]) {
						j++
					}
					if j < len(runes) && (runes[j] == '.' || runes[j] == ',' || runes[j] == ':' || runes[j] == ';' || runes[j] == '?' || runes[j] == '!' || runes[j] == '"') {
						i = j
					} else {
						output = append(output, ' ')
						i = nextQuote + 1
					}
				} else {
					output = append(output, ' ')
					i = nextQuote + 1
				}
			} else {
				output = append(output, '\'')
				i++
			}
		} else {
			output = append(output, runes[i])
			i++
		}
	}
	return strings.TrimSpace(string(output))
}

// formatDoubleQuotes форматирует двойные кавычки в тексте
func formatDoubleQuotes(input string) string {
	runes := []rune(input)
	var output []rune
	i := 0
	for i < len(runes) {
		if runes[i] == '"' {
			nextQuote := -1
			for j := i + 1; j < len(runes); j++ {
				if runes[j] == '"' {
					nextQuote = j
					break
				}
			}
			if nextQuote != -1 {
				if len(output) > 0 && !unicode.IsSpace(output[len(output)-1]) {
					output = append(output, ' ')
				}
				output = append(output, '"')
				content := strings.TrimSpace(string(runes[i+1 : nextQuote]))
				output = append(output, []rune(content)...)
				output = append(output, '"')
				if nextQuote+1 < len(runes) {
					j := nextQuote + 1
					for j < len(runes) && unicode.IsSpace(runes[j]) {
						j++
					}
					if j < len(runes) && (runes[j] == '.' || runes[j] == ',' || runes[j] == ':' || runes[j] == ';' || runes[j] == '?' || runes[j] == '!' || runes[j] == '"') {
						i = j
					} else {
						output = append(output, ' ')
						i = nextQuote + 1
					}
				} else {
					output = append(output, ' ')
					i = nextQuote + 1
				}
			} else {
				output = append(output, '"')
				i++
			}
		} else {
			output = append(output, runes[i])
			i++
		}
	}
	return strings.TrimSpace(string(output))
}

// formatQuotes форматирует одинарные и двойные кавычки в тексте
func formatQuotes(input string) string {
	input = formatSingleQuotes(input)
	input = formatDoubleQuotes(input)
	return input
}

// sanitizeText очищает текст от лишних пробелов и форматирует скобки
func sanitizeText(input string) string {
	input = regexp.MustCompile(`\s+`).ReplaceAllString(input, " ")
	input = regexp.MustCompile(`\(\s*([^\)]*?)\s*\)`).ReplaceAllString(input, `($1)`)
	return strings.TrimSpace(input)
}

// isSingleCharacter проверяет, является ли слово одной буквой
func isSingleCharacter(w string) bool {
	return len(w) == 1 && unicode.IsLetter(rune(w[0]))
}

// isVowel проверяет, является ли символ гласной буквой
func isVowel(c rune) bool {
	vowels := "aeiouAEIOU"
	return strings.ContainsRune(vowels, c)
}

// containsLetterOrDigit проверяет, содержит ли строка буквы или цифры
func containsLetterOrDigit(s string) bool {
	for _, char := range s {
		if char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char >= '0' && char <= '9' {
			return true
		}
	}
	return false
}

// extractModAndCount извлекает модификатор и количество из строки
func extractModAndCount(tag string) (string, int) {
	parts := strings.Split(strings.Trim(tag, "()"), ",")
	mod := strings.TrimSpace(parts[0])
	count := 1
	if len(parts) > 1 {
		fmt.Sscanf(strings.TrimSpace(parts[1]), "%d", &count)
	}
	return mod, count
}

// modifyPrefix изменяет регистр слов в префиксе
func modifyPrefix(mod string, count int, prefix string) string {
	words := strings.Fields(prefix)
	lastWordIndex := len(words) - 1
	switch mod {
	case "up":
		for i := 0; i < count && lastWordIndex-i >= 0; i++ {
			words[lastWordIndex-i] = strings.ToUpper(words[lastWordIndex-i])
		}
	case "low":
		for i := 0; i < count && lastWordIndex-i >= 0; i++ {
			words[lastWordIndex-i] = strings.ToLower(words[lastWordIndex-i])
		}
	case "cap":
		for i := 0; i < count && lastWordIndex-i >= 0; i++ {
			words[lastWordIndex-i] = capitalizeTitle(words[lastWordIndex-i])
		}
	}
	return strings.Join(words, " ")
}

// capitalizeTitle делает первую букву заглавной, а остальные строчными
func capitalizeTitle(s string) string {
	if len(s) > 0 {
		return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
	}
	return s
}

// transformText обрабатывает входной текст
func transformText(input string) string {
	input = sanitizeText(input)
	input = formatPunctuation(input)
	input = changeCase(input)
	input = hexadecimalToDecimal(input)
	input = binaryToDecimal(input)
	input = formatQuotes(input)
	input = adjustArticles(input)
	input = cleanWhitespace(input)
	return input
}

func main() {
	inputFile := "sample.txt"  // Имя входного файла
	outputFile := "result.txt" // Имя выходного файла
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()
	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer outFile.Close()
	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(outFile)
	for scanner.Scan() {
		line := scanner.Text()
		processedLine := transformText(line)            // Обработка текста
		_, _ = writer.WriteString(processedLine + "\n") // Запись в выходной файл
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
	writer.Flush()
	fmt.Println("Check", outputFile)
}
