// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 Canonical Ltd
// Copyright (C) 2018-2019 IOTech Ltd
// Copyright (C) 2021 Jiangxing Intelligence Ltd
//
// SPDX-License-Identifier: Apache-2.0

// Package driver this package provides an UART implementation of
// ProtocolDriver interface.
//
package driver

import (
	"fmt"
	"time"
	"strconv"
	"encoding/hex"
	"errors"
	"sync"

	dsModels "github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
)

type Driver struct {
	lc      logger.LoggingClient
	asyncCh chan<- *dsModels.AsyncValues
	monitor map[string]*UartMonitor
	transceiver map[string]*UartTransceiver
	locker         sync.Mutex
}



// Initialize performs protocol-specific initialization for the device
// service.
func (s *Driver) Initialize(lc logger.LoggingClient, asyncCh chan<- *dsModels.AsyncValues, deviceCh chan<- []dsModels.DiscoveredDevice) error {
	s.lc = lc
	s.asyncCh = asyncCh

	s.monitor = make(map[string]*UartMonitor)
	s.transceiver = make(map[string]*UartTransceiver)

	return nil
}

// HandleReadCommands triggers a protocol Read operation for the specified device.
func (s *Driver) HandleReadCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []dsModels.CommandRequest) (res []*dsModels.CommandValue, err error) {

	s.locker.Lock()
	defer s.locker.Unlock()
	s.lc.Infof("protocols: %v resource: %v attributes: %v", protocols, reqs[0].DeviceResourceName, reqs[0].Attributes)

	if len(reqs) == 1 {
		res = make([]*dsModels.CommandValue, 1)
		key_type_value := fmt.Sprintf("%v", reqs[0].Attributes["tpye"])
		if key_type_value == "transceiver" {
			key_dev_value := fmt.Sprintf("%v", reqs[0].Attributes["dev"])
			if _, ok := s.transceiver[key_dev_value]; ok == false {
				s.transceiver[key_dev_value] = NewUartTransceiver()
				s.transceiver[key_dev_value].rxbuf = nil
				s.transceiver[key_dev_value].rxlen = 0
			}
			rxbuf := hex.EncodeToString(s.transceiver[key_dev_value].rxbuf)
			cv, _ := dsModels.NewCommandValue(reqs[0].DeviceResourceName, common.ValueTypeString, rxbuf)
			s.transceiver[key_dev_value].rxbuf = nil
			s.transceiver[key_dev_value].rxlen = 0
			res[0] = cv
		} else if key_type_value == "monitor" {
			key_dev_value := fmt.Sprintf("%v", reqs[0].Attributes["dev"])
			key_baud_value, _ := strconv.Atoi(fmt.Sprintf("%v", reqs[0].Attributes["baud"]))
			if _, ok := s.monitor[key_dev_value]; ok {
				if s.monitor[key_dev_value].rxstatus == false {
					go s.monitor[key_dev_value].Start_Listen()
					time.Sleep(100 * time.Millisecond)
				}
			} else {
				s.monitor[key_dev_value] = NewUartMonitor(key_dev_value, key_baud_value)
				go s.monitor[key_dev_value].Start_Listen()
				time.Sleep(100 * time.Millisecond)
			}

			if s.monitor[key_dev_value].rxstatus {
				rxbuf := hex.EncodeToString(s.monitor[key_dev_value].rxbuf)
				cv, _ := dsModels.NewCommandValue(reqs[0].DeviceResourceName, common.ValueTypeString, rxbuf)
				s.monitor[key_dev_value].rxbuf = nil
				res[0] = cv
			} else {
				return nil, errors.New("[error]: Open serial fail")
			}
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
	s.locker.Lock()
	defer s.locker.Unlock()
	s.lc.Infof("Driver.HandleWriteCommands: protocols: %v, resource: %v, attribute: %v, parameters: %v", protocols, reqs[0].DeviceResourceName, reqs[0].Attributes, params)

	for i, r := range reqs {
		s.lc.Infof(r.DeviceResourceName)
		key_type_value := fmt.Sprintf("%v", reqs[0].Attributes["tpye"])
		if key_type_value == "transceiver" {
			if value, err := params[i].StringValue(); err == nil {
				key_dev_value := fmt.Sprintf("%v", reqs[0].Attributes["dev"])
				key_baud_value, _ := strconv.Atoi(fmt.Sprintf("%v", reqs[0].Attributes["baud"]))
				key_timeout_value, _ := strconv.Atoi(fmt.Sprintf("%v", reqs[0].Attributes["timeout"]))
				if _, ok := s.transceiver[key_dev_value]; ok == false {
					s.transceiver[key_dev_value] = NewUartTransceiver()
					s.transceiver[key_dev_value].rxbuf = nil
					s.transceiver[key_dev_value].rxlen = 0
				}
				txbuf, err := hex.DecodeString(value)
				if err != nil {
					return err
				}
				s.transceiver[key_dev_value].rxbuf, s.transceiver[key_dev_value].rxlen, err = uart_transceiver(key_dev_value, key_baud_value, key_timeout_value, txbuf)
				if err == nil {
					if s.transceiver[key_dev_value].rxlen == 0 {
						return errors.New("[log]: No response")
					}
				}
				return err
			} else {
				return err
			}
		}

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
	for k, _ := range s.monitor  {
		if s.monitor[k].rxstatus == true {
			s.monitor[k].Stop_Listen()
		}
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
	for k, _ := range s.monitor  {
		if s.monitor[k].rxstatus == true {
			s.monitor[k].Stop_Listen()
		}
	}
	return nil
}

// RemoveDevice is a callback function that is invoked
// when a Device associated with this Device Service is removed
func (s *Driver) RemoveDevice(deviceName string, protocols map[string]models.ProtocolProperties) error {
	s.lc.Debugf(fmt.Sprintf("Device %s is removed", deviceName))
	for k, _ := range s.monitor  {
		if s.monitor[k].rxstatus == true {
			s.monitor[k].Stop_Listen()
		}
	}
	return nil
}
