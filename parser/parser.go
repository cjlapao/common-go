package parser

import "errors"

// Parser parser structure
type Parser struct {
	// Map from string inputs to operator types
	Operators map[string]*Operator
	// Map from string inputs to function types
	Functions map[string]*Function
}

// EmptyParser create empty parser
func EmptyParser() *Parser {
	return &Parser{make(map[string]*Operator), make(map[string]*Function)}
}

// DefineOperator Adds an operator to the language. Provide the token, a precedence, and
// whether the operator is left, right, or not associative.
func (p *Parser) DefineOperator(token string, operands, assoc, precedence int) {
	p.Operators[token] = &Operator{token, assoc, operands, precedence}
}

// DefineFunction Adds a function to the language
func (p *Parser) DefineFunction(token string, params int) {
	p.Functions[token] = &Function{token, params}
}

func (p *Parser) Parse(tokens []*Token) (*ParseNode, error) {
	postfix, err := p.infixToPostfix(tokens)
	if err != nil {
		return nil, err
	}
	tree, err := p.postfixToTree(postfix)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

// InfixToPostfix Parses the input string of tokens using the given definitions of operators
// and functions. (Everything else is assumed to be a literal.) Uses the
// Shunting-Yard algorithm.
func (p *Parser) infixToPostfix(tokens []*Token) (*tokenQueue, error) {
	queue := tokenQueue{}
	stack := tokenStack{}
	wasLiteral := false // We use this bool to see if the last token was a literal

	for len(tokens) > 0 {
		token := tokens[0]
		tokens = tokens[1:]

		if _, ok := p.Functions[token.stringValue]; ok {
			// push functions onto the stack
			stack.push(token)
			wasLiteral = false
		} else if token.stringValue == "," {
			// function parameter separator, pop off stack until we see a "("
			for stack.peek().stringValue != "(" || stack.empty() {
				queue.enqueue(stack.pop())
			}
			// there was an error parsing
			if stack.empty() {
				return nil, errors.New("Parse error")
			}
			wasLiteral = false
		} else if o1, ok := p.Operators[token.stringValue]; ok {
			// push operators onto stack according to precedence
			if !stack.empty() {
				for o2, ok := p.Operators[stack.peek().stringValue]; ok &&
					(o1.Association == OpAssociationLeft && o1.Precedence <= o2.Precedence) ||
					(o1.Association == OpAssociationRight && o1.Precedence < o2.Precedence); {
					queue.enqueue(stack.pop())

					if stack.empty() {
						break
					}
					o2, ok = p.Operators[stack.peek().stringValue]
				}
			}
			stack.push(token)
			wasLiteral = false
		} else if token.stringValue == "(" {
			// push open parens onto the stack
			stack.push(token)
			wasLiteral = false
		} else if token.stringValue == ")" {
			// if we find a close paren, pop things off the stack
			for !stack.empty() && stack.peek().stringValue != "(" {
				queue.enqueue(stack.pop())
			}
			// there was an error parsing
			if stack.empty() {
				return nil, errors.New("parse error: mismatched parenthesis")
			}
			// pop off open paren
			stack.pop()
			// if next token is a function, move it to the queue
			if !stack.empty() {
				if _, ok := p.Functions[stack.peek().stringValue]; ok {
					queue.enqueue(stack.pop())
				}
			}
			wasLiteral = false
		} else {
			// if the last token was a literal it means we are trying to push 2 literals into the queue back to back
			// This will cause issues in the tree parsing. This is a rules violation and will throw an error
			if wasLiteral {
				return nil, errors.New("parse error: two literals found in a row")
			}
			// Token is a literal -- put it in the queue and set the bool to true
			queue.enqueue(token)
			wasLiteral = true
		}
	}

	// pop off the remaining operators onto the queue
	for !stack.empty() {
		if stack.peek().stringValue == "(" || stack.peek().stringValue == ")" {
			return nil, errors.New("parse error: mismatched parenthesis")
		}
		queue.enqueue(stack.pop())
	}

	return &queue, nil
}

// PostfixToTree Converts a Postfix token queue to a parse tree
func (p *Parser) postfixToTree(queue *tokenQueue) (*ParseNode, error) {
	stack := &nodeStack{}
	currNode := &ParseNode{}

	t := queue.Head
	for t != nil {
		t = t.Next
	}

	for !queue.empty() {
		// push the token onto the stack as a tree node
		currNode = &ParseNode{queue.dequeue(), nil, make([]*ParseNode, 0)}
		stack.push(currNode)

		if _, ok := p.Functions[stack.peek().Token.stringValue]; ok {
			// if the top of the stack is a function
			node, err := stack.pop()
			if err != nil {
				return nil, err
			}
			f := p.Functions[node.Token.stringValue]

			// pop off function parameters
			for i := 0; i < f.Params; i++ {
				childNode, childErr := stack.pop()
				if childErr != nil {
					return nil, childErr
				}
				// prepend children so they get added in the right order
				node.Children = append([]*ParseNode{childNode}, node.Children...)
			}

			if !checkChildType(node.Children) {
				return nil, errors.New("Cannot have literal and function/operator mismatch")
			}
			stack.push(node)
		} else if _, ok := p.Operators[stack.peek().Token.stringValue]; ok {
			// if the top of the stack is an operator
			node, err := stack.pop()
			if err != nil {
				return nil, err
			}
			o := p.Operators[node.Token.stringValue]

			// pop off operands
			for i := 0; i < o.Operands; i++ {
				// prepend children so they get added in the right order
				childNode, childErr := stack.pop()
				if childErr != nil {
					return nil, childErr
				}
				node.Children = append([]*ParseNode{childNode}, node.Children...)
			}
			if !checkChildType(node.Children) {
				return nil, errors.New("Cannot have literal and function/operator mismatch")
			}
			stack.push(node)
		}
	}

	return currNode, nil
}

// checkChildType Checks to make sure children types are compatible
func checkChildType(child []*ParseNode) bool {
	// Make sure we have 2 children in the tree
	if len(child) != 2 {
		return false
	}

	// Make sure that the token struct exists
	for _, c := range child {
		if c.Token == nil {
			return false
		}
	}

	// Children with the same type are compatible
	if child[0].Token.Type == child[1].Token.Type {
		return true
	}
	// If the first child is not an operator and function and
	// the second child is we have an invalid combination
	if (child[0].Token.Type != FilterTokenLogical &&
		child[0].Token.Type != FilterTokenFunc) &&
		(child[1].Token.Type == FilterTokenLogical ||
			child[1].Token.Type == FilterTokenFunc) {
		return false
	}
	// If the second child is not an operator and function and
	// the first child is we have an invalid combination
	if (child[0].Token.Type == FilterTokenLogical ||
		child[0].Token.Type == FilterTokenFunc) &&
		(child[1].Token.Type != FilterTokenLogical &&
			child[1].Token.Type != FilterTokenFunc) {
		return false
	}
	return true
}
