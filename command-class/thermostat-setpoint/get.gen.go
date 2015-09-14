// THIS FILE IS AUTO-GENERATED BY CCGEN
// DO NOT MODIFY

package thermostatsetpoint

import (
	"encoding/gob"
	"errors"
)

func init() {
	gob.Register(Get{})
}

// <no value>
type Get struct {
	Level struct {
		SetpointType byte
	}
}

func (cmd Get) CommandClassID() byte {
	return 0x43
}

func (cmd Get) CommandID() byte {
	return byte(CommandGet)
}

func (cmd *Get) UnmarshalBinary(data []byte) error {
	// According to the docs, we must copy data if we wish to retain it after returning

	payload := make([]byte, len(data))
	copy(payload, data)

	if len(payload) < 2 {
		return errors.New("Payload length underflow")
	}

	i := 2

	if len(payload) <= i {
		return errors.New("slice index out of bounds")
	}

	cmd.Level.SetpointType = (payload[i] & 0x0F)

	i += 1

	return nil
}

func (cmd *Get) MarshalBinary() (payload []byte, err error) {
	payload = make([]byte, 2)
	payload[0] = cmd.CommandClassID()
	payload[1] = cmd.CommandID()

	{
		var val byte

		val |= (cmd.Level.SetpointType) & byte(0x0F)

		payload = append(payload, val)
	}

	return
}