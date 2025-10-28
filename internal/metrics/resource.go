package metrics

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type ResourceMetrics struct {
	cpuUse    float64
	memUse    float64
	avgCpuUse float64
	avgMemUse float64

	count int64
	mu    sync.Mutex

	// For CPU calculation
	prevTotal uint64
	prevIdle  uint64
}

// GetSystemUsage gets the current system CPU and memory usage percentages using native /proc files
func (r *ResourceMetrics) GetSystemUsage() (currentCPU, avgCPU, currentMem, avgMem float64, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Get CPU usage
	currentCPU, err = r.getCPUUsage()
	if err != nil {
		return 0, 0, 0, 0, err
	}

	// Get memory usage
	currentMem, err = r.getMemoryUsage()
	if err != nil {
		return 0, 0, 0, 0, err
	}

	// Update running averages
	r.cpuUse = currentCPU
	r.memUse = currentMem
	r.count++
	r.avgCpuUse = (r.avgCpuUse*float64(r.count-1) + currentCPU) / float64(r.count)
	r.avgMemUse = (r.avgMemUse*float64(r.count-1) + currentMem) / float64(r.count)

	return r.cpuUse, r.avgCpuUse, r.memUse, r.avgMemUse, nil
}

func (r *ResourceMetrics) getCPUUsage() (float64, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return 0, fmt.Errorf("failed to read /proc/stat")
	}

	line := scanner.Text()
	if !strings.HasPrefix(line, "cpu ") {
		return 0, fmt.Errorf("unexpected format in /proc/stat")
	}

	fields := strings.Fields(line)[1:] // Skip "cpu"
	if len(fields) < 8 {
		return 0, fmt.Errorf("not enough fields in /proc/stat")
	}

	var total uint64
	for _, field := range fields {
		val, err := strconv.ParseUint(field, 10, 64)
		if err != nil {
			return 0, err
		}
		total += val
	}

	idle := fields[3] // idle is the 4th field (index 3)
	idleVal, err := strconv.ParseUint(idle, 10, 64)
	if err != nil {
		return 0, err
	}

	if r.prevTotal == 0 && r.prevIdle == 0 {
		// First call, just store and return 0
		r.prevTotal = total
		r.prevIdle = idleVal
		return 0, nil
	}

	totalDiff := total - r.prevTotal
	idleDiff := idleVal - r.prevIdle

	if totalDiff == 0 {
		return 0, nil
	}

	usage := float64(totalDiff-idleDiff) / float64(totalDiff) * 100.0

	r.prevTotal = total
	r.prevIdle = idleVal

	return usage, nil
}

func (r *ResourceMetrics) getMemoryUsage() (float64, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var total, available uint64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			total, err = parseMemValue(line)
			if err != nil {
				return 0, err
			}
		} else if strings.HasPrefix(line, "MemAvailable:") {
			available, err = parseMemValue(line)
			if err != nil {
				return 0, err
			}
			break // MemAvailable is after MemTotal
		}
	}

	if total == 0 {
		return 0, fmt.Errorf("MemTotal not found")
	}

	used := total - available
	percent := float64(used) / float64(total) * 100.0

	return percent, nil
}

func parseMemValue(line string) (uint64, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return 0, fmt.Errorf("invalid meminfo line: %s", line)
	}
	return strconv.ParseUint(fields[1], 10, 64)
}
