package iam

import (
	"fmt"
	"strings"
)

const (
	tagKeyUsername = "username"
	tagKeyGroups   = "groups"
	tagKeyType     = "type"
)

func tagKey(tagPrefix, key string) string {
	return fmt.Sprintf("%s/%s", tagPrefix, key)
}

func getTag(tags map[string]string, tagPrefix string, key string) string {
	return tags[tagKey(tagPrefix, key)]
}

func getK8sUsername(tags map[string]string, tagPrefix string) string {
	return getTag(tags, tagPrefix, tagKeyUsername)
}

func getK8sGroups(tags map[string]string, tagPrefix string, tagDelimiter string) []string {
	return strings.Split(getTag(tags, tagPrefix, tagKeyGroups), tagDelimiter)
}
