package main

import (
	"time"

	"google.golang.org/protobuf/types/known/structpb"
)

// Event representa o evento de log simplificado.
type Event struct {
	ID          string                     `json:"id"`
	EventTime   time.Time                  `json:"eventTime"`
	EventSource EventSource                `json:"eventSource"`
	Tags        map[string]string          `json:"tags"`
	Context     map[string]*structpb.Value `json:"context"`
	User        User                       `json:"user"`
	Session     Session                    `json:"session"`
	Request     map[string]*structpb.Value `json:"request"`
	Response    map[string]*structpb.Value `json:"response"`
}

// EventSource representa a origem do evento.
type EventSource struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Action string `json:"action"`
}

// User representa o usuário associado ao evento.
type User struct {
	ID    string `json:"id"`
	Login string `json:"login"`
}

// Session representa a sessão associada ao evento.
type Session struct {
	ID        string    `json:"id"`
	StartTime time.Time `json:"startTime"`
}
