package mathxf

var (
	ServerErr                    = New(-500, "internal error")
	ArgumentNotEnoughErr         = New(-501, "function '%s' requires %s arguments, but got %d")
	ArgumentNotNumberErr         = New(-502, "%s:argument '%v' not number ")
	ArgumentsOutPutErr           = New(-503, "function '%s' must have exactly 1 or 2 output arguments, the second argument must be of type error")
	ArgumentInputTypeErr         = New(-504, "function '%s' input argument %d must be of type %s or *mathxf.Value (not %T)")
	ArgumentVariadicInputTypeErr = New(-505, "function %s' variadic input argument must be of type %s or *mathxf.Value (not %T)")
	ArgumentInvalidErr           = New(-506, "function '%s' argument %d is invalid")
	ArgumentOutBoundsErr         = New(-507, "index out of bounds %s: 0-%d (index %d)")
	VariableInvalidErr           = New(-508, "variable '%s' is invalid")
	VariableNotFunctionErr       = New(-509, "variable '%s' is not a function")
	VariableNotAccessErr         = New(-510, "can't access a field by name on type %s (variable %s)")
	VariableCannotFunctionErr    = New(-511, "variable '%s' cannot be used as function")
	VariableCannotSetValueErr    = New(-512, "variable '%s' cannot be set value")
	WrapperUnclosedErr           = New(-513, "wrapper unclosed")
	UnexpectedTokenErr           = New(-514, "%s:unexpected token %v")
	MissingRightParenErr         = New(-515, "expect '%s' expected after expression")
	UnexpectedEofErr             = New(-516, "unexpected EOF")
	LexerTokenErr                = New(-517, "lexical analyzer error: %s")
	TokenNotIdentifierErr        = New(-518, "token '%s' is not an identifier")

	AssignObjectErr          = New(-519, "assign object error,Can only be 'Public'、'ResultMap'、Private objects;unexpected token %v")
	VariableAlreadyExistsErr = New(-520, "variable '%s' already exists,Cannot perform set operation")

	DivideZeroErr      = New(-522, "divide zero")
	UnknownOperatorErr = New(-523, "unknown operator %s")

	TagRegisteredErr       = New(-524, "tag '%s' is already registered")
	ConstRegisteredErr     = New(-525, "const '%s' is already exists")
	ResultKeyRegisteredErr = New(-526, "result key '%s' is already exists")
)
