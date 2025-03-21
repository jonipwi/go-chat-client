module github.com/jonipwi/go-chat-client

go 1.21

require (
	github.com/jonipwi/go-chat-client/state v0.0.0
	github.com/zhouhui8915/go-socket.io-client v0.0.0-20200925034401-83ee73793ba4
)

require (
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/smarty/assertions v1.16.0 // indirect
	github.com/zhouhui8915/engine.io-go v0.0.0-20150910083302-02ea08f0971f // indirect
)

replace (
	github.com/jonipwi/go-chat-client/commands => ./commands
	github.com/jonipwi/go-chat-client/events => ./events
	github.com/jonipwi/go-chat-client/state => ./state
	github.com/jonipwi/go-chat-client/utils => ./utils
)
