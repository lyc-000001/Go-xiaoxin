package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/xiaoxin/blog-backend/internal/controllers"
	"github.com/xiaoxin/blog-backend/internal/middleware"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// 初始化控制器
	userCtrl := controllers.NewUserController()
	articleCtrl := controllers.NewArticleController()
	categoryCtrl := controllers.NewCategoryController()
	uploadCtrl := controllers.NewUploadController()

	// 公开路由
	api := r.Group("/api/v1")

	// 用户相关
	api.POST("/register", userCtrl.Register)
	api.POST("/login", userCtrl.Login)

	// 文章相关（公开访问）
	api.GET("/articles", articleCtrl.GetArticleList)
	api.GET("/articles/:id", articleCtrl.GetArticle)

	// 分类相关（公开访问）
	api.GET("/categories", categoryCtrl.GetCategoryList)
	api.GET("/categories/:id", categoryCtrl.GetCategory)

	// 需要认证的路由
	auth := r.Group("/api/v1")
	auth.Use(middleware.JWTAuth())
	{
		// 用户相关
		auth.GET("/user/profile", userCtrl.GetProfile)
		auth.PUT("/user/profile", userCtrl.UpdateProfile)
		auth.PUT("/user/password", userCtrl.ChangePassword)

		// 文件上传
		auth.POST("/upload", uploadCtrl.UploadFile)

		// 文章相关（需要认证）
		auth.POST("/articles", articleCtrl.CreateArticle)
		auth.PUT("/articles/:id", articleCtrl.UpdateArticle)
		auth.DELETE("/articles/:id", articleCtrl.DeleteArticle)
		auth.POST("/articles/:id/like", articleCtrl.LikeArticle)
	}

	// 管理员路由
	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.JWTAuth(), middleware.RequireRole("admin"))
	{
		// 用户管理
		admin.GET("/users/:id", userCtrl.GetUserByID)

		// 分类管理
		admin.POST("/categories", categoryCtrl.CreateCategory)
		admin.PUT("/categories/:id", categoryCtrl.UpdateCategory)
		admin.DELETE("/categories/:id", categoryCtrl.DeleteCategory)
	}

	// 静态文件服务（上传的文件）
	r.Static("/uploads", "./uploads")
}
