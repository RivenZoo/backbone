package closure

import (
	"errors"
	"reflect"
)

var errUnSupportCreatorFunc = errors.New("resource creator must be func() T or func() (T, error)")

type ObjectCreator struct {
	// CreateFunc must be func() T or func() (T, error)
	CreateFunc interface{}
	// Receiver is optional
	// If set, it must be *T
	Receiver interface{}
}

func (c ObjectCreator) Validate() {
	v := reflect.Indirect(reflect.ValueOf(c.CreateFunc))
	if v.Type().Kind() != reflect.Func {
		panic(errUnSupportCreatorFunc)
	}
	t := v.Type()
	if t.NumIn() != 0 {
		panic(errUnSupportCreatorFunc)
	}
	if t.NumOut() <= 0 || t.NumOut() > 2 {
		panic(errUnSupportCreatorFunc)
	}
	if t.NumOut() == 2 {
		if t.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
			panic(errUnSupportCreatorFunc)
		}
	}
}

func (c ObjectCreator) Create() (interface{}, error) {
	fn := reflect.Indirect(reflect.ValueOf(c.CreateFunc))
	ret := fn.Call(nil)
	if len(ret) == 1 {
		return ret[0].Interface(), nil
	}
	if len(ret) == 2 {
		var err error
		if !ret[1].IsNil() {
			err = ret[1].Interface().(error)
		}
		return ret[0].Interface(), err
	}
	panic(errUnSupportCreatorFunc)
}

func (c ObjectCreator) SetReceiver(obj interface{}) {
	if c.Receiver != nil {
		reflect.ValueOf(c.Receiver).Elem().Set(reflect.ValueOf(obj))
	}
}
