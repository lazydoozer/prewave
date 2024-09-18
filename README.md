**Project Title**
Prewave alert term extractor

**Description**
Query prewave apis and extract specific term information from text

**Getting Started**
*Go enviornment*
1. Install the go framework https://go.dev/doc/install
2. Navigate to root of repository and run: _go run run main.go at_extractor.go at_processor.go_

*Docker*
1. Build the docker image and host application
	a. Navigate to root of repository and build image: docker build --tag prewave .
	b. Start the prewave container and expose port 8080 to port 8080 on the host.

**Output**
1. A results file(prewave_results.json) will be genenated in the root of the repo
2. A http server started on [::]:8080 will serve a /results endpoint which can be invoved via  curl http://localhost:8080/results

**Approach**
1. Retrieve the prewave query terms from the HTTP API
2. Process the retrieved query terms 
   a. Convert to lowercase for case-insensitive processing
   b. Ensure only unique terms are made available for processing
   c. if KeepOrder is not requried 
      i. Splits the wuery term's text into individual words by space
	  ii. convert to lowercase for case-insensitive processing
	  iii. Ensure only unique terms that have been split are made available for processing	  
3. Retrieve the prewave test alerts from the HTTP API
4. Perform analyis of retrieve alerts to determines in which alert a query term occurs.
	a. Process each content within an alert
	b. Converts the content’s text to lowercase to perform case-insensitive comparisons.
	c. For each term, check if the term exists as a whole word in the content.
		i. I have used regex.QuoteMeta is a function from the regexp package that escapes all special characters in a string so that the string can be used as a literal
	d. Count the number of occurances that the term matches 
5. Provide a readable summary of the analyis conducted
	a. via output file 
	b. via http by invoking curl localhost:8080/results




