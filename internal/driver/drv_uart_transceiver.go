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
	"time"
	"io"
	"github.com/tarm/serial"
)

type UartTransceiver struct {
	rxbuf []byte
	rxlen int
}

func NewUartTransceiver() *UartTransceiver {
	return &UartTransceiver{}
}

func uart_transceiver(dev string, baud int, timeout int, txbuf []byte) ([]byte, int, error) {
	config := &serial.Config{
		Name:        dev,
		Baud:        baud,
		ReadTimeout: time.Millisecond * time.Duration(timeout),
	}

	conn, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
		return nil, 0, err
	}
	defer conn.Close()

	_, err = conn.Write(txbuf)
	if err != nil {
		log.Println(err)
		return nil, 0, err
	}

	var buf []byte
	var buflen int
    re:
    for {
        b := make([]byte, 1024)
        lens, err := conn.Read(b)
        if err != nil {
            switch err {
            case io.EOF:
                if buf != nil {
					buflen =  len(buf)
					err = nil
                }
                break re
            default:
                log.Println(err)
                buf = nil
                break
            }
        }
        buf = append(buf, b[:lens]...)
    }

	return buf, buflen, err
}

