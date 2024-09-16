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
