package service

import (
	"io/ioutil"
	"net/http"

	"github.com/kobeHub/Pegasus-engine/registry"
)

// Basic imge struct
type repoInfo struct {
	Name string `json:"name,required"`
}

type pageInfo struct {
	Name string `json:"name,omitempty"`
	Page string `json:"pega,required"`
}

type createInfo struct {
	Name         string `json:"name,required"`
	Summary      string `json:"summary,required"`
	IsOverSea    bool   `json:"isOverSea,required"`
	DisableCache bool   `json:"disableCache,required"`
}

type ruleInfo struct {
	RepoName string `json:"repoName,required"`
	Location string `json:"location,required"`
	Tag      string `json:"tag,required"`
}

type buildInfo struct {
	RepoName    string `json:"repoName"`
	BuildRuleId string `json:"buildRuleId"`
}

type imageInfo struct {
	RepoName string `json:"repoName"`
	Tag      string `json:"tag"`
}

func GetRepo(w http.ResponseWriter, r *http.Request) {
	var info repoInfo
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondError(w, apiError{errorBadData, err}, "")
		return
	}

	if err := json.Unmarshal(jsonData, &info); err != nil {
		respondError(w, apiError{errorBadData, err}, "")
		return
	}
	if data, err := registry.GetRepo(info.Name); err != nil {
		respondError(w, apiError{errorInternal, err}, "")
	} else {
		respond(w, data)
	}
}

func GetRepoList(w http.ResponseWriter, r *http.Request) {
	var info pageInfo
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondError(w, apiError{errorBadData, err}, "")
		return
	}

	if err := json.Unmarshal(jsonData, &info); err != nil {
		respondError(w, apiError{errorBadData, err}, "")
		return
	}
	if data, err := registry.GetRepoList(info.Page); err != nil {
		respondError(w, apiError{errorInternal, err}, "")
	} else {
		respond(w, data)
	}
}

func GetRepoTags(w http.ResponseWriter, r *http.Request) {
	var info pageInfo
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondError(w, apiError{errorBadData, err}, "")
		return
	}

	if err := json.Unmarshal(jsonData, &info); err != nil {
		respondError(w, apiError{errorBadData, err}, "")
		return
	}
	if data, err := registry.GetRepoTags(info.Name, info.Page); err != nil {
		respondError(w, apiError{errorInternal, err}, "")
	} else {
		respond(w, data)
	}
}

func CreateRepo(w http.ResponseWriter, r *http.Request) {
	var info createInfo
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondError(w, apiError{errorBadData, err}, "")
		return
	}

	if err := json.Unmarshal(jsonData, &info); err != nil {
		respondError(w, apiError{errorBadData, err}, "")
		return
	}
	if data, err := registry.CreateRepo(info.Name, info.Summary, info.IsOverSea, info.DisableCache); err != nil {
		respondError(w, apiError{errorInternal, err}, "")
	} else {
		respond(w, data)
	}
}

func GetRepoBuildRule(w http.ResponseWriter, r *http.Request) {
	var info repoInfo
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondError(w, apiError{errorBadData, err}, "")
		return
	}

	if err := json.Unmarshal(jsonData, &info); err != nil {
		respondError(w, apiError{errorBadData, err}, "")
		return
	}
	if data, err := registry.GetRepoBuildRule(info.Name); err != nil {
		respondError(w, apiError{errorInternal, err}, "")
	} else {
		respond(w, data)
	}
}

func CreateRepoBuildRule(w http.ResponseWriter, r *http.Request) {
	var info ruleInfo
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondError(w, apiError{errorBadData, err}, "")
		return
	}

	if err := json.Unmarshal(jsonData, &info); err != nil {
		respondError(w, apiError{errorBadData, err}, "")
		return
	}
	if data, err := registry.CreateRepoBuildRule(info.RepoName, info.Location, info.Tag); err != nil {
		respondError(w, apiError{errorInternal, err}, "")
	} else {
		respond(w, data)
	}
}

func StartRepoBuildByRule(w http.ResponseWriter, r *http.Request) {
	var info buildInfo
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondError(w, apiError{errorBadData, err}, "")
		return
	}

	if err := json.Unmarshal(jsonData, &info); err != nil {
		respondError(w, apiError{errorBadData, err}, "")
		return
	}
	if data, err := registry.StartRepoBuildByRule(info.RepoName, info.BuildRuleId); err != nil {
		respondError(w, apiError{errorInternal, err}, "")
	} else {
		respond(w, data)
	}
}

func DeleteRepoBuildRule(w http.ResponseWriter, r *http.Request) {
	var info buildInfo
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondError(w, apiError{errorBadData, err}, "")
		return
	}

	if err := json.Unmarshal(jsonData, &info); err != nil {
		respondError(w, apiError{errorBadData, err}, "")
		return
	}
	if data, err := registry.DeleteRepoBuildRule(info.RepoName, info.BuildRuleId); err != nil {
		respondError(w, apiError{errorInternal, err}, "")
	} else {
		respond(w, data)
	}
}

func DeleteImage(w http.ResponseWriter, r *http.Request) {
	var info imageInfo
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondError(w, apiError{errorBadData, err}, "")
		return
	}

	if err := json.Unmarshal(jsonData, &info); err != nil {
		respondError(w, apiError{errorBadData, err}, "")
		return
	}
	if data, err := registry.DeleteImage(info.RepoName, info.Tag); err != nil {
		respondError(w, apiError{errorInternal, err}, "")
	} else {
		respond(w, data)
	}
}

func DeleteRepo(w http.ResponseWriter, r *http.Request) {
	var info repoInfo
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondError(w, apiError{errorBadData, err}, "")
		return
	}

	if err := json.Unmarshal(jsonData, &info); err != nil {
		respondError(w, apiError{errorBadData, err}, "")
		return
	}
	if data, err := registry.DeleteRepo(info.Name); err != nil {
		respondError(w, apiError{errorInternal, err}, "")
	} else {
		respond(w, data)
	}
}
