/**
 * 队列任务模块
 */

package main


// 执行队列查询任务（例：新开仓或平仓后需要 查询持仓、查询资金）频率不能低于1秒，否则会触发交易所频繁查询导致查询失败。
func RunQueryTask() {

    for {

        // 需要放在最上面，防止后面未休眠就跳出循环，会导致 cpu 占用 100%
        Sleep(1200)

        // 交易程序初始化流程走完了才能执行队列任务（否则会出现同一秒多次查询操作，导致查询失败）
        if !Ctp.IsTradeInitFinish {
            continue
        }

        // 每次只处理一个任务
        task := Queue.Poll()
        if task != nil {

            switch task {

                // 查询资金账户
                case "ReqQryTradingAccount":
                    iResult := TradeSpi.ReqQryTradingAccount()

                    // iResult 不等于0，则表示查询失败了，重新加入到队列中
                    if iResult != 0 {
                        Queue.PushUnique("ReqQryTradingAccount")
                    }

                // 查询持仓汇总
                case "ReqQryInvestorPosition":
                    iResult := TradeSpi.ReqQryInvestorPosition()
                    if iResult != 0 {
                        Queue.PushUnique("ReqQryInvestorPosition")
                    }

                // 查询持仓明细
                case "ReqQryInvestorPositionDetail":
                    iResult := TradeSpi.ReqQryInvestorPositionDetail()
                    if iResult != 0 {
                        Queue.PushUnique("ReqQryInvestorPositionDetail")
                    }
            }
        }
    }
}