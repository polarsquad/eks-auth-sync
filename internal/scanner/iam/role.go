package iam

import (
	"fmt"
	"log"
	"strings"

	"gitlab.com/polarsquad/eks-auth-sync/internal/mapping"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
)

func scanIAMRoles(svc iamiface.IAMAPI, accountID string, c *ScanConfig) (ms []*mapping.Role, err error) {
	ms = make([]*mapping.Role, 0, 100)
	tagPrefix := c.TagPrefix()
	groupSeparator := c.GroupSeparator()
	var marker *string

	// This is perhaps unnecessarily defensive but added just in case.
	// If the marker and truncation flag don't work for whatever reason,
	// there's an upper bound for how many times the markers are followed.
	for {
		var input iam.ListRolesInput
		var output *iam.ListRolesOutput

		if strings.TrimSpace(c.PathPrefix) != "" {
			input.PathPrefix = aws.String(c.PathPrefix)
		}
		input.Marker = marker

		output, err = svc.ListRoles(&input)
		if err != nil {
			return
		}

		for _, role := range output.Roles {
			var tags map[string]string
			tags, err = getTagsForRole(svc, role.RoleName, tagPrefix)
			if err != nil {
				return
			}
			m := createRoleMappingFromTags(accountID, *role.RoleName, tags, tagPrefix, groupSeparator)
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

func getTagsForRole(svc iamiface.IAMAPI, rolename *string, tagPrefix string) (tags map[string]string, err error) {
	tags = make(map[string]string, 100)
	var marker *string

	// This is perhaps unnecessarily defensive but added just in case.
	// If the marker and truncation flag don't work for whatever reason,
	// there's an upper bound for how many times the markers are followed.
	for {
		output, err := svc.ListRoleTags(&iam.ListRoleTagsInput{
			Marker:   marker,
			RoleName: rolename,
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

func createRoleMappingFromTags(
	accountID string,
	rolename string,
	tags map[string]string,
	tagPrefix string,
	groupSeparator string,
) *mapping.Role {
	roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", accountID, rolename)

	roleType := getTag(tags, tagPrefix, tagKeyType)
	if roleType == "node" {
		return mapping.Node(roleARN)
	}
	if roleType == "fargateNode" {
		return mapping.FargateNode(roleARN)
	}
	if roleType != "" && roleType != "user" {
		log.Printf("unknown role type [%s] found for role %s", roleType, roleARN)
		return nil
	}

	k8sUsername := getK8sUsername(tags, tagPrefix)
	if k8sUsername == "" {
		return nil
	}
	return &mapping.Role{
		RoleARN:  roleARN,
		Username: k8sUsername,
		Groups:   getK8sGroups(tags, tagPrefix, groupSeparator),
	}
}
