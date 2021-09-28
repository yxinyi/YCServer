package YTool

import (
	"github.com/yxinyi/YCServer/engine/YJson"
	ylog "github.com/yxinyi/YCServer/engine/YLog"
)

func JsonPrint(args_ ...interface{}) {
	for _, _arg_it := range args_ {
		ylog.Info("\n%s", YJson.GetPrintStr(_arg_it))
	}
}
