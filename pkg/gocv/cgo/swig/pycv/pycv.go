/* ----------------------------------------------------------------------------
 * This file was automatically generated by SWIG (https://www.swig.org).
 * Version 4.1.1
 *
 * Do not make changes to this file unless you know what you are doing - modify
 * the SWIG interface file instead.
 * ----------------------------------------------------------------------------- */

// source: pycv.swigcxx

package pycv

/*
#define intgo swig_intgo
typedef void *swig_voidp;

#include <stddef.h>
#include <stdint.h>


typedef long long intgo;
typedef unsigned long long uintgo;



typedef struct { char *p; intgo n; } _gostring_;
typedef struct { void* array; intgo len; intgo cap; } _goslice_;


typedef _gostring_ swig_type_1;
typedef _gostring_ swig_type_2;
typedef _gostring_ swig_type_3;
typedef _gostring_ swig_type_4;
typedef _gostring_ swig_type_5;
typedef _gostring_ swig_type_6;
typedef _gostring_ swig_type_7;
typedef _gostring_ swig_type_8;
typedef _gostring_ swig_type_9;
extern void _wrap_Swig_free_pycv_ffa8ea6aa3be9035(uintptr_t arg1);
extern uintptr_t _wrap_Swig_malloc_pycv_ffa8ea6aa3be9035(swig_intgo arg1);
extern void _wrap_LocalInitRequest_gpu_id_set_pycv_ffa8ea6aa3be9035(uintptr_t arg1, swig_intgo arg2);
extern swig_intgo _wrap_LocalInitRequest_gpu_id_get_pycv_ffa8ea6aa3be9035(uintptr_t arg1);
extern void _wrap_LocalInitRequest_sdk_dir_set_pycv_ffa8ea6aa3be9035(uintptr_t arg1, swig_type_1 arg2);
extern swig_type_2 _wrap_LocalInitRequest_sdk_dir_get_pycv_ffa8ea6aa3be9035(uintptr_t arg1);
extern void _wrap_LocalInitRequest_model_dir_set_pycv_ffa8ea6aa3be9035(uintptr_t arg1, swig_type_3 arg2);
extern swig_type_4 _wrap_LocalInitRequest_model_dir_get_pycv_ffa8ea6aa3be9035(uintptr_t arg1);
extern uintptr_t _wrap_new_LocalInitRequest_pycv_ffa8ea6aa3be9035(void);
extern void _wrap_delete_LocalInitRequest_pycv_ffa8ea6aa3be9035(uintptr_t arg1);
extern uintptr_t _wrap_new_LocalInitResponse_pycv_ffa8ea6aa3be9035(void);
extern void _wrap_delete_LocalInitResponse_pycv_ffa8ea6aa3be9035(uintptr_t arg1);
extern void _wrap_DoRequest_arg1_set_pycv_ffa8ea6aa3be9035(uintptr_t arg1, swig_type_5 arg2);
extern swig_type_6 _wrap_DoRequest_arg1_get_pycv_ffa8ea6aa3be9035(uintptr_t arg1);
extern void _wrap_DoRequest_arg2_set_pycv_ffa8ea6aa3be9035(uintptr_t arg1, swig_type_7 arg2);
extern swig_type_8 _wrap_DoRequest_arg2_get_pycv_ffa8ea6aa3be9035(uintptr_t arg1);
extern uintptr_t _wrap_new_DoRequest_pycv_ffa8ea6aa3be9035(void);
extern void _wrap_delete_DoRequest_pycv_ffa8ea6aa3be9035(uintptr_t arg1);
extern uintptr_t _wrap_new_DoResponse_pycv_ffa8ea6aa3be9035(void);
extern void _wrap_delete_DoResponse_pycv_ffa8ea6aa3be9035(uintptr_t arg1);
extern uintptr_t _wrap_new_Wrapped_PyImage_pycv_ffa8ea6aa3be9035(void);
extern void _wrap_delete_Wrapped_PyImage_pycv_ffa8ea6aa3be9035(uintptr_t arg1);
extern void _wrap_Wrapped_PyImage_GlobalInit_pycv_ffa8ea6aa3be9035(swig_type_9 arg1, swig_intgo arg2);
extern void _wrap_Wrapped_PyImage_GlobalRelease_pycv_ffa8ea6aa3be9035(void);
extern void _wrap_Wrapped_PyImage_Wrapped_PyImage_LocalInit_pycv_ffa8ea6aa3be9035(uintptr_t arg1, uintptr_t arg2, uintptr_t arg3);
extern void _wrap_Wrapped_PyImage_Wrapped_PyImage_Do_pycv_ffa8ea6aa3be9035(uintptr_t arg1, uintptr_t arg2, uintptr_t arg3);
#undef intgo
*/
import "C"

