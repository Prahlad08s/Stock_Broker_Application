package router

import (
	"authentication/business"
	serviceConstant "authentication/commons/constants"
	"authentication/docs"
	"authentication/handler"
	"authentication/repositories"
	"net/http"
	genericConstants "stock_broker_application/src/constants"
	"stock_broker_application/src/middleware/headerCheck"
	"stock_broker_application/src/utils/postgres"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func GetRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.New()
	router.Use(middlewares...)
	router.Use(gin.Recovery())
	connectionWithDb := postgres.GetPostGresClient().GormDb
	useDBMocks := true

	//dependency injection for signup
	userRepoSignUp := repositories.NewUserSignUpInstance()
	userSignUpService := business.NewSignUpService(userRepoSignUp)
	userSignUpController := handler.NewSignUpController(userSignUpService)
	//dependency injection for signin
	userDatabaseRepository := repositories.NewSignInRepository(connectionWithDb)
	signInService := business.NewSignInService(userDatabaseRepository)
	signInController := handler.NewSignInController(signInService)
	//Dependency Injection for OTP-Validation-Feature
	userRepository := repositories.NewUserRepository(connectionWithDb)
	otpService := business.NewOTPService(userRepository)
	otpValidationController := handler.NewOTPValidationController(otpService)
	//Dependency Injection for forgot-Password-Feature
	var forgotPasswordRepo repositories.ForgotPasswordRepository
	if useDBMocks {
		forgotPasswordRepo = repositories.NewForgotPasswordRepository(useDBMocks)
	} else {
		forgotPasswordRepo = repositories.NewForgotPasswordRepository(connectionWithDb)
	}
	forgotPasswordService := business.NewUsersService(forgotPasswordRepo)
	forgotPasswordController := handler.NewUsersController(forgotPasswordService)

	//Dependency Injection for change-password
	userDatabaseRepo := repositories.NewUserDBRepository(connectionWithDb)
	passwordService := business.NewChangePasswordService(userDatabaseRepo)
	changePasswordHandler := handler.NewChangePasswordController(passwordService)
	v1Routes := router.Group(genericConstants.RouterV1Config)
	{
		v1Routes.GET(serviceConstant.AuthenticationHealthCheck, func(c *gin.Context) {
			response := map[string]string{
				genericConstants.ResponseMessageKey: genericConstants.BFFResponseSuccessMessage,
			}
			c.JSON(http.StatusOK, response)
		})
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
		v1Routes.POST(serviceConstant.SignUp, userSignUpController.SignUp)
		v1Routes.POST(serviceConstant.SignIn, signInController.HandleSignIn)
		v1Routes.POST(serviceConstant.ValidateOTP, otpValidationController.HandleValidateOTP)
		v1Routes.PATCH(serviceConstant.ChangePassword, headerCheck.AuthMiddleware(), handler.HandleChangePassword(changePasswordHandler))
		v1Routes.POST(serviceConstant.ForgotPassword, forgotPasswordController.HandleForgotPassword)
		v1Routes.GET(serviceConstant.SwaggerRoute, ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	return router
}
