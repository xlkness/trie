/**
 * @Author: likun
 * @Title: todo
 * @Description: todo
 * @File: trie_test
 * @Date: 2025-01-17 17:58:26
 */
package trie

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

var words = []struct {
	word     string
	sensType int
}{
	{"杀戮", 1},
	{"杀戮游戏", 1},
	{"bwi", 2},
	{"bwin平台", 3},
}

func TestTrie(t *testing.T) {
	tree := New()
	for _, word := range words {
		tree.Insert(word.word, word.sensType)
	}

	_, _, ok := tree.Filter("杀")
	assert.Equal(t, false, ok)

	newString, results, ok := tree.Filter("杀戮1")
	assert.Equal(t, true, ok)
	assert.Equal(t, "**1", newString)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, 1, results[0].SensType)
	assert.Equal(t, "杀戮", results[0].MatchWord)

	newString, results, ok = tree.Filter("杀戮游戏戏")
	assert.Equal(t, true, ok)
	assert.Equal(t, "****戏", newString)
	assert.Equal(t, 2, len(results))
	assert.Equal(t, 1, results[0].SensType)
	assert.Equal(t, "杀戮", results[0].MatchWord)
	assert.Equal(t, 1, results[1].SensType)
	assert.Equal(t, "杀戮游戏", results[1].MatchWord)

	newString, results, ok = tree.Filter("bwin")
	assert.Equal(t, true, ok)
	assert.Equal(t, "***n", newString)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, 2, results[0].SensType)
	assert.Equal(t, "bwi", results[0].MatchWord)

	newString, results, ok = tree.Filter("bwi")
	assert.Equal(t, true, ok)
	assert.Equal(t, "***", newString)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, 2, results[0].SensType)
	assert.Equal(t, "bwi", results[0].MatchWord)

	newString, results, ok = tree.Filter("bwin平台1")
	resultsBin, _ := json.Marshal(&results)
	assert.Equal(t, true, ok)
	assert.Equal(t, "******1", newString, string(resultsBin))
	assert.Equal(t, 2, len(results), 2)
	assert.Equal(t, 2, results[0].SensType, 2)
	assert.Equal(t, "bwi", results[0].MatchWord)
	assert.Equal(t, 3, results[1].SensType, 3)
	assert.Equal(t, "bwin平台", results[1].MatchWord)

	newString, results, ok = tree.Filter("bwin杀戮平台")
	assert.Equal(t, true, ok)
	assert.Equal(t, "***n**平台", newString)
	assert.Equal(t, 2, len(results))
	assert.Equal(t, 2, results[0].SensType)
	assert.Equal(t, "bwi", results[0].MatchWord)
	assert.Equal(t, 1, results[1].SensType)
	assert.Equal(t, "杀戮", results[1].MatchWord)
}
