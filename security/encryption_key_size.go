package security

type EncryptionKeySize int64

const (
	Bit256 EncryptionKeySize = iota
	Bit384
	Bit512
	Bit1024
	Bit2048
	Bit4096
)

func (s EncryptionKeySize) String() string {
	switch s {
	case Bit256:
		return "256bits"
	case Bit384:
		return "384bits"
	case Bit512:
		return "512bits"
	case Bit1024:
		return "1024bits"
	case Bit2048:
		return "2048bits"
	case Bit4096:
		return "4096bits"
	default:
		return "unknown"
	}
}

func (s EncryptionKeySize) FromString(key string) EncryptionKeySize {
	switch key {
	case "256", "256bits":
		return Bit256
	case "384", "384bits":
		return Bit384
	case "512", "512bits":
		return Bit512
	case "1024", "1024bits":
		return Bit1024
	case "2048", "2048bits":
		return Bit2048
	case "4096", "4069bits":
		return Bit4096
	default:
		return Bit256
	}
}
