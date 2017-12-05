package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

/* solution by Jaro */

//Rec is loaded from csv
type Rec struct {
	Date     time.Time
	Open     float32
	High     float32
	Low      float32
	Close    float32
	Volume   float32
	AdjClose float32
}

var records = []Rec{}

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Funcs(fm).ParseFiles("tpl.gohtml"))
}

func main() {
	in, err := ioutil.ReadFile("table.csv")

	if err != nil {
		fmt.Println(err)
	}
	r := csv.NewReader(bytes.NewReader(in))

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if newrec, err := newRecord(record); err == nil {
			records = append(records, newrec)
		}
	}

	exerr := tpl.ExecuteTemplate(os.Stdout, "tpl.gohtml", records)
	if exerr != nil {
		log.Fatalln(exerr)
	}

}

func newRecord(r []string) (Rec, error) {
	d, err := parseDate(r[0])
	o, err := parseFloat(r[1])
	h, err := parseFloat(r[2])
	l, err := parseFloat(r[3])
	c, err := parseFloat(r[4])
	v, err := parseFloat(r[5])
	a, err := parseFloat(r[6])

	if err != nil {
		return Rec{}, err
	}

	return Rec{d, o, h, l, c, v, a}, nil
}

func parseDate(s string) (time.Time, error) {
	d, err := time.Parse("2006-01-02", s)

	if err != nil {
		return time.Now(), err
	}
	return d, nil
}

func parseFloat(s string) (float32, error) {
	f, err := strconv.ParseFloat(s, 32)

	if err != nil {
		return 0.0, err
	}

	return float32(f), nil
}

/* -- template func Map-- */

var fm = template.FuncMap{
	"fdateMDY": monthDayYear,
	"ftoInt":   floatToInt,
}

func monthDayYear(t time.Time) string {
	return t.Format("01-02-2006")
}

func floatToInt(f float32) int {
	return int(f)
}
