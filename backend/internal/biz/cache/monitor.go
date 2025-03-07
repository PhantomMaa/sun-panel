package cache

import (
	"sun-panel/internal/infra/kvcache"
	"sun-panel/internal/util/monitor"
	"time"
)

type Monitor struct {
	Cache kvcache.Cacher[interface{}]
}

const cacheSecond = 3

func (a *Monitor) GetCpuState() (any, error) {
	if v, ok := a.Cache.Get(monitor.SystemmonitorCpuInfo); ok {
		return v, nil
	}

	cpuInfo, err := monitor.GetCPUInfo()
	if err != nil {
		return nil, err
	}

	a.Cache.Set(monitor.SystemmonitorCpuInfo, cpuInfo, cacheSecond*time.Second)
	return cpuInfo, nil
}

func (a *Monitor) GetMemonyState() (any, error) {
	if v, ok := a.Cache.Get(monitor.SystemmonitorMemoryInfo); ok {
		return v, nil
	}

	memoryInfo, err := monitor.GetMemoryInfo()
	if err != nil {
		return nil, err
	}

	a.Cache.Set(monitor.SystemmonitorMemoryInfo, memoryInfo, cacheSecond*time.Second)
	return memoryInfo, nil
}

func (a *Monitor) GetDiskStateByPath(path string) (any, error) {
	disk := monitor.SystemmonitorDiskInfo + path
	if v, ok := a.Cache.Get(disk); ok {
		return v, nil
	}

	diskState, err := monitor.GetDiskInfoByPath(path)
	if err != nil {
		return nil, err
	}

	a.Cache.Set(disk, diskState, cacheSecond*time.Second)
	return diskState, nil
}
