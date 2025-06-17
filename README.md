# Prometheus Exporter for [ConfigCat](https://configcat.com)

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/dirsigler/configcat-exporter)](https://goreportcard.com/report/github.com/dirsigler/configcat-exporter)

This is an custom Prometheus Exporter for the https://configcat.com featureflag solution.

Allows your business to monitor all feature-flags across your Organization, Product and Configs.
Also includes metrics for Zombie flags (stale feature-flags).
Improve your software application hygiene and better understand the usage of your [configcat.com](https://configcat.com) setup.

## ⚙️ Metrics

The ConfigCat Prometheus Exporter supports all basic pre-configured types of incidents available in [configcat.com](https://configcat.com).

| Name                                | Label                                                    | Description                                        |
| ----------------------------------- | -------------------------------------------------------- | -------------------------------------------------- |
| `configcat_products_total`          |                                                          | Total number of ConfigCat products                 |
| `configcat_configs_total`           | `product_name`, `product_id`                             | Total number of ConfigCat configs per product      |
| `configcat_environments_total_`     | `product_name`, `product_id`                             | Total number of ConfigCat environments per product |
| `configcat_feature_flags_total`     | `product_name`, `product_id`, `config_name`, `config_id` | Total number of feature flags per config           |
| `configcat_scrape_errors_total`     |                                                          | Total number of ConfigCat API scrape errors        |
| `configcat_last_scrape_timestamp`   |                                                          | Unix timestamp of the last successful scrape       |
| `configcat_scrape_duration_seconds` |                                                          | Duration of ConfigCat API scrapes in seconds       |
| `configcat_zombie_flags_total`      | `product_name`, `product_id`, `config_name`, `config_id` | Total number of zombie flags per product           |

## 🚀 Deployment

> IMPORTANT: You have to provide the "CONFIGCAT_API_KEY="<MY_API_KEY>" environment variable to your deployment for the ConfigCat Prometheus Exporter to work.

---

With each [release](https://github.com/dirsigler/configcat-exporter/releases) I also provide a [secure by default](https://www.chainguard.dev/chainguard-images) Docker Image.

You can chose from:

- The Image on GitHub => [configcat-exporter on GitHub](https://github.com/dirsigler/configcat-exporter/pkgs/container/configcat-exporter)
- The Image on DockerHub => [configcat-exporter on DockerHub](https://hub.docker.com/repository/docker/dirsigler/configcat-exporter/general)

### Docker

```sh
docker run --rm \
--interactive --tty \
--env CONFIGCAT_API_KEY="<MY_API_KEY>" \
ghcr.io/dirsigler/configcat-exporter:latest
```

You can also enable a logger with Debug mode via the `--log.level=debug` flag.
See the available [configuration](#🚩-configuration)

## 🚩 Configuration

```sh
$ configcat-exporter --help

NAME:
   configcat-exporter - Export ConfigCat metrics to Prometheus

USAGE:
   configcat-exporter [global options] command [command options]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --configcat-api-key value  ConfigCat API key (can also be set via CONFIGCAT_API_KEY env var) [$CONFIGCAT_API_KEY]
   --port value               Port to listen on (default: 8080) [$PORT]
   --scrape-interval value    Scrape interval in seconds (default: 60) [$SCRAPE_INTERVAL]
   --log-level value          Log level (debug, info, warn, error) (default: "info") [$LOG_LEVEL]
   --organization-id value    ConfigCat Organization ID [$CONFIGCAT_ORGANIZATION_ID]
   --product-id value         ConfigCat Product ID [$CONFIGCAT_PRODUCT_ID]
   --help, -h                 show help
```

## 📝 License

Built with ☕️ and licensed via [Apache 2.0](./LICENSE)
