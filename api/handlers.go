package api

import (
	"database/sql"
	"fmt"
	db "github.com/erkkke/technodom_test/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (server *Server) createRedirect(ctx *gin.Context) {
	var req db.CreateRedirectParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		fmt.Printf("binding err : %v", err.Error())
		return
	}

	redirect, err := server.db.CreateRedirect(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, redirect)
	return
}

type getRedirectRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getRedirect(ctx *gin.Context) {
	var req getRedirectRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	redirect, err := server.db.GetRedirect(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, redirect)
}

type listRedirectRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=100"`
}

func (server *Server) listRedirect(ctx *gin.Context) {
	var req listRedirectRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListRedirectsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	redirects, err := server.db.ListRedirects(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, redirects)
}

func (server *Server) removeRedirect(ctx *gin.Context) {
	var req getRedirectRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	rsp, err := server.db.DeleteRedirect(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) updateRedirect(ctx *gin.Context) {
	var req db.UpdateRedirectParams
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	rsp, err := server.db.UpdateRedirect(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) redirectHandler(ctx *gin.Context) {
	link := ctx.Query("link")

	v, ok := server.cache.Get(link)
	if ok {
		ctx.JSON(http.StatusMovedPermanently, v)
		return
	}

	_, err := server.db.GetRedirectByActiveLink(ctx, link)
	if err != nil {
		if err == sql.ErrNoRows {
			activeLink, err := server.db.GetRedirectByHistoryLink(ctx, link)
			if err != nil {
				if err == sql.ErrNoRows {
					ctx.JSON(http.StatusNotFound, gin.H{"error": "Redirect does not found"})
					return
				}
			}
			if server.cache.Len() < 1000 {
				server.cache.Add(link, activeLink)
			}

			ctx.JSON(http.StatusMovedPermanently, activeLink)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
