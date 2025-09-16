package builddb

import (
	"fmt"
	"io"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// BackupDB crée une copie du fichier DB existant
func BackupDB(dbPath string) error {
	backDbPath := "backforum.db"
	// Vérifier si le fichier existe
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// Pas de fichier à sauvegarder
		return nil
	}
	// Vérifier si un backup existe déjà
	if _, err := os.Stat(backDbPath); err == nil {
		fmt.Println("Un backup existe déjà :", backDbPath)
		return nil
	}

	backupPath := "back" + dbPath // ex: forum.db.bak

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
