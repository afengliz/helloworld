package utils

type HeapNode struct {
	Key   string
	Count int
}

/**
大根堆
*/
type BigRootHeap []*HeapNode

// 初始化
func InitBigRootHeap(array ...*HeapNode) *BigRootHeap {
	intHeap := BigRootHeap{}
	for i := 0; i < len(array); i++ {
		intHeap = append(intHeap, array[i])
	}
	return &intHeap
}

func (p *BigRootHeap) Len() int {
	return len(*p)
}

func (p *BigRootHeap) Less(i, j int) bool {
	return (*p)[i].Count > (*p)[j].Count
}
func (p *BigRootHeap) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

func (p *BigRootHeap) Push(x interface{}) {
	*p = append(*p, x.(*HeapNode))
}

func (h *BigRootHeap) Pop() interface{} {
	// 最小值被放在最后一个位置，
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
func (h *BigRootHeap) Front() interface{} {
	old := *h
	n := len(old)
	if n <= 0 {
		return nil
	}
	return old[0]
}


