package Routes

import (
	"fmt"
	"strconv"
	"example.com/blog/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)
func InitializeRoutes(router *gin.Engine, db *gorm.DB) {
	router.POST("/signin", func(c *gin.Context) { SigninUser(c, db) })
	router.POST("/login", func(c *gin.Context) { LoginUser(c, db) })
	router.POST("/addBlog", func(c *gin.Context) { CreateBlog(c, db) })
	router.GET("/getBlogs/:id", func(c *gin.Context) { GetBlogs(c, db) })
	router.DELETE("/deleteBlog/:id/:blogId", func(c *gin.Context) { DeleteBlog(c, db) })
	router.PUT("/updateBlog/:id/:blogId", func(c *gin.Context) { UpdateBlog(c, db) })
}
func CreateBlog(c *gin.Context,db *gorm.DB) {
	var body struct{
		Title   string
		Content string
		UserId  int
	}
	c.Bind(&body);
	user := models.Blog{Title:body.Title,Content:body.Content,UserId:uint(body.UserId)}
    result := db.Create(&user) 
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create blog", "details": result.Error.Error()})
		return
	}
	c.JSON(200,gin.H{
		"message":"Blog added successfully",
	})
}
func SigninUser(c *gin.Context,db *gorm.DB) {
    var body struct{
		Email   string
		Password string
	}
	c.Bind(&body);
	var storeuser models.User
	result := db.Where("email = ?", body.Email).First(&storeuser)
	if result.Error == nil {
		c.JSON(400, gin.H{
			"message": "Email already registered",
		})
		return
	}
	user := models.User{Email:body.Email,Password:body.Password}
    result1 := db.Create(&user) 
	if result1.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(200, gin.H{
		"message": "signinsuccessfully",
	})
}
func LoginUser(c *gin.Context,db *gorm.DB) {
	var body struct{
		Email    string
		Password string
	}
	if err := c.ShouldBind(&body); err != nil {
        c.JSON(400, gin.H{
            "message": "Invalid request data",
        })
        return
    }
	var storeuser models.User
	result :=db.Where("email = ? AND password = ?",body.Email,body.Password).First(&storeuser);
	if result.Error != nil{
		c.JSON(400,gin.H{
			"message":"Login failed",
		})
		return;
	}
	c.JSON(200,gin.H{
		"message":"Login successful",
	})
	
}
func GetBlogs(c *gin.Context,db *gorm.DB) {
	id:=c.Param("id");
	parsedId,err:=strconv.Atoi(id);
	if err !=nil{
		c.JSON(200,gin.H{
			"message":"parsing error",
		})
	}
	var blogs []models.Blog
	result :=db.Preload("User").Where("user_Id = ?",parsedId).Find(&blogs);
	fmt.Println(blogs);
	if result.Error == nil{
		c.JSON(200,gin.H{
			"message":blogs,
		})
		return;
	}
	c.JSON(400,gin.H{
		"message":"no blogs found",
	})
}
func DeleteBlog(c *gin.Context,db *gorm.DB) {
	id:=c.Param("id");
	blogId := c.Param("blogId")
	parsedId,err:=strconv.Atoi(id);
	if err !=nil{
		c.JSON(200,gin.H{
			"message":"parsing error",
		})
	}
	parsedblogId,err:=strconv.Atoi(blogId);
	if err !=nil{
		c.JSON(200,gin.H{
			"message":"parsing error",
		})
	}
	var blog models.Blog
	result :=db.Where("user_id = ? AND id=?",parsedId,parsedblogId).First(&blog)
    if result.Error==nil{
          delete:=db.Delete(&blog);
		  if delete.Error==nil{
			c.JSON(400,gin.H{
				"message":"deleted successfully",
			})
			return;
		  }
	}
	c.JSON(200,gin.H{
		"message":"records not found",
	})

	
}
func UpdateBlog(c *gin.Context,db *gorm.DB) {
	var body struct{
		Title    string
		Content  string
	}
	c.Bind(&body);
	id:=c.Param("id");
	blogId := c.Param("blogId")
	parsedId,err:=strconv.Atoi(id);
	if err !=nil{
		c.JSON(200,gin.H{
			"message":"parsing error",
		})
	}
	parsedblogId,err:=strconv.Atoi(blogId);
	if err !=nil{
		c.JSON(200,gin.H{
			"message":"parsing error",
		})
	}
	var blog models.Blog
	result :=db.Where("user_id = ? AND id=?",parsedId,parsedblogId).First(&blog);
    if result.Error==nil{
        blog.Title=body.Title;
		blog.Content=body.Content;
		update := db.Save(&blog)
		if update.Error ==nil{
			c.JSON(200,gin.H{
				"message":"updated successfull",
			})
			return;
		}
	}
	c.JSON(200,gin.H{
		"message":"updated unsuccessfull",
	})
}