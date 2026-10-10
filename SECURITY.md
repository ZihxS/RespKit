# Security

Report suspected security vulnerabilities through [GitHub's private reporting form](https://github.com/ZihxS/RespKit/security/advisories/new) or email **m.saleh.solahudin@gmail.com**. Include the affected version and steps to reproduce when possible. Do not file public issues containing credentials, sensitive user data, SQL queries, production logs, or exploitable details.

Dependabot alerts are enabled for the repository. Maintainers should review alerts in each module, including dependencies used only by integration tests.

Never expose raw internal errors in production HTTP responses. `resp.WithExposeInternalErrors(true)` intentionally permits disclosure only for private/local development, and must remain disabled in public environments.

Diagnostic data may contain application-specific values. Review it before sharing logs or support bundles.
