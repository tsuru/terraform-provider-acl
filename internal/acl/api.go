package acl

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const userAgent = "Terraform-Provider-ACL"

func doProxyURLRequest(ctx context.Context, method, fullUrl string, body io.Reader, cli *clientImpl) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, fullUrl, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "bearer "+cli.token)
	req.Header.Set("User-Agent", userAgent)
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode < 200 || rsp.StatusCode >= 400 {
		data, _ := ioutil.ReadAll(rsp.Body)
		rsp.Body.Close()
		return nil, errors.Errorf("invalid status code %d: %q", rsp.StatusCode, string(data))
	}

	return rsp, nil
}

func doProxyRequest(ctx context.Context, method, service, instance, path string, body io.Reader, cli *clientImpl) (*http.Response, error) {
	fullUrl := fmt.Sprintf("%s/services/%s/proxy/%s?callback=%s",
		strings.TrimSuffix(cli.Host, "/"),
		service,
		instance,
		path,
	)
	log.Print("[DEBUG] making request to: ", fullUrl, " method=", method)
	return doProxyURLRequest(ctx, method, fullUrl, body, cli)
}
