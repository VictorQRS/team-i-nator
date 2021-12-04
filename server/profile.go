package main

type Profile int64

const (
	BackEnd Profile     = 1
	Web Profile         = 2
	Mobile Profile      = 3
	UX Profile          = 4
	DataScience Profile = 5
)

func toProfile(n int) Profile {
	switch n {
	case 1: return BackEnd
	case 2: return Web
	case 3: return Mobile
	case 4: return UX
	case 5: return DataScience
	default: panic("unknown profile")
	}
}