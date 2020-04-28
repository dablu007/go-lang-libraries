package apiclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dablu007/go-lang-libraries/logger"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// RestExecute : An abstraction to call a rest API
func RestExecute(r *http.Request) (int, string, error) {

	fmt.Println(r.Body.Read)
	client := &http.Client{}
	logger.SugarLogger.Debugf("Invoking API. URI: %s Method: %s", r.URL.String(), r.Method)

	requestedTime := time.Now()
	response, resperr := client.Do(r)
	if resperr != nil {
		logger.SugarLogger.Errorf("Error while invoking API. URI: %s Method: %s StatusCode: %d", r.URL.String(), r.Method, r.Response.StatusCode)
		return http.StatusInternalServerError, "", resperr
	}
	responseTime := time.Since(requestedTime)

	if response == nil {
		logger.SugarLogger.Errorf("No response when invoking API. Resource: %s Method: %s", r.URL.String(), r.Method)
		return http.StatusInternalServerError, "", errors.New("Empty response")
	}

	logger.SugarLogger.Infof(
		"Received response from Resource: %s Method: %s StatusCode: %s. Duration: %s",
		r.URL.String(),
		r.Method,
		response.Status,
		responseTime,
	)

	defer response.Body.Close()

	if response.StatusCode == 401 {
		logger.SugarLogger.Warnf("Unauthorized resource. URL: %s Method: %s", r, r.URL.String(), r.Method)
		// this error message is specific and maybe used for retrying with fresh token
		return response.StatusCode, "", errors.New("Unauthorized resource")
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.SugarLogger.Warnf("Unable to read response body")
		return response.StatusCode, "", errors.New("Unable to read response body")
	}

	responseString := string(responseData)
	logger.SugarLogger.Infof("Response Body is", map[string]string{
		"Body": responseString,
	})

	return response.StatusCode, responseString, nil
}

// CreateJSONRequest : Creates a JSON HTTP request
func CreateJSONRequest(httpMethod, absoluteURL string, accessToken string, requestBody interface{}) (*http.Request, error) {

	headers := make(map[string]string, 1)
	headers["Authorization"] = "Bearer " + accessToken

	u, urlerr := url.ParseRequestURI(absoluteURL)
	if urlerr != nil {
		logger.SugarLogger.Warnf("Invalid zest api request path. ResourcePath: %s", absoluteURL)
		return nil, urlerr
	}

	httpMethod = strings.ToUpper(httpMethod)

	requestJSONBody := ""
	if httpMethod != "GET" && requestBody != nil {
		jsonValue, merr := json.Marshal(requestBody)
		if merr != nil {
			logger.SugarLogger.Errorf("Unable to marshal request body. RequestURL: %s Method: %s", u.String(), httpMethod)
			return nil, errors.New("Unable to marshal request body")
		}
		requestJSONBody = string(jsonValue)
	}

	request, requestError := http.NewRequest(httpMethod, u.String(), strings.NewReader(requestJSONBody))

	if requestError != nil {
		logger.SugarLogger.Errorf("Unable to create request. URL: %s ErrorMessage: %s", u.String(), requestError.Error())
		return nil, requestError
	}

	for k, v := range headers {
		request.Header.Add(k, v)
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Content-Length", strconv.Itoa(len(requestJSONBody)))

	return request, nil
}
