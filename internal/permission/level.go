package permission

type Level string

const (
	LevelRead        Level = "read"
	LevelWrite       Level = "write"
	LevelDestructive Level = "destructive"
	LevelExecute     Level = "execute"
)
