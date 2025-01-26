package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Caesar Cipher encryption and decryption")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Choose an option\n")
		fmt.Println("1. Encrypt text\n2. Decrypt text with known shift\n3. Decrypt text with unknown shift\n4. Exit")
		option, _ := reader.ReadString('\n')
		switch option {
		case "1\n":
			encryptInput()
		case "2\n":
			res, _ := decryptInputKnownShift()
			fmt.Println("Decrypted text:", res)
		case "3\n":
			res, shift, confidenceScore, allResults := decryptInputUnknownShift()
			fmt.Printf("Decrypted text with shift of %d, with confidence score of %f: %s\n", shift, confidenceScore, res)
			fmt.Println("All possible results? (y/n):")
			seeAll, _ := reader.ReadString('\n')
			if strings.TrimSpace(seeAll) == "y" {
				for shift, text := range allResults {
					fmt.Printf("Shift %d: %s\n", shift, text)
				}
			}
			continue
		case "4\n":
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Invalid option, please try again.")
			continue
		}
	}
}

func encrypt(text string, shift int) string {
	encrypted := ""
	for _, char := range text {
		if char >= 'a' && char <= 'z' {
			encrypted += string((char-'a'+rune(shift))%26 + 'a')
		}
		if char >= 'A' && char <= 'Z' {
			encrypted += string((char-'A'+rune(shift))%26 + 'A')
			if char >= '0' && char <= '9' {
				encrypted += string((char-'0'+rune(shift))%10 + '0')
			}
		}
	}
	return encrypted
}

func encryptInput() {
	var shift int
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text to encrypt: ")
	text, _ := reader.ReadString('\n')
	fmt.Print("Enter shift value: ")
	fmt.Scanf("%d", &shift)
	encrypted := encrypt(text, shift)
	fmt.Println("Encrypted text:", encrypted)
}

func decryptInputKnownShift() (string, int) {
	var shift int
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text to decrypt: ")
	text, _ := reader.ReadString('\n')
	fmt.Print("Enter shift value: ")
	fmt.Scanf("%d", &shift)
	decrypted := decryptKnownShift(text, shift)
	return decrypted, shift
}
func decryptKnownShift(text string, shift int) string {
	decrypted := ""
	for _, char := range text {
		if char >= 'a' && char <= 'z' {
			decrypted += string((char-'a'-rune(shift)+26)%26 + 'a')
		}
		if char >= 'A' && char <= 'Z' {
			decrypted += string((char-'A'-rune(shift)+26)%26 + 'A')
		}
		if char >= '0' && char <= '9' {
			decrypted += string((char-'0'-rune(shift)+10)%10 + '0')
		}
	}
	return decrypted
}

func decryptInputUnknownShift() (string, int, float64, map[int]string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text to decrypt: ")
	text, _ := reader.ReadString('\n')
	decrypted := decryptUnknownShift(text)
	var baseText string
	var baseShift int
	var highestConfidence float64
	for shift, text := range decrypted {
		confidenceScore := englishScore(shift, text)
		// Set the highest confidence score with the results
		if confidenceScore[shift] > highestConfidence {
			highestConfidence = confidenceScore[shift]
			baseText = text
			baseShift = shift
		}
	}
	return baseText, baseShift, highestConfidence, decrypted
}

// Calculate the confidence score of the decrypted text based on the presence of common English words
func englishScore(shift int, text string) map[int]float64 {
	// Digraphs are weighted lowest as they are more likely to appear in any text
	englishDigraphs := []string{"th", "er", "on", "an", "re", "he", "in", "ed", "nd", "ha", "at", "en", "es", "of", "or", "nt", "ea", "ti", "to", "it", "st", "io", "le", "is", "ou", "ar", "as", "de", "rt", "ve"}
	totalDigraphs := 0
	digraphWeight := 0.5
	// Trigraphs have a higher weight than digraphs, but are not as high as three or four letter words
	englisTrigraphs := []string{"the", "and", "tha", "ent", "ion", "tio", "for", "nde", "has", "nce", "edt", "tis", "oft", "sth", "men"}
	totalTrigraphs := 0
	trigraphWeight := 1.0
	// Three letter words are weighted higher than digraphs and trigraphs, but lower than four letter words
	threeLetterWords := []string{"the", "and", "for", "are", "but", "not", "you", "all", "any", "can", "had", "her", "was", "one", "our", "out", "day", "get", "has", "him", "his", "how", "man", "new", "now", "old", "see", "two", "way", "who", "boy", "did", "its", "let", "put", "say", "she", "too", "use"}
	totalThreeLetterWords := 0
	threeLetterWeight := 2.0
	// Four letter words are weighted highest
	fourLetterWords := []string{"that", "have", "this", "with", "from", "they", "which", "will", "there", "about", "these", "other", "some", "into", "than", "more", "time"}
	totalFourLetterWords := 0
	fourLetterWeight := 3.0

	for _, word := range englisTrigraphs {
		if strings.Contains(strings.ToLower(text), word) {
			totalTrigraphs++
		}
	}

	for _, word := range englishDigraphs {
		if strings.Contains(strings.ToLower(text), word) {
			totalDigraphs++
		}
	}

	for _, word := range threeLetterWords {
		if strings.Contains(strings.ToLower(text), word) {
			totalThreeLetterWords++
		}
	}

	for _, word := range fourLetterWords {
		if strings.Contains(strings.ToLower(text), word) {
			totalFourLetterWords++
		}
	}

	finalScore := float64(totalTrigraphs)*trigraphWeight + float64(totalDigraphs)*digraphWeight + float64(totalThreeLetterWords)*threeLetterWeight + float64(totalFourLetterWords)*fourLetterWeight

	result := make(map[int]float64)
	result[shift] = finalScore
	return result
}

func decryptUnknownShift(text string) map[int]string {
	fmt.Println("Trying all possible shifts and showing results...")
	results := make(map[int]string)
	for shift := 1; shift < 26; shift++ {
		decrypted := decryptKnownShift(text, shift)
		results[shift] = decrypted
	}
	return results
}
