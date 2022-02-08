package configuration

import (
	"fmt"
	"strconv"

	"github.com/ydb-platform/ydb-kubernetes-operator/api/v1alpha1"
	"github.com/ydb-platform/ydb-kubernetes-operator/internal/configuration/schema"
	"gopkg.in/yaml.v3"
)

func generate(cr *v1alpha1.Storage) schema.Configuration {
	var hosts []schema.Host

	for i := 0; i < int(cr.Spec.Nodes); i++ {
		hosts = append(hosts, schema.Host{
			Host:         fmt.Sprintf("%v-%d", cr.GetName(), i),
			HostConfigID: 1, // TODO
			NodeID:       i + 1,
			Port:         v1alpha1.InterconnectPort,
			WalleLocation: schema.WalleLocation{
				Body:       12340 + i,
				DataCenter: "az-1",
				Rack:       strconv.Itoa(i),
			},
		})
	}

	return schema.Configuration{
		Hosts: hosts,
	}
}

func Build(cr *v1alpha1.Storage) (map[string]string, error) {
	var crdConfig map[string]interface{}
	generatedConfig := generate(cr)

	err := yaml.Unmarshal([]byte(cr.Spec.Configuration), &crdConfig)
	if err != nil {
		return nil, err
	}

	crdConfig["hosts"] = generatedConfig.Hosts

	data, err := yaml.Marshal(crdConfig)
	if err != nil {
		return nil, err
	}

	result := string(data)

	return map[string]string{
		v1alpha1.ConfigFileName: result,
	}, nil
}
