package route

import (
	"errors"
	"messanger/db"
	"messanger/internal"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

var (
	ErrBanned         = errors.New("you are banned")
	ErrNoSuchUser     = errors.New("no user is found")
	ErrFriendNotFound = errors.New("your friend is not found")
)

type Handler struct {
	logger *zap.Logger
	dbConn *db.DB
	client *fasthttp.Client
	// TODO: add rabbitMQ to keep messages, which couldn't be delivered
}

func New(logger *zap.Logger, db *db.DB) *Handler {
	return &Handler{
		logger: logger,
		dbConn: db,
		client: &fasthttp.Client{},
	}
}

func (h *Handler) Send(reqCtx *fasthttp.RequestCtx) {
	var (
		msg     internal.Message
		accFrom db.Account
		accTo   db.Account
	)

	data := reqCtx.Request.Body()
	if err := jsoniter.Unmarshal(data, &msg); err != nil {
		reqCtx.SetStatusCode(fasthttp.StatusInternalServerError)
		_, _ = reqCtx.Write([]byte(err.Error()))
		return
	}

	if err := h.dbConn.FindAccountByUserName(msg.FromName, &accFrom); err != nil {
		reqCtx.SetStatusCode(fasthttp.StatusInternalServerError)
		_, _ = reqCtx.Write([]byte(err.Error()))
	}

	if accFrom.ID <= 0 {
		reqCtx.SetStatusCode(fasthttp.StatusBadRequest)
		_, _ = reqCtx.Write([]byte(ErrNoSuchUser.Error()))
		return
	}

	if accFrom.Banned {
		reqCtx.SetStatusCode(fasthttp.StatusBadRequest)
		_, _ = reqCtx.Write([]byte(ErrBanned.Error()))
		return
	}

	if err := h.dbConn.FindAccountByUserName(msg.ToName, &accTo); err != nil {
		reqCtx.SetStatusCode(fasthttp.StatusInternalServerError)
		_, _ = reqCtx.Write([]byte(err.Error()))
	}

	if accTo.ID <= 0 {
		reqCtx.SetStatusCode(fasthttp.StatusNoContent)
		_, _ = reqCtx.Write([]byte(ErrFriendNotFound.Error()))
		return
	}

	h.dbConn.UpdateAccountHost(accFrom.ID, msg.FromHost)

	req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetRequestURI("http://" + accTo.Host + "/get_msg")
	req.SetBody(data)

	if err := h.client.Do(req, resp); err != nil {
		reqCtx.SetStatusCode(resp.StatusCode())
		_, _ = reqCtx.Write([]byte(err.Error()))
		return
	}
}
