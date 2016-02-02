[![Build Status](https://ci.matthewbrown.io/api/badges/mnbbrown/postmark/status.svg)](https://ci.matthewbrown.io/mnbbrown/postmark)

*postmark* is a golang client for [Postmark](https://postmarkapp.com).

It support the following API endpoints:

 - Email API - for sending email.
 - Server API - for getting server status - also useful for testing client configuration.

It also has a simple worker queue implementation for async sending emails.
