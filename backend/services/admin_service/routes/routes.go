package routes

import (
	"github.com/dath-241/coin-price-be-go/services/admin_service/config"
	"github.com/dath-241/coin-price-be-go/services/admin_service/controllers"
	"github.com/dath-241/coin-price-be-go/services/admin_service/middlewares"
	"github.com/dath-241/coin-price-be-go/services/admin_service/momo"
	"github.com/dath-241/coin-price-be-go/services/admin_service/repository"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {

	// Kết nối tới collection "User"
	collection := config.DB.Collection("User")

	// Tạo instance repository
	repo := &repository.MongoUserRepository{Collection: collection}

	// Kết nối tới collection "User"
	collectionpay := config.DB.Collection("OrderMoMo")

	// Tạo instance repository
	payrepo := &repository.MongoPaymentRepository{Collection: collectionpay}

	// Route User Management API (Quản lý tài khoản)
	userRoutes := r.Group("/api/v1/user")
	{
		// Áp dụng AuthMiddleware cho tất cả các route trong group này
		userRoutes.Use(middlewares.AuthMiddleware("VIP-0", "VIP-1", "VIP-2", "VIP-3", "Admin"))

		userRoutes.GET("/me", controllers.GetCurrentUserInfo(repo))   // Lấy thông tin tài khoản người dùng hiện tại.
		userRoutes.PUT("/me", controllers.UpdateUserProfile(repo))    // Chỉnh sửa thông tin tài khoản người dùng.
		userRoutes.DELETE("/me", controllers.DeleteCurrentUser(repo)) // Xóa tài khoản người dùng hiện tại.
		userRoutes.PUT("/me/change_password", controllers.ChangePassword(repo))
		userRoutes.PUT("/me/change_email", controllers.ChangeEmail(repo))
		userRoutes.GET("/me/payment-history", controllers.GetPaymentHistory(payrepo))
	}

	// Route cho xác thực
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/login", controllers.Login(repo)) // Người dùng đăng nhập
		authRoutes.POST("/google-login", controllers.GoogleLogin(repo))
		authRoutes.POST("/register", controllers.Register(repo))                                                                      // Người dùng đăng kí
		authRoutes.POST("/forgot-password", controllers.ForgotPassword(repo))                                                         // Người dùng quên mật khẩu
		authRoutes.POST("/reset-password", controllers.ResetPassword(repo))                                                           // Người dùng thay đổi mật khẩu
		authRoutes.POST("/refresh-token", controllers.RefreshToken())                                                             // Người dùng refresh token
		authRoutes.POST("/logout", middlewares.AuthMiddleware("VIP-0", "VIP-1", "VIP-2", "VIP-3", "Admin"), controllers.Logout()) // Người dùng đăng xuất
	}

	// Route cho quản trị viên (admin)
	adminRoutes := r.Group("/api/v1/admin/")
	{
		// Áp dụng AuthMiddleware cho tất cả các route trong group này
		adminRoutes.Use(middlewares.AuthMiddleware("Admin"))

		adminRoutes.GET("/users", controllers.GetAllUsers(repo))                  //Lay thong tin tat ca nguoi dung
		adminRoutes.GET("/user/:user_id", controllers.GetUserByAdmin(repo))       //Lay thong tin 1 user cu the
		adminRoutes.DELETE("/user/:user_id", controllers.DeleteUserByAdmin(repo)) //Xóa 1 người dùng
		adminRoutes.PUT("/user/:user_id/ban", controllers.BanAccount(repo)) //Ban 1 người dùn
		adminRoutes.PUT("/user/:user_id/active", controllers.ActiveAccount(repo)) //Ban 1 người dùn
		adminRoutes.GET("/payment-history", controllers.GetPaymentHistoryForAdmin(payrepo))
		adminRoutes.GET("/payment-history/:user_id", controllers.GetPaymentHistoryForUserByAdmin(payrepo))
	}

	// Route cho payment
	paymentRoutes := r.Group("/api/v1/payment")
	{
		paymentRoutes.POST("/vip-upgrade", middlewares.AuthMiddleware("VIP-0", "VIP-1", "VIP-2"), controllers.CreateVIPPayment())
		paymentRoutes.POST("/status", controllers.HandleQueryPaymentStatus())
		paymentRoutes.POST("/momo-callback", momo.MoMoCallback())
	}

	// return r
}
