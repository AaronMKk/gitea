package user

import (
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
	"github.com/opensourceways/xihe-server/utils"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	PrivateToken      = "PRIVATE-TOKEN"
	csrfToken         = "CSRF-Token"
	encodeUsername    = "encode-username"
	roleIndividuals   = "individuals"
	errorInvalidToken = "invalid_token"
)

var (
	respBadRequestBody = newResponseCodeMsg(
		errorBadRequestBody, "can't fetch request body",
	)
)

type baseController struct {
}

type accessController struct {
	RemoteAddr string      `json:"remote_addr"`
	Expiry     int64       `json:"expiry"`
	Role       string      `json:"role"`
	Payload    interface{} `json:"payload"`
}

func setCookie(ctx *gin.Context, key, val string, httpOnly bool, expireTime time.Time) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    val,
		Path:     "/",
		Expires:  expireTime,
		HttpOnly: httpOnly,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(ctx.Writer, cookie)
}

func (ctl baseController) checkUserApiTokenNoRefresh(
	ctx *gin.Context, allowVistor bool,
) (
	pl oldUserTokenPayload, visitor bool, ok bool,
) {
	return ctl.checkUserApiTokenBase(ctx, allowVistor, false)
}

func (ctl baseController) sendRespWithInternalError(ctx *gin.Context, data responseData) {
	private_log.Errorf("code: %s, err: %s", data.Code, data.Msg)

	ctx.JSON(http.StatusInternalServerError, data)
}

func (ctl baseController) sendCodeMessage(ctx *gin.Context, code string, err error) {
	if code == "" {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(code, err))
	}
}

func getCookieValue(ctx *gin.Context, key string) (string, error) {
	cookie, err := ctx.Request.Cookie(key)
	if err != nil {
		return "", nil
	}

	return cookie.Value, nil
}

func (ctl baseController) newRepo() repositories.Access {
	return repositories.NewAccessRepo(int(apiConfig.TokenExpiry - 10))
}

func (ctl *baseController) getCookieToken(ctx *gin.Context) (token string, err error) { // TODO add return value exist
	// get encode username
	u, err := getCookieValue(ctx, PrivateToken)
	if err != nil {
		return "", nil
	}

	// insert encode username to context
	ctx.Set(encodeUsername, u)

	// get token from redis
	token, err = ctl.newRepo().Get(u)
	if err != nil {
		return "", nil
	}

	return
}

func (ctl *baseController) getCSRFToken(ctx *gin.Context) (string, error) {
	v := ctx.GetHeader(csrfToken)
	return v, nil
}

func (ctl baseController) decryptData(s string) ([]byte, error) {
	dst, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return encryptHelper.Decrypt(dst)
}

func (ctl baseController) checkToken(
	ctx *gin.Context, token string, pl interface{},
) (ac accessController, tokenbyte []byte, ok bool) {
	tokenbyte, err := ctl.decryptData(token)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			newResponseCodeError(errorSystemError, err),
		)

		return
	}

	ac = accessController{
		Payload: pl,
	}

	if err := ac.initByToken(string(tokenbyte), apiConfig.TokenKey); err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			newResponseCodeError(errorSystemError, err),
		)

		return
	}

	if err = ac.verify([]string{roleIndividuals}); err != nil {
		ctx.JSON(
			http.StatusUnauthorized,
			newResponseCodeError(errorInvalidToken, err),
		)
		return
	}

	ok = true

	return
}

func (ctl baseController) decryptDataForCSRF(s string) ([]byte, error) {
	dst, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return encryptHelperCSRF.Decrypt(dst)
}

func (ctl baseController) checkCSRFToken(
	ctx *gin.Context, tokenbyte []byte, csrftoken string,
) (ok bool) {
	csrfbyte, err := ctl.decryptDataForCSRF(csrftoken)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			newResponseCodeError(errorSystemError, err),
		)

		return
	}

	if !verifyCSRFToken(tokenbyte, csrfbyte) {
		ctx.JSON(
			http.StatusUnauthorized,
			newResponseCodeError(errorInvalidToken, errors.New("token not allowed")),
		)

		return
	}

	ok = true

	return
}

// crypt for token
func (ctl baseController) encryptData(d string) (string, error) {
	t, err := encryptHelper.Encrypt([]byte(d))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(t), nil
}

// crypt for csrftoken
func (ctl baseController) encryptDataForCSRF(d string) (string, error) {
	t, err := encryptHelperCSRF.Encrypt([]byte(d))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(t), nil
}

