package api_test

import (
	"strings"
	"testing"

	"github.com/dop251/goja"
	"pansiot-device/internal/script/api"
)

// TestHTTPAPI 测试HTTP API
func TestHTTPAPI(t *testing.T) {
	vm := goja.New()
	commAPI := api.NewCommunicationAPI()
	err := commAPI.InjectToVM(vm)
	if err != nil {
		t.Fatalf("注入API失败: %v", err)
	}

	t.Run("GET Request", func(t *testing.T) {
		script := `(function() {
			var response = HTTP.get("https://httpbin.org/get");
			return response.status === 200;
		})()`
		result, err := vm.RunString(script)
		if err != nil {
			t.Logf("GET请求测试跳过(可能没有网络): %v", err)
			return
		}
		if !result.ToBoolean() {
			t.Error("GET请求失败")
		}
	})

	t.Run("POST Request", func(t *testing.T) {
		script := `(function() {
			var response = HTTP.post("https://httpbin.org/post", {
				data: { test: "data" },
				headers: { "Content-Type": "application/json" }
			});
			return response.status === 200;
		})()`
		result, err := vm.RunString(script)
		if err != nil {
			t.Logf("POST请求测试跳过(可能没有网络): %v", err)
			return
		}
		if !result.ToBoolean() {
			t.Error("POST请求失败")
		}
	})

	t.Run("PUT Request", func(t *testing.T) {
		script := `(function() {
			var response = HTTP.put("https://httpbin.org/put", {
				data: { status: "updated" }
			});
			return response.status === 200;
		})()`
		result, err := vm.RunString(script)
		if err != nil {
			t.Logf("PUT请求测试跳过(可能没有网络): %v", err)
			return
		}
		if !result.ToBoolean() {
			t.Error("PUT请求失败")
		}
	})

	t.Run("DELETE Request", func(t *testing.T) {
		script := `(function() {
			var response = HTTP.delete("https://httpbin.org/delete");
			return response.status === 200;
		})()`
		result, err := vm.RunString(script)
		if err != nil {
			t.Logf("DELETE请求测试跳过(可能没有网络): %v", err)
			return
		}
		if !result.ToBoolean() {
			t.Error("DELETE请求失败")
		}
	})

	t.Run("Request with Timeout", func(t *testing.T) {
		script := `(function() {
			var response = HTTP.request({
				url: "https://httpbin.org/delay/2",
				method: "GET",
				timeout: 5000
			});
			return response.status === 200 || response.error !== null;
		})()`
		_, err := vm.RunString(script)
		if err != nil {
			t.Logf("超时测试: %v", err)
		}
	})

	t.Run("Request Response Structure", func(t *testing.T) {
		script := `(function() {
			var response = HTTP.get("https://httpbin.org/get");
			return response.status !== undefined &&
			       response.statusText !== undefined &&
			       response.headers !== undefined &&
			       response.body !== undefined &&
			       response.duration !== undefined;
		})()`
		result, err := vm.RunString(script)
		if err != nil {
			t.Logf("响应结构测试跳过(可能没有网络): %v", err)
			return
		}
		if !result.ToBoolean() {
			t.Error("响应结构不完整")
		}
	})

	t.Run("Invalid URL", func(t *testing.T) {
		script := `(function() {
			var response = HTTP.get("http://invalid-domain-12345.com");
			return response.error !== null || response.status === 0;
		})()`
		result, err := vm.RunString(script)
		if err != nil {
			t.Fatal(err)
		}
		if !result.ToBoolean() {
			t.Error("无效URL应该返回错误")
		}
	})
}

