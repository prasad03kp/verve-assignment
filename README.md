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

Extension 3 - I have implemeted writing data to kafka at every minute. Before running the application, kafka needs to be set up by running below commands

```
brew install kafka

zookeeper-server-start /opt/homebrew/opt/kafka/libexec/config/zookeeper.properties

kafka-server-start /opt/homebrew/opt/kafka/libexec/config/server.properties

kafka-topics --create --topic test --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
```

To view the topic data, run the below command

`kafka-console-consumer --bootstrap-server localhost:9092 --topic test --from-beginning`
