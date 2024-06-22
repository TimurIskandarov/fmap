package fmap

import (
	"reflect"
	"strings"
	"unsafe"
)

type field struct {
	reflect.StructField
	structPath string
	parent     *field
}

type IField interface {
	GetName() string
	GetPkgPath() string
	GetType() reflect.Type
	GetTag() reflect.StructTag
	GetOffset() uintptr
	GetIndex() []int
	GetAnonymous() bool
	IsExported() bool

	// Get returns the value of the fields in the provided object.
	// It takes a parameter `obj` of type `interface{}`, representing the object.
	// It returns the value of the fields as an `interface{}`.
	Get(obj any) any

	// GetPtr returns the pointer to the field's value in the provided object.
	// It takes a parameter `obj` of type `any`, representing the pointer to object.
	// It returns the pointer to the field's value as an `any`.
	GetPtr(obj any) any
	// Set updates the value of the fields in the provided object with the provided value.
	// It takes two parameters:
	//   - obj: interface{}, representing the object pointer containing the field.
	//   - val: interface{}, representing the new value for the field.
	Set(obj any, val any)

	// GetStructPath returns the struct path of the field.
	// It returns the struct path as a string.
	GetStructPath() string

	// GetTagPath returns the path of the field's tag value with the given tag name.
	// It takes two parameters:
	//   - tag: string, representing the tag name.
	//   - ignoreParentTagMissing: bool, representing whether to ignore the missing parent tags or not.
	// It returns the tag value path as a string.
	GetTagPath(tag string, ignoreParentTagMissing bool) string
	// GetParent returns the parent field of the current field, if not exist return nil.
	GetParent() IField
}

func (f *field) GetName() string {
	return f.Name
}

func (f *field) GetPkgPath() string {
	return f.PkgPath
}

func (f *field) GetType() reflect.Type {
	return f.Type
}

func (f *field) GetTag() reflect.StructTag {
	return f.Tag
}

func (f *field) GetOffset() uintptr {
	return f.Offset
}

func (f *field) GetIndex() []int {
	return f.Index
}

func (f *field) GetAnonymous() bool {
	return f.Anonymous
}

func (f *field) IsExported() bool {
	return f.PkgPath == ""
}

func (f *field) GetStructPath() string {
	return f.structPath
}

func (f *field) GetParent() IField {
	return f.parent
}

func (f *field) GetPtr(obj interface{}) interface{} {
	return reflect.NewAt(f.Type, f.getPtr(obj)).Interface()
}

// Get returns the value of the fields in the provided object.
// It takes a parameter `obj` of type `interface{}`, representing the object.
// It returns the value of the fields as an `interface{}`.
func (f *field) Get(obj interface{}) interface{} {
	ptrToField := f.getPtr(obj)
	kind := f.Type.Kind()
	isPtr := false
	if kind == reflect.Ptr {
		isPtr = true
		kind = f.Type.Elem().Kind()
	}
	if isPtr {
		switch kind {
		case reflect.String:
			return getPtrValue[*string](ptrToField)
		case reflect.Int:
			return getPtrValue[*int](ptrToField)
		case reflect.Int8:
			return getPtrValue[*int8](ptrToField)
		case reflect.Int16:
			return getPtrValue[*int16](ptrToField)
		case reflect.Int32:
			return getPtrValue[*int32](ptrToField)
		case reflect.Int64:
			return getPtrValue[*int64](ptrToField)
		case reflect.Uint:
			return getPtrValue[*uint](ptrToField)
		case reflect.Uint8:
			return getPtrValue[*uint8](ptrToField)
		case reflect.Uint16:
			return getPtrValue[*uint16](ptrToField)
		case reflect.Uint32:
			return getPtrValue[*uint32](ptrToField)
		case reflect.Uint64:
			return getPtrValue[*uint64](ptrToField)
		case reflect.Float32:
			return getPtrValue[*float32](ptrToField)
		case reflect.Float64:
			return getPtrValue[*float64](ptrToField)
		case reflect.Bool:
			return getPtrValue[*bool](ptrToField)
		case reflect.Struct:
			return reflect.NewAt(f.Type, ptrToField).Elem().Interface()
		case reflect.Slice:
			return reflect.NewAt(f.Type, ptrToField).Elem().Interface()
		case reflect.Array:
			return reflect.NewAt(f.Type, ptrToField).Elem().Interface()
		default:
			panic("unhandled default case")
		}
	} else {
		switch kind {
		case reflect.String:
			return getPtrValue[string](ptrToField)
		case reflect.Int:
			return getPtrValue[int](ptrToField)
		case reflect.Int8:
			return getPtrValue[int8](ptrToField)
		case reflect.Int16:
			return getPtrValue[int16](ptrToField)
		case reflect.Int32:
			return getPtrValue[int32](ptrToField)
		case reflect.Int64:
			return getPtrValue[int64](ptrToField)
		case reflect.Uint:
			return getPtrValue[uint](ptrToField)
		case reflect.Uint8:
			return getPtrValue[uint8](ptrToField)
		case reflect.Uint16:
			return getPtrValue[uint16](ptrToField)
		case reflect.Uint32:
			return getPtrValue[uint32](ptrToField)
		case reflect.Uint64:
			return getPtrValue[uint64](ptrToField)
		case reflect.Float32:
			return getPtrValue[float32](ptrToField)
		case reflect.Float64:
			return getPtrValue[float64](ptrToField)
		case reflect.Bool:
			return getPtrValue[bool](ptrToField)
		case reflect.Struct:
			return reflect.NewAt(f.Type, ptrToField).Elem().Interface()
		case reflect.Slice:
			return reflect.NewAt(f.Type, ptrToField).Elem().Interface()
		case reflect.Array:
			return reflect.NewAt(f.Type, ptrToField).Elem().Interface()
		default:
			panic("unhandled default case")
		}
	}
}

