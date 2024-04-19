# goctp
上海期货交易所 ctp 接口 Golang版 (for linux64)
支持 testctp 回测

## 修改配置
    修改 ctp.go 文件的 SetTradeAccount 函数，写上对应的环境账号即可

## 订阅行情
    Subscribe := []string{"rb2410"}
    MdSpi.SubscribeTick(Subscribe)

## 开仓示例
    var Input InputOrderStruct

    Input.InstrumentID = "rb2410"
    Input.Direction    = Buy
    Input.Price        = 3600
    Input.Volume       = 1

    TradeSpi.OrderOpen(Input)

## 平仓示例
    var Input InputOrderStruct

    Input.InstrumentID = "rb2410"
    Input.Direction    = Sell
    Input.Price        = 3600
    Input.Volume       = 1

    TradeSpi.OrderClose(Input)

## 撤单示例
    TradeSpi.OrderCancel("rb2410", "报单编号")