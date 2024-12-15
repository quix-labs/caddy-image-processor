package CADDY_FILE_SERVER

import (
	"bytes"
	"fmt"
	"github.com/cespare/xxhash/v2"
	"net/url"
	"regexp"
	"sort"
)

func getProcessedImageEtag(initialEtag string, form *url.Values) string {
	// Return early if the initial ETag is empty
	if initialEtag == "" {
		return ""
	}

	re := regexp.MustCompile(`^(W/"|")(.*?)(")$`)
	matches := re.FindStringSubmatch(initialEtag)
	if len(matches) != 4 {
		return initialEtag
	}

	var params []string
	for key, values := range *form {
		params = append(params, key+"="+values[0])
	}

	// Sort the parameters to ensure consistent order
	sort.Strings(params)

	// Use a bytes.Buffer to join parameters efficiently
	var buffer bytes.Buffer
	for _, param := range params {
		buffer.WriteString(param)
	}

	// Generate the hash of the concatenated parameters
	hash := xxhash.New()
	_, err := hash.Write(buffer.Bytes())
	if err != nil {
		return ""
	}
	hashString := fmt.Sprintf("%x", hash.Sum(nil))
	return matches[1] + matches[2] + "-" + hashString + matches[3]
}
