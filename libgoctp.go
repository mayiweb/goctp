
package goctp

/*
#cgo linux LDFLAGS: -fPIC -L${SRCDIR}/api/v6.5.1_20200908_api_tradeapi_se_linux64 -Wl,-rpath,${SRCDIR}/api/v6.5.1_20200908_api_tradeapi_se_linux64 -lthostmduserapi_se -lthosttraderapi_se -lstdc++
#cgo linux CPPFLAGS: -fPIC -I${SRCDIR}/api/v6.5.1_20200908_api_tradeapi_se_linux64

// windows 不可用，留着记录
#cgo windows LDFLAGS: -fPIC -L${SRCDIR}/api/v6.5.1_20200908_tradeapi64_se_windows -Wl,-rpath,${SRCDIR}/api/v6.5.1_20200908_tradeapi64_se_windows ${SRCDIR}/api/v6.5.1_20200908_tradeapi64_se_windows/thostmduserapi_se.lib ${SRCDIR}/api/v6.5.1_20200908_tradeapi64_se_windows/thosttraderapi_se.lib -lthostmduserapi_se -lthosttraderapi_se
#cgo windows CPPFLAGS: -fPIC -I${SRCDIR}/api/v6.5.1_20200908_tradeapi64_se_windows -DISLIB -DWIN32 -DLIB_MD_API_EXPORT -DLIB_TRADER_API_EXPORT
*/
import "C"