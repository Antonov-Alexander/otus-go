package types

type Request struct {
	IP       int
	Login    string
	Password string
}

type Item interface{}
type Limit struct {
	Interval int
	Limit    int
}

type Config interface {
	Init(checkTypes []int) error
	GetCheckConfig(checkName int) (CheckConfig, bool)
}

type CheckConfig struct {
	CommonLimits []Limit
	ItemLimits   map[Item][]Limit
	BlackList    map[Item]struct{}
	WhiteList    map[Item]struct{}
}

type Check interface {
	Init(config CheckConfig, storage Storage) error
	GetItem(request Request) Item
	GetDefaultConfig() CheckConfig
	Check(request Request) error
	ClearCounter(Item)
	AddWhiteListItem(Item)
	AddBlackListItem(Item)
	RemoveWhiteListItem(Item)
	RemoveBlackListItem(Item)
}

type Storage interface {
	Inc(item Item, limit Limit) bool
	Reset(item Item)
}
