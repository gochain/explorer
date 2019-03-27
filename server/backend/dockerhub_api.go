package backend

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

const dockerHubUrl = "https://registry.hub.docker.com"

type Tag struct {
	Name string
}

var (
	LastUpdatedAt    time.Time
	CachedListOfTags []string
)

func GetSolcImageTags() ([]string, error) {
	duration := time.Since(LastUpdatedAt)
	if int(duration.Hours()) < 1 {
		return CachedListOfTags, nil
	}
	//update the list of the tags every hour
	tags, err := GetImageTags("ethereum/solc")
	if err != nil {
		return nil, err
	}
	LastUpdatedAt = time.Now()
	CachedListOfTags = tags
	return CachedListOfTags, nil
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
		if !(strings.Contains(elem.Name, "alpine") || elem.Name == "stable") {
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
