package domain

import "errors"

var (
    ErrNotFound           = errors.New("resource not found")
    ErrInvalidInput       = errors.New("invalid input")
    ErrUnauthorized       = errors.New("unauthorized access")
    ErrForbidden          = errors.New("forbidden access")
    ErrDocumentTooLarge   = errors.New("document size exceeds limit")
    ErrInvalidDocType     = errors.New("invalid document type")
    ErrVerificationFailed = errors.New("verification failed")
    ErrDuplicateDocument  = errors.New("duplicate document")
)

type Error struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

func NewError(code, message, details string) *Error {
    return &Error{
        Code:    code,
        Message: message,
        Details: details,
    }
}