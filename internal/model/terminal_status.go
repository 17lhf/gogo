package model

// TerminalStatus represents the status of a terminal.
type TerminalStatus string

const (
	TerminalStatusOffline  TerminalStatus = "offline"
	TerminalStatusOnline   TerminalStatus = "online"
	TerminalStatusDisabled TerminalStatus = "disabled"
	TerminalStatusEnabled  TerminalStatus = "enabled"
)

func (s TerminalStatus) String() string {
	return string(s)
}
