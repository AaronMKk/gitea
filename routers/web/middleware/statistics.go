package middleware

import (
	"code.gitea.io/gitea/event"
	"code.gitea.io/gitea/modules/context"
	"github.com/sirupsen/logrus"
)

// middleware for send statistics message about git clone
func GitCloneMessageMiddleware(ctx *context.Context) {
	logrus.Warnf("GitCloneMessageMiddleware Start...")
	if err := event.GitCloneSender(ctx); err != nil {
		logrus.Errorf("internal error of kafka sender: %s", err.Error())
	}
}
