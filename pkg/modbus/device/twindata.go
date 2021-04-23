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

package device

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/hoisie/mustache"
	"github.com/kubeedge/mappers-go/pkg/common"
	"github.com/kubeedge/mappers-go/pkg/modbus/driver"
	"github.com/kubeedge/mappers-go/pkg/modbus/globals"
	"k8s.io/klog"
)

// TwinData is the timer structure for getting twin/data.
type TwinData struct {
	Client             *driver.ModbusClient
	Name               string
	Type               string
	RegisterType       string
	Address            uint16
	Quantity           uint16
	Results            []byte
	Topic              string
	DeviceModel        string
	DeviceInstanceName string
}

// Run timer function.
func (td *TwinData) Run() error {
	var (
		err error
	)
	td.Results, err = td.Client.Get(td.RegisterType, td.Address, td.Quantity)
	// 访问失败之后，继续访问，访问10次，如果10次全部失败说明设备或者串口不可用，直接retuen
	if err != nil {
		for i := 0; i <= 9; i++ {
			if td.Results, err = td.Client.Get(td.RegisterType, td.Address, td.Quantity); err == nil {
				break
			}
			if i == 9 {
				content := mustache.Render(`
                节点编号: {{nodeName}}
                节点位置: {{location}}
                异常设备: {{abnormal}}
`, map[string]string{
					"nodeName": fmt.Sprintf("node-%v", strings.Split(td.DeviceInstanceName, "-")[2]),
					"location": globals.NodeDetail[fmt.Sprintf("node-%v", strings.Split(td.DeviceInstanceName, "-")[2])],
					"abnormal": strings.Split(td.DeviceInstanceName, "-")[1],
				})
				globals.DingTalkClient.Robot.SendText(content, globals.AtMobiles, false)
				return fmt.Errorf("设备不可用")
			}
		}
	}
	s1 := strings.Replace(fmt.Sprintf("%v", td.Results), "[", "", -1)
	s2 := strings.Replace(s1, "]", "", -1)
	splitS2 := strings.Split(s2, " ")
	var nodeName string
	if len(strings.Split(td.DeviceInstanceName, "-")) == 3 && strings.Split(td.DeviceInstanceName, "-")[2] != "" {
		nodeName = strings.Split(td.DeviceInstanceName, "-")[2]
	}
	if strings.Contains(td.DeviceInstanceName, "shutter01") {
		var lux, co2, pressure, temperature, humidity float64
		// 湿度
		ss1 := splitS2[0]
		ss2 := splitS2[1]
		humidity = Hex2Dec(ss1, ss2) / 10
		if humidity < 0 || humidity > 100 {
			content := mustache.Render(`
                节点编号: {{nodeName}}
                节点位置: {{location}}
                异常设备: {{abnormal}}
                异常属性: humidity[湿度]
                异常原因: {{value}}	
`, map[string]string{
				"nodeName": fmt.Sprintf("node-%v", strings.Split(td.DeviceInstanceName, "-")[2]),
				"location": globals.NodeDetail[fmt.Sprintf("node-%v", strings.Split(td.DeviceInstanceName, "-")[2])],
				"abnormal": strings.Split(td.DeviceInstanceName, "-")[1],
				"value":    fmt.Sprintf("%f", humidity),
			})
			globals.DingTalkClient.Robot.SendText(content, globals.AtMobiles, false)

		}
		// 温度
		ss3 := splitS2[2]
		ss4 := splitS2[3]
		temperature = Hex2Dec(ss3, ss4) / 10
		if temperature < -273 || temperature >= 6553.5 {
			content := mustache.Render(`
                节点编号: {{nodeName}}
                节点位置: {{location}}
                异常设备: {{abnormal}}
                异常属性: temperature[温度]
                异常原因: {{value}}
`, map[string]string{
				"nodeName": fmt.Sprintf("node-%v", strings.Split(td.DeviceInstanceName, "-")[2]),
				"location": globals.NodeDetail[fmt.Sprintf("node-%v", strings.Split(td.DeviceInstanceName, "-")[2])],
				"abnormal": strings.Split(td.DeviceInstanceName, "-")[1],
				"value":    fmt.Sprintf("%f", temperature),
			})
			globals.DingTalkClient.Robot.SendText(content, globals.AtMobiles, false)
		}
		// 光强
		ss5 := splitS2[4]
		ss6 := splitS2[5]
		ss7 := splitS2[6]
		ss8 := splitS2[7]
		lux = Hex2Dec(ss5, ss6, ss7, ss8)
		if lux <= 0 || lux > 5000 {
			content := mustache.Render(`
                节点编号: {{nodeName}}
                节点位置: {{location}}
                异常设备: {{abnormal}}
                异常属性: lux[光照强度]
                异常原因: {{value}}
`, map[string]string{
				"nodeName": fmt.Sprintf("node-%v", strings.Split(td.DeviceInstanceName, "-")[2]),
				"location": globals.NodeDetail[fmt.Sprintf("node-%v", strings.Split(td.DeviceInstanceName, "-")[2])],
				"abnormal": strings.Split(td.DeviceInstanceName, "-")[1],
				"value":    fmt.Sprintf("%f", lux),
			})
			globals.DingTalkClient.Robot.SendText(content, globals.AtMobiles, false)
		}
		// CO2
		ss11 := splitS2[14]
		ss12 := splitS2[15]
		co2 = Hex2Dec(ss11, ss12)
		if co2 <= 0 || co2 > 5000 {
			content := mustache.Render(`
                节点编号: {{nodeName}}
                节点位置: {{location}}
                异常设备: {{abnormal}}
                异常属性: co2[二氧化碳浓度]
                异常原因: {{value}}
`, map[string]string{
				"nodeName": fmt.Sprintf("node-%v", strings.Split(td.DeviceInstanceName, "-")[2]),
				"location": globals.NodeDetail[fmt.Sprintf("node-%v", strings.Split(td.DeviceInstanceName, "-")[2])],
				"abnormal": strings.Split(td.DeviceInstanceName, "-")[1],
				"value":    fmt.Sprintf("%f", co2),
			})
			globals.DingTalkClient.Robot.SendText(content, globals.AtMobiles, false)
		}
		// 大气压强
		ss13 := splitS2[22]
		ss14 := splitS2[23]
		pressure = Hex2Dec(ss13, ss14) / 10
		if pressure <= 0 || pressure > 1100 {
			content := mustache.Render(`
                节点编号: {{nodeName}}
                节点位置: {{location}}
                异常设备: {{abnormal}}
                异常属性: pressure[大气压强]
                异常原因: {{value}}
`, map[string]string{
				"nodeName": fmt.Sprintf("node-%v", strings.Split(td.DeviceInstanceName, "-")[2]),
				"location": globals.NodeDetail[fmt.Sprintf("node-%v", strings.Split(td.DeviceInstanceName, "-")[2])],
				"abnormal": strings.Split(td.DeviceInstanceName, "-")[1],
				"value":    fmt.Sprintf("%f", pressure),
			})
			globals.DingTalkClient.Robot.SendText(content, globals.AtMobiles, false)

		}
		klog.V(2).Info("---------湿度-----------", humidity)
		klog.V(2).Info("---------温度-----------", temperature)
		klog.V(2).Info("---------光强-----------", lux)
		klog.V(2).Info("---------二氧化碳浓度-----------", co2)
		klog.V(2).Info("---------大气压强-----------", pressure)
		globals.FBClient.Publish(td.DeviceInstanceName, fmt.Sprintf(`{"node":"%s", "__name__":"%s", "humidity":%f, "temperature":%f, "lux":%f, "co2":%f, "pressure":%f, "state":"%s"}`, nodeName, td.DeviceModel, humidity, temperature, lux, co2, pressure, td.Client.GetStatus()))
	} else if strings.Contains(td.DeviceInstanceName, "shutter02") {
		var pm2point5, pm10, noise float64
		// 噪音
		ss9 := splitS2[8]
		ss10 := splitS2[9]
		noise = Hex2Dec(ss9, ss10) / 10
		if noise <= 0 || noise >= 6553.5 {
			content := mustache.Render(`
                节点编号: {{nodeName}}
                节点位置: {{location}}
                异常设备: {{abnormal}}
                异常属性: noise[噪音]
                异常原因: {{value}}
`, map[string]string{
				"nodeName": fmt.Sprintf("node-%v", strings.Split(td.DeviceInstanceName, "-")[2]),
				"location": globals.NodeDetail[fmt.Sprintf("node-%v", strings.Split(td.DeviceInstanceName, "-")[2])],
				"abnormal": strings.Split(td.DeviceInstanceName, "-")[1],
				"value":    fmt.Sprintf("%f", noise),
			})
			globals.DingTalkClient.Robot.SendText(content, globals.AtMobiles, false)
		}
		// PM2.5
		ss15 := splitS2[40]
		ss16 := splitS2[41]
		pm2point5 = Hex2Dec(ss15, ss16)
		if pm2point5 <= 0 || pm2point5 >= 6553.5 {
			content := mustache.Render(`
                节点编号: {{nodeName}}
                节点位置: {{location}}
                异常设备: {{abnormal}}
                异常属性: pm2.5[pm2.5]
                异常原因: {{value}}
`, map[string]string{
				"nodeName": fmt.Sprintf("node-%v", strings.Split(td.DeviceInstanceName, "-")[2]),
				"location": globals.NodeDetail[fmt.Sprintf("node-%v", strings.Split(td.DeviceInstanceName, "-")[2])],
				"abnormal": strings.Split(td.DeviceInstanceName, "-")[1],
				"value":    fmt.Sprintf("%f", pm2point5),
			})
			globals.DingTalkClient.Robot.SendText(content, globals.AtMobiles, false)
		}
		// PM 10
		ss17 := splitS2[42]
		ss18 := splitS2[43]
		pm10 = Hex2Dec(ss17, ss18)
		if pm10 <= 0 || pm10 >= 6553.5 {
			content := mustache.Render(`
                节点编号: {{nodeName}}
                节点位置: {{location}}
                异常设备: {{abnormal}}
                异常属性: pm10[pm10]
                异常原因: {{value}}
`, map[string]string{
				"nodeName": fmt.Sprintf("node-%v", strings.Split(td.DeviceInstanceName, "-")[2]),
				"location": globals.NodeDetail[fmt.Sprintf("node-%v", strings.Split(td.DeviceInstanceName, "-")[2])],
				"abnormal": strings.Split(td.DeviceInstanceName, "-")[1],
				"value":    fmt.Sprintf("%f", pm10),
			})
			globals.DingTalkClient.Robot.SendText(content, globals.AtMobiles, false)
		}
		klog.V(2).Info("---------噪音-----------", noise)
		klog.V(2).Info("---------PM2.5-----------", pm2point5)
		klog.V(2).Info("---------PM10-----------", pm10)
		globals.FBClient.Publish(td.DeviceInstanceName, fmt.Sprintf(`{"node":"%s", "__name__":"%s", "nosie":%f, "pm2.5":%f, "pm10":%f, "state":"%s"}`, nodeName, td.DeviceModel, noise, pm2point5, pm10, td.Client.GetStatus()))
	} else if strings.Contains(td.DeviceInstanceName, "snow") {
		var snow float64
		ss1 := strings.Split(s2, " ")[0]
		ss2 := strings.Split(s2, " ")[1]
		snow = Hex2Dec(ss1, ss2)
		if snow != 0 && snow != 1 {
			content := mustache.Render(`
                节点编号: {{nodeName}}
                节点位置: {{location}}
                异常设备: {{abnormal}}
                异常属性: snow[雨雪]
                异常原因: {{value}}
`, map[string]string{
				"nodeName": fmt.Sprintf("node-%v", strings.Split(td.DeviceInstanceName, "-")[2]),
				"location": globals.NodeDetail[fmt.Sprintf("node-%v", strings.Split(td.DeviceInstanceName, "-")[2])],
				"abnormal": strings.Split(td.DeviceInstanceName, "-")[1],
				"value":    fmt.Sprintf("%f", snow),
			})
			globals.DingTalkClient.Robot.SendText(content, globals.AtMobiles, false)

		}
		klog.V(2).Info("----------雨雪--------------", snow)
		globals.FBClient.Publish(td.DeviceInstanceName, fmt.Sprintf(`{"__name__":"%s","snow":%f,"node":"%s","state":"%s"}`, td.DeviceModel, snow, nodeName, td.Client.GetStatus()))

	}
	var payload []byte
	if strings.Contains(td.Topic, "$hw") {
		if payload, err = common.CreateMessageTwinUpdate(td.Name, td.Type, strconv.Itoa(int(td.Results[0]))); err != nil {
			klog.Error("Create message twin update failed")
			return err
		}
	} else {
		if payload, err = common.CreateMessageData(td.Name, td.Type, strconv.Itoa(int(td.Results[0]))); err != nil {
			klog.Error("Create message data failed")
			return err
		}
	}
	if err = globals.MqttClient.Publish(td.Topic, payload); err != nil {
		klog.Error(err)
	}

	klog.V(2).Infof("Update value: %s, topic: %s", strconv.Itoa(int(td.Results[0])), td.Topic)
	return err
}

func Hex2Dec(vals ...string) float64 {
	var result float64
	for index, val := range vals {
		floatVal, _ := strconv.ParseFloat(val, 64)
		result += math.Pow(256, float64(len(vals)-index)-1) * floatVal
	}
	return result
}
