package sqwiggle

import (
	"fmt"
	"reflect"
)

// difference is a map that holds the difference between structs. The keys in
// the map are the struct field names where differences were found, and the values
// represent what was found in struct a and what was found in b.
type difference map[string]attr

// attr is a differing attribute (field) between the structs being compared
type attr struct {
	a interface{}
	b interface{}
}

var errNilInterface = fmt.Errorf("One of the interfaces is nil")

// compare compares two structs and returns a map containing the differences between them.
// This is not optimized for efficiency, and is only meant to be used in testing. It
// is not a recursive function, and structs within the given struct will be directly
// compared using reflect.DeepEqual.
func compare(a interface{}, b interface{}) (diff difference, err error) {
	// create an attribute data structure as a map of types keyed by a string.
	diff = difference(make(map[string]attr))

	if a == nil || b == nil {
		return diff, errNilInterface
	}

	// if a pointer to a struct is passed, get the type of the dereferenced object
	typA := reflect.TypeOf(a)
	aPtr := typA.Kind() == reflect.Ptr
	if aPtr {
		typA = typA.Elem()
	}

	// check that this is really a struct
	if typA.Kind() != reflect.Struct {
		return diff, fmt.Errorf("%v type not supported for comparison\n", typA.Kind())
	}

	typB := reflect.TypeOf(b)
	bPtr := typB.Kind() == reflect.Ptr
	if bPtr {
		typB = typB.Elem()
	}

	// check that B is also really a struct
	if typB.Kind() != reflect.Struct {
		return diff, fmt.Errorf("%v type not supported for comparison\n", typB.Kind())
	}

	// get the actual values of A and B, if they happen to be pointers
	valA := reflect.ValueOf(a)
	if aPtr {
		valA = valA.Elem()
	}
	valB := reflect.ValueOf(b)
	if bPtr {
		valB = valB.Elem()
	}

	// loop through struct A fields
	for i := 0; i < typA.NumField(); i++ {
		field := typA.Field(i)
		ad := getAttrDiff(field, valA, valB)
		if ad != nil {
			diff[field.Name] = *ad
		}
	}

	// now loop through struct B fields
	for i := 0; i < typB.NumField(); i++ {
		field := typB.Field(i)
		ad := getAttrDiff(field, valB, valA) // note: we switch A and B here
		if ad != nil {
			diff[field.Name] = attr{a: ad.b, b: ad.a} // now switch them back
		}
	}
	return diff, nil
}

func getAttrDiff(field reflect.StructField, valA, valB reflect.Value) *attr {
	if field.Anonymous {
		return nil
	}

	fa := valA.FieldByName(field.Name)
	fb := valB.FieldByName(field.Name)

	// if the name did not occur in struct B
	if !fb.IsValid() {
		return &attr{a: fa, b: fb}
	}

	// if the kind of A is not the same as the kind of B, we have a difference
	if !reflect.DeepEqual(fa.Kind(), fb.Kind()) {
		return &attr{a: fa, b: fb}
	}

	// if we can't interface either of the structs, and the property of being able
	// to interface is not the same between the two structs, we have a difference.
	// If they can both not be interfaced, we ignore it and continue.
	if fa.CanInterface() != fb.CanInterface() {
		return &attr{a: fa, b: fb}
	}
	if !fa.CanInterface() || !fb.CanInterface() {
		return nil
	}

	// if A is a pointer, we know by this point that B must be one too.
	// we derefence the pointer and compare the actual values.
	if fa.Kind() == reflect.Ptr {
		// check that the pointers are not nil
		if fa.IsNil() || fb.IsNil() {
			if fa.IsNil() != fb.IsNil() {
				return &attr{a: fa, b: fb}
			}
			return nil
		}

		fa = fa.Elem()
		fb = fb.Elem()
	}

	if !reflect.DeepEqual(fa.Interface(), fb.Interface()) {
		return &attr{a: fa.Interface(), b: fb.Interface()}
	}
	return nil
}
