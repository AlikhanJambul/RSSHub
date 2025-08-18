package models

type Config struct {
	DB            DB
	WorkerCount   int
	TimerInterval string
}
