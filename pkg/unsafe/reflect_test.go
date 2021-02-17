package unsafe

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

const (
	kindMask = (1 << 5) - 1

	flagKindMask    flag = 1<<flagKindWidth - 1
	flagKindWidth        = 5 // there are 27 kinds
	flagStickyRO = 1 << 5
	flagEmbedRO = 1 << 6
	flagIndir = 1 << 7
	flagMethod = 1 << 9
	flagMethodShift = 10
	flagRO = flagStickyRO | flagEmbedRO
)

type Value struct {
	typ *rtype
	ptr unsafe.Pointer
	flag
}

type flag uintptr

func (f flag) kind() reflect.Kind {
	return reflect.Kind(f & flagKindMask)
}

func (f flag) ro() flag {
	if f&flagRO != 0 { // if flags at << 5 & << 6 is 1
		return flagStickyRO // return 1 at << 5
	}
	return 0
}

type Person struct {
	Name string `json:"name"`
	Age int `json:"age"`
	Twitter string `json:"twitter"`
}

func (p *Person) Hi() {
	fmt.Println("hello there")
}

func TestMe(t *testing.T) {
	person := &Person{Name: "Lasse", Age: 23, Twitter: "ifndef_lmj"}
	fmt.Println(person)
	v := ValueOf(person)
	//v = Method(v, 0)
	Name(v, 0)
	fmt.Println(person)
}

//go:linkname firstmoduledata runtime.firstmoduledata
var firstmoduledata Moduledata

type Moduledata struct {
	pclntable    []byte
	ftab         []Functab
	filetab      []uint32
	findfunctab  uintptr
	minpc, maxpc uintptr

	text, etext           uintptr
	noptrdata, enoptrdata uintptr
	data, edata           uintptr
	bss, ebss             uintptr
	noptrbss, enoptrbss   uintptr
	end, gcdata, gcbss    uintptr
	types, etypes uintptr

	// Original type was []*_type
	typelinks []interface{}

	modulename string
	// Original type was []modulehash
	modulehashes []interface{}

	gcdatamask, gcbssmask Bitvector

	next *Moduledata
}

type Functab struct {
	entry   uintptr
	funcoff uintptr
}

type Bitvector struct {
	n        int32 // # of bits
	bytedata *uint8
}

func ValueOf(v interface{}) Value {
	val := reflect.ValueOf(v)
	//val.NumField()
	//val.Method(0)
	//val.Type().Field(0)
	val.Type().Name()
	t := (*emptyInterface)(unsafe.Pointer(&v))
	f := t.typ.kind & kindMask
	return Value{t.typ, t.word, flag(f)}
}

func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

type String struct {
	Data unsafe.Pointer
	Len  int
}

type name struct {
	bytes *byte
}

func (n name) tag() string {
	tl := n.tagLen()
	if tl == 0 {
		return ""
	}
	nl := n.nameLen()
	var b string
	hdr := (*stringStruct)(unsafe.Pointer(&b))
	hdr.str = unsafe.Pointer(n.data(3+nl+2))
	hdr.len = tl
	return b
}

func (n name) tagLen() int {
	if *n.data(0) & (1<<1) == 0 {
		return 0
	}
	off := 3 + n.nameLen()
	return int(uint16(*n.data(off)) << 8 | uint16(*n.data(off+1)))
}

func (n name) nameLen() int {
	return int(uint16(*n.data(1))<<8 | uint16(*n.data(2)))
}

func (n name) data(off int) *byte {
	return (*byte)(add(unsafe.Pointer(n.bytes), uintptr(off)))
}


func (n name) name() string {
	nl := int(uint16(*n.data(1))<<8 | uint16(*n.data(2)))
	var b string
	hdr := (*stringStruct)(unsafe.Pointer(&b))
	hdr.str = unsafe.Pointer(n.data(3))
	hdr.len = nl
	return b
}

func Name(val Value, i int) Value {
	field := resolveNameOff(unsafe.Pointer(val.typ), int32(val.typ.str)	)

	fmt.Println(field.name())
	// check index

	Field(val, i)
	return Value{}
}

func Field(val Value, i int) Value {
	// check the kind and ensure that it's not pointer ?
	ptrtt := (*ptrType)(unsafe.Pointer(val.typ)).elem
	tt := (*structType)(unsafe.Pointer(ptrtt))

	for i := 0; i < len(tt.fields); i++ {
		f := &tt.fields[i]
		ptr := add(val.ptr, f.offset())
		switch reflect.Kind(f.typ.kind & kindMask) {
		case reflect.String:
			fmt.Println(f.name.name(), resolveNameOff(unsafe.Pointer(f.typ), int32(f.typ.str)).name(), *(*string)(ptr), f.name.tag())
		case reflect.Int:
			*(*int)(ptr) = int(31)
			fmt.Println(f.name.name(), resolveNameOff(unsafe.Pointer(f.typ), int32(f.typ.str)).name(), *(*int64)(ptr), f.name.tag())
		default:
			fmt.Println(f.name.name(), resolveNameOff(unsafe.Pointer(f.typ), int32(f.typ.str)).name(), f.name.tag())
		}
	}
	return Value{}
}

// ptrType represents a pointer type.
type ptrType struct {
	rtype
	elem *rtype // pointer element (pointed at) type
}

type hex uint64

type stringStruct struct {
	str unsafe.Pointer
	len int
}

func resolveNameOff(ptr unsafe.Pointer, off int32) name {
	base := uintptr(ptr)
	for md := &firstmoduledata; md != nil; md = md.next {
		if base >= md.types && base < md.etypes {
			res := md.types + uintptr(off)
			if res > md.etypes {
				println("runtime: nameOff", hex(off), "out of range", hex(md.types), "-", hex(md.etypes))
				panic("runtime: Name offset out of range")
			}
			return name{(*byte)(unsafe.Pointer(res))}
		}
	}
	panic("oh no")
	return name{}
}

func Method(val Value, i int) Value {
	// check if actual struct with method
	fl := val.flag.ro() | (val.flag & flagIndir)
	fl |= flag(reflect.Func)
	// shift i (index of method) << 10 and set << 9'th bit
	fl |= flag(i)<<flagMethodShift | flagMethod
	return Value{ val.typ, val.ptr, fl}
}

type structType struct {
	rtype
	pkgPath name
	fields  []structField // sorted by offset
}

type structField struct {
	name        name    // Name is always non-empty
	typ         *rtype  // type of Name
	offsetEmbed uintptr // byte offset of Name<<1 | isEmbedded
}
func (f *structField) offset() uintptr {
	return f.offsetEmbed >> 1
}

type emptyInterface struct {
	typ  *rtype
	word unsafe.Pointer
}

type nameOff int32
type typeOff int32
type tflag uint8

type rtype struct {
	size       uintptr
	ptrdata    uintptr // number of bytes in the type that can contain pointers
	hash       uint32  // hash of type; avoids computation in hash tables
	tflag      tflag   // extra type information flags
	align      uint8   // alignment of variable with this type
	fieldAlign uint8   // alignment of struct Name with this type
	kind       uint8   // enumeration for C
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal     func(unsafe.Pointer, unsafe.Pointer) bool
	gcdata    *byte   // garbage collection data
	str       nameOff // string form
	ptrToThis typeOff // type for pointer to this type, may be zero
}
