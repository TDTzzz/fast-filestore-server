package meta

import (
	"log"
	"sort"
)

//文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

//元信息相关
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
	log.Println("sha1:", fmeta.FileSha1)
}

//通过sha1值获取文件的元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

//获得最新的文件
func GetLastFileMetas(count int) []FileMeta {
	fMetaArray := make([]FileMeta, len(fileMetas))
	for _, v := range fileMetas {
		fMetaArray = append(fMetaArray, v)
	}
	//把所有的文件按UpdatedAt进行排序
	sort.Sort(ByUploadTime(fMetaArray))
	return fMetaArray[0:count]
}

//删除
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
