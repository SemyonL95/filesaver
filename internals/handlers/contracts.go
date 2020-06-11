package handlers

import (
	"context"
	"io"
)

type FileStorage interface {
	OpenEntry(string) (io.ReadCloser, error)
	CreateEntry(string, io.Reader) (io.WriteCloser, error)
}

type Cache interface {
	Exists(context.Context, string, string) (bool, error)
	Set(context.Context, string, string) (bool, error)
	Remove(context.Context, string) (bool, error)
}
