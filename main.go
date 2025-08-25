package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CronField struct {
	Name   string
	Values []int
}

func parse(cronString string) ([]CronField, string) {
	parts := strings.Fields(cronString)
	fields := []CronField{
		{"minute", []int{}},
		{"hour", []int{}},
		{"day of month", []int{}},
		{"month", []int{}},
		{"day of week", []int{}},
	}

	fields[0].Values = parseField(parts[0], 0, 59)
	fields[1].Values = parseField(parts[1], 0, 23)
	fields[2].Values = parseField(parts[2], 1, 31)
	fields[3].Values = parseField(parts[3], 1, 12)
	fields[4].Values = parseField(parts[4], 1, 7)
	command := strings.Join(parts[5:], " ")
	return fields, command
}

func parseField(field string, min, max int) []int {
	if field == "*" {
		return generateRange(min, max, 1)
	}
	var values []int
	parts := strings.Split(field, ",")
	for _, part := range parts {
		partValues := parsePart(part, min, max)
		values = append(values, partValues...)
	}
	return values
}

func parsePart(part string, min, max int) []int {
	if strings.Contains(part, "/") {
		return parseStep(part, min, max)
	}
	if strings.Contains(part, "-") {
		return parseRange(part, min, max)
	}
	num, _ := strconv.Atoi(part)
	return []int{num}
}

func parseStep(part string, min, max int) []int {
	stepParts := strings.Split(part, "/")
	step, _ := strconv.Atoi(stepParts[1])
	var baseValues []int
	if stepParts[0] == "*" {
		baseValues = generateRange(min, max, 1)
	} else if strings.Contains(stepParts[0], "-") {
		baseValues = parseRange(stepParts[0], min, max)
	} else {
		num, _ := strconv.Atoi(stepParts[0])
		baseValues = []int{num}
	}
	var result []int
	for i := 0; i < len(baseValues); i += step {
		result = append(result, baseValues[i])
	}
	return result
}

func parseRange(part string, min, max int) []int {
	rangeParts := strings.Split(part, "-")
	start, _ := strconv.Atoi(rangeParts[0])
	end, _ := strconv.Atoi(rangeParts[1])
	return generateRange(start, end, 1)
}

func generateRange(start, end, step int) []int {
	var result []int
	for i := start; i <= end; i += step {
		result = append(result, i)
	}
	return result
}

func formatOutput(fields []CronField, command string) string {
	var result strings.Builder
	for _, field := range fields {
		fieldName := field.Name
		result.WriteString(fmt.Sprintf("%-13s ", fieldName))
		for i, value := range field.Values {
			if i > 0 {
				result.WriteString(" ")
			}
			result.WriteString(strconv.Itoa(value))
		}
		result.WriteString("\n")
	}
	result.WriteString(fmt.Sprintf("%-13s %s\n", "command", command))
	return result.String()
}

func main() {
	cronString := os.Args[1]
	fields, command := parse(cronString)
	output := formatOutput(fields, command)
	fmt.Print(output)
}
