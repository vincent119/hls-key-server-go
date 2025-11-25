# Security Policy

## Supported Versions

We release patches for security vulnerabilities. Currently supported versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of HLS Key Server seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### Please do NOT:

- Open a public GitHub issue for security vulnerabilities
- Disclose the vulnerability publicly before it has been addressed

### Please DO:

1. **Email us directly** at the maintainer's email or create a private security advisory on GitHub
2. **Provide detailed information** including:
   - Type of vulnerability (e.g., authentication bypass, path traversal, etc.)
   - Steps to reproduce the issue
   - Potential impact of the vulnerability
   - Any suggested fixes (optional)

### What to expect:

- **Initial Response**: Within 48 hours of your report
- **Status Update**: We will keep you informed about the progress
- **Fix Timeline**: We aim to release a fix within 7-14 days for critical issues
- **Credit**: We will credit you in the release notes (unless you prefer to remain anonymous)

## Security Best Practices

When deploying HLS Key Server, please follow these security recommendations:

### 1. JWT Configuration

```yaml
jwt:
  secret: "use-a-strong-random-secret-at-least-32-characters-long"
  expiration_hours: 168  # Adjust based on your needs
  header_key: "use-a-custom-header-name"
  header_value: "use-a-strong-random-value"
```

- Use a cryptographically strong random secret (minimum 32 characters)
- Change default header key/value in production
- Rotate JWT secrets regularly
- Set appropriate token expiration times

### 2. Key Storage

```bash
# Set proper file permissions
chmod 700 keys/
chmod 600 keys/*.key

# Keys directory should only be readable by the server process
chown server-user:server-group keys/
```

### 3. Network Security

- Deploy behind a reverse proxy (nginx, Caddy)
- Use HTTPS/TLS in production
- Implement rate limiting at the proxy level
- Configure CORS properly for your domains

### 4. Environment Configuration

```yaml
app:
  mode: "production"  # Never use debug mode in production
  port: "9090"
```

- Always use `production` or `release` mode in production
- Disable debug logging in production
- Use environment variables for sensitive configuration
- Never commit secrets to version control

### 5. Access Control

- Implement IP whitelisting if possible
- Use strong authentication credentials
- Regularly review and rotate access tokens
- Monitor access logs for suspicious activity

### 6. Docker Security

```dockerfile
# Run as non-root user
USER nonroot:nonroot

# Read-only file system
RUN chmod -R 400 /app/keys
```

### 7. Monitoring & Logging

- Enable structured logging
- Monitor for unusual access patterns
- Set up alerts for authentication failures
- Regularly review access logs

## Known Security Considerations

### Path Traversal Protection

The application includes built-in protection against path traversal attacks:
- Key names are validated to prevent `../` patterns
- File access is restricted to the configured keys directory

### JWT Security

- Uses HMAC-SHA256 for token signing
- Tokens include expiration times (configurable)
- Supports multiple validation methods (Bearer token, query param, cookie)

### Rate Limiting

**Recommendation**: Implement rate limiting at the reverse proxy level:

```nginx
# Example nginx configuration
limit_req_zone $binary_remote_addr zone=auth:10m rate=5r/s;
limit_req_zone $binary_remote_addr zone=keys:10m rate=100r/s;

location /api/v1/auth/ {
    limit_req zone=auth burst=10;
}

location /api/v1/hls/ {
    limit_req zone=keys burst=200;
}
```

## Security Audit History

| Date       | Auditor | Findings | Status |
|------------|---------|----------|--------|
| 2025-11-25 | Internal | Initial release | ✅ |

## Compliance

This project follows:
- OWASP API Security Top 10
- CWE Top 25 Most Dangerous Software Weaknesses
- Go Secure Coding Practices

## Contact

For security concerns, please contact:
- GitHub Security Advisories: [Create Advisory](https://github.com/vincent119/hls-key-server-go/security/advisories/new)
- Email: (Maintainer's email - update as needed)

## Acknowledgments

We appreciate the security research community's efforts to responsibly disclose vulnerabilities. Contributors who report valid security issues will be acknowledged in our security hall of fame (with permission).

---

**Last Updated**: 2025-11-25
