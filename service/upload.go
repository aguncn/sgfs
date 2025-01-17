package service

import (
	//"path"
	"strings"

	"github.com/LinkinStars/golang-util/gu"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/LinkinStars/sgfs/config"
	//"github.com/LinkinStars/sgfs/util/date_util"
)

func UploadFileHandler(ctx *fasthttp.RequestCtx) {
	// Get the file from the form
	header, err := ctx.FormFile("file")
	if err != nil {
		SendResponse(ctx, -1, "No file was found.", nil)
		return
	}

	// Check File Size
	if header.Size > int64(config.GlobalConfig.MaxUploadSize) {
		SendResponse(ctx, -1, "File size exceeds limit.", nil)
		return
	}

	// authentication
	token := string(ctx.FormValue("token"))
	if strings.Compare(token, config.GlobalConfig.OperationToken) != 0 {
		SendResponse(ctx, -1, "Token error.", nil)
		return
	}

	// Check upload File Path
	uploadSubPath := string(ctx.FormValue("uploadSubPath"))
	// 注释掉， 之前加了日期目录
	// visitPath := "/" + uploadSubPath + "/" + date_util.GetCurTimeFormat(date_util.YYYYMMDD)
	visitPath := "/" + uploadSubPath
	dirPath := config.GlobalConfig.UploadPath + visitPath
	if err := gu.CreateDirIfNotExist(dirPath); err != nil {
		zap.S().Error(err)
		SendResponse(ctx, -1, "Failed to create folder.", nil)
		return
	}
	// 注释掉，不要取什么后缀，直接取文件名
	// suffix := path.Ext(header.Filename)
	// filename := createFileName(suffix)
	filename := header.Filename

	fileAllPath := dirPath + "/" + filename

	/*
	   注释掉，有同名，就报错
	   // Guarantee that the filename does not duplicate
	   for {
	       if !gu.CheckPathIfNotExist(fileAllPath) {
	           break
	       }
	       filename = createFileName(suffix)
	       fileAllPath = dirPath + "/" + filename
	   }
	*/

	// Save file
	if err := fasthttp.SaveMultipartFile(header, fileAllPath); err != nil {
		zap.S().Error(err)
		SendResponse(ctx, -1, "Save file fail.", err.Error())
	}

	SendResponse(ctx, 1, "Save file success.", visitPath+"/"+filename)
	return
}

/*
注释掉，不需要重命名文件
func createFileName(suffix string) string {
    // Date and Time + _ + Random Number + File Suffix
    return date_util.GetCurTimeFormat(date_util.YYYYMMddHHmmss) + "_" + gu.GenerateRandomNumber(10) + suffix
}
*/
