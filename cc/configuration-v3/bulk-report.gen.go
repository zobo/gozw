// THIS FILE IS AUTO-GENERATED BY ZWGEN
// DO NOT MODIFY

package configurationv3

import (
	"encoding/binary"
	"encoding/gob"
	"errors"

	"github.com/gozwave/gozw/cc"
)

const CommandBulkReport cc.CommandID = 0x09

func init() {
	gob.Register(BulkReport{})
	cc.Register(cc.CommandIdentifier{
		CommandClass: cc.CommandClassID(0x70),
		Command:      cc.CommandID(0x09),
		Version:      3,
	}, NewBulkReport)
}

func NewBulkReport() cc.Command {
	return &BulkReport{}
}

// <no value>
type BulkReport struct {
	ParameterOffset uint16

	NumberOfParameters byte

	ReportsToFollow byte

	Properties1 struct {
		Size byte

		Handshake bool

		Default bool
	}

	Vg []BulkReportVg
}

type BulkReportVg struct {
	Parameter []byte
}

func (cmd BulkReport) CommandClassID() cc.CommandClassID {
	return 0x70
}

func (cmd BulkReport) CommandID() cc.CommandID {
	return CommandBulkReport
}

func (cmd BulkReport) CommandIDString() string {
	return "CONFIGURATION_BULK_REPORT"
}

func (cmd *BulkReport) UnmarshalBinary(data []byte) error {
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

	cmd.ParameterOffset = binary.BigEndian.Uint16(payload[i : i+2])
	i += 2

	if len(payload) <= i {
		return errors.New("slice index out of bounds")
	}

	cmd.NumberOfParameters = payload[i]
	i++

	if len(payload) <= i {
		return errors.New("slice index out of bounds")
	}

	cmd.ReportsToFollow = payload[i]
	i++

	if len(payload) <= i {
		return errors.New("slice index out of bounds")
	}

	cmd.Properties1.Size = (payload[i] & 0x07)

	cmd.Properties1.Handshake = payload[i]&0x40 == 0x40

	cmd.Properties1.Default = payload[i]&0x80 == 0x80

	i += 1

	for i < len(payload) {

		if len(payload) <= i {
			return errors.New("slice index out of bounds")
		}

		{
			length := (payload[131+2] >> 0) & 0x07
			cmd.Parameter = payload[i : i+int(length)]
			i += int(length)
		}

		vg := BulkReportVg{

			Parameter: parameter,
		}
		cmd.Vg = append(cmd.Vg, vg)
	}

	return nil
}

func (cmd *BulkReport) MarshalBinary() (payload []byte, err error) {
	payload = make([]byte, 2)
	payload[0] = byte(cmd.CommandClassID())
	payload[1] = byte(cmd.CommandID())

	{
		buf := make([]byte, 2)
		binary.BigEndian.PutUint16(buf, cmd.ParameterOffset)
		payload = append(payload, buf...)
	}

	payload = append(payload, cmd.NumberOfParameters)

	payload = append(payload, cmd.ReportsToFollow)

	{
		var val byte

		val |= (cmd.Properties1.Size) & byte(0x07)

		if cmd.Properties1.Handshake {
			val |= byte(0x40) // flip bits on
		} else {
			val &= ^byte(0x40) // flip bits off
		}

		if cmd.Properties1.Default {
			val |= byte(0x80) // flip bits on
		} else {
			val &= ^byte(0x80) // flip bits off
		}

		payload = append(payload, val)
	}

	for _, vg := range cmd.Vg {

		if vg.Parameter != nil && len(vg.Parameter) > 0 {
			payload = append(payload, vg.Parameter...)
		}

	}

	return
}
