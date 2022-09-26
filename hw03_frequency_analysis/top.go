package hw03frequencyanalysis

import (
	"math"
	"regexp"
	"sort"
	"strings"
)

var (
	symbolsRegexp = regexp.MustCompile(`[.,!:;'"]+`)
	hyphenRegexp  = regexp.MustCompile(`(-\s)+|(\s-)+`)
	spacesRegexp  = regexp.MustCompile(`[\s|\n]+`)
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
	counter := 0
	wordsSlice := make([]wordInfo, len(wordsCounts))
	for word, count := range wordsCounts {
		wordsSlice[counter] = wordInfo{word, count}
		counter++
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
	result := make([]string, len(resultSlice))
	for i, element := range resultSlice {
		result[i] = element.Word
	}

	return result
}

func getTextWords(text string, fineMode bool) []string {
	if fineMode {
		text = symbolsRegexp.ReplaceAllString(text, " ")
		text = hyphenRegexp.ReplaceAllString(text, " ")
		text = strings.ToLower(text)
	}

	text = spacesRegexp.ReplaceAllString(text, " ")
	text = strings.Trim(text, " ")
	return strings.Fields(text)
}
