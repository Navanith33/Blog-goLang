package main

import (
	"fmt"
	"log"

	"example.com/blog/Routes"
	"example.com/blog/models"
	"example.com/blog/supabase"
	"github.com/gin-gonic/gin"
	"github.com/supabase-community/auth-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)
var db *gorm.DB;
var client auth.Client;
func init(){
	er := godotenv.Load()
    if er != nil {
    log.Fatal("Error loading .env file")
  }
	dsn := "host=localhost user=navanithravi password=Navanith@56 dbname=todo port=5432 sslmode=disable TimeZone=Asia/Shanghai"
    var err error
    db,err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }
	fmt.Println("connected to database")
	err = db.AutoMigrate(&models.User{}, &models.Blog{})
	if err != nil {
	     fmt.Println("Error migrating models");
	}
    if err != nil {
        log.Fatalf("Cannot initialize client: %v", err)
    }
	fmt.Println("Tables migrated successfully")
}   
func main(){
    supabaseUrl := "dvdmrmtgniwwwynjwomc"
    supabaseKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImR2ZG1ybXRnbml3d3d5bmp3b21jIiwicm9sZSI6ImFub24iLCJpYXQiOjE3MzQ2Nzk0NjUsImV4cCI6MjA1MDI1NTQ2NX0.4h9kQiGTU__eDrVbJHGTsJGWEuplZ0956i0q0VfokQ8"
    client= supabase1.NewSupabaseClient(supabaseUrl,supabaseKey);
    fmt.Println(client);
	

	r := gin.Default()
	Routes.InitializeRoutes(r, db,client) 
	// r.POST("/signin",Routes.SigninUser)
    // r.GET("/login",Routes.LoginUser)
	// r.POST("addBlog",Routes.CreateBlog)
	// r.GET("getBlogs",Routes.GetBlogs)
    // r.DELETE("deleteBlog/:id/:blogId",Routes.DeleteBlog)
	// r.PUT("updateBlog/:id/:blogId",Routes.UpdateBlog)
	r.Run()


}