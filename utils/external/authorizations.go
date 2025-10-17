package external

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Mathis-Pain/Forum/utils/logs"
)

// loadEnv charge les variables d'environnement depuis un fichier .env
// Ce sont les variables qu'on utilise pour obtenir les accès GitHub, Google et Discord
func loadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <authorization.go> Erreur à l'ouverture du fichier %s : %v", filename, err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}
	defer file.Close()

	// Lecture ligne par ligne du fichier
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignore les lignes vides et les commentaires (lignes commençant par #)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Séparation de la ligne au premier '=' uniquement
		// Format attendu : CLE=valeur
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Définition de la variable d'environnement
			os.Setenv(key, value)
		}
	}

	// Vérification des erreurs de lecture du fichier
	if err := scanner.Err(); err != nil {
		logMsg := fmt.Sprintf("ERREUR : <authorization.go> Erreur à l'ouverture du fichier %s : %v", filename, err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}

	return nil
}
