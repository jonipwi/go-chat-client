module github.com/jonipwi/go-chat-client

go 1.21

require (
	github.com/gorilla/websocket v1.5.3
	github.com/jonipwi/go-chat-client/commands v0.0.0
	github.com/jonipwi/go-chat-client/events v0.0.0
	github.com/jonipwi/go-chat-client/state v0.0.0
	github.com/jonipwi/go-chat-client/utils v0.0.0
	github.com/stretchr/testify v1.10.0
	github.com/zhouhui8915/go-socket.io-client v0.0.0-20200925034401-83ee73793ba4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/zhouhui8915/engine.io-go v0.0.0-20150910083302-02ea08f0971f // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/jonipwi/go-chat-client/commands => ./commands
	github.com/jonipwi/go-chat-client/events => ./events
	github.com/jonipwi/go-chat-client/state => ./state
	github.com/jonipwi/go-chat-client/utils => ./utils
)
