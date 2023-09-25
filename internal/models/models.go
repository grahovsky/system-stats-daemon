package models

type LoadAverageInfo struct {
	Load1Min  float64
	Load5Min  float64
	Load15Min float64
}

type CpuInfo struct {
	Idle  uint64
	Total uint64
}
