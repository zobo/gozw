// THIS FILE IS AUTO-GENERATED BY CCGEN
// DO NOT MODIFY

package associationcommandconfiguration

// <no value>

type CommandConfigurationSet struct {
	GroupingIdentifier byte

	NodeId byte

	CommandLength byte

	CommandClassIdentifier byte

	CommandIdentifier byte

	CommandByte []byte
}

func ParseCommandConfigurationSet(payload []byte) CommandConfigurationSet {
	val := CommandConfigurationSet{}

	i := 2

	val.GroupingIdentifier = payload[i]
	i++

	val.NodeId = payload[i]
	i++

	val.CommandLength = payload[i]
	i++

	val.CommandClassIdentifier = payload[i]
	i++

	val.CommandIdentifier = payload[i]
	i++

	val.CommandByte = payload[i:]

	return val
}