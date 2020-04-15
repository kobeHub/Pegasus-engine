module pegasus-engine

go 1.13

require (
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.113
	github.com/json-iterator/go v1.1.9
	github.com/julienschmidt/httprouter v1.2.0
	github.com/kobeHub/Pegasus-engine v0.0.0-00010101000000-000000000000
	github.com/rs/xid v1.2.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.6.2
	github.com/tidwall/gjson v1.6.0
	k8s.io/api v0.0.0-20190819141258-3544db3b9e44
	k8s.io/apimachinery v0.0.0-20190817020851-f2f3a405f61d
	k8s.io/client-go v0.0.0-20190819141724-e14f31a72a77
)

replace github.com/kobeHub/Pegasus-engine => ./
