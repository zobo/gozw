package zwave

const (
	ZwNodeList              = 0x02
	ZwAddNodeToNetwork      = 0x4a
	ZwRemoveNodeFromNetwork = 0x4b
)

const (
	AddNodeAny                  = 1
	AddNodeController           = 2
	AddNodeSlave                = 3
	AddNodeExisting             = 4
	AddNodeStop                 = 5
	AddNodeStopFailed           = 6
	AddNodeStatusSecurityFailed = 9
)

const (
	AddNodeOptionNormalPower = 0x80
	AddNodeOptionNetworkWide = 0x40
)

const (
	RemoveNodeAny        = AddNodeAny
	RemoveNodeController = AddNodeController
	RemoveNodeSlave      = AddNodeSlave
	RemoveNodeStop       = AddNodeStop
)

const (
	RemoveNodeOptionNormalPower = AddNodeOptionNormalPower
	RemoveNodeOptionNetworkWide = AddNodeOptionNetworkWide
)

func GetNodeList() []byte {
	return []byte{ZwNodeList}
}

func EnterInclusionMode() []byte {
	return []byte{
		ZwAddNodeToNetwork,
		AddNodeAny | AddNodeOptionNormalPower | AddNodeOptionNetworkWide,
		0x01,
	}
}

func ExitInclusionMode() []byte {
	return []byte{
		ZwRemoveNodeFromNetwork,
		AddNodeStop,
		0x01,
	}
}

func EnterExclusionMode() []byte {
	return []byte{
		ZwRemoveNodeFromNetwork,
		RemoveNodeAny | RemoveNodeOptionNormalPower | RemoveNodeOptionNetworkWide,
		0x01,
	}
}

func ExitExclusionMode() []byte {
	return []byte{
		ZwRemoveNodeFromNetwork,
		RemoveNodeAny | RemoveNodeOptionNormalPower | RemoveNodeOptionNetworkWide,
		0x01,
	}
}
