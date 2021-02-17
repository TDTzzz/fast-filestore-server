package handler

import (
	"fast-filestore-server/common"
	"fast-filestore-server/meta"
	"fast-filestore-server/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

//Get:上传页面
func UploadHandler(c *gin.Context) {
	data, err := ioutil.ReadFile("./static/view/upload.html")
	if err != nil {
		c.String(404, "网页不存在")
		return
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, string(data))
}

//Post:处理文件上传
func DoUploadHandler(c *gin.Context) {
	errCode := 0
	defer func() {
		if errCode < 0 {
			c.JSON(http.StatusOK, gin.H{
				"Code": common.StatusErr,
				"Msg":  "Upload failed",
			})
		}
	}()

	//接受文件流并且存储到本地的目录
	file, head, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Printf("Failed to get data, err:%s\n", err.Error())
		errCode = -1
		return
	}
	defer file.Close()

	fileMeta := meta.FileMeta{
		FileName: head.Filename,
		Location: "/tmp/" + head.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		fmt.Printf("Failed to create file, err:%s\n", err.Error())
		errCode = -2
		return
	}
	defer newFile.Close()

	fileMeta.FileSize, err = io.Copy(newFile, file)
	if err != nil {
		fmt.Printf("Failed too save data into file, err:%s\n", err.Error())
		errCode = -3
		return
	}
	//如果不用这个 sha1值会不一样
	newFile.Seek(0, 0)
	fileMeta.FileSha1 = util.FileSha1(newFile)

	//游标重新回到文件头部
	//newFile.Seek(0, 0)
	//存储文件元信息
	meta.UpdateFileMeta(fileMeta)
	log.Println("success")
}

//查询文件 By SHA1
func GetFileBySha1(c *gin.Context) {
	sha1 := c.Param("sha1")

	fileMeta := meta.GetFileMeta(sha1)
	c.JSON(http.StatusOK, gin.H{
		"Code": common.StatusOK,
		"Msg":  "okk",
		"Data": fileMeta,
	})
}

//查询多个文件
func FileQueryHandler(c *gin.Context) {
	limitCnt, _ := strconv.Atoi(c.Param("limit"))

	fMetaArray := meta.GetLastFileMetas(limitCnt)
	c.JSON(http.StatusOK, gin.H{
		"Code": common.StatusOK,
		"Msg":  "okk",
		"Data": fMetaArray,
	})
}

//下载
func FileDownloadHandler(c *gin.Context) {
	errCode := common.StatusOK
	defer func() {

		if errCode == common.StatusOK {
			c.JSON(http.StatusOK, gin.H{
				"code": errCode,
				"msg":  "success",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": errCode,
				"msg":  "fail",
			})
		}
	}()
	fSha1 := c.Param("sha1")
	fm := meta.GetFileMeta(fSha1)

	f, err := os.Open(fm.Location)
	if err != nil {
		log.Println("file open err:", err)
		errCode = common.StatusErr
		return
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println("file read err:", err)
		errCode = common.StatusErr
		return
	}

	c.Header("content-disposition", "attachment; filename=\""+fm.FileName+"\"")
	c.Data(http.StatusOK, "application/octect-stream", data)
}

//文件重命名
func FileUpdateHandler(c *gin.Context) {
	fSha1 := c.Request.FormValue("sha1")
	newFileName := c.Request.FormValue("filename")

	curFileMeta := meta.GetFileMeta(fSha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	c.JSON(http.StatusOK, gin.H{
		"code": common.StatusOK,
		"msg":  "修改成功",
	})
}

func FileDeleteHandler(c *gin.Context) {
	fSha1 := c.Request.FormValue("sha1")

	//物理删除
	fMeta := meta.GetFileMeta(fSha1)
	os.Remove(fMeta.Location)

	meta.RemoveFileMeta(fSha1)

	c.JSON(http.StatusOK, gin.H{
		"code": common.StatusOK,
		"msg":  "删除成功",
	})
}
