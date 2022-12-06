package readTools

import (
	//"fmt"
	"baliance.com/gooxml/document"
	"log"
)


//读取word文档
func ReadWord(filePath string) string {
	doc, err := document.Open(filePath)
	if err != nil {
		log.Fatalf("error opening document: %s", err)
	}

	doc.Paragraphs()
	re := ""
	for _, para := range doc.Paragraphs() {
		//run为每个段落相同格式的文字组成的片段
		for _, run := range para.Runs() {
			re += run.Text()
		}
		re += "\n"
	}
	return re
}
