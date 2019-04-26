package asn1def

import (
	"bitbucket.org/vservices/dark/logger"
	"github.com/jansemmelink/asn1/parser"
)

//IModuleDefinition ...
//----------------------------------------------------------------------------
//===== Notation =====
// ModuleDefinition ::=
// 	ModuleIdentifier
// 	DEFINITIONS
// 	TagDefault
// 	ExtensionDefault
// 	" ::= "
// 	BEGIN
// 	ModuleBody
// 	END
//
//===== Example =====
// MAP-DialogueInformation {itu-t identified-organization (4) etsi (0) mobileDomain (0) gsm-Network (1) modules (3) map-DialogueInformation (3) version8 (8)}
// DEFINITIONS
// IMPLICIT TAGS
// ::=
// BEGIN
// ...
// END
//----------------------------------------------------------------------------
type IModuleDefinition interface {
	parser.ISeq
	//ModuleIdentifier() INotation //IModuleIdentifier
	TagDefault() ITagDefault
	ExtensionDefault() IExtensionDefault
	ModuleBody() IModuleBody
}

//ModuleDefinition ...
func ModuleDefinition() IModuleDefinition {
	return &moduleDefinition{
		ISeq: parser.Seq("moduleDefinition",
			ModuleIdentifier,
			// todo: parser.Keyword("DEFINITIONS"),
			// todo: TagDefault(),
			// todo: ExtensionDefault(),
			// todo: parser.Keyword("::="),
			// todo: parser.Block("moduleDefinitionBody", "BEGIN", "END", ModuleBody()),
		),
	}
}

type moduleDefinition struct {
	parser.ISeq
	//mi IModuleIdentifier
	td ITagDefault
	ed IExtensionDefault
	mb IModuleBody
}

// func (md moduleDefinition) ModuleIdentifier() IModuleIdentifier {
// 	return md.mi
// }

func (md moduleDefinition) TagDefault() ITagDefault {
	return md.td
}

func (md moduleDefinition) ExtensionDefault() IExtensionDefault {
	return md.ed
}

func (md moduleDefinition) ModuleBody() IModuleBody {
	return md.mb
}

func (md moduleDefinition) ParseV(log logger.ILogger, l parser.ILines) (parser.IValue, parser.ILines, error) {
	parsedValue, remain, err := md.ISeq.ParseV(log, l)
	if err != nil {
		return parsedValue, remain, log.Wrapf(err, "Failed")
	}

	log.Debugf("MD %T=%+v", parsedValue, parsedValue)
	return parsedValue, remain, nil
}

//ModuleIdentifier ...
//----------------------------------------------------------------------------
//===== Notation =====
// ModuleIdentifier ::=
// 	modulereference
// 	DefinitiveIdentifier
//
// DefinitiveIdentifier ::=
// 	" { " DefinitiveObjIdComponentList " } "
// 	|
// 	empty
//
//===== Example =====
// MAP-DialogueInformation {
// 	itu-t identified-organization (4) etsi (0) mobileDomain (0)
// 	gsm-Network (1) modules (3) map-DialogueInformation (3) version8 (8)}
//
//----------------------------------------------------------------------------
type moduleIdentifier struct {
	parser.ISeq
}

//ModuleIdentifier ...
var ModuleIdentifier *moduleIdentifier

func init() {
	ModuleIdentifier = &moduleIdentifier{
		ISeq: parser.Seq("moduleIdentifier",
			ModuleReference(),
			DefinitiveIdentifier(),
		),
	}
} //init()

func (mi moduleIdentifier) ParseV(log logger.ILogger, l parser.ILines) (parser.IValue, parser.ILines, error) {
	parsedValue, remain, err := mi.ISeq.ParseV(log, l)
	if err != nil {
		return nil, l, log.Wrapf(err, "Failed")
	}
	//log values:
	log.Debugf("MI Value:")
	for _, i := range parsedValue.Items() {
		log.Debugf("  %s: %+v", i.Name(), i.Value())
	}
	return parsedValue, remain, nil
} //moduleIdenfier.ParseV()

//IDefinitiveIdentifier ::= " { " DefinitiveObjIdComponentList " } " | empty
// i.e. list of space-separated name|number|name and number enclosed in braces
type IDefinitiveIdentifier interface {
	parser.IBlock
	//todo: method to get array of values
}

//DefinitiveIdentifier parser
func DefinitiveIdentifier() IDefinitiveIdentifier {
	return &definitiveIdentifier{
		IBlock: parser.Block("defIdBlock", "{", "}",
			parser.List(" ", 0, DefinitiveObjIDComponent()),
		),
	}
}

type definitiveIdentifier struct {
	parser.IBlock
}

//IDefinitiveObjIDComponent ::=
//		NameForm
//	|	DefinitiveNumberForm
//	|	DefinitiveNameAndNumberForm
type IDefinitiveObjIDComponent interface {
	parser.IChoice
}

//DefinitiveObjIDComponent parser
func DefinitiveObjIDComponent() IDefinitiveObjIDComponent {
	doidc := &definitiveObjIDComponent{
		IChoice: parser.Choice("definitiveObjIDComponent"),
	}
	doidc.IChoice.Add(parser.Number())
	doidc.IChoice.Add(parser.Seq("definitiveObjIDComponentNameAndNumber", Identifier(), parser.Block("defIDNr", "(", ")", parser.Number())))
	doidc.IChoice.Add(Identifier())
	return doidc
}

type definitiveObjIDComponent struct {
	parser.IChoice
}

//ITagDefault ...
//===== Notation =====
// TagDefault ::=
// 	EXPLICIT TAGS
// |IMPLICIT TAGS
// |AUTOMATIC TAGS
// |empty
type ITagDefault interface {
	parser.IOptional
}

//TagDefault ...
func TagDefault() ITagDefault {
	return &tagDefault{
		IOptional: parser.Optional(
			parser.Choice("tagDefaultOptions",
				parser.Seq("explicitTags", parser.Keyword("EXPLICIT"), parser.Keyword("TAGS")),
				parser.Seq("implicitTags", parser.Keyword("IMPLICIT"), parser.Keyword("TAGS")),
				parser.Seq("automaticTags", parser.Keyword("AUTOMATIC"), parser.Keyword("TAGS")),
			),
		),
	}
}

type tagDefault struct {
	parser.IOptional
}

func (td tagDefault) Value() string {
	return "NYI"
}

//IExtensionDefault ...
//===== Notation =====
// ExtensionDefault ::=
//  EXTENSIBILITY IMPLIED
// |empty
type IExtensionDefault interface {
	parser.IOptional
}

//ExtensionDefault ...
func ExtensionDefault() IExtensionDefault {
	return &extensionDefault{
		IOptional: parser.Optional(
			parser.Choice("extensionDefaultOptions",
				parser.Seq("EXTENSIBILITY IMPLIED", parser.Keyword("EXTENSIBILITY"), parser.Keyword("IMPLIED")),
			),
		),
	}
}

type extensionDefault struct {
	parser.IOptional
}

//IModuleBody ...
//===== Nodation =====
// ModuleBody ::=
// Exports Imports AssignmentList
// |
// empty
type IModuleBody interface {
	parser.IOptional
}

//ModuleBody ...
func ModuleBody() IModuleBody {
	return &moduleBody{
		IOptional: parser.Optional(
			parser.Seq("moduleBody", Exports(), Imports(), AssignmentList()),
		),
	}
}

type moduleBody struct {
	parser.IOptional
}

//IExports ...
//===== Notation =====
// Exports ::=
// 	EXPORTS SymbolsExported ";"
// |EXPORTS ALL ";"
// |empty
type IExports interface {
	parser.IOptional
	All() bool
	SymbolList() ISymbolsExported
}

