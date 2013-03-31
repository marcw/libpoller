package poller

func btou(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

func btos(b bool) string {
	if b {
		return "UP"
	}
	return "DOWN"
}