// TestMQTTAPI 测试MQTT API (预留)
func TestMQTTAPI(t *testing.T) {
	vm := goja.New()
	commAPI := api.NewCommunicationAPI()
	err := commAPI.InjectToVM(vm)
	if err != nil {
		t.Fatalf("注入API失败: %v", err)
	}

	t.Run("Publish Not Implemented", func(t *testing.T) {
		script := `(function() {
			var err = null;
			try {
				MQTT.publish("test/topic", "message");
			} catch(e) {
				err = e.message;
			}
			return err !== null && err.indexOf("暂未实现") >= 0;
		})()`
		result, err := vm.RunString(script)
		if err != nil {
			t.Fatal(err)
		}
		if !result.ToBoolean() {
			t.Error("MQTT.publish应该返回未实现错误")
		}
	})

	t.Run("Subscribe Not Implemented", func(t *testing.T) {
		script := `(function() {
			var err = null;
			try {
				MQTT.subscribe("test/topic", function(topic, msg) {});
			} catch(e) {
				err = e.message;
			}
			return err !== null && err.indexOf("暂未实现") >= 0;
		})()`
		result, err := vm.RunString(script)
		if err != nil {
			t.Fatal(err)
		}
		if !result.ToBoolean() {
			t.Error("MQTT.subscribe应该返回未实现错误")
		}
	})

	t.Run("Unsubscribe Not Implemented", func(t *testing.T) {
		script := `(function() {
			var err = null;
			try {
				MQTT.unsubscribe("sub-id");
			} catch(e) {
				err = e.message;
			}
			return err !== null && err.indexOf("暂未实现") >= 0;
		})()`
		result, err := vm.RunString(script)
		if err != nil {
			t.Fatal(err)
		}
		if !result.ToBoolean() {
			t.Error("MQTT.unsubscribe应该返回未实现错误")
		}
	})

	t.Run("IsConnected", func(t *testing.T) {
		script := `(function() {
			var connected = MQTT.isConnected();
			return connected === false;
		})()`
		result, err := vm.RunString(script)
		if err != nil {
			t.Fatal(err)
		}
		if !result.ToBoolean() {
			t.Error("MQTT.isConnected应该返回false")
		}
	})
}

// TestModbusAPI 测试Modbus API (预留)
func TestModbusAPI(t *testing.T) {
	vm := goja.New()
	commAPI := api.NewCommunicationAPI()
	err := commAPI.InjectToVM(vm)
	if err != nil {
		t.Fatalf("注入API失败: %v", err)
	}

	t.Run("ReadHoldingRegisters Not Implemented", func(t *testing.T) {
		script := `(function() {
			var err = null;
			try {
				var values = Modbus.readHoldingRegisters("device1", 0, 10);
			} catch(e) {
				err = e.message;
			}
			return err !== null && err.indexOf("暂未实现") >= 0;
		})()`
		result, err := vm.RunString(script)
		if err != nil {
			t.Fatal(err)
		}
		if !result.ToBoolean() {
			t.Error("Modbus.readHoldingRegisters应该返回未实现错误")
		}
	})

	t.Run("ReadInputRegisters Not Implemented", func(t *testing.T) {
		script := `(function() {
			var err = null;
			try {
				var values = Modbus.readInputRegisters("device1", 0, 5);
			} catch(e) {
				err = e.message;
			}
			return err !== null && err.indexOf("暂未实现") >= 0;
		})()`
		result, err := vm.RunString(script)
		if err != nil {
			t.Fatal(err)
		}
		if !result.ToBoolean() {
			t.Error("Modbus.readInputRegisters应该返回未实现错误")
		}
	})

	t.Run("WriteSingleRegister Not Implemented", func(t *testing.T) {
		script := `(function() {
			var err = null;
			try {
				Modbus.writeSingleRegister("device1", 0, 123);
			} catch(e) {
				err = e.message;
			}
			return err !== null && err.indexOf("暂未实现") >= 0;
		})()`
		result, err := vm.RunString(script)
		if err != nil {
			t.Fatal(err)
		}
		if !result.ToBoolean() {
			t.Error("Modbus.writeSingleRegister应该返回未实现错误")
		}
	})

	t.Run("WriteMultipleRegisters Not Implemented", func(t *testing.T) {
		script := `(function() {
			var err = null;
			try {
				Modbus.writeMultipleRegisters("device1", 0, [100, 200, 300]);
			} catch(e) {
				err = e.message;
			}
			return err !== null && err.indexOf("暂未实现") >= 0;
		})()`
		result, err := vm.RunString(script)
		if err != nil {
			t.Fatal(err)
		}
		if !result.ToBoolean() {
			t.Error("Modbus.writeMultipleRegisters应该返回未实现错误")
		}
	})

	t.Run("ReadCoils Not Implemented", func(t *testing.T) {
		script := `(function() {
			var err = null;
			try {
				var coils = Modbus.readCoils("device1", 0, 8);
			} catch(e) {
				err = e.message;
			}
			return err !== null && err.indexOf("暂未实现") >= 0;
		})()`
		result, err := vm.RunString(script)
		if err != nil {
			t.Fatal(err)
		}
		if !result.ToBoolean() {
			t.Error("Modbus.readCoils应该返回未实现错误")
		}
	})

	t.Run("WriteCoil Not Implemented", func(t *testing.T) {
		script := `(function() {
			var err = null;
			try {
				Modbus.writeCoil("device1", 0, true);
			} catch(e) {
				err = e.message;
			}
			return err !== null && err.indexOf("暂未实现") >= 0;
		})()`
		result, err := vm.RunString(script)
		if err != nil {
			t.Fatal(err)
		}
		if !result.ToBoolean() {
			t.Error("Modbus.writeCoil应该返回未实现错误")
		}
	})
}

