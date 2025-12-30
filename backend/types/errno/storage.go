// backend/types/errno/storage.go
package errno

import "net/http"

// 存储模块错误码(STORAGE前缀)
// 遵循ZKER统一错误码定义规范

const (
	ErrInitStorageFailedCode = 61001
	ErrUploadFileFailedCode  = 61002
	ErrDownloadFileFailedCode = 61003
)

var (
	// 存储初始化错误 (500)
	ErrInitStorageFailed = &BaseErrorCode{
		code:       "STORAGE61001",
		message:    "Failed to initialize storage",
		messageZH:  "存储初始化失败",
		messageEN:  "Failed to initialize storage",
		httpStatus: http.StatusInternalServerError,
	}

	// 文件上传错误 (500)
	ErrUploadFileFailed = &BaseErrorCode{
		code:       "STORAGE61002",
		message:    "Failed to upload file",
		messageZH:  "文件上传失败",
		messageEN:  "Failed to upload file",
		httpStatus: http.StatusInternalServerError,
	}

	// 文件下载错误 (500)
	ErrDownloadFileFailed = &BaseErrorCode{
		code:       "STORAGE61003",
		message:    "Failed to download file",
		messageZH:  "文件下载失败",
		messageEN:  "Failed to download file",
		httpStatus: http.StatusInternalServerError,
	}
)
