# Security Considerations

This document outlines security measures implemented in this project and provides recommendations for secure deployment and maintenance.

## Implemented Measures

- **JWT Security**:
    - The application enforces a minimum length of 32 characters for the `JWT_SECRET` environment variable. The application will not start if the secret is missing or too weak.
    - JWTs include standard claims like `exp` (expiration) and `iat` (issued at).
- **Admin Product Updates**:
    - The endpoint for updating product details (`PATCH /admin/products/:id`) now uses a specific request structure, preventing unintended or malicious updates to sensitive product fields.
- **Rate Limiting**:
    - Authentication endpoints (`/auth/login`, `/auth/register`) have global rate limiting (10 requests per minute, burst of 5) to protect against brute-force attacks.
- **File Upload Security**:
    - The admin product image upload endpoint (`POST /admin/products/:id/image`) validates uploaded files for:
        - Maximum file size (5MB).
        - Allowed MIME types (`image/jpeg`, `image/png`, `image/gif`).
    - Image uploads are handled via Cloudinary, offloading storage and some processing.
- **Security Headers**:
    - The application sets the following HTTP security headers for all responses to help mitigate common web vulnerabilities:
        - `X-Content-Type-Options: nosniff`
        - `X-Frame-Options: DENY`
        - `Content-Security-Policy: default-src 'self'; script-src 'self'; object-src 'none';`

## Recommendations for Secure Deployment & Maintenance

- **JWT Secret (`JWT_SECRET`)**:
    - **Strength**: Generate a cryptographically strong random string of at least 32 characters (64 characters is even better).
    - **Uniqueness**: Use a unique secret for each environment (development, staging, production).
    - **Rotation**: Implement a policy for regular rotation of the JWT secret.
    - **Storage**: Store the secret securely (e.g., using environment variables injected by your deployment platform, or a secret management system). Do not commit it to your repository.
- **CORS Origins (`CORS_ORIGINS`)**:
    - Configure `CORS_ORIGINS` precisely to only allow requests from your trusted frontend domains. Avoid using wildcard `*` in production if possible.
- **Dependency Vulnerability Scanning**:
    - Regularly scan your application's dependencies for known vulnerabilities. Use tools like `govulncheck` (`go install golang.org/x/vuln/cmd/govulncheck@latest && govulncheck ./...`).
    - Integrate dependency scanning into your CI/CD pipeline.
- **Keep Software Updated**:
    - Keep your Go version updated to the latest stable release.
    - Regularly update all project dependencies (`go get -u all`).
- **Input Validation**:
    - While the application uses Gin's binding for basic validation, ensure all user-supplied input is thoroughly validated and sanitized, especially if new input fields or complex data structures are introduced.
- **Error Handling**:
    - Ensure that error messages returned to clients do not leak sensitive internal information or stack traces.
- **Logging and Monitoring**:
    - Implement comprehensive logging to track important events, especially security-related ones (e.g., failed login attempts, unauthorized access attempts).
    - Monitor application logs and server activity for suspicious patterns.
- **HTTPS**:
    - Always use HTTPS in production to encrypt traffic between clients and your server.
- **Web Application Firewall (WAF)**:
    - Consider using a WAF to protect against common web attacks (SQLi, XSS, etc.) at the edge.
- **Principle of Least Privilege**:
    - Ensure that database users and any other integrated services have only the minimum necessary permissions.
- **Cloudinary Security**:
    - Review Cloudinary account settings for security best practices regarding API key management, allowed transformations, and access controls if you manage your Cloudinary account directly.

By following these recommendations and regularly reviewing your security posture, you can help maintain the security and integrity of your application.
