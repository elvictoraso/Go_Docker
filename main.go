package main

import (
	"encoding/json"
	"net/http"
	"regexp"
)

// Respuesta para el frontend
type PasswordResponse struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

// Función que calcula la fuerza de la contraseña
func CheckPasswordStrength(password string) (string, string) {
	if len(password) < 6 {
		return "Baja", "Muy corta. Es vulnerable a ataques de fuerza bruta inmediata."
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#\$%\^&\*]`).MatchString(password)

	// Calcular "entropía" basada en criterios cumplidos
	score := 0
	if hasUpper { score++ }
	if hasLower { score++ }
	if hasDigit { score++ }
	if hasSpecial { score++ }

	if len(password) >= 12 && score == 4 {
		return "Alta", "Excelente. Tardaría siglos en ser descifrada."
	} else if len(password) >= 8 && score >= 3 {
		return "Media", "Aceptable, pero puede mejorar agregando más variedad."
	}
	return "Baja", "Débil. Agrega mayúsculas, números o caracteres especiales."
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Simple HTML para la interfaz de usuario
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>Password Security Analyzer</title>
				<style>
					body { font-family: Arial, sans-serif; background: #121212; color: #fff; text-align: center; padding-top: 50px; }
					input { padding: 10px; width: 300px; font-size: 16px; border-radius: 5px; border: none; }
					#result { margin-top: 20px; font-size: 18px; font-weight: bold; }
				</style>
			</head>
			<body>
				<h2>Panel de Ciberseguridad: Analizador de Contraseñas</h2>
				<input type="password" id="pass" placeholder="Ingresa tu contraseña..." oninput="analyze()">
				<div id="result"></div>

				<script>
					function analyze() {
						let pass = document.getElementById('pass').value;
						fetch('/analyze?password=' + encodeURIComponent(pass))
							.then(res => res.json())
							.then(data => {
								let color = '#ff0000';
								if (data.level === 'Alta') {
									color = '#00ff00';
								} else if (data.level === 'Media') {
									color = '#ffaa00';
								}
								document.getElementById('result').innerHTML = '<span style="color:' + color + '">[Fuerza: ' + data.level + ']</span><br><p>' + data.message + '</p>';
							});
					}
				</script>
			</body>
			</html>
		`))
		return
	}

	// API Endpoint para analizar en tiempo real
	if r.Method == "GET" || r.Method == "POST" {
		password := r.URL.Query().Get("password")
		level, message := CheckPasswordStrength(password)
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PasswordResponse{Level: level, Message: message})
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/analyze", handler)
	println("Servidor corriendo en http://localhost:8080 ...")
	http.ListenAndServe(":8080", nil)
}
