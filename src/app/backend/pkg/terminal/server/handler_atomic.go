package server

import (
	"sync"
	"time"
)

// counter - Web Socket 연결 관리 정보 구조체
type Counter struct {
	duration    time.Duration
	zeroTimer   *time.Timer
	wg          sync.WaitGroup
	connections int
	mutex       sync.Mutex
}

// newCounter - 지저한 시간동안 유지할 Counter 인스턴스 생성
func NewCounter(duration time.Duration) *Counter {
	zeroTimer := time.NewTimer(duration)

	// 시간 구간이 0인 경우는 사용자가 전달되지 않도록 만료 이벤트 차단
	if duration == 0 {
		<-zeroTimer.C
	}

	return &Counter{
		duration:  duration,
		zeroTimer: zeroTimer,
	}
}

// add - 지정한 수의 연결 추가
func (counter *Counter) Add(n int) int {
	counter.mutex.Lock()
	defer counter.mutex.Unlock()

	if counter.duration > 0 {
		counter.zeroTimer.Stop()
	}
	counter.wg.Add(n)
	counter.connections += n

	return counter.connections
}

// done - 하나의 연결 해제
func (counter *Counter) Done() int {
	counter.mutex.Lock()
	defer counter.mutex.Unlock()

	counter.connections--
	counter.wg.Done()
	if counter.connections == 0 && counter.duration > 0 {
		counter.zeroTimer.Reset(counter.duration)
	}

	return counter.connections
}

// count - 관리중인 연결 수 반환
func (counter *Counter) count() int {
	counter.mutex.Lock()
	defer counter.mutex.Unlock()

	return counter.connections
}

// wait - 연결 수 조정 대기
func (counter *Counter) wait() {
	counter.wg.Wait()
}

// timer - 관리 중인 Timer 반환
func (counter *Counter) Timer() *time.Timer {
	return counter.zeroTimer
}
