package Routes

import (
	"fmt"

	"strconv"
	// "github.com/supabase-community/supabase-go"
	"example.com/blog/middlewares"
	"example.com/blog/models"
	"github.com/gin-gonic/gin"
	"github.com/supabase-community/auth-go"
	"github.com/supabase-community/auth-go/types"
	"gorm.io/gorm"
)
func InitializeRoutes(router *gin.Engine, db *gorm.DB,client  auth.Client) {
	router.POST("/signin", func(c *gin.Context) { SignupUser(c, db,client) })
	router.POST("/login", func(c *gin.Context) { LoginUser(c, db,client) })
	router.POST("/addBlog",middlewares.AuthMiddleware(), func(c *gin.Context) { CreateBlog(c, db) })
	router.GET("/getBlogs/:id", func(c *gin.Context) { GetBlogs(c, db) })
	router.DELETE("/deleteBlog/:id/:blogId",middlewares.AuthMiddleware(), func(c *gin.Context) { DeleteBlog(c, db) })
	router.PUT("/updateBlog/:id/:blogId",middlewares.AuthMiddleware(),func(c *gin.Context) { UpdateBlog(c, db) })
}
func CreateBlog(c *gin.Context,db *gorm.DB) {
	if userRole, exists := c.Get("user_role"); exists {
        if role, ok := userRole.(string); ok {
            if role != "Admin"{
				c.JSON(401, gin.H{"message": "unAuthorized"})
			}
			return;
        } 
    } else {
        c.JSON(404, gin.H{"error": "Role not found"})
    }
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
func SignupUser(c *gin.Context,db *gorm.DB,client auth.Client) {
    var body struct{
		Email   string
		Password string
		Role     string
	}
	c.Bind(&body);
	tokenRequest := types.SignupRequest{
        Email:body.Email,
		Phone:"",      
	    Password:body.Password,
		Data: map[string]interface{}{
			"role":body.Role,
		},
		
    }
    res,err:=client.Signup(tokenRequest);
	if err != nil {
		c.JSON(200, gin.H{
			"messsage":err.Error(),
		})
		return;
	
	}
    if res == nil {
		fmt.Println("res is nill");
		return
	}
	token := res.AccessToken
	if err !=nil{
		c.JSON(400, gin.H{
			"messsage":err.Error(),
		})
	}
	var storeuser models.User
	result := db.Where("email = ?", body.Email).First(&storeuser)
	if result.Error == nil {
		c.JSON(400, gin.H{
			"message": "Email already registered",
		})
		return
	}
	user := models.User{Email:body.Email,Password:body.Password,Role:body.Role}
    result1 := db.Create(&user) 
	if result1.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(200, gin.H{
		"Token":token,
		"message": "signinsuccessfully",
	})
}
func LoginUser(c *gin.Context,db *gorm.DB,client auth.Client) {
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
	res,err:=client.SignInWithEmailPassword(body.Email,body.Password);
	if err ==nil{
		c.JSON(200, gin.H{
            "message": "Login successful",
			"token":res.AccessToken,
        })
        return
	}
	// var storeuser models.User
	// result :=db.Where("email = ? AND password = ?",body.Email,body.Password).First(&storeuser);
	// if result.Error != nil{
	// 	c.JSON(400,gin.H{
	// 		"message":"Login failed",
	// 	})
	// 	return;
	// }
	c.JSON(401,gin.H{
		"message":"Login unsuccessful",
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
	if userRole, exists := c.Get("user_role"); exists {
        if role, ok := userRole.(string); ok {
            if role != "Admin"{
				c.JSON(401, gin.H{"message": "unAuthorized"})
			}
			return;
        } 
    } else {
        c.JSON(404, gin.H{"error": "Role not found"})
    }
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
	if userRole, exists := c.Get("user_role"); exists {
        if role, ok := userRole.(string); ok {
            if role != "Admin"{
				c.JSON(401, gin.H{"message": "unAuthorized"})
			}
			return;
        } 
    } else {
        c.JSON(404, gin.H{"error": "Role not found"})
    }
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