//Exports ...
func Exports() IExports {
	e := &exports{
		symbolsExported: SymbolsExported(),
		all:             false,
	}
	e.IOptional = parser.Optional(
		parser.Choice("exports",
			parser.Seq("exportedSymbols", parser.Keyword("EXPORTS"), e.symbolsExported, parser.Keyword(";")),
			parser.Seq("exportsAll", parser.Keyword("EXPORTS"), parser.Keyword("ALL"), parser.Keyword(";")),
		),
	)
	return e
}

type exports struct {
	parser.IOptional
	all             bool
	symbolsExported ISymbolsExported
}

func (e exports) All() bool {
	return e.all
}

func (e exports) SymbolList() ISymbolsExported {
	return e.symbolsExported
}

//ISymbolsExported ...
//===== Notation =====
// SymbolsExported ::=
//  SymbolList
// |empty
type ISymbolsExported interface {
	parser.IOptional
	SymbolList() ISymbolList
}

//SymbolsExported ...
func SymbolsExported() ISymbolsExported {
	se := &symbolsExported{
		symbolList: SymbolList(),
	}
	se.IOptional = parser.Optional(se.symbolList)
	return se
}

type symbolsExported struct {
	parser.IOptional
	symbolList ISymbolList
}

func (se symbolsExported) SymbolList() ISymbolList {
	return se.symbolList
}

//IImports ...
//===== Notation =====
// Imports ::=
//   IMPORTS SymbolsImported ";"
// | empty
type IImports interface {
	parser.IOptional
	SymbolsImported() ISymbolsImported
}

//Imports ...
func Imports() IImports {
	i := &imports{
		symbolsImported: SymbolsImported(),
	}
	i.IOptional = parser.Optional(
		parser.Seq("imports", parser.Keyword("IMPORTS"), i.symbolsImported, parser.Keyword(";")),
	)
	return i
}

type imports struct {
	parser.IOptional
	symbolsImported ISymbolsImported
}

func (i imports) SymbolsImported() ISymbolsImported {
	return i.symbolsImported
}

//ISymbolsImported ...
//===== Notation =====
// SymbolsImported ::=
//	SymbolsFromModuleList
// |empty
type ISymbolsImported interface {
	parser.IOptional
}

//SymbolsImported ...
func SymbolsImported() ISymbolsImported {
	return &symbolsImported{
		IOptional: parser.Optional(SymbolsFromModuleList()),
	}
}

type symbolsImported struct {
	parser.IOptional
}

//ISymbolsFromModuleList ...
type ISymbolsFromModuleList interface {
	parser.IList
}

//SymbolsFromModuleList ...
func SymbolsFromModuleList() ISymbolsFromModuleList {
	return &symbolsFromModuleList{
		IList: parser.List(" ", 1, SymbolsFromModule()),
	}
}

type symbolsFromModuleList struct {
	parser.IList
}

//ISymbolsFromModule ...
//===== Notation =====
// SymbolsFromModule ::= SymbolList FROM GlobalModuleReference
type ISymbolsFromModule interface {
	parser.ISeq
}

//SymbolsFromModule ...
func SymbolsFromModule() ISymbolsFromModule {
	return &symbolsFromModule{
		ISeq: parser.Seq("symbolsFromModule", SymbolList(), parser.Keyword("FROM"), GlobalModuleReference()),
	}
}

//IGlobalModuleReference ...
//===== Notation =====
// GlobalModuleReference ::= modulereference AssignedIdentifier
type IGlobalModuleReference interface {
	parser.ISeq
}

//GlobalModuleReference ...
func GlobalModuleReference() IGlobalModuleReference {
	return &globalModuleReference{
		ISeq: parser.Seq("globalModuleReference", ModuleReference(), AssignedIdentifier()),
	}
}

//IAssignedIdentifier ...
//===== Notation =====
// AssignedIdentifier ::=
//  ObjectIdentifierValue
// |DefinedValue
// |empty
type IAssignedIdentifier interface {
	parser.IOptional
}

//AssignedIdentifier ...
func AssignedIdentifier() IAssignedIdentifier {
	return &assignedIdentifier{
		IOptional: parser.Optional(
			parser.Choice("assignedIdentifierOptions",
				ObjectIdentifierValue(),
				DefinedValue(),
			),
		),
	}
}

type assignedIdentifier struct {
	parser.IOptional
}

type globalModuleReference struct {
	parser.ISeq
}

type symbolsFromModule struct {
	parser.ISeq
}

//IAssignmentList ...
type IAssignmentList interface {
	parser.IList
}

//AssignmentList ...
func AssignmentList() IAssignmentList {
	return &assignmentList{
		IList: parser.List(" ", 0, Assignment()),
	}
}

type assignmentList struct {
	parser.IList
}

//IAssignment ...
//===== Notation =====
// Assignment ::=
//  TypeAssignment
// |ValueAssignment
// |XMLValueAssignment
// |ValueSetTypeAssignment
// |ObjectClassAssignment
// |ObjectAssignment
// |ObjectSetAssignment
// |ParameterizedAssignment
type IAssignment interface {
	parser.IChoice
}

//Assignment ...
func Assignment() IAssignment {
	return &assignment{
		IChoice: parser.Choice("assignment",
			TypeAssignment(),
			ValueAssignment(),
			//todo: XMLValueAssignment(),
			//todo: ValueSetTypeAssignment(),
			//todo: ObjectClassAssignment(),
			//todo: ObjectAssignment(),
			//todo: ObjectSetAssignment(),
			//todo: ParameterizedAssignment(),
		),
	}
}

type assignment struct {
	parser.IChoice
}

//ITypeAssignment ...
type ITypeAssignment interface {
	parser.ISeq
}

//TypeAssignment ...
func TypeAssignment() ITypeAssignment {
	return &typeAssignment{
		ISeq: parser.Seq("typeAssignment", TypeReference(), parser.Keyword("::="), Type()),
	}
}

type typeAssignment struct {
	parser.ISeq
}

//IValueAssignment ...
type IValueAssignment interface {
	parser.ISeq
}

//ValueAssignment ...
func ValueAssignment() ITypeAssignment {
	return &valueAssignment{
		ISeq: parser.Seq("valueAssignment", ValueReference(), Type()),
	}
}

type valueAssignment struct {
	parser.ISeq
}

//IType ...
//===== Notation =====
// Type ::= BuiltinType | ReferencedType | ConstrainedType
type IType interface {
	parser.IChoice
}

//Type ...
func Type() IType {
	return &_type{
		IChoice: parser.Choice("type",
			BuiltinType(),
			ReferencedType(),
			ConstrainedType(),
		),
	}
}

type _type struct {
	parser.IChoice
}

//IBuiltinType ...
//===== Notation =====
// BuiltinType ::=
//  BitStringType
// |BooleanType
// |CharacterStringType
// |ChoiceType
// |EmbeddedPDVType
// |EnumeratedType
// |ExternalType
// |InstanceOfType
// |IntegerType
// |NullType
// |ObjectClassFieldType
// |ObjectIdentifierType
// |OctetStringType
// |RealType
// |RelativeOIDType
// |SequenceType
// |SequenceOfType
// |SetType
// |SetOfType
// |TaggedType
type IBuiltinType interface {
	parser.IChoice
}

