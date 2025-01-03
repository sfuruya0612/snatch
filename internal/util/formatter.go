package util

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type Column struct {
	Header string
	Width  int
}

type TableFormatter struct {
	columns []Column
	format  string // "tab" または "csv"
}

func NewTableFormatter(columns []Column, format string) *TableFormatter {
	return &TableFormatter{
		columns: columns,
		format:  format,
	}
}

func (f *TableFormatter) PrintHeader() {
	if f.format == "csv" {
		writer := csv.NewWriter(os.Stdout)
		headers := make([]string, len(f.columns))
		for i, col := range f.columns {
			headers[i] = col.Header
		}
		if err := writer.Write(headers); err != nil {
			fmt.Printf("Unable to write CSV header: %v\n", err)
			// TODO; エラーハンドリングちゃんとする
			os.Exit(1)
		}
		writer.Flush()
		return
	}

	format := f.createFormatString()
	headers := make([]interface{}, len(f.columns))
	for i, col := range f.columns {
		headers[i] = col.Header
	}
	fmt.Printf(format, headers...)
}

func (f *TableFormatter) PrintRows(rows [][]string) {
	if f.format == "csv" {
		writer := csv.NewWriter(os.Stdout)
		for _, row := range rows {
			if err := writer.Write(row); err != nil {
				fmt.Printf("Unable to write CSV row: %v\n", err)
				os.Exit(1)
			}
		}
		writer.Flush()
		return
	}

	format := f.createFormatString()
	for _, row := range rows {
		rowValues := make([]interface{}, len(row))
		for i, v := range row {
			rowValues[i] = v
		}
		fmt.Printf(format, rowValues...)
	}
}

func (f *TableFormatter) createFormatString() string {
	formats := make([]string, len(f.columns))
	for i, col := range f.columns {
		formats[i] = fmt.Sprintf("%%-%ds", col.Width)
	}
	return strings.Join(formats, "\t") + "\n"
}
