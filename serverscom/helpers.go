package serverscom

func expandIntList(elements []interface{}) []int {
	expandedIntList := make([]int, len(elements))
	for i, v := range elements {
		expandedIntList[i] = v.(int)
	}

	return expandedIntList
}

func expandedStringList(elements []interface{}) []string {
	expandedStringList := make([]string, len(elements))
	for i, v := range elements {
		expandedStringList[i] = v.(string)
	}

	return expandedStringList
}
