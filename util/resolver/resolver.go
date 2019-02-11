package resolver

import (
	"../local"
)

// GetHosts resolves different groups.
func GetHosts(s []string) (hosts []string, err error) {
	for _, v := range s {
		switch v[0:2] {
		case "L@":
			localGroup, err := local.GetHostFromLocalFile(string(v[2:]))
			if err != nil {
				return hosts, err
			}
			hosts = append(hosts, localGroup...)
		default:
			hosts = append(hosts, v)
		}
	}
	return hosts, err
}
