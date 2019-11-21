# SM

SM(String Match)字符串匹配的一个框架。

# Core Data Structure

## Trie tree

数据结构：
```golang
type Node struct {
	IsKey    bool
	Children map[uint8]*Node
	Height   int
	Value    interface{}
}
```

如果字符集全部是英文，可以用**uint8**，但是如果保存中文字符集，一个中文字符会多创建3个空节点，此时使用**rune**更合适。但是大部分情况是中英文混合的情况。如果分离key的粒度是个问题。

基数树可以解决这个问题，但是实现复杂，插入性能低。

## Radix tree

## AC Automaton 

# Replication
