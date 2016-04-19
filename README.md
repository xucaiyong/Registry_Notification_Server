Registry_Notification_Server
=================

Receive docker registry manifest relevant events and store them in MongoDB for Search function,Live-migration and Analysis

## Usage

```
./Registry_Notification_Server ./Configuration/config.yml
```

#### Required

  * `registry_notification_server` - Binary file(produced by go build)
  * `config.yml` - Configuration file(in the Configuration directory)

## Prerequisites
This server assumes the following:

  * The docker registry v2 only send manifest relevant events to notification server when pull or push action happens,and only one manifest pull or push event to be sent when one image pull or push action happens! This requires us to modify the code of registry v2 for deleting layer relevant events.[Here is my modification](https://github.com/duyanghao/distribution/commit/d4d06b8c413fa99488e210923e013c46508777fb).Now i want to simply explain for myself why i need to do that change.As you know,there are many events,including layer and manifest relevant events,sent by registry v2 when a image pull or push action happens.But After all,we only care about the images pushes and pulls actions,and hope there is only one notification when one push or pull action happens.And manifest is the thing which best represents the image push or pull action.So i do that change to delete layer relevant events and store only manifest relevant events notificaton!
  * This notification server is designed for internal use,so it uses http instead of https,but still,it can be easily expanded to https using [http.ListenAndServeTLS](https://golang.org/pkg/net/http/#Header)
  * This notification server is designed for token Authentication of registry v2,and i use [docker_auth](https://github.com/cesanta/docker_auth) as the token authentication server,so the database backend is MongoDB.

## History

Docker is a fantastic tool, there's no doubt about that. One of the main reasons
for it's success is that is provides a central place from which everybody can
grab pre-built images for famous tools or Linux distribution. This place is
called the [Docker Hub](https://hub.docker.com).

Besides the official Hub the Docker team gives everybody the tools to host their
own hub, aka *registry*. They used to have a registry written in python that has
been around for quite some time. It is often referred to as registry v1 because
it was the first incarnation of the registry (actually, the latest version ever
released was `0.9.1`).

With the advent of the [docker/distribution](http://github.com/docker/distribution)
project, the registry v1 must be considered deprecated. The docker/distribution
project also contains a registry component that is written in Go, just like the
rest of Docker. I read somewhere that the intention is to share more code
between all the different projects and to have the Docker team members be able
to work in a common programming language: Go. This registry component is often
referred to as registry v2.

The HTTP-API to query a registry has changed a lot from v1 to
[v2](https://github.com/docker/distribution/blob/master/docs/spec/api.md).
Although v2 responses contain some of the information that existed in v1
responses, this is only in the for compatibility with older Docker client
versions.

As an alternative to the compatibility information, the docker registry v2 has
introduced the concept of [registry event notifications to HTTP endpoints](https://github.com/docker/distribution/blob/mas
ter/docs/notifications.md).

Let me shamelessly copy some relevant information for you:

> The Registry supports sending webhook notifications in response to events
> happening within the registry. Notifications are sent in response to manifest
> pushes and pulls and layer pushes and pulls. These actions are serialized into
> events. The events are queued into a registry-internal broadcast system which
> queues and dispatches events to Endpoints.
> [...]
> Notifications are sent to endpoints via HTTP requests.

## Goal of this project and architecture
The registry notification server listens for events coming from a docker registry v2,Upon receiving an event,it inspects the event and inserts the pull or push records and repository informations into a Mongo database.The information stored in MongoDB support these functions:

  * You can search repository and tag(docker has not implemented this api when using token authentication)
  * You can search the images pushes and pulls records
  * Live Migration( [docker/migrator](https://github.com/duyanghao/migrator) has indeed implemented migration but not Live-Migration!)

## Example
Search function

1. Search home page(xxx/search/)
![](https://raw.githubusercontent.com/duyanghao/Registry_Notification_Server/master/images/Search_Png/home.png)
2. Search repository page(xxx/search/user/)
![](https://raw.githubusercontent.com/duyanghao/Registry_Notification_Server/master/images/Search_Png/repo.png)
3. Search repository result page(xxx/search/user/login/)
![](https://raw.githubusercontent.com/duyanghao/Registry_Notification_Server/master/images/Search_Png/repo_result.png)
4. Search repository+tag page(xxx/search/user/repo/)
![](https://raw.githubusercontent.com/duyanghao/Registry_Notification_Server/master/images/Search_Png/tag.png)
5. Search repository+tag result page(xxx/search/user/repo/login/)
![](https://raw.githubusercontent.com/duyanghao/Registry_Notification_Server/master/images/Search_Png/tag_result.png)

Analysis function

1. Analysis home page(xxx/analysis/)
![](https://raw.githubusercontent.com/duyanghao/Registry_Notification_Server/master/images/Analysis_Png/home.png)
2. Analysis record page(xxx/analysis/user/)
![](https://raw.githubusercontent.com/duyanghao/Registry_Notification_Server/master/images/Analysis_Png/record.png)
3. Analysis record result page(xxx/analysis/user/login/)
![](https://raw.githubusercontent.com/duyanghao/Registry_Notification_Server/master/images/Analysis_Png/record_result.png)

## Logging Notification Server Output
If you need to log the output from server, add `>> notify_server.log 2>&1 ` to the end of the command shown above to capture the output to a file of your choice.

## Build the tool on your own
These instructions walk you through compiling this project to create a single standalone binary that you can copy and run almost wherever you want.

```
mkdir -p ~/tmp
export GOPATH=~/tmp 
mkdir -p $GOPATH/src/github.com/docker
git clone https://github.com/duyanghao/distribution.git $GOPATH/src/github.com/docker/distribution
git clone https://github.com/duyanghao/Registry_Notification_Server.git $GOPATH/src/github.com/duyanghao/Registry_Notification_Server
go get gopkg.in/mgo.v2
go get gopkg.in/yaml.v2
cd $GOPATH/src/github.com/docker/distribution
git checkout fightingdu
GOPATH=$PWD/Godeps/_workspace:$GOPATH
cd ~/tmp/src/github.com/duyanghao/Registry_Notification_Server
go build

```
## Reference
There is already exists an implement of notification server,you can have good a reference here:[Docker Registry Event Collector](https://github.com/kwk/docker-registry-event-collector)

## Suggestion
It is really stupid to write websites using go,maybe you can use php to designed a good website with the data stored in MongoDB backend,that is a better choice!
