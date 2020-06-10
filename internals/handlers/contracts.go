package handlers

import (
	"context"
	"io"
)

type FileStorage interface {
	Get(string) (<-chan string, error)
	Put(string, io.Reader) error
}

type Cache interface {
	Exists(context.Context, string, string) (bool, error)
	Set(context.Context, string, string) (bool, error)
	Remove(context.Context, string) (bool, error)
}
