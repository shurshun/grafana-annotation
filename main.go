package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	Version = "0.1.0"
	cliApp  = cli.NewApp()
)

const (
	apiPath string = "api/annotations/graphite"
)

type GraphiteAnnotation struct {
	What string   `json:"what"`
	Tags []string `json:"tags"`
	When int64    `json:"when"`
	Data string   `json:"data"`
}

type jsonAnnotationResponse struct {
	Id      int    `json:id`
	Message string `json:message`
}

func init() {
	cliApp.Version = Version
	cliApp.Name = "annotation-poster"
	cliApp.Usage = "Tool to post graphite annotations to grafana"

	cliApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "data",
			Usage:  "Additional data.",
			EnvVar: "GRAFANA_DATA",
		},
		cli.StringFlag{
			Name:   "what",
			Usage:  "The What item to post.",
			EnvVar: "GRAFANA_WHAT",
		},
		cli.StringFlag{
			Name:   "tags",
			Usage:  "Tags.",
			EnvVar: "GRAFANA_TAGS",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "Bearer Token.",
			EnvVar: "GRAFANA_TOKEN",
		},
		cli.StringFlag{
			Name:   "uri",
			Usage:  "Example: https://some-grafana-host.tld",
			EnvVar: "GRAFANA_URI",
		},
	}
}

func main() {
	cliApp.Action = postAnnotation
	err := cliApp.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func postAnnotation(c *cli.Context) {
	annotation := NewGraphiteAnnotation(c.String("what"),
		strings.Split(c.String("tags"), ","),
		c.String("data"))

	annotation.post(c.String("uri"), c.String("token"))
}

func (a *GraphiteAnnotation) toJson() []byte {
	payload, _ := json.Marshal(a)
	return payload
}

func (a *GraphiteAnnotation) post(url string, token string) {

	completeUrl := fmt.Sprintf("%v/%v", url, apiPath)
	payload := a.toJson()
	log.Debug(string(payload))
	log.Debug(completeUrl)
	req, err := http.NewRequest("POST", completeUrl, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("unable to post to url %v. err=%v", completeUrl, err.Error())
	}
	defer resp.Body.Close()

	log.Info("response Status:", resp.Status)
	log.Info("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	var response jsonAnnotationResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Unable to parse response. response is %v. err=%v", string(body), err.Error())
	}
	if response.Id == 0 {
		// no id in response
		log.Fatalf("error sending annotation. Message is %v", response.Message)
	} else {
		log.Info("response Body:", string(body))
	}
}

func NewGraphiteAnnotation(what string, tags []string, data string) GraphiteAnnotation {
	now := time.Now()
	when := now.Unix()
	log.Debugf("new anotation with %v %v %v %v", what, when, tags, data)

	return GraphiteAnnotation{What: what, When: when, Tags: tags, Data: data}
}
