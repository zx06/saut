package assets_connector

import (
	"io"

	"context"
)

type AssetsConnector interface {
	io.ReadWriteCloser
	Attach(ctx context.Context) error
	WindowChange(h int, w int) error
}
