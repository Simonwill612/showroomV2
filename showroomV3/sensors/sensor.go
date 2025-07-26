package sensors

import (
	"fmt"
	"time"

	"github.com/d2r2/go-i2c"
	"github.com/stianeikeland/go-rpio/v4"
)

// GPIO for VL53L0X XSHUT
var (
	XSHUT1 = rpio.Pin(22) // GPIO22
	XSHUT2 = rpio.Pin(23) // GPIO23
)

// I2C Address
const (
	VL53L0X_DEFAULT_ADDR = 0x29
	VL53L0X_NEW_ADDR     = 0x30
)

// SPI MCP3008
const (
	CH_ACS1 = 0
	CH_ACS2 = 1
)

// Init initializes GPIO
func Init() error {
	if err := rpio.Open(); err != nil {
		return fmt.Errorf("failed to init GPIO: %v", err)
	}
	XSHUT1.Output()
	XSHUT2.Output()
	return nil
}

// SetXSHUT controls power to VL53L0X via XSHUT pin
func SetXSHUT(sensor int, state bool) {
	var pin rpio.Pin
	switch sensor {
	case 1:
		pin = XSHUT1
	case 2:
		pin = XSHUT2
	default:
		return
	}
	if state {
		pin.High()
	} else {
		pin.Low()
	}
}

// InitI2C initializes and returns the I2C bus
// InitI2C initializes and returns the I2C bus to sensor at given address
func InitI2C(addr uint8) (*i2c.I2C, error) {
	bus, err := i2c.NewI2C(addr, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to init I2C bus: %v", err)
	}
	return bus, nil
}

// ReadACS712 simulates reading current sensor
func ReadACS712(channel int) (int, error) {
	switch channel {
	case CH_ACS1:
		return 512, nil
	case CH_ACS2:
		return 530, nil
	default:
		return 0, fmt.Errorf("invalid channel")
	}
}

// ExampleVL53Init shows how to reset and enable sensors
func ExampleVL53Init() {
	fmt.Println("Resetting both sensors")
	SetXSHUT(1, false)
	SetXSHUT(2, false)
	time.Sleep(10 * time.Millisecond)

	fmt.Println("Starting sensor 1 (will set new address)")
	SetXSHUT(1, true)
	time.Sleep(10 * time.Millisecond)
	// TODO: change sensor 1 I2C address to VL53L0X_NEW_ADDR

	fmt.Println("Starting sensor 2 (default address)")
	SetXSHUT(2, true)
	time.Sleep(10 * time.Millisecond)
}

// Đọc chiều cao từ 1 cảm biến VL53L0X
func GetHeightVL53L0X(addr uint8) (int, error) {
	bus, err := i2c.NewI2C(addr, 1)
	if err != nil {
		return 0, fmt.Errorf("I2C init failed: %v", err)
	}
	defer bus.Close()

	// TODO: Đọc thực tế từ VL53L0X, đây là giả lập
	time.Sleep(20 * time.Millisecond)
	return 100, nil // giả lập 100cm
}

// Đọc chiều cao từ cả 2 cảm biến
func GetBothHeights() (int, int, error) {
	left, err1 := GetHeightVL53L0X(0x30)
	right, err2 := GetHeightVL53L0X(0x29)
	if err1 != nil || err2 != nil {
		return 0, 0, fmt.Errorf("read error: %v %v", err1, err2)
	}
	return left, right, nil
}
