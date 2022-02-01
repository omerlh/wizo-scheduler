package main
import . "github.com/ahmetb/go-linq/v3"

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
		"github.com/jessevdk/go-flags"
)

type Row struct {
	Kids string
	Location string
	Day string
	Teacher string
	Hour string
	Group string
}

const (
	groupColumn = 4
	hourColumn = 3
	teacherColumn = 1
	dayColumn = 2
	locationColumn = 5
	kidsColumn = 6
)

var (
	fileName string
	groupBy string
)

type Options struct {
	FileName string `short:"f" long:"file" description:"The name of the CSV file to parse" required:"true"`
	GroupBy string `short:"g" long:"groupBy" description:"The field to group by (teacher or group)" choice:"teacher" choice:"group" default:"teacher"`
}


func main() {
	var options Options

	parser := flags.NewParser(&options, flags.Default)

	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}

	f, err := os.Open(options.FileName)

	if err != nil {
		panic(err)
	}

	r := csv.NewReader(f)

	// Skip headers
	_, err = r.Read()

	if err != nil {
		panic(err)
	}

	rows := make([]Row, 1)

	for {

			record, err := r.Read()

			if err == io.EOF {
					break
			}

			if err != nil {
				panic(err)
			}

			row := Row{
				Kids: record[kidsColumn],
				Location: record[locationColumn],
				Day: record[dayColumn],
				Teacher: record[teacherColumn],
				Hour: record[hourColumn],
				Group: record[groupColumn],
			}

			rows = append(rows, row)
	}

	var grouped []interface{}

	From(rows).GroupBy(func(row interface{}) interface{} {
		if (options.GroupBy == "teacher") {
			return row.(Row).Teacher 
		} else {
			return row.(Row).Group
		}
		// author as key
	}, func(row interface{}) interface{} {
		return row // author as value
	}).ToSlice(&grouped)

	for _, v := range(grouped) {
		
		fileName := fmt.Sprintf("%s.html", v.(Group).Key)

		f, err := os.Create(fileName)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		f.WriteString(fmt.Sprintf(`<html dir="rtl" lang="he">
		<head>
		<meta charset="utf-8">
<style>
table {
  font-family: arial, sans-serif;
  border-collapse: collapse;
  width: 100%;
}

td, th {
  border: 1px solid #dddddd;
  text-align: right;
  padding: 8px;
}

td.bold {
	font-weight: bold;
	font-size: larger;
}

</style>
</head><body>`))

		f.WriteString(fmt.Sprintf("<h1>%s</h1>", v.(Group).Key))

		var byDay []interface{}

		From(v.(Group).Group).GroupBy(func(row interface{}) interface{} {
			return row.(Row).Day // author as key
		}, func(row interface{}) interface{} {
			return row // author as value
		}).ToSlice(&byDay)

		f.WriteString(fmt.Sprintf("<table><tr><th>זמן</th><th>מיקום</th><th>קבוצה</th><th>חניכות</th></tr>"))
		
		for _, day := range(byDay) {
			f.WriteString(fmt.Sprintf("<tr><td class=bold>%s</td><td></td><td></td><td></td></tr>\n", day.(Group).Key))

			var sorted []interface{}

			From(day.(Group).Group).OrderByDescending( // sort groups by its length
				func(group interface{}) interface{} {
					return len(group.(Row).Hour)
				}).ToSlice(&sorted)

			for _, hour := range(sorted) {
				f.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>\n", hour.(Row).Hour, hour.(Row).Location, hour.(Row).Group, hour.(Row).Kids))
			}
		}

		f.WriteString(fmt.Sprintf("</table></html></body>"))
	}

		
}