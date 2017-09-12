package main
import(
	"github.com/mitchellh/mapstructure"

	r "gopkg.in/gorethink/gorethink.v3"
	"fmt"
	"time"

)

const(
	ChannelStop = iota
	UserStop
	MessageStop
)

func addChannel(client *Client, data interface{}){
	var channel Channel

	err :=mapstructure.Decode(data,&channel)
	fmt.Println(channel)
	if err !=nil{
		client.send <- Message{"error",err.Error()}
		return
	}

	// TODO : insert RethinkDB
	err =r.Table("channel").Insert(channel).Exec(client.session);
	if err !=nil{
		client.send <- Message{"error",err.Error()}
		return
	}


}
func addChannelMessage(client *Client, data interface{}){
	var channelMessage ChannelMessage
	if err := mapstructure.Decode(data, &channelMessage); err!=nil{
		client.send <-Message{"error",err.Error()}
	}
	fmt.Println("ChannelMessage")
	fmt.Println(channelMessage)
	go func(){
		channelMessage.CreatedAt = time.Now()
		channelMessage.Author = client.userName
		if err:= r.Table("message").Insert(channelMessage).Exec(client.session); err!=nil{
			client.send <-Message{"error",err.Error()}
		}
	}()

}
func subscribeChannel(client *Client, data interface{}){
	stop := client.NewStopChannel(ChannelStop)
	result := make(chan r.ChangeResponse)
	cursor, err := r.Table("channel").Changes(r.ChangesOpts{IncludeInitial: true}).Run(client.session)
	if err != nil{
		client.send <- Message{"error", err.Error()}
		return
	}
	go func() {

		var change r.ChangeResponse
		for cursor.Next(&change){
			result <-change

		}

	}()
	////////////////////
	go func(){
		for{
			select{
				case <-stop:
					cursor.Close()
					return
				case change := <-result:
					if change.NewValue !=nil && change.OldValue ==nil{
						client.send <- Message{"channel add", change.NewValue}
						fmt.Println("sent channel add msg");
					}
			}
		}
	}()
}
func unsubscribeChannel(client *Client, data interface{}){
	client.StopForKey(ChannelStop)

}
func editUser(client *Client, data interface{}){
	var user User
	err :=mapstructure.Decode(data,&user);

	if err!=nil{
		client.send <- Message{"error", err.Error()}
		return
	}
	client.userName = user.Name
	go func(){
		err := r.Table("user").Get(client.id).Update(user).Exec(client.session)
		if err != nil{
			client.send <- Message{"error", err.Error()}
			return
		}



	}()
}

func subscribeUser(client *Client, data interface{}){
	go func(){
		stop := client.NewStopChannel(UserStop)
		cursor, err := r.Table("user").Changes(r.ChangesOpts{IncludeInitial:true}).Run(client.session)

		if err != nil{
			client.send <- Message{"error", err.Error()}
		}
		changeFeedHelper(cursor, "user", client.send, stop)
	}()
}
func unsubscribeUser(client *Client, data interface{}){
	client.StopForKey(UserStop)
}



func subscribeChannelMessage(client *Client, data interface{}){
	go func(){
		eventData := data.(map[string]interface{})
		val, ok:= eventData["channelId"]
		if !ok{
			return
		}
		channelId,ok:=val.(string)
		if !ok{
			return
		}
		stop :=client.NewStopChannel(MessageStop)
		//cursor, err := r.Table("message").
		//	OrderBy(r.OrderByOpts{Index:r.Desc("createdAt")}).
		//	Filter(r.Row.Field("channelId").Eq(channelId)).
		//	Changes(r.ChangesOpts{IncludeInitial:true}).
		//	Run(client.session)
		cursor, err := r.Table("message").OrderBy(r.OrderByOpts{Index:r.Desc("createdAt")}).Filter(func(uu r.Term) r.Term{
			return uu.Field("channelId").Eq(channelId)
		}).Changes(r.ChangesOpts{IncludeInitial:true}).Run(client.session)



		if err != nil{
			fmt.Println(err);
			client.send <-Message{"error",err.Error()}
		}
		changeFeedHelper(cursor, "message",client.send,stop)
	}()


}
func unsubscribeChannelMessage(client *Client, data interface{}){
	client.StopForKey(MessageStop)
}

func changeFeedHelper(cursor *r.Cursor, changeEventName string, send chan<- Message, stop <-chan bool){
	change := make(chan r.ChangeResponse)
	if(changeEventName=="message"){
		fmt.Println("Message allready executed")

	}
	cursor.Listen(change)
	for {
		eventName :=""
		var data interface{}
		select{
		case <-stop:
			cursor.Close()
			return
		case val := <-change:
			fmt.Println("Change channel executed")
			if val.NewValue !=nil&& val.OldValue ==nil{
				eventName = changeEventName +" add"
				data = val.NewValue
			} else if val.NewValue ==nil&& val.OldValue !=nil {
				eventName = changeEventName + " remove"
				data = val.OldValue
			} else if val.NewValue !=nil&& val.OldValue !=nil {

				eventName = changeEventName + " edit"
				data = val.NewValue
			}
			sendMessage := Message{eventName,data}
			fmt.Println("event name:"+eventName)
			fmt.Println(sendMessage)
			send <- sendMessage
		}
	}
}