import "unsafe"
import _ "runtime/cgo"
import "sync"
import "fmt"


type _ unsafe.Pointer



var Swig_escape_always_false bool
var Swig_escape_val interface{}


type _swig_fnptr *byte
type _swig_memberptr *byte


func getSwigcptr(v interface { Swigcptr() uintptr }) uintptr {
	if v == nil {
		return 0
	}
	return v.Swigcptr()
}


type _ sync.Mutex

//export cgo_panic__pycv_ffa8ea6aa3be9035
func cgo_panic__pycv_ffa8ea6aa3be9035(p *byte) {
	s := (*[1024]byte)(unsafe.Pointer(p))[:]
	for i, b := range s {
		if b == 0 {
			panic(string(s[:i]))
		}
	}
	panic(string(s))
}


type swig_gostring struct { p uintptr; n int }
func swigCopyString(s string) string {
  p := *(*swig_gostring)(unsafe.Pointer(&s))
  r := string((*[0x7fffffff]byte)(unsafe.Pointer(p.p))[:p.n])
  Swig_free(p.p)
  return r
}

func Swig_free(arg1 uintptr) {
	_swig_i_0 := arg1
	C._wrap_Swig_free_pycv_ffa8ea6aa3be9035(C.uintptr_t(_swig_i_0))
}

func Swig_malloc(arg1 int) (_swig_ret uintptr) {
	var swig_r uintptr
	_swig_i_0 := arg1
	swig_r = (uintptr)(C._wrap_Swig_malloc_pycv_ffa8ea6aa3be9035(C.swig_intgo(_swig_i_0)))
	return swig_r
}

type SwigcptrLocalInitRequest uintptr

func (p SwigcptrLocalInitRequest) Swigcptr() uintptr {
	return (uintptr)(p)
}

func (p SwigcptrLocalInitRequest) SwigIsLocalInitRequest() {
}

func (arg1 SwigcptrLocalInitRequest) SetGpu_id(arg2 int) {
	_swig_i_0 := arg1
	_swig_i_1 := arg2
	C._wrap_LocalInitRequest_gpu_id_set_pycv_ffa8ea6aa3be9035(C.uintptr_t(_swig_i_0), C.swig_intgo(_swig_i_1))
}

func (arg1 SwigcptrLocalInitRequest) GetGpu_id() (_swig_ret int) {
	var swig_r int
	_swig_i_0 := arg1
	swig_r = (int)(C._wrap_LocalInitRequest_gpu_id_get_pycv_ffa8ea6aa3be9035(C.uintptr_t(_swig_i_0)))
	return swig_r
}

func (arg1 SwigcptrLocalInitRequest) SetSdk_dir(arg2 string) {
	_swig_i_0 := arg1
	_swig_i_1 := arg2
	C._wrap_LocalInitRequest_sdk_dir_set_pycv_ffa8ea6aa3be9035(C.uintptr_t(_swig_i_0), *(*C.swig_type_1)(unsafe.Pointer(&_swig_i_1)))
	if Swig_escape_always_false {
		Swig_escape_val = arg2
	}
}

func (arg1 SwigcptrLocalInitRequest) GetSdk_dir() (_swig_ret string) {
	var swig_r string
	_swig_i_0 := arg1
	swig_r_p := C._wrap_LocalInitRequest_sdk_dir_get_pycv_ffa8ea6aa3be9035(C.uintptr_t(_swig_i_0))
	swig_r = *(*string)(unsafe.Pointer(&swig_r_p))
	var swig_r_1 string
 swig_r_1 = swigCopyString(swig_r) 
	return swig_r_1
}

