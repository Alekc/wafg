package cmds

type Rule struct {
	Name        string  `mapstructure:"name"`
	Description string  `mapstructure:"description"`
	Match       []Match `mapstructure:"match"`
	Action      []Match `mapstructure:"action"`
}

type Match struct {
	Field     string `mapstructure:"field"`
	Operation string `mapstructure:"op"`
	Value     string `mapstructure:"value"`
}
