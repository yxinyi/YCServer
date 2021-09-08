module YNode

go 1.16


require YMsg v0.0.0-00010101000000-000000000000
require YLog v0.0.0-00010101000000-000000000000
require YNet v0.0.0-00010101000000-000000000000
require queue v0.0.0-00010101000000-000000000000
replace YMsg => ../../YMsg
replace YLog => ../../YLog
replace YNet => ../../YNet
replace queue => ../../YTool/queue

require YDecode v0.0.0-00010101000000-000000000000
replace YDecode => ../../YDecode

require YUIDFactory v0.0.0-00010101000000-000000000000
replace YUIDFactory => ../../YTool/UIDFactory

require YEntity v0.0.0-00010101000000-000000000000
replace YEntity => ../YEntity
require YModule v0.0.0-00010101000000-000000000000
replace YModule => ../YModule
require YNode v0.0.0-00010101000000-000000000000
replace YNode => ../YNode