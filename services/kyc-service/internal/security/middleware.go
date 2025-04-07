package security

func SecurityMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Rate limiting with dynamic thresholds
        if !RateLimiter.Allow(GetClientID(c)) {
            c.AbortWithStatus(429)
            return
        }

        // AI-powered threat detection
        if score := ThreatDetector.Analyze(c.Request); score > ThresholdScore {
            AuditLogger.LogThreat(c)
            c.AbortWithStatus(403)
            return
        }

        // Document encryption
        if doc := c.GetDocument(); doc != nil {
            doc.Encrypt(EncryptionKey)
        }

        c.Next()
    }
}