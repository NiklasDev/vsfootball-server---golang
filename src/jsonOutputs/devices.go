package jsonOutputs

type DeviceOutput struct {
	Success        string
	Message        string
	Status         int
	Iosdevices     []string
	Androiddevices []string
}
