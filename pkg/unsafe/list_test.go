package unsafe

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

type List struct {
	_type *sliceType
	list unsafe.Pointer
	length int
	cap int
	slice *Slice
}

func (list *List) ElemType() string {
	return list._type.elem.Type()
}

func (list *List) Last() int {
	return list.length
}

func (list *List) Close() error {
	tmp := make([]*rtype, list.length)
	copy(tmp, *(*[]*rtype)(list.list))
	*(*[]*rtype)(list.list) = tmp
	return nil
}

func (list *List) Append(v interface{}) error {
	value := ValueOf(v)
	if value.Type() != list.ElemType() {
		return fmt.Errorf(
			"cannot append type of %s to list of type %s",
			value.Type(), list.ElemType(),
		)
	}
	if list.length >= list.cap {
		list.cap *= 2
		tmp := make([]*rtype, list.cap)
		copy(tmp, *(*[]*rtype)(list.list))
		*(*[]*rtype)(list.list) = tmp
	}

	idx := getSliceIndex(list.slice, list._type, list.length)
	set(value.typ.Kind(), idx, value.ptr)
	list.length += 1
	return nil
}

func (list *List) Sort() error {
	kind := list._type.elem.Kind()
	for i := 0; i < list.length-1; i++ {
		for j := 0; j < list.length - i - 1; j++ {
			a, b := list.get(j), list.get(j+1)
			if cmp(kind, a, b) {
				tmp := unsafe.Pointer(&rtype{})
				set(kind, tmp, a)
				set(kind, a, b)
				set(kind, b, tmp)
			}
		}
	}
	return nil
}

func cmp(kind reflect.Kind, a, b unsafe.Pointer) bool {
	switch kind {
	case reflect.Int:
		return *(*int)(a) > *(*int)(b)
	}
	return false
}

func (list *List) get(i int) unsafe.Pointer {
	return getSliceIndex(list.slice, list._type, i)
}

func (list *List) Set(i int, v interface{}) error {
	value := ValueOf(v)
	if value.Type() != list.ElemType() {
		return fmt.Errorf(
			"cannot append type of %s to list of type %s",
			value.Type(), list.ElemType(),
		)
	}
	if i >= list.length {
		return fmt.Errorf("cannot set index %d: out of range: %d", i, list.length)
	}
	idx := getSliceIndex(list.slice, list._type, i)
	set(value.typ.Kind(), idx, value.ptr)
	return nil
}

func set(kind reflect.Kind, a, b unsafe.Pointer) {
	switch kind {
	case reflect.Int:
		*(*int)(a) = *(*int)(b)
	case reflect.Struct:
		*(*rtype)(a) = *(*rtype)(b)
	}
}

func WrapList(v interface{}) (*List, error) {
	value := ValueOf(v)
	if value.typ.Kind() != reflect.Slice && value.typ.Kind() != reflect.Array {
		return nil, fmt.Errorf("cannot wrap type of %s as list", value.typ.Kind())
	}
	var tt *sliceType
	if value.typ.Kind() == reflect.Ptr {
		ptr := (*ptrType)(unsafe.Pointer(value.typ))
		tt = (*sliceType)(unsafe.Pointer(ptr.elem))
	} else {
		tt = (*sliceType)(unsafe.Pointer(value.typ))
	}

	s := (*Slice)(value.ptr)
	cap := 3
	*(*[]*rtype)(value.ptr) = make([]*rtype, cap)
	return &List{
		_type: tt,
		list: value.ptr,
		length: 0,
		cap: cap,
		slice: s,
	}, nil
}

func catch(t *testing.T, errs ...error) {
	for _, err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
}

type Equality uint8

const (
	Error Equality = iota
	Equal
	Greater
	Lesser
)


type Comparer interface {
	Compare(interface{}) Equality
}

type Rectangle struct {
	X, Y int
}

func (r Rectangle) Area() int {
	return r.X * r.Y
}

func (r Rectangle) Compare(v interface{}) Equality {
	value, ok := v.(Rectangle)
	if !ok {
		return Error
	}
	a, b := r.Area(), value.Area()
	if a == b {
		return Equal
	}
	if a > b {
		return Greater
	} else {
		return Lesser
	}
}

type Integer struct {
	value int
}

func (integer Integer) Compare(v interface{}) Equality {
	value, ok := v.(Integer)
	if !ok {
		return Error
	}

	if value.value == integer.value {
		return Equal
	}
	if integer.value > value.value {
		return Greater
	} else {
		return Lesser
	}
}

func TestInterface(t *testing.T) {
	val := ValueOf(Integer{})
	fmt.Println(val.Type())
	u := val.typ.uncommon()

	for _, method := range u.methods() {
		fmt.Println(method)
		fmt.Println(resolveNameOff(unsafe.Pointer(u), int32(method.name)).name())
		fmt.Println(resolveTypeOff(unsafe.Pointer(u), method.mtyp))
	}
}

func TestList(t *testing.T) {
	var arr []int
	l, err := WrapList(arr)
	if err != nil {
		t.Fatal(err)
	}
	_ = l
	fmt.Println(arr)
	catch(t,
		l.Append(4),
		l.Append(2),
		l.Append(1),
		l.Append(3),
		l.Close(),
		)
	fmt.Println(arr)
	l.Sort()
	fmt.Println(arr)
}