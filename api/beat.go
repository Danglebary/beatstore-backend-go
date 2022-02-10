package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/danglebary/beatstore-backend-go/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createBeatRequestUri struct {
	ID int32 `json:"id" uri:"id" binding:"required, min=1"`
}

type createBeatRequestParams struct {
	Title string `json:"title"`
	Genre string `json:"genre"`
	Key   string `json:"key"`
	Bpm   int16  `json:"bpm"`
	Tags  string `json:"tags"`
}

func (server *Server) createBeat(ctx *gin.Context) {
	var uid createBeatRequestUri
	var req createBeatRequestParams
	if err := ctx.ShouldBindUri(&uid); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	s3Key := "not implemented"

	arg := db.CreateBeatParams{
		CreatorID: uid.ID,
		Title:     req.Title,
		Genre:     req.Genre,
		Key:       req.Key,
		Bpm:       req.Bpm,
		Tags:      req.Tags,
		S3Key:     s3Key,
	}
	beat, err := server.store.CreateBeat(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, beat)
}

type updateBeatRequestUri struct {
	ID int32 `json:"id" uri:"id" binding:"required, min=1"`
}

type updateBeatRequestParams struct {
	Title string `json:"title" binding:"required"`
	Genre string `json:"genre" binding:"required"`
	Key   string `json:"key" binding:"required"`
	Bpm   int16  `json:"bpm" binding:"required, min=20, max=999"`
	Tags  string `json:"tags"`
}

func (server *Server) updateBeat(ctx *gin.Context) {
	var bid updateBeatRequestUri
	var req updateBeatRequestParams
	if err := ctx.ShouldBindUri(&bid); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	s3Key := "not implemented"

	arg := db.UpdateBeatParams{
		ID:    bid.ID,
		Title: req.Title,
		Genre: req.Genre,
		Key:   req.Key,
		Bpm:   req.Bpm,
		Tags:  req.Tags,
		S3Key: s3Key,
	}

	beat, err := server.store.UpdateBeat(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, beat)
}

type getBeatByIdRequest struct {
	ID int32 `json:"id" uri:"id" binding:"required, min=1"`
}

func (server *Server) getBeat(ctx *gin.Context) {
	var bid getBeatByIdRequest
	if err := ctx.ShouldBindUri(&bid); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	beat, err := server.store.GetBeatById(ctx, bid.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, beat)
}

type listBeatsByIdRequest struct {
	PageID   int32  `form:"page_id" binding:"required, min=1"`
	PageSize int32  `form:"page_size" binding:"required, min=5, max=10"`
	Order    string `form:"order" binding:"required, oneof=ID BPM KEY GENRE"`
	BpmMin   int16  `form:"min" binding:"min=20, max=998"`
	BpmMax   int16  `form:"max" binding:"min=21, max=999"`
	Key      string `form:"key"`
	Genre    string `form:"genre"`
}

func (server *Server) listBeatsById(ctx *gin.Context) {
	var req listBeatsByIdRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	switch req.Order {
	case "ID":
		arg := db.ListBeatsByIdParams{
			Limit:  req.PageSize,
			Offset: (req.PageID - 1) * req.PageSize,
		}
		beats, err := server.store.ListBeatsById(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, beats)
	case "BPM":
		arg := db.ListBeatsByBpmRangeParams{
			Bpm:    req.BpmMin,
			Bpm_2:  req.BpmMax,
			Limit:  req.PageSize,
			Offset: (req.PageID - 1) * req.PageSize,
		}
		beats, err := server.store.ListBeatsByBpmRange(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, beats)
	case "kEY":
		arg := db.ListBeatsByKeyParams{
			Key:    req.Key,
			Limit:  req.PageSize,
			Offset: (req.PageID - 1) * req.PageSize,
		}
		beats, err := server.store.ListBeatsByKey(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, beats)
	case "GENRE":
		arg := db.ListBeatsByGenreParams{
			Genre:  req.Genre,
			Limit:  req.PageSize,
			Offset: (req.PageID - 1) * req.PageSize,
		}
		beats, err := server.store.ListBeatsByGenre(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, beats)
	default:
		err := fmt.Errorf("invalid 'order' query param value")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
}

type listBeatsByCreatorIdRequestUri struct {
	ID int32 `json:"id" uri:"id" binding:"required, min=1"`
}

type listBeatsByCreatorIdRequestParams struct {
	PageID   int32  `form:"page_id" binding:"required, min=1"`
	PageSize int32  `form:"page_size" binding:"required, min=5, max=10"`
	Order    string `form:"order" binding:"required, oneof=ID BPM KEY GENRE"`
	BpmMin   int16  `form:"min" binding:"min=20, max=998"`
	BpmMax   int16  `form:"max" binding:"min=21, max=999"`
	Key      string `form:"key"`
	Genre    string `form:"genre"`
}

func (server *Server) listBeatsByCreatorId(ctx *gin.Context) {
	var uid listBeatsByCreatorIdRequestUri
	var req listBeatsByCreatorIdRequestParams
	if err := ctx.ShouldBindUri(&uid); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	switch req.Order {
	case "ID":
		arg := db.ListBeatsByCreatorIdParams{
			CreatorID: uid.ID,
			Limit:     req.PageSize,
			Offset:    (req.PageID - 1) * req.PageSize,
		}
		beats, err := server.store.ListBeatsByCreatorId(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, beats)
	case "BPM":
		arg := db.ListBeatsByCreatorIdAndBpmRangeParams{
			CreatorID: uid.ID,
			Bpm:       req.BpmMin,
			Bpm_2:     req.BpmMax,
			Limit:     req.PageSize,
			Offset:    (req.PageID - 1) * req.PageSize,
		}
		beats, err := server.store.ListBeatsByCreatorIdAndBpmRange(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, beats)
	case "KEY":
		arg := db.ListBeatsByCreatorIdAndKeyParams{
			CreatorID: uid.ID,
			Key:       req.Key,
			Limit:     req.PageSize,
			Offset:    (req.PageID - 1) * req.PageSize,
		}
		beats, err := server.store.ListBeatsByCreatorIdAndKey(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, beats)
	case "GENRE":
		arg := db.ListBeatsByCreatorIdAndGenreParams{
			CreatorID: uid.ID,
			Genre:     req.Genre,
			Limit:     req.PageSize,
			Offset:    (req.PageID - 1) * req.PageSize,
		}
		beats, err := server.store.ListBeatsByCreatorIdAndGenre(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, beats)
	default:
		err := fmt.Errorf("invalid 'order' query param value")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
}
