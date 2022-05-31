module omo.msa.acm

go 1.15

replace (
	golang.org/x/text => github.com/golang/text v0.3.3
	google.golang.org/grpc => github.com/grpc/grpc-go v1.26.0
)

require (
	github.com/casbin/casbin/v2 v2.7.2
	github.com/labstack/gommon v0.3.0
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/config/source/consul/v2 v2.9.1
	github.com/micro/go-plugins/logger/logrus/v2 v2.9.1
	github.com/micro/go-plugins/registry/consul/v2 v2.9.1
	github.com/micro/go-plugins/registry/etcdv3/v2 v2.9.1
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/tidwall/gjson v1.6.0
	github.com/xtech-cloud/omo-msp-acm v1.3.4
	github.com/xtech-cloud/omo-msp-status v1.0.1
	go.mongodb.org/mongo-driver v1.4.6
	google.golang.org/protobuf v1.24.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)
