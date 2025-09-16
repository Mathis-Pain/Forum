package builddb

import (
	"fmt"
	"io"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// BackupDB crée une copie du fichier DB existant
func BackupDB(dbPath string) error {
	backupPath := "backup_forum.db"

	// Vérifier si le fichier DB existe
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil // rien à sauvegarder
	}

	// Copier le contenu du fichier existant
	srcFile, err := os.Open(dbPath)
	if err != nil {
		return fmt.Errorf("erreur ouverture source: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("erreur création backup: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("erreur copie fichier: %w", err)
	}

	fmt.Println("Backup créé :", backupPath)
	return nil
}
