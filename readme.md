# Locust HTML Report Converter

Simple Go command line application to convert the results of a Locust load test into a HTML report.  
[Locust is a developer focused load testing tool](https://locust.io/)

The report will show all request summary statistics, history data as charts and failures

# Usage

The command takes three positional arguments:

- The input directory, which contains the three Locust CSV files.
- The prefix of the CSV files, can be an empty string.
- The output HTML file, which will be created or overwritten.

Example

```bash
./locust-reporter . "test" ./report.html
```

# Building & Running

Packr v2 is required to build a standalone binary

```bash
go get -u github.com/gobuffalo/packr/v2/packr2
```

Then run `./scripts/build.sh`

Alternatively run in place with `go run main.go`

# Screenshots

![image](https://user-images.githubusercontent.com/14982936/105252609-ac014e00-5b75-11eb-9e20-b97eb30208ee.png)

![image](https://user-images.githubusercontent.com/14982936/105252615-b02d6b80-5b75-11eb-9e82-573d8f329519.png)