//BuiltinType ...
func BuiltinType() IBuiltinType {
	return &builtinType{
		IChoice: parser.Choice("builtinType",
			BitStringType(),
			BooleanType(),
			CharacterStringType(),
			//todo: ChoiceType(), -- creates circular reference!
			EmbeddedPDVType(),
			EnumeratedType(),
			ExternalType(),
			// todo: InstanceOfType(),
			IntegerType(),
			NullType(),
			// todo: ObjectClassFieldType(),
			// todo: ObjectIdentifierType(),
			OctetStringType(),
			// todo: RealType(),
			// todo: RelativeOIDType(),
			// todo: SequenceType(), -- creates circular reference
			// todo: SequenceOfType(), --creates circular reference
			// todo: SetType(),
			// todo: SetOfType(),
			// todo: TaggedType(),
		),
	}
}

type builtinType struct {
	parser.IChoice
}

//IBitStringType ...
//===== Notation =====
// BitStringType ::=
//   BIT STRING
// | BIT STRING " { " NamedBitList " } "
type IBitStringType interface {
	parser.ISeq
}

//BitStringType ...
func BitStringType() IBitStringType {
	return &bitStringType{
		ISeq: parser.Seq("bitString",
			parser.Keyword("BIT"),
			parser.Keyword("STRING"),
			parser.Optional(parser.Block("bitStringBitsList", "{", "}", NamedBitList())),
		),
	}
}

type bitStringType struct {
	parser.ISeq
}

//INamedBitList ...
// NamedBitList ::=
//   NamedBit
// | NamedBitList "," NamedBit
type INamedBitList interface {
	parser.IList
}

//NamedBitList ...
func NamedBitList() INamedBitList {
	return &namedBitList{
		IList: parser.List(",", 1, NamedBit()),
	}
}

type namedBitList struct {
	parser.IList
}

//INamedBit ...
//===== Notation =====
// NamedBit ::=
//   identifier " ( " number " ) "
// | identifier " ( " DefinedValue " ) "
type INamedBit interface {
	parser.ISeq
}

//NamedBit ...
func NamedBit() INamedBit {
	return &namedBit{
		ISeq: parser.Seq("namedBit",
			Identifier(),
			parser.Block("namedBitValue", "(", ")",
				parser.Choice("namedBitValue",
					parser.Number(),
					DefinedValue(),
				),
			),
		),
	}
}

type namedBit struct {
	parser.ISeq
}

//IBooleanType ...
//===== Notation =====
//BooleanType ::= BOOLEAN
type IBooleanType interface {
	parser.IKeyword
}

//BooleanType ...
func BooleanType() IBooleanType {
	return &booleanType{
		IKeyword: parser.Keyword("BOOLEAN"),
	}
}

type booleanType struct {
	parser.IKeyword
}

//ICharacterStringType ...
//===== Notation =====
// CharacterStringType ::= RestrictedCharacterStringType | UnrestrictedCharacterStringType
type ICharacterStringType interface {
	parser.IChoice
}

//CharacterStringType ...
func CharacterStringType() ICharacterStringType {
	return &characterStringType{
		IChoice: parser.Choice("characterStringType",
			RestrictedCharacterStringType(),
			UnrestrictedCharacterStringType(),
		),
	}
}

type characterStringType struct {
	parser.IChoice
}

//IRestrictedCharacterStringType ...
//===== Notation =====
// RestrictedCharacterStringType ::=
//   BMPString
// | GeneralString
// | GraphicString
// | IA5String
// | ISO646String
// | NumericString
// | PrintableString
// | TeletexString
// | T61String
// | UniversalString
// | UTF8String
// | VideotexString
// | VisibleString
type IRestrictedCharacterStringType interface {
	parser.IChoice
}

//RestrictedCharacterStringType ...
func RestrictedCharacterStringType() IRestrictedCharacterStringType {
	return &restrictedCharacterStringType{
		IChoice: parser.Choice("restrictedCharacterStringType"),
		//BMPString(),
		// | GeneralString
		// | GraphicString
		// | IA5String
		// | ISO646String
		// | NumericString
		// | PrintableString
		// | TeletexString
		// | T61String
		// | UniversalString
		// | UTF8String
		// | VideotexString
		// | VisibleString

	}
}

type restrictedCharacterStringType struct {
	parser.IChoice
}

//IUnrestrictedCharacterStringType ...
type IUnrestrictedCharacterStringType interface {
	parser.ISeq
}

//UnrestrictedCharacterStringType ...
func UnrestrictedCharacterStringType() IUnrestrictedCharacterStringType {
	return &unrestrictedCharacterStringType{
		ISeq: parser.Seq("UnrestrictedCharacterStringType",
			parser.Keyword("CHARACTER"),
			parser.Keyword("STRING"),
		),
	}
}

type unrestrictedCharacterStringType struct {
	parser.ISeq
}

//IChoiceType ...
//===== Notation =====
//ChoiceType ::= CHOICE " { " AlternativeTypeLists " } "
type IChoiceType interface {
	parser.ISeq
}

//ChoiceType ...
func ChoiceType() IChoiceType {
	return &choiceType{
		ISeq: parser.Seq("choiceType",
			parser.Keyword("CHOICE"),
			parser.Block("choiceType", "{", "}", AlternativeTypeList())),
	}
}

type choiceType struct {
	parser.ISeq
}

//IAlternativeTypeList ...
type IAlternativeTypeList interface {
	parser.IList
}

//AlternativeTypeList ...
func AlternativeTypeList() IAlternativeTypeList {
	return &alternativeTypeList{
		IList: parser.List(",", 1, NamedType()),
	}
}

type alternativeTypeList struct {
	parser.IList
}

//IEmbeddedPDVType ...
type IEmbeddedPDVType interface {
	parser.ISeq
}

//EmbeddedPDVType ...
func EmbeddedPDVType() IEmbeddedPDVType {
	return &embeddedPdvType{
		ISeq: parser.Seq("embeddedPdvType",
			parser.Keyword("EMBEDDED"),
			parser.Keyword("PDV"),
		),
	}
}

type embeddedPdvType struct {
	parser.ISeq
}

//IEnumeratedType ...
//===== Notation =====
// ENUMERATED " { " Enumerations " } "
type IEnumeratedType interface {
	parser.ISeq
}

//EnumeratedType ...
func EnumeratedType() IEnumeratedType {
	return &enumeratedType{
		ISeq: parser.Seq("enumeratedType",
			parser.Keyword("ENUMERATED"),
			parser.Block("enumerated", "{", "}", Enumerations()),
		),
	}
}

type enumeratedType struct {
	parser.ISeq
}

//IEnumerations ...
//===== Notation =====
//	Enumerations ::=
//		RootEnumeration
//	|	RootEnumeration "," " ... " ExceptionSpec
//	|	RootEnumeration "," " ... " ExceptionSpec "," AdditionalEnumeration
type IEnumerations interface {
	parser.ISeq
}

//Enumerations ...
func Enumerations() IEnumerations {
	return &enumerations{
		ISeq: parser.Seq("enumerations",
			RootEnumeration(),
			//todo...
		),
	}
}

type enumerations struct {
	parser.ISeq
}

//IRootEnumeration ...
//===== Notation =====
//	RootEnumeration ::= Enumeration
type IRootEnumeration interface {
	IEnumeration
}

//RootEnumeration ...
func RootEnumeration() IRootEnumeration {
	return &rootEnumeration{
		IEnumeration: Enumeration(),
	}
}

type rootEnumeration struct {
	IEnumeration
}

//IEnumeration ...
//===== Notation =====
//	Enumeration ::= EnumerationItem | EnumerationItem "," Enumeration
type IEnumeration interface {
	parser.IList
}

//Enumeration ...
func Enumeration() IEnumeration {
	return &enumeration{
		IList: parser.List(",", 1, EnumerationItem()),
	}
}

type enumeration struct {
	parser.IList
}

