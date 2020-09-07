package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http"
	"strings"
	//"time"
)



func err_handler(err error) {
	fmt.Printf("err_handler, error:%s\n", err.Error())
	panic(err.Error())
}


func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "114.212.10.200:6379",
		Password: "",
		DB:       0,
	})

	defer client.Close()

	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Printf("ping error[%s]\n", err.Error())
		err_handler(err)
	}

	fmt.Printf("ping result: %s\n", pong)

	router := gin.Default()
	router.GET("/user", func(c *gin.Context) {
		fmt.Println(client)
		c.Header("Access-Control-Allow-Origin", "*")
		//name := c.Param("name")
		name2 :=strings.Split(c.Request.URL.String(), "=")[1]


		//name := c.GetHeader("name")
		//fmt.Printf("student id is :%v\n", name)
		fmt.Printf("url  is :%v\n", name2)
		// 此处的name 去查redis 学号做key , 序号做value
		value, err := client.Get(name2).Result()
		if err != nil {
			fmt.Printf("try get key[%v] error[%s]\n", name2, err.Error())
			// err_handler(err)
		}

		c.String(http.StatusOK, "%s", value)
		//fmt.Println(name)
		fmt.Printf("get key [%v] , value is [%v] \n", name2, value)
	})

	router.POST("/userIDs", func(context *gin.Context) {
		// 此处存储redis, 先全部删除，再重新存储
		context.Header("Access-Control-Allow-Origin", "*")

		context.String(http.StatusOK, "Hello simzhang deal IDs", )
		simMapString := context.PostForm("IDS")
		//fmt.Println("in iDs fun")
		fmt.Println(simMapString)
		strKeyValSingle := strings.Split(simMapString, "&")
		var i int
		for i = 0; i < len(strKeyValSingle); i++ {
			//fmt.Println(strKeyValSingle[i])
			kv := strings.Split(strKeyValSingle[i], "=")
			reKey := kv[0]
			reVal := kv[1]
			//fmt.Printf("key is : %v, and val is : %v \n", reKey, reVal)

			err := client.Set(reKey, reVal, 0).Err()
			if err != nil {
				fmt.Printf("try set key[%v] to value[%v] error[%s]\n",
					reKey, reVal, err.Error())
				err_handler(err)
			}
		}

	})
	router.Run(":8070") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}