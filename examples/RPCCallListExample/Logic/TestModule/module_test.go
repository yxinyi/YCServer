package TestModule

import (
	"github.com/json-iterator/go"
	"github.com/yxinyi/YCServer/engine/YMsg"
	"github.com/yxinyi/YCServer/engine/YNode"
	"testing"
)

func init() {
	YNode.Register(NewInfo(YNode.Obj()))
	YNode.Start()
}

func TestModule(t *testing.T) {
	{
		msg := &YMsg.S2S_rpc_msg{}
		msg.M_func_name = "Test_3"
		msg.M_func_parameter = make([][]byte, 0)
		{
			
			_bytes, _ := jsoniter.Marshal(1)
			msg.M_func_parameter = append(msg.M_func_parameter, _bytes)
		}
		{
			_bytes, _ := jsoniter.Marshal("123")
			msg.M_func_parameter = append(msg.M_func_parameter, _bytes)
		}
		YNode.RPCCall(msg)
	}
	{
		msg := &YMsg.S2S_rpc_msg{}
		msg.M_func_name = "Test_4"
		msg.M_func_parameter = make([][]byte, 0)
		{
			
			_bytes, _ := jsoniter.Marshal(&YMsg.TestParam{
				123,
				"TESTPARAMTER",
				[]int{
					1, 2, 3, 4, 5, 6, 7,
				},
			})
			msg.M_func_parameter = append(msg.M_func_parameter, _bytes)
		}
		YNode.RPCCall(msg)
	}
}

func BenchmarkModule(b_ *testing.B) {
	for _idx := 0; _idx < b_.N; _idx++ {
		{
			msg := &YMsg.S2S_rpc_msg{}
			msg.M_func_name = "Test"
			msg.M_func_parameter = make([][]byte, 0)
			
			YNode.RPCCall(msg)
		}
		{
			msg := &YMsg.S2S_rpc_msg{}
			msg.M_func_name = "Test_2"
			msg.M_func_parameter = make([][]byte, 0)
			{
				
				_bytes, _ := jsoniter.Marshal(1)
				msg.M_func_parameter = append(msg.M_func_parameter, _bytes)
			}
			YNode.RPCCall(msg)
		}
		{
			msg := &YMsg.S2S_rpc_msg{}
			msg.M_func_name = "Test_3"
			msg.M_func_parameter = make([][]byte, 0)
			{
				
				_bytes, _ := jsoniter.Marshal(1)
				msg.M_func_parameter = append(msg.M_func_parameter, _bytes)
			}
			{
				_bytes, _ := jsoniter.Marshal("123")
				msg.M_func_parameter = append(msg.M_func_parameter, _bytes)
			}
			YNode.RPCCall(msg)
		}
		{
			msg := &YMsg.S2S_rpc_msg{}
			msg.M_func_name = "Test_4"
			msg.M_func_parameter = make([][]byte, 0)
			{
				
				_bytes, _ := jsoniter.Marshal(&YMsg.TestParam{
					123,
					"TESTPARAMTER",
					[]int{
						1, 2, 3, 4, 5, 6, 7,
					},
				})
				msg.M_func_parameter = append(msg.M_func_parameter, _bytes)
			}
			YNode.RPCCall(msg)
		}
	}
}
