build :
	packr2
	go build -o bin/locust-reporter
	packr2 clean	

clean :
	rm -rf bin
	rm -rf *.csv
	rm -rf *.html
