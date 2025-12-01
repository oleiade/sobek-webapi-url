package url

import (
	"errors"
	"fmt"

	"github.com/grafana/sobek"
)

const iterableExtractorSource = `(function(obj) {
	const result = [];
	for (const item of obj) {
		if (!Array.isArray(item) || item.length !== 2) {
			throw new TypeError("Invalid argument");
		}
		result.push([String(item[0]), String(item[1])]);
	}
	return result;
})`

// RegisterRuntime exports the URL and URLSearchParams constructors
// into the provided sobek runtime.
func RegisterRuntime(rt *sobek.Runtime) error {
	if err := bindURL(rt); err != nil {
		return err
	}

	return bindURLSearchParams(rt)
}

// bindURL registers the URL constructor and static methods.
//
//nolint:funlen // This function is intentionally long as it defines all URL constructor logic in one place.
func bindURL(rt *sobek.Runtime) error {
	constructor := func(call sobek.ConstructorCall) *sobek.Object {
		// Get the input argument (required)
		inputArg := call.Argument(0)
		if isNullish(inputArg) {
			throwAsJSError(rt, invalidURLError())
		}

		input := inputArg.String()

		// Get the optional base argument
		var base string
		baseArg := call.Argument(1)
		if !isNullish(baseArg) {
			// base can be a string or a URL object
			if baseObj, ok := baseArg.Export().(*URL); ok {
				base = baseObj.Href()
			} else {
				base = baseArg.String()
			}
		}

		u, err := NewURL(input, base)
		if err != nil {
			throwAsJSError(rt, err)
		}

		return newURLObject(rt, u, call.This)
	}

	// Set the constructor
	if err := rt.Set("URL", constructor); err != nil {
		return fmt.Errorf("setting URL constructor: %w", err)
	}

	// Get the URL constructor object to add static methods
	urlConstructor := rt.Get("URL").ToObject(rt)

	// Add URL.canParse static method
	canParseFunc := func(call sobek.FunctionCall) sobek.Value {
		inputArg := call.Argument(0)
		// Convert undefined/null to "undefined"/"null" string as per JS behavior
		var input string
		if isNullish(inputArg) {
			input = "undefined"
		} else {
			input = inputArg.String()
		}

		var base string
		baseArg := call.Argument(1)
		if !isNullish(baseArg) {
			base = baseArg.String()
		}

		return rt.ToValue(CanParse(input, base))
	}

	if err := urlConstructor.Set("canParse", canParseFunc); err != nil {
		return fmt.Errorf("setting URL.canParse: %w", err)
	}

	// Add URL.parse static method
	parseFunc := func(call sobek.FunctionCall) sobek.Value {
		inputArg := call.Argument(0)
		// Convert undefined/null to "undefined"/"null" string as per JS behavior
		var input string
		if isNullish(inputArg) {
			input = "undefined"
		} else {
			input = inputArg.String()
		}

		var base string
		baseArg := call.Argument(1)
		if !isNullish(baseArg) {
			base = baseArg.String()
		}

		u := Parse(input, base)
		if u == nil {
			return sobek.Null()
		}

		// Create a new URL object
		obj := rt.NewObject()
		return newURLObject(rt, u, obj)
	}

	if err := urlConstructor.Set("parse", parseFunc); err != nil {
		return fmt.Errorf("setting URL.parse: %w", err)
	}

	return nil
}

