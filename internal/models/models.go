package models

type LoadAverageInfo struct {
	Load1Min  float64
	Load5Min  float64
	Load15Min float64
}

type CpuInfo struct {
	System uint64
	User   uint64
	Idle   uint64
}
