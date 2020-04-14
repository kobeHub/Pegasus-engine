package registry;

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

var client *sdk.Client

func InitClient() {
	var err error
	client, err = sdk.NewClientWithAccessKey("cn-hangzhou",
		viper.GetString("AccessKey"), viper.GetString("AccessSecret"))
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Alicloud access key initialize failed")
	}
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
	if len(body) != 0 {
		request.Content = []byte(body)
	}

	return request
}

func GetRepo(imageName string) (error) {
	request := buildRequest("GET", fmt.Sprintf("/repos/pegasus-registry/%s", imageName), "")
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		log.WithFields(log.Fields{
			"Fatal": err,
		}).Error("Alicloud request process error")
		return err
	}
	fmt.Print(response.GetHttpContentString())
	return nil
}