//IEnumerationItem ...
//===== Notation =====
//	EnumerationItem ::= identifier | NamedNumber
type IEnumerationItem interface {
	parser.IChoice
}

//EnumerationItem ...
func EnumerationItem() IEnumerationItem {
	return &enumerationItem{
		IChoice: parser.Choice("enumerationItem",
			Identifier(),
			NamedNumber(),
		),
	}
}

type enumerationItem struct {
	parser.IChoice
}

//INamedNumber ...
// NamedNumber ::=
//		identifier " ( " SignedNumber " ) "
//	|	identifier " ( " DefinedValue " ) "
type INamedNumber interface {
	parser.ISeq
}

//NamedNumber ...
func NamedNumber() INamedNumber {
	return &namedNumber{
		ISeq: parser.Seq("namedNumber",
			Identifier(),
			parser.Block("namedNumber", "(", ")",
				parser.Choice("namedNumberValue",
					parser.SignedNumber(),
					DefinedValue(),
				),
			),
		),
	}
}

type namedNumber struct {
	parser.ISeq
}

//IExternalType ...
type IExternalType interface {
	parser.IKeyword
}

//ExternalType ...
func ExternalType() IExternalType {
	return &externalType{
		IKeyword: parser.Keyword("EXTERNAL"),
	}
}

type externalType struct {
	parser.IKeyword
}

//IInstanceOfType ...
//"InstanceOfType" are specified in ITU-T Rec. X.681 | ISO/IEC 8824-2, 14.1 and Annex C.
type IInstanceOfType interface{}

//InstanceOfType ...
func InstanceOfType() IInstanceOfType {
	return &instanceOfType{}
}

type instanceOfType struct{}

//IIntegerType ...
//===== Notation =====
// IntegerType ::=
//   INTEGER
// | INTEGER " { " NamedNumberList " } "
type IIntegerType interface {
	parser.ISeq
}

//IntegerType ...
func IntegerType() IIntegerType {
	return &integerType{
		ISeq: parser.Seq("integerType",
			parser.Keyword("INTEGER"),
			parser.Optional(parser.Block("namedNumberListBlock", "{", "}", NamedNumberList())),
		),
	}
}

type integerType struct {
	parser.ISeq
}

//INamedNumberList ...
//===== Notation =====
//	NamedNumberList ::=
//		NamedNumber
//	|	NamedNumberList "," NamedNumber
type INamedNumberList interface {
	parser.IList
}

//NamedNumberList ...
func NamedNumberList() INamedNumberList {
	return &namedNumberList{
		IList: parser.List(",", 1, NamedNumber()),
	}
}

type namedNumberList struct {
	parser.IList
}

//INullType ...
type INullType interface {
	parser.IKeyword
}

//NullType ...
func NullType() INullType {
	return &nullType{
		IKeyword: parser.Keyword("NULL"),
	}
}

type nullType struct {
	parser.IKeyword
}

// todo: ObjectClassFieldType(),
// todo: ObjectIdentifierType(),

//IOctetStringType ...
type IOctetStringType interface {
	parser.ISeq
}

//OctetStringType ...
func OctetStringType() IOctetStringType {
	return &octetStringType{
		ISeq: parser.Seq("octetString",
			parser.Keyword("OCTET"),
			parser.Keyword("STRING"),
		),
	}
}

type octetStringType struct {
	parser.ISeq
}

// todo: RealType(),
// todo: RelativeOIDType(),

//ISequenceType ...
//===== Notation =====
// SequenceType ::=
//    SEQUENCE " { " " } "
//  | SEQUENCE " { " ExtensionAndException OptionalExtensionMarker " } "
//  | SEQUENCE " { " ComponentTypeLists " } "
type ISequenceType interface {
	parser.ISeq
}

//SequenceType ...
func SequenceType() ISequenceType {
	return &sequenceType{
		ISeq: parser.Seq("sequenceType",
			parser.Keyword("SEQUENCE"),
			parser.Block("sequenceType", "{", "}",
				parser.Optional(
					parser.Choice("sequenceType",
						ComponentTypeList(),
						parser.Seq("extAndExceptionWithOptExtMarker", ExtensionAndException(), OptionalExtensionMarker()),
					),
				),
			),
		),
	}
}

type sequenceType struct {
	parser.ISeq
}

//ISequenceOfType ...
//===== Notation =====
// SequenceOfType ::= SEQUENCE OF Type | SEQUENCE OF NamedType
type ISequenceOfType interface {
	parser.ISeq
}

//SequenceOfType ...
func SequenceOfType() ISequenceOfType {
	return &sequenceOfType{
		ISeq: parser.Seq("sequenceOfType",
			parser.Keyword("SEQUENCE"),
			parser.Keyword("OF"),
		),
	}
}

type sequenceOfType struct {
	parser.ISeq
}

//IComponentTypeList ...
// ComponentTypeList ::=
// 		ComponentType
// 	|	ComponentTypeList "," ComponentType
type IComponentTypeList interface {
	parser.IList
}

//ComponentTypeList ...
func ComponentTypeList() IComponentTypeList {
	return &componentTypeList{
		IList: parser.List(",", 1, ComponentType()),
	}
}

type componentTypeList struct {
	parser.IList
}

//IComponentType ...
//===== Notation =====
// ComponentType ::=
// 		NamedType
// 	|	NamedType OPTIONAL
// 	|	NamedType DEFAULT Value
// 	|	COMPONENTS OF Type
type IComponentType interface {
	parser.IChoice
}

//ComponentType ...
func ComponentType() IComponentType {
	return &componentType{
		IChoice: parser.Choice("componentType",
			parser.Seq("componentsOfType",
				parser.Keyword("COMPONENTS"),
				parser.Keyword("OF"),
				Type(),
			),
			parser.Seq("componentsUsingNamedType",
				NamedType(),
				parser.Optional(
					parser.Choice("componentsUsingNamedTypeDecorator",
						parser.Keyword("OPTIONAL"),
						parser.Seq("componentsUsingNamedTypeWithDefaultValue",
							parser.Keyword("DEFAULT"),
							Value(),
						),
					),
				),
			),
		),
	}
}

type componentType struct {
	parser.IChoice
}

//IExtensionAndException ...
//	ExtensionAndException ::= " ... " | " ... " ExceptionSpec
type IExtensionAndException interface {
	parser.ISeq
}

//ExtensionAndException ...
func ExtensionAndException() IExtensionAndException {
	return &extensionAndException{
		ISeq: parser.Seq("extensionAndException",
			parser.Keyword("..."),
			parser.Optional(ExceptionSpec()),
		),
	}
}

type extensionAndException struct {
	parser.ISeq
}

//IExceptionSpec ...
type IExceptionSpec interface {
	parser.IOptional
}

//ExceptionSpec ...
//===== Notation =====
//	ExceptionSpec ::= " ! " ExceptionIdentification | empty
func ExceptionSpec() IExceptionSpec {
	return &exceptionSpec{
		IOptional: parser.Optional(
			parser.Seq("exceptionSpec", parser.Keyword("!"), ExceptionIdentification()),
		),
	}
}

type exceptionSpec struct {
	parser.IOptional
}

//IExceptionIdentification ...
//===== Notation =====
// ExceptionIdentification ::=
//		SignedNumber
//	|	DefinedValue
//	|	Type " : " Value
type IExceptionIdentification interface {
	parser.IChoice
}

//ExceptionIdentification ...
func ExceptionIdentification() IExceptionIdentification {
	return &exceptionIdentification{
		IChoice: parser.Choice("exceptionIdentification",
			parser.SignedNumber(),
			DefinedValue(),
			parser.Seq("exceptionIdentificationTypeAndValue", Type(), parser.Keyword(":"), Value()),
		),
	}
}

