package main

import (
	"fmt"
	"html/template"
	"os"
	"strconv"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/gobuffalo/packr/v2"
	"github.com/gocarina/gocsv"
)

type Stat struct {
	Type          string  `csv:"Type"`
	Name          string  `csv:"Name"`
	CountReq      int     `csv:"Request Count"`
	CountFail     int     `csv:"Failure Count"`
	RespMedian    float64 `csv:"Median Response Time"`
	RespAvg       float64 `csv:"Average Response Time"`
	RespMin       float64 `csv:"Min Response Time"`
	RespMax       float64 `csv:"Max Response Time"`
	RateReq       float64 `csv:"Requests/s"`
	RateFail      float64 `csv:"Failures/s"`
	Percentile50  float64 `csv:"50%"`
	Percentile75  float64 `csv:"75%"`
	Percentile90  float64 `csv:"90%"`
	Percentile95  float64 `csv:"95%"`
	Percentile99  float64 `csv:"99%"`
	Percentile100 float64 `csv:"100%"`
}

type HistoryRow struct {
	Timestamp    string `csv:"Timestamp"`
	TimeFormated string `csv:"-"`
	UserCount    string `csv:"User Count"`
	Type         string `csv:"Type"`
	Name         string `csv:"Name"`
	RateReq      string `csv:"Requests/s"`
	RateFail     string `csv:"Failures/s"`
	RespAvg      string `csv:"Total Average Response Time"`
}

// TemplateData is our main data struct we populate from the CSVs
type TemplateData struct {
	Title           string
	Stats           []*Stat
	AggregatedStats Stat
	HistoryData     map[string][]*HistoryRow
}

func main() {
	fmt.Println("\n\033[36m╔════════════════════════════════════════════════╗")
	fmt.Println("║    \033[33m🦗 Locust HTML Report Converter 📜\033[36m   \033[35mv1.0   \033[36m║")
	fmt.Println("╚════════════════════════════════════════════════╝\033[0m")

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

	// Check args
	if len(os.Args) < 4 {
		fmt.Printf("\n💬 \033[31mERROR! Must supply three args, input directory, CSV prefix and output HTML file\n")
		os.Exit(1)
	}

	// Open history CSV
	historyFile, err := os.Open(fmt.Sprintf("%s/%s_stats_history.csv", os.Args[1], os.Args[2]))
	if err != nil {
		fmt.Println("💥 Input history file error", err)
		os.Exit(1)
	}
	defer historyFile.Close()

	// Open stats CSV
	statsFile, err := os.Open(fmt.Sprintf("%s/%s_stats.csv", os.Args[1], os.Args[2]))
	if err != nil {
		fmt.Println("💥 Input stats file error", err)
		os.Exit(1)
	}
	defer statsFile.Close()

	stats := []*Stat{}
	if err := gocsv.UnmarshalFile(statsFile, &stats); err != nil { // Load clients from file
		fmt.Println("💥 Stats CSV marshalling error", err)
		os.Exit(1)
	}
	historyAllRows := []*HistoryRow{}
	if err := gocsv.UnmarshalFile(historyFile, &historyAllRows); err != nil { // Load clients from file
		fmt.Println("💥 History CSV marshalling error", err)
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
		if key == " Aggregated" {
			key = "Aggregated"
		}
		histMap[key] = append(histMap[key], row)
	}

	// Open output HTML file
	outFile, err := os.Create(os.Args[3])
	if err != nil {
		fmt.Println("💥 Output file error", err)
		os.Exit(1)
	}

	// Build template data for passing to template
	templateData := TemplateData{}
	templateData.Stats = stats
	templateData.HistoryData = histMap
	templateData.Title = os.Args[2]

	// Pull out aggregated to make template easier
	for _, stat := range stats {
		if stat.Name == "Aggregated" {
			templateData.AggregatedStats = *stat
		}
	}

	fmt.Printf("\n📜 Done! Output HTML written to: %s\n", outFile.Name())

	// Render template into output fine, and that's it
	tmpl.Execute(outFile, templateData)
}
