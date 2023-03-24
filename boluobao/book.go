package boluobao

import (
	"github.com/VeronicaAlexia/BoluobaoAPI/boluobao/boluobaoapi"
	"github.com/VeronicaAlexia/BoluobaoAPI/pkg/threading"
	"github.com/tidwall/gjson"
	"strings"
)

type NovelInit struct {
	BookId string
}

func (book *NovelInit) NovelInfo() *gjson.Result {
	return boluobaoapi.NovelInformationAPI(book.BookId)

}

func (book *NovelInit) NovelCatalogue() []string {
	var chapter_id_list []string
	if response := boluobaoapi.NovelCatalogueAPI(book.BookId); response != nil {
		for _, volume := range response.Get("data.volumeList").Array() {
			for _, chapter := range volume.Get("chapterList").Array() {
				if chapter.Get("originNeedFireMoney").Int() == 0 {
					chapter_id_list = append(chapter_id_list, chapter.Get("chapId").String())
				}
			}
		}
		return chapter_id_list
	}
	return nil
}

func (book *NovelInit) NovelContent(chapter_id string) string {
	var contentText string
	if response := boluobaoapi.NovelContentAPI(chapter_id); response != nil {
		for _, content := range strings.Split(response.Get("data.expand.content").String(), "\n") {
			if content != "" {
				contentText += "　　" + content + "\n"
			}
		}
		return response.Get("data.title").String() + "\n\n" + contentText

	}
	return ""
}

func (book *NovelInit) NovelSearch(keyword string, lastPage int) []gjson.Result {
	var BookInfoList []gjson.Result
	Thread := threading.InitThreading(32)
	for i := 0; i < lastPage; i++ {
		Thread.Add()
		go func(page int) {
			defer Thread.Done()
			response := boluobaoapi.SearchAPI(keyword, page)
			if response != nil {
				//for _, data := range response.Data.Novels {
				for _, data := range response.Get("data.novels").Array() {
					BookInfoList = append(BookInfoList, data)
				}
			}
		}(i)
	}
	Thread.Wait() // wait for all goroutines to finish
	return BookInfoList
}