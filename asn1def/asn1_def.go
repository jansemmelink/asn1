package asn1def

import (
	"strconv"
	"strings"

	"bitbucket.org/vservices/dark/logger"
	"github.com/jansemmelink/asn1/parser"
)

var log = logger.New()

//New empty definition
func New() IDefinition {
	return definition{}
}

//IDefinition of ASN.1 data structures
type IDefinition interface {
	LoadFile(filename string) error
	//Unmarshal() error
}

type definition struct {
}

func (d definition) LoadFile(filename string) error {
	log.Debugf("Loading %s ...", filename)

	lines, err := parser.LinesFromFile(filename)
	if err != nil {
		return log.Wrapf(err, "Failed to read lines from file %s", filename)
	}

	log.Debugf("Read %d significant lines from %s:", lines.Count(), filename)

	//----------------------------------------------
	//construct the notation for an ASN.1 Definition
	//----------------------------------------------

	//identifiers are defined by this regex:
	identifier := parser.Regex("identifier", `[a-zA-Z][a-zA-Z0-9\-]*`)

	namedConstant := parser.Seq("namedConst",
		identifier,
		parser.Optional(
			parser.Block("constValue", "(", ")", parser.Number()),
		),
	)

	tag := parser.Block("tag", "[", "]", parser.Number())

	valueSpec := parser.Choice("valueSpec",
		parser.Number(),
		identifier)

	valueRange := parser.Block("range",
		"(", ")",
		parser.Seq("valueRange",
			valueSpec,
			parser.Optional(parser.Seq("valueRangeMax", parser.Keyword(".."), valueSpec)),
		),
	)

	sizeSpec := parser.Seq("sizeSpec",
		parser.Keyword("SIZE"),
		valueRange,
	)

	specTypeInteger := parser.Seq("specTypeInteger",
		parser.Keyword("INTEGER"),
		valueRange,
	)

	objectIdentifierKeywords := parser.Seq("objectIdKeywords",
		parser.Keyword("OBJECT"),
		parser.Keyword("IDENTIFIER"),
	)

	//Name of the ASN.1 Definition
	//Example:
	//	MAP-DialogueInformation {
	//		itu-t identified-organization (4) etsi (0) mobileDomain (0)
	//		gsm-Network (1) modules (3) map-DialogueInformation (3) version8 (8)
	//	}
	asn1defName := parser.Seq("asn1defName",
		identifier,
		parser.Block("numbers", "{", "}", parser.List(
			" ",
			1,
			parser.Seq("identifiers",
				parser.Regex("identifiers", `[a-zA-Z][a-zA-Z0-9- ]*`),
				parser.Block("value", "(", ")", parser.Number()),
			),
		)),
	)

	//union specification
	//Example:
	//	CHOICE {
	//		operationCode [0] OperationCode,
	//		errorCode [1] ErrorCode,
	//		userInfo [2] NULL
	//	}
	//
	//	CHOICE {
	//		allGPRSData NULL,
	//		contextIdList ContextIdList
	//	}
	//
	//	CHOICE {
	//		localValue INTEGER,
	// 		globalValue OBJECT IDENTIFIER
	//	}
	specTypeChoice := parser.Seq("specTypeChoice",
		parser.Keyword("CHOICE"),
		parser.Block("choiceItems", "{", "}", parser.List(
			",",
			1, //min
			parser.Seq("choiceItem",
				identifier,
				parser.Optional(tag),
				parser.Choice("oidOrId", objectIdentifierKeywords, identifier),
			),
		)),
	)

	//specification of enum as a ENUMERATED
	//Example:
	//	ENUMERATED {
	//		noReasonGiven (0),
	//		invalidDestinationReference (1),
	//		invalidOriginatingReference (2),
	//		encapsulatedAC-NotSupported (3) ,
	//		transportProtectionNotAdequate (4)}
	//	IST-SupportIndicator ::= ENUMERATED {
	//		basicISTSupported (0),
	//		istCommandSupported (1),
	//		...}
	specTypeEnumerated := parser.Seq("enumeratedSpec",
		parser.Keyword("ENUMERATED"),
		parser.Block("enumerators", "{", "}",
			parser.List(
				",",
				1, //min
				parser.Choice("enumeratorElem",
					namedConstant,
					parser.Keyword("..."), //after ... is extension - and unlisted values as possible but should be ignored
				),
			),
		),
	)

	//struct specification as a SEQUENCE
	//Example:
	//	::= SEQUENCE {
	// 		destinationReference [0] AddressString OPTIONAL,
	// 		originationReference [1] AddressString OPTIONAL,
	// 		...,
	// 		extensionContainer ExtensionContainer OPTIONAL
	// 		-- extensionContainer must not be used in version 2
	// 	}
	//
	//	::= SEQUENCE {
	//		extId											MAP-EXTENSION.&extensionId ({ExtensionSet}),
	//		extType											MAP-EXTENSION.&ExtensionType ({ExtensionSet}{@extId})	OPTIONAL
	//	}
	//
	//	::= SEQUENCE {
	// 	imsi											[0]  IMSI								OPTIONAL,
	// 	COMPONENTS OF SubscriberData,
	// 	extensionContainer								[14] ExtensionContainer					OPTIONAL,
	// 	... ,
	// 	naea-PreferredCI								[15] NAEA-PreferredCI					OPTIONAL,
	// 	-- naea-PreferredCI is included at the discretion of the HLR operator.
	// 	gprsSubscriptionData							[16] GPRSSubscriptionData				OPTIONAL,
	// 	roamingRestrictedInSgsnDueToUnsupportedFeature	[23] NULL								OPTIONAL,
	// 	networkAccessMode 								[24] NetworkAccessMode					OPTIONAL,
	// 	lsaInformation									[25] LSAInformation						OPTIONAL,
	// 	lmu-Indicator									[21] NULL								OPTIONAL,
	// 	lcsInformation									[22] LCSInformation						OPTIONAL,
	// 	istAlertTimer									[26] IST-AlertTimerValue				OPTIONAL,
	// 	superChargerSupportedInHLR						[27] AgeIndicator						OPTIONAL,
	// 	mc-SS-Info										[28] MC-SS-Info							OPTIONAL,
	// 	cs-AllocationRetentionPriority					[29] CS-AllocationRetentionPriority		OPTIONAL,
	// 	sgsn-CAMEL-SubscriptionInfo						[17] SGSN-CAMEL-SubscriptionInfo		OPTIONAL,
	// 	chargingCharacteristics							[18] ChargingCharacteristics			OPTIONAL
	// 	}
	classMemberRef := parser.Seq("classMemberRef",
		identifier,
		parser.Keyword(".&"),
		identifier,
		parser.Block("classMemberRefOptions", //e.g. ({ExtensionSet}) or ({ExtensionSet}{@extId})
			"(", ")",
			parser.List("", 1, parser.Block("option", "{", "}", parser.Regex("option", `[A-Za-z@][A-Za-z0-9]*`))),
		),
	)

	typeSpecification := parser.Seq("typeSpec",
		parser.Choice("memberType",
			objectIdentifierKeywords, //only: "OBJECT IDENTIFIER"
			classMemberRef,           //e.g. MAP-EXTENSION.&ExtensionType ({ExtensionSet}{@extId})
			identifier,               //e.g. ChargingCharacteristics
		),
		parser.Optional(parser.Keyword("OPTIONAL")),
	)

	structMemberSpecification := parser.Choice("structMemberSpec",
		parser.Keyword("..."),
		parser.Seq("componentsMember",
			parser.Keyword("COMPONENTS"),
			parser.Keyword("OF"),
			identifier,
		),
		parser.Seq("structMember",
			identifier,
			parser.Optional(tag),
			typeSpecification,
		),
	)

	//SEQUENCE is used for structs and lists
	//struct:
	specSequenceStruct := parser.Block("structMembers", "{", "}", parser.List(
		",",
		1,
		structMemberSpecification,
	))

	//list specification
	//Examples:
	//	TripletList ::= SEQUENCE SIZE (1..5) OF AuthenticationTriplet
	//	DestinationNumberLengthList ::= SEQUENCE SIZE (1..maxNumOfCamelDestinationNumberLengths) OF INTEGER(1..maxNumOfISDN-AddressDigits)
	listItemType := parser.Choice("listItemType",
		specTypeInteger,
		identifier,
	)
	specSequenceList := parser.Seq("listSpec",
		sizeSpec,
		parser.Keyword("OF"),
		listItemType,
	)

	specTypeSequence := parser.Seq("specTypeSequence",
		parser.Keyword("SEQUENCE"),
		parser.Choice("sequence",
			specSequenceStruct,
			specSequenceList,
		),
	)

	//class specification
	//Examples:
	// MAP-EXTENSION ::= CLASS {
	// 	&ExtensionType OPTIONAL,
	// 	&extensionId OBJECT IDENTIFIER }
	classMemberSpec := parser.Seq("classMemberSpec",
		parser.Keyword("&"),
		identifier,
		parser.Optional(objectIdentifierKeywords),
		parser.Optional(parser.Keyword("OPTIONAL")),
	)
	specTypeClass := parser.Seq("specTypeClass",
		parser.Keyword("CLASS"),
		parser.Block("classMembers",
			"{", "}",
			parser.List(",", 1, classMemberSpec),
		),
	)

	//specification of OBJECT IDENTIFIER
	//	<identifier> OBJECT IDENTIFIER ::= { ... }
	//Example:
	//map-DialogueAS OBJECT IDENTIFIER ::= {gsm-NetworkId as-Id map-DialoguePDU (1) version1 (1)}
	objectIdentifier := parser.Seq("objectId",
		objectIdentifierKeywords,
		parser.Keyword("::="),
		parser.Block("objectIdValues", "{", "}", parser.List(
			" ",
			1,
			namedConstant,
		)),
	)

	//OPERATION list specification:
	//Example:
	// Supported-MAP-Operations OPERATION ::= {updateLocation | cancelLocation | purgeMS |
	// 	provideSubscriberLocation | sendRoutingInfoForLCS | subscriberLocationReport |
	// 	secureTransportClass1 |secureTransportClass2 | secureTransportClass3 | secureTransportClass4}
	valueOperation1 := parser.Seq("operationSpec",
		parser.Keyword("OPERATION"),
		parser.Keyword("::="),
		parser.Block("operations", "{", "}", parser.List(
			"|",
			1,
			identifier,
		)),
	)

	//OPERATION specification:
	//Example:
	// updateLocation OPERATION ::= {
	// 	ARGUMENT UpdateLocationArg
	// 	RESULT UpdateLocationRes
	// 	ERRORS { systemFailure | dataMissing | unexpectedDataValue }
	// 	CODE local:2
	// }
	operArg := parser.Seq("operArg", parser.Keyword("ARGUMENT"), identifier)
	operRes := parser.Seq("operRes", parser.Keyword("RESULT"), identifier)
	operRet := parser.Seq("operRet", parser.Keyword("RETURN"), parser.Keyword("RESULT"), parser.Keyword("TRUE"))
	operErr := parser.Seq("operErr",
		parser.Keyword("ERRORS"),
		parser.Block("operationResults", "{", "}", parser.List("|", 1, identifier)),
	)
	operLink := parser.Seq("operLink",
		parser.Keyword("LINKED"),
		parser.Block("linked", "{", "}", identifier),
	)
	operCode := parser.Seq("operCode",
		parser.Keyword("CODE"),
		parser.Keyword("local"),
		parser.Keyword(":"),
		parser.Number(),
	)
	valueOperation2 := parser.Seq("operationSpec",
		parser.Keyword("OPERATION"),
		parser.Keyword("::="),
		parser.Block("operations", "{", "}", parser.Seq("operationDetails",
			parser.Optional(operArg),
			parser.Optional(operRes),
			parser.Optional(operRet),
			parser.Optional(operErr),
			parser.Optional(operLink),
			operCode,
		)),
	)

	//Example:
	// 	systemFailure ERROR ::= {
	//		PARAMETER
	//		SystemFailureParam
	//		-- optional
	//		CODE local:34
	//	}
	//	unknownMSC ERROR ::= {
	//		CODE local:3
	//	}
	valueError := parser.Seq("errorSpec",
		parser.Keyword("ERROR"),
		parser.Keyword("::="),
		parser.Block("errorDetails",
			"{", "}",
			parser.Seq("errorNameAndCode",
				parser.Optional(parser.Seq("errorParamId", parser.Keyword("PARAMETER"), identifier)),
				operCode,
			)),
	)

	//OCTET STRING
	//Examples:
	//	TBCD-STRING ::= OCTET STRING
	//	AgeIndicator ::= OCTET STRING (SIZE (1..6))
	//	RAND ::= OCTET STRING (SIZE (16))
	//	PermittedIntegrityProtectionAlgorithms ::= OCTET STRING (SIZE (1..maxPermittedIntegrityProtectionAlgorithmsLength))
	specTypeOctetString := parser.Seq("octetStrSpec",
		parser.Keyword("OCTET"),
		parser.Keyword("STRING"),
		parser.Optional(parser.Block("octStrOptions", "(", ")", sizeSpec)),
	)

	//BIT STRING
	//Examples:
	//	SupportedLCS-CapabilitySets ::= BIT STRING {
	//		lcsCapabilitySet1 (0),
	//		lcsCapabilitySet2 (1),
	//		lcsCapabilitySet3 (2) } (SIZE (2..16))
	specTypeBitString := parser.Seq("bitStrSpec",
		parser.Keyword("BIT"),
		parser.Keyword("STRING"),
		parser.Block("bitStringBits",
			"{", "}",
			parser.List(",", 1, namedConstant),
		),
		parser.Block("octStrOptions", "(", ")", sizeSpec),
	)

	//const integer value specification
	//Examples:
	//	maxPermittedIntegrityProtectionAlgorithmsLength INTEGER ::= 9
	valueInteger := parser.Seq("constIntSpec",
		parser.Choice("",
			parser.Keyword("INTEGER"),
			identifier,
		),
		parser.Keyword("::="),
		parser.Number(),
	)

	//const bit string value specification
	//Examples:
	// -- bits 87654321: group (bits 8765), and specific service
	// -- (bits 4321)
	// allSS SS-Code ::= '00000000'B
	// -- reserved for possible future use
	// -- all SS
	// allLineIdentificationSS SS-Code ::= '00010000'B
	// -- reserved for possible future use
	// -- all line identification SS
	// clip SS-Code ::= '00010001'B
	// -- calling line identification presentation
	bitsValue := parser.Regex("bitsValue", `[0-9A-F][0-9A-F]*`)
	valueBitString := parser.Seq("constBitStrSpec",
		identifier,
		parser.Keyword("::="),
		parser.Block("bitsValue", "'", "'", bitsValue),
		parser.Keyword("B"),
	)

	//value set is just written as this using literal "..." to indicate anything:
	//Example:
	//	ExtensionSet MAP-EXTENSION ::= { ... }
	//-- ExtensionSet is the set of all defined private extensions
	valueSet := parser.Seq("valueSet",
		identifier, //type of items in the set
		parser.Keyword("::="),
		parser.Block("valueSetBlock", "{", "}", parser.Keyword("...")),
	)

	//null
	//Example:
	//SuppressionOfAnnouncement ::= NULL
	typeSpecNull := parser.Seq("nullSpec",
		parser.Keyword("NULL"),
	)

	//numeric string is an ASN.1 type
	//there are several string types, including IA5String, VisibleString, PrintableString, and NumericString...
	//Examples:
	//	Password ::= NumericString
	//		(FROM ("0"|"1"|"2"|"3"|"4"|"5"|"6"|"7"|"8"|"9"))
	//		(SIZE (4))
	quotedCharacter := parser.Block("quotedCharacter", "\"", "\"", parser.Regex("character", `.`))
	typeSpecNumericString := parser.Seq("numStrSpec",
		parser.Keyword("NumericString"),
		parser.Block("numStrCharSet",
			"(", ")",
			parser.Seq("numStrCharSetSeq",
				parser.Keyword("FROM"),
				parser.Block("numStrCharSetChars",
					"(", ")",
					parser.List("|", 1, quotedCharacter),
				),
			)),
		parser.Block("octStrOptions", "(", ")", sizeSpec),
	)

	//new type, also to constrain another type
	//		<new-type> ::= <type> (SIZE(min[..max]))
	//Examples:
	//	HLR-Id ::= IMSI
	//	ISDN-AddressString ::= AddressString (SIZE (1..maxISDN-AddressLength))
	typeSpecSized := parser.Seq("typeSpecSized",
		identifier,
		parser.Optional(
			parser.Block("newTypeSizeSpec", "(", ")", sizeSpec),
		),
	)

	//specification of a definition can be written as one of the following:
	//	types are defined as:
	//		<identifier> ::= [tag] ...
	specType := parser.Seq("specType",
		parser.Keyword("::="),
		parser.Optional(tag),
		parser.Choice("specTypeOptions",
			specTypeChoice,        //CHOICE ...
			specTypeEnumerated,    //ENUMERATED ...
			specTypeSequence,      //SEQUENCE ...
			specTypeClass,         //CLASS ...
			specTypeOctetString,   //OCTET STRING ...
			specTypeBitString,     //BIT STRING ...
			specTypeInteger,       //INTEGER ...
			typeSpecNull,          //NULL
			typeSpecNumericString, //NumericString (...) ...
			typeSpecSized,         //<type-identifier> [size spec]
		))

	//	values are defined as:
	//		<identifier> <type> ::= ...
	typedValue := parser.Choice("typedValue",
		valueOperation1, //OPERATION ::= { oper1 | oper2 | ... }
		valueOperation2, //OPERATION ::= { ARGUMENT <identifier> RESULT <identifier> ... }
		valueError,      //ERROR ::= ...
		valueInteger,    //(INTEGER|<integer-type-identifier>) ::= ...
		valueBitString,  //(BIT STRING|<bit-string-type-identifier>) ::= ...
		valueSet,        //<type-identifier> ::= { ... }
	)

	specification := parser.Choice("specification",
		specType,         //::= <type specification>
		objectIdentifier, //OBJECT IDENTIFIER ::= { ... }
		typedValue,       //<typeName> ::= <value>
	)

	//each definition is written as:
	//	<identifier> <specification>
	definition := parser.Seq("definition",
		identifier,
		specification,
	)

	//EXPORTS
	//Example:
	// 	EXPORTS
	//		map-DialogueAS,
	//		MAP-DialoguePDU,
	//		map-ProtectedDialogueAS,
	//		MAP-ProtectedDialoguePDU
	//	;
	exports := parser.Seq("exports",
		parser.Keyword("EXPORTS"),
		parser.List(",", 1, identifier),
		parser.Keyword(";"),
	)

	//IMPORTS
	//	IMPORTS <names> FROM <definitionName> ... ;
	//Example:
	// IMPORTS
	// gsm-NetworkId, as-Id FROM MobileDomainDefinitions { itu-t (0) identified-organization (4) etsi (0) mobileDomain (0) mobileDomainDefinitions (0) version1 (1) }
	// AddressString FROM MAP-CommonDataTypes { itu-t identified-organization (4) etsi (0) mobileDomain (0) gsm-Network(1) modules (3) map-CommonDataTypes (18) version8 (8)}
	// ExtensionContainer FROM MAP-ExtensionDataTypes { itu-t identified-organization (4) etsi (0) mobileDomain (0) gsm-Network (1) modules (3) map-ExtensionDataTypes (21) version8 (8)}
	// SecurityHeader, ProtectedPayload FROM MAP-ST-DataTypes { itu-t identified-organization (4) etsi (0) mobileDomain (0) gsm-Network (1) modules (3) map-ST-DataTypes (27) version8 (8)}
	// ;

	commaSeparatedNames := parser.List(",", 1, identifier)

	importDefinition := parser.Seq("importDefinition",
		commaSeparatedNames,
		parser.Keyword("FROM"),
		asn1defName,
	)

	imports := parser.Seq("imports",
		parser.Keyword("IMPORTS"),
		parser.List(" ", 1, importDefinition),
		parser.Keyword(";"),
	)

	//Definitions are blocked as in this example
	//Example:
	//	DEFINITIONS IMPLICIT TAGS ::= BEGIN
	//			...
	//	END
	asn1defDefinitions := parser.Seq("asn1defDefinitions",
		parser.Keyword("DEFINITIONS"),
		parser.Optional(parser.Keyword("IMPLICIT")),
		parser.Optional(parser.Keyword("TAGS")),
		parser.Keyword("::="),
		parser.Block("definitions", "BEGIN", "END", parser.Seq("definitions",
			parser.Optional(exports),
			parser.Optional(imports),
			parser.List(" ", 0, definition)),
		),
	)

	//The full ASN.1 spec starts with the name then the definitions
	asn1def := parser.Seq("asn1def",
		asn1defName,
		asn1defDefinitions,
	)

	//loop several times over remaining lines
	//as one file may contain several asn1def's
	for lines.Count() > 0 {
		parsed, remainingLines, err := asn1def.Parse(lines)
		if err != nil {
			return log.Wrapf(err, "Failed to parse %s from file %s line %d", asn1def.Name(), filename, lines.LineNr())
		}

		log.Debugf("Parsed ASN.1 Definition: %s", parsed.Name())

		if remainingLines.LineNr() == lines.LineNr() {
			return log.Wrapf(nil, "After  parsing, same lines remain... nothing used!?!?!")
		}

		lines = remainingLines
	}
	return nil
} //definition.LoadFile()