// newURLObject creates a JS object wrapping a Go URL instance.
//
//nolint:funlen // This function is intentionally long as it defines all URL properties and methods.
func newURLObject(rt *sobek.Runtime, u *URL, obj *sobek.Object) *sobek.Object {
	// Create the searchParams object once and cache it
	searchParamsObj := newURLSearchParamsObject(rt, u.SearchParams())

	defineAccessor(rt, obj, "href",
		func(_ sobek.FunctionCall) sobek.Value {
			return rt.ToValue(u.Href())
		},
		func(call sobek.FunctionCall) sobek.Value {
			if len(call.Arguments) > 0 {
				if err := u.SetHref(call.Argument(0).String()); err != nil {
					throwAsJSError(rt, err)
				}
				// Update searchParams reference
				searchParamsObj = newURLSearchParamsObject(rt, u.SearchParams())
			}
			return sobek.Undefined()
		})

	defineAccessor(rt, obj, "origin",
		func(_ sobek.FunctionCall) sobek.Value {
			return rt.ToValue(u.Origin())
		},
		nil)

	defineAccessor(rt, obj, "protocol",
		func(_ sobek.FunctionCall) sobek.Value {
			return rt.ToValue(u.Protocol())
		},
		func(call sobek.FunctionCall) sobek.Value {
			if len(call.Arguments) > 0 {
				u.SetProtocol(call.Argument(0).String())
			}
			return sobek.Undefined()
		})

	defineAccessor(rt, obj, "username",
		func(_ sobek.FunctionCall) sobek.Value {
			return rt.ToValue(u.Username())
		},
		func(call sobek.FunctionCall) sobek.Value {
			if len(call.Arguments) > 0 {
				u.SetUsername(call.Argument(0).String())
			}
			return sobek.Undefined()
		})

	defineAccessor(rt, obj, "password",
		func(_ sobek.FunctionCall) sobek.Value {
			return rt.ToValue(u.Password())
		},
		func(call sobek.FunctionCall) sobek.Value {
			if len(call.Arguments) > 0 {
				u.SetPassword(call.Argument(0).String())
			}
			return sobek.Undefined()
		})

	defineAccessor(rt, obj, "host",
		func(_ sobek.FunctionCall) sobek.Value {
			return rt.ToValue(u.Host())
		},
		func(call sobek.FunctionCall) sobek.Value {
			if len(call.Arguments) > 0 {
				u.SetHost(call.Argument(0).String())
			}
			return sobek.Undefined()
		})

	defineAccessor(rt, obj, "hostname",
		func(_ sobek.FunctionCall) sobek.Value {
			return rt.ToValue(u.Hostname())
		},
		func(call sobek.FunctionCall) sobek.Value {
			if len(call.Arguments) > 0 {
				u.SetHostname(call.Argument(0).String())
			}
			return sobek.Undefined()
		})

	defineAccessor(rt, obj, "port",
		func(_ sobek.FunctionCall) sobek.Value {
			return rt.ToValue(u.Port())
		},
		func(call sobek.FunctionCall) sobek.Value {
			if len(call.Arguments) > 0 {
				u.SetPort(call.Argument(0).String())
			}
			return sobek.Undefined()
		})

	defineAccessor(rt, obj, "pathname",
		func(_ sobek.FunctionCall) sobek.Value {
			return rt.ToValue(u.Pathname())
		},
		func(call sobek.FunctionCall) sobek.Value {
			if len(call.Arguments) > 0 {
				u.SetPathname(call.Argument(0).String())
			}
			return sobek.Undefined()
		})

	defineAccessor(rt, obj, "search",
		func(_ sobek.FunctionCall) sobek.Value {
			return rt.ToValue(u.Search())
		},
		func(call sobek.FunctionCall) sobek.Value {
			if len(call.Arguments) > 0 {
				u.SetSearch(call.Argument(0).String())
				// Update searchParams reference
				searchParamsObj = newURLSearchParamsObject(rt, u.SearchParams())
			}
			return sobek.Undefined()
		})

	defineAccessor(rt, obj, "searchParams",
		func(_ sobek.FunctionCall) sobek.Value {
			return searchParamsObj
		},
		nil)

	defineAccessor(rt, obj, "hash",
		func(_ sobek.FunctionCall) sobek.Value {
			return rt.ToValue(u.Hash())
		},
		func(call sobek.FunctionCall) sobek.Value {
			if len(call.Arguments) > 0 {
				u.SetHash(call.Argument(0).String())
			}
			return sobek.Undefined()
		})

	// Define toString method
	toStringMethod := func(_ sobek.FunctionCall) sobek.Value {
		return rt.ToValue(u.String())
	}
	if err := obj.Set("toString", toStringMethod); err != nil {
		panic(rt.NewGoError(fmt.Errorf("defining toString method: %w", err)))
	}

	// Define toJSON method
	toJSONMethod := func(_ sobek.FunctionCall) sobek.Value {
		return rt.ToValue(u.ToJSON())
	}
	if err := obj.Set("toJSON", toJSONMethod); err != nil {
		panic(rt.NewGoError(fmt.Errorf("defining toJSON method: %w", err)))
	}

	return obj
}

