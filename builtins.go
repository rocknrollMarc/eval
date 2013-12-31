package eval

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	intType reflect.Type = reflect.TypeOf(int(0))
	f32 reflect.Type = reflect.TypeOf(float32(0))
	f64 reflect.Type = reflect.TypeOf(float64(0))
	c64 reflect.Type = reflect.TypeOf(complex64(0))
	c128 reflect.Type = reflect.TypeOf(complex128(0))
)

// For each parameter in a builtin function, a bool parameter is passed
// to indicate if the corresponding value is typed. The boolean(s) appear
// after entire builtin parameter list.
//
// Builtin functions must return the builtin function reflect.Value, a
// bool indicating if the return value is typed, and an error if there was one.
// The returned Value must be valid
var builtinFuncs = map[string] reflect.Value {
	"complex": reflect.ValueOf(func(r, i reflect.Value, rt, it bool) (reflect.Value, bool, error) {
		rr, rerr := assignableValue(r, f64, rt)
		ii, ierr := assignableValue(i, f64, it)
		if rerr == nil && ierr == nil {
			return reflect.ValueOf(complex(rr.Float(), ii.Float())), rt || it, nil
		}

		rr, rerr = assignableValue(r, f32, rt)
		ii, ierr = assignableValue(i, f32, it)
		if rerr == nil && ierr == nil {
			return reflect.ValueOf(complex64(complex(rr.Float(), ii.Float()))), rt || it, nil
		}
		return reflect.Zero(c128), false, ErrBadComplexArguments{r, i}
	}),
	"real": reflect.ValueOf(func(z reflect.Value, zt bool) (reflect.Value, bool, error) {
		if zz, err := assignableValue(z, c128, zt); err == nil {
			return reflect.ValueOf(real(zz.Complex())), zt, nil
		} else if zz, err := assignableValue(z, c64, zt); err == nil {
			return reflect.ValueOf(float32(real(zz.Complex()))), zt, nil
		} else {
			return reflect.Zero(f64), false, ErrBadBuiltinArgument{"real", z}
		}
	}),
	"imag": reflect.ValueOf(func(z reflect.Value, zt bool) (reflect.Value, bool, error) {
		if zz, err := assignableValue(z, c128, zt); err == nil {
			return reflect.ValueOf(imag(zz.Complex())), zt, nil
		} else if zz, err := assignableValue(z, c64, zt); err == nil {
			return reflect.ValueOf(float32(imag(zz.Complex()))), zt, nil
		} else {
			return reflect.Zero(f64), false, ErrBadBuiltinArgument{"imag", z}
		}
	}),
	"cap": reflect.ValueOf(builtinCap),
	"len": reflect.ValueOf(builtinLen),
	"panic": reflect.ValueOf(func(z reflect.Value, zt bool) (reflect.Value, bool, error) {
		panic(z.Interface())
		return reflect.ValueOf(nil), false, errors.New("Panic")
	}),
}

var builtinTypes = map[string] reflect.Type{
	"int": reflect.TypeOf(int(0)),
	"int8": reflect.TypeOf(int8(0)),
	"int16": reflect.TypeOf(int16(0)),
	"int32": reflect.TypeOf(int32(0)),
	"int64": reflect.TypeOf(int64(0)),

	"uint": reflect.TypeOf(uint(0)),
	"uint8": reflect.TypeOf(uint8(0)),
	"uint16": reflect.TypeOf(uint16(0)),
	"uint32": reflect.TypeOf(uint32(0)),
	"uint64": reflect.TypeOf(uint64(0)),

	"float32": reflect.TypeOf(float32(0)),
	"float64": reflect.TypeOf(float64(0)),

	"complex64": reflect.TypeOf(complex64(0)),
	"complex128": reflect.TypeOf(complex128(0)),

	"bool": reflect.TypeOf(bool(false)),
	"byte": reflect.TypeOf(byte(0)),
	"rune": RuneType,
	"string": reflect.TypeOf(""),

	"error": reflect.TypeOf(errors.New("")),
}

func builtinCap(v reflect.Value, vt bool) (reflect.Value, bool, error) {
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(v.Cap()), true, nil
	default:
		return reflect.Zero(intType), false,
		errors.New(fmt.Sprintf("invalid argument %v (type %v) for cap",
			v.Interface(), v.Type()))
	}
}

func builtinLen(z reflect.Value, zt bool) (reflect.Value, bool, error) {
	switch z.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return reflect.ValueOf(z.Len()), true, nil
	default:
		return reflect.Zero(intType), false, ErrBadBuiltinArgument{"len", z}
	}
}
