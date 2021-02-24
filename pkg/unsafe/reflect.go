package unsafe

import (
	"errors"
	"fmt"
	"github.com/Pungyeon/required/pkg/lexer"
	"github.com/Pungyeon/required/pkg/token"
	"reflect"
	"strconv"
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

	tflagExtraStar tflag = 1 << 1
	tflagUncommon tflag = 1 << 0
)

type Value struct {
	typ *rtype
	ptr unsafe.Pointer
	flag
}

func (v Value) Type() string {
	return v.typ.Type()
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

//go:linkname firstmoduledata runtime.firstmoduledata
var firstmoduledata Moduledata

type Moduledata struct {
	pclntable    []byte
	ftab         []struct{
		entry uintptr
		funcoff uintptr
	}
	filetab      []uint32
	findfunctab  uintptr
	minpc, maxpc uintptr

	text, etext           uintptr
	noptrdata, enoptrdata uintptr
	data, edata           uintptr
	bss, ebss             uintptr
	noptrbss, enoptrbss   uintptr
	end, gcdata, gcbss    uintptr
	types, etypes         uintptr

	textsectmap []struct {
		vaddr    uintptr // prelinked section vaddr
		length   uintptr // section length
		baseaddr uintptr // relocated section address
	}
	typelinks   []int32 // offsets from types
	itablinks   []*itab

	ptab []ptabEntry

	pluginpath string
	pkghashes  []modulehash

	modulename   string
	modulehashes []modulehash

	hasmain uint8 // 1 if module contains the main function, 0 otherwise

	gcdatamask, gcbssmask bitvector

	typemap map[typeOff]*rtype // offset to *_rtype in previous module

	bad bool // module failed to load and should be ignored

	next *Moduledata

}

type ptabEntry struct {
	name nameOff
	typ  typeOff
}

// ../cmd/compile/internal/gc/reflect.go:/^func.dumptabs.
type itab struct {
	inter *interfaceType
	_type *rtype
	hash  uint32 // copy of _type.hash. Used for type switches.
	_     [4]byte
	fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}

type modulehash struct {
	modulename   string
	linktimehash string
	runtimehash  *string
}

type Functab struct {
	entry   uintptr
	funcoff uintptr
}

type bitvector struct {
	n        int32 // # of bits
	bytedata *uint8
}

func ValueOf(v interface{}) Value {
	t := (*emptyInterface)(unsafe.Pointer(&v))
	f := t.typ.kind & kindMask
	return Value{t.typ, t.word, flag(f)}
}

func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
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

type Tag struct {
	FieldIndex        int
	FieldName         string
	Required          bool
	OmitIfEmpty       bool
	IsSet             bool
	RequiredInterface bool
}


func (t *Tag) addValue(value string) error {
	if t.FieldName == "" {
		t.FieldName = value
		return nil
	}
	switch value {
	case "required":
		t.Required = true
	case "omitifempty":
		t.OmitIfEmpty = true
	default:
		return fmt.Errorf("illegal tag value: `%s`", value)
	}
	return nil
}

var typeCache = map[string]map[string]Tag{}

func GetTags(val Value) (map[string]Tag, error) {
	tags, ok := typeCache[val.Type()]
	if ok {
		return tags, nil
	}
	var tt *structType
	if val.typ.Kind() == reflect.Ptr {
		ptrtt := (*ptrType)(unsafe.Pointer(val.typ)).elem
		tt = (*structType)(unsafe.Pointer(ptrtt))
	} else {
		tt = (*structType)(unsafe.Pointer(val.typ))
	}

	tags = make(map[string]Tag, 0)
	for i := 0; i < len(tt.fields); i++ {
		f := &tt.fields[i]

		t := reflect.StructTag(f.name.tag())
		jsonTag, ok := t.Lookup("json")
		if !ok {
			n := f.name.name()
			tags[toSnakeCase(n)] = Tag{ FieldIndex: i, FieldName: n }
		} else {
			tag, err := fromString(jsonTag, i)
			if err != nil {
				return tags, err
			}
			tags[jsonTag] = tag
		}
	}
	typeCache[val.Type()] = tags
	return tags, nil
}

var diff uint8 = 'a' - 'A'

func toSnakeCase(s string) string {
	var result string
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A'-1 && s[i] <= 'Z' {
			if i > 0 {
				result += "_"
			}
			result += string(s[i] + diff)
		} else {
			result += string(s[i])
		}
	}
	return result
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
	var tt *structType
	if val.typ.Kind() == reflect.Ptr {
		ptrtt := (*ptrType)(unsafe.Pointer(val.typ)).elem
		tt = (*structType)(unsafe.Pointer(ptrtt))
	} else {
		tt = (*structType)(unsafe.Pointer(val.typ))
	}

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

type sliceType struct {
	rtype
	elem *rtype // slice element type
}

type interfaceType struct {
	rtype
	pkgPath name      // import path
	methods []imethod // sorted by hash
}

type imethod struct {
	name nameOff // name of method
	typ  typeOff // .(*FuncType) underneath
}

type Slice struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

type hex uint64