type exceptionIdentification struct {
	parser.IChoice
}

//IOptionalExtensionMarker ...
//===== Notation =====
//	OptionalExtensionMarker ::= "," " ... " | empty
type IOptionalExtensionMarker interface {
	parser.IOptional
}

//OptionalExtensionMarker ...
func OptionalExtensionMarker() IOptionalExtensionMarker {
	return &optionalExtensionMarker{
		IOptional: parser.Optional(
			parser.Seq("optionalExtensionMarker",
				parser.Keyword(","),
				parser.Keyword("..."),
			),
		),
	}
}

type optionalExtensionMarker struct {
	parser.IOptional
}

//INamedType ...
//===== Notation =====
//	NamedType ::= identifier Type
type INamedType interface {
	parser.ISeq
}

//NamedType ...
func NamedType() INamedType {
	return &namedType{
		ISeq: parser.Seq("namedType", Identifier(), Type()),
	}
}

type namedType struct {
	parser.ISeq
}

// todo: SetType(),
// todo: SetOfType(),
// todo: TaggedType(),

//IReferencedType ...
// ReferencedType ::=
//  DefinedType
// |UsefulType
// |SelectionType
// |TypeFromObject
// |ValueSetFromObjects
type IReferencedType interface {
	parser.IChoice
}

//ReferencedType ...
func ReferencedType() IReferencedType {
	return &referencedType{
		IChoice: parser.Choice("referencedType",
			DefinedType(),
			UsefulType(),
			SelectionType(),
			//todo: TypeFromObject(),
			//todo: ValueSetFromObjects(),
		),
	}
}

type referencedType struct {
	parser.IChoice
}

//UsefulType ...
//===== Notation =====
//	UsefulType ::= typereference
func UsefulType() ITypeReference {
	return TypeReference()
}

//ISelectionType ...
//===== Notation =====
//	SelectionType ::= identifier "<" Type
type ISelectionType interface {
	parser.ISeq
}

//SelectionType ...
func SelectionType() ISelectionType {
	return &selectionType{
		ISeq: parser.Seq("selectionType",
			Identifier(),
			parser.Keyword("<"),
			Type(),
		),
	}
}

type selectionType struct {
	parser.ISeq
}

//IDefinedType ...
//===== Notation =====
// DefinedType ::=
// 		ExternalTypeReference
// 	|	Typereference
// 	|	ParameterizedType
// 	|	ParameterizedValueSetType
type IDefinedType interface {
	parser.IChoice
}

//DefinedType ...
func DefinedType() IDefinedType {
	return &definedType{
		IChoice: parser.Choice("definedType",
			ExternalTypeReference(),
			TypeReference(),
			//todo: ParameterizedType(),
			//todo: ParameterizedValueSetType(),
		),
	}
}

type definedType struct {
	parser.IChoice
}

//IExternalTypeReference ...
//===== Notation =====
//	ExternalTypeReference ::= modulereference " . " typereference
type IExternalTypeReference interface {
	parser.ISeq
}

//ExternalTypeReference ...
func ExternalTypeReference() IExternalTypeReference {
	return &externalTypeReference{
		ISeq: parser.Seq("externalTypeReference",
			ModuleReference(),
			parser.Keyword("."),
			TypeReference(),
		),
	}
}

type externalTypeReference struct {
	parser.ISeq
}

//IParameterizedType ...
//===== Notation =====
//The type identified by a "ParameterizedType" and "ParameterizedValueSetType",
// and the value identified by a "ParameterizedValue" are specified in ITU-T Rec. X.683 | ISO/IEC 8824-4.
//todo: type IParameterizedType interface{}

//ParameterizedType ...
// func ParameterizedType() IParameterizedType {
// 	return parameterizedType{}
// }

// type parameterizedType struct{}

//IConstrainedType ...
// ConstrainedType ::=
//  Type Constraint
// |TypeWithConstraint
type IConstrainedType interface {
	parser.IChoice
}

//ConstrainedType ...
func ConstrainedType() IConstrainedType {
	return &constrainedType{
		IChoice: parser.Choice("constrainedType",
			TypeConstraint(),
			TypeWithConstraint(),
		),
	}
}

type constrainedType struct {
	parser.IChoice
}

//ITypeConstraint ...
//===== Notation =====
//TypeConstraint ::= Type
type ITypeConstraint interface {
	IType
}

//TypeConstraint ...
func TypeConstraint() ITypeConstraint {
	return &typeConstraint{
		IType: Type(),
	}
}

type typeConstraint struct {
	IType
}

//ITypeWithConstraint ...
//===== Notation =====
// TypeWithConstraint ::=
// 		SET Constraint OF Type
// 	|	SET SizeConstraint OF Type
// 	|	SEQUENCE Constraint OF Type
// 	|	SEQUENCE SizeConstraint OF Type
// 	|	SET Constraint OF NamedType
// 	|	SET SizeConstraint OF NamedType
// 	|	SEQUENCE Constraint OF NamedType
// 	|	SEQUENCE SizeConstraint OF NamedType
type ITypeWithConstraint interface {
	parser.IChoice
}

//TypeWithConstraint ...
func TypeWithConstraint() ITypeWithConstraint {
	return &typeWithConstraint{
		IChoice: parser.Choice("typeWithConstraint",
			TypeWithConstraintSet(),
			TypeWithConstraintSequence(),
		),
	}
}

type typeWithConstraint struct {
	parser.IChoice
}

//ITypeWithConstraintSet ...
type ITypeWithConstraintSet interface {
	parser.ISeq
}

//TypeWithConstraintSet ...
func TypeWithConstraintSet() ITypeWithConstraintSet {
	return &typeWithConstraintSet{
		ISeq: parser.Seq("typeWithConstraintSet",
			parser.Keyword("SET"),
			parser.Choice("typeWithConstraintSetConstraintOptions", Constraint(), SizeConstraint()),
			parser.Keyword("OF"),
			parser.Choice("typeWithConstraintSetTypeOptions", Type(), NamedType()),
		),
	}
}

type typeWithConstraintSet struct {
	parser.ISeq
}

//ITypeWithConstraintSequence ...
type ITypeWithConstraintSequence interface {
	parser.ISeq
}

//TypeWithConstraintSequence ...
func TypeWithConstraintSequence() ITypeWithConstraintSequence {
	return &typeWithConstraintSequence{
		ISeq: parser.Seq("typeWithConstraintSequence",
			parser.Keyword("SEQUENCE"),
			parser.Choice("typeWithConstraintSetConstraintOptions", Constraint(), SizeConstraint()),
			parser.Keyword("OF"),
			parser.Choice("typeWithConstraintSetTypeOptions", Type(), NamedType()),
		),
	}
}

type typeWithConstraintSequence struct {
	parser.ISeq
}

//IConstraint ...
//===== Notation =====
//	Constraint ::= " ( " ConstraintSpec ExceptionSpec " ) "
type IConstraint interface {
	parser.IBlock
}

//Constraint ...
func Constraint() IConstraint {
	return &constraint{
		IBlock: parser.Block("constraint", "(", ")",
			parser.Seq("constraint",
				ConstraintSpec(),
				ExceptionSpec(),
			),
		),
	}
}

type constraint struct {
	parser.IBlock
}

//IConstraintSpec ...
//===== Notation =====
// ConstraintSpec ::=
// 		SubtypeConstraint
// 	|	GeneralConstraint
type IConstraintSpec interface {
	parser.IChoice
}

