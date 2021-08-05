package gamerecorder

import (
	"context"
	"encoding/json"
	"fmt"

	v1 "github.com/clocklear/battlesnake/lib/v1"
)

type StdOutGameRecorder struct{}

func (r StdOutGameRecorder) Start(ctx context.Context, req v1.GameRequest) error {
	data, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("START: %v\n\n", string(data))
	return nil
}

func (r StdOutGameRecorder) Move(ctx context.Context, req v1.GameRequest, move string) error {
	data, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("MOVE: %v, responded with '%v'\n", string(data), move)
	return nil
}

func (r StdOutGameRecorder) End(ctx context.Context, req v1.GameRequest) error {
	data, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("END: %v\n", string(data))
	return nil
}
