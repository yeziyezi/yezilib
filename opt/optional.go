package opt

import (
	"fmt"
	"reflect"
)

type Opt struct {
	data interface{}
}

func Empty() *Opt {
	return &Opt{data: nil}
}
func Of(data interface{}) *Opt {
	return &Opt{data: data}
}
func (opt *Opt) IfPresentT(f interface{}) {
	if opt.IsPresent() {
		opt.runFuncT(f)
	}
}

//IsEmpty 内部值是否为nil
//指针存在type不为nil,value为nil的情况，直接与nil比较结果是false，需要额外考虑
func (opt *Opt) IsEmpty() bool {
	v := reflect.ValueOf(opt.data)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.UnsafePointer {
		return v.IsNil()
	}
	return opt.data == nil
}
func (opt *Opt) IsPresent() bool {
	return !opt.IsEmpty()
}
func (opt *Opt) Or(opt2 *Opt) *Opt {
	if opt.IsPresent() {
		return opt
	}
	return opt2
}

//OrGet lazy evaluation
func (opt *Opt) OrGet(f func() *Opt) *Opt {
	if opt.IsPresent() {
		return opt
	}
	return f()
}

func (opt *Opt) OrElse(data interface{}) interface{} {
	if opt.IsPresent() {
		return opt.data
	}
	return data
}

//OrElseGet lazy evaluation
func (opt *Opt) OrElseGet(f func() interface{}) interface{} {
	if opt.IsPresent() {
		return opt.data
	}
	return f()
}

//Then lazy evaluation
func (opt *Opt) Then(f func(interface{}) interface{}) *Opt {
	if opt.IsPresent() {
		return Of(f(opt.data))
	}
	return opt
}

//ThenT f应当是func(T)V,T的类型应当根据上文推断
func (opt *Opt) ThenT(f interface{}) *Opt {
	if opt.IsPresent() {
		return Of(opt.callFuncT(f))
	}
	return opt
}
func (opt *Opt) Where(f func(interface{}) bool) *Opt {
	if opt.IsPresent() && f(opt.data) {
		return opt
	}
	return Empty()
}

//WhereT 应当是func(T)bool
func (opt *Opt) WhereT(f interface{}) *Opt {
	if opt.IsPresent() && opt.callFuncT(f).(bool) {
		return opt
	}
	return Empty()
}

func (opt *Opt) _callFuncT(f interface{}) []interface{} {
	funcT := reflect.ValueOf(f)
	params := []reflect.Value{reflect.ValueOf(opt.data)}[:]
	retValues := funcT.Call(params)
	ret := make([]interface{}, 0)
	for _, value := range retValues {
		ret = append(ret, value.Interface())
	}
	return ret
}
func (opt *Opt) callFuncT(f interface{}) interface{} {
	ret := opt._callFuncT(f)
	if len(ret) != 1 {
		panic(fmt.Sprintf("invalid return value nums %d,except 1.", len(ret)))
	}
	return ret[0]
}
func (opt *Opt) runFuncT(f interface{}) {
	opt._callFuncT(f)
}

//Get 返回实际保存的值
//使用该方法时请确保考虑到data为nil的情况
func (opt *Opt) Get() interface{} {
	return opt.data
}
