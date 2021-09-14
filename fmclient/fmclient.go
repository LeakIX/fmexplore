package fmclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type FmClient struct {
	url        string
	username   string
	password   string
	token      string
	httpClient *http.Client
}

var defaultLimit = 100
var defaultLimitString = strconv.Itoa(int(defaultLimit))

func (c *FmClient) Dump(database, layout string, writer io.Writer) (err error) {
	encoder := json.NewEncoder(writer)
	offset := 1
	for {
		reply, err := c.GetRecords(database, layout, strconv.Itoa(int(offset)), defaultLimitString)
		if err != nil {
			return err
		}
		for _, record := range reply.Response.Data {
			err = encoder.Encode(&record)
		}
		offset += defaultLimit
		if offset >= reply.Response.DataInfo.TotalRecordCount {
			return nil
		}
	}
}

func (c *FmClient) GetRecords(database, layout, offset, limit string) (FmReply, error) {
	queryUrl := c.url + "/fmi/data/v2/databases/" + database + "/layouts/" + layout + "/records"
	queryUrl += fmt.Sprintf("?_offset=%s&_limit=%s", offset, limit)
	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		return FmReply{}, err
	}
	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return FmReply{}, err
	}
	decoder := json.NewDecoder(resp.Body)
	var reply FmReply
	err = decoder.Decode(&reply)
	if err != nil {
		return FmReply{}, err
	}
	if len(reply.Messages) > 1 {
		return FmReply{}, errors.New("multiple errors")
	}
	if len(reply.Messages) < 1 {
		return FmReply{}, errors.New("unknown errors")
	}
	if reply.Messages[0].Message != "OK" {
		return FmReply{}, errors.New(reply.Messages[0].Message)
	}
	for rowIndex, row := range reply.Response.Data {
		for fieldName, fieldValue := range row.Fields {
			if fieldString, isString := fieldValue.(string); isString {
				var parsedJsonField interface{}
				err = json.Unmarshal([]byte(fieldString), &parsedJsonField)
				if err == nil {
					//it's json
					reply.Response.Data[rowIndex].Fields[fieldName] = parsedJsonField
				}
			}
		}
	}
	return reply, nil
}

func (c *FmClient) AuthDatabase(database string) (err error) {
	queryUrl := c.url + "/fmi/data/v2/databases/" + database + "/sessions"
	jsonPayload, err := json.Marshal(struct{}{})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", queryUrl, bytes.NewReader(jsonPayload))
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.username, c.password)
	req.Header.Add("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(resp.Body)
	var reply FmReply
	err = decoder.Decode(&reply)
	if err != nil {
		return err
	}
	if len(reply.Messages) > 1 {
		return errors.New("multiple errors")
	}
	if len(reply.Messages) < 1 {
		return errors.New("unknown errors")
	}
	if reply.Messages[0].Message != "OK" {
		return errors.New(reply.Messages[0].Message)
	}
	c.token = reply.Response.Token
	return nil
}

func GetFmClient(fullUrl string) *FmClient {
	fmUrl, err := url.Parse(fullUrl)
	if err != nil {
		log.Fatalln(err)
	}
	fmUrl.Path = ""
	username := ""
	password := ""
	if fmUrl.User != nil {
		username = fmUrl.User.Username()
		password, _ = fmUrl.User.Password()
		fmUrl.User = nil
	}
	return &FmClient{
		url:      fmUrl.String(),
		username: username,
		password: password,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
				ResponseHeaderTimeout: 20 * time.Second,
				ExpectContinueTimeout: 20 * time.Second,
			},
			Timeout: 20 * time.Second,
		},
	}
}

func (c *FmClient) GetDatabases() (databases []Database, err error) {
	queryUrl := c.url + "/fmi/data/v2/databases"
	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		return databases, err
	}
	req.SetBasicAuth(c.username, c.password)
	req.Header.Add("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return databases, err
	}
	decoder := json.NewDecoder(resp.Body)
	var reply FmReply
	err = decoder.Decode(&reply)
	if err != nil {
		return databases, err
	}
	if len(reply.Messages) > 1 {
		return databases, errors.New("multiple errors")
	}
	if len(reply.Messages) < 1 {
		return databases, errors.New("unknown errors")
	}
	if reply.Messages[0].Message != "OK" {
		return databases, errors.New(reply.Messages[0].Message)
	}
	return reply.Response.Databases, nil
}

func (c *FmClient) GetLayouts(database string) (layouts []Layout, err error) {
	queryUrl := c.url + "/fmi/data/v2/databases/" + database + "/layouts"
	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		return layouts, err
	}
	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return layouts, err
	}
	decoder := json.NewDecoder(resp.Body)
	var reply FmReply
	err = decoder.Decode(&reply)
	if err != nil {
		return layouts, err
	}
	if len(reply.Messages) > 1 {
		return layouts, errors.New("multiple errors")
	}
	if len(reply.Messages) < 1 {
		return layouts, errors.New("unknown errors")
	}
	if reply.Messages[0].Message != "OK" {
		return layouts, errors.New(reply.Messages[0].Message)
	}
	return reply.Response.Layouts, nil
}

type Database struct {
	Name string `json:"name"`
}

type Layout struct {
	Name     string `json:"name"`
	IsFolder bool   `json:"isFolder"`
}

type FmReply struct {
	Response struct {
		DataInfo struct {
			TotalRecordCount int `json:"totalRecordCount"`
		} `json:"dataInfo"`
		Data []struct {
			Id     string                 `json:"recordId"`
			Fields map[string]interface{} `json:"fieldData"`
		}
		Token     string     `json:"token"`
		Databases []Database `json:"databases"`
		Layouts   []Layout   `json:"layouts"`
	} `json:"response"`
	Messages []struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"messages"`
}

type FmAuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
