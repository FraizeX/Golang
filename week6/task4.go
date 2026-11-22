package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrBadFormat = errors.New("missing comma")

type Record struct {
	Name  string
	Grade int
}

func ParseAll(input string) ([]Record, error) {
	var records []Record
	var parseErrors []error

	if input == "" {
		return records, nil
	}

	// Один завершающий перевод строки считаем терминатором.
	input = strings.TrimSuffix(input, "\n")

	lines := strings.Split(input, "\n")

	for index, line := range lines {
		lineNumber := index + 1

		name, gradeText, found := strings.Cut(line, ",")

		if !found {
			parseErrors = append(
				parseErrors,
				fmt.Errorf(
					"line %d: %w",
					lineNumber,
					ErrBadFormat,
				),
			)

			continue
		}

		name = strings.TrimSpace(name)
		gradeText = strings.TrimSpace(gradeText)

		grade, err := strconv.Atoi(gradeText)

		if err != nil {
			parseErrors = append(
				parseErrors,
				fmt.Errorf(
					"line %d: %w",
					lineNumber,
					err,
				),
			)

			continue
		}

		records = append(records, Record{
			Name:  name,
			Grade: grade,
		})
	}

	return records, errors.Join(parseErrors...)
}

func main() {
	input := `Aruzhan, 95
Dias, 8O
Timur
Bota, 99999999999999999999
Ali, 70`

	records, err := ParseAll(input)

	fmt.Println("Корректные записи:")

	for _, record := range records {
		fmt.Printf("%s: %d\n", record.Name, record.Grade)
	}

	fmt.Println("\nВсе ошибки:")
	fmt.Println(err)

	fmt.Println("\nПроверка причин:")

	fmt.Println(
		"ErrBadFormat:",
		errors.Is(err, ErrBadFormat),
	)

	fmt.Println(
		"strconv.ErrSyntax:",
		errors.Is(err, strconv.ErrSyntax),
	)

	fmt.Println(
		"strconv.ErrRange:",
		errors.Is(err, strconv.ErrRange),
	)

	var numberError *strconv.NumError

	if errors.As(err, &numberError) {
		fmt.Println("\nПервая найденная NumError:")
		fmt.Println("Func:", numberError.Func)
		fmt.Println("Num:", numberError.Num)
		fmt.Println("Err:", numberError.Err)
	}

	fmt.Println("\nПроверка корректного ввода:")

	goodRecords, goodError := ParseAll(
		"Aruzhan, 95\nAli, 70",
	)

	fmt.Println(goodRecords)
	fmt.Println("Ошибка равна nil:", goodError == nil)
}