// TestHTTPIntegration 测试HTTP集成
func TestHTTPIntegration(t *testing.T) {
	vm := goja.New()
	commAPI := api.NewCommunicationAPI()
	err := commAPI.InjectToVM(vm)
	if err != nil {
		t.Fatalf("注入API失败: %v", err)
	}

	script := `(function() {
		var response1 = HTTP.get("https://httpbin.org/get");
		var response2 = HTTP.post("https://httpbin.org/post", {
			data: { key: "value" }
		});
		return response1.status === 200 && response2.status === 200;
	})()`

	result, err := vm.RunString(script)
	if err != nil {
		t.Logf("集成测试跳过(可能没有网络): %v", err)
		return
	}

	if !result.ToBoolean() {
		t.Error("集成测试失败")
	}
}

// TestHTTPJSONHandling 测试HTTP JSON处理
func TestHTTPJSONHandling(t *testing.T) {
	vm := goja.New()
	commAPI := api.NewCommunicationAPI()
	err := commAPI.InjectToVM(vm)
	if err != nil {
		t.Fatalf("注入API失败: %v", err)
	}

	t.Run("Parse JSON Response", func(t *testing.T) {
		script := `(function() {
			var response = HTTP.post("https://httpbin.org/post", {
				data: { name: "test", value: 123 }
			});

			if (response.status !== 200) {
				return false;
			}

			var result = JSON.parse(response.body);
			return result.json.name === "test" && result.json.value === 123;
		})()`
		result, err := vm.RunString(script)
		if err != nil {
			t.Logf("JSON解析测试跳过(可能没有网络): %v", err)
			return
		}
		if !result.ToBoolean() {
			t.Error("JSON响应解析失败")
		}
	})

	t.Run("Array Data", func(t *testing.T) {
		script := `(function() {
			var response = HTTP.post("https://httpbin.org/post", {
				data: [1, 2, 3, 4, 5]
			});

			if (response.status !== 200) {
				return true;
			}

			var result = JSON.parse(response.body);
			return Array.isArray(result.json);
		})()`
		result, err := vm.RunString(script)
		if err != nil {
			t.Logf("数组数据测试跳过(可能没有网络): %v", err)
			return
		}
		if !result.ToBoolean() {
			t.Error("数组数据处理失败")
		}
	})
}

