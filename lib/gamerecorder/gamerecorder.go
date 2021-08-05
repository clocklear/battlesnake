package gamerecorder

import (
	"context"

	v1 "github.com/clocklear/battlesnake/lib/v1"
)

type GameRecorder interface {
	Start(context.Context, v1.GameRequest) error
	Move(context.Context, v1.GameRequest, string) error
	End(context.Context, v1.GameRequest) error
}
