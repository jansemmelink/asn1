package def2

import (
	"bitbucket.org/vservices/dark/logger"
	"github.com/jansemmelink/asn1/parser"
)

//IOptional defines the functions that a user implementation
//must provide to implement an optional item, e.g.
//
//	type UserList struct {
//		*def2.Optional
//	}
//
//WARNING:	If you forgot the '*' and embedded just def2.Optional,
//the Parse() method will not be inherited because it requires
//a pointer parent, so parsing will fail with:
//		panic: interface conversion: userpkg.UserList is not def2.IOptional: missing method Parse
//
//The user list type must at least provide an implementation
//for method ItemType() as in this example:
//
//	func (ul UserList) ItemType() reflect.Type {
//		return reflect.TypeOf(UserItemType)
//	}
//
//and the item type must be IParsable, e.g.:
//
//	type UserItemType def2.Int
//
//To customise the list further, add new definitions for the
//other IOptional methods, e.g. changing the default item sepatator
//with this:
//
//	func (ul UserList) Sep() string {
//		return ";"
//	}
//
//You can also make generic user list types that can be re-used
//in other lists, e.g. you have many colon separated lists, with
//3 to 5 items, but the item types are different...
//Do this by defining a base type for them:
//
//	type ColonSepThreeToFive struct { *def2.Optional }
//	func (ColonSepThreeToFive) Sep() string { return ":" }
//	func (ColonSepThreeToFive) Min() int { return 3 }
//	func (ColonSepThreeToFive) Max() int { return 5 }
//
//Now define your other lists using this as the base:
//
//	type ListOfOrders struct { *ColonSepThreeToFive }
//	func (ListOfOrders) ItemType() reflect.Type { return reflect.TypeOf {Order} }
//
//	type ListOfCustomers struct { *ColonSepThreeToFive }
//	func (ListOfCustomers) ItemType() reflect.Type { return reflect.TypeOf {Customer} }
//
// type IOptional interface {
// 	IParsable
// }

//Optional is a parsable item that returns err=nil if it cannot parse
type Optional struct{}

//Parse method is inherited by all users types that embeds
//def2.Optional which makes your item optional.
func (opt Optional) Parse(log logger.ILogger, l parser.ILines, v IParsable) (parser.ILines, error) {

	log.Debugf("Parsing optional %T", v)
	return l, log.Wrapf(nil, "NYI")

	// remain := l

	// 	var err error
	// 	//try to parse another item
	// 	//allocate memory for it:
	// 	item := mem.NewX(userList.ItemType()).(IParsable)
	// 	if remain, err = item.Parse(log, remain, item); err != nil {
	// 		//failed to parse another item
	// 		//if list has no separators - this is just the end of the list
	// 		//but if list has a separator, and we got it (count>0), then this
	// 		//is an error because the list ended with a separator but no item
	// 		//following it.
	// 		if len(opt.items) > 0 && len(userList.Sep().String()) > 0 {
	// 			return l, log.Wrapf(nil, "%T.item[%d] not found after sep='%s' on line %d: %.32s ...", userList, len(opt.items), userList.Sep(), remain.LineNr(), remain.Next())
	// 		}
	// 		log.Debugf("Could not parse %T.item[%d] from %d: %.32s ...", userList, len(opt.items), remain.LineNr(), remain.Next())
	// 		break
	// 	}

	// 	//parsed, add to list
	// 	opt.items = append(opt.items, item)

	// if len(opt.items) < userList.Min() {
	// 	return l, log.Wrapf(nil, "%T has %d items (min=%d) on line %d: %.32s ...", userList, len(opt.items), userList.Min(), l.LineNr(), l.Next())
	// }
	// if userList.Max() > 0 && len(opt.items) > userList.Max() {
	// 	return l, log.Wrapf(nil, "%T has %d items (max=%d) on line %d: %.32s ...", userList, len(opt.items), userList.Min(), l.LineNr(), l.Next())
	// }

	// log.Debugf("%T parsed %d items:", userList, len(opt.items))
	// return remain, nil
} //Optional.Parse()
