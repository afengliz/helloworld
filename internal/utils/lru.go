package utils

import "sync"

// 链表节点
type node struct {
	// 键
	Key string
	// 值
	Val interface{}
	// 前一个节点
	last *node
	// 后一个节点
	next *node
}

func newNode(key string, val interface{}) *node {
	return &node{Key: key, Val: val}
}

// 双向链表
type link struct {
	// 头结点
	Head *node
	// 尾结点
	Tail *node
}

func newLink() *link {
	return &link{}
}

// 在链表尾巴添加新节点
func (l *link) add(node *node) {
	if node == nil {
		return
	}
	if l.Head == nil {
		l.Head = node
		l.Tail = node
	} else {
		l.Tail.next = node
		node.last = l.Tail
		l.Tail = node
	}
}

// 把节点移动到尾巴
func (l *link) moveNodeToTail(node *node) {
	if node == nil {
		return
	}
	if node == l.Tail {
		return
	}
	if l.Head == node {
		l.Head = l.Head.next
		l.Head.last = nil
	} else {
		node.next.last = node.last
		node.last.next = node.next
	}
	l.Tail.next = node
	node.last = l.Tail
	node.next = nil
	l.Tail = node
}

// 移除头部的一个节点
func (l *link) removeHead() *node {
	if l.Head == nil {
		return nil
	}
	head := l.Head
	if l.Head == l.Tail {
		l.Head = nil
		l.Tail = nil
		return head
	}
	l.Head = head.next
	l.Head.last = nil
	head.next = nil
	return head
}

type LRU interface {
	Put(key string, val interface{})
	Get(key string) interface{}

}

type lru struct {
	// 哈希表 用于存储节点所在的位置
	NodeMap map[string]*node
	// 链表
	NodeLink *link
	// 容量大小
	capacity int
	// 锁
	mutex sync.Mutex
}

func NewLRU(c int) LRU {
	myLru := lru{capacity: c}
	myLru.NodeLink = newLink()
	myLru.NodeMap = make(map[string]*node)
	return &myLru
}
func (m *lru) Get(key string) interface{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	cNode := m.NodeMap[key]
	if cNode == nil {
		return nil
	}
	m.NodeLink.moveNodeToTail(cNode)
	return cNode.Val
}

func (m *lru) Put(key string, val interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if cNode, ok := m.NodeMap[key]; ok {
		cNode.Val = val
		m.NodeLink.moveNodeToTail(cNode)
		return
	}
	nNode := newNode(key, val)
	m.NodeLink.add(nNode)
	m.NodeMap[nNode.Key] = nNode
	if len(m.NodeMap) == m.capacity+1 {
		head := m.NodeLink.removeHead()
		delete(m.NodeMap, head.Key)
	}
}
