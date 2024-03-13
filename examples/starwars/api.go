package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	baseEndpoint       = "https://starwars-databank-server.vercel.app/api/v1/"
	charactersEndpoint = "https://starwars-databank-server.vercel.app/api/v1/characters"
)

type onlineMsg bool

type character struct {
	ID          string `json:"_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image"`
}

type charactersResponseMsg struct {
	Info struct {
		Total int    `json:"total"`
		Limit int    `json:"limit"`
		Next  string `json:"next"`
		Prev  string `json:"prev"`
	} `json:"info"`

	Data []character `json:"data"`
}

type errorResponseMsg string

type swDbApi struct {
	page, limit int
}

func newSwDbApi(page, limit int) *swDbApi {
	return &swDbApi{
		page, limit,
	}
}

func (api *swDbApi) ping() tea.Msg {
	resp, err := http.Get(charactersEndpoint)
	if err != nil {
		return onlineMsg(false)
	}

	return onlineMsg(resp.StatusCode == http.StatusOK)
}

func (api *swDbApi) getCharacters() tea.Msg {
	baseURL, err := url.Parse(charactersEndpoint)
	if err != nil {
		return errorResponseMsg(err.Error())
	}

	query := baseURL.Query()

	query.Set("page", strconv.Itoa(api.page))
	query.Set("limit", strconv.Itoa(api.limit))
	baseURL.RawQuery = query.Encode()

	resp, err := http.Get(baseURL.String())
	if err != nil {
		return errorResponseMsg(err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errorResponseMsg(err.Error())
	}

	var response charactersResponseMsg
	err = json.Unmarshal(body, &response)
	if err != nil {
		return response
	}

	// Add an artificial lag in order to admire the loading spinner :)
	time.Sleep(1 * time.Second)

	return response
}
