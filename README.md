# chat
Make chat application with Golang based on Go Programming Blueprint Book By Mat Ryer with some modified oauth2 library and using Godeps dependency

## Requirements

[Go 1.7](https://golang.org/dl/)
Set your $GOPATH


Install [Godep](https://github.com/tools/godep) dependency tools,
```
go get github.com/tools/godep
cd <your/chat/path>
godep install
go build -o chat
```

To RUN the apps
```
./chat --addr=":3000"
```

Open Your Browser and login with your Google Account
```
http://localhost:3000
```
