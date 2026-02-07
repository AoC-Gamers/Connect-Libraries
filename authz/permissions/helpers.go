package permissions

import "sort"

func getPermissionName(permission uint64, names map[uint64]string) string {
	if name, ok := names[permission]; ok {
		return name
	}
	return "UNKNOWN_PERMISSION"
}

func getAllPermissionNames(mask uint64, names map[uint64]string) []string {
	perms := make([]uint64, 0, len(names))
	for perm := range names {
		if mask&perm != 0 {
			perms = append(perms, perm)
		}
	}

	sort.Slice(perms, func(i, j int) bool {
		return perms[i] < perms[j]
	})

	result := make([]string, 0, len(perms))
	for _, perm := range perms {
		result = append(result, names[perm])
	}
	return result
}

func getRolePermissions(role string, roles map[string]uint64) uint64 {
	if value, ok := roles[role]; ok {
		return value
	}
	return 0
}

func getRoleName(role string, names map[string]string) string {
	if value, ok := names[role]; ok {
		return value
	}
	return "Unknown"
}

func isRoleValid(role string, names map[string]string) bool {
	_, ok := names[role]
	return ok
}
