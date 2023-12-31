<p align="center" style="font-size:70px;" >
    <a href="https://www.invisiblelab.dev/">
    CertWatch  
    </a>
    <a href="https://www.invisiblelab.dev/">
    <img src="https://uploads-ssl.webflow.com/60057003af6cfb7362bab247/6005a8ba64602c1ef34c244f_brand.svg" width="100px" alt="InvisibleLab logo" />
    </a>
</p>

<h3 align="center">Web page SSL Certificate watcher package and cli tool.</h3>
<p align="center"> Get notified by email and/or slack when your certificates are close to expire. </p>
<br />

CertWatch is a open-source project. Just add some domains to the config file and the certification retrieval is ready to watch. After some email and slack configuration the notification system is ready to go!

This repo contains 2 tools:

1. The certwatch package, which you can import to your project to retrieve certificate information, calculate days until certificate expiration and send emails via smtp.

2. The certwatch cli, which uses the package to add domains to your config file, command to search a specific domains certificate or all from your list of domains in the config file.

## Getting Started

### Installation

Clone the repo and complete the config file.

### Configuration file

`certwatch.example.yaml` has all the configuration needed:

-   domains:
    you can manually add them or via the cli.
    -   name: is the domain url. It should include the subdomain, domain and top level domain without the protocol (e.g. [www.invisiblelab.dev](https://www.invisiblelab.dev/)).
    -   threshold: days until the certificate expires that you want to receive notification via slack and/or email.
-   refresh: days since last query you want to re-request the certificates. See the [Checking Certificates](https://github.com/invisiblelab-dev/certwatch#checking-certificates) section for deeper explanation.
-   notifications:
    -   email:
        -   username: smtp username given by email provider.
        -   password: smtp password given by email provider.
        -   smtphost: smtp host (e.g. for mailtrap test email is **sandbox.smtp.mailtrap.io**)
        -   port: smtphost port (e.g. smtphost +":"+ port = sandbox.smtp.mailtrap.io:45)
        -   from: email address that will be the from address of the email notifications.
        -   to: email address that will be the to address of the email notifications.
    -   slack:
        -   webhook: webhook url (e.g. https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX). For more information check [slack documentation](https://api.slack.com/messaging/webhooks#create_a_webhook).

## Adding a domain

To add a new domain, add it in your config file, should follow `certwatch.example.yaml` structure. Its location is by default be in the root of this repo but, if wanted, can be stored in any location but has to add the path location by the --config flag to the command in the terminal when running check-all.

```
go run ./cmd/certwatch check-all --config "/Users/user/Documents/certwatch.yaml"
```

## Running the cli

To use see the cli available commands run the following:

```bash
go run ./cmd/certwatch
```

To check hardcoded domains certificates:

```bash
go run ./cmd/certwatch check --domain example.com --domain example2.com
```

To check all the domains in your `certwatch.yaml` file, send notifications if needed and create a `queries.yaml` file with the desired info:

```bash
go run ./cmd/certwatch check-all
```

## Checking Certificates

To check the certificates of each domain, the package is making a GET request to the domains url and returning the certificate information. After that, depending on your configuration and **threshold** (environment variable in `certwatch.yaml`) established, the package is able to send a notification via email and/or slack to notify that the certificate is expired or will expire in $\leq$ **threshold** (days).

Because the certificates are retrieved via GET request and the certificates should not change very often, a cache file is implemented. The package will make a new GET request certificate retrieval only if there is a new domain on the `certwatch.yaml` file or if any of the certificates were last retrieved via GET request $\geq$ **refresh** days ago.
