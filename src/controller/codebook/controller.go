package codebook

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	"xweb"

	"app"

	"github.com/labstack/echo"
)

type BookConfig struct {
	showTitle     bool `json:"showArticleTitle"`
	tabSpaceCount int  `json:"spaceCount"`
}

type Book struct {
	id        string
	name      string
	path      string
	config    *BookConfig
	catalogue *Node
}

type CodeBookController struct {
	_catalogueMgr CatalogueMgr
	_book         map[string]*Book
}

var codeBookController *CodeBookController

func NewCodeBookController() *CodeBookController {
	if codeBookController == nil {
		codeBookController = &CodeBookController{}
	}
	return codeBookController
}

func (self *CodeBookController) Call(params interface{}) map[string]interface{} {
	result := make(map[string]interface{}, 0)
	aritcles := make([]map[string]string, 0)

	codeBookId := params.(string)
	book := self.BookItemById(codeBookId)

	log.Println("get the book urls ", codeBookId)

	if book == nil || book.catalogue == nil {
		log.Println("can not found the book id ", codeBookId)

		result["code"] = 0
		result["data"] = nil
		return result
	}

	node := book.catalogue
	for _, node := range node.Children {
		aritcle := make(map[string]string, 0)
		aritcle["title"] = node.Name
		aritcle["url"] = node.Url
		aritcles = append(aritcles, aritcle)
	}

	result["code"] = 0
	result["data"] = aritcles

	return result
}

func (self *CodeBookController) Init(e *echo.Echo, params map[string]string) {
	self.parseParams(params)
	registerRoute(e)

	// load all codebook
	for _, v := range self._book {
		bookObj := v
		if bookObj != nil {
			node := self._catalogueMgr.loadCategory("", bookObj.id, bookObj.path+"/SUMMARY.md")
			if node == nil {
				log.Println("error! can not load category for ", bookObj.id)
			} else {
				log.Println("success! load category for ", bookObj.id)
			}
			bookObj.catalogue = node
		}
	}

	xweb.RegisterPluginCall("codebook", self)
}

func (self *CodeBookController) catalogueMgr() *CatalogueMgr {
	return &(self._catalogueMgr)
}

func (self *CodeBookController) parseParams(params map[string]string) {
	self._book = make(map[string]*Book, 0)

	bookKeys, ok := params["books"]
	if ok {
		bookKeyArr := strings.Split(bookKeys, ";")
		for _, bookKey := range bookKeyArr {
			if len(bookKey) == 0 {
				continue
			}

			log.Println("check book id :", bookKey)

			section := "MBook:" + bookKey
			sectionMap := app.WebConfigInst().GetSectionValues(section)

			bookPath, ok := sectionMap["path"]
			if ok {
				var bookObj Book
				bookObj.id = bookKey
				bookObj.path = bookPath

				bookName, okName := sectionMap["name"]
				if okName {
					bookObj.name = bookName
				}

				configPath := bookPath + "/book.json"
				bookJsonBin, err := ioutil.ReadFile(configPath)
				if err == nil && bookJsonBin != nil {
					var bookJson = BookConfig{}
					err := json.Unmarshal(bookJsonBin, &bookJson)
					if err == nil {
						bookObj.config = &bookJson
					}
				} else {
					log.Println("no book.json or read error", configPath)
				}

				log.Println("load success for book ", bookKey)
				self._book[bookKey] = &bookObj
			} else {
				log.Println("can not get the sectoin for ", bookKey)
			}
		}
	}
}

func (self *CodeBookController) BookCount() int {
	return len(self._book)
}

func (self *CodeBookController) BookPath(bookName string) *string {
	bookObj, ok := self._book[bookName]
	if ok {
		return &(bookObj.path)
	}
	return nil
}

func (self *CodeBookController) BookItemById(bookId string) *Book {
	bookObj, ok := self._book[bookId]
	if ok {
		return bookObj
	}
	return nil
}

func (self *CodeBookController) BookItemList() []*Book {
	bookList := make([]*Book, 0)
	for _, v := range self._book {
		bookList = append(bookList, v)
	}
	return bookList
}

func (self *CodeBookController) BookItemByIndex(index int) *Book {
	tempIndex := 0
	for _, v := range self._book {
		if tempIndex == index {
			return v
		}
		tempIndex++
	}
	return nil
}

func (self *CodeBookController) CatalogueFile(bookName string) (int, string) {
	path := self.BookPath(bookName)
	if path == nil {
		return 1, ""
	}
	return 0, *path + "/SUMMARY.md"
}

func (self *CodeBookController) GetCodeBookList() []map[string]string {
	bookIds := make([]string, 0)

	codeBookList := make([]map[string]string, 0)
	for _, v := range self._book {
		bookIds = append(bookIds, v.id)
	}

	sort.Sort(sort.StringSlice(bookIds))
	log.Println(bookIds)

	for _, id := range bookIds {
		v := self._book[id]
		codeBook := make(map[string]string, 0)
		codeBook["id"] = v.id
		codeBook["name"] = v.name
		codeBook["url"] = "/" + v.id

		codeBookList = append(codeBookList, codeBook)
	}

	return codeBookList
}
