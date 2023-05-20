package gh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// map[avatar_url:https://avatars.githubusercontent.com/u/58223424?v=4 bio:I am full stack web and mobile application developer interested in TS, flutter and all the latest web technologies blog:shreyascodes.tech collaborators:0 company:<nil> created_at:2019-11-26T14:54:22Z disk_usage:7260 email:<nil> events_url:https://api.github.com/users/shreyassanthu77/events{/privacy} followers:1 followers_url:https://api.github.com/users/shreyassanthu77/followers following:4 following_url:https://api.github.com/users/shreyassanthu77/following{/other_user} gists_url:https://api.github.com/users/shreyassanthu77/gists{/gist_id} gravatar_id: hireable:<nil> html_url:https://github.com/shreyassanthu77 id:5.8223424e+07 location:India login:shreyassanthu77 name:Shreyas Mididoddi node_id:MDQ6VXNlcjU4MjIzNDI0 organizations_url:https://api.github.com/users/shreyassanthu77/orgs owned_private_repos:7 plan:map[collaborators:0 name:pro private_repos:9999 space:9.76562499e+08] private_gists:2 public_gists:3 public_repos:14 received_events_url:https://api.github.com/users/shreyassanthu77/received_events repos_url:https://api.github.com/users/shreyassanthu77/repos site_admin:false starred_url:https://api.github.com/users/shreyassanthu77/starred{/owner}{/repo} subscriptions_url:https://api.github.com/users/shreyassanthu77/subscriptions total_private_repos:7 twitter_username:Shreyassanthu77 two_factor_authentication:false type:User updated_at:2023-05-14T04:16:54Z url:https://api.github.com/users/shreyassanthu77]

type User struct {
	Email string `json:"email"`
	Id    int    `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

func request(
	AUTH_TOKEN string,
	method string,
	path string,
	data any,
	body any,
) (*int, error) {

	var reader io.Reader

	if body != nil {
		json_data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(json_data)
	}

	req, err := http.NewRequest(
		method,
		fmt.Sprintf("https://api.github.com/%s", path),
		reader,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+AUTH_TOKEN)
	req.Header.Set("User-Agent", "shreyassanthu77")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&data)

	return &resp.StatusCode, err
}

func Get_user(AUTH_TOKEN string) (*User, error) {

	var data *User

	status, err := request(
		AUTH_TOKEN,
		"GET",
		"user",
		&data,
		nil,
	)

	if err != nil {
		return nil, err
	}

	if *status != 200 && err == nil {
		return nil, fmt.Errorf("error: User not found")
	}

	if data.Email == "" {
		data.Email = fmt.Sprintf("%d+%s@users.noreply.github.com", data.Id, data.Login)
	}

	return data, err
}

func Get_Repo(
	AUTH_TOKEN string,
	user string,
	repo string,
) (map[string]interface{}, error) {
	var data map[string]interface{}

	status, err := request(
		AUTH_TOKEN,
		"GET",
		fmt.Sprintf("repos/%s/%s", user, repo),
		&data,
		nil,
	)

	if *status != 200 && err == nil {
		return nil, fmt.Errorf("error: Repo not found")
	}

	return data, err
}

func Create_repo(
	AUTH_TOKEN string,
	user string,
	owner string,
	name string,
	private bool,
) (map[string]interface{}, error) {
	var data map[string]interface{}

	var path string

	if user == owner {
		path = "user/repos"
	} else {
		path = fmt.Sprintf("orgs/%s/repos", owner)
	}

	status, err := request(
		AUTH_TOKEN,
		"POST",
		path,
		&data,
		map[string]interface{}{
			"name":    name,
			"private": private,
		},
	)

	if *status != 201 && err == nil {
		return nil, fmt.Errorf("error: Repo not created")
	}

	return data, err
}