//INotation ...
type INotation interface {
	//On success, Parse() returns remaining parsed item, remaining text and nil
	//but on error it returns nil parsed item, all the text and a descriptive error
	Parse(text string) (INotation, string, error)
}

//namedBrace parses: <name> { ... }
type namedBrace struct {
	//parser definition
	insideNotation INotation

	//parsed
	name       string
	insideItem INotation
}

func (n namedBrace) Parse(text string) (INotation, string, error) {
	if n.insideNotation == nil {
		return nil, text, log.Wrapf(nil, "namedBrace.Parse(inside=nil)")
	}

	//skip leading space
	next := strings.TrimLeft(text, " \t")

	//expect: "<name> {...}"
	parts := strings.SplitN(next, " ", 2)
	//log.Debugf("%d parts: %v", len(parts), parts)
	if len(parts) != 2 {
		return nil, text, log.Wrapf(nil, "not 2 parts")
	}
	n.name = parts[0]
	next, _ = skipOver(next, n.name)

	//expect "{...}"
	var ok bool
	if next, ok = skipOver(next, "{"); !ok {
		return nil, text, log.Wrapf(nil, "missing '{' after %s", n.name)
	}

	//parse inside...
	var err error
	if n.insideItem, next, err = n.insideNotation.Parse(next); err != nil {
		return nil, text, log.Wrapf(err, "failed to parse %s contents", n.name)
	}

	//expect "}"
	if next, ok = skipOver(next, "}"); !ok {
		return nil, text, log.Wrapf(nil, "missing '}' to end %s", n.name)
	}

	//success
	return n, next, nil
} //namedBrace.Parse()

