package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// Client implements all functionality necessary to communicate with
// the inngest server.
type Client interface {
	Credentials() []byte

	Login(ctx context.Context, email, password string) ([]byte, error)
	Workspaces(ctx context.Context) ([]Workspace, error)
	Actions(ctx context.Context, workspaceID uuid.UUID) error
	CreateAction(ctx context.Context, config string) (*Action, error)
}

type ClientOpt func(Client) Client

func New(opts ...ClientOpt) Client {
	c := &httpClient{
		Client: http.DefaultClient,
		api:    "https://api.inngest.com",
		ingest: "https://inn.gs",
	}

	for _, o := range opts {
		c = o(c).(*httpClient)
	}

	return c
}

// WithCredentials is used to configure a client with an existing JWT.
func WithAPI(api string) ClientOpt {
	return func(c Client) Client {
		if api == "" {
			return c
		}

		client := c.(*httpClient)
		client.api = api
		return client
	}
}

// WithCredentials is used to configure a client with an existing JWT.
func WithCredentials(creds []byte) ClientOpt {
	return func(c Client) Client {
		client := c.(*httpClient)
		client.creds = creds
		return client
	}
}

// httpClient represents a concrete HTTP implementation of a Client
type httpClient struct {
	*http.Client

	api    string
	ingest string
	creds  []byte
}

func (c httpClient) Credentials() []byte {
	return c.creds
}

func (c httpClient) Login(ctx context.Context, email, password string) ([]byte, error) {
	input := map[string]string{
		"email":    email,
		"password": password,
	}
	buf := jsonBuffer(ctx, input)

	req, err := c.NewRequest(http.MethodPost, "/v1/login", buf)
	if err != nil {
		return nil, fmt.Errorf("error creating login request: %s", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing login request: %s", err)
	}
	defer resp.Body.Close()

	type response struct {
		Error string
		JWT   string
	}

	r := &response{}
	if err = json.NewDecoder(resp.Body).Decode(r); err != nil {
		return nil, fmt.Errorf("invalid json response: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", r.Error)
	}

	return []byte(r.JWT), nil
}