package types

type Pokemon struct {
	ID             int        `json:"id"`
	BaseExperience int        `json:"base_experience"`
	Moves          []Moves    `json:"moves"`
	Name           string     `json:"name"`
	Species        BaseStruct `json:"species"`
	Sprites        Sprites    `json:"sprites"`
	Stats          []Stats    `json:"stats"`
	Types          []Types    `json:"types"`
}

type Moves struct {
	Move         BaseStruct            `json:"move"`
	VersionGroup []VersionGroupDetails `json:"version_group_details"`
}

type VersionGroupDetails struct {
	LevelLearnedAt int        `json:"level_learned_at"`
	LearnMethod    BaseStruct `json:"move_learn_method"`
	VersionGroup   BaseStruct `json:"version_group"`
}

type Sprites struct {
	BackUrl  string `json:"back_default"`
	FrontUrl string `json:"front_default"`
}

type Stats struct {
	BaseStat int        `json:"base_stat"`
	Effort   int        `json:"effort"`
	Stat     BaseStruct `json:"stat"`
}

type Types struct {
	Slot int        `json:"slot"`
	Type BaseStruct `json:"type"`
}

type BaseStruct struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
