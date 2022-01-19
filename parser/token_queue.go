package parser

import "errors"

type tokenQueue struct {
	Head *tokenQueueNode
	Tail *tokenQueueNode
}

type tokenQueueNode struct {
	Token *Token
	Prev  *tokenQueueNode
	Next  *tokenQueueNode
}

func (q *tokenQueue) enqueue(t *Token) {
	node := tokenQueueNode{t, q.Tail, nil}
	//fmt.Println(t.Value)

	if q.Tail == nil {
		q.Head = &node
	} else {
		q.Tail.Next = &node
	}

	q.Tail = &node
}

func (q *tokenQueue) dequeue() *Token {
	node := q.Head
	if node.Next != nil {
		node.Next.Prev = nil
	}
	q.Head = node.Next
	if q.Head == nil {
		q.Tail = nil
	}
	return node.Token
}

func (q *tokenQueue) empty() bool {
	return q.Head == nil && q.Tail == nil
}

type nodeStack struct {
	Head *nodeStackNode
}

type nodeStackNode struct {
	ParseNode *ParseNode
	Prev      *nodeStackNode
}

func (s *nodeStack) push(n *ParseNode) {
	node := nodeStackNode{n, s.Head}
	s.Head = &node
}

func (s *nodeStack) pop() (*ParseNode, error) {
	node := s.Head
	if node == nil {
		return nil, errors.New("child node not available")
	}
	s.Head = node.Prev
	return node.ParseNode, nil
}

func (s *nodeStack) peek() *ParseNode {
	return s.Head.ParseNode
}
