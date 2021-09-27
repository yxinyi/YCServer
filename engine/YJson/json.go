package YJson

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

var g_ghost_json = jsoniter.Config{TagKey: "SG", OnlyTaggedField: true}.Froze()
var g_save_json = jsoniter.Config{TagKey: "SA", OnlyTaggedField: true}.Froze()
var g_sync_self_json = jsoniter.Config{TagKey: "SS", OnlyTaggedField: true}.Froze()
var g_sync_other_json = jsoniter.Config{TagKey: "SO", OnlyTaggedField: true}.Froze()

var g_print_json = jsoniter.Config{SortMapKeys: true, MarshalFloatWith6Digits: true, IndentionStep: 4}.Froze()

func GetPrintStr(val interface{}) string {
	_print_str, _err := g_print_json.MarshalToString(val)
	if _err != nil {
		return fmt.Sprintf("Marshal err [%v]", _err.Error())
	}
	return _print_str
}

func UnMarshal(str_ string,val_ interface{}) (error) {
	return jsoniter.UnmarshalFromString(str_,val_)
}

func GhostMarshal(val_ interface{}) (string, error) {
	return g_ghost_json.MarshalToString(val_)
}



func SaveMarshal(val_ interface{}) (string, error) {
	return g_save_json.MarshalToString(val_)
}
func SyncOtherMarshal(val_ interface{}) (string, error) {
	return g_sync_other_json.MarshalToString(val_)
}
func SyncSelfMarshal(val_ interface{}) (string, error) {
	return g_sync_self_json.MarshalToString(val_)
}
