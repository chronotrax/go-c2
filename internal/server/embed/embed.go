package embed

import "embed"

//go:embed migrations/*.sql
var EmbededMigrations embed.FS