// bindURLSearchParams registers the URLSearchParams constructor.
//
//nolint:gocognit,nestif // Complex constructor logic to handle multiple input types as per WHATWG spec.
func bindURLSearchParams(rt *sobek.Runtime) error {
	constructor := func(call sobek.ConstructorCall) *sobek.Object {
		var sp *URLSearchParams

		initArg := call.Argument(0)

		if isNullish(initArg) {
			// No argument or undefined/null - create empty params
			sp = NewURLSearchParams()
		} else {
			// First check if it's a string
			exported := initArg.Export()
			if str, ok := exported.(string); ok {
				sp = NewURLSearchParamsFromString(str)
			} else if arr, ok := exported.([]interface{}); ok {
				// Array of pairs
				sp = NewURLSearchParams()
				for _, item := range arr {
					if pair, ok := item.([]interface{}); ok && len(pair) == 2 {
						key := fmt.Sprintf("%v", pair[0])
						value := fmt.Sprintf("%v", pair[1])
						sp.Append(key, value)
					} else if pair, ok := item.([]string); ok && len(pair) == 2 {
						sp.Append(pair[0], pair[1])
					} else {
						throwAsJSError(rt, NewError(TypeError, "Invalid argument"))
					}
				}
			} else {
				// Check if it has Symbol.iterator (like URLSearchParams or arrays)
				obj := initArg.ToObject(rt)
				iteratorMethod := obj.GetSymbol(sobek.SymIterator)

				if iteratorMethod != nil && !isNullish(iteratorMethod) {
					// Has iterator - iterate over it
					sp = NewURLSearchParams()
					iterator, err := rt.RunString(iterableExtractorSource)
					if err != nil {
						throwAsJSError(rt, NewError(TypeError, "Invalid argument"))
					}

					iterFn, ok := sobek.AssertFunction(iterator)
					if !ok {
						throwAsJSError(rt, NewError(TypeError, "Invalid argument"))
					}

					result, err := iterFn(sobek.Undefined(), initArg)
					if err != nil {
						throwAsJSError(rt, NewError(TypeError, "Invalid argument"))
					}

					if resultArr, ok := result.Export().([]interface{}); ok {
						for _, item := range resultArr {
							if pair, ok := item.([]interface{}); ok && len(pair) == 2 {
								sp.Append(fmt.Sprintf("%v", pair[0]), fmt.Sprintf("%v", pair[1]))
							}
						}
					}
				} else {
					// Try as record (object with string keys)
					sp = NewURLSearchParams()
					for _, key := range obj.Keys() {
						val := obj.Get(key)
						if val != nil {
							sp.Append(key, val.String())
						}
					}
				}
			}
		}

		return newURLSearchParamsObject(rt, sp)
	}

	return rt.Set("URLSearchParams", constructor)
}

