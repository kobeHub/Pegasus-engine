package registry

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var client *sdk.Client

var NS string

func InitClient() {
	var err error
	client, err = sdk.NewClientWithAccessKey("cn-hangzhou",
		viper.GetString("AccessKey"), viper.GetString("AccessSecret"))
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Alicloud access key initialize failed")
	}
	NS = viper.GetString("RepoNamespace")
	log.Info("Alicloud sdk client initilize successfully")
}

// Build basic method and body
func buildRequest(method, pathPattern, body string) *requests.CommonRequest {
	request := requests.NewCommonRequest()
	request.Method = method
	request.Scheme = "https" // https | http
	request.Domain = "cr.cn-hangzhou.aliyuncs.com"
	request.Version = "2016-06-07"
	request.PathPattern = pathPattern
	request.Headers["Content-Type"] = "application/json"
	request.Content = []byte(body)

	return request
}

// Get Repo details
func GetRepo(imageName string) (Repo, error) {
	var res Repo
	request := buildRequest("GET", fmt.Sprintf("/repos/%s/%s", NS, imageName), "")
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		log.WithFields(log.Fields{
			"Fatal": err,
		}).Error("Alicloud request process error")
		return res, err
	}
	raw := response.GetHttpContentString()
	res = parseRepo(gjson.Get(raw, "data.repo"))
	return res, nil
}

// Get repositorys in some page
func GetRepoList(page string) (RepoList, error) {
	var res RepoList
	request := buildRequest("GET", "/repos/"+NS, "")
	request.QueryParams["PageSize"] = "10"
	request.QueryParams["Page"] = page

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		log.WithFields(log.Fields{
			"Fatal": err,
		}).Error("Alicloud request process error")
		return res, err
	}
	data := gjson.Get(response.GetHttpContentString(), "data")
	pageSize := int32(data.Get("pageSize").Int())
	res = RepoList{
		Page:     int32(data.Get("page").Int()),
		PageSize: pageSize,
		Total:    int32(data.Get("total").Int()),
		Repos:    make([]Repo, 0, pageSize),
	}
	for _, repo := range data.Get("repos").Array() {
		res.Repos = append(res.Repos, parseRepo(repo))
	}
	return res, nil
}

// Get repos tags list
func GetRepoTags(repoName, page string) (TagsList, error) {
	var res TagsList
	request := buildRequest("GET", fmt.Sprintf("/repos/%s/%s/tags", NS, repoName), "")
	request.QueryParams["PageSize"] = "10"
	request.QueryParams["Page"] = page

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		log.WithFields(log.Fields{
			"Fatal": err,
		}).Error("Alicloud request process error")
		return res, err
	}

	data := gjson.Get(response.GetHttpContentString(), "data")
	pageSize := int32(data.Get("pageSize").Int())
	res = TagsList{
		Page:     int32(data.Get("page").Int()),
		PageSize: pageSize,
		Total:    int32(data.Get("total").Int()),
		Tags:     make([]Tag, 0, pageSize),
	}
	for _, tag := range data.Get("tags").Array() {
		res.Tags = append(res.Tags, parseTag(tag))
	}
	return res, nil
}

// Create a repo with build options, return repoId
func CreateRepo(name, summary string, isOverSea, disableCache bool) (int64, error) {
	var res int64
	body := fmt.Sprintf(`{
"Repo": {
    "Region": "cn-hangzhou",
    "RepoName": "%s",
    "RepoType": "PUBLIC",
    "Summary": "%s",
    "RepoNamespaceName": "%s",
    "RepoNamespace": "%s",
    "RepoBuildType": "AUTO_BUILD"
  },
  "RepoSource": {
    "Source": {
      "SourceRepoType": "GITHUB",
      "SourceRepoNamespace": "%s",
      "SourceRepoName": "%s"
    },
    "BuildConfig": {
      "IsAutoBuild": true,
      "IsOversea": %v,
      "IsDisableCache": %v
     }
   }
}`, name, summary, NS, NS, viper.GetString("GitAccount"),
		viper.GetString("GitSource"), isOverSea, disableCache)
	request := buildRequest("PUT", "/repos", body)

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		log.WithFields(log.Fields{
			"Fatal": err,
		}).Error("Alicloud request process error")
		return res, err
	}
	return gjson.Get(response.GetHttpContentString(), "data.repoId").Int(), nil
}

// Build rules
func GetRepoBuildRule(repoName string) ([]BuildRule, error) {
	var res []BuildRule
	request := buildRequest("GET", fmt.Sprintf("/repos/%s/%s/rules", NS, repoName), "")

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		log.WithFields(log.Fields{
			"Fatal": err,
		}).Error("Alicloud request process error")
		return res, err
	}
	raw := response.GetHttpContentString()
	for _, rule := range gjson.Get(raw, "data.buildRules").Array() {
		res = append(res, parseBuildRule(rule))
	}
	return res, nil
}

func CreateRepoBuildRule(repoName, location, tag string) (int64, error) {
	body := fmt.Sprintf(`
{
  "BuildRule": {
    "PushType": "GIT_BRANCH",
    "PushName": "master",
    "DockerfileLocation": "%s",
    "DockerfileName": "Dockerfile",
    "Tag": "%s"
  }
}`, location, tag)
	request := buildRequest("PUT", fmt.Sprintf("/repos/%s/%s/rules", NS, repoName), body)

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		log.WithFields(log.Fields{
			"Fatal": err,
		}).Error("Alicloud request process error")
		return 0., err
	}
	return gjson.Get(response.GetHttpContentString(), "data.buildRuleId").Int(), nil
}

func StartRepoBuildByRule(repoName, buildRuleId string) (bool, error) {
	request := buildRequest("PUT", fmt.Sprintf("/repos/%s/%s/rules/%v/build", NS, repoName, buildRuleId), "")

	_, err := client.ProcessCommonRequest(request)
	if err != nil {
		log.WithFields(log.Fields{
			"Fatal": err,
		}).Error("Alicloud request process error")
		return false, err
	}
	return true, nil
}

func DeleteRepoBuildRule(repoName, buildRuleId string) (bool, error) {
	request := buildRequest("DELETE", fmt.Sprintf("/repos/%s/%s/rules/%v", NS, repoName, buildRuleId), "")

	_, err := client.ProcessCommonRequest(request)
	if err != nil {
		log.WithFields(log.Fields{
			"Fatal": err,
		}).Error("Alicloud request process error")
		return false, err
	}
	return true, nil
}

func DeleteImage(repoName, tag string) (bool, error) {
	request := buildRequest("DELETE", fmt.Sprintf("/repos/%s/%s/tags/%s", NS, repoName, tag), "")

	_, err := client.ProcessCommonRequest(request)
	if err != nil {
		log.WithFields(log.Fields{
			"Fatal": err,
		}).Error("Alicloud request process error")
		return false, err
	}
	return true, nil
}

func DeleteRepo(repoName string) (bool, error) {
	request := buildRequest("DELETE", fmt.Sprintf("/repos/%s/%s", NS, repoName), "")

	_, err := client.ProcessCommonRequest(request)
	if err != nil {
		log.WithFields(log.Fields{
			"Fatal": err,
		}).Error("Alicloud request process error")
		return false, err
	}
	return true, nil
}
