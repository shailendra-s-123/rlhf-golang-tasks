func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.Header.Get("X-Username")
	password := r.Header.Get("X-Password")
	clientName := r.Header.Get("x-client-name")
	clientVersion := r.Header.Get("X-Client-Version")

	// Authenticate the user
	// check/concrete impl of authentication using external policies, LDAP
	// or local storage like in this example.
	authSuccess := auth.AuthenticateUser(username, password)

	if !authSuccess {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Authentication successful, generate a token
	token, err := generateToken(username)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Add additional permissions or metadata
	ctx := context.WithValue(r.Context(), models.HeaderClient, clientName)
	ctx = context.WithValue(ctx, models.HeaderClientVersion, clientVersion)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		log.Printf("Error writing response: %v", err)
	}
	nextHandler(ctx, w, r)
}