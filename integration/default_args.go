package integration

func doArgDefaulting(providedArgs, defaultArgs []string) []string {
	if len(providedArgs) > 0 {
		return providedArgs
	}
	return defaultArgs
}
