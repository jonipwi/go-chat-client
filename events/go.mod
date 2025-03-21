module github.com/jonipwi/go-chat-client/events

go 1.21

require (
	github.com/jonipwi/go-chat-client/state v0.0.0
	github.com/zhouhui8915/go-socket.io-client v0.0.0-20200925034401-83ee73793ba4
)

require (
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/gomodule/redigo v1.9.2 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	golang.org/x/net v0.17.0 // indirect
)

replace github.com/jonipwi/go-chat-client/state => ../state 