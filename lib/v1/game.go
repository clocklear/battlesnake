package v1

type Game struct {
	ID      string  `json:"id"`
	Ruleset Ruleset `json:"ruleset"`
	Timeout int32   `json:"timeout"`
}

type Ruleset struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

const (
	RulesetStandard    = "standard"
	RulesetSolo        = "solo"
	RulesetRoyale      = "royale"
	RulesetSquad       = "squad"
	RulesetConstrictor = "constrictor"
	RulesetWrapped     = "wrapped"
)
