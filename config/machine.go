package config

import (
	"golang.org/x/sys/windows/registry"
)

// GetMachineID retrieves the unique MachineGuid from the Windows registry.
// This is used to derive machine-specific encryption keys.
func GetMachineID() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Cryptography`, registry.QUERY_VALUE)
	if err != nil {
		return "default_shadowlog_machine_id_v1"
	}
	defer k.Close()

	id, _, err := k.GetStringValue("MachineGuid")
	if err != nil {
		return "default_shadowlog_machine_id_v1"
	}
	return id
}
