package parser

type tokenStack struct {
	Head *tokenStackNode
	Size int
}

type tokenStackNode struct {
	Token *Token
	Prev  *tokenStackNode
}

func (s *tokenStack) push(t *Token) {
	node := tokenStackNode{t, s.Head}
	s.Head = &node
	s.Size++
}

func (s *tokenStack) pop() *Token {
	node := s.Head
	s.Head = node.Prev
	s.Size--
	return node.Token
}

func (s *tokenStack) peek() *Token {
	return s.Head.Token
}

func (s *tokenStack) empty() bool {
	return s.Head == nil
}
