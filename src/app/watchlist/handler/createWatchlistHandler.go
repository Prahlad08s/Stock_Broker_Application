package handler

import (
	"encoding/json"
	genericConstants "stock_broker_application/src/constants"
	genericModel "stock_broker_application/src/models"
	"stock_broker_application/src/utils"
	"watchlist/business"
	serviceConstants "watchlist/commons/constants"
	serviceModel "watchlist/models"

	"github.com/gin-gonic/gin"
)

type CreateWatchlistController struct {
	service *business.CreateWatchlistService
}

func NewWatchlistController(service *business.CreateWatchlistService) *CreateWatchlistController {
	return &CreateWatchlistController{
		service: service,
	}
}

// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags CreateWatchlist
// @Accept json
// @Produce json
// @Security JWT
// @Param user body models.CreateWatchlist true "User details"
// @Success 201 {string} string "User created successfully"
// @Failure 400 {string} string "Bad request"
// @Router /v1/watchlist/create [post]
func (controller *CreateWatchlistController) HandleCreateWatchlist(ctx *gin.Context) {
	var createWatchlistRequest serviceModel.CreateWatchlist
	if err := ctx.BindJSON(&createWatchlistRequest); err != nil {
		errorMsgs := genericModel.ErrorMessage{Key: err.(*json.UnmarshalTypeError).Field, ErrorMessage: genericConstants.JsonBindingFieldError}
		utils.SendBadRequest(ctx, []genericModel.ErrorMessage{errorMsgs})
		return
	}
	if err := controller.service.CreateWatchlistService(createWatchlistRequest, ctx); err != nil {
		ErrMessage := genericModel.ErrorMessage{
			Key:          genericConstants.AccessToken,
			ErrorMessage: err.Error(),
		}
		utils.SendBadRequest(ctx, []genericModel.ErrorMessage{ErrMessage})
		ctx.Abort()
		return
	}
	utils.SendCreated(ctx, serviceConstants.WatchlistCreatedSuccessfully)

}
