// THIS FILE IS AUTO-GENERATED BY CCGEN
// DO NOT MODIFY

package ratetblmonitor

import "encoding/binary"

// <no value>

type RateTblCurrentDataGet struct {
	RateParameterSetId byte

	DatasetRequested uint32
}

func ParseRateTblCurrentDataGet(payload []byte) RateTblCurrentDataGet {
	val := RateTblCurrentDataGet{}

	i := 2

	val.RateParameterSetId = payload[i]
	i++

	val.DatasetRequested = binary.BigEndian.Uint32(payload[i : i+3])
	i += 3

	return val
}