//ConstraintSpec ...
func ConstraintSpec() IConstraintSpec {
	return &constraintSpec{
		IChoice: parser.Choice("constraintSpec",
			SubtypeConstraint(),
			//todo: GeneralConstraint(),
		),
	}
}

type constraintSpec struct {
	parser.IChoice
}

//ISubtypeConstraint ...
//	SubtypeConstraint ::= ElementSetSpecs
type ISubtypeConstraint interface {
	IElementSetSpecs
}

//SubtypeConstraint ...
func SubtypeConstraint() ISubtypeConstraint {
	return &subtypeConstraint{
		IElementSetSpecs: ElementSetSpecs(),
	}
}

type subtypeConstraint struct {
	IElementSetSpecs
}

//IElementSetSpecs ...
//===== Notation =====
// ElementSetSpecs ::=
// 		RootElementSetSpec
// 	|	RootElementSetSpec "," " ... "
// 	|	RootElementSetSpec "," " ... " "," AdditionalElementSetSpec
type IElementSetSpecs interface {
	parser.IChoice
}

//ElementSetSpecs ...
func ElementSetSpecs() IElementSetSpecs {
	return &elementSetSpecs{
		IChoice: parser.Choice("elementSetSpecs",
			RootElementSetSpec(),
			//todo: ....
		),
	}
}

type elementSetSpecs struct {
	parser.IChoice
}

//IRootElementSetSpec ...
//	RootElementSetSpec ::= ElementSetSpec
type IRootElementSetSpec interface {
	IElementSetSpec
}

//RootElementSetSpec ...
//	RootElementSetSpec ::= ElementSetSpec
func RootElementSetSpec() IRootElementSetSpec {
	return &rootElementSetSpec{
		IElementSetSpec: ElementSetSpec(),
	}
}

type rootElementSetSpec struct {
	IElementSetSpec
}

//IElementSetSpec ...
//===== Notation =====
//	ElementSetSpec ::= Unions | ALL Exclusions
type IElementSetSpec interface {
	parser.IChoice
}

//ElementSetSpec ...
func ElementSetSpec() IElementSetSpec {
	return &elementSetSpec{
		IChoice: parser.Choice("elementSetSpec",
			Unions(),
			parser.Seq("allElementSetSpecs",
				parser.Keyword("ALL"),
				Exclusions(),
			),
		),
	}
}

type elementSetSpec struct {
	parser.IChoice
}

//IExclusions ...
//	Exclusions ::= EXCEPT Elements
type IExclusions interface {
	parser.ISeq
}

//Exclusions ...
func Exclusions() IExclusions {
	return &exclusions{
		ISeq: parser.Seq("exclusions",
			parser.Keyword("EXCEPT"),
			Elements(),
		),
	}
}

type exclusions struct {
	parser.ISeq
}

//IElements ...
//===== Notation =====
//Elements ::=
//		SubtypeElements
//	|	ObjectSetElements
//	|	" ( " ElementSetSpec " ) "
type IElements interface {
	parser.IChoice
}

//Elements ...
func Elements() IElements {
	return &elements{
		IChoice: parser.Choice("elements",
			SubtypeElements(),
			//todo: ObjectSetElements(),
			parser.Block("elementsBlock", "(", ")", ElementSetSpec()),
		),
	}
}

type elements struct {
	parser.IChoice
}

//ISubtypeElements ...
// SubtypeElements ::=
// 		SingleValue
// 	|	ContainedSubtype
// 	|	ValueRange
// 	|	PermittedAlphabet
// 	|	SizeConstraint
// 	|	TypeConstraint
// 	|	InnerTypeConstraints
// 	|	PatternConstraint
type ISubtypeElements interface {
	parser.IChoice
}

//SubtypeElements ...
func SubtypeElements() ISubtypeElements {
	return &subtypeElements{
		IChoice: parser.Choice("subtypeElements",
			SingleValue(),
			ContainedSubtype(),
			ValueRange(),
			//todo: PermittedAlphabet(),
			SizeConstraint(),
			TypeConstraint(),
			//todo: InnerTypeConstraints(),
			//todo: PatternConstraint(),
		),
	}
}

type subtypeElements struct {
	parser.IChoice
}

//ISingleValue ...
type ISingleValue interface {
	IValue
}

//SingleValue ...
func SingleValue() ISingleValue {
	return &singleValue{
		IValue: Value(),
	}
}

type singleValue struct {
	IValue
}

//IContainedSubtype ...
type IContainedSubtype interface {
	parser.ISeq
}

//ContainedSubtype ...
func ContainedSubtype() IContainedSubtype {
	return &containedSubtype{
		ISeq: parser.Seq("containedSubtype",
			Includes(),
			Type(),
		),
	}
}

type containedSubtype struct {
	parser.ISeq
}

//IIncludes ...
type IIncludes interface {
	parser.IOptional
}

//Includes ...
func Includes() IIncludes {
	return &includes{
		IOptional: parser.Optional(parser.Keyword("INCLUDES")),
	}
}

type includes struct {
	parser.IOptional
}

//IValueRange ...
//	ValueRange ::= LowerEndPoint " .. " UpperEndPoint
type IValueRange interface {
	parser.ISeq
}

//ValueRange ...
func ValueRange() IValueRange {
	return &valueRange{
		ISeq: parser.Seq("valueRange",
			LowerEndPoint(),
			UpperEndPoint(),
		),
	}
}

type valueRange struct {
	parser.ISeq
}

//ILowerEndPoint ...
//===== Notation =====
//	LowerEndPoint ::= LowerEndValue | LowerEndValue "<"
type ILowerEndPoint interface {
	parser.ISeq
}

//LowerEndPoint ...
func LowerEndPoint() ILowerEndPoint {
	return &lowerEndpoint{
		ISeq: parser.Seq("lowerEndPoint",
			LowerEndValue(),
			parser.Optional(parser.Keyword("<")),
		),
	}
}

type lowerEndpoint struct {
	parser.ISeq
}

//IUpperEndPoint ...
//===== Notation =====
//	UpperEndPoint ::= UpperEndValue | "<" UpperEndValue
type IUpperEndPoint interface {
	parser.ISeq
}

//UpperEndPoint ...
func UpperEndPoint() IUpperEndPoint {
	return &upperEndpoint{
		ISeq: parser.Seq("upperEndpoint",
			parser.Optional(parser.Keyword("<")),
			UpperEndValue(),
		),
	}
}

type upperEndpoint struct {
	parser.ISeq
}

//ILowerEndValue ...
//===== Notation =====
//LowerEndValue ::= Value | MIN
type ILowerEndValue interface {
	parser.IChoice
}

//LowerEndValue ...
func LowerEndValue() ILowerEndValue {
	return &lowerEndValue{
		IChoice: parser.Choice("lowerEndValue",
			Value(),
			parser.Keyword("MIN"),
		),
	}
}

type lowerEndValue struct {
	parser.IChoice
}

//IUpperEndValue ...
//===== Notation =====
//UpperEndValue ::= Value | MAX
type IUpperEndValue interface {
	parser.IChoice
}

//UpperEndValue ...
func UpperEndValue() IUpperEndValue {
	return &upperEndValue{
		IChoice: parser.Choice("upperEndValue",
			Value(),
			parser.Keyword("MAX"),
		),
	}
}

type upperEndValue struct {
	parser.IChoice
}

//IUnions ...
//===== Notation =====
//	Unions ::=
// 		Intersections
//	|	UElems UnionMark Intersections
type IUnions interface {
	parser.IChoice
}

//Unions ...
func Unions() IUnions {
	return &unions{
		IChoice: parser.Choice("unions",
			Intersections(),
			parser.Seq("unions", UElems(), UnionMark(), Intersections()),
		),
	}
}

