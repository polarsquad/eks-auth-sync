package mapping

import (
	"gopkg.in/yaml.v2"
	k8sv1 "k8s.io/api/core/v1"
)

const (
	nodeUsername = "system:node:{{EC2PrivateDNSName}}"
	fargateNodeUsername = "system:node:{{SessionName}}"
)

var (
	nodeGroups = []string{
		"system:bootstrappers",
		"system:nodes",
	}
	fargateNodeGroups = append(nodeGroups, "system:node-proxier")
)

type Role struct {
	RoleARN  string   `yaml:"rolearn"`
	Username string   `yaml:"username"`
	Groups   []string `yaml:"groups"`
}

type User struct {
	UserARN  string   `yaml:"userarn"`
	Username string   `yaml:"username"`
	Groups   []string `yaml:"groups"`
}

func Node(roleARN string) *Role {
	return &Role{
		RoleARN:  roleARN,
		Username: nodeUsername,
		Groups:   nodeGroups,
	}
}

func FargateNode(roleARN string) *Role {
	return &Role{
		RoleARN:  roleARN,
		Username: fargateNodeUsername,
		Groups:   fargateNodeGroups,
	}
}

type All struct {
	Users []*User `yaml:"users"`
	Roles []*Role `yaml:"roles"`
}

func (m *All) IsEmpty() bool {
	return len(m.Users) == 0 && len(m.Roles) == 0
}

func (m *All) ToConfigMap() (*k8sv1.ConfigMap, error) {
	usersStr, err := toYAMLString(m.Users)
	if err != nil {
		return nil, err
	}
	rolesStr, err := toYAMLString(m.Roles)
	if err != nil {
		return nil, err
	}

	cm := &k8sv1.ConfigMap{}
	cm.Name = "aws-auth"
	cm.Namespace = "kube-system"
	cm.Data = map[string]string{
		"mapUsers": usersStr,
		"mapRoles": rolesStr,
	}
	return cm, nil
}

func (m *All) FromYAML(bs []byte) error {
	return yaml.Unmarshal(bs, m)
}

func (m *All) Append(mappings *All) {
	m.Users = append(m.Users, mappings.Users...)
	m.Roles = append(m.Roles, mappings.Roles...)
}

func toYAMLString(o interface{}) (string, error) {
	data, err := yaml.Marshal(o)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