func (arg1 SwigcptrLocalInitRequest) SetModel_dir(arg2 string) {
	_swig_i_0 := arg1
	_swig_i_1 := arg2
	C._wrap_LocalInitRequest_model_dir_set_pycv_ffa8ea6aa3be9035(C.uintptr_t(_swig_i_0), *(*C.swig_type_3)(unsafe.Pointer(&_swig_i_1)))
	if Swig_escape_always_false {
		Swig_escape_val = arg2
	}
}

func (arg1 SwigcptrLocalInitRequest) GetModel_dir() (_swig_ret string) {
	var swig_r string
	_swig_i_0 := arg1
	swig_r_p := C._wrap_LocalInitRequest_model_dir_get_pycv_ffa8ea6aa3be9035(C.uintptr_t(_swig_i_0))
	swig_r = *(*string)(unsafe.Pointer(&swig_r_p))
	var swig_r_1 string
 swig_r_1 = swigCopyString(swig_r) 
	return swig_r_1
}

func NewLocalInitRequest() (_swig_ret LocalInitRequest) {
	var swig_r LocalInitRequest
	swig_r = (LocalInitRequest)(SwigcptrLocalInitRequest(C._wrap_new_LocalInitRequest_pycv_ffa8ea6aa3be9035()))
	return swig_r
}

func DeleteLocalInitRequest(arg1 LocalInitRequest) {
	_swig_i_0 := getSwigcptr(arg1)
	C._wrap_delete_LocalInitRequest_pycv_ffa8ea6aa3be9035(C.uintptr_t(_swig_i_0))
}

type LocalInitRequest interface {
	Swigcptr() uintptr
	SwigIsLocalInitRequest()
	SetGpu_id(arg2 int)
	GetGpu_id() (_swig_ret int)
	SetSdk_dir(arg2 string)
	GetSdk_dir() (_swig_ret string)
	SetModel_dir(arg2 string)
	GetModel_dir() (_swig_ret string)
}

type SwigcptrLocalInitResponse uintptr

func (p SwigcptrLocalInitResponse) Swigcptr() uintptr {
	return (uintptr)(p)
}

func (p SwigcptrLocalInitResponse) SwigIsLocalInitResponse() {
}

func NewLocalInitResponse() (_swig_ret LocalInitResponse) {
	var swig_r LocalInitResponse
	swig_r = (LocalInitResponse)(SwigcptrLocalInitResponse(C._wrap_new_LocalInitResponse_pycv_ffa8ea6aa3be9035()))
	return swig_r
}

func DeleteLocalInitResponse(arg1 LocalInitResponse) {
	_swig_i_0 := getSwigcptr(arg1)
	C._wrap_delete_LocalInitResponse_pycv_ffa8ea6aa3be9035(C.uintptr_t(_swig_i_0))
}

type LocalInitResponse interface {
	Swigcptr() uintptr
	SwigIsLocalInitResponse()
}

type SwigcptrDoRequest uintptr

func (p SwigcptrDoRequest) Swigcptr() uintptr {
	return (uintptr)(p)
}

func (p SwigcptrDoRequest) SwigIsDoRequest() {
}

func (arg1 SwigcptrDoRequest) SetArg1(arg2 string) {
	_swig_i_0 := arg1
	_swig_i_1 := arg2
	C._wrap_DoRequest_arg1_set_pycv_ffa8ea6aa3be9035(C.uintptr_t(_swig_i_0), *(*C.swig_type_5)(unsafe.Pointer(&_swig_i_1)))
	if Swig_escape_always_false {
		Swig_escape_val = arg2
	}
}

func (arg1 SwigcptrDoRequest) GetArg1() (_swig_ret string) {
	var swig_r string
	_swig_i_0 := arg1
	swig_r_p := C._wrap_DoRequest_arg1_get_pycv_ffa8ea6aa3be9035(C.uintptr_t(_swig_i_0))
	swig_r = *(*string)(unsafe.Pointer(&swig_r_p))
	var swig_r_1 string
 swig_r_1 = swigCopyString(swig_r) 
	return swig_r_1
}

