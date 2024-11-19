package server

type Command string

const (
	Ping     Command = "ping"
	Echo     Command = "echo"
	Set      Command = "set"
	Get      Command = "get"
	Info     Command = "info"
	ReplConf Command = "replconf"
	PSync    Command = "psync"
)
