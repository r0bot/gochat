## Quickstart
The applications requires Go version 1.11+

All commands below must be executed from the root directory of the project

Start the server
```$xslt
    go run /cmd/server/main.go
```

Start a client
```$xslt
    go run /cmd/client/main.go
```

Run tests
```$xslt
    go test  ./...
```

## Client usage

The client is demonised and once started will continuously listen for user input.
Typing something in the stdin and pressing enter will send that as a message to the server.
The Client also supports commands using the "\\" character. For example typing `\exit` will quit the program.

#### Commands

`\exit` Will exit the process.

`\pm <clientId> <message>` will send a private message to that client only (clientId can be taken from the server log)

Simply typing something in the stdin and hitting enter will broadcast it as a message to all clients on the server.