package event

import (
	"time"

	"code.gitea.io/gitea/modules/context"
	"github.com/sirupsen/logrus"

	"code.gitea.io/gitea/event/infrastructure/kafka"
)

const (
	topicGitClone = "gitea_git_clone"
	topicGitLFS   = "gitea_git_lfs"

	typGitClone = "git_clone"

	fieldRepoName  = "repo_name"
	fieldRepoOwner = "repo_owner"
	fieldDoerEmail = "doer_email"
)

func GitCloneSender(ctx *context.Context) error {
	var doer string
	var email string

	if ctx.Doer == nil {
		doer = "anno"
		email = "anno"
	} else {
		doer = ctx.Doer.Name
		email = ctx.Doer.Email
	}

	msg := MsgNormal{
		Type: typGitClone,
		User: doer,
		Desc: "this message means someone clone the repository",
		Details: map[string]string{
			fieldRepoName:  ctx.Repo.Repository.Name,
			fieldRepoOwner: ctx.Repo.Repository.OwnerName,
			fieldDoerEmail: email,
		},
		CreatedAt: time.Now().Unix(),
	}

	logrus.Infof("pub msg %#v\n", msg)
	return kafka.Publish(topicGitClone, &msg, nil)
}

func GitLFSDownloadSender(msg *MsgNormal) error {
	return kafka.Publish(topicGitLFS, msg, nil)
}
