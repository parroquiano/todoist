package todoist_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"regexp"
	"testing"
	"time"

	"github.com/ides15/todoist"
	"github.com/ides15/todoist/types"
)

func TestNewClient_OK(t *testing.T) {
	_, err := todoist.NewClient("12345", nil)
	if err != nil {
		t.Fatalf("expected nil error, received %v", err)
	}
}

func TestNewClient_NilToken(t *testing.T) {
	_, err := todoist.NewClient("", nil)
	if err == nil {
		t.Fatalf("expected err, received %v", err)
	} else if err.Error() != types.ErrRequiredToken.Error() {
		t.Fatalf("expected %v, received %v", types.ErrRequiredToken.Error(), err)
	}
}

func TestNewRequest_OKURL(t *testing.T) {
	Setup()

	request, err := TestClient.NewRequest("*", nil, nil)
	if err != nil {
		t.Fatalf("expected nil error, received %v", err)
	}

	if request.URL.String() != todoist.DefaultBaseURL {
		t.Fatalf("expected %s, received %s", todoist.DefaultBaseURL, request.URL.String())
	}
}

func TestNewRequest_OKToken(t *testing.T) {
	Setup()

	request, _ := TestClient.NewRequest("*", nil, nil)
	defer request.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(request.Body)
	body := string(bodyBytes)

	re := regexp.MustCompile(`[^_]token=([^&\s]+)`)
	matches := re.FindStringSubmatch(body)

	if len(matches) < 1 {
		t.Fatalf("expected a matching token in body, received %s", body)
	} else if matches[1] != "12345" {
		t.Log(body)
		t.Fatalf("expected token %s, received %s", "12345", matches[1])
	}
}

func TestNewRequest_Bad(t *testing.T) {
	Setup()

	// ASCII control character will break `TestClient.NewRequest`
	TestClient.BaseURL = "\t"

	_, err := TestClient.NewRequest("*", nil, nil)
	if err == nil {
		t.Fatalf("expected err, received %v", err)
	}

	if err.Error() != types.ErrBuildRequest.Error() {
		t.Fatalf("expected %v, received %v", err.Error(), types.ErrBuildRequest.Error())
	}
}

func TestNewRequest_SyncToken(t *testing.T) {
	Setup()

	request, _ := TestClient.NewRequest("*", nil, nil)
	defer request.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(request.Body)
	body := string(bodyBytes)

	re := regexp.MustCompile(`sync_token=([^&\s]+)`)
	matches := re.FindStringSubmatch(body)

	if len(matches) < 1 {
		t.Fatalf("expected a matching sync_token in body, received %s", body)
	} else if matches[1] != "%2A" {
		t.Fatalf("expected synx_token '%s', received '%s'", "%2A", matches[1])
	}
}

func TestNewRequest_NilSyncToken(t *testing.T) {
	Setup()

	request, _ := TestClient.NewRequest("", nil, nil)
	defer request.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(request.Body)
	body := string(bodyBytes)

	re := regexp.MustCompile(`sync_token=([^&\s]+)`)
	matches := re.FindStringSubmatch(body)

	if len(matches) > 1 {
		t.Fatalf("expected no sync_token in body, received %s", body)
	}
}

func TestNewRequest_ContentType(t *testing.T) {
	Setup()

	request, _ := TestClient.NewRequest("*", nil, nil)

	expected := "application/x-www-form-urlencoded"
	if request.Header.Get("Content-Type") != expected {
		t.Fatalf("expected Content-Type of %s, received %s", expected, request.Header.Get("Content-Type"))
	}
}

func TestNewRequest_UserAgent(t *testing.T) {
	Setup()

	request, _ := TestClient.NewRequest("*", nil, nil)

	expected := todoist.DefaultUserAgent
	if request.Header.Get("User-Agent") != expected {
		t.Fatalf("expected User-Agent of %s, received %s", expected, request.Header.Get("User-Agent"))
	}
}

func TestNewRequest_Commands(t *testing.T) {
	Setup()

	commands := &[]types.Command{{
		Type: "project_add",
		Args: map[string]string{
			"arg": "test",
		},
		UUID:   "uuid",
		TempID: "tempID",
	}}

	request, _ := TestClient.NewRequest("*", commands, nil)
	defer request.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(request.Body)
	body := string(bodyBytes)

	re := regexp.MustCompile(`commands=([^&\s]+)`)
	matches := re.FindStringSubmatch(body)

	if len(matches) < 1 {
		t.Fatalf("expected matching commands in body, received %s", body)
	}
}

