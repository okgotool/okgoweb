package oktools

func GetMapByArray(pool map[string]string, roles []string) (map[string]string, error) {
	receivers := map[string]string{}

	for _, role := range roles {
		user, ok := pool[role]
		if ok {
			receivers[role] = user
		}
	}

	return receivers, nil
}

func HasKeysInLabels(labels map[string]string, keys []string) bool {
	for _, k := range keys {
		if _, ok := labels[k]; ok {
			return true
		}
	}

	return false
}
