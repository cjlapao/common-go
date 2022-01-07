package http_helper

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/language"
)

// DownloadFile Downloads a file from a url
func DownloadFile(url string, filepath string) error {
	resp, err := http.Get(url)

	helper.CheckError(err)

	if resp.StatusCode != http.StatusOK {
		return errors.New("Error downloading file from " + url + " with code " + fmt.Sprint(resp.StatusCode))
	}

	defer resp.Body.Close()

	out, err := os.Create(filepath)

	helper.CheckError(err)

	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return err
}

func GetHttpRequestLang(r *http.Request) language.Locale {
	queryLanguage := r.URL.Query().Get("lang")
	if queryLanguage == "" {
		queryLanguage = "en"
	}
	if r.Header.Get("Accept-Language") != "" {
		queryLanguage = r.Header.Get("Accept-Language")
	}
	var locale language.Locale
	return locale.FromString(queryLanguage)
}

func GetHttpRequestQueryValue(r *http.Request, key string) (interface{}, error) {
	keyExists := r.URL.Query().Has(key)
	if !keyExists {
		return nil, errors.New("key was not found")
	}

	keyValue := r.URL.Query().Get(key)
	return keyValue, nil
}

func GetHttpRequestPaginationQuery(r *http.Request) (skip int64, top int64, sortField string, sortOrder string) {
	sort := GetHttpRequestStrValue(r, "$orderby")
	sortArr := strings.Split(sort, " ")
	if len(sortArr) == 2 {
		sortField = sortArr[0]
		sortOrder = sortArr[1]
	}

	top = GetHttpRequestIntValue(r, "$top")
	skip = GetHttpRequestIntValue(r, "$skip")
	return skip, top, sortField, sortOrder
}

func GetHttpRequestFilterQuery(r *http.Request) (field string, value string) {
	filter := GetHttpRequestStrValue(r, "$filterby")
	filterArr := strings.Split(filter, " ")
	if len(filterArr) == 2 {
		field = filterArr[0]
		value = filterArr[1]
	}

	return field, value
}

func GetHttpRequestIntValue(r *http.Request, key string) int64 {
	queryValue, err := GetHttpRequestQueryValue(r, key)
	if err != nil {
		return -1
	}

	value, err := strconv.Atoi(fmt.Sprintf("%v", queryValue))

	if err != nil {
		return -1
	}

	return int64(value)
}

func GetHttpRequestStrValue(r *http.Request, key string) string {
	value, err := GetHttpRequestQueryValue(r, key)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%v", value)
}

func GetHttpRequestBoolValue(r *http.Request, key string, defValue bool) bool {
	queryValue, err := GetHttpRequestQueryValue(r, key)
	if err != nil {
		return defValue
	}

	value, err := strconv.ParseBool(fmt.Sprintf("%v", queryValue))

	if err != nil {
		return defValue
	}

	return value
}

func GetAuthorizationToken(request http.Header) (string, bool) {
	authHeader := strings.Split(request.Get("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		return "", false
	}

	return authHeader[1], true
}
