package api

import (
	"database/sql"
	"net/http"

	db "github.com/danglebary/beatstore-backend-go/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createLikeRequest struct {
	UserID int32 `json:"user_id" binding:"required,min=1"`
	BeatID int32 `json:"beat_id" binding:"required,min=1"`
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
	UserID int32 `uri:"uid" binding:"required,min=1"`
	BeatID int32 `uri:"bid" binding:"required,min=1"`
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

type listLikesByUserIDRequestUri struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type listLikesByUserIDRequestParams struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listLikesByUserID(ctx *gin.Context) {
	var uri listLikesByUserIDRequestUri
	var req listLikesByUserIDRequestParams

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListLikesByUserParams{
		UserID: uri.ID,
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

type listLikesByBeatIDRequestUri struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type listLikesByBeatIDRequestParams struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listLikesByBeatID(ctx *gin.Context) {
	var uri listLikesByBeatIDRequestUri
	var req listLikesByBeatIDRequestParams

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListLikesByBeatParams{
		BeatID: uri.ID,
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
