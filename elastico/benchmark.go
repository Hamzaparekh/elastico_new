package elastico

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/hamzaparekh/blockchain-sharding/crypto"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

var payloadSizes = []int{64, 256, 512, 1024} // Bytes

func RunBenchmark() {
	file, err := os.Create("benchmark.csv")
	if err != nil {
		fmt.Println("❌ Error creating CSV:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Iteration", "PayloadSize", "Latency(s)", "CPU(%)", "Memory(MB)"})

	iteration := 1
	for _, size := range payloadSizes {
		fmt.Printf("\n📦 Testing payload size: %d bytes\n", size)
		for i := 0; i < 5; i++ {
			fmt.Printf("⚙️ Iteration %d starting...\n", iteration)

			cpuBefore, _ := cpu.Percent(100*time.Millisecond, false)
			start := time.Now()

			runElasticoTransaction(size)

			latency := time.Since(start).Seconds()
			time.Sleep(300 * time.Millisecond)
			cpuAfter, _ := cpu.Percent(100*time.Millisecond, false)
			vm, _ := mem.VirtualMemory()

			avgCPU := (cpuBefore[0] + cpuAfter[0]) / 2
			usedMemMB := float64(vm.Used) / (1024 * 1024)

			fmt.Printf("✅ Iteration %d complete. Latency: %.4fs | CPU: %.2f%% | Mem: %.2fMB\n",
				iteration, latency, avgCPU, usedMemMB)

			writer.Write([]string{
				fmt.Sprintf("%d", iteration),
				fmt.Sprintf("%d", size),
				fmt.Sprintf("%.4f", latency),
				fmt.Sprintf("%.2f", avgCPU),
				fmt.Sprintf("%.2f", usedMemMB),
			})

			iteration++
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		fmt.Println("❌ CSV flush error:", err)
	}
}

// 🔁 Simulates a real Elastico transaction + PBFT round
func runElasticoTransaction(payloadSize int) {
	fmt.Println("🔁 Simulating Elastico transaction + PBFT...")

	addr := "localhost:9388"
	sk, err := crypto.NewKey()
	if err != nil {
		fmt.Println("❌ Failed to generate key:", err)
		return
	}

	proof := NewIDProof(addr, sk.D.Bytes())
	if !proof.Verify() {
		fmt.Println("❌ Identity proof verification failed")
		return
	}

	_ = proof.GetCommitteeNo()

	// Simulate committee-level consensus using dummy payload
	simulateCommitteePBFT(payloadSize)
}

func simulateCommitteePBFT(payloadSize int) {
	payload := make([]byte, payloadSize)
	for i := 0; i < payloadSize; i++ {
		payload[i] = byte(i % 256)
	}

	// Simulate 3 PBFT rounds: pre-prepare, prepare, commit
	for round := 0; round < 3; round++ {
		sum := 0
		for _, b := range payload {
			sum += int(b) % 7 // dummy computation
		}
		time.Sleep(50 * time.Millisecond) // simulate delay
	}
}