type stringStruct struct {
	str unsafe.Pointer
	len int
}

func resolveTypeOff(ptrInModule unsafe.Pointer, off typeOff) *rtype {
	if off == 0 {
		return nil
	}
	base := uintptr(ptrInModule)
	var md *Moduledata
	for next := &firstmoduledata; next != nil; next = next.next {
		if base >= next.types && base < next.etypes {
			md = next
			break
		}
	}
	if md == nil {
		panic("it's not worth it")
	}
	if t := md.typemap[off]; t != nil {
		return t
	}
	res := md.types + uintptr(off)
	if res > md.etypes {
		println("runtime: typeOff", hex(off), "out of range", hex(md.types), "-", hex(md.etypes))
	}
	return (*rtype)(unsafe.Pointer(res))
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
type textOff int32
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

func (r *rtype) Kind() reflect.Kind {
	return reflect.Kind(r.kind & kindMask)
}

func (r *rtype) Type() string {
	s := resolveNameOff(unsafe.Pointer(r), int32(r.str)).name()
	if r.tflag&tflagExtraStar != 0 {
		return s[1:]
	}
	return s
}

func (r *rtype) nameOff(off nameOff) string {
	return resolveNameOff(unsafe.Pointer(r), int32(off)).name()
}

func (r *rtype) typeOff(off typeOff) *rtype {
	return resolveTypeOff(unsafe.Pointer(r), off)
}

func (r *rtype) uncommon() *uncommonType {
	if r.tflag&tflagUncommon == 0 {
		return nil
	}
	switch r.Kind() {
	case reflect.Struct:
		return &(*structTypeUncommon)(unsafe.Pointer(r)).u
	default:
		type u struct {
			rtype
			u uncommonType
		}
		return &(*u)(unsafe.Pointer(r)).u
	}
}

type uncommonType struct {
	pkgPath nameOff // import path; empty for built-in types like int, string
	mcount  uint16  // number of methods
	xcount  uint16  // number of exported methods
	moff    uint32  // offset from this uncommontype to [mcount]method
	_       uint32  // unused
}

func (t *uncommonType) methods() []method {
	if t.mcount == 0 {
		return nil
	}
	return (*[1 << 16]method)(add(unsafe.Pointer(t), uintptr(t.moff)))[:t.mcount:t.mcount]
}

// Method on non-interface type
type method struct {
	name nameOff // name of method
	mtyp typeOff // method type (without receiver)
	ifn  textOff // fn used in interface call (one-word receiver)
	tfn  textOff // fn used for normal method call
}

type structTypeUncommon struct {
	structType
	u uncommonType
}

func fromString(input string, index int) (Tag, error) {
	tag := Tag{
		FieldIndex: index,
	}
	var previous int
	current := indexOfNextTag(input, 0)
	for current < len(input) {
		if input[current] == ',' {
			if err := tag.addValue(input[previous:current]); err != nil {
				return tag, err
			}
			current = indexOfNextTag(input, current)
			previous = current
		}
		current++
	}
	if previous != current {
		return tag, tag.addValue(input[previous:current])
	}
	return tag, nil
}

func indexOfNextTag(input string, current int) int {
	for input[current] == ' ' ||
		input[current] == '\n' ||
		input[current] == '\t' ||
		input[current] == ',' ||
		input[current] == '\r' {
		current++
	}
	return current
}

func Unmarshal(data []byte, v interface{}) error {
	p := &parser{ lexer: lexer.NewLexer(data) }
	return p.parse(ValueOf(v))
}

type parser struct {
	lexer *lexer.Lexer
}

func (p *parser) parse(val Value) error {
	// if interface do some stuff here
	for {
		tkn, err := p.lexer.Next()
		if err != nil {
			return err
		}
		return p.setValue(val, tkn)
	}
}

func (p *parser) setValue(val Value, tkn token.Token) error {
	switch tkn.Type {
	case token.OpenCurly:
		return p.parseObject(val)
	case token.OpenBrace:
		return p.parseArray(val)
	case token.Boolean:
		b, err := tkn.Bool()
		if err != nil {
		return err
	}
		*(*bool)(val.ptr) = b
	case token.String:
		*(*string)(val.ptr) = tkn.ToString()
	case token.Integer:
		var tt *rtype
		if val.typ.Kind() == reflect.Ptr {
		tt = (*ptrType)(unsafe.Pointer(val.typ)).elem
	} else {
		tt = val.typ
	}
		i, err := strconv.ParseInt(*(*string)(unsafe.Pointer(&tkn.Value)), 10, 64)
		if err != nil {
		return err
	}
		SetInt(val.ptr, tt.Kind(), i)
	case token.Float:
		var tt *rtype
		if val.typ.Kind() == reflect.Ptr {
		tt = (*ptrType)(unsafe.Pointer(val.typ)).elem
		} else {
			tt = val.typ
		}
			f, err := token.Ttof(tkn)
			if  err != nil {
			return err
		}
			// maybe check if it's float at all?
			if tt.Kind() == reflect.Float32 {
			*(*float32)(val.ptr) = float32(f)
		} else {
			*(*float64)(val.ptr) = f
		}
	}
	return nil
}

func (p *parser) parseArray(val Value) error {
	var tt *sliceType
	if val.typ.Kind() == reflect.Ptr {
		ptr := (*ptrType)(unsafe.Pointer(val.typ))
		tt = (*sliceType)(unsafe.Pointer(ptr.elem))
	} else {
		tt = (*sliceType)(unsafe.Pointer(val.typ))
	}

	s := (*Slice)(val.ptr)
	length := 3
	*(*[]*rtype)(val.ptr) = make([]*rtype, length)
	var i int
	for {
		tkn, err := p.lexer.Next()
		if err != nil {
			return err
		}
		if tkn.Type == token.ClosingBrace {
			return nil
		}
		if tkn.Type == token.Comma {
			continue
		}
		if i >= length {
			length *= 2
			tmp := make([]*rtype, length)
			copy(tmp, *(*[]*rtype)(val.ptr))
			*(*[]*rtype)(val.ptr) = tmp
		}
		index := getSliceIndex(s, tt, i)
		if err := p.setValue(Value{tt.elem, index, val.flag}, tkn); err != nil {
			return err
		}
		i++
	}
}

func grow(arr []*rtype, length int, index int) []*rtype {
	if index >= length {
		return make([]*rtype, length*2)
	}
	return arr
}


func getSliceIndex(slice *Slice, tt *sliceType, index int) unsafe.Pointer {
	return add(slice.Data, uintptr(index)*tt.elem.size)
}

func (p *parser) parseObject(val Value) error {
	tags, err := GetTags(val)
	if err != nil {
		return err
	}
	var tt *structType
	if val.typ.Kind() == reflect.Ptr {
		ptrtt := (*ptrType)(unsafe.Pointer(val.typ)).elem
		tt = (*structType)(unsafe.Pointer(ptrtt))
	} else {
		tt = (*structType)(unsafe.Pointer(val.typ))
	}
	for {
		tkn, err := p.lexer.Next()
		if err != nil {
			return err
		}
		if tkn.Type != token.String {
			return errors.New("expected field string")
		}
		field := *(*string)(unsafe.Pointer(&tkn.Value))
		tkn, err = p.lexer.Next()
		if err != nil {
			return err
		}

		if tkn.Type != token.Colon {
			return errors.New("expected comma")
		}
		tag, ok := tags[field]
		if !ok {
			return errors.New("no such field")
		}

		f := &tt.fields[tag.FieldIndex]
		ptr := add(val.ptr, f.offset())

		tkn, err = p.lexer.Next()
		if err != nil {
			return err
		}
		switch f.typ.Kind() {
		case reflect.Map:
		case reflect.Array, reflect.Slice:
			if err := p.parseArray(Value{f.typ, ptr, val.flag}); err != nil {
				return err
			}
		case reflect.Ptr:
		case reflect.Bool:
			b, err := tkn.Bool()
			if err != nil {
				return err
			}
			*(*bool)(ptr) = *(*bool)(unsafe.Pointer(&b))
		case reflect.String:
			*(*string)(ptr) = *(*string)(unsafe.Pointer(&tkn.Value))
		case reflect.Int, reflect.Float32, reflect.Float64,
			reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				i, err := strconv.ParseInt(*(*string)(unsafe.Pointer(&tkn.Value)), 10, 64)
				if err != nil {
					return err
				}
				SetInt(ptr, f.typ.Kind(), i)
		case reflect.Struct:
			if tkn.Type != token.OpenCurly {
				return errors.New("object must start with {")
			}
			if err := p.parseObject(Value{f.typ, ptr, val.flag}); err != nil {
				return err
			}
		default:
			return errors.New("unsupported type")
		}

		tkn, err = p.lexer.Next()
		if err != nil {
			return err
		}
		if tkn.Type != token.Comma {
			if tkn.Type != token.ClosingCurly {
				return fmt.Errorf("expected curly: %s", p.lexer.Previous())
			}
			return nil
		}
	}
}

func SetInt(ptr unsafe.Pointer, kind reflect.Kind, val int64) {
	switch kind {
	case reflect.Int:
		*(*int)(ptr) = int(val)
	case reflect.Int8:
		*(*int8)(ptr) = int8(val)
	case reflect.Int16:
		*(*int16)(ptr) = int16(val)
	case reflect.Int32:
		*(*int32)(ptr) = int32(val)
	case reflect.Int64:
		*(*int64)(ptr) = val
	case reflect.Uint:
		*(*uint)(ptr) = uint(val)
	case reflect.Uint8:
		*(*uint8)(ptr) = uint8(val)
	case reflect.Uint16:
		*(*uint16)(ptr) = uint16(val)
	case reflect.Uint32:
		*(*uint32)(ptr) = uint32(val)
	case reflect.Uint64:
		*(*uint64)(ptr) = uint64(val)
	case reflect.Interface:
		*(*interface{})(ptr) = int(val)
	default:
		panic(fmt.Sprintf("cannot set integer value of kind: %v", kind))
	}
}





