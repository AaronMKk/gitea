package middleware

import (
	"code.gitea.io/gitea/event"
	"code.gitea.io/gitea/modules/context"
	"code.gitea.io/gitea/modules/log"
)

// GitCloneMessageMiddleware middleware for send statistics message about git clone
func GitCloneMessageMiddleware(ctx *context.Context) {
	if err := event.GitCloneSender(ctx); err != nil {
		log.Fatal("internal error of kafka sender: %s", err.Error())
	}
}

// WebDownloadMiddleware middleware for send statistics message about git clone
func WebDownloadMiddleware(ctx *context.Context) {
	if err := event.WebDownloadSender(ctx); err != nil {
		log.Fatal("internal error of kafka sender: %s", err.Error())
	}
}
