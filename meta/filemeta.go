package meta

import (
	"fast-filestore-server/db"
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
func UpdateFileMeta(fMeta FileMeta) {
	fileMetas[fMeta.FileSha1] = fMeta
	log.Println("sha1:", fMeta.FileSha1)
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

//DB操作
func UpdateFileMetaDB(fMeta FileMeta) bool {
	return db.OnFileUploadFinished(fMeta.FileSha1, fMeta.FileName, fMeta.FileSize, fMeta.Location)
}

func GetFileMetaDB(fileSha1 string) (*FileMeta, error) {
	tFile, err := db.GetFileMeta(fileSha1)
	if tFile == nil || err != nil {
		return nil, err
	}
	fMeta := FileMeta{
		FileSha1: tFile.FileHash,
		FileName: tFile.FileName.String,
		FileSize: tFile.FileSize.Int64,
		Location: tFile.FileAddr.String,
	}
	return &fMeta, nil
}

func GetLastFileMetaDB(limit int) ([]FileMeta, error) {
	tFiles, err := db.GetFileMetaList(limit)
	if err != nil {
		return make([]FileMeta, 0), err
	}
	tFileMetas := make([]FileMeta, len(tFiles))
	for i := 0; i < len(tFileMetas); i++ {
		tFileMetas[i] = FileMeta{
			FileSha1: tFiles[i].FileHash,
			FileName: tFiles[i].FileName.String,
			FileSize: tFiles[i].FileSize.Int64,
			Location: tFiles[i].FileAddr.String,
		}
	}
	
	return tFileMetas, nil
}
