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
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	"github.com/edgexfoundry/device-sdk-go/v3/pkg/interfaces"
	dsModels "github.com/edgexfoundry/device-sdk-go/v3/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/models"
	"github.com/spf13/cast"
)

type Driver struct {
	sdk      interfaces.DeviceServiceSDK
	lc       logger.LoggingClient
	asyncCh  chan<- *dsModels.AsyncValues
	deviceCh chan<- []dsModels.DiscoveredDevice
	generic  map[string]*UartGeneric
}

// Initialize performs protocol-specific initialization for the device
// service.
func (s *Driver) Initialize(sdk interfaces.DeviceServiceSDK) error {
	s.sdk = sdk
	s.lc = sdk.LoggingClient()
	s.asyncCh = sdk.AsyncValuesChannel()
	s.deviceCh = sdk.DiscoveredDeviceChannel()

	s.generic = make(map[string]*UartGeneric)

	return nil
}

// Start runs device service startup tasks after the SDK has been completely
// initialized. This allows device service to safely use DeviceServiceSDK
// interface features in this function call
func (s *Driver) Start() error {
	return nil
}

// HandleReadCommands triggers a protocol Read operation for the specified device.
func (s *Driver) HandleReadCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []dsModels.CommandRequest) (res []*dsModels.CommandValue, err error) {

	const castError = "Failed to parse %s reading: %v"
	const createCommandValueError = "Failed to create %s reading: %v"

	res = make([]*dsModels.CommandValue, len(reqs))

	var deviceLocation string
	var baudRate int

	for i, protocol := range protocols {
		deviceLocation = fmt.Sprintf("%v", protocol["deviceLocation"])
		baudRate, _ = cast.ToIntE(protocol["baudRate"])

		s.lc.Debugf("Driver.HandleReadCommands(): protocol = %v, device location = %v, baud rate = %v", i, deviceLocation, baudRate)
	}

	for i, req := range reqs {
		s.lc.Debugf("Driver.HandleReadCommands(): protocols: %v resource: %v attributes: %v", protocols, req.DeviceResourceName, req.Attributes)

		// Get the value type from device profile
		valueType := req.Type
		s.lc.Debugf("Driver.HandleReadCommands(): value type = %v", valueType)

		key_type_value := fmt.Sprintf("%v", req.Attributes["type"])

		if key_type_value == "generic" {
			key_maxbytes_value, _ := cast.ToIntE(req.Attributes["maxbytes"])
			key_timeout_value, _ := cast.ToIntE(req.Attributes["timeout"])

			// check device is already initialized
			if _, ok := s.generic[deviceLocation]; ok {
				s.lc.Debugf("Driver.HandleReadCommands(): Device %v is already initialized with baud - %v, maxbytes - %v, timeout - %v", s.generic[deviceLocation], baudRate, key_maxbytes_value, key_timeout_value)
			} else {
				// initialize device for the first time
				s.generic[deviceLocation] = NewUartGeneric(deviceLocation, baudRate, key_timeout_value, s.lc)
				s.generic[deviceLocation].rxbuf = nil

				s.lc.Debugf("Driver.HandleReadCommands(): Device %v initialized for the first time with baud - %v, maxbytes - %v, timeout - %v", s.generic[deviceLocation], baudRate, key_maxbytes_value, key_timeout_value)
			}

			if err := s.generic[deviceLocation].GenericUartRead(key_maxbytes_value, s.lc); err != nil {
				return nil, fmt.Errorf("Driver.HandleReadCommands(): Reading UART failed: %v", err)
			}

			rxbuf := hex.EncodeToString(s.generic[deviceLocation].rxbuf)
			s.lc.Debugf("Driver.HandleReadCommands(): Received Data = %s", rxbuf)

			// Pass the received values to higher layers
			// Handle data based on the value type mentioned in device profile
			var cv *dsModels.CommandValue
			switch valueType {
			case common.ValueTypeInt8:
				value, err := strconv.ParseInt(rxbuf, 16, 8)
				if err != nil {
					return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
				}
				cv, err = dsModels.NewCommandValue(req.DeviceResourceName, valueType, int8(value))
				if err != nil {
					return nil, fmt.Errorf(createCommandValueError, req.DeviceResourceName, err)
				}
			case common.ValueTypeInt16:
				value, err := strconv.ParseInt(rxbuf, 16, 16)
				if err != nil {
					return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
				}
				cv, err = dsModels.NewCommandValue(req.DeviceResourceName, valueType, int16(value))
				if err != nil {
					return nil, fmt.Errorf(createCommandValueError, req.DeviceResourceName, err)
				}
			case common.ValueTypeString:
				cv, err = dsModels.NewCommandValue(req.DeviceResourceName, valueType, rxbuf)
				if err != nil {
					return nil, fmt.Errorf(createCommandValueError, req.DeviceResourceName, err)
				}
			default:
				return nil, fmt.Errorf("Driver.HandleReadCommands(): Unsupported value type: %v", valueType)
			}

			s.generic[deviceLocation].rxbuf = nil
			res[i] = cv
			s.lc.Debugf("Driver.HandleReadCommands(): Response = %v", res[i])
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

	var deviceLocation string
	var baudRate int

	for i, protocol := range protocols {
		deviceLocation = fmt.Sprintf("%v", protocol["deviceLocation"])
		baudRate, _ = cast.ToIntE(protocol["baudRate"])

		s.lc.Debugf("Driver.HandleWriteCommands(): protocol = %v, device location = %v, baud rate = %v", i, deviceLocation, baudRate)
	}

	for i, req := range reqs {
		s.lc.Debugf("Driver.HandleWriteCommands(): deviceResourceName = %v", req.DeviceResourceName)
		s.lc.Debugf("Driver.HandleWriteCommands(): protocols: %v, resource: %v, attribute: %v, parameters: %v", protocols, req.DeviceResourceName, req.Attributes, params)

		// Get the value type from device profile
		valueType := req.Type
		s.lc.Debugf("Driver.HandleWriteCommands(): value type = %v", valueType)

		var value interface{}
		var err error

		key_type_value := fmt.Sprintf("%v", req.Attributes["type"])
		if key_type_value == "generic" {
			switch valueType {
			case common.ValueTypeInt8:
				value, err = params[i].Int8Value()
			case common.ValueTypeInt16:
				value, err = params[i].Int16Value()
			case common.ValueTypeString:
				value, err = params[i].StringValue()
			default:
				return fmt.Errorf("Driver.HandleWriteCommands(): Unsupported value type: %v", valueType)
			}

			if err != nil {
				return err
			}

			s.lc.Debugf("Driver.HandleWriteCommands(): %s= %v", valueType, value)

			key_timeout_value, err := cast.ToIntE(req.Attributes["timeout"])
			if err != nil {
				return err
			}

			// initialize the device if it is not initialized already
			if _, ok := s.generic[deviceLocation]; !ok {
				s.generic[deviceLocation] = NewUartGeneric(deviceLocation, baudRate, key_timeout_value, s.lc)
			}

			// Pass the values to the UART end device
			// Handle data based on the value type mentioned in device profile
			var txbuf []byte
			switch valueType {
			case common.ValueTypeInt8:
				txbuf = make([]byte, 1)
				txbuf[0] = byte(value.(int8))
			case common.ValueTypeInt16:
				txbuf = make([]byte, 2)
				// binary.LittleEndian.PutUint16(txbuf, uint16(value.(int16)))
				binary.BigEndian.PutUint16(txbuf, uint16(value.(int16)))
			case common.ValueTypeString:
				// Decode the string to hex format
				txbuf, err = hex.DecodeString(value.(string))
				if err != nil {
					return fmt.Errorf("Driver.HandleWriteCommands(): String decode failed: %v", err)
				}
			default:
				return fmt.Errorf("Driver.HandleWriteCommands(): Unsupported value type: %v", valueType)
			}

			s.lc.Debugf("Driver.HandleWriteCommands(): txbuf = %v", txbuf)

			//Write to UART device
			txlen, err := s.generic[deviceLocation].GenericUartWrite(txbuf, s.lc)

			if err == nil {
				s.lc.Debugf("Driver.HandleWriteCommands(): tx length = %v", txlen)
			}

			return err
		}
	}

	return nil
}

// Discover triggers protocol specific device discovery, asynchronously writes
// the results to the channel which is passed to the implementation via
// ProtocolDriver.Initialize()
func (s *Driver) Discover() error {
	return fmt.Errorf("Discover function is yet to be implemented!")
}

// ValidateDevice triggers device's protocol properties validation, returns error
// if validation failed and the incoming device will not be added into EdgeX
func (s *Driver) ValidateDevice(device models.Device) error {
	protocol, ok := device.Protocols["UART"]
	if !ok {
		return errors.New("Missing 'UART' protocols")
	}

	deviceLocation, ok := protocol["deviceLocation"]
	if !ok {
		return errors.New("Missing 'deviceLocation' information")
	} else if deviceLocation == "" {
		return errors.New("deviceLocation must not empty")
	}

	baudRate, ok := protocol["baudRate"]
	if !ok {
		return errors.New("Missing 'baudRate' information")
	} else if baudRate == "" {
		return errors.New("baudRate must not empty")
	}

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
