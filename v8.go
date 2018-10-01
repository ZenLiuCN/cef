package cef

/*
#include <stdlib.h>
#include "string.h"
#include "include/capi/cef_app_capi.h"
#include "include/capi/cef_client_capi.h"
#include "cef_helpers.h"
*/
import "C"

import (
	"fmt"
	"reflect"
	//"strconv"
	"unsafe"
)

type V8Value *C.cef_v8value_t
type V8Callback func([]V8Value)
type V8Handler func(*Browser, []V8Value) interface{}

var V8Callbacks map[string]V8Callback
var V8Handlers map[string]V8Handler
var V8Extensions map[string]string

func init() {
	if V8Handlers == nil {
		V8Handlers = make(map[string]V8Handler)
	}
	if V8Extensions == nil {
		V8Extensions = make(map[string]string)
	}
}

//export go_RenderProcessHandlerOnWebKitInitialized
func go_RenderProcessHandlerOnWebKitInitialized(handler *C.cef_v8handler_t) {
	logger.Println("go_RenderProcessHandlerOnWebKitInitialized")

	init_cef_handlers()
	for ext, code := range V8Extensions {
		C.cef_register_extension(CEFString(ext), CEFString(code), handler)
	}

}

//2018-10-01 add to register mulit extentisons
func RegisterJsExtensionHandler(ext, code string, handlers map[string]V8Handler) {
	if V8Extensions == nil {
		V8Extensions = make(map[string]string)
	}
	if _, ok := V8Extensions[ext]; ok {
		V8Extensions[ext] += code
	} else {
		V8Extensions[ext] = code
	}
	for name, handler := range handlers {
		RegisterV8Handler(name, handler)
	}
}

//2018-10-01 add to register mulit extentisons
func RegisterJsExtensionCallback(ext, code string, handlers map[string]V8Callback) {
	if V8Extensions == nil {
		V8Extensions = make(map[string]string)
	}
	if _, ok := V8Extensions[ext]; ok {
		V8Extensions[ext] += code
	} else {
		V8Extensions[ext] = code
	}
	for name, handler := range handlers {
		RegisterV8Callback(name, handler)
	}
}

//export go_V8HandlerExecute
func go_V8HandlerExecute(browserId int, browser *C.cef_browser_t, name *C.cef_string_t, object *C.cef_v8value_t, argsCount C.size_t, args **C.cef_v8value_t, retval **C.cef_v8value_t, exception *C.cef_string_t) int {
	callbackName := CEFToGoString(name)
	argsN := int(argsCount)
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(args)),
		Len:  argsN,
		Cap:  argsN,
	}
	arguments := *(*[]V8Value)(unsafe.Pointer(&hdr))
	gbrowser := &Browser{browserId, browser, nil}
	if cb, ok := V8Handlers[callbackName]; ok {
		r := cb(gbrowser, arguments)
		switch v := r.(type) {
		case string:
			//fmt.Println("is string", v)
			*retval = C.cef_v8value_create_string(CEFString(v))
		case float64:
			//fmt.Println("is float", int64(v))
			*retval = C.cef_v8value_create_double(C.double(v))
		case int:
			//fmt.Println("is int", v)
			*retval = C.cef_v8value_create_int(C.int32(v))
		case []interface{}:
			//fmt.Println("is an array:")
			for i, u := range v {
				fmt.Println(i, u)
			}
		case nil:
			//fmt.Println("is nil", "null")
			*retval = C.cef_v8value_create_null()
		case map[string]interface{}:
			//fmt.Println("is an map:")
			//print_json(vv)
		default:
			fmt.Println("is of a type I don't know how to handle ")
		}
		return 1
	} else {
		return 0
	}
}

//export go_V8HandlerExecute2
func go_V8HandlerExecute2(browserId int, browser *C.cef_browser_t, name *C.cef_string_t, object *C.cef_v8value_t, argsCount C.size_t, args **C.cef_v8value_t, retval **C.cef_v8value_t, exception *C.cef_string_t) int {
	argsN := int(argsCount)
	if argsN < 1 {
		return 0
	}
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(args)),
		Len:  argsN,
		Cap:  argsN,
	}
	arguments := *(*[]V8Value)(unsafe.Pointer(&hdr))
	handler1, arguments := arguments[0], arguments[1:]
	handler := unsafe.Pointer(uintptr(V8ValueToInt32(handler1)))
	logger.Printf("BrowserId: %v Handler: %v Args: %v\n", browserId, handler, arguments)
	//fmt.Printf("Browser: %v\n", browser)
	gName := CEFToGoString(name)
	gbrowser := &Browser{browserId, browser, nil}

	if handler, ok := V8Handlers[gName]; ok {
		r := handler(gbrowser, arguments)
		switch v := r.(type) {
		case string:
			//fmt.Println("is string", v)
			*retval = C.cef_v8value_create_string(CEFString(v))
		case float64:
			//fmt.Println("is float", int64(v))
			*retval = C.cef_v8value_create_double(C.double(v))
		case int:
			//fmt.Println("is int", v)
			*retval = C.cef_v8value_create_int(C.int32(v))
		case []interface{}:
			//fmt.Println("is an array:")
			for i, u := range v {
				fmt.Println(i, u)
			}
		case nil:
			//fmt.Println("is nil", "null")
			*retval = C.cef_v8value_create_null()
		case map[string]interface{}:
			//fmt.Println("is an map:")
			//print_json(vv)
		default:
			fmt.Println("is of a type I don't know how to handle ")
		}
		return 1
	}

	return 0
}

func RegisterV8Handler(name string, handler V8Handler) {
	if V8Handlers == nil {
		V8Handlers = make(map[string]V8Handler)
	}
	V8Handlers[name] = handler
}

func RegisterV8Callback(name string, callback V8Callback) {
	if V8Callbacks == nil {
		V8Callbacks = make(map[string]V8Callback)
	}
	V8Callbacks[name] = callback
}

func V8ValueToInt32(v V8Value) int32 {
	return int32(C.v8ValueToInt32((*C.cef_v8value_t)(v)))
}

func V8ValueToString(v V8Value) string {
	return CEFToGoString(C.v8ValueToString((*C.cef_v8value_t)(v)))
}

func V8ValueToBool(v V8Value) bool {
	return int(C.v8ValueToBool((*C.cef_v8value_t)(v))) == 1
}
