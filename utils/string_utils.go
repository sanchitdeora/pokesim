package utils

import(
	"strings"
)

func ToCapitalizeFirstLetterOfEachWord(phrase string) string {
	if len(phrase) == 0 {
		return phrase
	}

	words := strings.Split(phrase, " ")
	var updatedWords []string
	if len(words) > 1 {
		for _, word := range words {
			updatedWords = append(updatedWords, toCapitalizeFirstLetter(word))
		}
		return strings.Join(updatedWords, "")
	}

	words = strings.Split(phrase, "-")
	if len(words) > 1 {
		for _, word := range words {
			updatedWords = append(updatedWords, toCapitalizeFirstLetter(word))
		}
		return strings.Join(updatedWords, " ")
	}

	return toCapitalizeFirstLetter(phrase)
}

func toCapitalizeFirstLetter(word string) string {
	if len(word) == 0 {
		return word
	}
	return strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
}
