package user

import (
	"github.com/opensourceways/xihe-server/app"
)

const (
	errorSystemError         = "system_error"
	errorDuplicateCreating   = "duplicate_creating"
	errorResourceNotExists   = "resource_not_exists"
	errorConcurrentUpdating  = "concurrent_updateing"
	errorExccedMaxNum        = "exceed_max_num"
	errorUpdateLFSFile       = "update_lfs_file"
	errorPreviewLFSFile      = "preview_lfs_file"
	errorUnavailableRepoFile = "unavailable_repo_file"
)

// responseData is the response data to client
type responseData struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// ErrorDuplicateCreating
type ErrorDuplicateCreating struct {
	error
}

// ErrorConcurrentUpdating
type ErrorConcurrentUpdating struct {
	error
}

// ErrorResourceNotExists
type ErrorResourceNotExists struct {
	error
}

func newResponseCodeMsg(code, msg string) responseData {
	return responseData{
		Code: code,
		Msg:  msg,
	}
}

func newResponseError(err error) responseData {
	code := errorSystemError

	switch err.(type) {
	case ErrorDuplicateCreating:
		code = errorDuplicateCreating

	case ErrorResourceNotExists:
		code = errorResourceNotExists

	case ErrorConcurrentUpdating:
		code = errorConcurrentUpdating

	case app.ErrorExceedMaxRelatedResourceNum:
		code = errorExccedMaxNum

	case app.ErrorUpdateLFSFile:
		code = errorUpdateLFSFile

	case app.ErrorUnavailableRepoFile:
		code = errorUnavailableRepoFile

	case app.ErrorPreviewLFSFile:
		code = errorPreviewLFSFile

	default:

	}

	return responseData{
		Code: code,
		Msg:  err.Error(),
	}
}

func newResponseCodeError(code string, err error) responseData {
	return responseData{
		Code: code,
		Msg:  err.Error(),
	}
}
