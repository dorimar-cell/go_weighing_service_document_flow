package dto

import "time"

// StationDTO — весовая и настройки её терминала.
type StationDTO struct {
	TerminalID       int64      `json:"terminalId"`
	WeightRoomID     int64      `json:"weightRoomId"`
	WeightRoomName   string     `json:"weightRoomName"`
	PlatformNumber   int64      `json:"platformNumber"`
	TerminalModel    string     `json:"terminalModel"`
	InterfaceType    string     `json:"interfaceType"`
	ComName          string     `json:"comName"`
	Baudrate         string     `json:"baudrate"`
	Parity           string     `json:"parity"`
	Databits         string     `json:"databits"`
	Stopbits         string     `json:"stopbits"`
	MultiplierWeight float64    `json:"multiplierWeight"`
	IsModbus         bool       `json:"isModbus"`
	IP               string     `json:"ip"`
	Port             string     `json:"port"`
	UseSocket        bool       `json:"useSocket"`
	SocketIP         string     `json:"socketIp"`
	SocketPort       string     `json:"socketPort"`
	Enable           bool       `json:"enable"`
	TerminalInterval string     `json:"terminalInterval"`
	UpdatedAt        *time.Time `json:"updatedAt"`
}

// SaveStationRequest — правка настроек терминала.
type SaveStationRequest struct {
	PlatformNumber   int64   `json:"platformNumber"`
	TerminalModel    string  `json:"terminalModel"`
	InterfaceType    string  `json:"interfaceType"`
	ComName          string  `json:"comName"`
	Baudrate         string  `json:"baudrate"`
	Parity           string  `json:"parity"`
	Databits         string  `json:"databits"`
	Stopbits         string  `json:"stopbits"`
	MultiplierWeight float64 `json:"multiplierWeight"`
	IsModbus         bool    `json:"isModbus"`
	IP               string  `json:"ip"`
	Port             string  `json:"port"`
	UseSocket        bool    `json:"useSocket"`
	SocketIP         string  `json:"socketIp"`
	SocketPort       string  `json:"socketPort"`
	Enable           bool    `json:"enable"`
	TerminalInterval string  `json:"terminalInterval"`
}
