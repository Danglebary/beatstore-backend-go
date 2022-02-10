package api

import (
	"database/sql"
	"net/http"

	db "github.com/danglebary/beatstore-backend-go/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createLikeRequest struct {
	UserID int32 `json:"user_id" binding:"required, min=1"`
	BeatID int32 `json:"beat_id" binding:"required, min=1"`
}

func (server *Server) createLike(ctx *gin.Context) {
	var req createLikeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateLikeParams{
		UserID: req.UserID,
		BeatID: req.BeatID,
	}
	like, err := server.store.CreateLike(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, like)
}

type getLikeRequest struct {
	UserID int32 `uri:"uid" binding:"required, min=1"`
	BeatID int32 `uri:"bid" binding:"required, min=1"`
}

func (server *Server) getLike(ctx *gin.Context) {
	var req getLikeRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.GetLikeByUserAndBeatParams{
		UserID: req.UserID,
		BeatID: req.BeatID,
	}
	like, err := server.store.GetLikeByUserAndBeat(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	ctx.JSON(http.StatusOK, like)
}

type listLikesByBeatIdRequestUri struct {
	ID int32 `uri:"uid" binding:"required, min=1"`
}

type listLikesByBeatIdRequestParams struct {
	PageID   int32 `form:"page_id" binding:"required, min=1"`
	PageSize int32 `form:"page_size" binding:"required, min=5, max=10"`
}

func (server *Server) listLikesByBeatId(ctx *gin.Context) {
	var bid listLikesByBeatIdRequestUri
	var req listLikesByBeatIdRequestParams
	if err := ctx.ShouldBindUri(&bid); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.ListLikesByBeatParams{
		BeatID: bid.ID,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	likes, err := server.store.ListLikesByBeat(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, likes)
}

type listLikesByUserIdRequestUri struct {
	ID int32 `uri:"uid" binding:"required, min=1"`
}

type listLikesByUserIdRequestParams struct {
	PageID   int32 `form:"page_id" binding:"required, min=1"`
	PageSize int32 `form:"page_size" binding:"required, min=5, max=10"`
}

func (server *Server) listLikesByUserId(ctx *gin.Context) {
	var uid listLikesByUserIdRequestUri
	var req listLikesByUserIdRequestParams
	if err := ctx.ShouldBindUri(&uid); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.ListLikesByUserParams{
		UserID: uid.ID,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	likes, err := server.store.ListLikesByUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, likes)
}
