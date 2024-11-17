package routes

import (
	"github.com/dath-241/coin-price-be-go/services/admin_service/controllers"
	"github.com/dath-241/coin-price-be-go/services/admin_service/middlewares"
	"github.com/dath-241/coin-price-be-go/services/admin_service/momo"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	// Route User Management API (Quản lý tài khoản)
	userRoutes := r.Group("/api/v1/user")
	{
		// Áp dụng AuthMiddleware cho tất cả các route trong group này
		userRoutes.Use(middlewares.AuthMiddleware("VIP-0", "VIP-1", "VIP-2", "VIP-3", "Admin"))

		userRoutes.GET("/me", controllers.GetCurrentUserInfo())   // Lấy thông tin tài khoản người dùng hiện tại.
		userRoutes.PUT("/me", controllers.UpdateCurrentUser())    // Chỉnh sửa thông tin tài khoản người dùng.
		userRoutes.DELETE("/me", controllers.DeleteCurrentUser()) // Xóa tài khoản người dùng hiện tại.
		userRoutes.GET("/me/payment-history", controllers.GetPaymentHistory())
	}

	// Route cho xác thực
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/login", controllers.Login()) // Người dùng đăng nhập
		authRoutes.POST("/google-login", controllers.GoogleLogin)
		authRoutes.POST("/register", controllers.Register())                                                                      // Người dùng đăng kí
		authRoutes.POST("/forgot-password", controllers.ForgotPassword())                                                         // Người dùng quên mật khẩu
		authRoutes.POST("/reset-password", controllers.ResetPassword())                                                           // Người dùng thay đổi mật khẩu
		authRoutes.POST("/refresh-token", controllers.RefreshToken())                                                             // Người dùng refresh token
		authRoutes.POST("/logout", middlewares.AuthMiddleware("VIP-0", "VIP-1", "VIP-2", "VIP-3", "Admin"), controllers.Logout()) // Người dùng đăng xuất
	}

	// Route cho quản trị viên (admin)
	adminRoutes := r.Group("/api/v1/admin/")
	{
		// Áp dụng AuthMiddleware cho tất cả các route trong group này
		adminRoutes.Use(middlewares.AuthMiddleware("Admin"))

		adminRoutes.GET("/users", controllers.GetAllUsers())                  //Lay thong tin tat ca nguoi dung
		adminRoutes.GET("/user/:user_id", controllers.GetUserByAdmin())       //Lay thong tin 1 user cu the
		adminRoutes.DELETE("/user/:user_id", controllers.DeleteUserByAdmin()) //Xóa 1 người dùng
		adminRoutes.GET("/payment-history", controllers.GetPaymentHistoryForAdmin())
	}

	// Route cho payment
	paymentRoutes := r.Group("/api/v1/payment")
	{
		paymentRoutes.POST("/vip-upgrade", middlewares.AuthMiddleware("VIP-0", "VIP-1", "VIP-2"), controllers.CreateVIPPayment())
		paymentRoutes.POST("/confirm", controllers.ConfirmPaymentHandlerSuccess())
		paymentRoutes.POST("/status", controllers.HandleQueryPaymentStatus())
		paymentRoutes.POST("/momo-callback", momo.MoMoCallback())
	}

	// return r
}
