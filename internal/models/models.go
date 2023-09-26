package models

type LoadAverageInfo struct {
	Load1Min  float64
	Load5Min  float64
	Load15Min float64
}

type CpuInfo struct {
	User   float64
	System float64
	Idle   float64
}

type DiskInfo struct {
	Kbt float64
	Tps float64
}
