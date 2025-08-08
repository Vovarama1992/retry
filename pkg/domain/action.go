package domain

import "time"

// Action — доменная сущность действия (в том числе визит)
type Action struct {
	ID             int64
	ActionTypeID   int
	ActionTypeName string // опционально, может быть пустым
	VisitID        string
	Source         string
	IPAddress      string
	Timestamp      time.Time
}

// ActionType — справочник типов действий
type ActionType struct {
	ID   int
	Name string
}
