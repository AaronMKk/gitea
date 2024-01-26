package event

import (
	"strconv"
	"time"

	"code.gitea.io/gitea/event/infrastructure/kafka"
	"code.gitea.io/gitea/modules/context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Constants for topics, types, and field names.
const (
	// TODO change topic name to a more reasonable name, prob access_statistic
	topicGitClone = "gitea_git_clone"

	typGitClone    = "git_clone"
	typWebDownload = "web_download"

	fieldRepoName      = "repo_name"
	fieldRepoOwner     = "repo_owner"
	fieldDoerEmail     = "doer_email"
	fieldRepoId        = "repo_id"
	fieldRepoUserEmail = "repo_user_email"

	requestIp    = "request_IP"
	requestId    = "request_ID"
	requestProto = "request_proto"
	requestAgent = "http_user_agent"
	requestPath  = "path"
)

// getUserInfo extracts user information from the context.
func getUserInfo(ctx *context.Context) (string, string) {
	if ctx.Doer == nil {
		return "anno", "anno"
	}
	return ctx.Doer.Name, ctx.Doer.Email
}

// prepareMessage constructs the message to be sent.
func prepareMessage(ctx *context.Context, messageType string, description string) MsgNormal {
	doer, email := getUserInfo(ctx)

	// Generate a new UUID for the request.
	requestUuid := uuid.New().String()

	return MsgNormal{
		Type: messageType,
		User: doer,
		Desc: description,
		Details: map[string]string{
			fieldRepoName:      ctx.Repo.Repository.Name,
			fieldRepoOwner:     ctx.Repo.Repository.OwnerName,
			fieldDoerEmail:     email,
			fieldRepoId:        strconv.FormatInt(ctx.Repo.Repository.ID, 10),
			requestProto:       ctx.Req.Proto,
			requestIp:          ctx.Req.RemoteAddr,
			requestAgent:       ctx.Req.UserAgent(),
			fieldRepoUserEmail: ctx.ContextUser.Email,
			requestPath:        ctx.Repo.RepoLink,
			requestId:          requestUuid,
		},
		CreatedAt: time.Now().Unix(),
	}
}

// GitCloneSender handles sending Git clone events to Kafka.
func GitCloneSender(ctx *context.Context) error {
	msg := prepareMessage(ctx, typGitClone, "this message means someone cloned the repository")
	logrus.Infof("Publishing Git clone message: %#v\n", msg)
	return kafka.Publish(topicGitClone, &msg, nil)
}

// WebDownloadSender handles sending Web download events to Kafka.
func WebDownloadSender(ctx *context.Context) error {
	msg := prepareMessage(ctx, typWebDownload, "this message means someone downloaded a file from the website")
	logrus.Infof("Publishing Web download message: %#v\n", msg)
	return kafka.Publish(topicGitClone, &msg, nil)
}
