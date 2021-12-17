package mongodb

type MongoOperator int64

const (
	NONE MongoOperator = iota
	AND
	NOR
	OR
	EQUALS
	GREATERTHAN
	LOWERTHAN
)

func (o MongoOperator) String() string {
	switch o {
	case AND:
		return "$and"
	case NOR:
		return "$nor"
	case OR:
		return "$or"
	case NONE:
		return "$none"
	case EQUALS:
		return "$eq"
	case GREATERTHAN:
		return "$gt"
	case LOWERTHAN:
		return "$lt"
	default:
		return "$and"
	}
}
