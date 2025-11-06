package getdata

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Taille maximale autorisée en octets (ici 700 Ko)
const MaxImageSize = 200 * 1024

func GetImg(w http.ResponseWriter, r *http.Request) (string, error) {
	var imagePath string

	// Limite la taille totale du corps HTTP
	r.Body = http.MaxBytesReader(w, r.Body, MaxImageSize)

	// Parse la requête multipart (0 = aucun buffer en mémoire, tout sur disque)
	if err := r.ParseMultipartForm(0); err != nil {
		return "", fmt.Errorf("image trop lourde ou erreur dans la requête : %v", err)
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		// Pas de fichier uploadé : ok, on renvoie ""
		return "", nil
	}
	defer file.Close()

	// Vérifie la taille réelle du fichier
	if handler.Size > MaxImageSize {
		return "", fmt.Errorf("image trop lourde : %.2f Ko", float64(handler.Size)/1024)
	}

	// Lit les premiers octets pour détecter le type MIME
	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		return "", fmt.Errorf("erreur lecture fichier : %v", err)
	}
	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" &&
		filetype != "image/gif" && filetype != "image/svg+xml" {
		return "", fmt.Errorf("format d’image non supporté (JPEG/PNG/GIF/SVG uniquement)")
	}
	file.Seek(0, 0) // remet le curseur au début pour la copie

	// Nettoyage du nom de fichier
	cleanFilename := handler.Filename
	cleanFilename = strings.ReplaceAll(cleanFilename, " ", "_")
	cleanFilename = strings.ReplaceAll(cleanFilename, "'", "")
	cleanFilename = strings.ReplaceAll(cleanFilename, "\"", "")

	// Crée le dossier si nécessaire
	os.MkdirAll("./static/uploads", os.ModePerm)

	// Chemin final pour sauvegarde
	imagePath = fmt.Sprintf("./static/uploads/msg_%d_%s", time.Now().Unix(), cleanFilename)
	dst, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("impossible de créer le fichier : %v", err)
	}
	defer dst.Close()

	// Copie du contenu
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("erreur lors de l'import du fichier : %v", err)
	}

	// Supprime le "." initial pour correspondre à ton chemin static
	imagePath = strings.TrimPrefix(imagePath, ".")
	fmt.Println("Image enregistrée :", imagePath)

	return imagePath, nil
}
