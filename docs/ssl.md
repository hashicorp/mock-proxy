## SSL Certificates and Mocking HTTPS requests

Using Squid's SSL Bump configuration, mock-proxy can also act as an
`https_proxy` and successfully mock upstream requests to HTTPS endpoints.

It does so in this local configuration using a self-signed certificate in
`/certs`. This self-signed cert is automatically trusted for local dev, and
should work out of the box.

If configuring mock-proxy in another environment, you will need to volume mount
a self-signed certificate to `/etc/squid/ssl_cert/ca.pem`, and trust that
certificate on any system attempting to use mock-proxy as an `https_proxy`.

To generate these certificates, use the script: `/hack/gen-certs.sh`. This may
also be useful example code if you need to incorporate self-signed certs into
another system using mock-proxy.
