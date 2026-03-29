package router

import "strings"

var lambdaRestPrefixes = []string{
	"/2015-03-31/",
	"/2016-08-19/",
	"/2017-03-31/",
	"/2017-10-31/",
	"/2018-10-31/",
	"/2019-09-30/",
	"/2020-01-01/",
	"/2020-06-30/",
	"/2021-01-01/",
	"/2021-10-01/",
	"/2022-01-01/",
	"/2022-07-01/",
	"/2023-01-01/",
	"/2023-07-01/",
	"/2024-01-01/",
	"/2025-01-01/",
}

func IsLambdaRestPath(path string) bool {
	for _, prefix := range lambdaRestPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}
