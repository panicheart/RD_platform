package collectors

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"rdp-platform/rdp-api/models"
	"rdp-platform/rdp-api/services"
)

// SystemCollector collects system metrics periodically
type SystemCollector struct {
	monitorService *services.MonitorService
	interval       time.Duration
	stopChan       chan bool
	isRunning      bool
	dataDir        string
}

// NewSystemCollector creates a new SystemCollector
func NewSystemCollector(monitorService *services.MonitorService, interval time.Duration) *SystemCollector {
	if interval <= 0 {
		interval = 60 * time.Second
	}

	return &SystemCollector{
		monitorService: monitorService,
		interval:       interval,
		stopChan:       make(chan bool),
		isRunning:      false,
		dataDir:        "/var/lib/rdp/metrics",
	}
}

// Start begins the metrics collection loop
func (c *SystemCollector) Start() {
	if c.isRunning {
		return
	}

	c.isRunning = true
	go c.collectionLoop()
	fmt.Println("System metrics collector started")
}

// Stop halts the metrics collection
func (c *SystemCollector) Stop() {
	if !c.isRunning {
		return
	}

	c.isRunning = false
	close(c.stopChan)
	fmt.Println("System metrics collector stopped")
}

// IsRunning returns whether the collector is running
func (c *SystemCollector) IsRunning() bool {
	return c.isRunning
}

// collectionLoop runs the periodic collection
func (c *SystemCollector) collectionLoop() {
	// Collect immediately on start
	c.collect()

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.collect()
		case <-c.stopChan:
			return
		}
	}
}

// collect gathers all system metrics and saves them
func (c *SystemCollector) collect() {
	ctx := context.Background()

	metric := &models.SystemMetric{
		Timestamp: time.Now(),
	}

	// Collect CPU usage (simplified)
	c.collectCPU(metric)

	// Collect memory usage
	c.collectMemory(metric)

	// Collect disk usage
	c.collectDisk(metric)

	// Collect basic network stats (placeholder)
	c.collectNetwork(metric)

	// Save to database
	if err := c.monitorService.CreateSystemMetric(ctx, metric); err != nil {
		fmt.Printf("Failed to save system metric: %v\n", err)
	}
}

// collectCPU collects CPU usage percentage using Go runtime
func (c *SystemCollector) collectCPU(metric *models.SystemMetric) {
	// Simplified CPU usage estimation based on Goroutines
	// In production, use gopsutil library
	numGoroutine := runtime.NumGoroutine()

	// Very rough estimation: more goroutines = higher CPU usage
	cpuUsage := float64(numGoroutine) * 0.5
	if cpuUsage > 100 {
		cpuUsage = 100
	}

	metric.CPUUsage = cpuUsage
}

// collectMemory collects memory usage statistics
func (c *SystemCollector) collectMemory(metric *models.SystemMetric) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Use Go runtime memory stats as approximation
	metric.MemoryUsed = int64(m.Sys)
	metric.MemoryTotal = int64(m.Sys) + 1024*1024*1024 // Assume 1GB additional
	if metric.MemoryTotal > 0 {
		metric.MemoryUsage = float64(m.Sys) / float64(metric.MemoryTotal) * 100
	}
}

// collectDisk collects disk usage statistics
func (c *SystemCollector) collectDisk(metric *models.SystemMetric) {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "/"
	}

	// Try to get disk usage using platform-specific implementation
	usage := getDiskUsage(cwd)
	metric.DiskTotal = usage.Total
	metric.DiskUsed = usage.Used
	if usage.Total > 0 {
		metric.DiskUsage = float64(usage.Used) / float64(usage.Total) * 100
	}
}

// DiskUsage represents disk usage information
type DiskUsage struct {
	Total int64
	Used  int64
	Free  int64
}

// getDiskUsage returns disk usage for the given path (platform-specific)
func getDiskUsage(path string) DiskUsage {
	// This is a placeholder - platform-specific implementation needed
	// For now, return dummy data
	return DiskUsage{
		Total: 100 * 1024 * 1024 * 1024, // 100 GB
		Used:  50 * 1024 * 1024 * 1024,  // 50 GB
		Free:  50 * 1024 * 1024 * 1024,  // 50 GB
	}
}

// collectNetwork collects network I/O statistics (placeholder)
func (c *SystemCollector) collectNetwork(metric *models.SystemMetric) {
	// Network metrics require platform-specific implementation
	// For now, leave as 0 or implement using /proc/net/dev on Linux
	// In production, use gopsutil library
	metric.NetworkIn = 0
	metric.NetworkOut = 0
}

// GetSystemInfo returns static system information
func GetSystemInfo() map[string]interface{} {
	info := make(map[string]interface{})

	// Go runtime info
	info["go_version"] = runtime.Version()
	info["go_os"] = runtime.GOOS
	info["go_arch"] = runtime.GOARCH
	info["num_cpu"] = runtime.NumCPU()
	info["num_goroutine"] = runtime.NumGoroutine()

	// Memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	info["memory_alloc"] = m.Alloc
	info["memory_sys"] = m.Sys
	info["memory_heap_alloc"] = m.HeapAlloc
	info["memory_heap_sys"] = m.HeapSys

	return info
}

// GetDiskPartitions returns disk partition information (simplified)
func GetDiskPartitions() ([]map[string]interface{}, error) {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "/"
	}

	usage := getDiskUsage(cwd)

	result := []map[string]interface{}{
		{
			"device":       "root",
			"mountpoint":   cwd,
			"fstype":       "unknown",
			"total":        usage.Total,
			"used":         usage.Used,
			"free":         usage.Free,
			"used_percent": float64(usage.Used) / float64(usage.Total) * 100,
		},
	}

	return result, nil
}
