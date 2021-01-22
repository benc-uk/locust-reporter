package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"html/template"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/gobuffalo/packr/v2"
	"github.com/gocarina/gocsv"
)

// Stat holds a row from the stats CSV
type Stat struct {
	Type          string  `csv:"Type"`
	Name          string  `csv:"Name"`
	CountReq      int     `csv:"Request Count"`
	CountFail     int     `csv:"Failure Count"`
	RespMedian    float64 `csv:"Median Response Time"`
	RespAvg       float64 `csv:"Average Response Time"`
	RespMin       float64 `csv:"Min Response Time"`
	RespMax       float64 `csv:"Max Response Time"`
	SizeAvg       float64 `csv:"Average Content Size"`
	RateReq       float64 `csv:"Requests/s"`
	RateFail      float64 `csv:"Failures/s"`
	Percentile50  float64 `csv:"50%"`
	Percentile75  float64 `csv:"75%"`
	Percentile90  float64 `csv:"90%"`
	Percentile95  float64 `csv:"95%"`
	Percentile99  float64 `csv:"99%"`
	Percentile100 float64 `csv:"100%"`
}

// Failure holds a record from failures CSV
type Failure struct {
	Method      string `csv:"Method"`
	Name        string `csv:"Name"`
	Error       string `csv:"Error"`
	Occurrences int    `csv:"Occurrences"`
}

// HistoryRow holds a row from the stats CSV
type HistoryRow struct {
	Timestamp     string  `csv:"Timestamp"`
	TimeFormated  string  `csv:"-"`
	CountUser     int     `csv:"User Count"`
	Type          string  `csv:"Type"`
	Name          string  `csv:"Name"`
	RateReq       float64 `csv:"Requests/s"`
	RateFail      float64 `csv:"Failures/s"`
	Percentile50  float64 `csv:"50%"`
	Percentile75  float64 `csv:"75%"`
	Percentile90  float64 `csv:"90%"`
	Percentile95  float64 `csv:"95%"`
	Percentile99  float64 `csv:"99%"`
	Percentile100 float64 `csv:"100%"`
	CountReq      float64 `csv:"Total Request Count"`
	CountFail     float64 `csv:"Total Failure Count"`
	RespMedian    float64 `csv:"Total Median Response Time"`
	RespAvg       float64 `csv:"Total Average Response Time"`
	RespMin       float64 `csv:"Total Min Response Time"`
	RespMax       float64 `csv:"Total Max Response Time"`
	SizeAvg       float64 `csv:"Total Average Response Time"`
}

// TemplateData is our main data struct we populate from the CSVs
type TemplateData struct {
	Title           string
	Stats           []*Stat
	Failures        []*Failure
	AggregatedStats Stat
	HistoryData     map[string][]*HistoryRow
}

func main() {
	fmt.Println("\n\033[36m╔════════════════════════════════════════════════╗")
	fmt.Println("║    \033[33m🦗 Locust HTML Report Converter 📜\033[36m   \033[35mv1.0   \033[36m║")
	fmt.Println("╚════════════════════════════════════════════════╝\033[0m")

	var inDir = flag.String("dir", ".", "Directory holding input Locust CSV files")
	var csvPrefix = flag.String("prefix", "", "Prefix for CSV files, required")
	var outFilename = flag.String("outfile", "./out.html", "Output HTML filename")
	flag.Parse()
	if *csvPrefix == "" {
		fmt.Printf("\n🚫 CSV prefix not specified, please add -prefix\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// We package the template into the binary using packr
	// Hopefully when Go 1.16 is release we can use an embed
	templateBox := packr.New("Templates Box", "./templates")
	templateString, err := templateBox.FindString("report.tmpl")
	if err != nil {
		fmt.Println("💥 Template unboxing error", err)
		os.Exit(1)
	}
	tmpl, err := template.New("").Funcs(sprig.FuncMap()).Parse(templateString)
	if err != nil {
		fmt.Println("💥 Template file error", err)
		os.Exit(1)
	}

	// Open CSV files
	historyFile, err := os.Open(fmt.Sprintf("%s/%s_stats_history.csv", *inDir, *csvPrefix))
	if err != nil {
		fmt.Println("📃 History CSV file error", err)
		os.Exit(1)
	}
	defer historyFile.Close()

	statsFile, err := os.Open(fmt.Sprintf("%s/%s_stats.csv", *inDir, *csvPrefix))
	if err != nil {
		fmt.Println("📃 Stats CSV file error", err)
		os.Exit(1)
	}
	defer statsFile.Close()

	failureFile, err := os.Open(fmt.Sprintf("%s/%s_failures.csv", *inDir, *csvPrefix))
	if err != nil {
		fmt.Println("📃 Failure CSV file error", err)
		os.Exit(1)
	}
	defer statsFile.Close()

	// Marshall CSVs into memory
	stats := []*Stat{}
	if err := gocsv.UnmarshalFile(statsFile, &stats); err != nil {
		fmt.Println("🚽 Stats CSV marshalling error", err)
		os.Exit(1)
	}

	historyAllRows := []*HistoryRow{}
	gocsv.UnmarshalFileWithErrorHandler(historyFile, func(err *csv.ParseError) bool {
		// We need to deal with any values which are "N/A" which can be in the CSV
		// These would normally trigger an error, we can safely ignore and value will remain as zero
		if strings.Contains(err.Err.Error(), "parsing \"N/A\"") {
			return true
		}
		return false
	}, &historyAllRows)

	failures := []*Failure{}
	if err := gocsv.UnmarshalFile(failureFile, &failures); err != nil {
		fmt.Println("🚽 Failure CSV marshalling error", err)
		os.Exit(1)
	}

	// Make a map of history rows, keyed on combo of request type and id
	histMap := make(map[string][]*HistoryRow)
	for _, row := range historyAllRows {
		// Format timestamps, easier to do this here than in the template
		timeInt, err := strconv.ParseInt(row.Timestamp, 10, 64)
		if err != nil {
			continue
		}
		tm := time.Unix(timeInt, 0)
		row.TimeFormated = tm.Format("15:04:05")

		// Build a key from type + name (as name isn't unique)
		key := fmt.Sprintf("%s %s", row.Type, row.Name)
		// Aggregated is a special case, it never has a type
		if key == " Aggregated" {
			key = "Aggregated"
		}

		// Push row into correct array in the map
		histMap[key] = append(histMap[key], row)
	}

	// Open output HTML file
	outFile, err := os.Create(*outFilename)
	if err != nil {
		fmt.Println("💥 Output file error", err)
		os.Exit(1)
	}

	// Build template data struct for passing to template
	templateData := TemplateData{}
	templateData.Stats = stats
	templateData.Failures = failures
	templateData.HistoryData = histMap
	templateData.Title = *csvPrefix

	// Pull out aggregated stats to make template easier
	for _, stat := range stats {
		if stat.Name == "Aggregated" {
			templateData.AggregatedStats = *stat
		}
	}

	// Render template into output file
	tmpl.Execute(outFile, templateData)
	fmt.Printf("\n📜 Done! Output HTML written to: %s\n", outFile.Name())
}
