package config

var Live = []byte(`
{  
  "env":"live",
  "server":{  
    "port":8080,
    "version":"0.0.1"
  },
  "database":{  
    "host":"localhost:27017",
    "name":"banner",
    "user":"luciano",
    "password":"luciano",
	"timeout":10,
    "enabled":false
  }
}
`)
