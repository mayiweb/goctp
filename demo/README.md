# goctp
上海期货交易所 ctp 接口 Golang版 (for linux64)

## 修改配置
    修改 ctp.go 文件 107 行的 SetTradeAccount 函数，写上对应的环境账号即可

## 订阅行情
    Subscribe := []string{"rb2001"}
    MdSpi.SubscribeMarketData(Subscribe)

## 开仓示例
    var Input InputOrderStruct

    Input.InstrumentID = "rb2001"
    Input.Direction    = OrderBuy
    Input.Price        = 3600
    Input.Volume       = 1

    TraderSpi.OrderOpen(Input)

## 平仓示例
    var Input InputOrderStruct

    Input.InstrumentID = "rb2001"
    Input.Direction    = OrderBuy
    Input.Price        = 3600
    Input.Volume       = 1

    TraderSpi.OrderClose(Input)

## 撤单示例
    TraderSpi.OrderCancel("rb2001", "报单编号")