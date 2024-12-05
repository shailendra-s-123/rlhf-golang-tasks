func (j *JWTProcessor) VerifyToken(tok *jwt.Token) (*models.UserToken, error) {
	// claim contains the IP address of the caller to enforce Security protection by IP address restrictions
	claims, ok := tok.Claims.(*jwt.CustomClaims)

	if ok && tok.Valid {
		log.Println("HELLOOOO", claims.Username)
		responseObject := models.UserToken{Username: claims.Username}
		return &responseObject, nil
	} else {
		return nil, errors.New("Invalid token")
	}
}

func handleValidateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	token := r.FormValue("token")

	if token == "" {
		http.Error(w, "Token parameter missing", http.StatusBadRequest)
		return
	}

	jp := jwtUtil.JWTProcessor{}
	userToken, err := jp.ValidateToken(token)

	if err != nil {
		log.Printf("Token validation error: %v", err)
		if utils.IsTokenExpiredError(err) {
			http.Error(w, "Token expired", http.StatusUnauthorized)
		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userToken); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}