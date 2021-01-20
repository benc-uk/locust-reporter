package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/Masterminds/sprig/v3"
	"github.com/gobuffalo/packr/v2"
	"github.com/gocarina/gocsv"
)

type Stat struct {
	Type        string  `csv:"Type"`
	Name        string  `csv:"Name"`
	CountReq    int     `csv:"Request Count"`
	CountFail   int     `csv:"Failure Count"`
	RespMedian  float64 `csv:"Median Response Time"`
	RespAvg     float64 `csv:"Average Response Time"`
	RespMin     float64 `csv:"Min Response Time"`
	RespMax     float64 `csv:"Max Response Time"`
	RateRequest float64 `csv:"Requests/s"`
	RateFail    float64 `csv:"Failures/s"`
}

// TemplateData is our main data struct we populate from the CSVs
type TemplateData struct {
	Title string
	Stats []*Stat
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

	// Open output HTML file
	outFile, err := os.Create(os.Args[3])
	if err != nil {
		fmt.Println("💥 Output file error", err)
		os.Exit(1)
	}

	templateData := TemplateData{}
	templateData.Stats = stats
	templateData.Title = os.Args[2]

	//spew.Dump(templateData)

	fmt.Printf("\n📜 Done! Output HTML written to: %s\n", outFile.Name())
	// Render template into output fine, and that's it
	tmpl.Execute(outFile, templateData)
}
