package main

import (
	"encoding/json"
	"net/http"
	"regexp"
)

// Estructura mejorada para los criterios met/no met
type CriteriaMet struct {
	Length      bool `json:"length"`
	Uppercase   bool `json:"uppercase"`
	Lowercase   bool `json:"lowercase"`
	Number      bool `json:"number"`
	SpecialChar bool `json:"specialChar"`
}

// Respuesta para el frontend
type PasswordResponse struct {
	Level    string      `json:"level"`
	Score    int         `json:"score"` // 0-4 para el medidor visual
	Message  string      `json:"message"`
	Criteria CriteriaMet `json:"criteria"`
}

// Función principal de análisis
func AnalyzePasswordStrength(password string) PasswordResponse {
	if len(password) == 0 {
		return PasswordResponse{} // Devolver vacío si no hay entrada
	}

	criteria := CriteriaMet{
		Length:      len(password) >= 8,
		Uppercase:   regexp.MustCompile(`[A-Z]`).MatchString(password),
		Lowercase:   regexp.MustCompile(`[a-z]`).MatchString(password),
		Number:      regexp.MustCompile(`[0-9]`).MatchString(password),
		SpecialChar: regexp.MustCompile(`[!@#\$%\^&\*]`).MatchString(password),
	}

	// Calcular "entropía" numérica para el medidor visual
	score := 0
	if criteria.Length { score++ }
	if criteria.Uppercase { score++ }
	if criteria.Lowercase { score++ }
	if criteria.Number { score++ }
	if criteria.SpecialChar { score++ }

	level := "Baja"
	message := "Débil. Agrega mayúsculas, números o caracteres especiales."
	
	if len(password) < 6 {
		level = "Muy Baja"
		message = "Contraseña demasiado corta. Vulnerable a ataques inmediatos."
		score = 1 // Forzar representación visual baja
	} else if len(password) >= 12 && score >= 4 {
		level = "Alta"
		message = "Excelente. Contraseña robusta y tardaría siglos en ser descifrada."
	} else if len(password) >= 8 && score >= 3 {
		level = "Media"
		message = "Aceptable, pero puede mejorar agregando más variedad."
	}

	return PasswordResponse{
		Level:    level,
		Score:    score,
		Message:  message,
		Criteria: criteria,
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Simple HTML para la interfaz de usuario MODERNA Y CORREGIDA
	if r.Method == "GET" && r.URL.Path == "/" {
		// ¡IMPORTANTE! Forzar Charset UTF-8 para evitar caracteres extraños (Ñ)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8"> <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Panel de Ciberseguridad: Analizador de Contraseñas</title>
    <style>
        :root {
            --bg-color: #121212; --card-color: #1e1e1e; --text-primary: #ffffff; --text-secondary: #aaaaaa;
            --low: #ff4d4d; --med: #ffa64d; --high: #4dff4d; --accent: #2196f3;
        }
        body { font-family: 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; background-color: var(--bg-color); color: var(--text-primary); margin: 0; display: flex; align-items: center; justify-content: center; min-height: 100vh; padding: 20px; box-sizing: border-box; }
        .container { background-color: var(--card-color); padding: 40px; border-radius: 12px; box-shadow: 0 10px 25px rgba(0,0,0,0.5); width: 100%; max-width: 450px; text-align: center; }
        h1 { font-size: 24px; margin-bottom: 30px; font-weight: 600; color: var(--text-primary); }
        .input-group { margin-bottom: 25px; position: relative; }
        input[type="password"] { width: 100%; padding: 12px 15px; font-size: 16px; border-radius: 8px; border: 2px solid #333; background-color: rgba(255,255,255,0.05); color: var(--text-primary); box-sizing: border-box; transition: border-color 0.3s; }
        input[type="password"]:focus { outline: none; border-color: var(--accent); }
        .meter-container { height: 8px; width: 100%; background-color: #333; border-radius: 4px; margin-bottom: 15px; overflow: hidden; position: relative; }
        #meter-bar { height: 100%; width: 0%; background-color: #555; transition: width 0.3s, background-color 0.3s; border-radius: 4px; }
        #result-status { font-weight: bold; font-size: 18px; text-transform: uppercase; letter-spacing: 1px; margin-bottom: 5px; min-height: 22px; }
        #result-message { color: var(--text-secondary); font-size: 14px; margin-bottom: 25px; min-height: 40px; }
        .requirements { text-align: left; background-color: rgba(0,0,0,0.2); padding: 15px; border-radius: 8px; font-size: 14px; color: var(--text-secondary); }
        .requirements ul { list-style: none; padding: 0; margin: 0; }
        .requirements li { margin-bottom: 8px; display: flex; align-items: center; }
        .requirements li::before { content: "❌"; margin-right: 10px; font-size: 12px; opacity: 0.5; }
        .requirements li.met { color: var(--high); }
        .requirements li.met::before { content: "✅"; opacity: 1; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Panel de Ciberseguridad:<br>Analizador de Contraseñas</h1>
        
        <div class="input-group">
            <input type="password" id="pass" placeholder="Ingresa tu contraseña para analizarla..." oninput="analyze()" autocomplete="off">
        </div>

        <div class="result-area">
            <div class="meter-container">
                <div id="meter-bar"></div>
            </div>
            <div id="result-status"></div>
            <p id="result-message"></p>
        </div>

        <div class="requirements">
            <strong>Requisitos mínimos:</strong>
            <ul>
                <li id="req-length">Mínimo 8 caracteres</li>
                <li id="req-upper">Al menos una mayúscula (A-Z)</li>
                <li id="req-lower">Al menos una minúscula (a-z)</li>
                <li id="req-number">Al menos un número (0-9)</li>
                <li id="req-special">Carácter especial (!@#$...)</li>
            </ul>
        </div>
    </div>

    <script>
        function analyze() {
            let pass = document.getElementById('pass').value;
            const bar = document.getElementById('meter-bar');
            const statusEl = document.getElementById('result-status');
            const messageEl = document.getElementById('result-message');
            const reqs = { length: 'req-length', upper: 'req-upper', lower: 'req-lower', number: 'req-number', special: 'req-special' };

            // Resetear si está vacío
            if (pass.length === 0) {
                bar.style.width = '0%'; bar.style.backgroundColor = '#555';
                statusEl.innerText = ''; messageEl.innerText = '';
                Object.values(reqs).forEach(id => document.getElementById(id).classList.remove('met'));
                return;
            }

            // Llamada AJAX mejorada
            fetch('/analyze?password=' + encodeURIComponent(pass))
                .then(res => res.json())
                .then(data => {
                    // 1. Actualizar Medidor Visual
                    let scorePercentage = (data.score / 5) * 100;
                    bar.style.width = scorePercentage + '%';
                    
                    // Colores del medidor
                    if (data.level === 'Alta') { bar.style.backgroundColor = 'var(--high)'; statusEl.style.color = 'var(--high)'; }
                    else if (data.level === 'Media') { bar.style.backgroundColor = 'var(--med)'; statusEl.style.color = 'var(--med)'; }
                    else { bar.style.backgroundColor = 'var(--low)'; statusEl.style.color = 'var(--low)'; }

                    // 2. Actualizar Textos
                    statusEl.innerText = 'Fuerza: ' + data.level;
                    messageEl.innerText = data.message;

                    // 3. Actualizar Lista de Requisitos
                    updateReq('req-length', data.criteria.length);
                    updateReq('req-upper', data.criteria.uppercase);
                    updateReq('req-lower', data.criteria.lowercase);
                    updateReq('req-number', data.criteria.number);
                    updateReq('req-special', data.criteria.specialChar);
                });
        }

        function updateReq(id, isMet) {
            const el = document.getElementById(id);
            if (isMet) el.classList.add('met');
            else el.classList.remove('met');
        }
    </script>
</body>
</html>
		`))
		return
	}

	// API Endpoint para analizar en tiempo real (JSON)
	if r.Method == "GET" && r.URL.Path == "/analyze" {
		password := r.URL.Query().Get("password")
		response := AnalyzePasswordStrength(password)
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/analyze", handler)
	println("Servidor moderno corriendo en http://localhost:8080 ...")
	http.ListenAndServe(":8080", nil)
}
