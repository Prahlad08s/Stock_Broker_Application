package tests

import (
	"authentication/business"
	"authentication/handler"
	"authentication/models"
	"authentication/repositories"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ForgotPasswordSuite struct {
	suite.Suite
	mockRepo *repositories.MockForgotPasswordRepository
	service  business.ForgotPasswordService // Corrected type definition
	router   *gin.Engine
	ctx      *gin.Context
}

func (suite *ForgotPasswordSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()
	suite.mockRepo = new(repositories.MockForgotPasswordRepository)
	suite.service = business.NewForgotPasswordService(suite.mockRepo)

	// Setup Gin routes
	forgotPasswordMockHandler := handler.NewForgotPasswordController(suite.service)
	suite.router.POST("/v1/auth/forgot-password", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json")
		forgotPasswordMockHandler.HandleForgotPassword(ctx)
	})
}

func (suite *ForgotPasswordSuite) SetupTest() {
	// Create a context to use in the tests
	w := httptest.NewRecorder()
	suite.ctx, _ = gin.CreateTestContext(w)
}

func (suite *ForgotPasswordSuite) TestHandleForgotPassword_Success() {
	request := models.ForgotPasswordRequest{
		Email:         "test@example.com",
		PanCardNumber: "ABCDE1234F",
		NewPassword:   "NewSecurePassword1",
	}

	suite.mockRepo.On("VerifyAndUpdatePassword", request.Email, request.PanCardNumber, request.NewPassword).Return(nil)

	reqBody, _ := sonic.Marshal(request)
	suite.ctx.Request = httptest.NewRequest("POST", "/v1/auth/forgot-password", bytes.NewBuffer(reqBody))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, suite.ctx.Request)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.JSONEq(suite.T(), `{"message":"Password updated successfully"}`, w.Body.String())
}

func (suite *ForgotPasswordSuite) TestHandleForgotPassword_InternalServerError() {
	request := models.ForgotPasswordRequest{
		Email:         "test@example.com",
		PanCardNumber: "ABCDE1234F",
		NewPassword:   "NewSecurePassword1",
	}

	suite.mockRepo.On("VerifyAndUpdatePassword", request.Email, request.PanCardNumber, request.NewPassword).Return(assert.AnError)

	reqBody, _ := sonic.Marshal(request)
	suite.ctx.Request = httptest.NewRequest("POST", "/v1/auth/forgot-password", bytes.NewBuffer(reqBody))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, suite.ctx.Request)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.JSONEq(suite.T(), `{"error":"invalid user data."}`, w.Body.String())
}

func (suite *ForgotPasswordSuite) TestHandleForgotPassword_BadRequest() {
	request := models.ForgotPasswordRequest{
		Email: "test@example.com", // Missing required fields
	}

	reqBody, _ := sonic.Marshal(request)
	suite.ctx.Request = httptest.NewRequest("POST", "/v1/auth/forgot-password", bytes.NewBuffer(reqBody))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, suite.ctx.Request)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func TestForgotPasswordSuite(t *testing.T) {
	suite.Run(t, new(ForgotPasswordSuite))
}
