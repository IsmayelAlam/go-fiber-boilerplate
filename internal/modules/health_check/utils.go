package healthCheck

import (
	"errors"
	"runtime"
)

func (s *healthCheckService) dbCheck() error {
	if err := s.db.Ping(); err != nil {
		return err
	}
	return nil
}

func (s *healthCheckService) memoryHeapCheck() error {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats) // Collect memory statistics

	heapAlloc := memStats.HeapAlloc            // Heap memory currently allocated
	heapThreshold := uint64(300 * 1024 * 1024) // Example threshold: 300 MB

	// If the heap allocation exceeds the threshold, return an error
	if heapAlloc > heapThreshold {
		return errors.New("heap memory usage too high")
	}
	return nil
}

func (s *healthCheckService) addServiceStatus(serviceList *[]HealthCheck, name string, isUp bool, message *string) {
	status := "Up"

	if !isUp {
		status = "Down"
	}

	*serviceList = append(*serviceList, HealthCheck{
		Name:    name,
		Status:  status,
		IsUp:    isUp,
		Message: message,
	})
}
