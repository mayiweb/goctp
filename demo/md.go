package main

/*
#include <stdio.h>
#include <stdlib.h>
*/
import "C"

import (
    "github.com/mayiweb/goctp"
    "fmt"
    "log"
    "unsafe"
)

// 获得行情请求编号
func (p *FtdcMdSpi) GetMdRequestId() int {
    Ctp.MdRequestId += 1
    return Ctp.MdRequestId
}

// 当客户端与交易后台通信连接断开时，该方法被调用。当发生这个情况后，API会自动重新连接，客户端可不做处理。
// 服务器已断线，该函数也会被调用。【api 会自动初始化程序，并重新登陆】
func (p *FtdcMdSpi) OnFrontDisconnected(nReason int) {
    log.Println("行情服务器已断线，尝试重新连接中...")
}

// 当客户端与交易后台建立起通信连接时（还未登录前），该方法被调用。
func (p *FtdcMdSpi) OnFrontConnected() {

    MdStr := "=================================================================================================\n" +
             "= 行情模块初始化成功，API 版本：" + goctp.CThostFtdcMdApiGetApiVersion() + "\n" +
             "================================================================================================="
    fmt.Println(MdStr)

    // 登录（如果行情模块在交易模块后初始化则直接登录行情）
    if Ctp.IsTraderInit {
        p.ReqUserLogin()
    }
}

// 行情用户登录
func (p *FtdcMdSpi) ReqUserLogin() {
    log.Println("行情系统账号登陆中...")

    req := goctp.NewCThostFtdcReqUserLoginField()
    req.SetBrokerID(BrokerID)
    req.SetUserID(InvestorID)
    req.SetPassword(Password)

    iResult := Ctp.MdApi.ReqUserLogin(req, p.GetMdRequestId())

    if iResult != 0 {
        ReqFailMsg("发送用户登录请求失败！", iResult)
    }
}

// 登录请求响应
func (p *FtdcMdSpi) OnRspUserLogin(pRspUserLogin goctp.CThostFtdcRspUserLoginField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

    if bIsLast && !p.IsErrorRspInfo(pRspInfo) {
        log.Printf("行情系统登陆成功，当前交易日： %v\n", Ctp.MdApi.GetTradingDay())
    }
}

// 订阅行情
func (p *FtdcMdSpi) SubscribeMarketData(InstrumentID []string) int {

    if len(InstrumentID) == 0 {
        log.Println("没有指定需要订阅的行情数据")
        return 0
    }

    fmt.Println("")
    log.Println("订阅行情数据中...")

    args := make([]*C.char, 0)
    for _, v := range InstrumentID {
        char := C.CString(v)
        defer C.free(unsafe.Pointer(char))
        args = append(args, char)
    }

    iResult := Ctp.MdApi.SubscribeMarketData((*string)(unsafe.Pointer(&args[0])), len(InstrumentID))

    if iResult != 0 {
        ReqFailMsg("发送订阅行情请求失败！", iResult)
    }

    return iResult
}

// 订阅行情应答
func (p *FtdcMdSpi) OnRspSubMarketData(pSpecificInstrument goctp.CThostFtdcSpecificInstrumentField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
    if !p.IsErrorRspInfo(pRspInfo) {
        log.Printf("订阅合约 %v 行情数据成功！\n", pSpecificInstrument.GetInstrumentID())
    }
}

// 退订行情
func (p *FtdcMdSpi) UnSubscribeMarketData(InstrumentID []string) int {

    if len(InstrumentID) == 0 {
        log.Println("没有指定需要退订的行情数据")
        return 0
    }

    fmt.Println("")
    log.Println("退订行情数据中...")

    args := make([]*C.char, 0)
    for _, v := range InstrumentID {
        char := C.CString(v)
        defer C.free(unsafe.Pointer(char))
        args = append(args, char)
    }

    iResult := Ctp.MdApi.UnSubscribeMarketData((*string)(unsafe.Pointer(&args[0])), len(InstrumentID))

    if iResult != 0 {
        ReqFailMsg("发送退订行情请求失败！", iResult)
    }

    return iResult
}

// 退订行情应答
func (p *FtdcMdSpi) OnRspUnSubMarketData(pSpecificInstrument goctp.CThostFtdcSpecificInstrumentField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
    if !p.IsErrorRspInfo(pRspInfo) {
        log.Printf("取消订阅 %v 行情数据成功！\n", pSpecificInstrument.GetInstrumentID())
    }
}

// 深度行情通知
func (p *FtdcMdSpi) OnRtnDepthMarketData(pDepthMarketData goctp.CThostFtdcDepthMarketDataField) {
    fmt.Printf("%v 合约：%v \t最新价：%v \t买一价：%v \t卖一价：%v \t买一量：%v \t卖一量：%v\n", pDepthMarketData.GetUpdateTime(), pDepthMarketData.GetInstrumentID(), pDepthMarketData.GetLastPrice(), pDepthMarketData.GetBidPrice1(), pDepthMarketData.GetAskPrice1(), pDepthMarketData.GetBidVolume1(), pDepthMarketData.GetAskVolume1())
}

// 行情系统错误通知
func (p *FtdcMdSpi) IsErrorRspInfo(pRspInfo goctp.CThostFtdcRspInfoField) bool {

    rspInfo := fmt.Sprintf("%v", pRspInfo)

    // 容错处理 pRspInfo ，部分响应函数中，pRspInfo 为 0
    if rspInfo == "0" {
        return false

    } else {

        // 如果ErrorID != 0, 说明收到了错误的响应
        bResult := (pRspInfo.GetErrorID() != 0)
        if bResult {
            // pRspInfo.GetErrorMsg 为 GBK 编码需要自行转成 utf8
            log.Printf("ErrorID=%v ErrorMsg=%v\n", pRspInfo.GetErrorID(), ConvertToString(pRspInfo.GetErrorMsg(), "gbk", "utf-8"))
        }

        return bResult
    }
}