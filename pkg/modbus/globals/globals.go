/*
Copyright 2020 The KubeEdge Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package globals

import (
	"github.com/kubeedge/mappers-go/pkg/common"
	"github.com/kubeedge/mappers-go/pkg/modbus/driver"
)

// ModbusDev is the modbus device configuration and client information.
type ModbusDev struct {
	Instance     common.DeviceInstance
	ModbusClient *driver.ModbusClient
}

var MqttClient common.MqttClient
var FBClient common.MqttClient
var DingTalkClient common.DingTalkClient
var NodeDetail map[string]string
var AtMobiles = []string{"18626860751"}

const (
	Modbus = "modbus"
)

func init() {
	NodeDetail = make(map[string]string)
	NodeDetail["node-4e92b0ae0c01024c3be1"] = "点4，8号"
	NodeDetail["node-93217f62d5c4fd7221d7"] = "22.283114452,N,113.736734930,E（点3 10号）"
	NodeDetail["node-9d8f72a02af0f80f113a"] = "22.283249247,N,113.731630695,E（点8 9号)"
	NodeDetail["node-242c5790f805ed7d06f8"] = "22.282939860,N,113.731631707,E （点7 11号）"
}
