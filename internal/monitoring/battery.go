package monitoring

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// BatterySnapshot holds sysfs power_supply data for the main battery pack.
type BatterySnapshot struct {
	SupplyName string

	Status         string
	Capacity       int    // 0–100, -1 if unknown
	CapacityLevel  string // e.g. Normal, High
	EnergyNowUWh   int64  // micro watt-hours; 0 if absent
	EnergyFullUWh  int64
	EnergyDesignUWh int64
	PowerNowUW     int64 // microwatts
	VoltageNowUV   int64
	CycleCount     int
	Technology     string
	Manufacturer   string
	ModelName      string

	ACOnline *bool // nil if no AC supply node
	ACName   string
}

// ReadBattery finds the first power_supply of type Battery and reads common attributes.
func ReadBattery() (*BatterySnapshot, error) {
	const base = "/sys/class/power_supply"
	ents, err := os.ReadDir(base)
	if err != nil {
		return nil, fmt.Errorf("monitoring/battery: %w", err)
	}

	var batName string
	for _, e := range ents {
		name := e.Name()
		t := readTrimFile(filepath.Join(base, name, "type"))
		if strings.EqualFold(t, "Battery") {
			batName = name
			break
		}
	}
	if batName == "" {
		return nil, fmt.Errorf("monitoring/battery: nenhuma bateria encontrada em %s", base)
	}

	dir := filepath.Join(base, batName)
	s := &BatterySnapshot{
		SupplyName:    batName,
		Status:        readTrimFile(filepath.Join(dir, "status")),
		CapacityLevel: readTrimFile(filepath.Join(dir, "capacity_level")),
		Technology:    readTrimFile(filepath.Join(dir, "technology")),
		Manufacturer:  readTrimFile(filepath.Join(dir, "manufacturer")),
		ModelName:     readTrimFile(filepath.Join(dir, "model_name")),
	}
	if c, ok := readIntFile(filepath.Join(dir, "capacity")); ok {
		s.Capacity = c
	} else {
		s.Capacity = -1
	}
	s.EnergyNowUWh, _ = readInt64File(filepath.Join(dir, "energy_now"))
	s.EnergyFullUWh, _ = readInt64File(filepath.Join(dir, "energy_full"))
	s.EnergyDesignUWh, _ = readInt64File(filepath.Join(dir, "energy_full_design"))
	s.PowerNowUW, _ = readInt64File(filepath.Join(dir, "power_now"))
	s.VoltageNowUV, _ = readInt64File(filepath.Join(dir, "voltage_now"))
	if cy, ok := readIntFile(filepath.Join(dir, "cycle_count")); ok {
		s.CycleCount = cy
	}

	// Prefer Mains/USB adapter with "online"
	for _, e := range ents {
		name := e.Name()
		if name == batName {
			continue
		}
		t := readTrimFile(filepath.Join(base, name, "type"))
		if t != "Mains" && t != "USB" && !strings.HasPrefix(strings.ToUpper(name), "ADP") {
			continue
		}
		onPath := filepath.Join(base, name, "online")
		if st, err := os.Stat(onPath); err != nil || st.IsDir() {
			continue
		}
		if v, ok := readIntFile(onPath); ok {
			b := v == 1
			s.ACOnline = &b
			s.ACName = name
			break
		}
	}

	return s, nil
}

// HealthPercent returns energy_full / energy_full_design * 100, or ok=false.
func (s *BatterySnapshot) HealthPercent() (pct float64, ok bool) {
	if s == nil || s.EnergyFullUWh <= 0 || s.EnergyDesignUWh <= 0 {
		return 0, false
	}
	return float64(s.EnergyFullUWh) * 100 / float64(s.EnergyDesignUWh), true
}

// PowerWatts returns power in watts (signed semantics same as sysfs).
func (s *BatterySnapshot) PowerWatts() (w float64, ok bool) {
	if s == nil || s.PowerNowUW == 0 {
		return 0, false
	}
	return float64(s.PowerNowUW) / 1e6, true
}

func readTrimFile(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

func readIntFile(path string) (int, bool) {
	s := readTrimFile(path)
	if s == "" {
		return 0, false
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return n, true
}

func readInt64File(path string) (int64, bool) {
	s := readTrimFile(path)
	if s == "" {
		return 0, false
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false
	}
	return n, true
}
