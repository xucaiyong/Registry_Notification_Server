package Handler_event

import (
	"fmt"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/notifications"
	"github.com/duyanghao/Registry_Notification_Server/Configuration"
	"github.com/duyanghao/Registry_Notification_Server/Data_strcut"
	"gopkg.in/mgo.v2"
	"net/http"
)

//Insert the Manifest pull record into MongoDB analysis_notify
func ProcessPullEvent(w http.ResponseWriter, r *http.Request, e notifications.Event, c *Configuration.Config) error {
	if e.Target.MediaType != schema2.MediaTypeManifest {
		return fmt.Errorf("Wrong event.Target.MediaType: \"%s\". Expected: \"%s\"", e.Target.MediaType, schema2.MediaTypeManifest)
	}
	//create MongoDB Session
	session, err := mgo.DialWithInfo(&c.Analysis_config.Db_info)
	if err != nil {
		return fmt.Errorf("Failed to create Analysis_config MongoDB session: %s", err)
	}
	//collection
	collection := session.DB(c.Analysis_config.Db_info.Database).C(c.Analysis_config.Collection)

	repo_tmp := fmt.Sprintf("%s:%s", e.Target.Repository, e.Target.Tag)
	tmp := &Data_strcut.Cnt_analysis{
		Src:       e.Request.Addr,
		Timestamp: e.Timestamp,
		Action:    e.Action,
		Repo:      repo_tmp,
		User:      e.Actor.Name,
	}
	err = collection.Insert(tmp)
	if err != nil {
		return fmt.Errorf("Failed to insert pull record: %s", err)
	}
	fmt.Printf("INFO: pull\n")
	fmt.Printf("INFO: %s \n", e)
	return nil
}
