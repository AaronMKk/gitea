// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2020 The Gitea Authors.
// SPDX-License-Identifier: MIT

package user

import (
	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain/authing"
	userlogincli "github.com/opensourceways/xihe-server/user/infrastructure/logincli"
	"net/http"

	activities_model "code.gitea.io/gitea/models/activities"
	user_model "code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/modules/context"
	"code.gitea.io/gitea/routers/api/v1/utils"
	"code.gitea.io/gitea/services/convert"

	userapp "github.com/opensourceways/xihe-server/user/app"
	userrepo "github.com/opensourceways/xihe-server/user/domain/repository"
)

const (
	errorBadRequestBody   = "bad_request_body"
	errorBadRequestHeader = "bad_request_header"
)

type UserController struct {
	baseController

	repo     userrepo.User
	auth     authing.User
	s        userapp.UserService
	email    userapp.EmailService
	register userapp.RegService
}

// Search search users
func Search(ctx *context.APIContext) {
	// swagger:operation GET /users/search user userSearch
	// ---
	// summary: Search for users
	// produces:
	// - application/json
	// parameters:
	// - name: q
	//   in: query
	//   description: keyword
	//   type: string
	// - name: uid
	//   in: query
	//   description: ID of the user to search for
	//   type: integer
	//   format: int64
	// - name: page
	//   in: query
	//   description: page number of results to return (1-based)
	//   type: integer
	// - name: limit
	//   in: query
	//   description: page size of results
	//   type: integer
	// responses:
	//   "200":
	//     description: "SearchResults of a successful search"
	//     schema:
	//       type: object
	//       properties:
	//         ok:
	//           type: boolean
	//         data:
	//           type: array
	//           items:
	//             "$ref": "#/definitions/User"

	listOptions := utils.GetListOptions(ctx)

	uid := ctx.FormInt64("uid")
	var users []*user_model.User
	var maxResults int64
	var err error

	switch uid {
	case user_model.GhostUserID:
		maxResults = 1
		users = []*user_model.User{user_model.NewGhostUser()}
	case user_model.ActionsUserID:
		maxResults = 1
		users = []*user_model.User{user_model.NewActionsUser()}
	default:
		users, maxResults, err = user_model.SearchUsers(ctx, &user_model.SearchUserOptions{
			Actor:       ctx.Doer,
			Keyword:     ctx.FormTrim("q"),
			UID:         uid,
			Type:        user_model.UserTypeIndividual,
			ListOptions: listOptions,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]any{
				"ok":    false,
				"error": err.Error(),
			})
			return
		}
	}

	ctx.SetLinkHeader(int(maxResults), listOptions.PageSize)
	ctx.SetTotalCountHeader(maxResults)

	ctx.JSON(http.StatusOK, map[string]any{
		"ok":   true,
		"data": convert.ToUsers(ctx, ctx.Doer, users),
	})
}

// GetInfo get user's information
func GetInfo(ctx *context.APIContext) {
	// swagger:operation GET /users/{username} user userGet
	// ---
	// summary: Get a user
	// produces:
	// - application/json
	// parameters:
	// - name: username
	//   in: path
	//   description: username of user to get
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/User"
	//   "404":
	//     "$ref": "#/responses/notFound"

	if !user_model.IsUserVisibleToViewer(ctx, ctx.ContextUser, ctx.Doer) {
		// fake ErrUserNotExist error message to not leak information about existence
		ctx.NotFound("GetUserByName", user_model.ErrUserNotExist{Name: ctx.Params(":username")})
		return
	}
	ctx.JSON(http.StatusOK, convert.ToUser(ctx, ctx.ContextUser, ctx.Doer))
}

// GetAuthenticatedUser get current user's information
func GetAuthenticatedUser(ctx *context.APIContext) {
	// swagger:operation GET /user user userGetCurrent
	// ---
	// summary: Get the authenticated user
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/responses/User"

	ctx.JSON(http.StatusOK, convert.ToUser(ctx, ctx.Doer, ctx.Doer))
}

// GetUserHeatmapData is the handler to get a users heatmap
func GetUserHeatmapData(ctx *context.APIContext) {
	// swagger:operation GET /users/{username}/heatmap user userGetHeatmapData
	// ---
	// summary: Get a user's heatmap
	// produces:
	// - application/json
	// parameters:
	// - name: username
	//   in: path
	//   description: username of user to get
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/UserHeatmapData"
	//   "404":
	//     "$ref": "#/responses/notFound"

	heatmap, err := activities_model.GetUserHeatmapDataByUser(ctx, ctx.ContextUser, ctx.Doer)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "GetUserHeatmapDataByUser", err)
		return
	}
	ctx.JSON(http.StatusOK, heatmap)
}