func TestNewRequest_NilCommands(t *testing.T) {
	Setup()

	request, _ := TestClient.NewRequest("*", nil, nil)
	defer request.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(request.Body)
	body := string(bodyBytes)

	re := regexp.MustCompile(`commands=([^&\s]+)`)
	matches := re.FindStringSubmatch(body)

	if len(matches) > 1 {
		t.Fatalf("expected no commands in body, received %s", body)
	}
}

func TestNewRequest_ResourceTypes(t *testing.T) {
	Setup()

	resourceTypes := &[]string{"resource_type"}

	request, _ := TestClient.NewRequest("*", nil, resourceTypes)
	defer request.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(request.Body)
	body := string(bodyBytes)

	re := regexp.MustCompile(`resource_types=([^&\s]+)`)
	matches := re.FindStringSubmatch(body)

	if len(matches) < 1 {
		t.Fatalf("expected matching resource_types in body, received %s", body)
	}
}

func TestNewRequest_NilResourceTypes(t *testing.T) {
	Setup()

	request, _ := TestClient.NewRequest("*", nil, nil)
	defer request.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(request.Body)
	body := string(bodyBytes)

	re := regexp.MustCompile(`resource_types=([^&\s]+)`)
	matches := re.FindStringSubmatch(body)

	if len(matches) > 1 {
		t.Fatalf("expected no resource_types in body, received %s", body)
	}
}

func TestDo_RequestOK(t *testing.T) {
	Setup()

	TestClient.BaseURL = TestServer.URL

	request, _ := TestClient.NewRequest("*", nil, nil)
	_, err := TestClient.Do(context.Background(), request, nil)
	if err != nil {
		t.Fatalf("expected no err, received %v", err)
	}
}

func TestDo_RequestContextCancel(t *testing.T) {
	Setup()

	TestClient.BaseURL = TestServer.URL
	d := time.Now().Add(1 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	cancel()

	request, _ := TestClient.NewRequest("*", nil, nil)
	_, err := TestClient.Do(ctx, request, nil)
	if err == nil {
		t.Fatalf("expected context cancelled error, received %v", err)
	}
}

func TestDo_RequestError(t *testing.T) {
	Setup()

	TestClient.BaseURL = TestServer.URL

	request, _ := TestClient.NewRequest("*", nil, nil)

	// Force error from `TestClient.Do`
	request.URL = nil
	_, err := TestClient.Do(context.Background(), request, nil)
	if err == nil {
		t.Fatalf("expected Request.URL error, received %v", err)
	}
}

func TestDo_AUTH_CSRF_ERRORResponse(t *testing.T) {
	Setup()

	TestClient.BaseURL = TestServer.URL + "/AUTH_CSRF_ERROR"

	request, _ := TestClient.NewRequest("*", nil, nil)
	_, err := TestClient.Do(context.Background(), request, nil)

	if e, ok := err.(*types.HTTPError); ok {
		expected := "AUTH_CSRF_ERROR"
		if e.ErrorTag != expected {
			t.Fatalf("expected %s, received %s", expected, e.ErrorTag)
		}
	} else {
		t.Fatalf("expected *types.HTTPError, received %v (%T)", err, err)
	}
}

func TestDo_AUTH_INVALID_TOKENResponse(t *testing.T) {
	Setup()

	TestClient.BaseURL = TestServer.URL + "/AUTH_INVALID_TOKEN"

	request, _ := TestClient.NewRequest("*", nil, nil)
	_, err := TestClient.Do(context.Background(), request, nil)

	if e, ok := err.(*types.HTTPError); ok {
		expected := "AUTH_INVALID_TOKEN"
		if e.ErrorTag != expected {
			t.Fatalf("expected %s, received %s", expected, e.ErrorTag)
		}
	} else {
		t.Fatalf("expected *types.HTTPError, received %v (%T)", err, err)
	}
}

func TestDo_InvalidErrorResponse(t *testing.T) {
	Setup()

	TestClient.BaseURL = TestServer.URL + "/invalid-error"

	request, _ := TestClient.NewRequest("*", nil, nil)
	_, err := TestClient.Do(context.Background(), request, nil)

	if e, ok := err.(*types.HTTPError); ok {
		t.Fatalf("expected %v, received %v", types.ErrUnknown, e)
	}
}

func TestDo_VIsIOWriter(t *testing.T) {
	Setup()

	TestClient.BaseURL = TestServer.URL

	request, _ := TestClient.NewRequest("*", nil, nil)

	w := bytes.NewBufferString("new buffer string")

	_, err := TestClient.Do(context.Background(), request, w)
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
}

func TestDo_ResBodyHasEOF(t *testing.T) {
	Setup()

	TestClient.BaseURL = TestServer.URL + "/empty-response-body"
	request, _ := TestClient.NewRequest("*", nil, nil)

	res, err := TestClient.Do(context.Background(), request, nil)
	t.Log(res.Raw.Body)
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
}
