// -*- Mode: Go; indent-tabs-mode: t -*-
//
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
	"io"
	"log"
	"time"

	"github.com/edgexfoundry/go-mod-core-contracts/v3/clients/logger"
	"github.com/tarm/serial"
)

// UartGeneric is a structure
// config is a pointer to Config structure in github.com/tarm/serial package. Config contains the information needed to open a serial port.
// conn is a pointer to Port structure in github.com/tarm/serial package.
// rxbuf is byte array
// enable is set to true when new device is added
// portStatus is set to true when the autoevent is being executed
type UartGeneric struct {
	config     *serial.Config
	conn       *serial.Port
	rxbuf      []byte
	enable     bool
	portStatus bool
}

// NewUartGeneric is a function
// takes device name and baud rate as arguments and
// returns a pointer to UartGeneric structure
func NewUartGeneric(dev string, baud int, timeout int, lc logger.LoggingClient) *UartGeneric {
	config := &serial.Config{
		Name:        dev,
		Baud:        baud,
		ReadTimeout: time.Second * time.Duration(timeout),
	}
	var err error

	conn, err := serial.OpenPort(config)
	if err != nil {
		lc.Errorf("NewUartGeneric(): Open serial %s failed: %v", config.Name, err)
	}
	return &UartGeneric{config: config, conn: conn, enable: true, portStatus: false}
}

// GenericUartRead method
func (dev *UartGeneric) GenericUartRead(maxbytes int, lc logger.LoggingClient) error {
	var buf []byte

	// GO uart packages read a maximum of 16 bytes in a single shot
	// so using maxbytes & need to repeatedly read until max bytes
	// are read or EOF is reached
	readCount := (maxbytes / 16) + 1

	lc.Debugf("GenericUartRead(): readCount = %v", readCount)

	// We don't want next auto-event to interrupt when the current one is
	// still executing
	if dev.portStatus {
		lc.Errorf("GenericUartRead(): Exit - Device busy... Read request dropped for %s", dev.config.Name)
		return nil
	}

	dev.portStatus = true

	// Allow up to 128 to be read but also note how many actually were read.
	b := make([]byte, 128)

	for i := 1; i <= readCount; i++ {
		lens, err := dev.conn.Read(b)

		if err != nil {
			if err == io.EOF {
				lc.Debugf("GenericUartRead(): %v - Finished reading!", err)
				break
			}
			lc.Errorf("GenericUartRead(): Exit - Error = %v", err)

			dev.portStatus = false

			dev.conn.Flush()

			return err
		}

		lc.Debugf("GenericUartRead(): Number of bytes read = %v, buf = %s", lens, b)

		// Copy the content of buf to device rxbuf
		buf = append(buf, b[:lens]...)

		dev.rxbuf = append(dev.rxbuf, buf[:]...)
		lc.Debugf("GenericUartRead(): dev.rxbuf = %s", dev.rxbuf)
		buf = nil
	}

	dev.portStatus = false
	dev.conn.Flush()
	lc.Debugf("GenericUartRead(): Exit - Success")

	return nil
}

// GenericUartWrite method
func (dev *UartGeneric) GenericUartWrite(txbuf []byte, lc logger.LoggingClient) (int, error) {

	dev.conn.Flush()

	length, err := dev.conn.Write(txbuf)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	lc.Debugf("GenericUartWrite(): Number of bytes transmitted = %d\n", length)

	return length, err
}
