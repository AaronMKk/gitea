package event

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"

	"code.gitea.io/gitea/event/infrastructure/kafka"
	"code.gitea.io/gitea/modules/context"
)

// Constants for topics, types, and field names.
const (
	topicGitClone  = "gitea_git_clone"
	typGitClone    = "git_clone"
	typWebDownload = "web_download"
)

type Detail struct {
	RepoName      string `json:"repo_name"`
	RepoOwner     string `json:"repo_owner"`
	DoerEmail     string `json:"doer_email"`
	RepoId        string `json:"repo_id"`
	RepoUserEmail string `json:"repo_user_email"`
	RequestIP     string `json:"request_IP"`
	RequestID     string `json:"request_ID"`
	RequestProto  string `json:"request_proto"`
	RequestAgent  string `json:"http_user_agent"`
	RequestPath   string `json:"path"`
}

type MsgNormal struct {
	Type      string      `json:"type"`
	User      string      `json:"user"`
	Desc      string      `json:"desc"`
	Details   interface{} `json:"details"`
	CreatedAt int64       `json:"created_at"`
}

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

	detail := Detail{
		RepoName:      ctx.Repo.Repository.Name,
		RepoOwner:     ctx.Repo.Repository.OwnerName,
		DoerEmail:     email,
		RepoId:        strconv.FormatInt(ctx.Repo.Repository.ID, 10),
		RequestProto:  ctx.Req.Proto,
		RequestIP:     ctx.Req.RemoteAddr,
		RequestAgent:  ctx.Req.UserAgent(),
		RepoUserEmail: ctx.Repo.Repository.Owner.Email,
		RequestPath:   ctx.Repo.RepoLink,
		RequestID:     requestUuid,
	}

	return MsgNormal{
		Type:      messageType,
		User:      doer,
		Desc:      description,
		Details:   detail,
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
