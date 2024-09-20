# verve-assignment
Verve interview assignment

## To run the application, use either of below method

### 1. Build executable and execute
- Build project using command `go build -o ../bin` inside `/src` directory
- Execute the binary `bin/verve-assignment`

### 2. Run main method
- Run command `cd src && go run main.go`

By default, application will start in port 3002. If you want to start application in a different port, set environment variable `PORT` to the correct port.
`export PORT=<My Port>`

Application requires memcached to be up and running at localhost port 11211 (`127.0.0.1:11211` )
If you are using MAC, execute below commands to start memcached at above port.

1. `brew install memcached`
2. `memcached -d -l 127.0.0.1 -p 11211`

Bruno API collection for the same is present in `bruno collection` folder.

Note:

Extension 1 of the assignment is present in `extension-1` branch of the same repository. It will make a post call with `freeText` in the body.

Extension 2 should be implicit with shared cache across running instances of application.

Extension 3 - With ongoing release activities and hackathons in my current organization, I found very little time to work on this assignment. Hence I am unable to complete streaming service integration.