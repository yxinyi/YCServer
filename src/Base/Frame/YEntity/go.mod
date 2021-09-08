module YEntity

go 1.16

require YMsg v0.0.0-00010101000000-000000000000
require YLog v0.0.0-00010101000000-000000000000
require YEntity v0.0.0-00010101000000-000000000000
require YModule v0.0.0-00010101000000-000000000000
require YNet v0.0.0-00010101000000-000000000000
require YDecode v0.0.0-00010101000000-000000000000
require queue v0.0.0-00010101000000-000000000000
replace YEntity => ../YEntity
replace YModule => ../YModule
replace YMsg => ../../YMsg
replace YLog => ../../YLog
replace YNet => ../../YNet
replace queue => ../../YTool/queue
replace YDecode => ../../YDecode
