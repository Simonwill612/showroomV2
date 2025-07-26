package showroomV

import (
	"sync"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

var (
	relayPin = rpio.Pin(2)
	pwmL1    = rpio.Pin(6)
	pwmR1    = rpio.Pin(7)
	lEn1     = rpio.Pin(4)
	rEn1     = rpio.Pin(5)
	pwmL2    = rpio.Pin(8)
	pwmR2    = rpio.Pin(9)
	lEn2     = rpio.Pin(10)
	rEn2     = rpio.Pin(11)
	acs1Pin  = 1 // analog channel
	acs2Pin  = 3 // analog channel

	pulsesPerCm = 200.0

	mu            sync.Mutex
	isMoving      bool
	moveDir       string
	pulseCount1   int64
	pulseCount2   int64
	currentHeight = 80
)

func InitTableGPIO() error {
	if err := rpio.Open(); err != nil {
		return err
	}
	relayPin.Output()
	pwmL1.Output()
	pwmR1.Output()
	lEn1.Output()
	rEn1.Output()
	pwmL2.Output()
	pwmR2.Output()
	lEn2.Output()
	rEn2.Output()
	relayPin.High()
	return nil
}

func moveUp() {
	mu.Lock()
	defer mu.Unlock()
	lEn1.High()
	rEn1.High()
	pwmL1.High()
	pwmR1.Low()
	lEn2.High()
	rEn2.High()
	pwmL2.High()
	pwmR2.Low()
	isMoving = true
	moveDir = "UP"
}

func moveDown() {
	mu.Lock()
	defer mu.Unlock()
	lEn1.High()
	rEn1.High()
	pwmL1.Low()
	pwmR1.High()
	lEn2.High()
	rEn2.High()
	pwmL2.Low()
	pwmR2.High()
	isMoving = true
	moveDir = "DOWN"
}

func stopAll() {
	mu.Lock()
	defer mu.Unlock()
	pwmL1.Low()
	pwmR1.Low()
	pwmL2.Low()
	pwmR2.Low()
	lEn1.High()
	rEn1.High()
	lEn2.High()
	rEn2.High()
	isMoving = false
	moveDir = ""
}

func GetCurrentHeight() int {
	mu.Lock()
	defer mu.Unlock()
	return 80 + int(float64(pulseCount1)/pulsesPerCm)
}

// Giả lập tăng pulse khi nâng/hạ (thực tế cần interrupt encoder)
func SimulatePulse(dir string, durationMs int) {
	mu.Lock()
	defer mu.Unlock()
	delta := int64(float64(durationMs) / 10) // mỗi 10ms = 1 pulse
	if dir == "UP" {
		pulseCount1 += delta
		pulseCount2 += delta 
	} else if dir == "DOWN" {
		pulseCount1 -= delta
		pulseCount2 -= delta
	}
} 

// Hàm nâng/hạ bàn theo thời gian
func MoveTable(dir string, durationMs int) {
	if dir == "UP" {
		moveUp()
	} else {
		moveDown()
	}
	go func() {
		time.Sleep(time.Duration(durationMs) * time.Millisecond)
		stopAll()
		SimulatePulse(dir, durationMs)
	}()
}

// Hàm nâng/hạ bàn đến chiều cao cụ thể
func MoveToHeight(target int) {
	cur := GetCurrentHeight()
	if target > cur {
		delta := target - cur
		duration := int(float64(delta) * pulsesPerCm * 10) // giả lập
		MoveTable("UP", duration)
	} else if target < cur {
		delta := cur - target
		duration := int(float64(delta) * pulsesPerCm * 10)
		MoveTable("DOWN", duration)
	}
}

// // Hàm đọc dòng điện (giả lập)
// func ReadCurrent() float64 {
// 	// TODO: Đọc từ MCP3008/ACS712
// 	return 2.0 // amps
// }

// // Hàm kiểm tra quá dòng
// func CheckOvercurrent() bool {
// 	return ReadCurrent() > 15.0
// }
