// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 Canonical Ltd
// Copyright (C) 2018-2019 IOTech Ltd
// Copyright (C) 2021 Jiangxing Intelligence Ltd
// Copyright (C) 2022 HCL Technologies Ltd
//
// SPDX-License-Identifier: Apache-2.0

// Package driver this package provides an UART implementation of
// ProtocolDriver interface.
//
// CONTRIBUTORS              COMPANY
//===============================================================
// 1. Sathya Durai           HCL Technologies
// 2. Sudhamani Bijivemula   HCL Technologies
// 3. Vediyappan Villali     HCL Technologies
// 4. Vijay Annamalaisamy    HCL Technologies
//
//

package driver

import (
    "fmt"
    "encoding/hex"

    dsModels "github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
    "github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
    "github.com/edgexfoundry/go-mod-core-contracts/v2/models"
    "github.com/edgexfoundry/go-mod-core-contracts/v2/common"
    "github.com/spf13/cast"
)

type Driver struct {
    lc      logger.LoggingClient
    asyncCh chan<- *dsModels.AsyncValues
    generic map[string]*UartGeneric
}



// Initialize performs protocol-specific initialization for the device
// service.
func (s *Driver) Initialize(lc logger.LoggingClient, asyncCh chan<- *dsModels.AsyncValues, deviceCh chan<- []dsModels.DiscoveredDevice) (error) {
    s.lc = lc
    s.asyncCh = asyncCh

    s.generic = make(map[string]*UartGeneric)

    return nil
}

// HandleReadCommands triggers a protocol Read operation for the specified device.
func (s *Driver) HandleReadCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []dsModels.CommandRequest) (res []*dsModels.CommandValue, err error) {
	
	res = make([]*dsModels.CommandValue, len(reqs))

    for i, req := range reqs {
        s.lc.Infof(fmt.Sprintf("Driver.HandleReadCommands(): protocols: %v resource: %v attributes: %v", protocols, req.DeviceResourceName, req.Attributes))
		
        key_type_value := fmt.Sprintf("%v", req.Attributes["type"])
		
        if key_type_value == "generic" {
            key_dev_value := fmt.Sprintf("%v", req.Attributes["dev"])

            key_baud_value, _ := cast.ToIntE(fmt.Sprintf("%v", req.Attributes["baud"]))
            key_maxbytes_value, _ := cast.ToIntE(fmt.Sprintf("%v", req.Attributes["maxbytes"]))
            key_timeout_value, _  := cast.ToIntE(fmt.Sprintf("%v", req.Attributes["timeout"]))

            // check device is already initialized
            if _, ok := s.generic[key_dev_value]; ok {
                s.lc.Infof(fmt.Sprintf("Driver.HandleReadCommands(): Device %v is already initialized with baud - %v, maxbytes - %v, timeout - %v", s.generic[key_dev_value], key_baud_value, key_maxbytes_value, key_timeout_value))
            } else {
                // initialize device for the first time
                s.generic[key_dev_value] = NewUartGeneric(key_dev_value, key_baud_value, key_timeout_value)
                s.generic[key_dev_value].rxbuf = nil

		s.lc.Infof(fmt.Sprintf("Driver.HandleReadCommands(): Device %v initialized for the first time with baud - %v, maxbytes - %v, timeout - %v", s.generic[key_dev_value], key_baud_value, key_maxbytes_value, key_timeout_value))
            }
            
            // Call the UART read function
            s.generic[key_dev_value].GenericUartRead(key_maxbytes_value)

            rxbuf := hex.EncodeToString(s.generic[key_dev_value].rxbuf)
            s.lc.Infof(fmt.Sprintf("Driver.HandleReadCommands(): Received Data =  %s", rxbuf))
			
            // Pass the received values to higher layers
            cv, _ := dsModels.NewCommandValue(req.DeviceResourceName, common.ValueTypeString, rxbuf)
            s.generic[key_dev_value].rxbuf = nil
            res[i] = cv
            s.lc.Infof(fmt.Sprintf("Driver.HandleReadCommands(): Response = %v", res[i]))
        }
    }

    return res, nil
}

// HandleWriteCommands passes a slice of CommandRequest struct each representing
// a ResourceOperation for a specific device resource.
// Since the commands are actuation commands, params provide parameters for the individual
// command.
func (s *Driver) HandleWriteCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []dsModels.CommandRequest,
    params []*dsModels.CommandValue) error {

    for i, req := range reqs {
        s.lc.Infof(fmt.Sprintf("Driver.HandleWriteCommands(): deviceResourceName = %v", req.DeviceResourceName))
		s.lc.Infof(fmt.Sprintf("Driver.HandleWriteCommands(): protocols: %v, resource: %v, attribute: %v, parameters: %v", protocols, req.DeviceResourceName, req.Attributes, params))
		
        key_type_value := fmt.Sprintf("%v", req.Attributes["type"])
		
        if key_type_value == "generic" {
            if value, err := params[i].StringValue(); err == nil {
                key_dev_value := fmt.Sprintf("%v", req.Attributes["dev"])
                key_baud_value, _ := cast.ToIntE(fmt.Sprintf("%v", req.Attributes["baud"]))
                key_timeout_value, _  := cast.ToIntE(fmt.Sprintf("%v", req.Attributes["timeout"]))

                // initialize the device if it is not initialized already
                if _, ok := s.generic[key_dev_value]; ok == false {
                    s.generic[key_dev_value] = NewUartGeneric(key_dev_value, key_baud_value, key_timeout_value)
                }
			    
                // decode the string in hex format
                txbuf, err := hex.DecodeString(value)
                if err != nil {
                    return err
                }
			    
                //Write to UART device
                txlen, err :=s.generic[key_dev_value].GenericUartWrite(txbuf)
			    
                if err == nil {
                    s.lc.Infof(fmt.Sprintf("Driver.HandleWriteCommands(): tx length = %v", txlen))
                }

                return err
            } else {
                return err
            }
        }
    }
    s.lc.Infof(fmt.Sprintf("Driver.HandleWriteCommands(): return nil"))
    return nil
}

// Stop the protocol-specific DS code to shutdown gracefully, or
// if the force parameter is 'true', immediately. The driver is responsible
// for closing any in-use channels, including the channel used to send async
// readings (if supported).
func (s *Driver) Stop(force bool) error {
    // Then Logging Client might not be initialized
    if s.lc != nil {
        s.lc.Debugf(fmt.Sprintf("Driver.Stop called: force=%v", force))
    }
    return nil
}

// AddDevice is a callback function that is invoked
// when a new Device associated with this Device Service is added
func (s *Driver) AddDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
    s.lc.Debugf(fmt.Sprintf("a new Device is added: %s", deviceName))
    return nil
}

// UpdateDevice is a callback function that is invoked
// when a Device associated with this Device Service is updated
func (s *Driver) UpdateDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
    s.lc.Debugf(fmt.Sprintf("Device %s is updated", deviceName))
    return nil
}

// RemoveDevice is a callback function that is invoked
// when a Device associated with this Device Service is removed
func (s *Driver) RemoveDevice(deviceName string, protocols map[string]models.ProtocolProperties) error {
    s.lc.Debugf(fmt.Sprintf("Device %s is removed", deviceName))
    return nil
}

