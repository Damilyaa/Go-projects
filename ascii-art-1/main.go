package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func main() {
	var textFromOutside string

	// Проверка аргументов командной строки.
	if len(os.Args) > 1 {
		textFromOutside = strings.Join(os.Args[1:], " ")
		if strings.Count(textFromOutside, "\"")%2 != 0 {
			fmt.Println("Unclosed quotation mark detected. Waiting for the closing mark...")
			reader := bufio.NewReader(os.Stdin)
			for {
				input, _ := reader.ReadString('\n')
				textFromOutside += input
				if strings.Count(textFromOutside, "\"")%2 == 0 {
					break
				}
				fmt.Println("Still waiting for closing quotation mark...")
			}
		}
		firstQuoteIndex := strings.Index(textFromOutside, "\"")
		if firstQuoteIndex != -1 {
			textFromOutside = textFromOutside[firstQuoteIndex+1:]
		}
		textFromOutside = strings.Trim(textFromOutside, "\"")
	}

	// Если есть текст для обработки.
	if textFromOutside != "" {
		fileLines := ReadStandardTxt()
		asciiTemplates := return2dASCIIArray(fileLines)
		printAllStringASCII(textFromOutside, asciiTemplates)
	}

	// Проверка хэшей файлов.
	if err := checkFileHashes(); err != nil {
		fmt.Println("Error:", err)
		return
	}
}

// Читает строки из файла и возвращает их в виде среза строк.
func ReadStandardTxt() []string {
	readFile, err := os.Open("banners/standard.txt")
	if err != nil {
		fmt.Println("Error: Unable to open the file.")
		os.Exit(1)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	return fileLines
}

// Преобразует строки из файла в двумерный массив ASCII-шаблонов.
func return2dASCIIArray(fileLines []string) [][]string {
	var asciiTemplates [][]string
	counter := 0
	var tempAsciArray []string

	for _, line := range fileLines {
		counter++
		if counter != 1 {
			tempAsciArray = append(tempAsciArray, line)
		}
		if counter == 9 {
			asciiTemplates = append(asciiTemplates, tempAsciArray)
			counter = 0
			tempAsciArray = nil
		}
	}
	return asciiTemplates
}

// Печатает символы с использованием ASCII-шаблонов.
func printMultipleCharacter(s string, asciiTemplates [][]string) {
	tempIntArrLetter, err := returnAsciiCodeInt(s)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < 8; i++ {
		for _, v := range tempIntArrLetter {
			fmt.Print(asciiTemplates[v][i])
		}
		fmt.Println()
	}
}

// Возвращает массив ASCII-кодов для строки.
func returnAsciiCodeInt(s string) ([]int, error) {
	var tempIntArrLetter []int
	for _, v := range s {
		if !unicode.IsPrint(v) || v < 32 || v > 126 {
			return nil, fmt.Errorf("Error: Non-ASCII character detected!")
		}
		tempIntArrLetter = append(tempIntArrLetter, (int(v) - 32))
	}
	return tempIntArrLetter, nil
}

// Обрабатывает текст и печатает его с использованием ASCII-шаблонов.
func printAllStringASCII(text string, asciiTemplates [][]string) {
	text = strings.ReplaceAll(text, "\\n", "\n")

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if line == "" && i == len(lines)-1 {
			break
		}
		if line == "" {
			fmt.Println()
			continue
		}
		if strings.HasPrefix(line, "|") {
			fmt.Println()
			continue
		}
		printMultipleCharacter(line, asciiTemplates)
	}
}

// Ожидаемые хэши для проверки целостности файлов.
var expectedHashes = map[string]string{
	"banners/standard.txt": "e194f1033442617ab8a78e1ca63a2061f5cc07a3f05ac226ed32eb9dfd22a6bf",
}

// Вычисляет хэш-сумму файла.
func calculateFileHash(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("Error: Unable to open file %v", err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := bufio.NewReader(file).WriteTo(hasher); err != nil {
		return "", fmt.Errorf("Error: Unable to calculate hash %v", err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// Проверяет хэш-суммы файлов.
func checkFileHashes() error {
	for filename, expectedHash := range expectedHashes {
		actualHash, err := calculateFileHash(filename)
		if err != nil {
			return fmt.Errorf("Error checking hash for %s: %v", filename, err)
		}

		if actualHash != expectedHash {
			return fmt.Errorf("Hash mismatch for %s: expected %s, got %s", filename, expectedHash, actualHash)
		}
	}
	return nil
}
