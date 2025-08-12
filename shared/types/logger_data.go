package types

import "time"

type LoggerData struct {
	// SchemaVersion defines the version of the data schema.
	SchemaVersion int `json:"schema_version"`

	// DeviceID uniquely identifies the logger device.
	DeviceID string `json:"device_id"`

	// Timestamp is the time when the data was recorded (UTC).
	Timestamp time.Time `json:"timestamp"`

	// FirmwareVersion is the firmware version running on the device.
	FirmwareVersion string `json:"firmware_version"`

	// MeterReading contains all the electrical measurement data.
	MeterReading MeterReading `json:"meter_reading"`

	// Events is a list of important events recorded by the logger.
	Events []Event `json:"events"`

	// Location is the physical geographic position of the logger.
	Location Location `json:"location"`

	// SignalStrengthDBM represents the wireless signal strength in dBm.
	SignalStrengthDBM int `json:"signal_strength_dbm"`

	// BatteryLevelPercent shows the remaining battery percentage (0-100).
	BatteryLevelPercent uint8 `json:"battery_level_percent"`

	// Status indicates the current device status (e.g., "OK", "Warning").
	Status string `json:"status"`
}

type MeterReading struct {
	// Voltage contains voltage measurements for three phases.
	Voltage ThreePhase `json:"voltage"`

	// Current contains current measurements for three phases.
	Current ThreePhase `json:"current"`

	// Frequency in Hertz (Hz) of the electrical supply.
	Frequency float32 `json:"frequency"`

	// PowerFactor is the ratio of active power to apparent power.
	PowerFactor float32 `json:"power_factor"`

	// ActivePowerKW is the real power in kilowatts.
	ActivePowerKW float32 `json:"active_power_kw"`

	// ReactivePowerKVAR is the reactive power in kilovolt-amperes reactive.
	ReactivePowerKVAR float32 `json:"reactive_power_kvar"`

	// ApparentPowerKVA is the apparent power in kilovolt-amperes.
	ApparentPowerKVA float32 `json:"apparent_power_kva"`

	// EnergyConsumedKWH is the cumulative consumed energy in kilowatt-hours.
	EnergyConsumedKWH float64 `json:"energy_consumed_kwh"`

	// EnergyExportedKWH is the cumulative exported energy in kilowatt-hours.
	EnergyExportedKWH float32 `json:"energy_exported_kwh"`
}

type ThreePhase struct {
	// PhaseA measurement for phase A.
	PhaseA float32 `json:"phase_a"`

	// PhaseB measurement for phase B.
	PhaseB float32 `json:"phase_b"`

	// PhaseC measurement for phase C.
	PhaseC float32 `json:"phase_c"`
}

type Event struct {
	// Type is the event type (e.g., "OVERVOLTAGE", "POWER_OUTAGE").
	Type string `json:"type"`

	// Value is the event value, if applicable (e.g., voltage level).
	Value float32 `json:"value"`

	// Timestamp is when the event occurred.
	Timestamp time.Time `json:"timestamp"`
}

type Location struct {
	// Latitude of the logger's physical location.
	Latitude float64 `json:"latitude"`

	// Longitude of the logger's physical location.
	Longitude float64 `json:"longitude"`
}
