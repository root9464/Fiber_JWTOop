package root

/* func identifyExpiredTokens(db *gorm.DB) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			expiredTokenIDs, err := findExpiredTokens(db)
			if err != nil {
				log.Println("Error identifying expired tokens:", err)
				// Handle the error
			} else {
				// Process the expired tokens, e.g., invalidate them in the database
				for _, tokenID := range expiredTokenIDs {
					// Invalidate the token
				}
			}
		}
	}
}

func findExpiredTokens(db *gorm.DB) ([]string, error) {
	var expiredTokens []auth.Token

	currentTime := time.Now()
	err := db.Where("expiry <= ?", currentTime).Find(&expiredTokens).Error
	if err != nil {
		return nil, err
	}

	var expiredTokenIDs []string
	for _, token := range expiredTokens {
		expiredTokenIDs = append(expiredTokenIDs, token.JwtAccessToken, token.RefreshToken)
	}

	return expiredTokenIDs, nil
} */
