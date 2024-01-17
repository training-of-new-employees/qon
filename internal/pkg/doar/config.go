package doar

// Config - конфигурация для отправителя писем.
type SenderConfig struct {
	Mode SenderMode

	// Используется при режимах SmptMode и ApiMode
	SenderEmail string

	// Используется при режиме SmptMode
	SenderPassword string

	// Используется при режиме ApiMode
	SenderApiKey string
}
