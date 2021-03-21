// Copyright 2021 James Grant

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 		http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package udmx

import (
	"fmt"

	"github.com/google/gousb"
)

type Udmx struct {
	c *gousb.Context
	d *gousb.Device
}

const (
	manufacturerString = "www.anyma.ch"
	productString      = "uDMX"
)

func NewUdmxForId(vid, pid gousb.ID) (*Udmx, error) {
	c := gousb.NewContext()

	devs, err := c.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return desc.Vendor == vid || desc.Product == pid
	})
	if err != nil {
		for _, d := range devs {
			d.Close()
		}
		_ = c.Close()
		return nil, err
	}

	for i, dev := range devs {
		manufacturer, err := dev.Manufacturer()
		if err != nil {
			return nil, err
		}
		product, err := dev.Product()
		if err != nil {
			return nil, err
		}
		if manufacturer == manufacturerString && product == productString {
			for j, d := range devs {
				if i != j {
					err := d.Close()
					if err != nil {
						return nil, fmt.Errorf("problem closing unwanted device : %v", err)
					}
				}
			}
			return &Udmx{
				c: c,
				d: dev,
			}, nil
		}
	}

	if err := c.Close(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("unable to find device")
}

func NewUdmx() (*Udmx, error) {
	return NewUdmxForId(0x16C0, 0x05DC)
}

const (
	cmd_SetSingleChannel = 1
	cmd_SetChannelRange  = 2
)

func (dmx *Udmx) SetSingleChannel(channel, value uint8) error {
	_, err := dmx.d.Control(
		gousb.ControlVendor|gousb.ControlDevice|gousb.ControlOut,
		cmd_SetSingleChannel,
		uint16(value),
		uint16(channel),
		[]byte{})
	return err
}

func (dmx *Udmx) SetChannelRange(channel uint8, values []uint8) error {
	_, err := dmx.d.Control(
		gousb.ControlVendor|gousb.ControlDevice|gousb.ControlOut,
		cmd_SetChannelRange,
		uint16(len(values)),
		uint16(channel),
		values)
	return err
}

func (dmx *Udmx) Close() error {
	if err := dmx.d.Close(); err != nil {
		return err
	}

	if err := dmx.c.Close(); err != nil {
		return err
	}

	return nil
}
