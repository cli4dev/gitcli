package gitlabs

import (
	"encoding/json"
	"errors"
	"fmt"
	xurl "net/url"
	"strings"

	"github.com/micro-plat/lib4go/net/http"
)

var NOT_FOUND = errors.New("服务未找到")

func GetRepositories(url string) ([]*Repository, error) {
	if !strings.Contains(url, "://") {
		url = "https://" + url
	}
	path := make([]*Repository, 0, 1)
	reps, err := getChildren(url)
	if err == NOT_FOUND {
		path = append(path, NewRepository(url))
		return path, nil
	}
	if err != nil {
		return nil, err
	}
	for _, rep := range reps {
		switch rep.Type {
		case "project":
			path = append(path, rep)
		case "group":
			r, err := GetRepositories(rep.FullPath)
			if err != nil {
				return nil, err
			}
			path = append(path, r...)
		}
	}
	return path, nil
}

func getChildren(url string) ([]*Repository, error) {
	u, _ := xurl.Parse(url)
	path := fmt.Sprintf("%s://%s/groups%s/-/children.json", u.Scheme, u.Host, u.Path)
	client, _ := http.NewHTTPClient()
	result, status, err := client.Request("get", path, "", "utf-8", nil)
	if err != nil {
		return nil, err
	}
	if status != 200 || strings.Contains(result, "<!DOCTYPE html>") {
		return nil, NOT_FOUND
	}
	groups := make([]*Repository, 0, 1)
	if err := json.Unmarshal([]byte(result), &groups); err != nil {
		return nil, err
	}
	for _, g := range groups {
		g.FullPath = fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, g.Path)
	}
	return groups, nil
}
