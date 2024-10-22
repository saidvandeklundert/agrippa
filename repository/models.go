package repository

type Database int

// Use GetDatabaseByType to select the database in a type-safe way
const (
	APPL_DB         Database = iota
	ASIC_DB         Database = iota
	COUNTERS_DB     Database = iota
	LOGLEVEL_DB     Database = iota
	CONFIG_DB       Database = iota
	PFC_WD_DB       Database = iota
	STATE_DB        Database = iota
	INTERNAL_AMAZON Database = iota
)

type DeviceMetadata struct {
	BufferModel string `redis:"buffer_model"`
	Hwsku       string `redis:"hwsku"`
	SwitchType  string `redis:"switch_type"`
	Mac         string `redis:"mac"`
	Hostname    string `redis:"hostname"`
}

/*
The model from APPL_DB that represents the output of LLDP_LOC_CHASSIS
*/
type LldpLocalChassis struct {
	Id                 string `redis:"lldp_loc_chassis_id"`
	SystemName         string `redis:"lldp_loc_sys_name"`
	IdSubtype          string `redis:"lldp_loc_chassis_id_subtype"`
	SystemDescription  string `redis:"lldp_loc_sys_desc"`
	ManagementAddress  string `redis:"lldp_loc_man_addr"`
	SystemCapEnabled   string `redis:"lldp_loc_sys_cap_enabled"`
	SystemCapSupported string `redis:"lldp_loc_sys_cap_supported"`
}
