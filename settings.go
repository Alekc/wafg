package waf

type Settings struct {
	ListenPort int
}

func loadSettings() Settings{
	obj := Settings{
		ListenPort:80,
	}
	return obj
}