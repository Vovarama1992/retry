package domain

import (
	"encoding/json"
	"time"
)

// Action — доменная сущность действия (в том числе визит)
type Action struct {
	ID             int64
	ActionTypeID   int
	ActionTypeName string // опционально, может быть пустым
	VisitID        string
	SessionID      string
	Source         string
	IPAddress      string
	Timestamp      time.Time
	Meta           json.RawMessage // хранит meta как JSON, можно парсить в map при необходимости
}

// ActionType — справочник типов действий
type ActionType struct {
	ID   int
	Name string
}
