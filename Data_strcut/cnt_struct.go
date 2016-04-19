package Data_strcut

import (
	"time"
)

type Cnt_analysis struct {
	Src       string
	Timestamp time.Time
	Action    string
	Repo      string
	User      string
}

type Cnt_repo struct {
	User string
	Repo string
	Tag  string
}

type Cnt_user struct {
	Username string
	Password string
}

type ACLEntry struct {
	Match MatchConditions `yaml:"match"`
}

type MatchConditions struct {
	Account string `yaml:"account,omitempty" json:"account,omitempty"`
	Name    string `yaml:"name,omitempty" json:"name,omitempty"`
}
