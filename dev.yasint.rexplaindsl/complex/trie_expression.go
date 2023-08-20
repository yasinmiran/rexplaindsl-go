package complex

import (
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"rexplaindsl/dev.yasint.rexplaindsl/api"
	"rexplaindsl/dev.yasint.rexplaindsl/util"
)

var nilKey string = ""
var nilNode node = node{
	nodes: treemap.NewWithStringComparator(),
}

// TrieExpression structure

type TrieExpression struct {
	root node
}

func NewTrieExpression() TrieExpression {
	return TrieExpression{
		// using string as the key instead runes
		root: node{treemap.NewWith(utils.StringComparator)},
	}
}

func (trie TrieExpression) Insert(word string) {
	var current = trie.root
	for i := 0; i < len(word); i++ {
		var c = string(word[i])
		if !current.containsKey(c) {
			_n := node{nodes: treemap.NewWith(utils.StringComparator)}
			current.put(c, _n)
		}
		current = current.get(c)
	}
	current.put(nilKey, nilNode) // nil as terminator
}

func (trie TrieExpression) InsertAll(words []string) {
	for _, v := range words {
		trie.Insert(v)
	}
}

func (trie TrieExpression) ToRegex() string {
	return *trie.root.ToRegex()
}

// Node construct

type node struct {
	// We use a forward reference of a node and
	// tree map returns a new pointer allocated
	// to this leaf node.
	nodes *treemap.Map
}

func (n node) containsKey(key string) bool {
	if _, ok := n.nodes.Get(key); ok {
		return true
	}
	return false
}

func (n node) put(key string, no node) {
	n.nodes.Put(key, no)
}

func (n node) get(key string) node {
	val, found := n.nodes.Get(key)
	if found {
		return val.(node)
	}
	panic("key not found!")
}

func (n node) size() int {
	return n.nodes.Size()
}

func (n node) ToRegex() *string {

	if n.containsKey(nilKey) && n.size() == 1 {
		return nil
	}

	var alternations []string
	var charClasses []string
	var hasOptionals = false
	var hasCharacterClasses = false

	n.nodes.Each(func(key interface{}, value interface{}) {
		escaped := util.AsRegexLiteral(key.(string))
		if value != nil {
			subExpression := value.(node).ToRegex()
			if subExpression != nil {
				alternations = append(alternations, escaped+*subExpression)
			} else {
				charClasses = append(charClasses, escaped)
			}
		} else {
			hasOptionals = true
		}
	})

	hasCharacterClasses = len(alternations) == 0

	if len(charClasses) > 0 {
		if len(charClasses) == 1 {
			alternations = append(alternations, charClasses[0])
		} else {
			setExpression := api.OPEN_SQUARE_BRACKET
			for _, v := range charClasses {
				setExpression += v
			}
			setExpression += api.CLOSE_SQUARE_BRACKET
			alternations = append(alternations, setExpression)
		}
	}

	var expression = "" // final expression

	if len(alternations) == 1 {
		expression += alternations[0]
	} else {
		expression += api.PAREN_OPEN + api.QUESTION_MARK + api.COLON
		for i, v := range alternations {
			expression += v
			if i != len(alternations)-1 {
				expression += api.ALTERNATION
			}
		}
		expression += api.PAREN_CLOSE
	}

	if hasOptionals {
		if hasCharacterClasses {
			expression += api.QUESTION_MARK
			return &expression
		} else {
			var temp = api.PAREN_OPEN + api.QUESTION_MARK + api.COLON +
				expression + api.PAREN_CLOSE + api.QUESTION_MARK
			return &temp
		}
	}

	return &expression

}
