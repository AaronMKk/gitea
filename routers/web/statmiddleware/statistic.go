package statmiddleware

import (
	"code.gitea.io/gitea/modules/context"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/messagequeue"
)

// GitCloneMessageMiddleware middleware for send statistics message about git clone
func GitCloneMessageMiddleware(ctx *context.Context) {
	if err := messagequeue.GitCloneSender(ctx); err != nil {
		log.Info("internal error of kafka sender: %s", err.Error())
	}
}

// WebDownloadMiddleware middleware for send statistics message about git clone
func WebDownloadMiddleware(ctx *context.Context) {
	if err := messagequeue.WebDownloadSender(ctx); err != nil {
		log.Info("internal error of kafka sender: %s", err.Error())
	}
}
