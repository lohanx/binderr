# binderr
```go
package main

import (
    "errors"
    "github.com/gin-gonic/gin"
    "github.com/lohanx/binderr"
    "net/http"
)

var msgErrors = map[string]map[string]error{
    "name":{
        "required":errors.New("name is required"),
    },
    "email":{
        "required":errors.New("email is required"),
        "email":errors.New("email format error"),
    },
}

type Person struct {
    Name string `form:"name" binding:"required"`
    Email string `form:"email" binding:"required,email"`
}

func main() {
    router := gin.New()
    router.GET("/get-person",func(ctx *gin.Context){
        var person Person
        if err := ctx.ShouldBind(&person);err != nil {
            ctx.JSON(http.StatusBadRequest,gin.H{
                "code":-1,
                "err":binderr.New(err,msgErrors).FirstError().Error(),
            })
            return
        }
        ctx.JSON(http.StatusOK,gin.H{
            "code":0,
            "data":person,
        })
    })
    router.Run("127.0.0.1:8080")
}
```