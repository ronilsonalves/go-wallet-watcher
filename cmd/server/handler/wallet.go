package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ronilsonalves/go-wallet-watcher/internal/wallet"
	"github.com/ronilsonalves/go-wallet-watcher/pkg/web"
	"net/http"
	"strconv"
)

type walletHandler struct {
	s wallet.Service
}

func NewWalletHandler(s wallet.Service) *walletHandler {
	return &walletHandler{s: s}
}

// GetWalletByAddress get wallet info balance from a given address
// @BasePath /api/v1
// GetWalletByAddress godoc
// @Summary Get a wallet info and balance from a wallet address
// @Schemes
// @Description Get wallet info from an address
// @Tags Wallets
// @Accept json
// @Produce json
// @Param address path string true "Wallet address"
// @Success 200 {object} wallet.DTO
// @Failure 400 {object} web.errResponse
// @Router /eth/wallets/{address} [GET]
func (h *walletHandler) GetWalletByAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ap := ctx.Param("address")
		w, err := h.s.GetWalletBalanceByAddress(ap)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}
		web.OKResponse(ctx, http.StatusOK, w)
	}
}

// GetTransactionsByAddress retrieves up to 10000 transactions by given adrress in a paggeable response
// @BasePath /api/v1
// GetTransactionsByAddress godoc
// @Summary Retrieves up to 10000 transactions by given adrress in a paggeable response
// @Schemes
// @Description Get wallet info from an address
// @Tags Wallets
// @Accept json
// @Produce json
// @Param address path string true "Wallet address"
// @Param page query string false "Page number"
// @Param pageSize query string false "Items per page"
// @Success 200 {object} web.PageableResponse
// @Failure 400 {object} web.errResponse
// @Router /eth/wallets/{address}/transactions [GET]
func (h *walletHandler) GetTransactionsByAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		address := ctx.Param("address")
		page := ctx.Query("page")
		size := ctx.Query("pageSize")

		if len(page) == 0 {
			page = "1"
		}
		if len(size) == 0 {
			size = "10"
		}

		if _, err := strconv.Atoi(page); err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", fmt.Sprintf("Invalid page param. Verify page value: %s", page))
			return
		}

		if _, err := strconv.Atoi(size); err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", fmt.Sprintf("Invalid pageSize param. Verify pageSize value: %s", size))
			return
		}

		response, err := h.s.GetTransactionsByAddress(address, page, size)
		if err != nil {
			web.BadResponse(ctx, http.StatusBadRequest, "error", err.Error())
			return
		}

		pageableResponse := web.PageableResponse{}

		pageableResponse.Page = page
		pageableResponse.Items = size
		pageableResponse.Data = response

		web.OKResponse(ctx, http.StatusOK, pageableResponse)
	}
}
