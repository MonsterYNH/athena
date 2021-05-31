module athena

go 1.16

require (
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.5
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5
	golang.org/x/sys v0.0.0-20210525143221-35b2ab0089ea // indirect
	google.golang.org/genproto v0.0.0-20200423170343-7949de9c1215
	google.golang.org/grpc v1.38.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect

)

replace (
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5
	github.com/golang/protobuf => github.com/golang/protobuf v1.3.2
	google.golang.org/grpc => google.golang.org/grpc v1.27.0
)
