package stub

import (
	"gitlab.com/polarsquad/eks-auth-sync/test/testdata"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
)

const (
	testUsersListMarker = "moreusers"
	testUserTagsMarker  = "moreusertags"
	testRolesListMarker = "moreusers"
	testRoleTagsMarker  = "moreusertags"
)

type IAM struct {
	iamiface.IAMAPI
	PathPrefix string
	Users      []*iam.User
	Roles      []*iam.Role
}

func NewIAM() *IAM {
	return &IAM{
		PathPrefix: testdata.PathPrefix,
		Users:      testdata.Users,
		Roles:      testdata.Roles,
	}
}

func (i *IAM) ListUsers(input *iam.ListUsersInput) (output *iam.ListUsersOutput, err error) {
	output = &iam.ListUsersOutput{
		IsTruncated: aws.Bool(false),
	}

	// If the path prefix doesn't match, don't return anything.
	if *input.PathPrefix != i.PathPrefix {
		return
	}

	// If there's 2 or less users given, just return them all.
	if len(i.Users) <= 2 {
		output.Users = i.Users
		return
	}

	// If the marker is set, return the last two entries
	if input.Marker != nil && *input.Marker == testUsersListMarker {
		output.Users = i.Users[2:]
		return
	}

	// Otherwise, return the "remaining" entries (i.e. everything excluding the first two) with the marker
	output.IsTruncated = aws.Bool(true)
	output.Marker = aws.String(testUsersListMarker)
	output.Users = i.Users[:2]
	return
}

func (i *IAM) ListUserTags(input *iam.ListUserTagsInput) (output *iam.ListUserTagsOutput, err error) {
	output = &iam.ListUserTagsOutput{
		IsTruncated: aws.Bool(false),
	}
	for _, user := range i.Users {
		if *user.UserName == *input.UserName {
			if len(user.Tags) <= 2 {
				output.Tags = user.Tags
				return
			}
			if input.Marker != nil && *input.Marker == testUserTagsMarker {
				output.Tags = user.Tags[2:]
				return
			}
			output.IsTruncated = aws.Bool(true)
			output.Marker = aws.String(testUserTagsMarker)
			output.Tags = user.Tags[:2]
		}
	}
	return
}

func (i *IAM) ListRoles(input *iam.ListRolesInput) (output *iam.ListRolesOutput, err error) {
	output = &iam.ListRolesOutput{
		IsTruncated: aws.Bool(false),
	}

	// If the path prefix doesn't match, don't return anything.
	if *input.PathPrefix != i.PathPrefix {
		return
	}

	// If there's 2 or less users given, just return them all.
	if len(i.Users) <= 2 {
		output.Roles = i.Roles
		return
	}

	// If the marker is set, return the last two entries
	if input.Marker != nil && *input.Marker == testRolesListMarker {
		output.Roles = i.Roles[2:]
		return
	}

	// Otherwise, return the "remaining" entries (i.e. everything excluding the first two) with the marker
	output.IsTruncated = aws.Bool(true)
	output.Marker = aws.String(testRolesListMarker)
	output.Roles = i.Roles[:2]
	return
}

func (i *IAM) ListRoleTags(input *iam.ListRoleTagsInput) (output *iam.ListRoleTagsOutput, err error) {
	output = &iam.ListRoleTagsOutput{
		IsTruncated: aws.Bool(false),
	}
	for _, role := range i.Roles {
		if *role.RoleName == *input.RoleName {
			if len(role.Tags) <= 2 {
				output.Tags = role.Tags
				return
			}
			if input.Marker != nil && *input.Marker == testRoleTagsMarker {
				output.Tags = role.Tags[2:]
				return
			}
			output.IsTruncated = aws.Bool(true)
			output.Marker = aws.String(testRoleTagsMarker)
			output.Tags = role.Tags[:2]
		}
	}
	return
}
