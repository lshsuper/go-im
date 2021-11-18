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
		Pwd: "nicaia",
		PoolSize: 100,
		Addr: ":6379",
		Tag:"Simple",
	})

	var redisProvider=xredis.RedisFatory.Client("Simple")

	server := socketio.NewServer(&engineio.Options{
		SessionIDGenerator: models.UIDGenerator{},
	})

	server.OnConnect("/", func(s socketio.Conn) error {
		log.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/","join", func(s socketio.Conn, msg string){
		s.Join(msg)
		s.SetContext(msg)
		redisProvider.HashMSet(context.Background(),msg, map[string]interface{}{cast.ToString(s.ID()):cast.ToString(s.ID())})
		server.BroadcastToRoom("/",msg,"join_msg",cast.ToString(s.ID()))

	})

	server.OnEvent("/", "send_to_group", func(s socketio.Conn, req models.SendToGroupRequest)  {
		s.SetContext(req)

		server.BroadcastToRoom("/",req.GroupID,"get_group_msg",models.GroupMsgDTO{
			Msg: req.Msg,
			UserID: cast.ToString(s.ID()),
		})
	})


	server.OnEvent("/","get_group_member_count",func(s socketio.Conn, req models.GroupMemberCountRequest)models.GroupMemberCountDTO {

		 num,_:=redisProvider.HashLen(context.Background(),req.GroupID)
		 return models.GroupMemberCountDTO{
			 Count: num,
		 }

	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		s.LeaveAll()
		log.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		//s.Leave()
		g:=s.Context()
        s.Close()
		s.LeaveAll()
		log.Println("closed", g)
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
