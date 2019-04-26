package asn1def

import (
	"github.com/jansemmelink/asn1/parser"
	//"github.com/prometheus/common/log"
)

//IIdentifier ...
//===== Notation =====
//11.3 Identifiers
//Same as typeReference, but starts with lowercase letter
type IIdentifier interface {
	parser.IRegex
}

var identifierNotation IIdentifier

//Identifier parser
func Identifier() IIdentifier {
	if identifierNotation == nil {
		identifierNotation = &identifier{
			IRegex: parser.Regex("identifier", `[a-z][a-zA-Z0-9-]*[a-zA-Z0-9]`),
		}
	}
	return identifierNotation
}

type identifier struct {
	parser.IRegex
}

//ITypeReference (11.2)
//consist of an arbitrary number (one or more) of letters, digits, and hyphens.
//The initial character shall be an upper-case letter.
//A hyphen shall not be the last character.
//A hyphen shall not be immediately followed by another hyphen
type ITypeReference interface {
	parser.IRegex
}

//TypeReference parser
func TypeReference() ITypeReference {
	return &identifier{
		IRegex: parser.Regex("typeReference", `[A-Z][a-zA-Z0-9-]*[a-zA-Z0-9]`),
	}
}

//IValueReference ...
//===== Notation =====
//11.4 Value Reference
//Name of lexical item – valuereference
//A "valuereference" shall consist of the sequence of characters specified for an "identifier" in 11.3.
//In analysing an instance of use of this notation, a "valuereference" is
//distinguished from an "identifier" by the context in which it appears.
type IValueReference interface {
	parser.IRegex
}

//ValueReference parser
func ValueReference() IValueReference {
	return &identifier{
		IRegex: parser.Regex("valueReference", `[a-z][a-zA-Z0-9-]*[a-zA-Z0-9]`),
	}
}

//IModuleReference (11.5)
//Name of lexical item – modulereference
//A "modulereference" shall consist of the sequence of characters specified for a "typereference" in 11.2.
//In analysing an instance of use of this notation, a "modulereference" is
//distinguished from a "typereference" by the context in which it appears.
// type IModuleReference interface {
// 	parser.IRegex
// }

var mr IIdentifier

//ModuleReference parser
func ModuleReference() IIdentifier {
	if mr == nil {
		mr = &identifier{
			IRegex: parser.Regex("moduleReference", `[A-Z][a-zA-Z0-9-]*[a-zA-Z0-9]`),
		}
		log.Debugf("mr = %p = %s", mr, mr.Name())
	}
	return mr
}

//IIdentifierList ...
type IIdentifierList interface {
	parser.IList
}

//IdentifierList ...
func IdentifierList() IIdentifierList {
	return &identifierList{
		IList: parser.List(",", 1, Identifier()),
	}
}

type identifierList struct {
	parser.IList
}
