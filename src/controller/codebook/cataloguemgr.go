package codebook

import (
	"io/ioutil"
	"log"
	"strings"
)

/* 目录管理器：缓存了所有codebook的目录 */

type Node struct {
	Name     string
	HtmlName string
	Url      string // 跳转地址
	Level    int    // 结点等级
	Type     int    // 结点类型,0:根结点，虚结点，不展示。1:分区，以#开始，没有url,一般用在顶层，用做section区分，他与后面跟进的内容，并不是父子关系（展示上没有收起，展开的逻辑）。2：以*开始，没有url,表示中间章节，仅用来表示父子关系 3：以*开始，有url，表示某一具体的节

	ParentNode *Node
	Children   []*Node
}

type CatalogueMgr struct {
	rootNodeList map[string]*Node
}

func (self *CatalogueMgr) initNode() *Node {
	rootNode := new(Node)
	rootNode.Name = "ROOT"
	rootNode.Url = ""
	rootNode.Level = -1
	rootNode.Type = 0
	rootNode.Children = make([]*Node, 0)
	rootNode.ParentNode = nil
	return rootNode
}

func (self *CatalogueMgr) loadCategory(prefix, name, file string) *Node {
	if self.rootNodeList == nil {
		self.rootNodeList = make(map[string]*Node, 0)
	}

	contentbin, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("read file error, can not load the file's SUMMARY:", file)
		return nil
	}

	prefixUrl := ""
	if len(prefix) > 0 {
		prefixUrl = "/" + prefix + "/" + name
	} else {
		prefixUrl = "/" + name
	}

	//initRootNode()
	rootNode := self.initNode()
	self.rootNodeList[name] = rootNode

	content := string(contentbin)
	lines := strings.Split(content, "\n")
	if lines == nil || len(lines) == 0 {
		log.Println("the content is empty")
		return nil
	}

	currentNode := rootNode
	for _, line := range lines {
		tempLine := strings.Trim(line, " ")
		if len(tempLine) == 0 {
			continue
		}
		lineNode := self.processLine(prefixUrl, rootNode, currentNode, line)
		if lineNode != nil {
			currentNode = lineNode
		}
	}

	return rootNode
}

func getSpaceCount(line string) int {
	count := 0
	len := len(line)
	for i := 0; i < len; i++ {
		if line[i] == ' ' {
			count++
		} else {
			break
		}
	}
	return count
}

func setParentNode(current, newNode, rootNode *Node) {
	if newNode.Level == current.Level {
		newNode.ParentNode = current.ParentNode
		newNode.ParentNode.Children = append(newNode.ParentNode.Children, newNode)
	} else if newNode.Level > current.Level {
		newNode.ParentNode = current
		newNode.ParentNode.Children = append(newNode.ParentNode.Children, newNode)
	} else if newNode.Level < current.Level {
		// 一直回朔到根
		for {
			if newNode.Level < current.Level || current == nil {
				current = current.ParentNode
			} else {
				break
			}
		}

		if current != nil {
			newNode.ParentNode = current.ParentNode
		}

		if newNode.ParentNode == nil {
			newNode.ParentNode = rootNode
		}

		newNode.ParentNode.Children = append(newNode.ParentNode.Children, newNode)
	}
}

func (self *CatalogueMgr) parseNode(line, prefixUrl string, node *Node) {
	pos := strings.Index(line, "[")
	posend := strings.Index(line, "](")
	if pos >= 0 && posend > 0 && pos < posend {
		node.Name = line[pos+1 : posend]
		node.HtmlName = node.Name
	} else {
		log.Println("tt", pos, posend, line)
	}

	pos = strings.Index(line, "](")
	posend = strings.LastIndex(line, ")")
	if pos > 0 && posend > 0 && pos < posend {
		node.Type = 3
		node.Url = line[pos+2 : posend]
		if !strings.HasPrefix(node.Url, "/") {
			node.Url = prefixUrl + "/" + node.Url
		}

		// .md -> .html
		log.Println("url :" + node.Url)
		if strings.HasSuffix(node.Url, ".md") {
			node.Url = node.Url[0 : len(node.Url)-3]
			node.Url = node.Url + ".html"
		}
	}

	if len(node.Name) == 0 {
		node.Name = line
		node.HtmlName = node.Name
	}
}

// lastNode 为上一行解析后的node,即当前node
func (self *CatalogueMgr) processLine(prefixUrl string, rootNode, current *Node, line string) *Node {
	tempNode := new(Node)
	tempNode.Children = make([]*Node, 0)
	tempNode.Level = getSpaceCount(line) / 4

	line = strings.TrimSpace(line)

	if strings.HasPrefix(line, "#") {
		line := strings.Trim(line, "#")
		self.parseNode(line, prefixUrl, tempNode)
		if len(tempNode.Url) == 0 {
			tempNode.Type = 1
		} else {
			tempNode.Type = 3
		}
	} else if strings.HasPrefix(line, "*") {
		line := strings.Trim(line, "*")
		self.parseNode(line, prefixUrl, tempNode)
		if len(tempNode.Url) == 0 {
			tempNode.Type = 2
		} else {
			tempNode.Type = 3
		}
	} else {
		// error
		tempNode = nil
		log.Println("error, this format is not support", line)
	}

	if tempNode != nil {
		setParentNode(current, tempNode, rootNode)
	}

	return tempNode
}

func (self *CatalogueMgr) GetRootNode(name string) *Node {
	v, ok := self.rootNodeList[name]
	if ok {
		return v
	}
	return nil
}

func (self *CatalogueMgr) renderRootNode(name, page string) string {
	v, ok := self.rootNodeList[name]
	if ok {
		log.Println("find the root node, will render : ", name, page)
		return renderRootNodeImpl(name, page, v)
	}
	log.Println("can not find the root node")
	return ""
}

func renderRootNodeImpl(name, page string, node *Node) string {
	result := ""

	if node.ParentNode != nil {
		hasChild := len(node.Children) > 0
		if node.Type == 1 {
			result = result + "<li class=\"chapter\"><span class=\"cb_section\">" + node.HtmlName + "</span></li>"
		} else if node.Type == 2 {
			if hasChild {
				result = result + "<li class=\"chapter\"><span class=\"cb_section_2\"><span class=\"glyphicon glyphicon-menu-down\" aria-hidden=\"true\"></span>&nbsp;" + node.HtmlName + "</span></li>"
			} else {
				result = result + "<li class=\"chapter\"><span class=\"cb_section_2\">" + node.HtmlName + "</span></li>"
			}
		} else if node.Type == 3 {
			urlActive := strings.Compare(node.Url, page) == 0
			urlClass := "chapterselected"
			if !urlActive {
				urlClass = ""
			}
			if hasChild {
				result = result + "<li class=\"chapter " + urlClass + "\"><a href=\"" + node.Url + "\"><span class=\"glyphicon glyphicon-menu-down\" aria-hidden=\"true\"></span>&nbsp;" + node.HtmlName + "</a></li>"
			} else {
				result = result + "<li class=\"chapter " + urlClass + "\"><a href=\"" + node.Url + "\">" + node.HtmlName + "</a></li>"
			}
		} else {
		}
	}

	if len(node.Children) == 0 {
		return result
	}

	result = result + "<ul class=\"articles\">"
	for i := 0; i < len(node.Children); i++ {
		result = result + renderRootNodeImpl(name, page, node.Children[i])
	}
	result = result + "</ul>"

	return result
}

func buildNode() *Node {
	return nil
}
