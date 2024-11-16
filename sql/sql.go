package sql

import (
	"embed"
	"io/fs"
)

//go:embed migrations
var migrations embed.FS

func GetMigrations() fs.FS {
	// Access the migrations directory
	//files, err := fs.ReadDir(migrations, "migrations")
	//if err != nil {
	//	log.Fatalf("Failed to read migrations directory: %v", err)
	//}
	//
	//for _, file := range files {
	//	if !file.IsDir() {
	//		// Read each file's content
	//		content, err := fs.ReadFile(migrations, filepath.Join("migrations", file.Name()))
	//		if err != nil {
	//			log.Printf("Failed to read file %s: %v", file.Name(), err)
	//			continue
	//		}
	//
	//		// Print file name and content
	//		fmt.Printf("File: %s\nContent:\n%s\n\n", file.Name(), string(content))
	//	}
	//}

	return migrations
}
