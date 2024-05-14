package api

import (
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/go-chi/render"
	"github.com/heyztb/lists/internal/models"
)

func HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	this, err := os.Executable()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
	}
	checksum, err := exec.Command("sha256sum", this).Output()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Internal server error",
		})
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, &models.HealthcheckResponse{
		Status:   http.StatusOK,
		Checksum: strings.Split(string(checksum), " ")[0],
	})
}
