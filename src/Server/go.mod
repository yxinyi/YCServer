module YServer

go 1.16

require (
	YMsg v0.0.0-00010101000000-000000000000
	YNet v0.0.0-00010101000000-000000000000
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/lestrrat/go-file-rotatelogs v0.0.0-20180223000712-d3151e2a480f
	github.com/lestrrat/go-strftime v0.0.0-20180220042222-ba3bf9c1d042 // indirect
	github.com/yxinyi/YEventBus v0.1.2
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.18.1
	queue v0.0.0-00010101000000-000000000000
)

replace YNet => ../Base/YNet

replace YMsg => ../Base/YMsg

replace queue => ../Base/YTool/queue
