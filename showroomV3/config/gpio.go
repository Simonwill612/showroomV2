package config

const ( 
	// Relay Pins
	RELAY_PIN_1 = 17 // GPIO17 - Pin 11
	RELAY_PIN_2 = 27 // GPIO27 - Pin 13

	// Motor 1 (BTS7960)
	PWM_L1 = 12 // GPIO12 - Pin 32 (PWM0)
	PWM_R1 = 13 // GPIO13 - Pin 33 (PWM1)
	L_EN1  = 5  // GPIO5  - Pin 29
	R_EN1  = 6  // GPIO6  - Pin 31

	// Motor 2 (BTS7960)
	PWM_L2 = 18 // GPIO18 - Pin 12 (PWM0)
	PWM_R2 = 19 // GPIO19 - Pin 35 (PWM1)
	L_EN2  = 20 // GPIO20 - Pin 38
	R_EN2  = 21 // GPIO21 - Pin 40

	// VL53L0X XSHUT control
	VL53L0X_1_XSHUT = 22 // GPIO22 - Pin 15
	VL53L0X_2_XSHUT = 23 // GPIO23 - Pin 16

	// MCP3008 SPI (assumed)
	SPI_MOSI = 10 // GPIO10 - Pin 19
	SPI_MISO = 9  // GPIO9  - Pin 21
	SPI_CLK  = 11 // GPIO11 - Pin 23
	SPI_CS   = 8  // GPIO8  - Pin 24

	// MCP3008 Channels for ACS712 current sensors
	ACS1_CHANNEL = 0 // CH0
	ACS2_CHANNEL = 1 // CH1

	// I2C pins (shared for both VL53L0X)
	I2C_SDA = 2 // GPIO2 - Pin 3
	I2C_SCL = 3 // GPIO3 - Pin 5
)
