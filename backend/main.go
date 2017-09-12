package main
import(
	"net/http"
	r "gopkg.in/gorethink/gorethink.v3"
	"log"

	"time"
)

type Channel struct {
	Id string `json:"id" gorethink:"id,omitempty"`
	Name string `json:"name" gorethink:"name"`
}
type User struct{
	Id string `json:"id" gorethink:"id,omitempty"`
	Name string `json:"name" gorethink:"name"`
}
type ChannelMessage struct {
	Id string `gorethink:"id,omitempty"`
	ChannelId string `gorethink:"channelId"`
	Body string `gorethink:"body"`
	Author string `gorethink: "author"`
	CreatedAt time.Time `gorethink:"createdAt"`
}
func main(){
	session, err :=r.Connect(r.ConnectOpts{
		Address:"localhost:28015",
		Database:"realtime_message",
	})
	if err!=nil{
		log.Panic(err.Error())

	}

	router := NewRouter(session)

	router.Handle("channel add", addChannel)
	router.Handle("channel subscribe",subscribeChannel)
	router.Handle("channel unsubscribe", unsubscribeChannel)

	router.Handle("user edit", editUser)
	router.Handle("user subscribe", subscribeUser)
	router.Handle("user unsubscribe", unsubscribeUser)

	router.Handle("message add", addChannelMessage)
	router.Handle("message subscribe", subscribeChannelMessage)
	router.Handle("message unsubscribe", unsubscribeChannelMessage)


	http.Handle("/",router)

	http.ListenAndServe(":4000",nil)
}
