package hw03frequencyanalysis

import (
	"math"
	"regexp"
	"sort"
	"strings"
)

func Top10(text string, fineMode bool) []string {
	words := getTextWords(text, fineMode)
	if len(words) == 0 {
		return []string{}
	}

	// подсчитываем кол-во слов
	wordsCounts := make(map[string]int)
	for _, word := range words {
		wordsCounts[word]++
	}

	type wordInfo struct {
		Word  string
		Count int
	}

	// конвертим в слайс
	wordsSlice := make([]wordInfo, 0)
	for word, count := range wordsCounts {
		wordsSlice = append(wordsSlice, wordInfo{word, count})
	}

	// сортируем
	sort.Slice(wordsSlice, func(i, j int) bool {
		return (wordsSlice[i].Count == wordsSlice[j].Count && wordsSlice[i].Word < wordsSlice[j].Word) ||
			wordsSlice[i].Count > wordsSlice[j].Count
	})

	// обрезаем
	resultSliceLength := int(math.Min(10, float64(len(wordsSlice))))
	resultSlice := wordsSlice[:resultSliceLength]

	// конвертим в слайс строк
	result := make([]string, 0)
	for _, element := range resultSlice {
		result = append(result, element.Word)
	}

	return result
}

func getTextWords(text string, fineMode bool) []string {
	if fineMode {
		text = regexp.MustCompile(`[.,!:;'"]+`).ReplaceAllString(text, " ")
		text = regexp.MustCompile(`(-\s)+|(\s-)+`).ReplaceAllString(text, " ")
		text = strings.ToLower(text)
	}

	text = regexp.MustCompile(`[\s|\n]+`).ReplaceAllString(text, " ")
	text = strings.Trim(text, " ")
	return strings.Fields(text)
}
