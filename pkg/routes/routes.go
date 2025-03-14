package routes

import (
	"aesthetics/pkg/handlers"
	"aesthetics/smtp"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupRoutes(r *gin.Engine, db *sql.DB, smtpClient *smtp.SMTPClient, redisClient *redis.Client) {
	v1 := r.Group("/api/v1")
	{
		v1.GET("/", handlers.HomePage)

		sexGroup := v1.Group("/:gender")
		{
			sexGroup.GET("/", handlers.GetGender)

			categoryGroup := sexGroup.Group("/:category")
			{

				categoryGroup.GET("/", handlers.GetCategory(db))

				subCategoryGroup := categoryGroup.Group("/:subcategory")
				{
					subCategoryGroup.GET("/", handlers.GetSubCategory(db))

					productGroup := subCategoryGroup.Group("/:productID")
					{
						productGroup.GET("/", handlers.GetProducts(db))
					}
				}
			}
		}

		v1.GET("/cart", handlers.AuthMiddleware(db))               // Получить корзину пользователя
		v1.POST("/cart/add", handlers.AuthMiddleware(db))          // Добавить товар в корзину
		v1.PUT("/cart/update/:id", handlers.AuthMiddleware(db))    // Обновить количество товаров в корзине
		v1.DELETE("/cart/remove/:id", handlers.AuthMiddleware(db)) // Удалить товар из корзины

		// ЗАКАЗ
		v1.GET("/orders", handlers.AuthMiddleware(db))            // Получить все заказы пользователя
		v1.GET("/orders/:id", handlers.AuthMiddleware(db))        // Получить заказы по ID
		v1.POST("/orders", handlers.AuthMiddleware(db))           // Создать новый заказ
		v1.PUT("/orders/:id/cancel", handlers.AuthMiddleware(db)) // Отменить заказ

		// ПОЛЬЗОВАТЕЛЬ
		v1.POST("/register", handlers.RegisterPage(db)) // Зарегистрироваться
		v1.POST("/login", handlers.LoginPage(db))       // Залогиниться

		protected := v1.Group("/profile")
		protected.Use(handlers.AuthMiddleware(db))
		protected.GET("/", handlers.GetProfile(db)) // Просмотреть профиль
		protected.PUT("/")                          // Обновить профиль
		protected.POST("/address")                  // Обновить профиль

		// Email
		v1.POST("/subscribe", handlers.HandleEmail(smtpClient))

		//Admin panel
		v1.GET("admin/users", handlers.GetUsers(db))
		v1.PUT("admin/users/:id", handlers.UpdateUser(db))
		v1.DELETE("admin/users/:id", handlers.DeleteUser(db))

		v1.GET("admin/category", handlers.GetCategory(db))           // Получение конкретной категории по id
		v1.PUT("admin/category/:id", handlers.UpdateCategory(db))    // Обновление категории по id
		v1.DELETE("admin/category/:id", handlers.DeleteCategory(db)) // Удаление категории по id
		v1.POST("admin/category", handlers.CreateCategory(db))

		v1 := r.Group("/api/v1")
		{
			admin := v1.Group("/admin")
			{
				admin.GET("/subcategories", handlers.GetSubCategories(db))
				admin.GET("/subcategories/:id", handlers.GetSubCategory(db))
				admin.POST("/subcategories", handlers.CreateSubCategory(db))
				admin.PUT("/subcategories/:id", handlers.UpdateSubCategory(db))
				admin.DELETE("/subcategories/:id", handlers.DeleteSubCategory(db))
			}
		}

		v1.GET("admin/products", handlers.GetProducts(db))
		v1.GET("admin/products/:productID", handlers.GetProduct(db))
		v1.POST("admin/products", handlers.AddProduct(db))
		v1.PUT("admin/products/:id", handlers.UpdateProduct(db))
		v1.DELETE("admin/products/:id", handlers.DeleteProduct(db))
		v1.POST("admin/products/refresh", handlers.RefreshProducts(db))
	}
}
