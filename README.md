# tools
go programming tools

1. #### field validate

   It provides a convenient method to verify the validity of the struct. 

   ```go
   type User struct {
      Name  string `validate:"min=4; max=16; noSpace"`
      Email string `validate:"email; noSpace"` 
   }
   
   user := User{
       ...
   }
   err := validate.Validate(&user)
   ```

   Built in tags:

   | number |                    |
   | ------ | ------------------ |
   | eq     | Equals             |
   | gt     | GreaterThan        |
   | gte    | GreaterThanOrEqual |
   | lt     | LessThan           |
   | lte    | LessThanOrEqual    |
   | ne     | NotEqual           |

   | String  |               |
   | ------- | ------------- |
   | min     | min length    |
   | max     | max length    |
   | length  | length        |
   | noSpace | with no space |

   | other |                      |
   | ----- | -------------------- |
   | email | Email address format |
   | url   | urladdress format    |

   

2. ### http client

   Based on net/http package, provide chain operation of HTTP request.

   ```go
   	user := User{}
   	c := client.New("192.168.0.1:8080/api/v1/user")
   	err := c.AddHeader("Authorization", "Authorization token").
   		AddQuery("param1", "1").
   		AddQuery("param2", "2").
   		Get().
   		Do(&user)
   	if err != nil {
   		panic(err)
   	}
   
   
   // Prehandler: Server data preprocessing method
   type PreHandler interface {
   	PreHandler(body []byte) ([]byte, error)
   }
   
   type User struct {
   	Name string
   }
   
   type HttpResponse struct {
   	Code int
   	Msg  string
   	Data interface{}
   }
   
   func (r HttpResponse) PreHandler(body []byte) (data []byte,err error) {
   	...
   	return 
   }
   
   
   func main() {
   	
   	user := User{}
   	c := client.New("192.168.0.1:8080/api/v1/user")
   	err := c.AddHeader("Authorization", "Authorization token").
   		AddQuery("param1", "1").
   		AddQuery("param2", "2").
   		Get().
   		AddPreHandler(&HttpResponse{}).
   		Do(&user)
   	if err != nil {
   		panic(err)
   	}
   }
   
   ```



3. ### ssh

   SSH client based on websocket. The SSH Connection code comes from [kubekey]([kubesphere/kubekey: Provides a flexible, rapid and convenient way to install Kubernetes only, both Kubernetes and KubeSphere, and related cloud-native add-ons. It is also an efficient tool to scale and upgrade your cluster. (github.com)](https://github.com/kubesphere/kubekey))

   ```go
   connection, err := ssh.NewConnection(ssh.Cfg{
   	Username: "root",
   	Password: "...your password",
   	Address:  "...your ssh address",
   	Port:     22,
   })
   if err != nil {
   	return
   }
   
   if err = connection.SshClient(c); err != nil {
   	//api.BadRequest(c, err)
   	return
   }
   ```

   