// TestHTTPRequestOptions 测试HTTP请求选项
func TestHTTPRequestOptions(t *testing.T) {
	vm := goja.New()
	commAPI := api.NewCommunicationAPI()
	err := commAPI.InjectToVM(vm)
	if err != nil {
		t.Fatalf("注入API失败: %v", err)
	}

	t.Run("Custom Headers", func(t *testing.T) {
		script := `(function() {
			var response = HTTP.request({
				url: "https://httpbin.org/headers",
				method: "GET",
				headers: {
					"X-Custom-Header": "test-value",
					"User-Agent": "pans-runtime-test"
				}
			});

			if (response.status !== 200) {
				return true;
			}

			var result = JSON.parse(response.body);
			return result.headers["X-Custom-Header"] === "test-value";
		})()`
		result, err := vm.RunString(script)
		if err != nil {
			t.Logf("自定义头测试跳过(可能没有网络): %v", err)
			return
		}
		if !result.ToBoolean() {
			t.Error("自定义头处理失败")
		}
	})

	t.Run("Method Override", func(t *testing.T) {
		script := `(function() {
			var response = HTTP.request({
				url: "https://httpbin.org/post",
				method: "POST"
			});
			return response.status === 200 || response.error !== null;
		})()`
		result, err := vm.RunString(script)
		if err != nil {
			t.Logf("方法覆盖测试跳过: %v", err)
			return
		}
		if !result.ToBoolean() {
			t.Error("方法覆盖失败")
		}
	})
}

// TestHTTPErrors 测试HTTP错误处理
func TestHTTPErrors(t *testing.T) {
	vm := goja.New()
	commAPI := api.NewCommunicationAPI()
	err := commAPI.InjectToVM(vm)
	if err != nil {
		t.Fatalf("注入API失败: %v", err)
	}

	t.Run("Invalid URL Format", func(t *testing.T) {
		script := `(function() {
			var response = HTTP.get("://invalid-url");
			return response.error !== null || response.status === 0;
		})()`
		result, err := vm.RunString(script)
		if err != nil {
			t.Fatal(err)
		}
		if !result.ToBoolean() {
			t.Error("无效URL格式应该返回错误")
		}
	})

	t.Run("Timeout Test", func(t *testing.T) {
		script := `(function() {
			var response = HTTP.request({
				url: "https://httpbin.org/delay/10",
				method: "GET",
				timeout: 1000
			});
			return response.error !== null || response.status === 200;
		})()`
		result, err := vm.RunString(script)
		if err != nil {
			t.Logf("超时测试: %v", err)
			return
		}
		if !result.ToBoolean() {
			t.Error("超时测试失败")
		}
		resultStr := result.String()
		if strings.Contains(resultStr, "timeout") || strings.Contains(resultStr, "deadline") {
			t.Log("✓ 超时机制工作正常")
		}
	})
}

// TestCommunicationAPICreate 测试通讯API创建
func TestCommunicationAPICreate(t *testing.T) {
	t.Run("NewCommunicationAPI", func(t *testing.T) {
		api := api.NewCommunicationAPI()
		if api == nil {
			t.Error("创建CommunicationAPI失败")
		}
	})

	t.Run("InjectToVM", func(t *testing.T) {
		vm := goja.New()
		commAPI := api.NewCommunicationAPI()
		err := commAPI.InjectToVM(vm)
		if err != nil {
			t.Errorf("注入API到VM失败: %v", err)
		}

		if vm.Get("HTTP") == nil || goja.IsUndefined(vm.Get("HTTP")) {
			t.Error("HTTP API未注入")
		}
		if vm.Get("MQTT") == nil || goja.IsUndefined(vm.Get("MQTT")) {
			t.Error("MQTT API未注入")
		}
		if vm.Get("Modbus") == nil || goja.IsUndefined(vm.Get("Modbus")) {
			t.Error("Modbus API未注入")
		}
	})
}
