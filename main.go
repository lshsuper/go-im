package main

import (
	"context"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/lshsuper/go-pkg/src/xredis"
	"github.com/lshsuper/go-pkg/src/xredis/domain"
	"github.com/spf13/cast"
	"go-im/handlers"
	"go-im/models"
	"log"
	"net/http"
)

func main() {
	router := gin.New()

	xredis.Register()
    xredis.RedisFatory.Set(domain.RedisOpt{
		Pwd: "honghe@2020",
		PoolSize: 100,
		Addr: ":6379",
		Tag:"Simple",
	})
	var redisProvider=xredis.RedisFatory.Client("Simple")

	server := socketio.NewServer(&engineio.Options{
		SessionIDGenerator: models.UIDGenerator{},

	})


	server.OnConnect("/", func(s socketio.Conn) error {


		s.SetContext("")
		log.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/","join", func(s socketio.Conn, msg string) string{
		s.Join(msg)
		redisProvider.HashMSet(context.Background(),msg, map[string]interface{}{cast.ToString(s.ID()):cast.ToString(s.ID())})
		  //s.SetContext("加入聊天室成功...")
		return "加入聊天室 "+msg
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})


	server.OnEvent("/","getGroupNumber",func(s socketio.Conn, msg string)string {

		 num,_:=redisProvider.HashLen(context.Background(),msg)
		 return cast.ToString(num)
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		//s.Leave()
		log.Println("closed", reason)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()

	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))
	router.StaticFS("/public", http.Dir("./asset"))

	router.POST("/im/authorization", handlers.Authorization)

	if err := router.Run(":8000"); err != nil {
		log.Fatal("failed run app: ", err)
	}
}