func (f *field) GetTagPath(tag string, ignoreParentTagMissing bool) string {
	tagPath := ""
	if val, ok := f.Tag.Lookup(tag); ok {
		vals := strings.Split(val, ",")
		if len(vals) > 0 && len(vals[0]) > 0 {
			tagPath = vals[0]
		}
	}
	if tagPath == "" {
		return tagPath
	}
	if f.parent == nil {
		return tagPath
	}
	parentTag := f.parent.GetTagPath(tag, ignoreParentTagMissing)
	if parentTag == "" && !ignoreParentTagMissing {
		return ""
	}
	if parentTag == "" {
		return tagPath
	}
	return parentTag + "." + tagPath
}

// getPtr returns a pointer to the field's value in the provided configuration object.
// It takes a parameter `conf` of type `any`, representing the configuration object.
// It returns an `unsafe.Pointer` to the `field's` value in the configuration object.
func (f *field) getPtr(obj interface{}) unsafe.Pointer {
	confPointer := ((*[2]unsafe.Pointer)(unsafe.Pointer(&obj)))[1]
	ptToField := unsafe.Add(confPointer, f.Offset)
	return ptToField
}

func setPtrValue[T any](ptr unsafe.Pointer, val any) {
	valSet := (*T)(ptr)
	*valSet = val.(T)
}

func getPtrValue[T any](ptr unsafe.Pointer) T {
	return *(*T)(ptr)
}

// Set updates the value of the fields in the provided object with the provided value.
// It takes two parameters:
//   - obj: interface{}, representing the object containing the field.
//   - val: interface{}, representing the new value for the field.
//
// The Set method uses the getPtr method to get a pointer to the fields in the object.
// It then performs a type switch on the kind of the fields to determine its type, and sets the value accordingly.
// The supported fields types are string, int, and bool.
// If the fields type is not one of the supported types, it panics with the message "unhandled default case".
func (f *field) Set(obj interface{}, val interface{}) {
	ptrToField := f.getPtr(obj)
	kind := f.Type.Kind()
	isPtr := false
	if kind == reflect.Ptr {
		isPtr = true
		kind = f.Type.Elem().Kind()
	}
	if isPtr {
		switch kind {
		case reflect.String:
			setPtrValue[*string](ptrToField, val)
		case reflect.Int:
			setPtrValue[*int](ptrToField, val)
		case reflect.Int8:
			setPtrValue[*int8](ptrToField, val)
		case reflect.Int16:
			setPtrValue[*int16](ptrToField, val)
		case reflect.Int32:
			setPtrValue[*int32](ptrToField, val)
		case reflect.Int64:
			setPtrValue[*int64](ptrToField, val)
		case reflect.Uint:
			setPtrValue[*uint](ptrToField, val)
		case reflect.Uint8:
			setPtrValue[*uint8](ptrToField, val)
		case reflect.Uint16:
			setPtrValue[*uint16](ptrToField, val)
		case reflect.Uint32:
			setPtrValue[*uint32](ptrToField, val)
		case reflect.Uint64:
			setPtrValue[*uint64](ptrToField, val)
		case reflect.Float32:
			setPtrValue[*float32](ptrToField, val)
		case reflect.Float64:
			setPtrValue[*float64](ptrToField, val)
		case reflect.Bool:
			setPtrValue[*bool](ptrToField, val)
		default:
			dest := reflect.NewAt(f.Type, ptrToField)
			dest = dest.Elem()
			source := reflect.ValueOf(val)
			dest.Set(source)
		}
	} else {
		switch kind {
		case reflect.String:
			setPtrValue[string](ptrToField, val)
		case reflect.Int:
			setPtrValue[int](ptrToField, val)
		case reflect.Int8:
			setPtrValue[int8](ptrToField, val)
		case reflect.Int16:
			setPtrValue[int16](ptrToField, val)
		case reflect.Int32:
			setPtrValue[int32](ptrToField, val)
		case reflect.Int64:
			setPtrValue[int64](ptrToField, val)
		case reflect.Uint:
			setPtrValue[uint](ptrToField, val)
		case reflect.Uint8:
			setPtrValue[uint8](ptrToField, val)
		case reflect.Uint16:
			setPtrValue[uint16](ptrToField, val)
		case reflect.Uint32:
			setPtrValue[uint32](ptrToField, val)
		case reflect.Uint64:
			setPtrValue[uint64](ptrToField, val)
		case reflect.Float32:
			setPtrValue[float32](ptrToField, val)
		case reflect.Float64:
			setPtrValue[float64](ptrToField, val)
		case reflect.Bool:
			setPtrValue[bool](ptrToField, val)
		default:
			dest := reflect.NewAt(f.Type, ptrToField)
			dest = dest.Elem()
			source := reflect.ValueOf(val)
			dest.Set(source)
		}
	}
}
