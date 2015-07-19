package webapp

import (
	"net/url"
)

type User struct {
	UserId   string `mgo:"_id"`
	UserName string
	Password string
}

type ViewModel struct {
	Model interface{}
	User  UserManager
	Html  *Html
}

func Form2M(values url.Values) map[string]interface{} {
	out := make(map[string]interface{})
	for k, v := range values {
		if len(v) == 0 {
			continue
		}
		if len(v) == 1 {
			out[k] = v[0]
			continue
		}
		out[k] = v
	}
	return out
}
