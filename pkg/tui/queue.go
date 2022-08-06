package tui

import (
	"container/heap"
)

// priorityQueue implements heap.Interface and holds Views.
type priorityQueue []*View

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x any) {
	item := x.(*View)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

func (pq *priorityQueue) PushView(v *View) {
	heap.Push(pq, v)
}

func (pq *priorityQueue) PopView() *View {
	return heap.Pop(pq).(*View)
}

func newQueue() priorityQueue {
	pq := make(priorityQueue, 0, 16)
	heap.Init(&pq)
	return pq
}
