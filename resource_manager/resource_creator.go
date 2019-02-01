package resource_manager

import "reflect"

type ResourceCreator struct {
	// CreateFunc must be func() T or func() (T, error)
	// T must implements Closable.
	// If T implements Initializable, Init will be called.
	CreateFunc interface{}
	// Receiver is optional
	// If set, it must be *T
	Receiver interface{}
}

func (c ResourceCreator) validate() {
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

func (c ResourceCreator) create() (Closable, error) {
	fn := reflect.Indirect(reflect.ValueOf(c.CreateFunc))
	ret := fn.Call(nil)
	if len(ret) == 1 {
		obj, ok := ret[0].Interface().(Closable)
		if !ok {
			panic(errUnSupportCreatorFunc)
		}
		return obj, nil
	}
	if len(ret) == 2 {
		obj, ok := ret[0].Interface().(Closable)
		if !ok {
			panic(errUnSupportCreatorFunc)
		}
		var err error
		if !ret[1].IsNil() {
			err = ret[1].Interface().(error)
		}
		return obj, err
	}
	panic(errUnSupportCreatorFunc)
}

func (c ResourceCreator) setReceiver(obj Closable) {
	if c.Receiver != nil {
		reflect.ValueOf(c.Receiver).Elem().Set(reflect.ValueOf(obj))
	}
}
