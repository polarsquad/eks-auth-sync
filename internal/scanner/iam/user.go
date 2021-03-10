package iam

import (
	"fmt"
	"strings"

	"gitlab.com/polarsquad/eks-auth-sync/internal/mapping"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
)

func scanIAMUsers(svc iamiface.IAMAPI, accountID string, c *ScanConfig) (ms []*mapping.User, err error) {
	ms = make([]*mapping.User, 0, 100)
	tagPrefix := c.TagPrefix()
	groupSeparator := c.GroupSeparator()
	var marker *string

	// This is perhaps unnecessarily defensive but added just in case.
	// If the marker and truncation flag don't work for whatever reason,
	// there's an upper bound for how many times the markers are followed.
	for i := 0; i < 1000; i++ {
		var input iam.ListUsersInput
		var output *iam.ListUsersOutput

		if strings.TrimSpace(c.PathPrefix) != "" {
			input.PathPrefix = aws.String(c.PathPrefix)
		}
		input.Marker = marker

		output, err = svc.ListUsers(&input)
		if err != nil {
			return
		}

		for _, user := range output.Users {
			var tags map[string]string
			tags, err = getTagsForUser(svc, user.UserName, tagPrefix)
			if err != nil {
				return
			}
			m := createUserMappingFromTags(accountID, *user.UserName, tags, tagPrefix, groupSeparator)
			if m != nil {
				ms = append(ms, m)
			}
		}

		if *output.IsTruncated {
			marker = output.Marker
		} else {
			break
		}
	}
	return
}

func getTagsForUser(svc iamiface.IAMAPI, username *string, tagPrefix string) (tags map[string]string, err error) {
	tags = make(map[string]string, 100)
	var marker *string

	// This is perhaps unnecessarily defensive but added just in case.
	// If the marker and truncation flag don't work for whatever reason,
	// there's an upper bound for how many times the markers are followed.
	for {
		output, err := svc.ListUserTags(&iam.ListUserTagsInput{
			Marker:   marker,
			UserName: username,
		})
		if err != nil {
			return nil, err
		}

		for _, tag := range output.Tags {
			if strings.HasPrefix(*tag.Key, tagPrefix) {
				tags[*tag.Key] = *tag.Value
			}
		}

		if *output.IsTruncated {
			marker = output.Marker
		} else {
			break
		}
	}
	return tags, nil
}

func createUserMappingFromTags(
	accountID string,
	username string,
	tags map[string]string,
	tagPrefix string,
	groupSeparator string,
) *mapping.User {
	k8sUsername := getK8sUsername(tags, tagPrefix)
	if k8sUsername == "" {
		return nil
	}
	return &mapping.User{
		UserARN:  fmt.Sprintf("arn:aws:iam::%s:user/%s", accountID, username),
		Username: k8sUsername,
		Groups:   getK8sGroups(tags, tagPrefix, groupSeparator),
	}
}