func (arg1 SwigcptrDoRequest) SetArg2(arg2 string) {
	_swig_i_0 := arg1
	_swig_i_1 := arg2
	C._wrap_DoRequest_arg2_set_pycv_ffa8ea6aa3be9035(C.uintptr_t(_swig_i_0), *(*C.swig_type_7)(unsafe.Pointer(&_swig_i_1)))
	if Swig_escape_always_false {
		Swig_escape_val = arg2
	}
}

func (arg1 SwigcptrDoRequest) GetArg2() (_swig_ret string) {
	var swig_r string
	_swig_i_0 := arg1
	swig_r_p := C._wrap_DoRequest_arg2_get_pycv_ffa8ea6aa3be9035(C.uintptr_t(_swig_i_0))
	swig_r = *(*string)(unsafe.Pointer(&swig_r_p))
	var swig_r_1 string
 swig_r_1 = swigCopyString(swig_r) 
	return swig_r_1
}

func NewDoRequest() (_swig_ret DoRequest) {
	var swig_r DoRequest
	swig_r = (DoRequest)(SwigcptrDoRequest(C._wrap_new_DoRequest_pycv_ffa8ea6aa3be9035()))
	return swig_r
}

func DeleteDoRequest(arg1 DoRequest) {
	_swig_i_0 := getSwigcptr(arg1)
	C._wrap_delete_DoRequest_pycv_ffa8ea6aa3be9035(C.uintptr_t(_swig_i_0))
}

type DoRequest interface {
	Swigcptr() uintptr
	SwigIsDoRequest()
	SetArg1(arg2 string)
	GetArg1() (_swig_ret string)
	SetArg2(arg2 string)
	GetArg2() (_swig_ret string)
}

type SwigcptrDoResponse uintptr

func (p SwigcptrDoResponse) Swigcptr() uintptr {
	return (uintptr)(p)
}

func (p SwigcptrDoResponse) SwigIsDoResponse() {
}

func NewDoResponse() (_swig_ret DoResponse) {
	var swig_r DoResponse
	swig_r = (DoResponse)(SwigcptrDoResponse(C._wrap_new_DoResponse_pycv_ffa8ea6aa3be9035()))
	return swig_r
}

func DeleteDoResponse(arg1 DoResponse) {
	_swig_i_0 := getSwigcptr(arg1)
	C._wrap_delete_DoResponse_pycv_ffa8ea6aa3be9035(C.uintptr_t(_swig_i_0))
}

type DoResponse interface {
	Swigcptr() uintptr
	SwigIsDoResponse()
}

type SwigcptrWrapped_PyImage uintptr

func (p SwigcptrWrapped_PyImage) Swigcptr() uintptr {
	return (uintptr)(p)
}

func (p SwigcptrWrapped_PyImage) SwigIsWrapped_PyImage() {
}

func NewWrapped_PyImage() (_swig_ret Wrapped_PyImage) {
	var swig_r Wrapped_PyImage
	swig_r = (Wrapped_PyImage)(SwigcptrWrapped_PyImage(C._wrap_new_Wrapped_PyImage_pycv_ffa8ea6aa3be9035()))
	return swig_r
}

func DeleteWrapped_PyImage(arg1 Wrapped_PyImage) {
	_swig_i_0 := getSwigcptr(arg1)
	C._wrap_delete_Wrapped_PyImage_pycv_ffa8ea6aa3be9035(C.uintptr_t(_swig_i_0))
}

func Wrapped_PyImageGlobalInit(arg1 string, arg2 int) {
	_swig_i_0 := arg1
	_swig_i_1 := arg2
	C._wrap_Wrapped_PyImage_GlobalInit_pycv_ffa8ea6aa3be9035(*(*C.swig_type_9)(unsafe.Pointer(&_swig_i_0)), C.swig_intgo(_swig_i_1))
	if Swig_escape_always_false {
		Swig_escape_val = arg1
	}
}

func Wrapped_PyImageGlobalRelease() {
	C._wrap_Wrapped_PyImage_GlobalRelease_pycv_ffa8ea6aa3be9035()
}

