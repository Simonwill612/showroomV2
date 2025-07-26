package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
)

// const eventName = "Emit" 
 
func main() {
	loadEnv()

	rpcURL := os.Getenv("WS_RPC_URL")
	contractAddr := os.Getenv("CONTRACT_ADDRESS")
	abiFile := os.Getenv("CONTRACT_ABI_FILE")

	client := connectToClient(rpcURL)
	defer client.Close()

	parsedABI := parseContractABIFromFile(abiFile)
	logs := make(chan types.Log)
	sub := subscribeToEvents(client, logs, contractAddr)

	fmt.Println("✅ Listening for Emit events from tableHandleAI...")

	handleLogs(client, parsedABI, logs, sub)
}

// --- Load .env
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load .env file: %v", err)
	}
}

// --- Connect to WebSocket
func connectToClient(url string) *ethclient.Client {
	rpcClient, err := rpc.DialWebsocket(context.Background(), url, "")
	if err != nil {
		log.Fatalf("❌ Failed to connect to WebSocket: %v", err)
	}
	return ethclient.NewClient(rpcClient)
}

// --- Read ABI from file
func parseContractABIFromFile(path string) abi.ABI {
	abiBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("❌ Failed to read ABI file: %v", err)
	}
	parsed, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		log.Fatalf("❌ Failed to parse ABI: %v", err)
	}
	return parsed
}

// --- Subscribe to ALL events
func subscribeToEvents(client *ethclient.Client, logs chan types.Log, contractAddr string) ethereum.Subscription {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(contractAddr)},
		// Không cần Topics, sẽ nhận tất cả event
	}
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalf("❌ Failed to subscribe to logs: %v", err)
	}
	return sub
}

// --- Handle ALL event logs
func handleLogs(client *ethclient.Client, parsedABI abi.ABI, logs chan types.Log, sub ethereum.Subscription) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			chainID, err := client.ChainID(context.Background())
			if err != nil {
				log.Printf("❌ Failed to get chain ID: %v", err)
			} else {
				fmt.Printf("🔗 Chain ID: %v\n", chainID)
			}

		case err := <-sub.Err():
			log.Printf("❌ Subscription error: %v", err)
			time.Sleep(10 * time.Second)

		case vLog := <-logs:
			event, err := parsedABI.EventByID(vLog.Topics[0])
			if err != nil {
				log.Printf("⚠️ Unknown event: %v", err)
				continue
			}

			// In ra tên event và dữ liệu
			fmt.Printf("📢 Event: %s\n", event.Name)
			switch event.Name {
			case "RaiseTable", "LowerTable", "StopTable":
				var evt struct {
					TableId string
				}
				err := parsedABI.UnpackIntoInterface(&evt, event.Name, vLog.Data)
				if err == nil {
					log.Printf("➡️ %s: tableId = %s", event.Name, evt.TableId)
				}
			case "HeightSet":
				var evt struct {
					TableId string
					Height  uint64
				}
				err := parsedABI.UnpackIntoInterface(&evt, event.Name, vLog.Data)
				if err == nil {
					log.Printf("📏 HeightSet: tableId = %s | height = %d", evt.TableId, evt.Height)
				}
			case "RelayToggled":
				var evt struct {
					TableId string
					On      bool
				}
				err := parsedABI.UnpackIntoInterface(&evt, event.Name, vLog.Data)
				if err == nil {
					log.Printf("🔌 RelayToggled: tableId = %s | on = %v", evt.TableId, evt.On)
				}
			case "TableCommand":
				var evt struct {
					TableId string
					Command string
				}
				err := parsedABI.UnpackIntoInterface(&evt, event.Name, vLog.Data)
				if err == nil {
					log.Printf("📝 TableCommand: TableId = %s | command = %s", evt.TableId, evt.Command)
				}
			case "EmitTableData":
				var evt struct {
					TableId   string
					HeightsCm string
				}
				err := parsedABI.UnpackIntoInterface(&evt, event.Name, vLog.Data)
				if err == nil {
					log.Printf("📦 EmitTableData: tableId = %s | heightsCm = %s", evt.TableId, evt.HeightsCm)
					go handleHeightChange(evt.TableId, evt.HeightsCm)

				}
			default:
				log.Printf("⚠️ Event %s không được xử lý tự động", event.Name)
			}
		}
	}
}

// --- Logic xử lý chiều cao bàn thông minh
func handleHeightChange(tableId string, heights string) {
	if tableId != "1" {
		log.Println("ℹ️ Bỏ qua TableId không hợp lệ:", tableId)
		return
	}

	parts := strings.Split(heights, ",")
	if len(parts) < 2 {
		log.Println("⚠️ heightsCm không hợp lệ:", heights)
		return
	}

	currentHeight, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	targetHeight, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err1 != nil || err2 != nil {
		log.Printf("❌ Lỗi phân tích chiều cao: %v | %v\n", err1, err2)
		return
	}

	diff := targetHeight - currentHeight
	if diff <= 0 {
		log.Println("⛔️ Không nâng vì chiều cao hiện tại lớn hơn hoặc bằng mục tiêu.")
		return
	}

	// Tính thời gian nâng
	durationMs := 2000
	if diff > 10 {
		durationMs += ((diff - 10) / 10) * 2000
	}
	log.Printf("🔼 Nâng từ %dcm lên %dcm | Thời gian nâng: %dms\n", currentHeight, targetHeight, durationMs)

	// 👉 Tại đây, bạn sẽ tích hợp GPIO hoặc command nội bộ để điều khiển motor

	time.Sleep(time.Duration(durationMs) * time.Millisecond)
	log.Println("⏹️ Dừng sau nâng")

	// Hạ sau 60 giây
	time.AfterFunc(60*time.Second, func() {
		downTarget := targetHeight / 2
		downDiff := targetHeight - downTarget
		downDurationMs := 2000
		if downDiff > 10 {
			downDurationMs += ((downDiff - 10) / 10) * 2000
		}
		log.Printf("⏬ Hạ từ %dcm về %dcm sau 60s | Thời gian hạ: %dms\n", targetHeight, downTarget, downDurationMs)

		// 👉 Thêm điều khiển hạ motor tại đây
		time.Sleep(time.Duration(downDurationMs) * time.Millisecond)
		log.Println("⏹️ Dừng sau hạ")
	})
}
