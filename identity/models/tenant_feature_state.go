package models

type TenantFeatureState int64

const (
	Unknown TenantFeatureState = iota
	Provisioning
	Provisioned
	Enabling
	Enabled
	Disabling
	Disabled
	Deleting
	Deleted
	Error
)

func (l TenantFeatureState) String() string {
	switch l {
	case 0:
		return "unknown"
	case 1:
		return "provisioning"
	case 2:
		return "provisioned"
	case 3:
		return "enabling"
	case 4:
		return "enabled"
	case 5:
		return "disabling"
	case 6:
		return "disabled"
	case 7:
		return "deleting"
	case 8:
		return "deleted"
	case 9:
		return "error"
	default:
		return "unknown"
	}
}

func (l TenantFeatureState) FromString(value string) TenantFeatureState {
	switch value {
	case "unknown":
		return Unknown
	case "provisioning":
		return Provisioning
	case "provisioned":
		return Provisioned
	case "enabling":
		return Enabling
	case "enabled":
		return Enabled
	case "disabling":
		return Disabling
	case "disabled":
		return Disabled
	case "deleting":
		return Deleting
	case "deleted":
		return Deleted
	default:
		return Unknown
	}
}
