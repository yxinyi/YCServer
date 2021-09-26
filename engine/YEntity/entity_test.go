package YEntity

import (
	"github.com/yxinyi/YCServer/engine/YAttr"
	"testing"
)

type TmpStruct struct {
	M_val int
}

func TestEntity(t *testing.T) {
	RegisterEntityAttr("Test",
		YAttr.Define("AttrPanel_1",
			YAttr.Tmpl("attr_1", TmpStruct{M_val: 9}, true, true, true, true),
		),
		YAttr.Define("AttrPanel_2",
			YAttr.Tmpl("attr_1", []uint32{1,2,3}, true, true, true, true),
			YAttr.Tmpl("attr_2", uint32(2), true, true, true, true),
			YAttr.Tmpl("attr_3", uint32(3), true, true, true, true),
			YAttr.Tmpl("attr_4", uint32(4), true, true, true, true),
			YAttr.Tmpl("attr_5", uint32(5), true, true, true, true),
		))
	_entity_1 := NewWithUID("Test", 1)
	{
		_tmp, _err := _entity_1.GetAttr("AttrPanel_1.attr_1").(*TmpStruct)
		if !_err {
			t.Errorf("[%v]", _err)
			return
		}
		t.Logf("[%v]", _tmp.M_val)
		_tmp.M_val = 100
	}
	{
		_tmp, _ := _entity_1.GetAttr("AttrPanel_1.attr_1").(*TmpStruct)
		t.Logf("[%v]", _tmp.M_val)
	}
	{
		_tmp, _ := _entity_1.GetAttr("AttrPanel_2.attr_1").(*[]uint32)
		t.Logf("[%v]", *_tmp)
	}
	{
		_tmp, _ := _entity_1.GetAttr("AttrPanel_2.attr_2").(*uint32)
		t.Logf("[%v]", *_tmp)
		*_tmp = 1000
	}
	{
		_tmp, _ := _entity_1.GetAttr("AttrPanel_2.attr_3").(*uint32)
		t.Logf("[%v]", *_tmp)
	}
	{
		_tmp, _ := _entity_1.GetAttr("AttrPanel_2.attr_4").(*uint32)
		t.Logf("[%v]", *_tmp)
	}
	{
		_tmp, _ := _entity_1.GetAttr("AttrPanel_2.attr_5").(*uint32)
		t.Logf("[%v]", *_tmp)
	}
	{
		_tmp, _ := _entity_1.GetAttr("AttrPanel_2.attr_2").(*uint32)
		t.Logf("[%v]", *_tmp)
	}
}
