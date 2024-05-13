package handler

import (
	"log"
	"net/http"

	genericConstants "stock_broker_application/src/constants"
	"watchlist/business"
	"watchlist/commons/constants"

	"github.com/gin-gonic/gin"
)

type NewGetWatchlistController interface {
	HandleGetWatchlist(ctx *gin.Context)
}

type getWatchlistController struct {
	service business.NewGetWatchlistsService
}

func NewGetWatchListController(service business.NewGetWatchlistsService) NewGetWatchlistController {
	return &getWatchlistController{
		service: service,
	}
}

// @Summary Get the list of WatchLists
// @Description Handler function to fetch the user's watchlist data.
// @Produce json
// @Success 200 {object} map[string]interface{} "Returns the user's watchlist data"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security JWT
// @Router /v1/watchlist/list [get]
func (controller *getWatchlistController) HandleGetWatchlist(context *gin.Context) {

	watchlistData, err := controller.service.NewGetWatchlistsService(context)
	if len(watchlistData) == 0 {
		log.Println("watchlist not found")
		context.JSON(http.StatusNotFound, gin.H{genericConstants.GenericJSONErrorMessage: constants.WatchlistNotFoundError})
		return
	}

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{genericConstants.GenericJSONErrorMessage: err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{genericConstants.WatchlistTable: watchlistData})

}
