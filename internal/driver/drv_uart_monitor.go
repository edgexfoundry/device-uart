// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2021 Jiangxing Intelligence Ltd
//
// SPDX-License-Identifier: Apache-2.0

// Package driver this package provides an UART implementation of
// ProtocolDriver interface.
//
package driver

import (
	"log"
	"io"
	"os"
	"time"
	"github.com/tarm/serial"
)

type UartMonitor struct {
	config *serial.Config
	conn *serial.Port
	rxbuf []byte
	rxlen int
	enable bool
	rxstatus bool
}

func NewUartMonitor(dev string, baud int) *UartMonitor {
	config := &serial.Config{
		Name:        dev,
		Baud:        baud,
		ReadTimeout: time.Millisecond * 10,
	}

	return &UartMonitor{config: config, enable: true, rxstatus: false}
}

func exist(path string) bool {
    _, err := os.Stat(path)
    if err != nil {
        if os.IsExist(err) {
            return true
        }
        return false
    }
    return true
}

func (dev *UartMonitor) Start_Listen() error {
	var err error
	dev.conn, err = serial.OpenPort(dev.config)
	if err != nil {
		log.Printf("Open serial %s fail\n", dev.config.Name)
		dev.rxstatus = false
		return err
	}
	defer dev.conn.Close()
	dev.rxstatus = true

	var buf []byte
	re:
	for {
		b := make([]byte, 1024)
		lens, err := dev.conn.Read(b)
		if err != nil {
			switch err {
			case io.EOF:
				if dev.enable == false {
					break re
				}
				if exist(dev.config.Name)  == false {
					break re
				}
				if buf != nil {
					dev.rxbuf = nil
					dev.rxbuf = append(dev.rxbuf, buf[:len(buf)]...)
					buf = nil
					break
				}
			default:
				log.Println(err)
				buf = nil
				break
			}
		}
		buf = append(buf, b[:lens]...)
	}
	dev.rxstatus = false
	log.Printf("Stop listen serial %s\n", dev.config.Name)
	return nil
}

func (dev *UartMonitor) Stop_Listen() error {
	dev.enable = false
	return nil
}

