package backend

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type DockerHubAPI struct {
	lastUpdatedSolcAt    time.Time
	cachedListOfSolcTags []string
}

const dockerHubUrl = "https://registry.hub.docker.com"

type Tag struct {
	Name string
}

func (api *DockerHubAPI) GetSolcImageTags() ([]string, error) {
	duration := time.Since(api.lastUpdatedSolcAt)
	if int(duration.Hours()) < 1 {
		return api.cachedListOfSolcTags, nil
	}
	//update the list of the tags every hour
	tags, err := GetImageTags("ethereum/solc")
	if err != nil {
		return nil, err
	}
	api.lastUpdatedSolcAt = time.Now()
	api.cachedListOfSolcTags = tags
	return api.cachedListOfSolcTags, nil
}

func GetImageTags(imageFullName string) (tags []string, err error) {
	url := dockerHubUrl + "/v1/repositories/" + imageFullName + "/tags"
	var response []Tag
	err = parseJson(url, &response)
	if err != nil {
		return nil, err
	}
	tags = append(tags, filterTags(response)...)
	return tags, nil
}

func filterTags(tags []Tag) (res []string) {
	for _, elem := range tags {
		if !(strings.Contains(elem.Name, "alpine") || elem.Name == "stable" || elem.Name == "nightly") {
			res = append(res, elem.Name)
		}
	}
	return res
}
func parseJson(url string, response interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(response)
	if err != nil {
		return err
	}
	return nil
}
