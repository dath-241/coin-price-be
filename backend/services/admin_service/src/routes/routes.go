package routes

import (
	"backend/services/admin_service/src/middlewares"
	"backend/services/admin_service/src/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Route User Management API (Quản lý tài khoản)
	user := r.Group("/api/v1/user")
	{
		user.GET("/me", middlewares.AuthMiddleware("VIP-0", "VIP-1", "VIP-2", "VIP-3", "Admin"), controllers.GetCurrentUserInfo())          // Lấy thông tin tài khoản người dùng hiện tại.
		user.PUT("/me", middlewares.AuthMiddleware("VIP-0", "VIP-1", "VIP-2", "VIP-3", "Admin"), controllers.UpdateCurrentUser())           // Chỉnh sửa thông tin tài khoản người dùng.
		user.DELETE("/me", middlewares.AuthMiddleware("VIP-0", "VIP-1", "VIP-2", "VIP-3", "Admin"), controllers.DeleteCurrentUser())           // Xóa tài khoản người dùng hiện tại.
	}

	// Route cho xác thực
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/login", controllers.Login()) // Người dùng đăng nhập
		authRoutes.POST("/google-login", controllers.GoogleLogin)
		authRoutes.POST("/register", controllers.Register())              // Người dùng đăng kí
		authRoutes.POST("/forgot-password", controllers.ForgotPassword()) // Người dùng quên mật khẩu
		authRoutes.POST("/reset-password", controllers.ResetPassword())   // Người dùng thay đổi mật khẩu
		authRoutes.POST("/refresh-token", middlewares.AuthMiddleware("VIP-0", "VIP-1", "VIP-2", "VIP-3", "Admin"), controllers.RefreshToken())     // Người dùng refresh token
		authRoutes.POST("/logout", middlewares.AuthMiddleware("VIP-0", "VIP-1", "VIP-2", "VIP-3", "Admin"), controllers.Logout())                  // Người dùng đăng xuất
	}

	// Route cho quản trị viên (admin)
	adminRoutes := r.Group("/api/v1/admin/")
	{
		adminRoutes.GET("/users", middlewares.AuthMiddleware("Admin"), controllers.GetAllUsers())            //Lay thong tin tat ca nguoi dung
		//adminRoutes.GET("/user/:user_id", middlewares.AuthMiddleware("Admin"), controllers.GetUserByAdmin()) //Lay thong tin 1 user cu the
		//adminRoutes.DELETE("/user/:user_id", middlewares.AuthMiddleware("Admin"), controllers.DeleteUser())  //Xóa 1 người dùng
	}

	// Route cho payment 
	paymentRoutes := r.Group("/api/v1/payment")
	{
		paymentRoutes.POST("/vip-upgrade", middlewares.AuthMiddleware("VIP-0", "VIP-1", "VIP-2"), controllers.CreateVIPPayment())
		paymentRoutes.POST("/confirm", controllers.ConfirmPaymentHandlerSuccess())
	}

	return r
}
