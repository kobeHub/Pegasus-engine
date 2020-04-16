package registry

import (
	"fmt"

	"github.com/tidwall/gjson"
)

// Define Repo status
type RepoStatus string

const (
	REPO_NORMAL   RepoStatus = "NORMAL"
	REPO_DELETING RepoStatus = "DELETING"
	REPO_ALL      RepoStatus = "ALL"
)

// Repo present a image repository
type Repo struct {
	Id        int64      `json:"id"`
	Name      string     `json:"name"`
	Summary   string     `json:"summary"`
	Status    RepoStatus `json:"status"`
	Downloads int32      `json:"downloads"`
	Url       string     `json:"url"`
}

// A page list `Repo`s
type RepoList struct {
	Page     int32  `json:"page"`
	PageSize int32  `json:"pageSize"`
	Total    int32  `json:"total"`
	Repos    []Repo `json:"repos"`
}

// Image tags
type Tag struct {
	ImageId string `json:"imageId"`
	Tag     string `json:"tag"`
	Size    int64  `json:"size"`
	Status  string `json:"status"`
	Digest  string `json:"digest"`
}

type TagsList struct {
	Page     int32 `json:"page"`
	PageSize int32 `json:"pageSize"`
	Total    int32 `json:"total"`
	Tags     []Tag `json:"tags"`
}

// Build rule
type BuildRule struct {
	BuildRuleId string `json:"buildRuleId"`
	ImageTag string `json:"imageTag"`
}

// Build from alicloud response
func parseRepo(data gjson.Result) Repo {
	name := data.Get("repoName").String()
	return Repo{
		Id:        data.Get("repoId").Int(),
		Name:      name,
		Summary:   data.Get("summary").String(),
		Status:    RepoStatus(data.Get("repoStatus").String()),
		Downloads: int32(data.Get("downloads").Int()),
		Url: fmt.Sprintf("%s/%s/%s",
			data.Get("repoDomainList.public").String(),
			"pegasus-registry",
			name),
	}
}

func parseTag(data gjson.Result) Tag {
	return Tag{
		ImageId: data.Get("imageId").String(),
		Tag:     data.Get("tag").String(),
		Size:    data.Get("imageSize").Int(),
		Status:  data.Get("status").String(),
		Digest:  data.Get("digest").String(),
	}
}

func parseBuildRule(data gjson.Result) BuildRule {
	return BuildRule {
		BuildRuleId: data.Get("buildRuleId").String(),
		ImageTag: data.Get("imageTag").String(),
	}
}