// newURLSearchParamsObject creates a JS object wrapping a Go URLSearchParams instance.
//
//nolint:gocognit,cyclop,funlen // This function is intentionally complex as it defines all URLSearchParams methods.
func newURLSearchParamsObject(rt *sobek.Runtime, sp *URLSearchParams) *sobek.Object {
	obj := rt.NewObject()

	// Set Symbol.toPrimitive for proper string conversion (params + '')
	toPrimitiveMethod := func(_ sobek.FunctionCall) sobek.Value {
		return rt.ToValue(sp.String())
	}
	if err := obj.SetSymbol(sobek.SymToPrimitive, rt.ToValue(toPrimitiveMethod)); err != nil {
		panic(rt.NewGoError(fmt.Errorf("defining Symbol.toPrimitive: %w", err)))
	}

	// append method
	appendMethod := func(call sobek.FunctionCall) sobek.Value {
		if len(call.Arguments) < 2 {
			return sobek.Undefined()
		}
		key := call.Argument(0).String()
		value := call.Argument(1).String()
		sp.Append(key, value)
		return sobek.Undefined()
	}
	if err := obj.Set("append", appendMethod); err != nil {
		panic(rt.NewGoError(err))
	}

	// delete method
	deleteMethod := func(call sobek.FunctionCall) sobek.Value {
		if len(call.Arguments) < 1 {
			return sobek.Undefined()
		}
		key := call.Argument(0).String()
		if len(call.Arguments) > 1 && !isNullish(call.Argument(1)) {
			sp.DeletePair(key, call.Argument(1).String())
		} else {
			sp.DeleteAll(key)
		}
		return sobek.Undefined()
	}
	if err := obj.Set("delete", deleteMethod); err != nil {
		panic(rt.NewGoError(err))
	}

	// get method
	getMethod := func(call sobek.FunctionCall) sobek.Value {
		if len(call.Arguments) < 1 {
			return sobek.Null()
		}
		key := call.Argument(0).String()
		value, found := sp.Get(key)
		if !found {
			return sobek.Null()
		}
		return rt.ToValue(value)
	}
	if err := obj.Set("get", getMethod); err != nil {
		panic(rt.NewGoError(err))
	}

	// getAll method
	getAllMethod := func(call sobek.FunctionCall) sobek.Value {
		if len(call.Arguments) < 1 {
			return rt.NewArray()
		}
		key := call.Argument(0).String()
		values := sp.GetAll(key)
		return rt.ToValue(values)
	}
	if err := obj.Set("getAll", getAllMethod); err != nil {
		panic(rt.NewGoError(err))
	}

	// has method
	hasMethod := func(call sobek.FunctionCall) sobek.Value {
		if len(call.Arguments) < 1 {
			return rt.ToValue(false)
		}
		key := call.Argument(0).String()
		if len(call.Arguments) > 1 && !isNullish(call.Argument(1)) {
			return rt.ToValue(sp.HasPair(key, call.Argument(1).String()))
		}
		return rt.ToValue(sp.HasKey(key))
	}
	if err := obj.Set("has", hasMethod); err != nil {
		panic(rt.NewGoError(err))
	}

	// set method
	setMethod := func(call sobek.FunctionCall) sobek.Value {
		if len(call.Arguments) < 2 {
			return sobek.Undefined()
		}
		key := call.Argument(0).String()
		value := call.Argument(1).String()
		sp.Set(key, value)
		return sobek.Undefined()
	}
	if err := obj.Set("set", setMethod); err != nil {
		panic(rt.NewGoError(err))
	}

	// sort method
	sortMethod := func(_ sobek.FunctionCall) sobek.Value {
		sp.Sort()
		return sobek.Undefined()
	}
	if err := obj.Set("sort", sortMethod); err != nil {
		panic(rt.NewGoError(err))
	}

	// toString method
	toStringMethod := func(_ sobek.FunctionCall) sobek.Value {
		return rt.ToValue(sp.String())
	}
	if err := obj.Set("toString", toStringMethod); err != nil {
		panic(rt.NewGoError(err))
	}

	// forEach method
	forEachMethod := func(call sobek.FunctionCall) sobek.Value {
		if len(call.Arguments) < 1 {
			return sobek.Undefined()
		}

		callback, ok := sobek.AssertFunction(call.Argument(0))
		if !ok {
			throwAsJSError(rt, NewError(TypeError, "Callback is not a function"))
		}

		thisArg := sobek.Undefined()
		if len(call.Arguments) > 1 {
			thisArg = call.Argument(1)
		}

		sp.ForEach(func(value, key string) {
			_, err := callback(thisArg, rt.ToValue(value), rt.ToValue(key), obj)
			if err != nil {
				panic(err)
			}
		})

		return sobek.Undefined()
	}
	if err := obj.Set("forEach", forEachMethod); err != nil {
		panic(rt.NewGoError(err))
	}

	// entries method - returns an iterator
	entriesMethod := func(_ sobek.FunctionCall) sobek.Value {
		return sliceIterator(rt, entriesToInterfaces(sp.Entries()))
	}
	if err := obj.Set("entries", entriesMethod); err != nil {
		panic(rt.NewGoError(err))
	}

	// keys method - returns an iterator
	keysMethod := func(_ sobek.FunctionCall) sobek.Value {
		return sliceIterator(rt, sp.Keys())
	}
	if err := obj.Set("keys", keysMethod); err != nil {
		panic(rt.NewGoError(err))
	}

	// values method - returns an iterator
	valuesMethod := func(_ sobek.FunctionCall) sobek.Value {
		return sliceIterator(rt, sp.Values())
	}
	if err := obj.Set("values", valuesMethod); err != nil {
		panic(rt.NewGoError(err))
	}

	// size property (getter)
	defineAccessor(rt, obj, "size",
		func(_ sobek.FunctionCall) sobek.Value {
			return rt.ToValue(sp.Size())
		},
		nil)

	// Symbol.iterator - make URLSearchParams iterable
	// Returns the same as entries()
	iteratorMethod := func(_ sobek.FunctionCall) sobek.Value {
		return sliceIterator(rt, entriesToInterfaces(sp.Entries()))
	}
	if err := obj.SetSymbol(sobek.SymIterator, rt.ToValue(iteratorMethod)); err != nil {
		panic(rt.NewGoError(fmt.Errorf("defining Symbol.iterator: %w", err)))
	}

	return obj
}

