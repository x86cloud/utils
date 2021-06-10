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