func ListUserActivityFeeds(ctx *context.APIContext) {
	// swagger:operation GET /users/{username}/activities/feeds user userListActivityFeeds
	// ---
	// summary: List a user's activity feeds
	// produces:
	// - application/json
	// parameters:
	// - name: username
	//   in: path
	//   description: username of user
	//   type: string
	//   required: true
	// - name: only-performed-by
	//   in: query
	//   description: if true, only show actions performed by the requested user
	//   type: boolean
	// - name: date
	//   in: query
	//   description: the date of the activities to be found
	//   type: string
	//   format: date
	// - name: page
	//   in: query
	//   description: page number of results to return (1-based)
	//   type: integer
	// - name: limit
	//   in: query
	//   description: page size of results
	//   type: integer
	// responses:
	//   "200":
	//     "$ref": "#/responses/ActivityFeedsList"
	//   "404":
	//     "$ref": "#/responses/notFound"

	includePrivate := ctx.IsSigned && (ctx.Doer.IsAdmin || ctx.Doer.ID == ctx.ContextUser.ID)
	listOptions := utils.GetListOptions(ctx)

	opts := activities_model.GetFeedsOptions{
		RequestedUser:   ctx.ContextUser,
		Actor:           ctx.Doer,
		IncludePrivate:  includePrivate,
		OnlyPerformedBy: ctx.FormBool("only-performed-by"),
		Date:            ctx.FormString("date"),
		ListOptions:     listOptions,
	}

	feeds, count, err := activities_model.GetFeeds(ctx, opts)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "GetFeeds", err)
		return
	}
	ctx.SetTotalCountHeader(count)

	ctx.JSON(http.StatusOK, convert.ToActivities(ctx, feeds, ctx.Doer))
}

func (ctl baseController) checkUserApiToken(
	ctx *gin.Context, allowVistor bool,
) (
	pl oldUserTokenPayload, visitor bool, ok bool,
) {
	return ctl.checkUserApiTokenBase(ctx, allowVistor, true)
}

// @Summary		SendBindEmail
// @Description	send code to user
// @Tags			User
// @Accept			json
// @Success		201	{object}			app.UserDTO
// @Failure		500	system_error		system	error
// @Failure		500	duplicate_creating	create	user	repeatedly
// @Router			/v1/user/email/sendbind [post]
func (ctl *UserController) SendBindEmail(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "send code to user")

	req := EmailSend{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	cmd, err := req.toCmd(pl.DomainAccount())
	if err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	if code, err := ctl.email.SendBindEmail(&cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, "success")
	}
}

// @Summary		BindEmail
// @Description	bind email according the code
// @Tags			User
// @Param			body	body	userCreateRequest	true	"body of creating user"
// @Accept			json
// @Success		201	{object}			app.UserDTO
// @Failure		400	bad_request_body	can't	parse		request	body
// @Failure		400	bad_request_param	some	parameter	of		body	is	invalid
// @Failure		500	system_error		system	error
// @Failure		500	duplicate_creating	create	user	repeatedly
// @Router			/v1/user/email/bind [post]
func (ctl *UserController) BindEmail(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "bind email according to the code")

	req := EmailCode{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	cmd, err := req.toCmd(pl.DomainAccount())
	if err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	// create new token
	f := func() (token, csrftoken string) {
		user, err := ctl.s.GetByAccount(pl.DomainAccount())
		if err != nil {
			return
		}

		payload := oldUserTokenPayload{
			Account:                 user.Account,
			Email:                   user.Email,
			PlatformToken:           user.Platform.Token,
			PlatformUserNamespaceId: user.Platform.NamespaceId,
		}

		token, csrftoken, err = ctl.newApiToken(ctx, payload)
		if err != nil {
			return
		}

		return
	}

	if code, err := ctl.email.VerifyBindEmail(&cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		token, csrftoken := f()
		if token != "" {
			ctl.setRespToken(ctx, token, csrftoken, pl.Account)
		}

		ctl.sendRespOfPost(ctx, "success")
	}
}

func AddRouterForUserController(
	rg *gin.RouterGroup,
	us userapp.UserService,
	repo userrepo.User,
	auth authing.User,
	login app.LoginService,
	register userapp.RegService,
) {
	ctl := UserController{
		auth: auth,
		repo: repo,
		s:    us,
		email: userapp.NewEmailService(
			auth, userlogincli.NewLoginCli(login),
			us,
		),
		register: register,
	}

	// email
	rg.GET("/v1/user/check_email", checkUserEmailMiddleware(&ctl.baseController))
	rg.POST("/v1/user/email/sendbind", ctl.SendBindEmail)
	rg.POST("/v1/user/email/bind", ctl.BindEmail)
}
