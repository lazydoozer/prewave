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
   &emsp; a. Navigate to root of repository and build image: _docker build --tag prewave .  _  
   &emsp; b. Start the prewave container and expose port 8080 to port 8080 on the host: _docker run --publish 8080:8080 prewave_  

**Output**
1. A results file(prewave_results.json) will be genenated in the root of the repo
2. A http server started on [::]:8080 will serve a /results endpoint which can be invoved via  curl http://localhost:8080/results

**Approach**
1. Retrieve the prewave query terms from the HTTP API
2. Process the retrieved query terms   
   &emsp; a. Convert to lowercase for case-insensitive processing  
   &emsp; b. Ensure only unique terms are made available for processing    
   &emsp; c. if KeepOrder is not requried   
      &emsp; &emsp; &emsp; i. Split the query term's text into individual words by space  
      &emsp; &emsp; &emsp; ii. convert to lowercase for case-insensitive processing  
      &emsp; &emsp; &emsp; iii. Ensure only unique terms that have been split are made available for processing  	  
3. Retrieve the prewave test alerts from the HTTP API  
4. Perform analyis of retrieve alerts to determines in which alert a query term occurs  
   &emsp; a. Process each content within an alert  
   &emsp; b. Converts the content’s text to lowercase to perform case-insensitive comparisons  
   &emsp; c. For each term, check if the term exists as a whole word in the content    
      &emsp; &emsp; &emsp; i. regex.QuoteMeta which is a function that escapes all special characters in a string so that the string can be used as a literal  
   &emsp; d. Count the number of occurances that the term matches 
6. Provide a readable summary of the analyis conducted  
   &emsp; a. via output file at root of repo: prewave_results.json  
   &emsp; b. via http by invoking curl localhost:8080/results  




