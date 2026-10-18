package texttask

import "testing"

func TestTextAnalyzer(t *testing.T) {
	testCases := []struct {
		name string
		text string
		want Result
	}{
		{
			name: "English text",
			text: "Hello Go world",
			want: Result{WordCount: 3, CharacterCount: 14, LineCount: 1},
		},
		{
			name: "Cyrillic text",
			text: "Привет мир",
			want: Result{WordCount: 2, CharacterCount: 10, LineCount: 1},
		},
		{
			name: "mixed whitespace",
			text: " Go\tGo\nGo ",
			want: Result{WordCount: 3, CharacterCount: 10, LineCount: 2},
		},
		{
			name: "Unicode space and emoji",
			text: "Go\u00a0🙂",
			want: Result{WordCount: 2, CharacterCount: 4, LineCount: 1},
		},
		{
			name: "punctuation stays in a word",
			text: "hello,world!",
			want: Result{WordCount: 1, CharacterCount: 12, LineCount: 1},
		},
		{
			name: "trailing newline",
			text: "Go\n",
			want: Result{WordCount: 1, CharacterCount: 3, LineCount: 2},
		},
		{
			name: "Windows newline",
			text: "Go\r\nGo",
			want: Result{WordCount: 2, CharacterCount: 6, LineCount: 2},
		},
		{
			name: "empty text in the standalone analyzer",
			text: "",
			want: Result{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := (TextAnalyzer{}).Process(testCase.text)
			if got != testCase.want {
				t.Fatalf("Process(%q) = %+v, want %+v", testCase.text, got, testCase.want)
			}
		})
	}
}
