/**
 * @Author: likun
 * @Title: todo
 * @Description: todo
 * @File: trie
 * @Date: 2025-01-17 17:52:53
 */
package trie

import (
	"bytes"
)

type node struct {
	no         int
	depth      int
	c          rune
	sensType   int // 敏感词类型，例如宗教信仰、政治、涉黄等
	isFullWord bool
	children   []*node
}

func (n *node) Insert(parentDepth int, reader *bytes.Reader, sensType int) {
	curRune, size, err := reader.ReadRune()
	if err != nil || size == 0 {
		n.isFullWord = true
		n.sensType = sensType
		return
	}

	var find bool
	for _, children := range n.children {
		if children.c == curRune {
			children.Insert(children.depth, reader, sensType)
			find = true
			break
		}
	}
	if !find {
		children := &node{no: len(n.children) + 1, depth: parentDepth + 1, c: curRune}
		n.children = append(n.children, children)
		children.Insert(children.depth, reader, sensType)
	}

}

func (n *node) FilterChildren(runeList []rune, records *filterRecords) {
	if n.isFullWord {
		// 找到当前字符就已经组成违禁词了，返回找到
		record := &FilterRecord{
			MatchWord:    records.matchContent,
			MatchRuneNum: records.matchRuneNum,
			SensType:     n.sensType,
		}
		records.records = append(records.records, record)
	}

	if len(runeList) == 0 {
		// 过滤内容被找完了还没找到完整的词，返回没找到
		return
	}

	curRune := runeList[0]

	for _, children := range n.children {
		if children.c == curRune {
			records.matchRuneNum += 1
			records.matchContent += string(curRune)
			children.FilterChildren(runeList[1:], records)
			continue
		}
	}

	return
}

type Tree struct {
	root *node
}

func New() *Tree {
	tree := &Tree{
		root: &node{c: -1},
	}
	return tree
}
func (t *Tree) Insert(s string, sensType int) {
	r := bytes.NewReader([]byte(s))
	t.root.Insert(-1, r, sensType)
}

func (t *Tree) Filter(s string) (string, []*FilterRecord, bool) {
	rawRuneList := make([]rune, 0, len(s)/3)
	replacedRuneList := make([]rune, 0, len(s)/3)

	for _, v := range s {
		rawRuneList = append(rawRuneList, v)
		replacedRuneList = append(replacedRuneList, v)
	}

	totalRecords := make([]*FilterRecord, 0)
	filterFlag := false

	for i := range rawRuneList {
		records := &filterRecords{}
		t.root.FilterChildren(rawRuneList[i:], records)
		if len(records.records) > 0 {
			filterFlag = true
			for _, record := range records.records {
				for j := i; j < i+record.MatchRuneNum; j++ {
					replacedRuneList[j] = '*'
				}
			}
			totalRecords = append(totalRecords, records.records...)
		}
	}

	return string(replacedRuneList), totalRecords, filterFlag
}

type FilterRecord struct {
	MatchWord    string
	MatchRuneNum int
	SensType     int
}

type filterRecords struct {
	records      []*FilterRecord
	matchContent string
	matchRuneNum int
}
