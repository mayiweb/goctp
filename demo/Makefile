
build:
	go mod tidy -v
	go build -v -x -linkshared -o ctp *.go

error:
	# 如出现类似错误，去掉参数 -linkshared
	# WORK=/tmp/go-build***
	# readELFNote failed: EOF
	# make: *** [Makefile:7：build] 错误 1