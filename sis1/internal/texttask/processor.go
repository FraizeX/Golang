package texttask

import (
	"strings"
	"unicode/utf8"
)

type Processor interface {
	Process(text string) Result
}

type TextAnalyzer struct{}

func (TextAnalyzer) Process(text string) Result {
	lineCount := 0
	if text != "" {
		lineCount = strings.Count(text, "\n") + 1
	}

	return Result{
		WordCount:      len(strings.Fields(text)),
		CharacterCount: utf8.RuneCountInString(text),
		LineCount:      lineCount,
	}
}
