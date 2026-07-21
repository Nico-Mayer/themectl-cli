package integration

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func setJSONCString(config, key, value string) (string, error) {
	quoted, _ := json.Marshal(value) // marshaling a string never fails

	re := regexp.MustCompile(`("` + regexp.QuoteMeta(key) + `"\s*:\s*)"[^"]*"`)
	if re.MatchString(config) {
		repl := `${1}` + strings.ReplaceAll(string(quoted), "$", "$$")
		return re.ReplaceAllString(config, repl), nil
	}

	end := strings.LastIndex(config, "}")
	if end < 0 {
		return "", fmt.Errorf("no object found in config")
	}
	head := strings.TrimRight(config[:end], " \t\r\n")
	if head == "" {
		return "", fmt.Errorf("no object found in config")
	}

	sep := ",\n"
	if last := head[len(head)-1]; last == '{' || last == ',' {
		sep = "\n"
	}
	return head + sep + "  \"" + key + "\": " + string(quoted) + "\n" + config[end:], nil
}
