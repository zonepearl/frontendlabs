# 90-Day Cloud/Hybrid DevOps Engineering Plan

> 📖 **Read this on [frontendlabs.xyz](https://frontendlabs.xyz/kubernetes/plan/)**: the official edition, with one page per chapter, search, and dark mode. <!-- frontendlabs-notice -->

> **Companion guide:** [Docker & Kubernetes — The Real-Life Field Guide](./real-life-k8s-guide.md) explains the *how* and *why* behind every week of this plan.

**Target:** Sr. Engineering Manager transitioning to hands-on Cloud/Hybrid DevOps
**Starting skill level:** Novice
**Goal:** Intermediate proficiency across Linux, containers, Kubernetes, CI/CD, observability, and multi-cloud operations

---

## Guiding Principles

- Learn by doing: every concept gets a hands-on lab
- Build one cumulative project (a sample microservice app) that you deploy, monitor, break, and recover throughout the 90 days
- Weekly retros: 30 minutes every Sunday to review gaps, update notes
- Use a local kind/minikube cluster + one cloud account (start with AWS or GCP free tier)

---

## Phase 1 — Foundation (Days 1–30)

### Week 1-2: Linux & Shell Scripting

**Goal:** Be comfortable navigating, automating, and troubleshooting Linux systems from the command line.

**Topics:**
- Filesystem hierarchy: `/etc`, `/var`, `/proc`, `/sys`, `/tmp`
- File permissions, ownership, ACLs (`chmod`, `chown`, `setfacl`)
- Process management: `ps`, `top`, `htop`, `kill`, `systemctl`, `journalctl`
- Networking basics: `ip`, `ss`, `netstat`, `curl`, `wget`, `dig`, `nslookup`, `traceroute`
- Package management: `apt`/`yum`/`dnf`
- Users, groups, sudo, SSH key management
- Cron jobs and scheduling
- Bash scripting: variables, loops, conditionals, functions, error handling (`set -euo pipefail`)
- Text processing: `grep`, `awk`, `sed`, `cut`, `sort`, `uniq`, `jq`
- Log files: `tail -f`, `grep` patterns, `/var/log/syslog`, `journalctl -u`

**Hands-on Labs:**
- [ ] Write a Bash script that monitors disk usage and emails/Slacks an alert above 80%
- [ ] Write a log-rotation script using `logrotate` config
- [ ] Automate SSH key distribution to 3 VMs using a Bash loop
- [ ] Parse a sample application log file with `awk`/`jq` to extract error rates

**Resources:**
- `man` pages and `tldr` CLI tool (install it)
- Linux Command Line by William Shotts (free online)
- `overthewire.org/wargames/bandit` — gamified Linux practice

---

### Week 3: Docker Hands-On

**Goal:** Build, run, debug, and publish container images confidently.

**Topics:**
- Container vs VM mental model
- Docker architecture: daemon, client, image layers, registry
- Dockerfile best practices: layer caching, multi-stage builds, `.dockerignore`, non-root users
- `docker build`, `docker run`, `docker exec`, `docker logs`, `docker inspect`
- Volumes and bind mounts
- Networking: bridge, host, overlay, port mapping
- Docker Compose: multi-service local stacks
- Image scanning: `docker scout` / `trivy`
- Registry: Docker Hub, ECR, GCR, Artifact Registry

**Hands-on Labs:**
- [ ] Containerize a Node.js or Python Flask app with a multi-stage Dockerfile
- [ ] Build a Docker Compose stack: app + Postgres + Redis
- [ ] Scan the image with `trivy`; fix at least one HIGH vulnerability
- [ ] Push to a private ECR or GCR registry
- [ ] Use `docker stats` and `docker top` to observe resource usage live

**Resources:**
- Docker official docs (docs.docker.com)
- *Docker Deep Dive* by Nigel Poulton (short, practical)

---

### Week 4: Kubernetes Fundamentals

**Goal:** Deploy and manage workloads on a local Kubernetes cluster.

**Topics:**
- K8s architecture: control plane (API server, etcd, scheduler, controller-manager), nodes, kubelet, kube-proxy
- Core objects: Pod, ReplicaSet, Deployment, Service (ClusterIP/NodePort/LoadBalancer), ConfigMap, Secret, Namespace
- `kubectl` essentials: `get`, `describe`, `apply`, `delete`, `exec`, `logs`, `port-forward`, `rollout`
- YAML manifest authoring
- Labels, selectors, annotations
- Resource requests and limits
- Liveness and readiness probes
- `kubectl explain` — your best friend

**Hands-on Labs:**
- [ ] Stand up a local cluster with `kind` or `minikube`
- [ ] Deploy the sample app (from Week 3) as a Deployment + Service
- [ ] Scale replicas up/down; observe pod scheduling
- [ ] Trigger a rolling update; observe rollout status
- [ ] Deliberately break a probe; watch pod restart behavior

**Resources:**
- *Kubernetes in Action* by Marko Luksa (the K8s bible for practitioners)
- killer.sh / killercoda.com — browser-based K8s labs

---

## Phase 2 — Intermediate Operations (Days 31–60)

### Week 5: Kubernetes — Advanced Workloads & Networking

**Topics:**
- StatefulSets, DaemonSets, Jobs, CronJobs
- Ingress controllers (NGINX Ingress, Traefik)
- Network policies
- Persistent Volumes, PVCs, StorageClasses
- Helm: install, upgrade, rollback, write your own chart
- Kustomize: base + overlays pattern

**Hands-on Labs:**
- [ ] Deploy a StatefulSet for Postgres with PVC
- [ ] Install NGINX Ingress; route traffic to two services by path
- [ ] Package the sample app as a Helm chart with `values.yaml` overrides
- [ ] Apply a NetworkPolicy that restricts pod-to-pod traffic

---

### Week 6: Observability — Logs & Metrics

**Goal:** Follow what's happening in a live system; detect and diagnose issues.

**Topics:**
- The three pillars: Logs, Metrics, Traces
- Prometheus: scrape configs, PromQL basics, alerting rules
- Grafana: dashboards, data sources, alert notifications
- Node Exporter, kube-state-metrics, cAdvisor
- Log aggregation: Loki + Promtail stack (or EFK: Elasticsearch + Fluentd + Kibana)
- `kubectl logs --previous`, `stern` for multi-pod tailing
- Distributed tracing concept: OpenTelemetry, Jaeger (overview)
- SLIs, SLOs, error budgets — the language of reliability

**Hands-on Labs:**
- [ ] Deploy kube-prometheus-stack via Helm; import K8s dashboard in Grafana
- [ ] Write a PromQL query: CPU usage per pod, memory limits hit
- [ ] Set up an alert rule: fire when pod restarts > 3 in 10 minutes
- [ ] Deploy Loki + Promtail; search app logs from Grafana
- [ ] Instrument the sample app with a `/metrics` endpoint (use `prom-client` or `prometheus_client`)

**Resources:**
- *Observability Engineering* by Charity Majors, Liz Fong-Jones, George Miranda
- Prometheus docs: prometheus.io/docs

---

### Week 7: CI/CD Pipelines

**Goal:** Automate build → test → containerize → deploy with no manual steps.

**Topics:**
- CI/CD concepts: pipeline stages, artifact promotion, environment gates
- GitHub Actions: workflows, jobs, steps, secrets, reusable workflows
- GitLab CI (optional if org uses it): `.gitlab-ci.yml`
- Pipeline patterns: build on PR, deploy to staging on merge to main, deploy to prod on tag
- Image tagging strategy: `sha`, `semver`, `latest` (avoid `latest` in prod)
- GitOps with ArgoCD: App of Apps pattern, sync policies, auto-rollback on drift
- Secret management in pipelines: GitHub Secrets, HashiCorp Vault, AWS Secrets Manager

**Hands-on Labs:**
- [ ] Build a GitHub Actions pipeline: lint → test → docker build/push → deploy to kind
- [ ] Add image vulnerability scanning step (trivy); fail on CRITICAL
- [ ] Set up ArgoCD; connect it to a GitOps repo; deploy sample app via ArgoCD
- [ ] Implement branch-based environment promotion: `dev` → `staging` → `prod` namespaces
- [ ] Store a DB password in Vault or AWS Secrets Manager; inject it via External Secrets Operator

---

### Week 8: Rollback Strategies

**Goal:** Know exactly how to safely and quickly undo bad deployments.

**Topics:**
- Kubernetes rollout rollback: `kubectl rollout undo`
- Deployment history and revision tracking
- Blue/Green deployments
- Canary deployments with Argo Rollouts or Flagger
- Feature flags as a rollback mechanism (LaunchDarkly concept)
- Helm rollback: `helm rollback <release> <revision>`
- ArgoCD sync rollback: reverting Git commit + ArgoCD sync
- Database migration rollback strategies (Flyway, Liquibase)
- Incident runbooks: define your rollback decision tree

**Hands-on Labs:**
- [ ] Deploy a broken image version; execute `kubectl rollout undo`; verify recovery
- [ ] Install Argo Rollouts; implement a canary: 10% → 50% → 100% traffic shift
- [ ] Simulate a bad Helm release; run `helm rollback`
- [ ] Write a 1-page rollback runbook for the sample app

---

## Phase 3 — Production Readiness (Days 61–90)

### Week 9: Backup & Stateful Data Protection

**Goal:** Know what to back up, how, and how to verify restores.

**Topics:**
- What needs backing up: etcd (cluster state), PVCs (application data), ConfigMaps/Secrets
- etcd backup and restore: `etcdctl snapshot save/restore`
- Velero: cluster backup/restore, scheduled backups, namespace migration
- Volume snapshots: CSI VolumeSnapshot API
- Database-specific backup: `pg_dump`, MySQL dump, mongodump + automation scripts
- Backup retention policies, encryption at rest, off-site storage (S3, GCS)
- Backup verification: automated restore tests
- RPO (Recovery Point Objective) and RTO (Recovery Time Objective) definitions

**Hands-on Labs:**
- [ ] Take an etcd snapshot on a kind cluster; simulate cluster loss; restore from snapshot
- [ ] Install Velero with MinIO (local S3); schedule hourly namespace backups
- [ ] Backup Postgres PVC with Velero; delete namespace; restore and verify data
- [ ] Write a cron-based `pg_dump` script that uploads to S3 with a 7-day retention policy

---

### Week 10: Disaster Recovery & High Availability

**Goal:** Design and validate systems that survive zone failures, node failures, and full-region outages.

**Topics:**
- HA cluster topology: multi-master control plane, etcd quorum (3 or 5 nodes)
- Pod Disruption Budgets (PDB)
- Node affinity, pod anti-affinity, topology spread constraints
- Cluster Autoscaler and Karpenter (AWS)
- Multi-region/multi-cluster strategies: active-active vs active-passive
- Global load balancing: AWS Route 53 / GCP Cloud DNS with health checks
- Cross-region data replication: RDS Multi-AZ, Cloud SQL HA, Redis Sentinel
- Disaster recovery runbook: RTO/RPO targets, communication plan, blameless postmortem template
- Chaos engineering introduction: LitmusChaos, Chaos Monkey

**Hands-on Labs:**
- [ ] Configure PDB for the sample app (maxUnavailable: 1)
- [ ] Set pod anti-affinity so replicas spread across nodes
- [ ] Simulate a node failure (`kubectl delete node`); verify zero-downtime
- [ ] Run a LitmusChaos pod-delete experiment; observe Grafana during chaos
- [ ] Write a DR runbook: "Region us-east-1 is unavailable — steps to fail over to us-west-2"

---

### Week 11: Cloud Platforms — AWS, GCP, Azure

**Goal:** Navigate and operate managed services across the three major clouds.

**Topics (parallel track — pick one cloud deep, survey the other two):**

**AWS (primary):**
- EKS: cluster creation, node groups, Fargate profiles, IAM Roles for Service Accounts (IRSA)
- ECR, ECS (awareness), Lambda (awareness)
- VPC, subnets, security groups, NACLs, VPC peering
- IAM: roles, policies, least-privilege, instance profiles
- RDS, ElastiCache, S3, CloudFront
- CloudWatch: Logs, Metrics, Alarms, Container Insights
- Load balancers: ALB, NLB, AWS Load Balancer Controller
- Terraform / AWS CDK for infra-as-code

**GCP (secondary):**
- GKE Autopilot vs Standard; Workload Identity
- Cloud Run, Cloud Functions (awareness)
- VPC, Cloud NAT, Shared VPC
- Cloud Monitoring, Cloud Logging, Error Reporting
- Artifact Registry, Cloud Build

**Azure (tertiary):**
- AKS, ACR, Azure DevOps Pipelines
- Azure Monitor, Log Analytics Workspace
- Managed identities, RBAC

**Hands-on Labs:**
- [ ] Stand up an EKS cluster with Terraform (or eksctl); deploy sample app
- [ ] Configure IRSA so pods can read from S3 without storing credentials
- [ ] Set up CloudWatch Container Insights on EKS
- [ ] Deploy the same app to GKE; compare CLI/console experience with AWS
- [ ] Use Terraform to provision VPC + EKS + RDS as code; `terraform plan/apply/destroy`

**Resources:**
- *Terraform: Up & Running* by Yevgeniy Brikman
- AWS Well-Architected Framework (free whitepaper)
- GCP Architecture Center (cloud.google.com/architecture)

---

### Week 12: Integration Week — Full-Stack Ops System

**Goal:** Wire everything learned into one coherent, production-grade ops system.

**Deliverables:**
- [ ] Sample microservice app running on EKS (or GKE) via ArgoCD GitOps
- [ ] Full Prometheus + Grafana + Loki observability stack
- [ ] GitHub Actions CI pipeline: lint → test → build → push → trigger ArgoCD sync
- [ ] Canary rollout via Argo Rollouts with automatic rollback on error rate spike
- [ ] Velero scheduled backups to S3 with verified restore
- [ ] DR runbook documented and tested (simulate region failure)
- [ ] All infra defined in Terraform; `terraform apply` from scratch in < 30 min
- [ ] On-call runbook: links to dashboards, rollback steps, escalation path

---

## Recommended Books

### Core Reading (read in order)

| # | Book | Why |
|---|------|-----|
| 1 | *The Linux Command Line* — William Shotts | Free online; best foundation for shell work |
| 2 | *Docker Deep Dive* — Nigel Poulton | Short (200 pages), practical, up-to-date |
| 3 | *Kubernetes in Action* (2nd Ed) — Marko Luksa | The most thorough K8s book; read alongside Phase 1-2 |
| 4 | *The DevOps Handbook* — Kim, Humble, Debois, Willis | Mental model for DevOps culture and flow |
| 5 | *Site Reliability Engineering* — Google (free online) | Core SRE thinking: SLOs, error budgets, toil reduction |
| 6 | *Terraform: Up & Running* (3rd Ed) — Yevgeniy Brikman | Infra-as-code; AWS/GCP/Azure examples |
| 7 | *Observability Engineering* — Majors, Fong-Jones, Miranda | Modern observability beyond metrics |

### Supplemental

| Book | Why |
|------|-----|
| *Continuous Delivery* — Jez Humble, David Farley | CI/CD pipeline design principles |
| *Designing Distributed Systems* — Brendan Burns (free PDF from Microsoft) | Patterns for resilient, scalable systems |
| *Cloud Native DevOps with Kubernetes* — Arundel, Domingus | Practical K8s operations in production |
| *The Phoenix Project* — Kim, Behr, Spafford | Fiction; builds intuition for ops culture fast |
| *Learning eBPF* — Liz Rice | Advanced Linux observability (save for post-90 days) |

---

## Weekly Rhythm

| Day | Activity |
|-----|----------|
| Mon–Thu | Study + hands-on lab (2–3 hours/day) |
| Fri | Build on the cumulative project; apply the week's skills end-to-end |
| Sat | Read (book or documentation) |
| Sun | 30-min retro: what's gaps, update notes, plan next week |

---

## Tools to Install on Day 1

```bash
# CLI tools
brew install kubectl helm kind terraform kustomize
brew install jq yq fzf bat ripgrep tldr stern
brew install awscli
curl -sSL https://install.python-poetry.org | python3 -  # if Python work
brew install trivy          # image scanning
brew install argocd         # ArgoCD CLI
brew install velero         # backup CLI

# VSCode extensions
# - Kubernetes (ms-kubernetes-tools)
# - Docker
# - HashiCorp Terraform
# - YAML
# - GitLens
```

---

## Success Metrics at Day 90

- [ ] Can deploy a multi-service app to Kubernetes from scratch, using only a Git repo and Terraform
- [ ] Can diagnose a pod crash using logs, metrics, and `kubectl describe` in under 5 minutes
- [ ] Can execute a safe rollback (Helm, K8s, ArgoCD) without service downtime
- [ ] Can restore a namespace from Velero backup within 15 minutes
- [ ] Has a working DR runbook tested against a simulated failure
- [ ] Understands cloud bill basics: what EKS, RDS, NAT Gateway, and data transfer cost
- [ ] Can have a credible technical conversation with platform/SRE engineers about architectural tradeoffs

---

*Last updated: 2026-05-19*
