package monitoring

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// MemorySnapshot holds values from /proc/meminfo (kB unless noted).
type MemorySnapshot struct {
	MemTotalKB     uint64
	MemFreeKB      uint64
	MemAvailableKB uint64
	BuffersKB      uint64
	CachedKB       uint64
	SReclaimableKB uint64
	ShmemKB        uint64
	SwapTotalKB    uint64
	SwapFreeKB     uint64
}

// ReadMemory parses /proc/meminfo.
func ReadMemory() (*MemorySnapshot, error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, fmt.Errorf("monitoring/memory: %w", err)
	}
	defer f.Close()

	s := &MemorySnapshot{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		i := strings.IndexByte(line, ':')
		if i < 0 {
			continue
		}
		key := strings.TrimSpace(line[:i])
		rest := strings.TrimSpace(line[i+1:])
		fields := strings.Fields(rest)
		if len(fields) < 1 {
			continue
		}
		val, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			continue
		}
		switch key {
		case "MemTotal":
			s.MemTotalKB = val
		case "MemFree":
			s.MemFreeKB = val
		case "MemAvailable":
			s.MemAvailableKB = val
		case "Buffers":
			s.BuffersKB = val
		case "Cached":
			s.CachedKB = val
		case "SReclaimable":
			s.SReclaimableKB = val
		case "Shmem":
			s.ShmemKB = val
		case "SwapTotal":
			s.SwapTotalKB = val
		case "SwapFree":
			s.SwapFreeKB = val
		}
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("monitoring/memory: %w", err)
	}
	if s.MemTotalKB == 0 {
		return nil, fmt.Errorf("monitoring/memory: MemTotal ausente ou inválido")
	}
	return s, nil
}

// UsedApproxKB is total minus free minus buffers minus cache (similar to `free` "used").
func (s *MemorySnapshot) UsedApproxKB() uint64 {
	if s == nil {
		return 0
	}
	cache := s.CachedKB + s.BuffersKB + s.SReclaimableKB
	if s.MemTotalKB <= s.MemFreeKB+cache {
		return 0
	}
	return s.MemTotalKB - s.MemFreeKB - cache
}
