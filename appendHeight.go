package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
  // 结合前缀树 进行敏感词的标注
	flags := []string{"大", "大家", "大家好"} // 敏感词
	str := splitString("我是大家好abc爱大人", flags) // 进行的标注
	fmt.Println(str)
	// 我是<highlight>大家好</highlight>abc爱<highlight>大</highlight>人
}

type TrieNode struct {
	Next   map[string]*TrieNode
	Height int
	End    bool
}

type appendPos struct {
	front   bool
	backend bool
}

func buildTrieNode(flags []string) map[string]*TrieNode {
	trie := make(map[string]*TrieNode)
	first := trie
	for _, v := range flags {
		// 拆分词语
		runeV := []rune(v)
		maxL := len(runeV)
		for i := 0; i < maxL; i++ {
			trieVal := string(runeV[i])
			if trieRow, ok := trie[trieVal]; ok {
				trie = trieRow.Next
				continue
			}
			newTrieNode := &TrieNode{
				Next: make(map[string]*TrieNode),
				End:  false,
			}
			if i == maxL-1 {
				newTrieNode.Height = maxL
				newTrieNode.End = true
			}
			trie[trieVal] = newTrieNode
		}
	}
	return first
}

func splitString(str string, flags []string) string {
	if len(flags) == 0 {
		return str
	}
	trieTree := buildTrieNode(flags)
	//确定敏感词位置
	pos := trieTreeSearch(str, trieTree)
	// 敏感词位置去重
	uniquePos := uniqueSlice(pos)
	// 敏感词位置确定前后追加位置
	mapPos := confirmHeight(uniquePos)
	// 追加高亮标记
	return appendHeight(str, mapPos)
}

func trieTreeSearch(str string, trieTree map[string]*TrieNode) []int {
	runeS := []rune(str)
	maxL := len(runeS)
	firstCheckTrieTree := trieTree // 由于后面trieTree会递归，保留开始的位置
	continueTrieTree := []map[string]*TrieNode{}
	re := []int{} // 记录敏感词的位置
	for i := 0; i < maxL; i++ {
		checkWord := string(runeS[i])
		// 在前缀树中遍历
		checkTrieTree, ok := trieTree[checkWord]
		if !ok && len(continueTrieTree) == 0 {
			continue
		}
		if ok && checkTrieTree.End {
			// 根据深度 追加位置
			insertI := i
			for j := 1; j <= checkTrieTree.Height; j++ {
				re = append(re, insertI)
				insertI--
			}
			trieTree = firstCheckTrieTree // 回归
			// 放到continue中 为了后续的检查
			continueTrieTree = append(continueTrieTree, checkTrieTree.Next)
			continue
		}
		if ok && !checkTrieTree.End {
			continueTrieTree = append(continueTrieTree, checkTrieTree.Next)
		}
		//判断是否在continue中
		if len(continueTrieTree) > 0 {
			delPos := []int{}
			for m, tmpContinueCheck := range continueTrieTree {
				if tmpContinueCheckRow, ok := tmpContinueCheck[checkWord]; ok {
					if !tmpContinueCheckRow.End {
						// 剔除
						delPos = append(delPos, m)
					} else {
						// 根据深度 追加位置
						insertI := i
						for j := 1; j <= tmpContinueCheckRow.Height; j++ {
							re = append(re, insertI)
							insertI--
						}
					}
					continueTrieTree = append(continueTrieTree, tmpContinueCheckRow.Next)
				} else {
					// 剔除
					delPos = append(delPos, m)
				}
			}
			// 删除位置
			if len(delPos) > 0 {
				for j := len(delPos) - 1; j >= 0; j-- {
					trueDelPos := delPos[j]
					continueTrieTree = append(continueTrieTree[0:trueDelPos], continueTrieTree[trueDelPos+1:]...)
				}
			}
		}
	}
	return re
}

// 正整数去重
func uniqueSlice(nums []int) []int {
	// 排序
	sort.Ints(nums)
	newNums := []int{}
	pre := -1
	for _, num := range nums {
		if num != pre {
			newNums = append(newNums, num)
			pre = num
		}
	}
	return newNums
}

func confirmHeight(pos []int) map[int]appendPos {
	mapPos := make(map[int]appendPos)
	// 特殊情况 全部是连续的
	if pos[len(pos)-1]-pos[0] == len(pos)-1 {
		mapPos[pos[0]] = appendPos{
			front:   true,
			backend: false,
		}
		mapPos[pos[len(pos)-1]] = appendPos{
			front:   false,
			backend: true,
		}
	} else {
		for i := 0; i < len(pos); {
			if i == len(pos)-1 {
				if _, ok := mapPos[pos[i]]; !ok {
					mapPos[pos[i]] = appendPos{
						front:   true,
						backend: true,
					}
				}
				break
			}
			for j := i + 1; j < len(pos); {
				if pos[j]-pos[i] == j-i {
					j++
					continue
				} else {
					if pp, ok := mapPos[pos[i]]; ok {
						pp.front = true
					} else {
						mapPos[pos[i]] = appendPos{
							front:   true,
							backend: false,
						}
					}

					if pp, ok := mapPos[pos[j-1]]; ok {
						pp.backend = true
					} else {
						mapPos[pos[j-1]] = appendPos{
							front:   false,
							backend: true,
						}
					}
					i = j
					break
				}
			}
		}
	}
	return mapPos
}

// 在指定位置追加高亮
func appendHeight(str string, mapPos map[int]appendPos) string {
	strArr := []string{}
	runeStr := []rune(str)
	for k, v := range runeStr {
		if kMapPos, ok := mapPos[k]; ok {
			if kMapPos.front {
				strArr = append(strArr, "<highlight>")
			}
			strArr = append(strArr, string(v))
			if kMapPos.backend {
				strArr = append(strArr, "</highlight>")
			}
		} else {
			strArr = append(strArr, string(v))
		}
	}

	//fmt.Println(strArr)
	return strings.Join(strArr, "")
}
