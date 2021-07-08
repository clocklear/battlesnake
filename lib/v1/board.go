package v1

type Board struct {
	Height  int           `json:"height"`
	Width   int           `json:"width"`
	Food    CoordList     `json:"food"`
	Hazards CoordList     `json:"hazards"`
	Snakes  []Battlesnake `json:"snakes"`
}