//parses: ".. .. ... (#)"
type textBracketNumber struct {
	text   string
	number int
}

func (n textBracketNumber) Parse(text string) (INotation, string, error) {
	//skip leading space
	next := strings.TrimLeft(text, " \t")

	//expect text followed by a '('
	parts := strings.SplitN(next, "(", 2)
	//log.Debugf("%d parts: %v", len(parts), parts)
	if len(parts) != 2 {
		return nil, text, log.Wrapf(nil, "no opening bracket")
	}
	n.text = parts[0]
	next = strings.TrimLeft(parts[1], " \t") //i.e. the text after the '('

	parts = strings.SplitN(next, ")", 2)
	if len(parts) != 2 {
		return nil, text, log.Wrapf(nil, "no closing bracket")
	}

	var err error
	n.number, err = strconv.Atoi(parts[0])
	if err != nil {
		return nil, text, log.Wrapf(err, "Invalid number")
	}

	next = strings.TrimLeft(parts[1], " \t") //i.e. the text after the ')'
	return n, next, nil
}

//parses space separated values terminated by some terminator text and with minimum nr of entries
type spaceSeparatedList struct {
	//parser definition
	terminator   string
	min          int
	itemNotation INotation

	//parsed items
	items []INotation
}

func (n spaceSeparatedList) Parse(text string) (INotation, string, error) {
	next := text
	var err error
	terminatorLen := len(n.terminator)
	for {
		next = strings.TrimLeft(next, " \t")
		nextLen := len(next)
		if terminatorLen > 0 && nextLen >= terminatorLen && next[0:terminatorLen] == n.terminator {
			log.Debugf("End of spaceSeparatedList")
			break
		}

		var newItem INotation
		newItem, next, err = n.itemNotation.Parse(next)
		if err != nil {
			break
		}

		if n.items == nil {
			n.items = make([]INotation, 0)
		}
		n.items = append(n.items, newItem)
	}

	if len(n.items) < n.min {
		return nil, text, log.Wrapf(err, "%d items < min %d", len(n.items), n.min)
	}

	return n, next, nil
} //spaceSeparatedList.Parse()

