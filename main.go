package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Usage msg
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: randstr [-d path] [-l length] [-o s(symbols)/n(numbers)/l(lowercase)/u(uppercase)]")
	}

	lengthFlag := flag.Int("l", 6, "Length of the random string")
	optionsFlag := flag.String("o", "nl", "Character sets: s (symbols), n (numbers), l (lowercase), u (uppercase)")
	dirFlag := flag.String("d", "", "Directory to check for filename collisions")
	flag.Parse()

	// Validate options flag contains only allowed characters (s, n, l, u)
	for _, ch := range *optionsFlag {
		if ch != 's' && ch != 'n' && ch != 'l' && ch != 'u' {
			fmt.Fprintf(os.Stderr, "Error: Invalid character '%c' in options flag. Use only combinations of s, n, l, u.\n", ch)
			os.Exit(1)
		}
	}

	// Build character set based on options
	var charset string
	if strings.Contains(*optionsFlag, "s") {
		charset += "!@#$%^&*()_+-=[]{}|;:,.<>?"
	}
	if strings.Contains(*optionsFlag, "n") {
		charset += "0123456789"
	}
	if strings.Contains(*optionsFlag, "l") {
		charset += "abcdefghijklmnopqrstuvwxyz"
	}
	if strings.Contains(*optionsFlag, "u") {
		charset += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}

	if charset == "" {
		fmt.Fprintln(os.Stderr, "Error: Character set cannot be empty.")
		os.Exit(1)
	}

	// Load existing file basenames into an in-memory map if -d is provided
	existingFiles := make(map[string]bool)
	if *dirFlag != "" {
		if _, err := os.Stat(*dirFlag); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: Directory '%s' does not exist.\n", *dirFlag)
			os.Exit(1)
		}

		entries, err := os.ReadDir(*dirFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading directory: %v\n", err)
			os.Exit(1)
		}

		for _, entry := range entries {
			name := entry.Name()
			ext := filepath.Ext(name)
			baseName := strings.TrimSuffix(name, ext)
			existingFiles[baseName] = true
		}
	}

	// Generate a unique random string
	for {
		randStr, err := generateRandomString(charset, *lengthFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating random string: %v\n", err)
			os.Exit(1)
		}

		if *dirFlag == "" || !existingFiles[randStr] {
			fmt.Println(randStr)
			break
		}
	}
}

func generateRandomString(charset string, length int) (string, error) {
	bytes := make([]byte, length)
	max := big.NewInt(int64(len(charset)))
	for i := range length {
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		bytes[i] = charset[num.Int64()]
	}
	return string(bytes), nil
}
