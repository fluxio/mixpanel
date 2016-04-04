// Go API client library to consume mixpanel.com analytics data according
// to the following specification of Mixpanel's Data Export API.
//
// https://mixpanel.com/docs/api-documentation/data-export-api
package mixpanel

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"time"
)

const (
	endpoint           = "https://mixpanel.com/api"
	version            = "2.0"
	apiKeyParam        = "api_key"
	expireParam        = "expire"
	formatParam        = "format"
	sigParam           = "sig"
	formatDefault      = "json"
	missingApiKeyError = "Mixpanel package requires ApiKey to be provided."
)

// Mixpanel Client structure requires api_key and api_secret associated with
// your account
type Client struct {
	ApiKey, ApiSecret string
}

// Mixpanel Client's method which sends a request to the Mixpanel's server
func (m Client) Request(
	methods []string,
	params map[string](interface{}),
	http_method string,
	format string,
) ([]byte, error) {
	if m.ApiKey == "" {
		return nil, errors.New(missingApiKeyError)
	}
	if params == nil {
		params = map[string](interface{}){}
	}
	params[apiKeyParam] = m.ApiKey
	params[expireParam] = time.Now().UTC().Add(time.Duration(time.Minute * 10)).Unix()
	if format == "" {
		params[formatParam] = formatDefault
	} else {
		params[formatParam] = format
	}

	// Removing signature in case it was provided as it should not be present
	// while calculating hash signature as specified by the API:
	//https://mixpanel.com/docs/api-documentation/data-export-api#libs-python
	delete(params, sigParam)

	jsonParams, err := jsonifyParams(params)
	if err != nil {
		return nil, fmt.Errorf("Failed to jsonify params: %v", err)
	}

	sig, err := hashArgs(jsonParams, m.ApiSecret)
	if err != nil {
		return nil, fmt.Errorf("Failed to hash the arguments: %v", err)
	}
	jsonParams[sigParam] = sig

	// Forming a path by appending endpoint, version and methods passed as
	// part of the request, for example, annotations/create
	methodsRequest := []string{endpoint, version}
	methodsRequest = append(methodsRequest, methods...)
	urlRequest, err := url.Parse(strings.Join(methodsRequest, "/") + "/")
	if err != nil {
		return nil, fmt.Errorf("Failed to parse the path: %v", err)
	}

	// Encoding data using UTF-8
	encodedData, err := encodeParams(jsonParams)
	if err != nil {
		return nil, err
	}

	var data io.Reader
	if http_method == "GET" {
		urlRequest.RawQuery = encodedData
		data = nil
	} else {
		data = strings.NewReader(encodedData)
	}
	req, err := http.NewRequest(http_method, urlRequest.String(), data)
	if err != nil {
		return nil, fmt.Errorf("Failed to create a new http.Request: %v", err)
	}

	// Sending a request.
	timeout := time.Duration(2 * time.Minute)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to execute a Mixpanel request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read the response: %v", err)
	}
	return respBody, nil
}

// Helper function which encodes parameters in UTF-8
func encodeParams(params map[string]string) (string, error) {
	result := url.Values{}
	for k, v := range params {
		result.Add(k, v)
	}
	return result.Encode(), nil
}

// Hashing function according to the Data Export API which can be found at
// https://mixpanel.com/docs/api-documentation/data-export-api
func hashArgs(args map[string]string, apiSecret string) (string, error) {
	// Data Export API requires sorting arguments before hashing.
	var order []string
	for k := range args {
		order = append(order, k)
	}
	sort.Strings(order)

	argsJoined := ""
	for _, k := range order {
		argsJoined = fmt.Sprintf("%s%s=%s", argsJoined, k, string(args[k]))
	}

	hashed := md5.Sum([]byte(argsJoined + apiSecret))
	return fmt.Sprintf("%x", hashed), nil
}

// Helper function to iterate through a map of arguments jsonifing them at the same time.
func jsonifyParams(params map[string]interface{}) (map[string]string, error) { //, saveResult func(k, v string)) error {
	jsonifiedParams := map[string]string{}

	// Mixpanel API requires individual strings not to be JSONified, i.e param1=value1
	// instead of param1="value1"
	for k, v := range params {
		switch reflect.TypeOf(v).Kind() {
		case reflect.String:
			helper, _ := v.(string)
			jsonifiedParams[k] = helper
		default:
			bytes, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("Failed to Marshal parameters: %v", err)
			}
			jsonifiedParams[k] = string(bytes)
		}
	}
	return jsonifiedParams, nil
}
