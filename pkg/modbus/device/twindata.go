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
	if err != nil {
		klog.Error("Get register failed: ", err)
		return err
	}
	s1 := strings.Replace(fmt.Sprintf("%v", td.Results), "[", "", -1)
	s2 := strings.Replace(s1, "]", "", -1)
	var nodeName string
	if len(strings.Split(td.DeviceInstanceName, "-")) == 3 && strings.Split(td.DeviceInstanceName, "-")[2] != "" {
		nodeName = strings.Split(td.DeviceInstanceName, "-")[2]
	}
	if strings.Contains(td.DeviceInstanceName, "device-shutter01") {
		var lux, co2, pressure, temperature, humidity float64
		if td.Name == "lux" {
			ss1 := strings.Split(s2, " ")[0]
			ss2 := strings.Split(s2, " ")[1]
			ss3 := strings.Split(s2, " ")[2]
			ss4 := strings.Split(s2, " ")[3]
			lux = Hex2Dec(ss1, ss2, ss3, ss4)
			fmt.Println("----------光强----------", lux)
			globals.FBClient.Publish(td.DeviceModel, fmt.Sprintf(`{"__name__":"%s","lux":%f,"node":"%s","state":"%s"}`, td.DeviceModel, lux, nodeName, td.Client.GetStatus()))
		} else if td.Name == "co2" {
			ss1 := strings.Split(s2, " ")[0]
			ss2 := strings.Split(s2, " ")[1]
			co2 = Hex2Dec(ss1, ss2)
			fmt.Println("----------co2----------", co2)
			globals.FBClient.Publish(td.DeviceModel, fmt.Sprintf(`{"__name__":"%s","co2":%f,"node":"%s","state":"%s"}`, td.DeviceModel, co2, nodeName, td.Client.GetStatus()))
		} else if td.Name == "pressure" {
			ss1 := strings.Split(s2, " ")[0]
			ss2 := strings.Split(s2, " ")[1]
			pressure = Hex2Dec(ss1, ss2)
			fmt.Println("----------压强----------", pressure/10)
			globals.FBClient.Publish(td.DeviceModel, fmt.Sprintf(`{"__name__":"%s","pressure":%f,"node":"%s","state":"%s"}`, td.DeviceModel, pressure/10, nodeName, td.Client.GetStatus()))
		} else if td.Name == "temperature" {
			ss1 := strings.Split(s2, " ")[0]
			ss2 := strings.Split(s2, " ")[1]
			humidity = Hex2Dec(ss1, ss2)
			fmt.Println("----------湿度----------", humidity/10)
			globals.FBClient.Publish(td.DeviceModel, fmt.Sprintf(`{"__name__":"%s","humidity":%f,"node":"%s","state":"%s"}`, td.DeviceModel, humidity/10, nodeName, td.Client.GetStatus()))
		} else if td.Name == "humidity" {
			ss1 := strings.Split(s2, " ")[2]
			ss2 := strings.Split(s2, " ")[3]
			temperature = Hex2Dec(ss1, ss2)
			fmt.Println("----------温度----------", temperature/10)
			globals.FBClient.Publish(td.DeviceModel, fmt.Sprintf(`{"__name__":"%s","temperature":%f,"node":"%s","state":"%s"}`, td.DeviceModel, temperature/10, nodeName, td.Client.GetStatus()))
		}
	} else if strings.Contains(td.DeviceInstanceName, "device-shutter02") {
		var pm2point5, pm10, noise float64
		if td.Name == "pm2.5" {
			ss1 := strings.Split(s2, " ")[0]
			ss2 := strings.Split(s2, " ")[1]
			pm2point5 = Hex2Dec(ss1, ss2)
			fmt.Println("----------pm2.5--------------", pm2point5)
			globals.FBClient.Publish(td.DeviceModel, fmt.Sprintf(`{"__name__":"%s","pm2point5":%f,"node":"%s","state":"%s"}`, td.DeviceModel, pm2point5, nodeName, td.Client.GetStatus()))
		} else if td.Name == "pm10" {
			ss1 := strings.Split(s2, " ")[2]
			ss2 := strings.Split(s2, " ")[3]
			pm10 = Hex2Dec(ss1, ss2)
			fmt.Println("----------pm10--------------", pm10)
			globals.FBClient.Publish(td.DeviceModel, fmt.Sprintf(`{"__name__":"%s","pm2point5":%f,"node":"%s","state":"%s"}`, td.DeviceModel, pm10, nodeName, td.Client.GetStatus()))
		} else if td.Name == "noise" {
			ss1 := strings.Split(s2, " ")[0]
			ss2 := strings.Split(s2, " ")[1]
			noise = Hex2Dec(ss1, ss2)
			fmt.Println("----------噪音--------------", noise/10)
			globals.FBClient.Publish(td.DeviceModel, fmt.Sprintf(`{"__name__":"%s","noise":%f,"node":"%s","state":"%s"}`, td.DeviceModel, noise/10, nodeName, td.Client.GetStatus()))
		}
	} else if strings.Contains(td.DeviceInstanceName, "device-snow") {
		var snow float64
		ss1 := strings.Split(s2, " ")[0]
		ss2 := strings.Split(s2, " ")[1]
		snow = Hex2Dec(ss1, ss2)
		fmt.Println("----------雨雪--------------", snow)
		globals.FBClient.Publish(td.DeviceModel, fmt.Sprintf(`{"__name__":"%s","snow":%f,"node":"%s","state":"%s"}`, td.DeviceModel, snow, nodeName, td.Client.GetStatus()))
	}
	// construct payload
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
