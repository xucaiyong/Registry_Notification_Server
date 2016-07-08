package Handler_event

import (
	"fmt"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/notifications"
	"github.com/duyanghao/Registry_Notification_Server/Configuration"
	"github.com/duyanghao/Registry_Notification_Server/Data_strcut"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"os/exec"
	"strings"
	"sync"
)

var repo_lock = new(sync.Mutex)

//Insert the push Manifest record into MongoDB analysis_notify and provide support for searching function and Migration
func ProcessPushEvent(w http.ResponseWriter, r *http.Request, e notifications.Event, c *Configuration.Config) error {
	if e.Target.MediaType != schema2.MediaTypeManifest {
		return fmt.Errorf("Wrong event.Target.MediaType: \"%s\". Expected: \"%s\"", e.Target.MediaType, schema2.MediaTypeManifest)
	}
	//Insert the repository(+tag) information into MongoDB search_notify to provide support for searching function
	repo_string := strings.Split(e.Target.Repository, "/")
	if len(repo_string) != 2 {
		return fmt.Errorf("Failed to string.Split e.Target.repository")
	}
	session, err := mgo.DialWithInfo(&c.Search_repo.Db_info)
	if err != nil {
		return fmt.Errorf("Failed to create Search_repo MongoDB session: %s", err)
	}
	//collection
	collection := session.DB(c.Search_repo.Db_info.Database).C(c.Search_repo.Collection)

	//lock area
	repo_lock.Lock()
	num, err := collection.Find(bson.M{"user": repo_string[0], "repo": repo_string[1], "tag": e.Target.Tag}).Count()
	if err != nil {
		repo_lock.Unlock()
		return fmt.Errorf("Failed to find record: %s", err)
	}
	if num == 0 {
		tmp_repo := &Data_strcut.Cnt_repo{
			User: repo_string[0],
			Repo: repo_string[1],
			Tag:  e.Target.Tag,
		}
		err = collection.Insert(tmp_repo)
		if err != nil {
			repo_lock.Unlock()
			return fmt.Errorf("Failed to insert push record: %s", err)
		}
	}
	repo_lock.Unlock()

	//Insert the push Manifest record into MongoDB analysis_notify
	session, err = mgo.DialWithInfo(&c.Analysis_config.Db_info)
	if err != nil {
		return fmt.Errorf("Failed to create Analysis_config MongoDB session: %s", err)
	}
	//collection
	collection = session.DB(c.Analysis_config.Db_info.Database).C(c.Analysis_config.Collection)

	repo_tmp := fmt.Sprintf("%s:%s", e.Target.Repository, e.Target.Tag)
	tmp_analysis := &Data_strcut.Cnt_analysis{
		Src:       e.Request.Addr,
		Timestamp: e.Timestamp,
		Action:    e.Action,
		Repo:      repo_tmp,
		User:      e.Actor.Name,
	}
	err = collection.Insert(tmp_analysis)
	if err != nil {
		return fmt.Errorf("Failed to insert push record: %s", err)
	}

	//Live_migration
	cmd := exec.Command("./Migration/migration.sh", repo_tmp)
	err = cmd.Start()
	if err != nil {
		fmt.Printf("ERROR: Migration failure %s\n", err)
	}

	fmt.Println("INFO: Push")
	fmt.Printf("INFO: %s \n", e)
	return nil
}
