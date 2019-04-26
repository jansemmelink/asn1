package def2

import (
	"reflect"

	"bitbucket.org/vservices/dark/logger"
	"github.com/jansemmelink/asn1/parser"
	"github.com/jansemmelink/mem"
)

//IList defines the functions that a user implementation
//must provide to implement a list
//Note that Sep(), Min() and Max() has default implementation
//that the userList will inherit when defined as:
//
//	type UserList struct {
//		*def2.List
//	}
//
//WARNING:	If you forgot the '*' and embedded just def2.List,
//the Parse() method will not be inherited because it requires
//a pointer parent, so parsing will fail with:
//		panic: interface conversion: userpkg.UserList is not def2.IList: missing method Parse
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
//other IList methods, e.g. changing the default item sepatator
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
//	type ColonSepThreeToFive struct { *def2.List }
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
type IList interface {
	IParsable
	Sep() IToken
	Min() int
	Max() int
	ItemType() reflect.Type
}

//List is a parsable sequence that can repeat after a separator
type List struct {
	items []IParsable
}

//Sep defaults to a comma "," token
func (List) Sep() IToken {
	return NewToken(",")
}

//Min defaults to 0 to allow empty list
func (List) Min() int {
	return 0
}

//Max defaults to any valye <= 0 to any nr of elements
func (List) Max() int {
	return -1
}

//Items method is not part of the interface IList because
//it is not required for List.Parse() to work, but the user
//will inherit it by embedding *def2.List and it can be used
//to access the list of parsed items.
func (lst List) Items() []IParsable {
	return lst.items
}

//Parse method is inherited by all users lists that embed
//*def2.List which makes your lists parsable. But parsing
//will fail if your list type embeds def2.List instead of *def2.List!
func (lst *List) Parse(log logger.ILogger, l parser.ILines, v IParsable) (parser.ILines, error) {
	remain := l

	//lst is the base List definition, that does not implement
	//the user's customization method ItemType()
	//we get the user implementation from argument v:
	var userList IList
	userList = reflect.ValueOf(v).Elem().Interface().(IList)
	log.Debugf("=====[ PARSING LIST %T: sep='%s', min=%d, max=%d, item=%s ]=====", userList, userList.Sep(), userList.Min(), userList.Max(), userList.ItemType().Name())

	//start with empty list
	lst.items = make([]IParsable, 0)
	for userList.Max() <= 0 || len(lst.items) < userList.Max() {
		var err error

		//if not first, expect optional separator before next item
		if len(lst.items) > 0 && len(userList.Sep().String()) > 0 {
			if remain, err = userList.Sep().Parse(log, remain, nil); err != nil {
				log.Debugf("%T no sep='%s' after item[%d] on line %d: %.32s ...", userList, userList.Sep(), len(lst.items), remain.LineNr(), remain.Next())
				break
			}
		} //if not first and expects separator between items

		//try to parse another item
		//allocate memory for it:
		item := mem.NewX(userList.ItemType()).(IParsable)
		if remain, err = item.Parse(log, remain, item); err != nil {
			//failed to parse another item
			//if list has no separators - this is just the end of the list
			//but if list has a separator, and we got it (count>0), then this
			//is an error because the list ended with a separator but no item
			//following it.
			if len(lst.items) > 0 && len(userList.Sep().String()) > 0 {
				return l, log.Wrapf(nil, "%T.item[%d] not found after sep='%s' on line %d: %.32s ...", userList, len(lst.items), userList.Sep(), remain.LineNr(), remain.Next())
			}
			log.Debugf("Could not parse %T.item[%d] from %d: %.32s ...", userList, len(lst.items), remain.LineNr(), remain.Next())
			break
		}

		//parsed, add to list
		lst.items = append(lst.items, item)
	} //while count < max

	if len(lst.items) < userList.Min() {
		return l, log.Wrapf(nil, "%T has %d items (min=%d) on line %d: %.32s ...", userList, len(lst.items), userList.Min(), l.LineNr(), l.Next())
	}
	if userList.Max() > 0 && len(lst.items) > userList.Max() {
		return l, log.Wrapf(nil, "%T has %d items (max=%d) on line %d: %.32s ...", userList, len(lst.items), userList.Min(), l.LineNr(), l.Next())
	}

	log.Debugf("%T parsed %d items:", userList, len(lst.items))
	return remain, nil
} //List.Parse()
