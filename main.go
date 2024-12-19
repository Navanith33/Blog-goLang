package main

import (
	"fmt"
	"log"
	"example.com/blog/models"
	"example.com/blog/Routes"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
var db *gorm.DB;
func init(){
	dsn := "host=localhost user=navanithravi password=Navanith@56 dbname=todo port=5432 sslmode=disable TimeZone=Asia/Shanghai"
    var err error
    db,err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }
	fmt.Println("connected to database")
	err = db.AutoMigrate(&models.User{}, &models.Blog{})
	if err != nil {
		log.Fatalf("Error migrating models", err)
	}
	fmt.Println("Tables migrated successfully")
}   
func main(){
    
	r := gin.Default()
	Routes.InitializeRoutes(r, db) 
	// r.POST("/signin",Routes.SigninUser)
    // r.GET("/login",Routes.LoginUser)
	// r.POST("addBlog",Routes.CreateBlog)
	// r.GET("getBlogs",Routes.GetBlogs)
    // r.DELETE("deleteBlog/:id/:blogId",Routes.DeleteBlog)
	// r.PUT("updateBlog/:id/:blogId",Routes.UpdateBlog)
	r.Run()


}