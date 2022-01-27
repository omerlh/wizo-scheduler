package main
import . "github.com/ahmetb/go-linq/v3"

import (
    "encoding/csv"
    "fmt"
    "io"
    "log"
    "os"
)

type Row struct {
	Kids string
	Location string
	Day string
	Teacher string
	Hour string
	Group string
}

func main() {
    argsWithoutProg := os.Args[1:]

    f, err := os.Open(argsWithoutProg[0])

    if err != nil {

        log.Fatal(err)
    }

    r := csv.NewReader(f)

		_, err = r.Read()

		if err != nil {
			log.Fatal(err)
		}

		rows := make([]Row, 1)

    for {

        record, err := r.Read()

        if err == io.EOF {
            break
        }

        if err != nil {
            log.Fatal(err)
        }

        row := Row{
					Kids: record[9],
					Location: record[8],
					Day: record[2],
					Teacher: record[1],
					Hour: record[3],
					Group: record[4],
				}

				rows = append(rows, row)
    }

		var grouped []interface{}

		From(rows).GroupBy(func(row interface{}) interface{} {
			return row.(Row).Teacher // author as key
		}, func(row interface{}) interface{} {
			return row // author as value
		}).ToSlice(&grouped)

		for _, v := range(grouped) {
			
			fileName := fmt.Sprintf("%s.txt", v.(Group).Key)

			f, err := os.Create(fileName)
			if err != nil {
				panic(err)
			}

			defer f.Close()

			f.WriteString(fmt.Sprintf("%s", v.(Group).Key))

			var byDay []interface{}

			From(v.(Group).Group).GroupBy(func(row interface{}) interface{} {
				return row.(Row).Day // author as key
			}, func(row interface{}) interface{} {
				return row // author as value
			}).ToSlice(&byDay)

			
			for _, day := range(byDay) {
				f.WriteString(fmt.Sprintf("%s\n", day.(Group).Key))

				var sorted []interface{}

				From(day.(Group).Group).OrderByDescending( // sort groups by its length
					func(group interface{}) interface{} {
						return len(group.(Row).Hour)
					}).ToSlice(&sorted)

				for _, hour := range(sorted) {
					f.WriteString(fmt.Sprintf("%s\t%s\t%s\n", hour.(Row).Hour, hour.(Row).Location, hour.(Row).Group))
				}
			}

			fmt.Println("apd:", )
		}

		
}