package wage

import "io"

type Module interface {
	Open(input any) error
	io.Closer
}

type Export = func() Module
