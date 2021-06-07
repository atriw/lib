package adt

type Key interface {
	Less(interface{}) bool
	Equal(interface{}) bool
}

type ADT interface {
	Insert(key Key, value interface{})
	Search(key Key) interface{}
	Delete(key Key) interface{}
	Length() int
}

type Validate interface {
	Validate() bool
}
