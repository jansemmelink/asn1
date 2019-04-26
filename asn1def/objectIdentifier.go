package asn1def

import "github.com/jansemmelink/asn1/parser"

//IObjectIdentifierValue ...
//===== Notation =====
//ObjectIdentifierValue ::=
// " { " ObjIDComponentsList " } "
//|" { " DefinedValue ObjIDComponentsList " } "
type IObjectIdentifierValue interface {
	parser.IBlock
}

//ObjectIdentifierValue ...
func ObjectIdentifierValue() IObjectIdentifierValue {
	return &objectIdentifierValue{
		IBlock: parser.Block("objectIdentifierValue", "{", "}",
			parser.Choice("objectIdentifierValueOptions",
				ObjIDComponentsList(),
				parser.Seq("objectIdentifierValue", parser.Optional(DefinedValue()), ObjIDComponentsList()),
			),
		),
	}
}

type objectIdentifierValue struct {
	parser.IBlock
}

//IObjIDComponentsList ...
//===== Notation =====
// ObjIDComponentsList ::=
//  ObjIDComponents
// |ObjIDComponents ObjIDComponentsList
type IObjIDComponentsList interface {
	parser.IList
}

//ObjIDComponentsList ...
func ObjIDComponentsList() IObjIDComponentsList {
	return &objIDComponentsList{
		IList: parser.List(" ", 1, ObjIDComponents()),
	}
}

type objIDComponentsList struct {
	parser.IList
}

//IObjIDComponents ...
//===== Notation =====
// ObjIDComponents ::=
//   NameForm
// | NumberForm
// | NameAndNumberForm
// | DefinedValue
type IObjIDComponents interface {
	parser.IChoice
}

//ObjIDComponents ...
func ObjIDComponents() IObjIDComponents {
	return &objIDComponents{
		IChoice: parser.Choice("objIDComponentsOptions",
			NameAndNumberForm(),
			NumberForm(),
			NameForm(),
			DefinedValue(),
		),
	}
}

type objIDComponents struct {
	parser.IChoice
}

//INameForm ...
//===== Notation =====
// NameForm ::= identifier
type INameForm interface {
	IIdentifier
}

//NameForm ...
func NameForm() INameForm {
	return &nameForm{
		IIdentifier: Identifier(),
	}
}

type nameForm struct {
	IIdentifier
}

//INumberForm ...
//===== Notation =====
// NumberForm ::= number | DefinedValue
type INumberForm interface {
	parser.INumber
}

//NumberForm ...
func NumberForm() INumberForm {
	return &numberForm{
		INumber: parser.Number(),
	}
}

type numberForm struct {
	parser.INumber
}

//INameAndNumberForm ...
//===== Notation =====
// NameAndNumberForm ::= identifier " ( " NumberForm " ) "
type INameAndNumberForm interface {
	parser.ISeq
}

//NameAndNumberForm ...
func NameAndNumberForm() INameAndNumberForm {
	return &nameAndNumberForm{
		ISeq: parser.Seq("NameAndNumberForm",
			Identifier(), parser.Block("identifierNumber", "(", ")", NumberForm()),
		),
	}
}

type nameAndNumberForm struct {
	parser.ISeq
}

//IDefinedValue ...
//===== Notation =====
// IDefinedValue ::=
//   ExternalValueReference
// | valuereference
// | ParameterizedValue
type IDefinedValue interface {
	parser.IChoice
}

//DefinedValue ...
func DefinedValue() IDefinedValue {
	return &definedValue{
		IChoice: parser.Choice("definedValue",
			ExternalValueReference(),
			ValueReference(),
			//todo: ParameterizedValue(),
		),
	}
}

type definedValue struct {
	parser.IChoice
}

//IExternalValueReference ...
//===== Notation =====
// ExternalValueReference ::= modulereference " . " valuereference
type IExternalValueReference interface {
	parser.ISeq
}

//ExternalValueReference ...
func ExternalValueReference() IExternalValueReference {
	return &externalValueReference{
		ISeq: parser.Seq("ExternalValueReference", ModuleReference(), parser.Keyword("."), ValueReference()),
	}
}

type externalValueReference struct {
	parser.ISeq
}

//IParameterizedValue ...
//todo...
//===== Notation =====
//The type identified by a "ParameterizedType" and "ParameterizedValueSetType", and the value identified by a "ParameterizedValue" are specified in ITU-T Rec. X.683 | ISO/IEC 8824-4.
