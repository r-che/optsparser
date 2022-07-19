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
	// TODO Add required flag
	optType		string
	short		string
}

type dummyWriter struct {}
func (dw *dummyWriter) Write(p []byte) (int, error) {
	return len(p), nil
}
