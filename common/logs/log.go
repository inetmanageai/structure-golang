package logs

// Implement PORT
type AppLog interface {
	// Info logging เวลามีการทำงานต่างๆ
	Info(msg string)

	// Debug Logging เวลาทดสอบการทำงาน
	Debug(msg string)

	// Warning Logging เวลามีความผิดผลาด
	Warning(msg string)

	// Error Logging เมื่อเกิดความผิดพลาดร้ายแรง
	Error(msg interface{})
}
