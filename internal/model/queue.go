package model

import "container/list"

func NewRequest() *Request {
	return &Request{}
}

type Request struct {
	ArrivedAt          int64
	ProcessedTime      int64
	ProcessingStarted  int64
	ProcessingFinished int64
}

func (r Request) WaitTime() int64 {
	return r.ProcessingStarted - r.ArrivedAt
}

type Queue struct {
	queue *list.List
	list  []*Request
}

func NewQueue() *Queue {
	return &Queue{
		queue: list.New(),
		list:  make([]*Request, 0),
	}
}

func (q *Queue) Push(r *Request) {
	q.list = append(q.list, r)
	q.queue.PushBack(len(q.list) - 1)
}

func (q *Queue) Pop() (int, *Request) {
	elem := q.queue.Front()
	defer q.queue.Remove(elem)

	reqId, _ := elem.Value.(int)
	return reqId, q.list[reqId]
}

func (q *Queue) Requests() []*Request {
	return q.list
}

func (q Queue) Len() int64 {
	return int64(q.queue.Len())
}
