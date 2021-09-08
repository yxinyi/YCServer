module YServer

go 1.16

require (
	YLog v0.0.0-00010101000000-000000000000
	YModule v0.0.0-00010101000000-000000000000
	YMsg v0.0.0-00010101000000-000000000000
	YNet v0.0.0-00010101000000-000000000000
	YNode v0.0.0-00010101000000-000000000000
	YTimer v0.0.0-00010101000000-000000000000
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/lestrrat/go-envload v0.0.0-20180220120943-6ed08b54a570 // indirect
	github.com/lestrrat/go-file-rotatelogs v0.0.0-20180223000712-d3151e2a480f // indirect
	github.com/lestrrat/go-strftime v0.0.0-20180220042222-ba3bf9c1d042 // indirect
	github.com/tebeka/strftime v0.1.5 // indirect
	github.com/yxinyi/YEventBus v0.1.2
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	queue v0.0.0-00010101000000-000000000000
)

replace YNet => ../Base/YNet

replace YMsg => ../Base/YMsg

replace YTimer => ../Base/YTimer

replace YNode => ../Base/Frame/YNode

replace YModule => ../Base/Frame/YModule

replace YEntity => ../Base/Frame/YEntity

replace queue => ../Base/YTool/queue

replace YLog => ../Base/YLog

require YDecode v0.0.0-00010101000000-000000000000

replace YDecode => ../Base/YDecode

require YUIDFactory v0.0.0-00010101000000-000000000000

replace YUIDFactory => ../Base/YTool/UIDFactory
