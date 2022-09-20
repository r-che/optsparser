package optsparser

const sepPrefix = "\u0000\u0000separator\u0000\u0000"
const optIndent = "    "
const helpIndent = optIndent + "  "

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
	typeSeparator	=	"sep"
)

type optDescr struct {
	optType		string
	short		string
}
