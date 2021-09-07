package TestModule

import (
	"YMsg"
	"YServer/Frame/YNode"
	"encoding/json"
	"testing"
)
func init(){
	YNode.Register(newInfo())
	YNode.Start()
}

func TestModule(t *testing.T) {
	{
		msg := &YMsg.S2S_rpc_msg{}
		msg.M_func_name = "Test_3"
		msg.M_func_parameter = make([][]byte, 0)
		{
			
			_bytes, _ := json.Marshal(1)
			msg.M_func_parameter = append(msg.M_func_parameter, _bytes)
		}
		{
			_bytes, _ := json.Marshal("123")
			msg.M_func_parameter = append(msg.M_func_parameter, _bytes)
		}
		YNode.Dispatch(msg)
	}
	{
		msg := &YMsg.S2S_rpc_msg{}
		msg.M_func_name = "Test_4"
		msg.M_func_parameter = make([][]byte, 0)
		{
			
			_bytes, _ := json.Marshal(&YMsg.TestParam{
				123,
				"TESTPARAMTER",
				[]int{
					1,2,3,4,5,6,7,
				},
			})
			msg.M_func_parameter = append(msg.M_func_parameter, _bytes)
		}
		YNode.Dispatch(msg)
	}
}

func BenchmarkModule(b_ *testing.B) {
	for _idx := 0 ;_idx < b_.N;_idx++ {
		{
			msg := &YMsg.S2S_rpc_msg{}
			msg.M_func_name = "Test"
			msg.M_func_parameter = make([][]byte, 0)
			
			YNode.Dispatch(msg)
		}
		{
			msg := &YMsg.S2S_rpc_msg{}
			msg.M_func_name = "Test_2"
			msg.M_func_parameter = make([][]byte, 0)
			{
				
				_bytes, _ := json.Marshal(1)
				msg.M_func_parameter = append(msg.M_func_parameter, _bytes)
			}
			YNode.Dispatch(msg)
		}
		{
			msg := &YMsg.S2S_rpc_msg{}
			msg.M_func_name = "Test_3"
			msg.M_func_parameter = make([][]byte, 0)
			{
				
				_bytes, _ := json.Marshal(1)
				msg.M_func_parameter = append(msg.M_func_parameter, _bytes)
			}
			{
				_bytes, _ := json.Marshal("123")
				msg.M_func_parameter = append(msg.M_func_parameter, _bytes)
			}
			YNode.Dispatch(msg)
		}
		{
			msg := &YMsg.S2S_rpc_msg{}
			msg.M_func_name = "Test_4"
			msg.M_func_parameter = make([][]byte, 0)
			{
				
				_bytes, _ := json.Marshal(&YMsg.TestParam{
					123,
					"TESTPARAMTER",
					[]int{
						1,2,3,4,5,6,7,
					},
				})
				msg.M_func_parameter = append(msg.M_func_parameter, _bytes)
			}
			YNode.Dispatch(msg)
		}
	}
}
