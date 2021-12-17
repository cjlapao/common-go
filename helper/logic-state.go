package helper

type LogicState int64

const (
	Unknown LogicState = iota
	Enabled
	Disabled
)

func (l LogicState) String() string {
	switch l {
	case 0:
		return "unknown"
	case 1:
		return "enabled"
	case 2:
		return "disabled"
	default:
		return "unknown"
	}
}

func (l LogicState) FromString(value string) LogicState {
	switch value {
	case "unknown":
		return Unknown
	case "enabled":
		return Enabled
	case "disabled":
		return Disabled
	default:
		return Disabled
	}
}
