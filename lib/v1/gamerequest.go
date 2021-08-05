package v1

type GameRequest struct {
	Game  Game        `json:"game"`
	Turn  int         `json:"turn"`
	Board Board       `json:"board"`
	You   Battlesnake `json:"you"`
}

type BoardState struct {
	Turn  int         `json:"turn"`
	Board Board       `json:"board"`
	You   Battlesnake `json:"you"`
}

func (gr GameRequest) ToBoardState() BoardState {
	return BoardState{
		Turn:  gr.Turn,
		Board: gr.Board,
		You:   gr.You,
	}
}
