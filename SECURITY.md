# Security

Report suspected security vulnerabilities privately to **m.saleh.solahudin@gmail.com**. Please do not file public issues containing credentials, sensitive user data, SQL queries, production logs, or exploitable details. Include the affected version and steps to reproduce when possible.

Never expose raw internal errors in production HTTP responses. `resp.WithExposeInternalErrors(true)` intentionally permits disclosure only for private/local development, and must remain disabled in public environments.

Diagnostic data may contain application-specific values. Review it before sharing logs or support bundles.
