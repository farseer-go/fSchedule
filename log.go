package fSchedule

type Enum int

const (
	Trace Enum = iota
	Debug
	Information
	Warning
	Error
	Critical
	NoneLevel
)