// throwAsJSError converts an error to a JS exception and panics.
func throwAsJSError(rt *sobek.Runtime, err error) {
	var urlErr *Error
	if errors.As(err, &urlErr) {
		panic(urlErr.JSError(rt))
	}
	panic(rt.NewGoError(err))
}

// isNullish returns true if the value is null or undefined.
func isNullish(v sobek.Value) bool {
	return v == nil || sobek.IsUndefined(v) || sobek.IsNull(v)
}

func defineAccessor(rt *sobek.Runtime, obj *sobek.Object, name string,
	getter func(call sobek.FunctionCall) sobek.Value,
	setter func(call sobek.FunctionCall) sobek.Value,
) {
	var getterValue sobek.Value
	var setterValue sobek.Value
	if getter != nil {
		getterValue = rt.ToValue(getter)
	}
	if setter != nil {
		setterValue = rt.ToValue(setter)
	}
	if err := obj.DefineAccessorProperty(name, getterValue, setterValue, sobek.FLAG_FALSE, sobek.FLAG_TRUE); err != nil {
		panic(rt.NewGoError(fmt.Errorf("defining %s property: %w", name, err)))
	}
}

func sliceIterator(rt *sobek.Runtime, data interface{}) sobek.Value {
	arr := rt.ToValue(data).ToObject(rt)
	iteratorFn := arr.GetSymbol(sobek.SymIterator)
	if fn, ok := sobek.AssertFunction(iteratorFn); ok {
		iter, _ := fn(arr)
		return iter
	}
	return arr
}

func entriesToInterfaces(entries [][2]string) []interface{} {
	result := make([]interface{}, len(entries))
	for i, entry := range entries {
		result[i] = []interface{}{entry[0], entry[1]}
	}
	return result
}
