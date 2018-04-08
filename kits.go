package kits

import (
    "fmt"
    "reflect"
)

var UnassignableErr = fmt.Errorf("pkgs: Find of unassignable value. Use a pointer to interface")
var NonInterfaceErr = fmt.Errorf("pkgs: Find of non-interface value. Use a pointer to interface")
var NotFoundErr = fmt.Errorf("pkgs: matching package not found")

// Kit is a single component of code, exposing the top level API. Only the
// functions exposed by the registered Kit are accessible via find. Generally
// there should be a single Kit per Go package, however sometimes you may prefer
// to have more than one, in cases where you will want ot register them in a
// different order for overrides.
type Kit interface{}

// global registeries
var kits = []Kit{}
var uninited = []Kit{} // list of registered but yet uninitialized kits

// Register a new kit to be accessible via Find(). All registerations must
// happen upon initialization, before the call to Init() or Find()
func Register(k Kit) {
    kits = append(kits, k)
    uninited = append(uninited, k)
}

// Init all of the registered, but yet uninitialized kits, by running their
// Init() function, if exists. This is a safe place to call all of the Find()
// functions because we're guaranteed that all of the relevant kits were
// registered and are thus accessible.
func Init() {
    var k Kit
    for len(uninited) > 0 {
        k, uninited = uninited[0], uninited[1:]
        initer, ok := k.(interface{ Init() })
        if ok {
            initer.Init()
        }
    }
}

// MustFind is similar to Find but panics on error.
func MustFind(e interface{}) {
    err := Find(e)
    if err != nil {
        panic(err)
    }
}

// Find a kit that implements the provided interface. The input to this
// function must be a pointer to an interface. All of the registered kits are
// tested against the provided interface to determine if they implement it. If
// a matching kit is found, it's assigned to the pointer argument. Otherwise a
// NotFoundErr is returned.
//
// In case of conflicts, where multiple kits implements the provided interface,
// the most recently registered kit is returned. This is useful for cases where
// you don't have a preference as to which specific kits is desired, and allows
// for kit overrides.
//
// Alternatively, the argument can be a pointer to a slice of interfaces, in
// which case all of the matching kits will be assigned to the slice, allowing
// the user to devise their own approach to prioritizing which specific kit
// should be used, normally by using examining other functions exposed on the
// kits to differentiate between them (like Version(), Type(), etc.). It's up
// to the individual kits to provide the API required for such prioritization.
func Find(e interface{}) error {
    if e == nil {
        return UnassignableErr
    }

    v := reflect.ValueOf(e)
    t := v.Type()
    if t.Kind() != reflect.Ptr {
        return UnassignableErr
    }

    // pointer of..
    v, t = v.Elem(), t.Elem()
    vals, err := find(t)
    if err != nil {
        return err
    }

    lastv := vals[len(vals) - 1]
    v.Set(lastv)
    return nil
}

// finds all of the kits that implements the provided type.
func find(t reflect.Type) ([]reflect.Value, error) {
    if t.Kind() != reflect.Interface {
        return nil, NonInterfaceErr
    }

    var matched []reflect.Value
    for _, k := range kits {
        v := reflect.ValueOf(k)
        if v.Type().Implements(t) {
            matched = append(matched, v)
        }
    }

    if matched == nil {
        return nil, fmt.Errorf("kits: matching package not found: %s", t)
    }

    return matched, nil
}
