# goctp
This library pre-builds the CTP required shared-objects using cgo.

## Prerequisites
    install go
    install swig

## Swig file
```
%% goctp.swigcxx

%module(directors="1") goctp
%{
#include "./api/ThostFtdcUserApiDataType.h"
#include "./api/ThostFtdcUserApiStruct.h"
#include "./api/ThostFtdcTraderApi.h"
#include "./api/ThostFtdcMdApi.h"
%}

%feature("director") CThostFtdcMdSpi;
%feature("director") CThostFtdcTraderSpi;

%include "./api/ThostFtdcUserApiDataType.h"
%include "./api/ThostFtdcUserApiStruct.h"
%include "./api/ThostFtdcTraderApi.h"
%include "./api/ThostFtdcMdApi.h"
```

If you want to reproduce the same results, copy the above file excluding the `%% goctp.swigcxx` line
and put the file directly under the project folder and then run `swig -c++ -cgo -go -intgosize 64 goctp.swigcxx`.

## 部署发布
    生成 ctp 文件后，使用 ldd ctp 命令查看依赖关系，并将如下文件放在同一文件夹（注意需要有执行权限）:
        ctp
        libruntime,sync-atomic.so
        libthostmduserapi_se.so
        libthosttraderapi_se.so

    将部署目录路径写入 /etc/ld.so.conf 文件最后一行，并执行 /sbin/ldconfig 命令