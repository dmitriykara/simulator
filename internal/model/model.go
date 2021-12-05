package model

import (
	"math"
)

type Model struct {
	ArrivalTimeFunc   func() float64
	OperationTimeFunc func() float64
	workersNum        int64
	Queue             *Queue
	requests          map[int]*Request
	currentTime       int64
	arrivalTime       int64
}

func NewModel(workersNum int64, arrivalTimeFunc, operationTimeFunc func() float64) *Model {
	return &Model{
		ArrivalTimeFunc:   arrivalTimeFunc,
		OperationTimeFunc: operationTimeFunc,
		workersNum:        workersNum,
		Queue:             NewQueue(),
		requests:          make(map[int]*Request),
	}
}

type State struct {
	Processing    int64
	Waiting       int64
	ResponseTimes []int64
}

func (m *Model) GetCurrentState() *State {
	var (
		state = &State{
			ResponseTimes: make([]int64, 0),
		}
	)

	if m.currentTime == m.arrivalTime {
		// New request to the system
		newRequest := &Request{
			ArrivedAt:     m.currentTime,
			ProcessedTime: int64(math.Ceil(m.OperationTimeFunc())),
		}
		m.Queue.Push(newRequest)

		// Calculate next time when new request will come
		m.arrivalTime = m.currentTime + int64(math.Ceil(m.ArrivalTimeFunc()))
	}

	// Start service from request queue
	var i int64
	for i = 0; i < m.workersNum; i++ {
		if m.workersNum <= int64(len(m.requests)) {
			break
		}
		if m.Queue.Len() == 0 {
			break
		}
		requestID, request := m.Queue.Pop()
		request.ProcessingStarted = m.currentTime
		request.ProcessingFinished = m.currentTime + request.ProcessedTime - 1
		m.requests[requestID] = request
	}

	state.Processing = int64(len(m.requests))
	state.Waiting = m.Queue.Len()

	// Stop service
	for requestID, request := range m.requests {
		if request.ProcessingFinished == m.currentTime {
			delete(m.requests, requestID)
			state.ResponseTimes = append(state.ResponseTimes, request.WaitTime()+request.ProcessedTime)
		}
	}
	m.currentTime++

	return state
}