func (arg1 SwigcptrWrapped_PyImage) Wrapped_PyImage_LocalInit(arg2 LocalInitRequest, arg3 LocalInitResponse) {
	_swig_i_0 := arg1
	_swig_i_1 := getSwigcptr(arg2)
	_swig_i_2 := getSwigcptr(arg3)
	C._wrap_Wrapped_PyImage_Wrapped_PyImage_LocalInit_pycv_ffa8ea6aa3be9035(C.uintptr_t(_swig_i_0), C.uintptr_t(_swig_i_1), C.uintptr_t(_swig_i_2))
}

func (arg1 SwigcptrWrapped_PyImage) Wrapped_PyImage_Do(arg2 DoRequest, arg3 DoResponse) {
	_swig_i_0 := arg1
	_swig_i_1 := getSwigcptr(arg2)
	_swig_i_2 := getSwigcptr(arg3)
	C._wrap_Wrapped_PyImage_Wrapped_PyImage_Do_pycv_ffa8ea6aa3be9035(C.uintptr_t(_swig_i_0), C.uintptr_t(_swig_i_1), C.uintptr_t(_swig_i_2))
}

type Wrapped_PyImage interface {
	Swigcptr() uintptr
	SwigIsWrapped_PyImage()
	Wrapped_PyImage_LocalInit(arg2 LocalInitRequest, arg3 LocalInitResponse)
	Wrapped_PyImage_Do(arg2 DoRequest, arg3 DoResponse)
}



type PyImage interface {
   Wrapped_PyImage  
   LocalInit(LocalInitRequest) (LocalInitResponse, error);
   Do(DoRequest) (DoResponse, error);
}

func NewPyImage() PyImage {
     e := NewWrapped_PyImage()
     return e.(PyImage)
}

// catch will recover from a panic and store the recover message to the error 
// parameter. The error must be passed by reference in order to be returned to the
// calling function.
func catch(err *error) {
    if r := recover(); r != nil {
        *err = fmt.Errorf("%v", r)
    }
}

func GlobalInit(model_dir string, gpu_id int) (err error) {
   defer catch(&err)
   Wrapped_PyImageGlobalInit(model_dir, gpu_id)
   return
}

func GlobalRelease() (err error){
   defer catch(&err)
   Wrapped_PyImageGlobalRelease()
   return
}

func (arg SwigcptrWrapped_PyImage) LocalInit(req LocalInitRequest) (resp LocalInitResponse, err error) {
   defer catch(&err)
   resp = NewLocalInitResponse()
   arg.Wrapped_PyImage_LocalInit(req, resp)
   return
}

func (arg SwigcptrWrapped_PyImage) Do(req DoRequest) (resp DoResponse, err error) {
   defer catch(&err)
   resp = NewDoResponse()
   arg.Wrapped_PyImage_Do(req, resp)
   return
}





var swigDirectorTrack struct {
	sync.Mutex
	m map[int]interface{}
	c int
}

func swigDirectorAdd(v interface{}) int {
	swigDirectorTrack.Lock()
	defer swigDirectorTrack.Unlock()
	if swigDirectorTrack.m == nil {
		swigDirectorTrack.m = make(map[int]interface{})
	}
	swigDirectorTrack.c++
	ret := swigDirectorTrack.c
	swigDirectorTrack.m[ret] = v
	return ret
}

func swigDirectorLookup(c int) interface{} {
	swigDirectorTrack.Lock()
	defer swigDirectorTrack.Unlock()
	ret := swigDirectorTrack.m[c]
	if ret == nil {
		panic("C++ director pointer not found (possible	use-after-free)")
	}
	return ret
}

func swigDirectorDelete(c int) {
	swigDirectorTrack.Lock()
	defer swigDirectorTrack.Unlock()
	if swigDirectorTrack.m[c] == nil {
		if c > swigDirectorTrack.c {
			panic("C++ director pointer invalid (possible memory corruption")
		} else {
			panic("C++ director pointer not found (possible use-after-free)")
		}
	}
	delete(swigDirectorTrack.m, c)
}


