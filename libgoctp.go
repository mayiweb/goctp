// +build linux,cgo windows,cgo

// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goctp

/*
#cgo linux LDFLAGS: -fPIC -L${SRCDIR}/api/v6.3.15_20190220_api_tradeapi_se_linux64 -Wl,-rpath,${SRCDIR}/api/v6.3.15_20190220_api_tradeapi_se_linux64 -lthostmduserapi_se -lthosttraderapi_se -lstdc++
#cgo linux CPPFLAGS: -fPIC -I${SRCDIR}/api/v6.3.15_20190220_api_tradeapi_se_linux64
#cgo windows LDFLAGS: -fPIC -L${SRCDIR}/api/6.3.15_20190220_tradeapi64_se_windows -Wl,-rpath,${SRCDIR}/api/6.3.15_20190220_tradeapi64_se_windows ${SRCDIR}/api/6.3.15_20190220_tradeapi64_se_windows/thostmduserapi.lib ${SRCDIR}/api/6.3.15_20190220_tradeapi64_se_windows/thosttraderapi.lib -lthostmduserapi -lthosttraderapi
#cgo windows CPPFLAGS: -fPIC -I${SRCDIR}/api/6.3.15_20190220_tradeapi64_se_windows -DISLIB -DWIN32 -DLIB_MD_API_EXPORT -DLIB_TRADER_API_EXPORT
*/
import "C"