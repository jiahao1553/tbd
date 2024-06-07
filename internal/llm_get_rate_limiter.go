package internal

import "time"

func (o *OpenAI) GetRateLimiter() (semaphore chan struct{}, limiter *time.Ticker) {
	return getLimiter(500)
}

func (g *Groq) GetRateLimiter() (semaphore chan struct{}, limiter *time.Ticker) {
	return getLimiter(30)
}

func (a *Anthropic) GetRateLimiter() (semaphore chan struct{}, limiter *time.Ticker) {
	return getLimiter(5)
}

func getLimiter(mr int) (semaphore chan struct{}, limiter *time.Ticker) {
	i := time.Minute
	semaphore = make(chan struct{}, (mr / 2))
	limiter = time.NewTicker(i / time.Duration(mr/3))
	return semaphore, limiter
}
