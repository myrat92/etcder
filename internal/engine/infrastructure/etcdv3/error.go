package etcdv3

import "fmt"

const (
	ErrCodeKeyNotFound int = iota + 1
)

var errCodeToMessage = map[int]string{
	ErrCodeKeyNotFound: "key not found",
}

type Error struct {
	Code               int
	Key                string
	ResourceVersion    int64
	AdditionalErrorMsg string
}

func (e Error) Error() string {
	return fmt.Sprintf("Etcd3Error: %s, Code: %d, Key: %s, ResourceVersion: %d, AdditionalErrorMsg: %s",
		errCodeToMessage[e.Code], e.Code, e.Key, e.ResourceVersion, e.AdditionalErrorMsg)
}

func NewKeyNotFoundError(key string, rv int64) *Error {
	return &Error{
		Code:            ErrCodeKeyNotFound,
		Key:             key,
		ResourceVersion: rv,
	}
}
