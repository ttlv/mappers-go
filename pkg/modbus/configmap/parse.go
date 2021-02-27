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

package configmap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/kubeedge/mappers-go/pkg/common"
	"github.com/kubeedge/mappers-go/pkg/modbus/globals"
	"k8s.io/klog"
)

// Parse parse the configmap.
func Parse(path string,
	devices map[string]*globals.ModbusDev,
	dms map[string]common.DeviceModel,
	protocols map[string]common.Protocol) error {
	var deviceProfile common.DeviceProfile

	jsonFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(jsonFile, &deviceProfile); err != nil {
		return err
	}

	for _, deviceInstance := range deviceProfile.DeviceInstances {
		for index, protocol := range deviceProfile.Protocols {
			var protocolFound bool
			if protocol.Protocol != globals.Modbus {
				continue
			}
			if deviceInstance.ProtocolName == protocol.Name {
				protocolFound = true
				deviceInstance.PProtocol = protocol
				break
			}
			if !protocolFound && index == len(deviceProfile.Protocols)-1 {
				return fmt.Errorf("Protocol not found")
			}
		}
		for _, propertyVisitor := range deviceInstance.PropertyVisitors {
			modelName := propertyVisitor.ModelName
			propertyName := propertyVisitor.PropertyName
			for index, deviceModel := range deviceProfile.DeviceModels {
				var deviceModelFound bool
				if modelName == deviceModel.Name {
					deviceModelFound = true
					for index, property := range deviceModel.Properties {
						var propertyFound bool
						if propertyName == property.Name {
							propertyFound = true
							propertyVisitor.PProperty = property
							break
						}
						if !propertyFound && index == len(deviceModel.Properties)-1 {
							return fmt.Errorf("Property not found")
						}
					}
					break
				}
				if !deviceModelFound && index == len(deviceProfile.DeviceModels)-1 {
					return fmt.Errorf("Device model not found")
				}
			}
		}
		for _, twin := range deviceInstance.Twins {
			propertyName := twin.PropertyName
			for index, propertyVisitor := range deviceInstance.PropertyVisitors {
				var propertyNameFound bool
				if propertyName == propertyVisitor.PropertyName {
					propertyNameFound = true
					twin.PVisitor = &propertyVisitor
					break
				}
				if !propertyNameFound && index == len(deviceInstance.PropertyVisitors)-1 {
					return fmt.Errorf("PropertyVisitor not found")
				}
			}
		}

		deviceInstance.Datas.Properties = deviceInstance.Properties
		deviceInstance.Datas.Topic = deviceInstance.Topic

		for _, property := range deviceInstance.Datas.Properties {
			propertyName := property.PropertyName
			for index, propertyVisitor := range deviceInstance.PropertyVisitors {
				var PropertyNameFound bool
				if propertyName == propertyVisitor.PropertyName {
					PropertyNameFound = true
					property.PVisitor = &propertyVisitor
					break
				}
				if !PropertyNameFound && index == len(deviceInstance.PropertyVisitors)-1 {
					return fmt.Errorf("PropertyVisitor not found")
				}
			}
		}
		devices[deviceInstance.ID] = new(globals.ModbusDev)
		devices[deviceInstance.ID].Instance = deviceInstance
		klog.V(4).Info("Instance: ", deviceInstance.ID, deviceInstance)
	}
	for _, deviceModel := range deviceProfile.DeviceModels {
		dms[deviceModel.Name] = deviceModel
	}

	for _, protocol := range deviceProfile.Protocols {
		protocols[protocol.Name] = protocol
	}
	return nil

}
