package resolver

import (
	"fmt"
)

// GetHosts resolves different groups.
func GetHosts(s []string) []string {
	var hosts []string
	for _, v := range s {
		switch v[0:2] {
		case "G@":
			fmt.Println("G@")
		default:
			hosts = append(hosts, v)
		}
	}
	return hosts
}
