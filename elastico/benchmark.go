package elastico

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/hamzaparekh/blockchain-sharding/crypto"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/process"
)

var payloadSizes = []int{64, 128, 256, 512, 1024, 2048, 4096, 8192}
var transactionsPerIteration = 5

func RunBenchmark() {
	file, err := os.Create("benchmark.csv")
	if err != nil {
		fmt.Println("❌ Error creating CSV:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Iteration", "PayloadSize", "NumTransactions", "Latency(s)", "CPU(%)", "Memory(MB)"})

	iteration := 1
	for _, size := range payloadSizes {
		fmt.Printf("\n📦 Testing payload size: %d bytes (Transactions per iteration: %d)\n", size, transactionsPerIteration)
		for i := 0; i < 50; i++ {
			fmt.Printf("⚙️ Iteration %d starting... (%d transactions)\n", iteration, transactionsPerIteration)

			cpuBefore, _ := cpu.Percent(100*time.Millisecond, false)
			start := time.Now()

			for t := 0; t < transactionsPerIteration; t++ {
				runElasticoTransaction(size)
			}

			latency := time.Since(start).Seconds()
			time.Sleep(300 * time.Millisecond)
			cpuAfter, _ := cpu.Percent(100*time.Millisecond, false)

			pid := int32(os.Getpid())
			proc, err := process.NewProcess(pid)
			if err != nil {
				fmt.Println("❌ Error fetching process info:", err)
			}
			memInfo, err := proc.MemoryInfo()
			if err != nil {
				fmt.Println("❌ Error fetching memory info:", err)
			}

			avgCPU := (cpuBefore[0] + cpuAfter[0]) / 2
			usedMemMB := float64(memInfo.RSS) / (1024 * 1024)

			fmt.Printf("Iteration %d complete. Latency: %.4fs (for %d tx) | CPU: %.2f%% | Mem: %.2fMB\n",
				iteration, latency, transactionsPerIteration, avgCPU, usedMemMB)

			writer.Write([]string{
				fmt.Sprintf("%d", iteration),
				fmt.Sprintf("%d", size),
				fmt.Sprintf("%d", transactionsPerIteration),
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

	simulateCommitteePBFT(payloadSize)
}

func simulateCommitteePBFT(payloadSize int) {
	payload := make([]byte, payloadSize)
	for i := 0; i < payloadSize; i++ {
		payload[i] = byte(i % 256)
	}

	for round := 0; round < 3; round++ {
		sum := 0
		for _, b := range payload {
			sum += int(b) % 7
		}
		time.Sleep(50 * time.Millisecond)
	}
}
