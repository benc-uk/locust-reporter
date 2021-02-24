# Locust HTML Report Converter

Simple Go command line application to convert the results of a Locust load test into a HTML report.  
[Locust is a developer focused load testing tool](https://locust.io/)

The report will show three main sets of data:

- Summary statistics for all requests
- History data, with a range of metrics charted over the lifetime of the load test
- List of failures by request

This project uses Go templates, [Sprig](http://masterminds.github.io/sprig/) and [Packr](https://github.com/gobuffalo/packr)

![](https://img.shields.io/github/license/benc-uk/locust-reporter)
![](https://img.shields.io/github/last-commit/benc-uk/locust-reporter)
![](https://img.shields.io/github/release/benc-uk/locust-reporter)
![](https://img.shields.io/github/checks-status/benc-uk/locust-reporter/main)

# Getting Started

## Installing

Install with `go get`

```bash
go get -u github.com/benc-uk/locust-reporter/cmd/locust-reporter
```

Or download a pre-compiled binary (Linux x64)

```bash
wget https://github.com/benc-uk/locust-reporter/releases/download/v1.2.2/locust-reporter
```

## Usage

The command takes the following arguments:

- `-dir` - The input directory, which contains the three Locust CSV files. Default is current directory.
- `-prefix` - The prefix of the CSV files, **required parameter**
- `-outfile` - The output HTML file, which will be created or overwritten. Default is `out.html`
- `-failures` - Include log of failures in the report, can result in very large output. Default is `false`

Example

```bash
./locust-reporter -help

╔════════════════════════════════════════════════╗
║   🦗 Locust HTML Report Converter 📜  v1.2.2  ║
╚════════════════════════════════════════════════╝

  -dir string
        Directory holding input Locust CSV files (default ".")
  -failures
        Include failures in report, can result in very large output
  -outfile string
        Output HTML filename (default "./out.html")
  -prefix string
        Prefix for CSV files, required
```

# Building & Running

Makefile reference

```txt
build                Build binary executable, into bin directory
clean                Clean up
help                 This help message :)
lint-fix             Lint & format, will try to fix errors and modify code
lint                 Lint & format, will not fix but sets exit code on error
run                  Run locally
```

# Screenshots

![image](https://user-images.githubusercontent.com/14982936/105252609-ac014e00-5b75-11eb-9e20-b97eb30208ee.png)

![image](https://user-images.githubusercontent.com/14982936/105252615-b02d6b80-5b75-11eb-9e82-573d8f329519.png)

# Known Issues

A very long test with high number of data points may take some processing to display

# Change Log

See [complete change log](./CHANGELOG.md)

# License

This project uses the MIT software license. See [full license file](./LICENSE)