type unions struct {
	parser.IChoice
}

//IIntersections ...
//	Intersections ::=
//		IntersectionElements
//	|	IElems IntersectionMark IntersectionElements
type IIntersections interface {
	parser.IChoice
}

//Intersections ...
func Intersections() IIntersections {
	return &intersections{
		IChoice: parser.Choice("intersections",
			IntersectionElements(),
			parser.Seq("intersections", IElems(), IntersectionMark(), IntersectionElements()),
		),
	}
}

type intersections struct {
	parser.IChoice
}

//IIntersectionElements ...
//===== Notation =====
//	IntersectionElements ::= Elements | Elems Exclusions
type IIntersectionElements interface {
	parser.IChoice
}

//IntersectionElements ...
func IntersectionElements() IIntersectionElements {
	return &intersectionElements{
		IChoice: parser.Choice("intersectionElements",
			Elements(),
			parser.Seq("intersectionElementsWithExclusions", Elems(), Exclusions()),
		),
	}
}

type intersectionElements struct {
	parser.IChoice
}

//IiElems ...
//===== Notation =====
//	Elems ::= Elements
type IiElems interface {
	IElements
}

//Elems ...
func Elems() IiElems {
	return &elems{
		IElements: Elements(),
	}
}

type elems struct {
	IElements
}

//IiIElems ...
//===== Notation =====
//	IElems ::= Intersections
type IiIElems interface {
	IIntersections
}

//IElems ...
func IElems() IiIElems {
	return &ielems{
		IIntersections: Intersections(),
	}
}

type ielems struct {
	IIntersections
}

//IUElems ...
//===== Notation =====
//	UElems ::= Unions
type IUElems interface {
	IUnions
}

//UElems ...
func UElems() IUElems {
	return &uelems{
		IUnions: Unions(),
	}
}

type uelems struct {
	IUnions
}

//IIntersectionMark ...
//===== Notation =====
//	IntersectionMark ::= " ^ " | INTERSECTION
type IIntersectionMark interface {
	parser.IChoice
}

//IntersectionMark ...
func IntersectionMark() IIntersectionMark {
	return &intersectionMark{
		IChoice: parser.Choice("intersectionMark",
			parser.Keyword("^"),
			parser.Keyword("INTERSECTION"),
		),
	}
}

type intersectionMark struct {
	parser.IChoice
}

//IUnionMark ...
//===== Notation =====
//	UnionMark ::= "|" |	UNION
type IUnionMark interface {
	parser.IChoice
}

//UnionMark ...
func UnionMark() IUnionMark {
	return &unionMark{
		IChoice: parser.Choice("unionMark",
			parser.Keyword("|"),
			parser.Keyword("UNION"),
		),
	}
}

type unionMark struct {
	parser.IChoice
}

//IGeneralConstraint ...
//The "GeneralConstraint" is defined in ITU-T Rec. X.682 | ISO/IEC 8824-3, 8.1.
//type IGeneralConstraint interface{}

//GeneralConstraint ...
// func GeneralConstraint() IGeneralConstraint {
// 	return &generalConstraint{}
// }

//type generalConstraint struct{}

//ISizeConstraint ...
//===== Notation =====
//	SizeConstraint ::= SIZE Constraint
type ISizeConstraint interface {
	parser.ISeq
}

//SizeConstraint ...
func SizeConstraint() ISizeConstraint {
	return &sizeConstraint{
		ISeq: parser.Seq("sizeConstraint",
			parser.Keyword("SIZE"),
			Constraint(),
		),
	}
}

type sizeConstraint struct {
	parser.ISeq
}

//IValue ...
//===== Notation =====
//	Value ::= BuiltinValue | ReferencedValue | ObjectClassFieldValue
type IValue interface {
	parser.IChoice
}

//Value ...
func Value() IValue {
	return &value{
		IChoice: parser.Choice("value",
			BuiltinValue(),
			ReferencedValue(),
			//ObjectClassFieldValue(),
		),
	}
}

type value struct {
	parser.IChoice
}

//IBuiltinValue ...
// BuiltinValue ::=
// 	BitStringValue
// |BooleanValue
// |CharacterStringValue
// |ChoiceValue
// |EmbeddedPDVValue
// |EnumeratedValue
// |ExternalValue
// |InstanceOfValue
// |IntegerValue
// |NullValue
// |ObjectIdentifierValue
// |OctetStringValue
// |RealValue
// |RelativeOIDValue
// |SequenceValue
// |SequenceOfValue
// |SetValue
// |SetOfValue
// |TaggedValue
type IBuiltinValue interface {
	parser.IChoice
}

//BuiltinValue ...
func BuiltinValue() IBuiltinValue {
	return &builtinValue{
		IChoice: parser.Choice("builtinValue",
			BitStringValue(),
			BooleanValue(),
			CharacterStringValue(),
			ChoiceValue(),
			EmbeddedPdvValue(),
			EnumeratedValue(),
			ExternalValue(),
			//todo: ???InstanceOfValue(),
			IntegerValue(),
			NullValue(),
			ObjectIdentifierValue(),
			OctetStringValue(),
			//todo: RealValue(),
			//todo: RelativeOIDValue(),
			SequenceValue(),
			//todo: SequenceOfValue(),
			//todo: SetValue(),
			//todo: SetOfValue(),
			TaggedValue(),
		),
	}
}

type builtinValue struct {
	parser.IChoice
}

//IBitStringValue ...
//===== Notation =====
// BitStringValue ::=
// 	bstring
// |hstring
// |" { " IdentifierList " } "
// |" { " " } "
// |CONTAINING Value
type IBitStringValue interface {
	parser.IChoice
}

//BitStringValue ...
func BitStringValue() IBitStringValue {
	return &bitStringValue{
		IChoice: parser.Choice("bitStringValue",
			BString(),
			HString(),
			parser.Block("bitStringValueIdentifierListBlock", "{", "}", parser.Optional(IdentifierList())),
			parser.Seq("bitStringContainingValue", parser.Keyword("CONTAINING"), Value()),
		),
	}
}

type bitStringValue struct {
	parser.IChoice
}

//IBooleanValue ...
type IBooleanValue interface {
	parser.IChoice
}

//BooleanValue ...
func BooleanValue() IBooleanValue {
	return &booleanValue{
		IChoice: parser.Choice("booleanValue",
			parser.Keyword("TRUE"),
			parser.Keyword("FALSE"),
		),
	}
}

type booleanValue struct {
	parser.IChoice
}

//ICharacterStringValue ...
//CharacterStringValue ::=
// 	RestrictedCharacterStringValue
// |	UnrestrictedCharacterStringValue
type ICharacterStringValue interface {
	parser.IChoice
}

//CharacterStringValue ...
func CharacterStringValue() ICharacterStringValue {
	return &characterStringValue{
		IChoice: parser.Choice("characterStringValue",
			RestrictedCharacterStringValue(),
			UnrestrictedCharacterStringValue(),
		),
	}
}

type characterStringValue struct {
	parser.IChoice
}

//IRestrictedCharacterStringValue ...
//RestrictedCharacterStringValue ::= cstring | CharacterStringList | Quadruple | Tuple
type IRestrictedCharacterStringValue interface {
	parser.IChoice
}

//RestrictedCharacterStringValue ...
func RestrictedCharacterStringValue() IRestrictedCharacterStringValue {
	return &restrictedCharacterStringValue{
		IChoice: parser.Choice("restrictedCharacterStringValue",
			CString(),
			CharacterStringList(),
			//todo: Quadruple(),
			//todo: Tuple(),
		),
	}
}

type restrictedCharacterStringValue struct {
	parser.IChoice
}

