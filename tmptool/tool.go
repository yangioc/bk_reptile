package tmptool

func YearConver(year int, region string) int {
	switch region {
	case "tw":
		return year - 1911
	default:
		return year
	}
}