//parses: "<keyword> ..."
type keywordText struct {
	//parser specification
	keyword    string
	followedBy INotation
	//output
	followingItem INotation
}

func (n keywordText) Parse(text string) (INotation, string, error) {
	next := strings.TrimLeft(text, " \t")
	var ok bool
	if next, ok = skipOver(next, n.keyword); !ok {
		return nil, text, log.Wrapf(nil, "Keyword '%s' not specified", n.keyword)
	}

	//continue with following
	if n.followedBy != nil {
		var err error
		n.followingItem, next, err = n.followedBy.Parse(next)
		if err != nil {
			return nil, text, log.Wrapf(err, "Failed to parse %s-following item", n.keyword)
		}
	}
	return n, next, nil
}

//parsed any of specified notations in order
type anyOf struct {
	//parser specification
	notations []INotation
	//output
	item INotation
}

func (n anyOf) Parse(text string) (INotation, string, error) {
	for _, itemNotation := range n.notations {
		var next string
		var err error
		n.item, next, err = itemNotation.Parse(text)
		if err == nil {
			return n, next, nil
		}
	}
	return nil, text, log.Wrapf(nil, "No item notation matched")
}

//parse sequence of notations
type seq struct {
	//parser specification
	notations []INotation
	//output
	items []INotation
}

func (n seq) Parse(text string) (INotation, string, error) {
	n.items = make([]INotation, len(n.notations))
	next := text
	for seqNr, itemNotation := range n.notations {
		var err error
		n.items[seqNr], next, err = itemNotation.Parse(next)
		if err != nil {
			return nil, text, log.Wrapf(err, "Failed to parse seq[%d]", seqNr)
		}
	}
	return n, next, nil
}

func skipOver(s string, what string) (string, bool) {
	l := len(what)
	if len(s) < l || s[0:l] != what {
		return s, false
	}
	next := strings.TrimLeft(s[l:], " \t")
	//log.Debugf("skipOver(%s) -> %.20s...", what, next)
	return next, true
}
