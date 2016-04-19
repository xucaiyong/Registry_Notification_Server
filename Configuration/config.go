package Configuration

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type Server_Config struct{
	Address string    `yaml:"address,omitempty"`
	Port    uint      `yaml:"port,omitempty"`
}

type Db_Config  struct{
	Db_info mgo.DialInfo `yaml:"dial_info,omitempty"`	
	Collection string  `yaml:"collection,omitempty"`
}

type Config struct{
	Server Server_Config `yaml:"server,omitempty"`
	Search_user Db_Config `yaml:"search_user,omitempty"`
	Search_repo Db_Config `yaml:"search_repo,omitempty"`
	Analysis_config Db_Config `yaml:"analysis_config,omitempty"`
	Mongo_auth  Db_Config `yaml:"mongo_auth,omitempty"`
}

// GetEndpointConnectionString builds and returns a string with the IP and port
// separated by a colon. Nothing special but anyway.
func (c Config) GetEndpointConnectionString() string {
	return fmt.Sprintf("%s:%d", c.Server.Address, c.Server.Port)
}

//validate the server configuration!
func validate_server(c *Server_Config) error{
	if c.Address == "" {
		return fmt.Errorf("Server address must not be empty")
	} 
	return nil
}

//validate the MongoDB configuration!
func validate_db(c *Db_Config) error{
	if len(c.Db_info.Addrs) == 0 {
		return fmt.Errorf("db.addrs must not be empty")
	}
	if c.Db_info.Timeout == 0 {
		c.Db_info.Timeout = 10 * time.Second	
	}
	if c.Db_info.Database == "" {
		return fmt.Errorf("db.database must not be empty")
	}
        if c.Db_info.Username == "" {
		return fmt.Errorf("db.Username must not be empty")
	}
        if c.Db_info.Password  == "" {
		return fmt.Errorf("db.password_file must not be empty")
	}
	if c.Collection == "" {
		return fmt.Errorf("db.collection is required")
	}
	return nil
}
 
//validate the configuration
func validate(c *Config) error {
	if err := validate_server(&c.Server); err != nil {
		return err
	}
	if err := validate_db(&c.Search_user); err != nil {
		return err
	}
	if err:= validate_db(&c.Search_repo); err != nil {
		return err
	}
	if err:= validate_db(&c.Analysis_config); err != nil {
		return err
	}
	if err:= validate_db(&c.Mongo_auth); err != nil {
		return err
	}
	return nil
}

// LoadConfig parses all flags from the command line and returns
// an initialized Settings object and an error object if any. For instance if it
// cannot find the SSL certificate file or the SSL key file it will set the
// returned error appropriately.
func LoadConfig(path string) (*Config, error) {
	//fmt.Println("Starting Loading configuration!")
	c := &Config{}
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read config %s: %s", path, err)
	}
	if err = yaml.Unmarshal(contents, c); err != nil {
		return nil, fmt.Errorf("Failed to parse config: %s", err)
	}
	if err = validate(c); err != nil {
		return nil, fmt.Errorf("Invalid config: %s", err)
	}
	//fmt.Println("Loading configuration done!")
	return c, nil
}

func test_server(c *Server_Config){
	fmt.Printf("%s\n",c.Address)
	fmt.Printf("%d\n",c.Port)
}
func test_db(c *Db_Config){
	fmt.Printf("%s\n",c.Db_info.Addrs)
	fmt.Printf("%d\n",c.Db_info.Timeout)
	fmt.Printf("%s\n",c.Db_info.Database)
	fmt.Printf("%s\n",c.Db_info.Username)
	fmt.Printf("%s\n",c.Db_info.Password)
	fmt.Printf("%s\n",c.Collection)
}

/*
func main(){
	path := "/root/Notification_Server/src/github.com/duyanghao/Registry_Notification_Server/config.yml"
	config, err := LoadConfig(path)
	if err !=nil {
		fmt.Printf("%s\n",err)
		return
	}
	fmt.Println("start config test")
	//server config test
	fmt.Println("server config test:")
	test_server(&config.Server)
	//db config test
	//test user db
	fmt.Println("user db config test:")
	test_db(&config.Search_user)
	//test repo db
	fmt.Println("repo db config test:")
	test_db(&config.Search_repo)
	//test analysis db
	fmt.Println("analysis db config test:")
	test_db(&config.Analysis_config)
	//test docker auth db
	fmt.Println("docker auth db config test:")
	test_db(&config.Mongo_auth)
	
	fmt.Println("end config test")
	
}
*/