//IUnrestrictedCharacterStringValue ...
type IUnrestrictedCharacterStringValue interface {
	ISequenceValue
}

//UnrestrictedCharacterStringValue ...
func UnrestrictedCharacterStringValue() IUnrestrictedCharacterStringValue {
	return &unrestrictedCharacterStringValue{
		ISequenceValue: SequenceValue(),
	}
}

type unrestrictedCharacterStringValue struct {
	ISequenceValue
}

//IChoiceValue ...
type IChoiceValue interface {
	parser.ISeq
}

//ChoiceValue ...
func ChoiceValue() IChoiceValue {
	return &choiceValue{
		ISeq: parser.Seq("choiceValue",
			Identifier(),
			parser.Keyword(":"),
			Value(),
		),
	}
}

type choiceValue struct {
	parser.ISeq
}

//IEmbeddedPdvValue ...
//===== Notation =====
// EmbeddedPDVValue ::= SequenceValue
type IEmbeddedPdvValue interface {
	ISequenceValue
}

//EmbeddedPdvValue ...
func EmbeddedPdvValue() IEmbeddedPdvValue {
	return &embeddedPdvValue{
		ISequenceValue: SequenceValue(),
	}
}

type embeddedPdvValue struct {
	ISequenceValue
}

//IEnumeratedValue ...
//===== Notation =====
//	EnumeratedValue ::= identifier
type IEnumeratedValue interface {
	IIdentifier
}

//EnumeratedValue ...
func EnumeratedValue() IEnumeratedValue {
	return &enumeratedValue{
		IIdentifier: Identifier(),
	}
}

type enumeratedValue struct {
	IIdentifier
}

//IExternalValue ...
//===== Notation =====
//	ExternalValue ::= SequenceValue
type IExternalValue interface {
	ISequenceValue
}

//ExternalValue ...
func ExternalValue() IExternalValue {
	return &externalValue{
		ISequenceValue: SequenceValue(),
	}
}

type externalValue struct {
	ISequenceValue
}

//IInstanceOfValue ...
//===== Notation =====
//???
//type IInstanceOfValue interface{}

//InstanceOfValue ...
//func InstanceOfValue() IInstanceOfValue { return &instanceOfValue{} }

//type instanceOfValue struct{}

//IIntegerValue ...
//===== Notation =====
//IntegerValue ::= SignedNumber | identifier
type IIntegerValue interface {
	parser.IChoice
}

//IntegerValue ...
func IntegerValue() IIntegerValue {
	return &integerValue{
		IChoice: parser.Choice("integerValue",
			parser.SignedNumber(),
			Identifier(),
		),
	}
}

type integerValue struct {
	parser.IChoice
}

//INullValue ...
//===== Notation =====
//	NullValue ::= NULL
type INullValue interface {
	parser.IKeyword
}

//NullValue ...
func NullValue() INullValue {
	return &nullValue{
		IKeyword: parser.Keyword("NULL"),
	}
}

type nullValue struct {
	parser.IKeyword
}

//IOctetStringValue ...
// OctetStringValue ::=
//		bstring
//	|	hstring
//	|	CONTAINING Value
type IOctetStringValue interface {
	parser.IChoice
}

//OctetStringValue ...
func OctetStringValue() IOctetStringValue {
	return &octetStringValue{
		IChoice: parser.Choice("octetStringValue",
			BString(),
			HString(),
			parser.Seq("octetStringValueContaining", parser.Keyword("CONTAINING"), Value()),
		),
	}
}

type octetStringValue struct {
	parser.IChoice
}

// func RealValue() IRealValue               { return &realValue{} }
// func RelativeOIDValue() IRelativeOIDValue { return &relativeOIDValue{} }

//ISequenceValue ...
//===== Notation =====
//SequenceValue ::=
// 	" { " ComponentValueList " } "
// |	" { " " } "
type ISequenceValue interface {
	parser.IBlock
}

//SequenceValue ...
func SequenceValue() ISequenceValue {
	return &sequenceValue{
		IBlock: parser.Block("sequenceValue", "{", "}",
			parser.Optional(ComponentValueList()),
		),
	}
}

type sequenceValue struct {
	parser.IBlock
}

//IComponentValueList ...
// ComponentValueList ::=
// 		NamedValue
// 	|	ComponentValueList "," NamedValue
type IComponentValueList interface {
	parser.IList
}

//ComponentValueList ...
func ComponentValueList() IComponentValueList {
	return &componentValueList{
		IList: parser.List(",", 1, NamedValue()),
	}
}

type componentValueList struct {
	parser.IList
}

// func SequenceOfValue() ISequenceOfValue   { return &sequenceOfValue{} }
// func SetValue() ISetValue                 { return &setValue{} }
// func SetOfValue() ISetOfValue             { return &setOfValue{} }

//INamedValue ...
//===== Notation =====
//	NamedValue ::= identifier Value
type INamedValue interface {
	parser.ISeq
}

//NamedValue ...
func NamedValue() INamedValue {
	return &namedValue{
		ISeq: parser.Seq("namedValue", Identifier(), Value()),
	}
}

type namedValue struct {
	parser.ISeq
}

//ITaggedValue ...
//===== Notation =====
//	TaggedValue ::= Value
type ITaggedValue interface {
	IValue
}

//TaggedValue ...
func TaggedValue() ITaggedValue {
	return &taggedValue{
		IValue: Value(),
	}
}

type taggedValue struct {
	IValue
}

//IReferencedValue ...
//===== Notation =====
//ReferencedValue ::=
// 	DefinedValue
// |	ValueFromObject
type IReferencedValue interface {
	parser.IChoice
}

//ReferencedValue ...
func ReferencedValue() IReferencedValue {
	return &referencedValue{
		IChoice: parser.Choice("referencedValue",
			DefinedValue(),
			//todo: ValueFromObject(),
		),
	}
}

type referencedValue struct {
	parser.IChoice
}

//IObjectClassFieldValue ...
//NOTE 1 – "ObjectClassFieldValue" and "XMLObjectClassFieldValue" are defined in ITU-T Rec. X.681 | ISO/IEC 8824-2, 14.6.
//type IObjectClassFieldValue interface{}

//ObjectClassFieldValue ...
// func ObjectClassFieldValue() IObjectClassFieldValue {
// 	return &objectClassFieldValue{}
// }

//type objectClassFieldValue struct{}

//ISymbolList ...
//===== Notation =====
// SymbolList ::=
// Symbol
// | SymbolList "," Symbol
type ISymbolList interface {
	parser.IList
}

//SymbolList ...
func SymbolList() ISymbolList {
	return &symbolList{
		IList: parser.List(",", 1, Symbol()),
	}
}

type symbolList struct {
	parser.IList
}

//ISymbol ...
//===== Notation =====
//Symbol ::=
//  Reference
// |ParameterizedReference
type ISymbol interface {
	parser.IChoice
}

//Symbol ...
func Symbol() ISymbol {
	return &symbol{
		IChoice: parser.Choice("symbolOptions",
			Reference(),
			//ParameterizedReference(),
		),
	}
}

type symbol struct {
	parser.IChoice
}

//IReference ...
//===== Notation =====
// Reference ::=
// typereference
// |valuereference
// |objectclassreference
// |objectreference
// |objectsetreference
type IReference interface {
	parser.IChoice
}

//Reference ...
func Reference() IReference {
	return &reference{
		IChoice: parser.Choice("referenceOptions",
			TypeReference(),
			ValueReference(),
			// todo: ObjectClassReference(),
			// todo: ObjectReference(),
			// todo: ObjectSetReference(),
		),
	}
}

type reference struct {
	parser.IChoice
}
