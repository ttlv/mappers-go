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
		nodeName = fmt.Sprintf("node-%s", strings.Split(td.DeviceInstanceName, "-")[2])
	}
	if td.Name == "acceleration" {
		ss1 := strings.Split(s2, " ")[0]
		ss2 := strings.Split(s2, " ")[1]
		ss3 := strings.Split(s2, " ")[2]
		ss4 := strings.Split(s2, " ")[3]
		ss5 := strings.Split(s2, " ")[4]
		ss6 := strings.Split(s2, " ")[5]
		axh, _ := strconv.Atoi(ss1)
		axl, _ := strconv.Atoi(ss2)
		ayh, _ := strconv.Atoi(ss3)
		ayl, _ := strconv.Atoi(ss4)
		azh, _ := strconv.Atoi(ss5)
		azl, _ := strconv.Atoi(ss6)
		k := 16.0
		accX := float64(axh<<8|axl) / 32768.0 * k
		accY := float64(ayh<<8|ayl) / 32768.0 * k
		accZ := float64(azh<<8|azl) / 32768.0 * k
		if accX >= k {
			accX -= 2 * k
		}
		if accY >= k {
			accY -= 2 * k
		}
		if accZ >= k {
			accZ -= 2 * k
		}
		klog.V(2).Info("---------accX-----------", accX)
		klog.V(2).Info("---------accY-----------", accY)
		klog.V(2).Info("---------accZ-----------", accZ)
		globals.FBClient.Publish(td.DeviceModel, fmt.Sprintf(`{"__name__":"%s","accX":%f,"accY":%f,"accZ":%f,"node":"%s","state":"%s"}`, td.DeviceModel, accX, accY, accZ, nodeName, td.Client.GetStatus()))
	} else if td.Name == "angularVelocity" {
		ss1 := strings.Split(s2, " ")[0]
		ss2 := strings.Split(s2, " ")[1]
		ss3 := strings.Split(s2, " ")[2]
		ss4 := strings.Split(s2, " ")[3]
		ss5 := strings.Split(s2, " ")[4]
		ss6 := strings.Split(s2, " ")[5]
		wxh, _ := strconv.Atoi(ss1)
		wxl, _ := strconv.Atoi(ss2)
		wyh, _ := strconv.Atoi(ss3)
		wyl, _ := strconv.Atoi(ss4)
		wzh, _ := strconv.Atoi(ss5)
		wzl, _ := strconv.Atoi(ss6)
		k := 2000.0
		wX := float64(wxh<<8|wxl) / 32768.0 * k
		wY := float64(wyh<<8|wyl) / 32768.0 * k
		wZ := float64(wzh<<8|wzl) / 32768.0 * k
		if wX >= k {
			wX -= 2 * k
		}
		if wY >= k {
			wY -= 2 * k
		}
		if wZ >= k {
			wZ -= 2 * k
		}
		klog.V(2).Info("---------wX-----------", wX)
		klog.V(2).Info("---------wY-----------", wY)
		klog.V(2).Info("---------wZ-----------", wZ)
		globals.FBClient.Publish(td.DeviceModel, fmt.Sprintf(`{"__name__":"%s","wX":%f,"wY":%f,"wZ":%f,"node":"%s","state":"%s"}`, td.DeviceModel, wX, wY, wZ, nodeName, td.Client.GetStatus()))
	} else if td.Name == "angular" {
		ss1 := strings.Split(s2, " ")[0]
		ss2 := strings.Split(s2, " ")[1]
		ss3 := strings.Split(s2, " ")[2]
		ss4 := strings.Split(s2, " ")[3]
		ss5 := strings.Split(s2, " ")[4]
		ss6 := strings.Split(s2, " ")[5]
		rollH, _ := strconv.Atoi(ss1)
		rollL, _ := strconv.Atoi(ss2)
		pitchH, _ := strconv.Atoi(ss3)
		pitchL, _ := strconv.Atoi(ss4)
		yawH, _ := strconv.Atoi(ss5)
		YawL, _ := strconv.Atoi(ss6)
		k := 180.0
		roll := float64(rollH<<8|rollL) / 32768.0 * k
		pitch := float64(pitchH<<8|pitchL) / 32768.0 * k
		yaw := float64(yawH<<8|YawL) / 32768.0 * k
		if roll >= k {
			roll -= 2 * k
		}
		if pitch >= k {
			pitch -= 2 * k
		}
		if yaw >= k {
			yaw -= 2 * k
		}
		klog.V(2).Info("---------roll-----------", roll)
		klog.V(2).Info("---------pitch-----------", pitch)
		klog.V(2).Info("---------yaw-----------", yaw)
		globals.FBClient.Publish(td.DeviceModel, fmt.Sprintf(`{"__name__":"%s","Roll":%f,"Pitch":%f,"Yaw":%f,"node":"%s","state":"%s"}`, td.DeviceModel, roll, pitch, yaw, nodeName, td.Client.GetStatus()))
	} else if td.Name == "magnetic" {
		ss1 := strings.Split(s2, " ")[0]
		ss2 := strings.Split(s2, " ")[1]
		ss3 := strings.Split(s2, " ")[2]
		ss4 := strings.Split(s2, " ")[3]
		ss5 := strings.Split(s2, " ")[4]
		ss6 := strings.Split(s2, " ")[5]
		hxH, _ := strconv.Atoi(ss1)
		hxL, _ := strconv.Atoi(ss2)
		hyH, _ := strconv.Atoi(ss3)
		hyL, _ := strconv.Atoi(ss4)
		hzH, _ := strconv.Atoi(ss5)
		hzL, _ := strconv.Atoi(ss6)
		k := 1.0
		hX := float64(hxH<<8 | hxL)
		hY := float64(hyH<<8 | hyL)
		hZ := float64(hzH<<8 | hzL)
		if hX >= k {
			hX -= 2 * k
		}
		if hY >= k {
			hY -= 2 * k
		}
		if hZ >= k {
			hZ -= 2 * k
		}
		klog.V(2).Info("---------hX-----------", hX)
		klog.V(2).Info("---------hY-----------", hY)
		klog.V(2).Info("---------hZ-----------", hZ)
		globals.FBClient.Publish(td.DeviceModel, fmt.Sprintf(`{"__name__":"%s","Hx":%f,"Hy":%f,"Hz":%f,"node":"%s","state":"%s"}`, td.DeviceModel, hX, hY, hZ, nodeName, td.Client.GetStatus()))
	} else if td.Name == "element" {
		ss1 := strings.Split(s2, " ")[0]
		ss2 := strings.Split(s2, " ")[1]
		ss3 := strings.Split(s2, " ")[2]
		ss4 := strings.Split(s2, " ")[3]
		ss5 := strings.Split(s2, " ")[4]
		ss6 := strings.Split(s2, " ")[5]
		ss7 := strings.Split(s2, " ")[6]
		ss8 := strings.Split(s2, " ")[7]
		q0H, _ := strconv.Atoi(ss1)
		q0L, _ := strconv.Atoi(ss2)
		q1H, _ := strconv.Atoi(ss3)
		q1L, _ := strconv.Atoi(ss4)
		q2H, _ := strconv.Atoi(ss5)
		q2L, _ := strconv.Atoi(ss6)
		q3H, _ := strconv.Atoi(ss7)
		q3L, _ := strconv.Atoi(ss8)
		k := 1.0
		q0 := float64(q0H<<8|q0L) / 32768.0
		q1 := float64(q1H<<8|q1L) / 32768.0
		q2 := float64(q2H<<8|q2L) / 32768.0
		q3 := float64(q3H<<8|q3L) / 32768.0
		if q0 >= k {
			q0 -= 2 * k
		}
		if q1 >= k {
			q1 -= 2 * k
		}
		if q2 >= k {
			q2 -= 2 * k
		}
		if q3 >= k {
			q3 -= 2 * k
		}
		klog.V(2).Info("---------q0-----------", q0)
		klog.V(2).Info("---------q1-----------", q1)
		klog.V(2).Info("---------q2-----------", q2)
		klog.V(2).Info("---------q3-----------", q3)
		globals.FBClient.Publish(td.DeviceModel, fmt.Sprintf(`{"__name__":"%s","Q0":%f,"Q1":%f,"Q2":%f,"Q3":%f,"node":"%s",state":"%s"}`, td.DeviceModel, q0, q1, q2, q3, nodeName, td.Client.GetStatus()))
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
