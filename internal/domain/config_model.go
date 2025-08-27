package domain

type Config struct {
	DB            DB
	WorkerCount   int32
	TimerInterval string
}
