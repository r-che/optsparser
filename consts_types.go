package optsparser

const (
	typeBool		=	"bool"
	typeString		=	"string"
	typeInt			=	"int"
	typeInt64		=	"int64"
	typeUint		=	"uint"
	typeUint64		=	"uint64"
	typeFloat64		=	"float64"
	typeDuration	=	"duration"
	typeVal			=	"value"
)

type optDescr struct {
	optType		string
	short		string
}
