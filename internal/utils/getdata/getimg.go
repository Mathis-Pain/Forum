package getdata

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func GetImg(w http.ResponseWriter, r *http.Request) (string, error) {

	var imagePath string
	// Limite la taille maximale du corps de la requête à 20MB
	r.Body = http.MaxBytesReader(w, r.Body, 20<<20) // 20 * 1024 * 1024
	// verifie la taille du fichier reçu 20mb
	err := r.ParseMultipartForm(0)
	if err != nil {
		return "", fmt.Errorf("image trop lourde")

	}

	file, handler, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		// Vérifie le type MIME (lit les 512 prmeier octet du fichier pour obtenir
		// une "signature magique" (magic number) propre au type de fichier (par exemple \xFF\xD8\xFF pour un JPEG).)
		buff := make([]byte, 512)
		file.Read(buff)
		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/gif" && filetype != "image/svg" {
			return "", fmt.Errorf("format d’image non supporté (JPEG/PNG/GIF/SVG uniquement)")
		}
		file.Seek(0, 0)
		// cree le chemin de destination du fichier avant de l'importer
		os.MkdirAll("./static/uploads", os.ModePerm)
		imagePath = fmt.Sprintf("./static/uploads/msg_%d_%s", time.Now().Unix(), handler.Filename)
		dst, err := os.Create(imagePath)
		if err != nil {

			return "", fmt.Errorf("fimage path no created")
		}
		defer dst.Close()
		// copie le contenu d'un fichier uploadé vers un fichier sur le serveur
		_, err = io.Copy(dst, file)
		if err != nil {
			return "", fmt.Errorf("image path no import")
		}
	}
	// pour retirer le point devant le / dans le chemin static
	imagePath = strings.TrimPrefix(imagePath, ".")
	fmt.Println(imagePath)
	return imagePath, nil
}
