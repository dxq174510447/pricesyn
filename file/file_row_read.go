package file

import (
	"bufio"
	"os"
)

type FileRowRead struct {
	name string
	reader *bufio.Reader
}

func (f *FileRowRead) Parse(filename string) error {
	//file := "/Users/klook/code/company/klook-np2p-supplier/scripts/age_analyze/mapping.csv"
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	buf := bufio.NewReader(file)

	f.name = filename
	f.reader = buf
	return nil
}

func (f* FileRowRead) ScanCsvRow(fn func(int,[]string)) error{
	var row []byte
	var rowError error

	var result [][]string
	var rowNum int = 0
	for {
		row, _, rowError = f.reader.ReadLine()
		if row == nil || rowError != nil {
			break
		}
		rowStr := string(row)

		var colBegin int = 0
		var colEnd int = 0
		var rowResult []string
		for {
			if colEnd >= len(rowStr){
				// 结束
				break
			}
			// 列是否有双引号
			hasQuote := rowStr[colBegin:colBegin+1] == "\""
			quoteHasEnd := !hasQuote

			if hasQuote {
				colEnd = colEnd +1
			}

			var currentChar string
			var preChar string
			for {

				if colEnd >= len(rowStr) {
					//结束
					if hasQuote {
						rowResult = append(rowResult,rowStr[colBegin+1:colEnd-1])
					}else{
						rowResult = append(rowResult,rowStr[colBegin:colEnd])
					}
					break
				}

				preChar = currentChar
				currentChar = rowStr[colEnd:colEnd+1]
				if (!quoteHasEnd) && currentChar == "\"" && preChar != "\\" {
					quoteHasEnd = true
				}else if currentChar == "," && quoteHasEnd {
					if hasQuote {
						rowResult = append(rowResult,rowStr[colBegin+1:colEnd-1])
					}else{
						rowResult = append(rowResult,rowStr[colBegin:colEnd])
					}
					colBegin = colEnd+1
					colEnd = colBegin
					preChar = ""
					currentChar = ""
					break
				}
				colEnd = colEnd + 1
			}
		}
		result = append(result,rowResult)
		fn(rowNum,rowResult)
		rowNum++
	}
	return nil
}
