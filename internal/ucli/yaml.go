package ucli

import "fmt"

func printYAMLList(key string, values []string, indent string) {
	if len(values) == 0 {
		fmt.Printf("%s%s: []\n", indent, key)
		return
	}

	fmt.Printf("%s%s:\n", indent, key)
	for _, v := range values {
		fmt.Printf("%s  - %s\n", indent, v)
	}
}
