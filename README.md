# Insider Go Backend API

Modern, scalable Go backend API with comprehensive monitoring stack.

## 🚀 Features

- RESTful API with JWT authentication
- PostgreSQL database with migrations
- Credit/Debit transaction system
- Rate limiting and middleware
- Comprehensive monitoring with Prometheus & Grafana
- Container monitoring with cAdvisor
- System metrics with Node Exporter

## 📊 Monitoring Stack

### Available Dashboards

- **Go Application Metrics** - HTTP performance, response times, status codes
  - 🔗 [Live Dashboard](https://snapshots.raintank.io/dashboard/snapshot/TLMHjQAecHpSSpOJAWvxA3apgfNp9LCZ)
- **cAdvisor Container Metrics** - Docker container monitoring
  - 🔗 [Live Dashboard](https://snapshots.raintank.io/dashboard/snapshot/TLMHjQAecHpSSpOJAWvxA3apgfNp9LCZ)
- **Node Exporter System Metrics** - CPU, memory, disk, network
  - 🔗 [Live Dashboard](https://snapshots.raintank.io/dashboard/snapshot/Dy01bPhC3luXWP5fY9YZDPJhb2c6uKuz)
- **PostgreSQL Database Metrics** - Database performance and health
  - 🔗 [Live Dashboard](https://snapshots.raintank.io/dashboard/snapshot/run9uBWH7oeyCImgzx90I6gK0T28aCRU?orgId=0&refresh=10s)

### Quick Access

- Grafana: http://localhost:3000 (admin/admin)
- Prometheus: http://localhost:9090
- Application: http://localhost:8080
- cAdvisor: http://localhost:8081
- Node Exporter: http://localhost:9100/metrics
- PostgreSQL Exporter: http://localhost:9187/metrics

## 🛠 Quick Start

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop services
docker-compose down
```
