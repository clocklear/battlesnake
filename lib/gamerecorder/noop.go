package gamerecorder

import (
	"context"

	v1 "github.com/clocklear/battlesnake/lib/v1"
)

type NoopGameRecorder struct{}

func (r NoopGameRecorder) Start(context.Context, v1.GameRequest) error {
	return nil
}

func (r NoopGameRecorder) Move(context.Context, v1.GameRequest, string) error {
	return nil
}

func (r NoopGameRecorder) End(context.Context, v1.GameRequest) error {
	return nil
}
