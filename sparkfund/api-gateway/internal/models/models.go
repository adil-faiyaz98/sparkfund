package models

import (
	"net/http"
	"time"
)

type ServiceInstance struct {
	ID       string
	Address  string
	Healthy  bool
	LastSeen int64
}

type RequestContext struct {
	UserID    string
	ServiceID string
	Path      string
	Method    string
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Body    []byte      `json:"-"`
	Headers http.Header `json:"-"`
}

type Request struct {
	Method  string
	Path    string
	Headers http.Header
	Body    []byte
}

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