func (ctl baseController) refreshDoubleToken(ac accessController) (token, csrftoken string) {
	decode, err := ac.refreshToken(apiConfig.TokenExpiry, apiConfig.TokenKey)
	if err == nil {
		if t, err := ctl.encryptData(decode); err == nil {
			token = t
		}
	}

	if w := ac.genCSRFToken(decode); len(w) != 0 {
		if t, err := ctl.encryptDataForCSRF(w); err == nil {
			csrftoken = t
		}
	}

	return
}

func (ctl baseController) setRespCookieToken(ctx *gin.Context, token, username string) error {
	// encrpt username
	u, err := ctl.encryptData(username)
	if err != nil {
		return err
	}

	// insert redis
	if err = ctl.newRepo().Insert(u, token); err != nil {
		return err
	}

	// set expire time for old token
	var ok, oldusername = false, ""
	o, exist := ctx.Get(encodeUsername)
	if exist {
		if oldusername, ok = o.(string); !ok {
			return errors.New("encode username illegal")
		}
	}

	if err = ctl.newRepo().Expire(oldusername, 3); err != nil {
		return err
	}

	// set cookie
	setCookie(ctx, PrivateToken, u, true, utils.ExpiryReduceSecond(apiConfig.TokenExpiry))

	return nil
}

func (ctl baseController) setRespCSRFToken(ctx *gin.Context, token string) {
	setCookie(ctx, csrfToken, token, false, utils.ExpiryReduceSecond(apiConfig.TokenExpiry))
}

func (ctl baseController) setRespToken(ctx *gin.Context, token, csrftoken, username string) error {
	ctl.setRespCSRFToken(ctx, csrftoken)

	if err := ctl.setRespCookieToken(ctx, token, username); err != nil {
		return err
	}

	return nil
}

func (ctl baseController) checkApiToken(
	ctx *gin.Context, token string, csrftoken string, pl interface{}, refresh bool,
) (ok bool) {
	ac, tokenbyte, ok := ctl.checkToken(ctx, token, pl)
	if !ok {
		return
	}

	if ok = ctl.checkCSRFToken(ctx, tokenbyte, csrftoken); !ok {
		return
	}

	if !refresh {
		return
	}

	token, csrftoken = ctl.refreshDoubleToken(ac)

	payload, good := ac.Payload.(*oldUserTokenPayload)
	if !good {
		return good
	}

	if err := ctl.setRespToken(ctx, token, csrftoken, payload.Account); err != nil {
		logrus.Debugf("set resp token error: %s", err.Error())
	}

	return
}

func (ctl baseController) checkUserApiTokenBase(
	ctx *gin.Context, allowVistor bool, refresh bool,
) (
	pl oldUserTokenPayload, visitor bool, ok bool,
) {
	token, err := ctl.getCookieToken(ctx)
	if err != nil {
		return
	}

	csrftoken, err := ctl.getCSRFToken(ctx)
	if err != nil {
		return
	}

	if token == "" || csrftoken == "" {
		// try best to grab userinfo
		if token != "" {
			_, _, _ = ctl.checkToken(ctx, token, &pl)
		}
		if allowVistor {
			visitor = true
			ok = true
		} else {
			ctx.JSON(
				http.StatusUnauthorized,
				newResponseCodeMsg(errorBadRequestHeader, "no token"),
			)
		}

		return
	}

	ok = ctl.checkApiToken(ctx, token, csrftoken, &pl, refresh)

	return
}

func (ctl baseController) sendBadRequestBody(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, respBadRequestBody)
}

func newResponseData(data interface{}) responseData {
	return responseData{
		Data: data,
	}
}

func (ctl baseController) sendRespOfPost(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusCreated, newResponseData(data))
}

func (ctl baseController) getRemoteAddr(ctx *gin.Context) (string, error) {
	ips := ctx.Request.Header.Get("x-forwarded-for")

	for _, item := range strings.Split(ips, ", ") {
		if net.ParseIP(item) != nil {
			return item, nil
		}
	}

	return "", errors.New("can not fetch client ip")
}

func (ctl baseController) newApiToken(ctx *gin.Context, pl interface{}) (
	token string, csrftoken string, err error,
) {
	addr, err := ctl.getRemoteAddr(ctx)
	if err != nil {
		return "", "", err
	}

	ac := &accessController{
		Expiry:     utils.Expiry(apiConfig.TokenExpiry),
		Role:       roleIndividuals,
		Payload:    pl,
		RemoteAddr: addr,
	}

	t, err := ac.newToken(apiConfig.TokenKey)
	if err != nil {
		return "", "", err
	}

	ct := ac.genCSRFToken(t)

	if token, err = ctl.encryptData(t); err != nil {
		return
	}

	if csrftoken, err = ctl.encryptDataForCSRF(ct); err != nil {
		return
	}

	return
}
