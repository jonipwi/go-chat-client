module github.com/jonipwi/go-chat-client/commands

go 1.21

require (
	github.com/jonipwi/go-chat-client/state v0.0.0-00010101000000-000000000000
	github.com/jonipwi/go-chat-client/utils v0.0.0-00010101000000-000000000000
)

replace (
	github.com/jonipwi/go-chat-client/state => ../state
	github.com/jonipwi/go-chat-client/utils => ../utils
) 