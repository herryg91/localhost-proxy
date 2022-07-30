package helpers

import "strings"

func StripUrlPrefix(url string, prefix string) string {
	slash_encode := "{^&!#"
	prefix_encode := strings.ReplaceAll(prefix, "/", slash_encode)
	url_encode := strings.ReplaceAll(url, "/", slash_encode)

	trimmed_url_encode := strings.TrimPrefix(url_encode, prefix_encode)
	trimmed_url := strings.ReplaceAll(trimmed_url_encode, slash_encode, "/")

	return trimmed_url
}
