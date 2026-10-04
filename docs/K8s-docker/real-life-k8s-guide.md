# Docker & Kubernetes — The Real-Life Field Guide (Beginner → Production)

> 📖 **Read this on [frontendlabs.xyz](https://frontendlabs.xyz/kubernetes/)**: the official edition, with one page per chapter, search, and dark mode. <!-- frontendlabs-notice -->

> A hands-on, example-driven path from "what is a container" to running a
> production-grade, auto-scaling, self-healing, fully observable cloud system
> on Kubernetes.
>
> Every topic comes with runnable snippets, a **Real-world example** or
> **War story**, and the **best practice** you would be expected to follow on
> a real platform team. One system is built across the entire guide:
> by the last chapter you have it running locally, on Kubernetes, behind
> GitOps, with OpenTelemetry, Prometheus, Loki, Tempo, Grafana, autoscaling,
> canary releases, backups, and chaos tests.

---

> **Where this fits:** it picks up where step 1 of the series,
> [Operating Systems, Linux & Containers](../os-linux/real-life-os-guide.md),
> leaves off. Its Phase 3 builds containers from namespaces and cgroups and
> introduces Kubernetes; this guide takes the same ideas to production.
>
> **Read alongside:** [Go](../Golang/real-life-golang-guide.md) (every service here is Go) ·
> [Security Engineering in Depth](../security/real-life-security-guide-v1.md) (cloud and Kubernetes security, beyond Part VI) ·
> [Scale, Load & Performance Testing](../scale-perf/real-life-scale-guide.md) (the load-testing theory behind §27.7).
>
> **Day by day:** the [90-day DevOps plan](./k8s-docker-plan.md) schedules this guide's labs week by week.

---

## How to use this guide

- **Part 0 — The lab:** a scripted Ubuntu VM on your Mac (Lima) that
  holds the whole project, plus a milestone runbook and a method for
  verifying this guide as you build.
- **Part I — Containers & Docker (beginner → advanced):** what a container
  really is (namespaces, cgroups, layered filesystems), the Docker toolchain,
  writing production Dockerfiles, networking, storage, Compose, security and
  supply chain, debugging. Ends with the whole system running under
  `docker compose` with local observability.
- **Part II — Kubernetes core (beginner → intermediate):** the control loop,
  a multi-node local cluster, Pods, Deployments, Services, config, storage,
  StatefulSets, Jobs. Ends with the system deployed from raw manifests.
- **Part III — Production workloads (intermediate → advanced):** resources
  and QoS, health checks, graceful shutdown, self-healing, scheduling,
  Gateway API + TLS, Kustomize and Helm, autoscaling (HPA, KEDA, VPA,
  Karpenter).
- **Part IV — Observability:** OpenTelemetry instrumentation, the Collector,
  Prometheus, Loki, Tempo, Grafana, SLOs and burn-rate alerting.
- **Part V — Debugging:** a symptom-driven playbook for every common
  failure, ephemeral debug containers, node-level debugging, profiling.
- **Part VI — Security:** RBAC, Pod Security, NetworkPolicy, secrets
  management, admission policy and image signature verification.
- **Part VII — Delivery:** CI with GitHub Actions, GitOps with Argo CD,
  canary releases with Argo Rollouts, database migrations.
- **Part VIII — Reliability operations:** backup and disaster recovery,
  chaos engineering, cost and capacity, moving to a managed cloud (EKS).
- **Part IX — Capstone & reference:** assemble and verify the full system,
  the production-readiness checklist, cheat sheets, pitfalls, and further
  reading.

Conventions:

- Commands run in a Linux shell: inside the lab VM (§0, Path A) or natively on
  macOS (Path B). Commands marked **[mac]** run on the Mac host. Native macOS
  differences are called out where they matter. `$` is omitted from
  commands so you can copy them directly.
- Kubernetes examples target **Kubernetes 1.33+** (native sidecars,
  `preStop` sleep action, Gateway API). Version-specific behavior is called
  out where it matters.
- Application code is **Go** (the language of Kubernetes itself and the
  most common choice for cloud services); a few snippets show Python for
  comparison.
- **⚠️ Caveat (verify)** boxes mark content written from knowledge of the
  docs but **not yet executed or checked against the latest release**:
  chart value keys, CRD fields that move between versions, version-specific
  feature status, and metric names. Each box says how to check it. Treat
  everything *without* a box as believed correct but still untested until
  you run it in the lab (§0). The "how to verify this guide" checklist is in
  §0.6.
- **Real-world example** boxes ground a concept in something you will meet
  on the job. **War story** boxes describe a failure mode that has taken down
  real production systems. **Best practice** boxes are the rule you should
  default to unless you have a reason not to.
- This guide is the *theory + build* companion to
  [`k8s-docker-plan.md`](./k8s-docker-plan.md) (the 90-day schedule). The
  plan tells you *when*; this guide tells you *how* and *why*.

---

## The system we build: "Orbit"

Orbit is a small but realistic order-processing platform. It is deliberately
shaped like the systems you run at work: a synchronous HTTP path, an
asynchronous queue path, a relational database, a cache, and background
workers that have to scale on queue depth rather than CPU.

```
                                   ┌──────────────────────── Kubernetes cluster ───────────────────────────┐
                                   │                                                                         │
  client ──HTTPS──► Gateway API ───┼──► catalog (Go, 3+ pods) ──► Redis (cache)                              │
  (k6 / curl)      (Envoy Gateway) │        │                                                                │
                   + cert-manager  │        └──────────────► Postgres (CloudNativePG, 3 instances)           │
                                   │                              ▲                                          │
                                   ├──► orders (Go, 3+ pods) ─────┘                                          │
                                   │        │  XADD orders.created                                           │
                                   │        ▼                                                                │
                                   │     Redis Stream ◄── XREADGROUP ── worker (Go, 0..N pods, KEDA-scaled)  │
                                   │                                                                         │
                                   │  every pod ──OTLP──► OTel Collector (agent DS) ──► OTel gateway ──┐     │
                                   │                                                                    │     │
                                   │      Prometheus ◄─metrics──────────────────────────────────────────┤     │
                                   │      Loki       ◄─logs─────────────────────────────────────────────┤     │
                                   │      Tempo      ◄─traces───────────────────────────────────────────┘     │
                                   │      Grafana (dashboards, SLOs)  Alertmanager ──► Slack / PagerDuty      │
                                   │                                                                         │
                                   │  Argo CD (GitOps)  Argo Rollouts (canary)  KEDA  Kyverno  Velero        │
                                   └─────────────────────────────────────────────────────────────────────────┘
```

| Component | What it does | What it teaches |
|---|---|---|
| `catalog` | `GET /products`, `GET /products/{id}`; read-heavy, Redis cache-aside | stateless scaling, HPA on CPU, cache patterns |
| `orders` | `POST /orders` writes Postgres, publishes an event to a Redis Stream | transactions, idempotency, canary releases, RPS-based scaling |
| `worker` | consumes `orders.created`, "charges payment", marks order paid | queue-driven autoscaling (KEDA, scale-to-zero), retries, poison messages |
| Postgres | system of record (CloudNativePG operator) | StatefulSets, operators, backups, failover |
| Redis | cache + stream | StatefulSet basics, persistence trade-offs |
| Observability | OTel Collector, Prometheus, Loki, Tempo, Grafana | the three signals, correlation, SLOs |
| Platform | Envoy Gateway, cert-manager, Argo CD, Argo Rollouts, KEDA, Kyverno, Velero | the "platform layer" every real cluster grows |

Final repository layout (you build it up chapter by chapter):

```
orbit/
├── go.mod
├── cmd/
│   ├── catalog/main.go
│   ├── orders/main.go
│   └── worker/main.go
├── internal/
│   ├── platform/         # shared: config, logging, telemetry, http server, health
│   ├── catalog/
│   ├── orders/
│   └── worker/
├── migrations/           # SQL migrations (goose)
├── Dockerfile            # one multi-stage Dockerfile, SERVICE build-arg
├── .dockerignore
├── compose.yaml          # local dev stack
├── load/orders.js        # k6 load test
├── lab/                  # Lima VM: create-vm.sh, bootstrap.sh, extract-blocks.py (§0)
├── deploy/
│   ├── base/             # Kustomize base: one folder per service
│   └── overlays/{dev,staging,prod}
├── platform/             # cluster add-ons as Argo CD Applications (Helm values)
└── .github/workflows/ci.yaml
```

---

## Prerequisites and tooling

You need a machine with 4+ CPU cores and 16 GB RAM for the full stack
(8 GB works if you skip the observability stack until Part IV).

There are two ways to set up, and both are supported throughout:

- **Path A (recommended): a Linux lab VM on your Mac.** The Mac only
  needs Lima (`brew install lima`). Every tool below is installed *inside*
  an Ubuntu VM by a script. See **§0** (Part 0) for the full setup and the
  milestone runbook for building the project.
- **Path B: native macOS.** Install the tools directly with Homebrew, as
  below. Everything except §1.1 (which needs a Linux kernel) works.
  Differences are noted where they matter (e.g. reaching kind
  LoadBalancers, §43.1).

Path B tool install (macOS):

```bash
# macOS (Homebrew). Linux users: use your package manager or the upstream binaries.
brew install --cask docker            # Docker Desktop (or: brew install colima docker docker-buildx)
brew install kubectl kind helm kustomize k9s stern kubectx jq yq
brew install trivy cosign             # image scanning and signing
brew install k6                       # load testing
brew install argocd                   # Argo CD CLI
brew install go                       # Go 1.25+

# sanity check
docker version && kubectl version --client && kind version && helm version
```

> **Best practice:** pin tool versions per project (e.g. with `mise` or
> `asdf` and a committed `.tool-versions`). "Works on my machine" bugs in
> Kubernetes are very often a `kubectl`/`helm` version skew. `kubectl` is
> supported within **one minor version** of the API server.

---

## Table of contents

**Part 0 — The Lab Environment**

0. [Lab environment: a Linux VM on your Mac](#0-lab-environment-a-linux-vm-on-your-mac)

**Part I — Containers & Docker**

1. [What a container actually is](#1-what-a-container-actually-is)
2. [Docker architecture and core CLI](#2-docker-architecture-and-core-cli)
3. [Images and Dockerfiles from first principles](#3-images-and-dockerfiles-from-first-principles)
4. [Production-grade Dockerfiles](#4-production-grade-dockerfiles)
5. [Docker networking](#5-docker-networking)
6. [Docker storage: volumes, bind mounts, tmpfs](#6-docker-storage-volumes-bind-mounts-tmpfs)
7. [Docker Compose: the whole system on your laptop](#7-docker-compose-the-whole-system-on-your-laptop)
8. [Docker security and the software supply chain](#8-docker-security-and-the-software-supply-chain)
9. [Debugging containers](#9-debugging-containers)
10. [Registries, tagging, and multi-arch builds](#10-registries-tagging-and-multi-arch-builds)
11. [Part I capstone: Orbit under Compose, fully working](#11-part-i-capstone-orbit-under-compose-fully-working)

**Part II — Kubernetes Core**

12. [Kubernetes architecture: the control loop](#12-kubernetes-architecture-the-control-loop)
13. [Pods: the unit of everything](#13-pods-the-unit-of-everything)
14. [Deployments, ReplicaSets, and rollouts](#14-deployments-replicasets-and-rollouts)
15. [Services, DNS, and the networking model](#15-services-dns-and-the-networking-model)
16. [Configuration: ConfigMaps, Secrets, and the Downward API](#16-configuration-configmaps-secrets-and-the-downward-api)
17. [Storage: Volumes, PVs, PVCs, StorageClasses](#17-storage-volumes-pvs-pvcs-storageclasses)
18. [StatefulSets, DaemonSets, Jobs, and CronJobs](#18-statefulsets-daemonsets-jobs-and-cronjobs)
19. [Namespaces, labels, and organizing a cluster](#19-namespaces-labels-and-organizing-a-cluster)
20. [Part II capstone: Orbit from raw manifests](#20-part-ii-capstone-orbit-from-raw-manifests)

**Part III — Production Workloads**

21. [Health checks: probes done right](#21-health-checks-probes-done-right)
22. [Self-healing and auto-recovery](#22-self-healing-and-auto-recovery)
23. [Resources, QoS, and right-sizing](#23-resources-qos-and-right-sizing)
24. [Scheduling: placement for resilience](#24-scheduling-placement-for-resilience)
25. [Traffic in: Gateway API, TLS, and cert-manager](#25-traffic-in-gateway-api-tls-and-cert-manager)
26. [Packaging: Kustomize and Helm](#26-packaging-kustomize-and-helm)
27. [Autoscaling: pods, sizes, and nodes](#27-autoscaling-pods-sizes-and-nodes)

**Part IV — Observability**

28. [Observability strategy](#28-observability-strategy)
29. [Instrumenting Orbit with OpenTelemetry](#29-instrumenting-orbit-with-opentelemetry)
30. [Deploying the observability stack](#30-deploying-the-observability-stack)
31. [Metrics, dashboards, and alerting](#31-metrics-dashboards-and-alerting)
32. [SLOs and burn-rate alerting](#32-slos-and-burn-rate-alerting)

**Part V — Debugging**

33. [The Kubernetes debugging playbook](#33-the-kubernetes-debugging-playbook)

**Part VI — Security**

34. [Securing the cluster and the workloads](#34-securing-the-cluster-and-the-workloads)

**Part VII — Delivery**

35. [CI: build, test, scan, sign, publish](#35-ci-build-test-scan-sign-publish)
36. [GitOps with Argo CD](#36-gitops-with-argo-cd)
37. [Progressive delivery: canary with automatic rollback](#37-progressive-delivery-canary-with-automatic-rollback)
38. [Database migrations in a continuously deployed system](#38-database-migrations-in-a-continuously-deployed-system)

**Part VIII — Reliability Operations**

39. [Backup and disaster recovery](#39-backup-and-disaster-recovery)
40. [Chaos engineering](#40-chaos-engineering)
41. [Cost and capacity](#41-cost-and-capacity)
42. [Going to the cloud: Orbit on EKS](#42-going-to-the-cloud-orbit-on-eks)

**Part IX — Capstone and Reference**

43. [Capstone: assemble and verify the whole system](#43-capstone-assemble-and-verify-the-whole-system)
44. [Production-readiness checklist](#44-production-readiness-checklist)
45. [Cheat sheets](#45-cheat-sheets)
46. [Pitfalls that cause real outages](#46-pitfalls-that-cause-real-outages)
47. [Further reading: the real-life references](#47-further-reading-the-real-life-references)

---

# Part 0 — The Lab Environment

## 0. Lab environment: a Linux VM on your Mac

You can do most of this guide natively on macOS with Docker Desktop (see
*Prerequisites*, path B). The recommended path is a **dedicated Ubuntu VM**
on your Mac that holds the whole project: Docker, kind, every CLI, the Orbit
repo, and all its clusters. This section builds that VM with a script, so
you can delete and recreate it at any time.

### 0.1 Why a Linux VM instead of native macOS

| | Native macOS (Docker Desktop) | Lima Ubuntu VM (this lab) |
|---|---|---|
| §1.1 namespaces/cgroups by hand | not possible (no Linux kernel) | works natively |
| kind LoadBalancer IPs | often unreachable from the host (needs port-forwards) | reachable inside the VM; `cloud-provider-kind` runs as a service |
| `tc`/`netem`, `tcpdump`, `nsenter`, `crictl`, `/proc` | partial, inside Docker's hidden VM | first-class |
| Matches production nodes | no (macOS userland, BSD tools) | yes (Ubuntu, systemd, cgroup v2, GNU tools) |
| Isolation | tools and clusters mixed into your Mac | everything in one disposable VM; `limactl delete` resets it |
| Resource control | Docker Desktop's VM setting | explicit CPUs/RAM/disk per VM |
| Docker Desktop license | needed for larger companies | not needed (Docker CE inside the VM) |

**Lima** (a CNCF project, `brew install lima`) runs Linux VMs on macOS with
Apple's Virtualization.framework (`vz`). It forwards guest ports to your
Mac's `localhost` automatically, and it's scriptable. Alternatives are in
§0.8.

### 0.2 Size the VM

| Mac RAM | VM CPUs / RAM | What runs comfortably |
|---|---|---|
| 16 GB | 4 / 10 GiB | Parts I–III fully; Part IV with a 2-node kind cluster and Loki/Tempo at minimum sizes |
| 32 GB | 6 / 16 GiB | **the full system** (4-node kind + observability + Argo + chaos): the default below |
| 64 GB+ | 8 / 24 GiB | full system plus a second "DR" cluster for §39.5 |

Disk: 100 GiB. It's a sparse file, so it only uses what is actually
written.

On Apple Silicon the VM is **arm64**. Everything in this guide is
multi-arch (distroless, Postgres, Redis, Envoy, Grafana stack, Argo). For
the odd amd64-only image, the bootstrap enables QEMU emulation through
binfmt.

### 0.3 Create the VM

In the wiki, the lab scripts ship **next to this guide** in `K8s-docker/lab/` (`create-vm.sh`, `bootstrap.sh`,
`extract-blocks.py`). They're identical to the listings below, so you can create the VM
before the Orbit repo exists. Later, copy `lab/` into the repo (the layout above shows it there).
Run this script **on the Mac**, from this directory:

```bash
#!/usr/bin/env bash
# lab/create-vm.sh: [mac] create and provision the Orbit lab VM
set -euo pipefail
VM=${VM:-orbit}
CPUS=${CPUS:-6}
MEMORY=${MEMORY:-16}      # GiB
DISK=${DISK:-100}         # GiB
HERE=$(cd "$(dirname "$0")" && pwd)

command -v limactl >/dev/null || brew install lima

if ! limactl list -q | grep -qx "$VM"; then
  limactl create --name="$VM" \
    --cpus="$CPUS" --memory="$MEMORY" --disk="$DISK" \
    --vm-type=vz --mount-type=virtiofs \
    --containerd=none \
    --tty=false \
    template:ubuntu-24.04
fi
limactl start "$VM"

limactl copy "$HERE/bootstrap.sh" "$VM":/tmp/bootstrap.sh
limactl shell "$VM" bash /tmp/bootstrap.sh

echo
echo "Done. Enter the lab with:   limactl shell $VM"
echo "VS Code Remote-SSH host:    add 'Include ~/.lima/$VM/ssh.config' to ~/.ssh/config, connect to 'lima-$VM'"
```

```bash
chmod +x lab/create-vm.sh && ./lab/create-vm.sh        # [mac] ~5-10 min the first time
limactl list                                            # [mac] NAME orbit, STATUS Running
limactl shell orbit                                     # [mac] -> you are now in Ubuntu
```

Choices in that command:

- `--vm-type=vz` + `virtiofs`: Apple's native hypervisor, the fastest
  option on Apple Silicon and recent Intel Macs.
- `--containerd=none`: Lima's templates can start their own rootless
  containerd/nerdctl. We install the real Docker Engine instead (kind and
  the guide expect `docker`), so we turn Lima's off to avoid two runtimes.
- By default Lima mounts your Mac home directory **read-only** in the VM.
  Keep the Orbit repo on the VM's own disk (`~/orbit` inside the VM). It's
  faster, and it avoids UID/permission mismatches with bind mounts.

> ⚠️ **Caveat (verify):** the flags above (`--cpus`, `--memory`, `--disk`,
> `--vm-type`, `--mount-type`, `--containerd=none`, `--tty=false`) and
> `template:ubuntu-24.04` were **checked against Lima 2.2.1** (`limactl create --help`,
> `--list-templates`). The VM itself hasn't been booted with them yet. Older Lima 1.x
> spells templates `template://ubuntu-24.04`. If `limactl create` rejects something
> on your version, run `limactl create --help` and `limactl create --list-templates`,
> or put the same settings in a YAML file (`cpus`, `memory`, `disk`,
> `vmType: vz`, `containerd: {system: false, user: false}`) and run
> `limactl create --name=orbit orbit.yaml`. The `~/.lima/<vm>/ssh.config`
> path is where current Lima versions write it; check with
> `limactl list --format '{{.Dir}}'`.

### 0.4 Bootstrap the toolchain inside the VM

`lab/bootstrap.sh` is idempotent: safe to re-run to upgrade tools. It
detects the architecture, so the same script works on Intel Macs and on
Linux cloud VMs.

```bash
#!/usr/bin/env bash
# lab/bootstrap.sh: run INSIDE the VM as your normal user (uses sudo). Idempotent.
set -euo pipefail
ARCH=$(dpkg --print-architecture)                  # arm64 | amd64
BIN=/usr/local/bin
log() { printf '\n\033[1;34m>> %s\033[0m\n' "$*"; }
gh_latest() { curl -fsSL "https://api.github.com/repos/$1/releases/latest" | jq -r .tag_name; }
install_url() { sudo curl -fsSLo "$BIN/$1" "$2" && sudo chmod +x "$BIN/$1"; }

log "base packages"
sudo apt-get update -y
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y \
  ca-certificates curl git jq make unzip build-essential bash-completion \
  postgresql-client redis-tools dnsutils iproute2 tcpdump netcat-openbsd

log "kernel settings (multi-node kind, Loki/Tempo, chaos)"
sudo tee /etc/sysctl.d/99-orbit-lab.conf >/dev/null <<'EOF'
fs.inotify.max_user_watches=524288
fs.inotify.max_user_instances=512
vm.max_map_count=262144
EOF
sudo sysctl --system >/dev/null

log "Docker Engine (CE)"
if ! command -v docker >/dev/null; then
  curl -fsSL https://get.docker.com | sudo sh
fi
sudo usermod -aG docker "$USER"
# emulate amd64 for the rare single-arch image (registration does not survive a reboot; re-run bootstrap)
sudo docker run --privileged --rm tonistiigi/binfmt --install amd64 >/dev/null

log "Go"
GO_VERSION=${GO_VERSION:-$(curl -fsSL 'https://go.dev/VERSION?m=text' | head -1)}
if [ "$(/usr/local/go/bin/go version 2>/dev/null | awk '{print $3}')" != "$GO_VERSION" ]; then
  sudo rm -rf /usr/local/go
  curl -fsSL "https://go.dev/dl/${GO_VERSION}.linux-${ARCH}.tar.gz" | sudo tar -C /usr/local -xz
fi
grep -q '/usr/local/go/bin' ~/.bashrc || echo 'export PATH=/usr/local/go/bin:$HOME/go/bin:$PATH' >> ~/.bashrc
export PATH=/usr/local/go/bin:$HOME/go/bin:$PATH

log "kubectl + helm"
install_url kubectl "https://dl.k8s.io/release/$(curl -fsSL https://dl.k8s.io/release/stable.txt)/bin/linux/${ARCH}/kubectl"
curl -fsSL https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

log "Go-installable CLIs"
for pkg in \
  sigs.k8s.io/kind@latest \
  sigs.k8s.io/cloud-provider-kind@latest \
  sigs.k8s.io/kustomize/kustomize/v5@latest \
  github.com/stern/stern@latest \
  github.com/ahmetb/kubectx/cmd/kubectx@latest \
  github.com/ahmetb/kubectx/cmd/kubens@latest \
  github.com/mikefarah/yq/v4@latest \
  github.com/yannh/kubeconform/cmd/kubeconform@latest \
  github.com/pressly/goose/v3/cmd/goose@latest \
  go.k6.io/k6@latest ; do
  echo "   go install $pkg"; go install "$pkg"
done

log "release-binary CLIs"
install_url argocd                "https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-${ARCH}"
install_url kubectl-argo-rollouts "https://github.com/argoproj/argo-rollouts/releases/latest/download/kubectl-argo-rollouts-linux-${ARCH}"
install_url cosign                "https://github.com/sigstore/cosign/releases/latest/download/cosign-linux-${ARCH}"
curl -fsSL "https://github.com/derailed/k9s/releases/latest/download/k9s_Linux_${ARCH}.tar.gz" | sudo tar -C $BIN -xz k9s
curl -fsSL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sudo sh -s -- -b $BIN
V=$(gh_latest vmware-tanzu/velero)
curl -fsSL "https://github.com/vmware-tanzu/velero/releases/download/${V}/velero-${V}-linux-${ARCH}.tar.gz" \
  | sudo tar -C $BIN -xz --strip-components=1 "velero-${V}-linux-${ARCH}/velero"
curl -fsSL https://github.com/cloudnative-pg/cloudnative-pg/raw/main/hack/install-cnpg-plugin.sh | sudo sh -s -- -b $BIN

log "cloud-provider-kind as a service (LoadBalancer IPs for kind)"
sudo install -m 0755 "$HOME/go/bin/cloud-provider-kind" $BIN/
sudo tee /etc/systemd/system/cloud-provider-kind.service >/dev/null <<'EOF'
[Unit]
Description=cloud-provider-kind: LoadBalancer services for kind clusters
After=docker.service
Requires=docker.service

[Service]
ExecStart=/usr/local/bin/cloud-provider-kind
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
sudo systemctl daemon-reload
sudo systemctl enable --now cloud-provider-kind

log "shell conveniences"
grep -q 'orbit-lab' ~/.bashrc || cat >> ~/.bashrc <<'EOF'
# orbit-lab
source <(kubectl completion bash); alias k=kubectl; complete -o default -F __start_kubectl k
source <(helm completion bash)
export KUBE_EDITOR=vim
EOF

log "versions"
docker version --format 'docker {{.Server.Version}}' 2>/dev/null || sudo docker version --format 'docker {{.Server.Version}}'
go version; kubectl version --client | head -1; helm version --short; kind version
echo "Open a NEW shell (limactl shell orbit) so the docker group and PATH apply."
```

Check it from a new shell:

```bash
limactl shell orbit                          # [mac]
docker run --rm hello-world                  # Docker works without sudo
kind version && kubectl version --client && k9s version --short
systemctl status cloud-provider-kind --no-pager | head -3
```

> ⚠️ **Caveat (verify):** the bootstrap pulls **latest** releases, which is
> convenient for a lab but not reproducible. Once it works, pin versions
> (write the output of the "versions" step into `lab/versions.env` and
> install from it). Things to check: release asset names (`k9s_Linux_arm64.tar.gz`,
> `argocd-linux-arm64`, Velero's tarball layout); `get-helm-3` (Helm 4 has
> its own installer if you move to Helm 4); and `go install` paths for
> tools that restructure their modules. Building `k6` with `go install`
> takes a few minutes. If it fails, use Grafana's apt repo or
> `docker run --rm -i --network host grafana/k6 run - < load/orders.js`.

### 0.5 Daily workflow in the lab

**Editing code.** Use VS Code **Remote-SSH** into `lima-orbit` and open
`~/orbit`. The editor UI runs on the Mac; files, terminals, Go tooling and
Docker run in the VM. (Alternatively, edit on the Mac in a writable Lima
mount. It's slower for large trees and has UID quirks.)

**Reaching services from your Mac browser.** Lima forwards any port the VM
listens on at `127.0.0.1`/`0.0.0.0` to your Mac's `localhost`.
`kubectl port-forward` inside the VM is therefore enough:

```bash
# inside the VM (run each in its own terminal, or append & to background them)
kubectl -n monitoring port-forward svc/kps-grafana 3000:80          # Mac: http://localhost:3000
kubectl -n argocd port-forward svc/argocd-server 8088:80            # Mac: http://localhost:8088
kubectl -n argo-rollouts port-forward svc/argo-rollouts-dashboard 3100:3100
GW_SVC=$(kubectl -n envoy-gateway-system get svc \
  -l gateway.envoyproxy.io/owning-gateway-name=public -o jsonpath='{.items[0].metadata.name}')
kubectl -n envoy-gateway-system port-forward svc/$GW_SVC 8443:443   # Mac: https://orbit.localtest.me:8443
```

Inside the VM, LoadBalancer IPs from `cloud-provider-kind` work directly:
`curl -sk --resolve orbit.localtest.me:443:$GW https://orbit.localtest.me/api/products`
(§25.3).

**Lifecycle.**

```bash
limactl stop orbit            # [mac] pause the lab (frees RAM); kind clusters come back with Docker on start
limactl start orbit           # [mac]
limactl delete orbit && ./lab/create-vm.sh     # [mac] full reset: ~10 min, the VM is cattle too
kind delete cluster --name orbit               # [vm]  reset just the cluster
```

If a kind cluster misbehaves after a VM restart (stale node IPs, CoreDNS
crash-looping), delete and recreate it. Because everything is in Git and
driven by Argo CD (§36), that takes minutes. This is §39's "clusters are
cattle" principle applied to your laptop.

> ⚠️ **Caveat (verify):** `limactl snapshot` has historically worked only
> with the QEMU VM type, not `vz`. Lima 2.x also has `limactl clone`
> (stop the VM, then `limactl clone orbit orbit-golden`) to keep a known-good
> copy, but test it before relying on it. The scripts plus Git are the
> dependable reset path.

### 0.6 Verifying this guide as you go

Most content in this guide is unexecuted until you run it. The lab is
where you check it, in four layers, cheapest first:

**1. Static checks.** Extract the fenced code blocks and validate each
with its own tool:

```python
#!/usr/bin/env python3
# lab/extract-blocks.py: dump every fenced block of the guide into lab/blocks/NNN.<lang>
import pathlib, re, sys
src = pathlib.Path(sys.argv[1]).read_text()
out = pathlib.Path("lab/blocks"); out.mkdir(parents=True, exist_ok=True)
for i, (lang, body) in enumerate(re.findall(r"^```(\w*)\n(.*?)^```", src, re.S | re.M)):
    (out / f"{i:03d}.{lang or 'txt'}").write_text(body)
print(f"wrote {i + 1} blocks to {out}")
```

```bash
python3 lab/extract-blocks.py ~/wiki/K8s-docker/real-life-k8s-guide.md
# Kubernetes YAML (including CRDs, via the community CRD schema catalog)
kubeconform -strict -summary -ignore-missing-schemas \
  -schema-location default \
  -schema-location 'https://raw.githubusercontent.com/datreeio/CRDs-catalog/main/{{.Group}}/{{.ResourceKind}}_{{.ResourceAPIVersion}}.json' \
  lab/blocks/*.yaml
# Go: assemble the snippets into ~/orbit as you go through §4, §11, §29, then:
cd ~/orbit && go mod tidy && go build ./... && go vet ./...
# Dockerfiles
docker build --check . ; docker run --rm -i hadolint/hadolint < Dockerfile
# Prometheus rules (take .spec out of each PrometheusRule)
yq '.spec' rule.yaml > /tmp/rules.yaml && docker run --rm -v /tmp:/w --entrypoint promtool prom/prometheus check rules /w/rules.yaml
# Collector configs (take the `config:` block out of the Helm values)
yq '.config' platform/otel-gateway/values.yaml > /tmp/otel.yaml && \
  docker run --rm -v /tmp:/w otel/opentelemetry-collector-contrib validate --config=/w/otel.yaml
# Helm values -> rendered manifests -> schema check
helm template kps prometheus-community/kube-prometheus-stack -f platform/kube-prometheus-stack/values.yaml | kubeconform -strict -ignore-missing-schemas
# Compose, Terraform
docker compose config -q ; (cd infra/envs/prod && terraform init -backend=false && terraform validate)
```

(Many YAML blocks are excerpts with `# ...` placeholders, so expect some
failures that are just incomplete snippets. The complete manifests you
assemble in `~/orbit/deploy` are what must pass.)

**2. Server-side dry runs.** Once a CRD's operator is installed,
`kubectl apply --dry-run=server -f file.yaml` checks the YAML against the
**exact** CRD schema and admission webhooks in your cluster. That is the
definitive check for every ⚠️ box about CRD fields.

**3. Behavior.** Run the milestone checks in §0.7 and the acceptance tests
in §43.2. For anything involving metrics, check that the query returns
data in Grafana Explore before you trust an alert, a KEDA trigger, or a
canary analysis built on it.

**4. Sources.** For each ⚠️ box, compare with the project's own docs or
release notes for the version you pinned, and record the result in
`lab/VERIFIED.md` (date, version, pass/fix). Re-run layer 1 in CI on the
repo so regressions show up when you bump versions.

### 0.7 Lab runbook: build the full project, milestone by milestone

Everything below runs **inside the VM**. Each milestone has a hard "done
when" check. Don't move on until it passes.

| # | Milestone | Chapters | Key commands | Done when |
|---|---|---|---|---|
| M0 | Lab VM ready | §0.3–0.4 | `./lab/create-vm.sh` | `docker run hello-world` works without sudo; `kind version` prints |
| M1 | Container by hand | §1.1 | `unshare ... chroot` + cgroup `memory.max` | you're PID 1 in a new namespace; 100 MB allocation gets `Killed` |
| M2 | Orbit code + images | §4, §11 | `go build ./...`, `docker build --build-arg SERVICE=...` | three images ~10 MB, non-root (`docker inspect -f '{{.Config.User}}'`) |
| M3 | Compose stack | §7, §11 | `docker compose --profile obs up -d --build --wait` | order goes `pending` → `paid`; trace visible in otel-lgtm Grafana (localhost:3000 on the Mac) |
| M4 | kind cluster | §12.4 | `kind create cluster --config kind-orbit.yaml` | 3 workers, one per zone label; `kubectl top nodes` works |
| M5 | Raw manifests | §18, §20 | CNPG, Redis StatefulSet, `kubectl apply -f deploy/raw/` | §20 exercises 1–4 behave as described |
| M6 | Production workloads | §21–24 | golden manifest for all 3 services | rolling restart during `k6` shows 0 failed requests |
| M7 | Gateway + TLS | §25 | Envoy Gateway, cert-manager, HTTPRoute | `curl -sk https://orbit.localtest.me/api/products` via the LB IP |
| M8 | Kustomize | §26 | `deploy/base` + `overlays/dev` | `kubectl diff -k deploy/overlays/dev` is empty after apply |
| M9 | Autoscaling | §27 | KEDA, HPA, `k6 run load/orders.js` | orders scales out; worker 0 → N → 0 |
| M10 | Observability | §29–32 | kps, Loki, Tempo, Collector agent + gateway | alert → dashboard → trace → logs without typing a pod name |
| M11 | Security | §34 | PSA restricted, NetworkPolicies, Kyverno | acceptance tests 2, 3, 12 (§43.2) |
| M12 | CI + GitOps | §35–36 | push the repo to GitHub (Actions + GHCR); Argo CD in kind pulls it | a `git push` reaches the dev namespace with no manual step |
| M13 | Canary | §37 | Argo Rollouts + analysis | a bad release aborts by itself at ≤ 25% |
| M14 | Migrations | §38 | PreSync migration Job | expand/contract rename done with zero errors under load |
| M15 | Backup/DR | §39 | MinIO in-cluster as "S3", CNPG backups, Velero | delete the namespace, restore, measure RTO |
| M16 | Chaos | §40 | Chaos Mesh experiments H1–H3 | hypotheses hold, or you've filed fixes |
| M17 | Capstone | §43 | all 14 acceptance tests | all pass on a freshly created cluster |

M12 note: GitHub Actions builds the images on GitHub's runners and pushes
them to GHCR. kind pulls them from there (make the packages public, or
create an imagePullSecret, §16.3). For a fully offline lab, run a local
registry next to kind (kind's docs: *Local Registry*) and point Argo CD at
a Git server in the cluster (e.g. Gitea), but the GitHub path is closer
to real life.

### 0.8 Alternatives to Lima

| Tool | Notes |
|---|---|
| **Multipass** (Canonical) | `multipass launch 24.04 --name orbit --cpus 6 --memory 16G --disk 100G --cloud-init cloud-init.yaml`; then run the same `bootstrap.sh`. Simple; Ubuntu only |
| **UTM** | GUI on top of QEMU/Apple Virtualization; good if you want a desktop Linux; manual setup |
| **OrbStack** | fast, polished Linux machines (`orb create ubuntu orbit`) and Docker; needs a paid license for business use |
| **Colima** | built on Lima, focused on providing a Docker/containerd runtime rather than a general-purpose VM |
| **Cloud VM** | an 8 vCPU / 32 GB Linux VM (e.g. EC2 `m7g.2xlarge`) running the same `bootstrap.sh`: useful if your Mac is short on RAM. Stop it when idle |

---

# Part I — Containers & Docker

## 1. What a container actually is

A container is **not** a lightweight VM. It is an ordinary Linux process
that the kernel has been told to lie to. Three kernel features do the lying:

| Feature | What it isolates/limits | Seen from inside the container |
|---|---|---|
| **Namespaces** | what the process can *see*: PIDs, mounts, network, hostname (UTS), IPC, users, cgroups | "I am PID 1", "I have my own `eth0`", "`/` is my image" |
| **cgroups (v2)** | what the process can *use*: CPU, memory, IO, PIDs | "I get killed at 256 MiB", "I get throttled above 0.5 CPU" |
| **Layered filesystem (overlayfs)** | a root filesystem assembled from read-only image layers + one writable layer | "`/usr/bin/myapp` exists", writes vanish when the container is removed |

Plus some hardening: **capabilities** (root split into ~40 privileges, most
dropped), **seccomp** (syscall filter), and **AppArmor/SELinux** (mandatory
access control).

```
  VM                                    Container
  ┌──────────────┐                      ┌──────────────┐ ┌──────────────┐
  │ app          │                      │ app          │ │ app          │
  │ libs         │                      │ libs (image) │ │ libs (image) │
  │ guest kernel │                      └──────┬───────┘ └──────┬───────┘
  ├──────────────┤                             │ namespaces + cgroups
  │ hypervisor   │                      ┌──────┴────────────────┴───────┐
  ├──────────────┤                      │        ONE host kernel        │
  │ host kernel  │                      └───────────────────────────────┘
  └──────────────┘
```

Consequences you will rely on later:

- **Containers start in milliseconds** because nothing boots; a process is
  just exec'd.
- **The kernel is shared.** A kernel exploit escapes every container on the
  host. That is why we drop capabilities, run as non-root, and use seccomp.
- **On macOS/Windows, Docker runs a Linux VM** behind the scenes. Your
  containers share *that* VM's kernel, not macOS's.
- **Memory limits are enforced by the kernel OOM killer**, not by your
  runtime. Exceed the cgroup limit and the process gets `SIGKILL`: exit code
  137, no chance to clean up.

### 1.1 Build a "container" by hand (Linux only)

Doing this once makes everything else make sense. On a Linux VM:

```bash
# 1. get a root filesystem (an "image") - here, Alpine's minirootfs
mkdir -p /tmp/rootfs && cd /tmp/rootfs
curl -sSL https://dl-cdn.alpinelinux.org/alpine/v3.20/releases/x86_64/alpine-minirootfs-3.20.3-x86_64.tar.gz | tar xz

# 2. start a shell in new PID, mount, UTS, IPC and network namespaces, chrooted into it
sudo unshare --pid --fork --mount-proc=/tmp/rootfs/proc \
             --mount --uts --ipc --net \
             chroot /tmp/rootfs /bin/sh

# inside:
hostname orbit-box && hostname   # changes only this UTS namespace
ps aux                           # you are PID 1; host processes are invisible
ip addr                          # only a down loopback: new network namespace
```

Now add a cgroup memory limit from another terminal on the host:

```bash
sudo mkdir /sys/fs/cgroup/demo
echo 50M | sudo tee /sys/fs/cgroup/demo/memory.max
echo <PID-of-the-unshare-shell> | sudo tee /sys/fs/cgroup/demo/cgroup.procs
# inside the box, allocate 100 MB and watch the kernel kill it:
#   head -c 100m /dev/zero | tail     -> Killed
```

That is, give or take, what `runc` does when Docker or Kubernetes starts a
container.

> ⚠️ **Caveat (verify):** the Alpine URL pins `3.20.3` and `x86_64`. If it 404s, take the current file name from `dl-cdn.alpinelinux.org/alpine/latest-stable/releases/<arch>/`. In the Apple Silicon lab VM (§0), use `aarch64`. The cgroup path assumes cgroup v2 (`stat -fc %T /sys/fs/cgroup` prints `cgroup2fs`).

### 1.2 The standards: OCI, containerd, runc

```
docker CLI ──► dockerd ──► containerd ──► containerd-shim ──► runc ──► your process
kubectl ──► API server ──► kubelet ──CRI──► containerd ──► shim ──► runc ──► your process
```

- **OCI image spec**: what an image is (a manifest + config + tarball
  layers, content-addressed by SHA-256). Any OCI image runs on any OCI
  runtime: Docker, containerd, CRI-O, Podman.
- **OCI runtime spec**: how to run one (`runc`, `crun`, `gVisor`'s `runsc`,
  Kata Containers).
- **CRI**: the gRPC interface the kubelet uses to talk to a runtime.
  Kubernetes removed its Docker shim in 1.24; clusters run **containerd** or
  **CRI-O** directly. Images built with `docker build` work unchanged
  because they are OCI images.

> **Real-world example:** "We removed Docker from our nodes, will our images
> break?" was a common panic in 2022. They didn't: only things that talked to
> `/var/run/docker.sock` *on the node* (old log agents, Docker-in-Docker CI)
> broke. Your build tool and your node runtime are independent choices.

---

## 2. Docker architecture and core CLI

### 2.1 The moving parts

| Piece | Role |
|---|---|
| `docker` CLI | client; talks to the daemon over a Unix socket or TCP |
| `dockerd` | the daemon: images, networks, volumes, the API |
| **BuildKit** | the build engine behind `docker build`/`docker buildx` (parallel stages, cache mounts, secrets, multi-arch) |
| `containerd` | image pull/unpack, container lifecycle |
| registry | stores images: Docker Hub, GHCR, ECR, GAR, ACR, Harbor |

> **War story:** access to `/var/run/docker.sock` is **root on the host**.
> Anyone who can `docker run -v /:/host` owns the machine. Teams that mount
> the socket into CI jobs or "monitoring" containers have handed out root.
> Treat the socket like an SSH key to root.

### 2.2 The 20 commands you will use daily

```bash
# images
docker pull nginx:1.27-alpine
docker images                                  # or: docker image ls
docker image history nginx:1.27-alpine         # layers and their sizes
docker image inspect nginx:1.27-alpine | jq '.[0].Config'
docker rmi <image>                             # remove
docker image prune -a                          # remove unused images

# containers
docker run -d --name web -p 8080:80 nginx:1.27-alpine   # detached, port 8080 -> 80
docker ps                                      # running; -a for all
docker logs -f --tail 100 web                  # follow logs
docker exec -it web sh                         # shell inside a running container
docker stop web                                # SIGTERM, then SIGKILL after 10s
docker kill web                                # SIGKILL now
docker rm -f web                               # stop + remove
docker inspect web | jq '.[0].State'           # status, exit code, OOMKilled flag
docker stats                                   # live CPU/mem/net/io per container
docker top web                                 # processes inside, seen from the host
docker cp web:/etc/nginx/nginx.conf ./         # copy files out (or in)

# housekeeping
docker system df                               # what is using disk
docker system prune                            # remove stopped containers, dangling images, unused networks
```

### 2.3 The `run` flags that matter in production

```bash
docker run -d \
  --name orders \
  --restart unless-stopped \           # auto-recovery: restart on crash and after daemon restart
  --memory 256m --memory-swap 256m \   # hard limit, no swap
  --cpus 0.5 \                         # CFS quota: half a core
  --pids-limit 200 \                   # fork-bomb protection
  --read-only --tmpfs /tmp \           # immutable root filesystem
  --cap-drop ALL \                     # no Linux capabilities
  --security-opt no-new-privileges \   # setuid binaries can't escalate
  --user 65532:65532 \                 # non-root
  --health-cmd 'wget -qO- http://localhost:8080/readyz || exit 1' \
  --health-interval 10s --health-retries 3 \
  -e DATABASE_URL \                    # pass through from your shell, not hard-coded
  -p 127.0.0.1:8080:8080 \             # bind to localhost only, not 0.0.0.0
  ghcr.io/acme/orbit-orders:1.4.2
```

Restart policies: `no` (default), `on-failure[:max]`, `always`,
`unless-stopped`. This is single-host auto-recovery. Kubernetes replaces it
with controllers (§22).

> **Best practice:** publish ports to `127.0.0.1` on dev machines and
> servers. `-p 5432:5432` on a cloud VM binds to every interface, and
> **Docker's iptables rules bypass `ufw`**. Many leaked databases were
> exposed exactly like this.

---

## 3. Images and Dockerfiles from first principles

### 3.1 Layers and the build cache

Each filesystem-changing instruction (`RUN`, `COPY`, `ADD`) creates a layer.
Layers are cached by **instruction + inputs**. Once one layer's cache is
invalidated, *every later layer is rebuilt*. Hence the golden rule:

> **Best practice:** order instructions from **least to most frequently
> changing**. Copy dependency manifests and install dependencies *before*
> copying source code.

```dockerfile
# BAD: any source change re-downloads every dependency
COPY . .
RUN go mod download && go build -o /app ./cmd/orders

# GOOD: dependency layer is cached until go.mod/go.sum change
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app ./cmd/orders
```

Deleting files in a later layer **does not shrink the image**: the bytes
still exist in the earlier layer. This is also how secrets leak:

```dockerfile
# LEAKS: the key is in layer N even though layer N+1 deletes it
COPY id_rsa /root/.ssh/id_rsa
RUN git clone git@github.com:acme/private.git && rm /root/.ssh/id_rsa
```

`docker history --no-trunc` or `dive` (a layer explorer) will show it to
anyone who pulls the image. Use BuildKit secret mounts instead (§4.4).

### 3.2 Instruction reference with the gotchas

| Instruction | Use | Gotcha |
|---|---|---|
| `FROM image:tag AS name` | base image, names a stage | pin by tag *and* digest in prod: `FROM golang:1.25@sha256:...` |
| `ARG` | build-time variable | visible in `docker history`; never put secrets in it |
| `ENV` | runtime env var | baked into the image; also visible |
| `WORKDIR` | sets cwd (creates it) | prefer over `RUN cd ...` |
| `COPY` | copy from context/stage | prefer over `ADD`; use `--chown`, `--chmod`, `--link` |
| `ADD` | COPY + URL fetch + tar auto-extract | surprising; use only for the auto-extract |
| `RUN` | execute at build | chain and clean in the same layer: `apt-get update && apt-get install -y --no-install-recommends x && rm -rf /var/lib/apt/lists/*` |
| `USER` | run as this uid | use numeric `65532:65532` so Kubernetes `runAsNonRoot` can verify it |
| `EXPOSE` | documentation only | does not publish anything |
| `HEALTHCHECK` | Docker-level health | **ignored by Kubernetes**; use probes there |
| `ENTRYPOINT` | the executable | use exec form `["..."]` |
| `CMD` | default args to ENTRYPOINT | overridden by `docker run image args` |
| `STOPSIGNAL` | signal sent on stop | default `SIGTERM`; nginx uses `SIGQUIT` for graceful stop |

### 3.3 PID 1, signals, and shell form vs exec form

This one bug causes thousands of "my container takes 10 seconds to stop" and
"in-flight requests fail on every deploy" tickets.

```dockerfile
# shell form: runs `/bin/sh -c "node server.js"`. sh is PID 1.
CMD node server.js
# exec form: node IS PID 1 and receives SIGTERM directly.
CMD ["node", "server.js"]
```

With shell form, `docker stop` / Kubernetes sends `SIGTERM` to `sh`, which
**does not forward it** to your app. After the grace period (10 s in Docker,
30 s in Kubernetes) the app is `SIGKILL`ed mid-request.

Also, **PID 1 has no default signal handlers** in Linux: a process that
doesn't explicitly handle `SIGTERM` ignores it when it is PID 1. And PID 1
must reap zombie child processes. Options:

1. Exec form + an app that handles `SIGTERM` (Go: `signal.NotifyContext`).
2. If your app spawns children (shell scripts, Puppeteer, legacy apps), use a
   tiny init: `docker run --init` or `ENTRYPOINT ["tini", "--", "app"]`.
3. If you need a shell wrapper script, end it with `exec "$@"` so the app
   replaces the shell.

```sh
#!/bin/sh
# docker-entrypoint.sh
set -eu
# ... templating, waiting, etc.
exec "$@"        # app becomes PID 1, receives signals
```

### 3.4 `ENTRYPOINT` + `CMD` together

```dockerfile
ENTRYPOINT ["/app/orders"]
CMD ["--port=8080"]
# docker run img                 -> /app/orders --port=8080
# docker run img --port=9090     -> /app/orders --port=9090
```

In Kubernetes, `command:` overrides `ENTRYPOINT` and `args:` overrides
`CMD`. Mixing that up is a classic.

---

## 4. Production-grade Dockerfiles

### 4.1 Start the Orbit codebase

Before containerizing, we need something to containerize. The shared
`platform` package is written once and reused by all three services; it
already contains the hooks (health, graceful shutdown, telemetry) that later
chapters rely on.

```bash
mkdir -p orbit && cd orbit
go mod init github.com/acme/orbit
go get github.com/jackc/pgx/v5 github.com/redis/go-redis/v9
```

`internal/platform/config.go`: twelve-factor config from the environment.

```go
package platform

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Env returns the value of key or def if unset.
func Env(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

// MustEnv returns the value of key or exits: fail fast on missing config,
// so a misconfigured pod crashes at start instead of serving errors.
func MustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		fmt.Fprintf(os.Stderr, "missing required env var %s\n", key)
		os.Exit(1)
	}
	return v
}

func EnvDuration(key string, def time.Duration) time.Duration {
	if d, err := time.ParseDuration(os.Getenv(key)); err == nil {
		return d
	}
	return def
}

func EnvInt(key string, def int) int {
	if n, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return n
	}
	return def
}
```

`internal/platform/health.go`: separate liveness and readiness (§21 explains
why they must differ).

```go
package platform

import (
	"context"
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"
)

// Check is a named readiness dependency check.
type Check struct {
	Name string
	Fn   func(ctx context.Context) error
}

type Health struct {
	draining atomic.Bool
	checks   []Check
}

func NewHealth(checks ...Check) *Health { return &Health{checks: checks} }

// SetDraining flips readiness to false during shutdown so the load balancer
// stops sending new traffic before the server closes.
func (h *Health) SetDraining() { h.draining.Store(true) }

// Livez answers "is this process wedged?". It never checks dependencies:
// a database outage must not cause every pod to be restarted.
func (h *Health) Livez(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// Readyz answers "should I receive traffic right now?".
func (h *Health) Readyz(w http.ResponseWriter, r *http.Request) {
	if h.draining.Load() {
		http.Error(w, "draining", http.StatusServiceUnavailable)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 800*time.Millisecond)
	defer cancel()
	failed := map[string]string{}
	for _, c := range h.checks {
		if err := c.Fn(ctx); err != nil {
			failed[c.Name] = err.Error()
		}
	}
	if len(failed) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(failed)
		return
	}
	_, _ = w.Write([]byte("ok"))
}
```

`internal/platform/server.go`: an HTTP server with timeouts and graceful
shutdown that follows the Kubernetes termination sequence (§13.5).

```go
package platform

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type ServerConfig struct {
	Addr            string        // ":8080"
	DrainDelay      time.Duration // time to keep serving after SIGTERM while endpoints update
	ShutdownTimeout time.Duration // max time to finish in-flight requests
}

// Run serves h until SIGTERM/SIGINT, then drains gracefully.
func Run(cfg ServerConfig, h http.Handler, health *Health) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second, // slowloris protection
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       90 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("http server listening", "addr", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err // failed to bind etc.
	case <-ctx.Done():
	}

	// 1. Fail readiness so endpoints/LBs remove us.
	health.SetDraining()
	slog.Info("shutdown signal received, draining", "delay", cfg.DrainDelay)
	// 2. Keep serving while that propagates (kube-proxy, ingress, LB).
	time.Sleep(cfg.DrainDelay)
	// 3. Stop accepting, finish in-flight requests.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	slog.Info("shutdown complete")
	return nil
}
```

`cmd/catalog/main.go` (first version; we add the database and telemetry
later):

```go
package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/acme/orbit/internal/platform"
)

var version = "dev" // set with -ldflags at build time

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	health := platform.NewHealth()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /livez", health.Livez)
	mux.HandleFunc("GET /readyz", health.Readyz)
	mux.HandleFunc("GET /products", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1, "name": "Rocket", "price_cents": 9900},
		})
	})
	mux.HandleFunc("GET /version", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(version))
	})

	err := platform.Run(platform.ServerConfig{
		Addr:            ":" + platform.Env("PORT", "8080"),
		DrainDelay:      platform.EnvDuration("DRAIN_DELAY", 5*time.Second),
		ShutdownTimeout: platform.EnvDuration("SHUTDOWN_TIMEOUT", 20*time.Second),
	}, mux, health)
	if err != nil {
		slog.Error("server failed", "err", err)
		os.Exit(1)
	}
}
```

### 4.2 The production Dockerfile (Go, multi-stage, distroless, non-root)

One Dockerfile builds every service; `SERVICE` selects which `cmd/` to
compile. This is the standard monorepo pattern.

```dockerfile
# syntax=docker/dockerfile:1.7
# ---- build stage --------------------------------------------------------
ARG GO_VERSION=1.25
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-bookworm AS build
WORKDIR /src

# 1) dependencies: cached until go.mod/go.sum change
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# 2) source + build. TARGETOS/TARGETARCH are set by buildx for multi-arch.
COPY . .
ARG SERVICE
ARG VERSION=dev
ARG TARGETOS TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" \
      -o /out/app ./cmd/${SERVICE}

# ---- runtime stage ------------------------------------------------------
# distroless/static: CA certs, tzdata, /etc/passwd with a nonroot user. No shell, no package manager.
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/app /app
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/app"]
```

Line by line, why it is production-grade:

- `# syntax=docker/dockerfile:1.7` pins the Dockerfile frontend, so builds
  are reproducible across Docker versions.
- `--platform=$BUILDPLATFORM` runs the compiler natively and
  **cross-compiles** for `TARGETARCH`. Multi-arch (amd64 + arm64) builds
  take seconds, without QEMU emulation.
- **Cache mounts** keep the Go module and build caches across builds
  without putting them in a layer.
- `CGO_ENABLED=0` gives a static binary, so it runs on `distroless/static`
  or `scratch`.
- `-trimpath -s -w` gives reproducible paths and drops debug symbols
  (smaller). Keep symbols if you rely on `pprof` symbolization in prod;
  Go's runtime still has function names for stack traces either way.
- **distroless:nonroot**: no shell means an attacker with RCE has no `sh`,
  `curl` or `apt`. CVE scanners report very few findings. Includes CA
  certificates (needed for outbound TLS) and tzdata.
- `USER 65532:65532`: numeric, so Kubernetes can enforce `runAsNonRoot`.

Build and run:

```bash
docker build --build-arg SERVICE=catalog --build-arg VERSION=0.1.0 -t orbit-catalog:0.1.0 .
docker run --rm -p 127.0.0.1:8080:8080 orbit-catalog:0.1.0
curl localhost:8080/products
docker images orbit-catalog      # ~10 MB
```

Typical image sizes for the same Go service:

| Base | Size | Shell? | Typical CVE count |
|---|---|---|---|
| `golang:1.25` (single-stage) | ~900 MB | yes | dozens–hundreds |
| `debian:bookworm-slim` | ~80 MB | yes | some |
| `alpine:3.20` | ~15 MB | yes (busybox) | few |
| `gcr.io/distroless/static:nonroot` | ~10 MB | no | ~0 |
| `scratch` | ~8 MB | no | 0 (but no CA certs/tzdata/users unless you copy them) |

> **Real-world example:** a platform team mandates distroless or Chainguard
> base images for all services. The CVE backlog drops from thousands of
> "fix available" findings to nearly zero, and patching becomes "rebuild
> weekly" instead of triaging tickets. The cost: you need `kubectl debug`
> (§33.3) instead of `kubectl exec ... sh`.

### 4.3 Python equivalent (for comparison)

```dockerfile
# syntax=docker/dockerfile:1.7
FROM python:3.12-slim AS build
ENV PIP_NO_CACHE_DIR=1 PIP_DISABLE_PIP_VERSION_CHECK=1
WORKDIR /app
RUN python -m venv /venv
ENV PATH=/venv/bin:$PATH
COPY requirements.txt .
RUN --mount=type=cache,target=/root/.cache/pip pip install -r requirements.txt
COPY src/ ./src/

FROM python:3.12-slim
ENV PATH=/venv/bin:$PATH PYTHONUNBUFFERED=1 PYTHONDONTWRITEBYTECODE=1
RUN useradd --uid 10001 --no-create-home app
COPY --from=build /venv /venv
COPY --from=build --chown=10001:10001 /app/src /app/src
WORKDIR /app
USER 10001
EXPOSE 8000
# gunicorn as PID 1 handles SIGTERM gracefully (finishes requests up to --graceful-timeout)
ENTRYPOINT ["gunicorn", "src.main:app", "-k", "uvicorn.workers.UvicornWorker", \
            "-b", "0.0.0.0:8000", "--workers", "2", "--graceful-timeout", "20"]
```

`PYTHONUNBUFFERED=1` matters: without it, logs sit in a buffer and appear
late (or never, if the container is killed).

### 4.4 Secrets at build time: never `ARG`, always `--mount=type=secret`

```dockerfile
# private Go modules: a token is needed during download only
RUN --mount=type=secret,id=gh_token \
    git config --global url."https://x-access-token:$(cat /run/secrets/gh_token)@github.com/".insteadOf "https://github.com/" \
 && go mod download
```

```bash
docker build --secret id=gh_token,env=GH_TOKEN -t orbit-orders .
```

The secret is mounted only for that `RUN` and never written to a layer. One
caveat: the `git config` line above *does* write the token to
`/root/.gitconfig` in that layer. This is fine in a build stage that is
discarded, but never do it in the final stage. For SSH:
`RUN --mount=type=ssh go mod download` with `docker build --ssh default`.

### 4.5 `.dockerignore`: the build context

`docker build .` sends the context to the builder. Without a
`.dockerignore`, you ship `.git`, `node_modules`, local `.env` files
(secrets!), and test fixtures, and every change to them busts the
`COPY . .` cache.

```gitignore
# .dockerignore
.git
.github
**/*.md
**/.env*
**/*_test.go
deploy/
load/
bin/
tmp/
compose*.yaml
Dockerfile*
```

> **Best practice:** use an allow-list for large repos: `*` first, then
> `!go.mod`, `!go.sum`, `!cmd/`, `!internal/`.

### 4.6 Reproducibility and provenance

- Pin base images by digest. Let Renovate or Dependabot bump the digests
  with PRs so you still get patches:
  `FROM gcr.io/distroless/static-debian12:nonroot@sha256:...`
- Tag images with the **git SHA** (immutable) and optionally semver. **Never
  deploy `:latest`**: you cannot tell what's running or roll back to it.
- Attach an SBOM and provenance at build time:
  `docker buildx build --sbom=true --provenance=mode=max ...` (§8.3).
- Add OCI labels so tools can trace an image back to its source:

```dockerfile
LABEL org.opencontainers.image.source="https://github.com/acme/orbit" \
      org.opencontainers.image.revision="${GIT_SHA}"
```

### 4.7 Dockerfile linting

```bash
docker build --check .                         # BuildKit's built-in checks (Docker 27+)
docker run --rm -i hadolint/hadolint < Dockerfile
```

Hadolint catches missing `--no-install-recommends`, unpinned apt packages,
shell-form `CMD`, `cd` instead of `WORKDIR`, and more. Run it in CI (§35).

> ⚠️ **Caveat (verify):** `docker build --check` needs a recent Docker Engine/Buildx (introduced around Docker 27). If the flag is unknown, upgrade or rely on hadolint. The image-size and CVE-count table gives typical orders of magnitude, not measurements; check yours with `docker images` and `trivy image`.

---

## 5. Docker networking

### 5.1 Network drivers

| Driver | What it is | Use for |
|---|---|---|
| `bridge` (default `docker0`) | private L2 bridge on the host + NAT; containers get 172.17.x.x | quick one-offs; **no DNS between containers** on the default bridge |
| user-defined `bridge` | same, but with **embedded DNS** (container name → IP) and isolation | every multi-container app (Compose creates one automatically) |
| `host` | no network namespace; container uses the host's stack | max performance, network tools; port conflicts become your problem |
| `none` | loopback only | batch jobs that must not touch the network |
| `overlay` | VXLAN across hosts (Swarm) | multi-host Swarm; in Kubernetes the CNI plugin does this job |
| `macvlan`/`ipvlan` | container gets an IP on the physical LAN | legacy apps that need to look like a physical host |

```bash
docker network create orbit-net
docker run -d --name redis --network orbit-net redis:7-alpine
docker run --rm --network orbit-net redis:7-alpine redis-cli -h redis ping   # PONG: DNS by name
docker network inspect orbit-net | jq '.[0].Containers'
```

### 5.2 How port publishing works

`-p 8080:80` installs a DNAT rule in iptables (or nftables): traffic to
`host:8080` is rewritten to `container_ip:80`. From inside a container,
`localhost` is **the container itself**, not your laptop. To reach a
service on the host from a container, use `host.docker.internal` (Docker
Desktop, or add `--add-host=host.docker.internal:host-gateway` on Linux).

> **War story:** "The app works locally but can't connect to the DB in
> Docker." The app was configured with `DB_HOST=localhost`. In a container
> that means the container's own loopback. Use the service name (`db`) on a
> shared network.

### 5.3 Debugging container networking

```bash
# netshoot: a container full of network tools, joined to another container's network namespace
docker run --rm -it --network container:orders nicolaka/netshoot
#   ss -tlnp            -> what is listening
#   dig redis           -> DNS
#   curl -v http://catalog:8080/readyz
#   tcpdump -i eth0 port 5432
```

The `--network container:<name>` trick is the same idea as the ephemeral
debug containers you'll use in Kubernetes (§33.3).

---

## 6. Docker storage: volumes, bind mounts, tmpfs

| Type | Syntax | Managed by | Use for |
|---|---|---|---|
| **Named volume** | `-v pgdata:/var/lib/postgresql/data` | Docker (`/var/lib/docker/volumes`) | database data, anything that must outlive the container |
| **Bind mount** | `-v $(pwd)/src:/app/src` | you (host path) | dev hot-reload, injecting config files |
| **tmpfs** | `--tmpfs /tmp:size=64m` | memory | scratch space with a read-only root FS, secrets that must not hit disk |

```bash
docker volume create pgdata
docker run -d --name db -e POSTGRES_PASSWORD=dev -v pgdata:/var/lib/postgresql/data postgres:17
docker rm -f db                      # container gone ...
docker run -d --name db -e POSTGRES_PASSWORD=dev -v pgdata:/var/lib/postgresql/data postgres:17
# ... data survives

# backup a volume: tar it from a throwaway container
docker run --rm -v pgdata:/data -v $(pwd):/backup alpine tar czf /backup/pgdata.tgz -C /data .
```

Gotchas:

- **The writable container layer is not for data.** It is slow
  (copy-on-write) and deleted with the container.
- **Bind-mount UID mismatch:** a container running as uid 65532 can't write
  to a host dir owned by uid 501. Fix ownership, or use named volumes.
- **Never back up a live database by copying its files.** Use `pg_dump` or
  the database's own backup tool; a file copy of a running DB can be
  corrupt.
- On macOS, bind mounts cross the VM boundary and are slower. Use
  Compose `develop.watch` (§7.3) or synchronized file shares for big trees.

---

## 7. Docker Compose: the whole system on your laptop

Compose is the right tool for **local development and CI integration
tests**. It is not a production orchestrator for multi-host systems.

### 7.1 Orbit's `compose.yaml`

```yaml
# compose.yaml
name: orbit

x-service-defaults: &svc
  restart: unless-stopped
  read_only: true
  tmpfs: [/tmp]
  cap_drop: [ALL]
  security_opt: [no-new-privileges:true]
  environment: &svc-env
    OTEL_EXPORTER_OTLP_ENDPOINT: http://lgtm:4317
    OTEL_RESOURCE_ATTRIBUTES: service.namespace=orbit,deployment.environment.name=local
    DATABASE_URL: postgres://orbit:orbit@db:5432/orbit?sslmode=disable
    REDIS_ADDR: redis:6379
  depends_on: &svc-deps
    db:    { condition: service_healthy }
    redis: { condition: service_healthy }
    migrate: { condition: service_completed_successfully }

services:
  catalog:
    <<: *svc
    build: { context: ., args: { SERVICE: catalog } }
    environment: { <<: *svc-env, OTEL_SERVICE_NAME: catalog }
    ports: ["127.0.0.1:8081:8080"]
    develop:
      watch:
        - { action: rebuild, path: ./internal }
        - { action: rebuild, path: ./cmd/catalog }

  orders:
    <<: *svc
    build: { context: ., args: { SERVICE: orders } }
    environment: { <<: *svc-env, OTEL_SERVICE_NAME: orders }
    ports: ["127.0.0.1:8082:8080"]
    develop:
      watch:
        - { action: rebuild, path: ./internal }
        - { action: rebuild, path: ./cmd/orders }

  worker:
    <<: *svc
    build: { context: ., args: { SERVICE: worker } }
    environment: { <<: *svc-env, OTEL_SERVICE_NAME: worker }
    deploy: { replicas: 2 }

  migrate:
    image: ghcr.io/pressly/goose:latest   # pin a version in real projects
    command: ["-dir", "/migrations", "postgres", "postgres://orbit:orbit@db:5432/orbit?sslmode=disable", "up"]
    volumes: ["./migrations:/migrations:ro"]
    depends_on: { db: { condition: service_healthy } }
    restart: "no"

  db:
    image: postgres:17
    environment:
      POSTGRES_USER: orbit
      POSTGRES_PASSWORD: orbit
      POSTGRES_DB: orbit
    volumes: [pgdata:/var/lib/postgresql/data]
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U orbit -d orbit"]
      interval: 5s
      timeout: 3s
      retries: 10
    ports: ["127.0.0.1:5432:5432"]

  redis:
    image: redis:7-alpine
    command: ["redis-server", "--appendonly", "yes"]
    volumes: [redisdata:/data]
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      retries: 10

  # Grafana + Prometheus + Loki + Tempo + OTel Collector in ONE container, for local dev only
  lgtm:
    image: grafana/otel-lgtm:latest
    ports: ["127.0.0.1:3000:3000", "127.0.0.1:4317:4317", "127.0.0.1:4318:4318"]
    profiles: [obs]

volumes:
  pgdata:
  redisdata:
```

What this demonstrates:

- **YAML anchors (`x-` extension fields + `<<: *svc`)** keep the hardening
  flags DRY.
- **`depends_on` with conditions.** `service_healthy` waits for the
  healthcheck, and `service_completed_successfully` waits for the migration
  job. Plain `depends_on` only waits for the container to *start*, which is
  the source of "connection refused on boot" flakiness.
- **Profiles:** `docker compose --profile obs up` adds the observability
  stack only when you want it.
- **Production hardening even locally** (read-only, cap_drop) catches "app
  writes to `/`" bugs on day one instead of in the cluster.

> ⚠️ **Caveat (verify):** check that the `ghcr.io/pressly/goose` image exists and pin a version. If it doesn't, build a two-line image (`FROM golang` + `go install github.com/pressly/goose/v3/cmd/goose@<version>`, copied into distroless) and use that. `grafana/otel-lgtm` is a dev-only image; pin a tag. Validate the file with `docker compose config`.

### 7.2 Migrations

`migrations/00001_init.sql` (goose format):

```sql
-- +goose Up
CREATE TABLE products (
  id          BIGSERIAL PRIMARY KEY,
  name        TEXT NOT NULL,
  price_cents INTEGER NOT NULL CHECK (price_cents >= 0),
  stock       INTEGER NOT NULL DEFAULT 0
);
CREATE TABLE orders (
  id              UUID PRIMARY KEY,
  idempotency_key TEXT UNIQUE NOT NULL,
  product_id      BIGINT NOT NULL REFERENCES products(id),
  quantity        INTEGER NOT NULL CHECK (quantity > 0),
  status          TEXT NOT NULL DEFAULT 'pending',   -- pending | paid | failed
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX orders_status_idx ON orders(status);
INSERT INTO products(name, price_cents, stock) VALUES
  ('Rocket', 9900, 100), ('Satellite', 49900, 20), ('Telescope', 19900, 50);

-- +goose Down
DROP TABLE orders;
DROP TABLE products;
```

### 7.3 Daily workflow

```bash
docker compose up -d --build          # build + start everything
docker compose --profile obs up -d    # include Grafana/LGTM at http://localhost:3000
docker compose watch                  # rebuild on file change (develop.watch)
docker compose ps                     # status + health
docker compose logs -f orders worker  # tail selected services
docker compose exec db psql -U orbit  # shell into the DB
docker compose up -d --scale worker=5 # scale a service
docker compose down                   # stop (keeps volumes)
docker compose down -v                # stop AND delete volumes (fresh DB)
```

> **Best practice:** Compose is also your **integration test harness** in
> CI. `docker compose up -d --wait` blocks until every healthcheck passes,
> then you run the tests against it. It is the same stack developers run,
> so CI failures reproduce locally.

---

## 8. Docker security and the software supply chain

Container security comes down to three questions. What's **in** the image
(vulnerabilities, secrets)? **Who built it** (provenance, signatures)? And
what can it **do** at runtime (privileges)?

### 8.1 Runtime hardening checklist

| Control | Docker flag | Kubernetes field (§34) |
|---|---|---|
| Non-root user | `--user 65532` | `runAsNonRoot: true`, `runAsUser` |
| Drop capabilities | `--cap-drop ALL` | `capabilities.drop: ["ALL"]` |
| No privilege escalation | `--security-opt no-new-privileges` | `allowPrivilegeEscalation: false` |
| Read-only root FS | `--read-only` | `readOnlyRootFilesystem: true` |
| Seccomp | default profile is on | `seccompProfile.type: RuntimeDefault` |
| Resource limits | `--memory --cpus --pids-limit` | `resources.limits` |
| Never | `--privileged`, `-v /var/run/docker.sock`, `--pid=host`, `--net=host` | `privileged`, `hostPID`, `hostNetwork`, hostPath |

**Rootless mode** (`dockerd-rootless-setuptool.sh install`, or Podman by
default) runs the daemon itself as a normal user, so a container escape
lands as an unprivileged user.

### 8.2 Scanning images

```bash
trivy image --severity HIGH,CRITICAL --ignore-unfixed orbit-catalog:0.1.0
trivy image --exit-code 1 --severity CRITICAL orbit-catalog:0.1.0   # CI gate
trivy fs --scanners secret,misconfig .                              # secrets + Dockerfile/K8s misconfig
docker scout cves orbit-catalog:0.1.0                               # Docker's scanner
```

> **Best practice:** gate CI on **CRITICAL with a fix available**, not on
> every finding. A gate nobody can pass gets disabled. Track the rest with
> a scheduled re-scan of what's actually *running*: new CVEs are published
> against images you shipped last month.

### 8.3 SBOMs, provenance, and signing

```bash
# Build with SBOM + SLSA provenance attestations attached to the image in the registry
docker buildx build --platform linux/amd64,linux/arm64 \
  --sbom=true --provenance=mode=max \
  -t ghcr.io/acme/orbit-orders:$(git rev-parse --short HEAD) --push .

# Inspect the SBOM later (e.g. "are we affected by the new CVE in libfoo?")
docker buildx imagetools inspect ghcr.io/acme/orbit-orders:abc123 --format '{{ json .SBOM }}'

# Sign by DIGEST with Sigstore cosign (keyless in CI via OIDC; key-based shown here)
cosign generate-key-pair
cosign sign --key cosign.key ghcr.io/acme/orbit-orders@sha256:<digest>
cosign verify --key cosign.pub ghcr.io/acme/orbit-orders@sha256:<digest>
```

In CI you'll use **keyless signing**: the GitHub Actions OIDC token proves
"this image was built by workflow X in repo Y", and in Part VI a Kyverno
policy refuses to run anything that isn't signed that way (§34.6).

> **Real-world example:** after the 2020 SolarWinds and 2021 Codecov
> incidents, "prove where this binary came from" became a compliance
> requirement (US EO 14028, SLSA). SBOM + provenance + signature
> verification at admission is now the expected baseline at regulated
> companies.

> ⚠️ **Caveat (verify):** the `imagetools inspect --format '{{ json .SBOM }}'` template depends on the Buildx version and on the SBOM having been attached. If it prints nothing, use `docker buildx imagetools inspect <ref> --format '{{ json . }}'` or `docker scout sbom <ref>`.

---

## 9. Debugging containers

### 9.1 Exit codes you must recognize

| Code | Meaning | Typical cause |
|---|---|---|
| 0 | clean exit | a "server" that exits 0 means it thought it was done (wrong CMD, batch mode) |
| 1 | app error | panic, unhandled exception, bad config (check logs) |
| 125 | Docker itself failed | bad `docker run` flag |
| 126 | command not executable | missing `+x`, wrong architecture binary |
| 127 | command not found | typo in ENTRYPOINT, missing shell in distroless (`sh -c` on distroless!) |
| 137 | 128 + 9 = **SIGKILL** | **OOM killed** (check `OOMKilled: true`) or killed after the stop grace period |
| 139 | 128 + 11 = SIGSEGV | native crash, CGO, wrong libc (glibc binary on Alpine/musl) |
| 143 | 128 + 15 = SIGTERM | stopped normally by `docker stop`/Kubernetes |

```bash
docker inspect orders --format '{{.State.ExitCode}} OOMKilled={{.State.OOMKilled}} {{.State.Error}}'
```

### 9.2 The debugging ladder

```bash
docker logs --tail 200 orders                  # 1. what did it say before dying?
docker inspect orders | jq '.[0].State, .[0].Config.Env, .[0].Mounts'   # 2. config as actually applied
docker events --since 10m --filter container=orders   # 3. die/oom/restart events
docker stats --no-stream                       # 4. resources
docker exec -it orders sh                      # 5. inside (not possible on distroless...)

# 6. ...so attach a toolbox that shares the target's namespaces
docker run --rm -it --pid container:orders --network container:orders nicolaka/netshoot
#   ps aux, ls /proc/1/root/app, ss -tnp, curl localhost:8080/readyz

# 7. container exits too fast to exec? override the entrypoint
docker run --rm -it --entrypoint sh orbit-orders:debug     # with a debug image variant
# 8. or snapshot its filesystem
docker create --name snap orbit-orders:0.1.0 && docker export snap | tar -t | head
```

> **Best practice:** build a `:debug` variant only when you need it
> (`gcr.io/distroless/static-debian12:debug-nonroot` has busybox). Don't
> ship shells in production images "just in case". §33.3 shows ephemeral
> containers, which make that unnecessary.

### 9.3 "Works locally, fails in the container" checklist

1. Binding to `127.0.0.1` inside the container. Bind to `0.0.0.0` (or
   `:8080` in Go) or the port publish can't reach it.
2. `localhost` used to reach another container.
3. Architecture mismatch: an arm64 Mac building images for amd64 clusters
   (`exec format error`). Build with `--platform`.
4. Missing CA certificates (`x509: certificate signed by unknown authority`)
   on `scratch`.
5. Timezone or locale assumptions (`tzdata` missing).
6. File permissions with a non-root user.
7. App writes to a read-only root FS. Give it a `tmpfs`/`emptyDir` for that
   path.

---

## 10. Registries, tagging, and multi-arch builds

### 10.1 Tagging strategy

| Tag | Mutable? | Use |
|---|---|---|
| `sha-3f9c2a1` | no | what you deploy; traceable to a commit |
| `1.4.2` | should not be | releases |
| `1.4`, `1` | yes | convenience for consumers of public images |
| `latest` | yes | never in deployment manifests |
| `@sha256:...` (digest) | **immutable by definition** | what GitOps should pin; what you sign |

> **Best practice:** enable **tag immutability** in your registry (ECR,
> GAR, Harbor support it), so `1.4.2` can never silently point to a
> different image.

### 10.2 Multi-arch with buildx

```bash
docker buildx create --name orbit-builder --use --bootstrap
docker buildx build --platform linux/amd64,linux/arm64 \
  --build-arg SERVICE=orders \
  --cache-from type=registry,ref=ghcr.io/acme/orbit-orders:buildcache \
  --cache-to   type=registry,ref=ghcr.io/acme/orbit-orders:buildcache,mode=max \
  -t ghcr.io/acme/orbit-orders:sha-$(git rev-parse --short HEAD) --push .
docker buildx imagetools inspect ghcr.io/acme/orbit-orders:sha-...   # shows the manifest list
```

Multi-arch matters because arm64 nodes (AWS Graviton, GCP Axion/Tau T2A,
Azure Cobalt) are typically 20–40% cheaper per unit of performance. A
multi-arch image lets the scheduler place a pod on either.

The **registry cache** (`--cache-to type=registry`) gives ephemeral CI
runners a warm cache. It commonly takes a 6-minute build down to under
1 minute.

### 10.3 Private registry credentials

```bash
echo "$GHCR_TOKEN" | docker login ghcr.io -u <user> --password-stdin
# Kubernetes pulls via an imagePullSecret (§16) or, on cloud, node/workload identity (better: no static creds)
```

---

## 11. Part I capstone: Orbit under Compose, fully working

Now complete the three services. These are the final versions; Part IV adds
telemetry with a few lines in `main.go`.

`internal/orders/handler.go`: idempotent order creation with an outbox-style
publish.

```go
package orders

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const Stream = "orders.created"

type Handler struct {
	DB  *pgxpool.Pool
	RDB *redis.Client
}

type createReq struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Idempotency-Key")
	if key == "" {
		http.Error(w, "Idempotency-Key header required", http.StatusBadRequest)
		return
	}
	var req createReq
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&req); err != nil || req.Quantity <= 0 {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	id, created, err := h.insert(r.Context(), key, req)
	if err != nil {
		slog.ErrorContext(r.Context(), "create order failed", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if created {
		// Publish after commit. If this fails, a reconciler (or an outbox table) re-publishes
		// pending orders; the worker is idempotent, so duplicates are harmless.
		if err := h.RDB.XAdd(r.Context(), &redis.XAddArgs{
			Stream: Stream, MaxLen: 100_000, Approx: true,
			Values: map[string]any{"order_id": id},
		}).Err(); err != nil {
			slog.WarnContext(r.Context(), "publish failed; will be re-driven", "order_id", id, "err", err)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	if created {
		w.WriteHeader(http.StatusCreated)
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"id": id, "status": "pending"})
}

// insert returns the existing order for a repeated idempotency key instead of creating a duplicate.
func (h *Handler) insert(ctx context.Context, key string, req createReq) (string, bool, error) {
	id := uuid.NewString()
	err := h.DB.QueryRow(ctx,
		`INSERT INTO orders (id, idempotency_key, product_id, quantity)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (idempotency_key) DO NOTHING
		 RETURNING id`, id, key, req.ProductID, req.Quantity).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) { // conflict: already exists
		err = h.DB.QueryRow(ctx, `SELECT id FROM orders WHERE idempotency_key=$1`, key).Scan(&id)
		return id, false, err
	}
	return id, err == nil, err
}
```

`internal/worker/worker.go`: a consumer group with acknowledgements,
retry via pending-entry reclaim, and graceful stop.

```go
package worker

import (
	"context"
	"errors"
	"log/slog"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Worker struct {
	DB       *pgxpool.Pool
	RDB      *redis.Client
	Stream   string
	Group    string
	Consumer string // pod name: unique per replica
}

func (w *Worker) Run(ctx context.Context) error {
	// Create the group once; BUSYGROUP means it already exists.
	if err := w.RDB.XGroupCreateMkStream(ctx, w.Stream, w.Group, "0").Err(); err != nil &&
		!strings.HasPrefix(err.Error(), "BUSYGROUP") {
		return err
	}
	for ctx.Err() == nil {
		// Reclaim messages another (crashed) consumer held for > 60s: at-least-once delivery.
		claimed, _, _ := w.RDB.XAutoClaim(ctx, &redis.XAutoClaimArgs{
			Stream: w.Stream, Group: w.Group, Consumer: w.Consumer,
			MinIdle: time.Minute, Start: "0-0", Count: 10,
		}).Result()
		for _, m := range claimed {
			w.handle(ctx, m)
		}

		res, err := w.RDB.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group: w.Group, Consumer: w.Consumer, Streams: []string{w.Stream, ">"},
			Count: 10, Block: 5 * time.Second,
		}).Result()
		if errors.Is(err, redis.Nil) || errors.Is(err, context.Canceled) {
			continue
		}
		if err != nil {
			slog.Error("xreadgroup", "err", err)
			time.Sleep(time.Second)
			continue
		}
		for _, s := range res {
			for _, m := range s.Messages {
				w.handle(ctx, m)
			}
		}
	}
	return nil
}

func (w *Worker) handle(ctx context.Context, m redis.XMessage) {
	orderID, _ := m.Values["order_id"].(string)
	// Use a context that survives SIGTERM so an in-flight message can finish.
	hctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
	defer cancel()

	time.Sleep(time.Duration(50+rand.IntN(150)) * time.Millisecond) // "call the payment provider"

	// Idempotent state transition: only pending -> paid.
	if _, err := w.DB.Exec(hctx,
		`UPDATE orders SET status='paid', updated_at=now() WHERE id=$1 AND status='pending'`, orderID); err != nil {
		slog.Error("process order failed; will be retried", "order_id", orderID, "err", err)
		return // not ACKed -> reclaimed later by XAUTOCLAIM
	}
	w.RDB.XAck(hctx, w.Stream, w.Group, m.ID)
	slog.Info("order paid", "order_id", orderID)
}
```

Production notes embedded in this code (each one maps to a later chapter):

- **Idempotency keys** make client retries safe. Retries happen constantly
  in Kubernetes: rolling updates, pod evictions, LB timeouts.
- **At-least-once + idempotent consumer** is the realistic delivery
  contract. Exactly-once is a property you build, not one you get from the
  queue.
- **`Consumer` = pod name** (from the Downward API, §16.4), so reclaimed
  messages can be traced to the pod that dropped them.
- A message that fails forever would loop forever. In production, check
  the delivery count (`XPENDING`) and move a message to a **dead-letter
  stream** after N attempts. Left as an exercise; the pattern is in §43.

The `main.go` files wire it together. `cmd/orders/main.go` (the worker's
main is the same shape: it builds a `worker.Worker` and calls `Run(ctx)`
under `signal.NotifyContext` instead of `platform.Run`):

```go
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/acme/orbit/internal/orders"
	"github.com/acme/orbit/internal/platform"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	ctx := context.Background()

	pcfg, err := pgxpool.ParseConfig(platform.MustEnv("DATABASE_URL"))
	if err != nil {
		slog.Error("bad DATABASE_URL", "err", err)
		os.Exit(1)
	}
	pcfg.MaxConns = int32(platform.EnvInt("DB_MAX_CONNS", 10)) // x replicas must stay < Postgres max_connections
	db, err := pgxpool.NewWithConfig(ctx, pcfg)                 // lazy: does not fail if the DB is down at boot
	if err != nil {
		slog.Error("db pool", "err", err)
		os.Exit(1)
	}
	defer db.Close()
	rdb := redis.NewClient(&redis.Options{Addr: platform.MustEnv("REDIS_ADDR")})
	defer rdb.Close()

	health := platform.NewHealth(
		platform.Check{Name: "postgres", Fn: db.Ping},
		platform.Check{Name: "redis", Fn: func(ctx context.Context) error { return rdb.Ping(ctx).Err() }},
	)
	h := &orders.Handler{DB: db, RDB: rdb}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /livez", health.Livez)
	mux.HandleFunc("GET /readyz", health.Readyz)
	mux.HandleFunc("POST /orders", h.Create)

	if err := platform.Run(platform.ServerConfig{
		Addr:            ":" + platform.Env("PORT", "8080"),
		DrainDelay:      platform.EnvDuration("DRAIN_DELAY", 5*time.Second),
		ShutdownTimeout: platform.EnvDuration("SHUTDOWN_TIMEOUT", 20*time.Second),
	}, mux, health); err != nil {
		slog.Error("server failed", "err", err)
		os.Exit(1)
	}
}
```

Note the **boot behavior**: the process starts even if Postgres is down
(the pool connects lazily), stays *alive*, and reports *not ready* until
the dependency is reachable. That is exactly what Kubernetes probes expect.
A service that crashes on boot when its DB is unavailable produces
crash-loop storms during every database failover.

Run it:

```bash
docker compose --profile obs up -d --build --wait
curl -s localhost:8081/products | jq
curl -s -XPOST localhost:8082/orders -H 'Idempotency-Key: demo-1' \
     -H 'Content-Type: application/json' -d '{"product_id":1,"quantity":2}' | jq
docker compose exec db psql -U orbit -c 'select id,status from orders'   # pending -> paid within ~1s
```

**Part I checkpoint:** a 10 MB distroless, non-root, read-only image per
service; a reproducible local stack with healthchecks and migrations; and
you can explain what happens on `docker stop`.

> ⚠️ **Caveat (verify):** the Go code in §4.1 and §11 was written for this guide and has not been compiled. Assemble the repo and run `go mod tidy && go build ./... && go vet ./...` before trusting it. Known spots to check: the return values of go-redis v9 `XAutoClaim`, and the `uuid` dependency (`go get github.com/google/uuid`).

---

# Part II — Kubernetes Core

## 12. Kubernetes architecture: the control loop

### 12.1 The one idea: declarative desired state + reconciliation

You don't tell Kubernetes *what to do*. You tell it *what should be true*
("3 replicas of orders:1.4.2 exist"). Controllers then loop forever:

```
        ┌──────────────────────────────────────────────┐
        │  observe actual state  (watch the API server) │
        │            │                                  │
        │            ▼                                  │
        │  diff against desired state (spec vs status)  │
        │            │                                  │
        │            ▼                                  │
        │  act to converge (create/delete/update)       │
        └────────────┴──────────────────────────────────┘
                  repeat forever
```

Every "self-healing" behavior is just this loop. A pod dies, so actual is
2, desired is 3, and a new pod is created. A node disappears, its pods are
marked gone, and they get replaced elsewhere. When you read Kubernetes
objects, read **`spec`** (what you asked for) vs **`status`** (what the
controller observed).

### 12.2 Components

```
                       CONTROL PLANE                                          WORKER NODE (xN)
  ┌───────────────────────────────────────────────────────┐        ┌──────────────────────────────────┐
  │ kube-apiserver  ◄──── kubectl / controllers / kubelets │        │ kubelet   ── CRI ──► containerd  │
  │      │  (the only thing that talks to etcd)            │◄──────►│   │ runs pods, probes, reports   │
  │      ▼                                                 │        │ kube-proxy (or eBPF CNI)         │
  │ etcd (Raft, 3 or 5 members: the cluster's database)    │        │   │ Service VIP -> pod IPs       │
  │ kube-scheduler (assigns pods to nodes)                 │        │ CNI plugin (pod network)         │
  │ kube-controller-manager (Deployment, ReplicaSet, Node, │        │ CSI node plugin (volumes)        │
  │   Job, EndpointSlice, ... controllers)                 │        └──────────────────────────────────┘
  │ cloud-controller-manager (LBs, node lifecycle on cloud)│
  └───────────────────────────────────────────────────────┘
```

| Component | If it dies... |
|---|---|
| API server | no changes possible; **running pods keep running** |
| etcd (quorum lost) | API read-only/unavailable; running pods keep running |
| scheduler | new pods stay `Pending` |
| controller-manager | no self-healing, no rollouts, no endpoint updates |
| kubelet on a node | after ~40 s the node is `NotReady`; pods are evicted after a toleration timeout (default 300 s) |
| CoreDNS | service discovery fails: the most common "everything is broken" incident |

> **Real-world example:** managed Kubernetes (EKS, GKE, AKS) runs the
> control plane for you. You operate the nodes, the add-ons, and your
> workloads. Etcd backup, API server scaling and upgrades become the
> provider's problem. That's why almost everyone uses managed control
> planes in production.

### 12.3 What happens on `kubectl apply -f deployment.yaml`

1. `kubectl` sends the object to the **API server**: authn → authz (RBAC) →
   **mutating admission** (defaults, sidecar injection) → schema validation →
   **validating admission** (policies) → persisted to **etcd**.
2. The **Deployment controller** sees a new Deployment and creates a
   **ReplicaSet**.
3. The **ReplicaSet controller** creates 3 **Pod** objects (with no node
   assigned yet).
4. The **scheduler** filters nodes (resources, taints, affinity), scores
   them, and binds each pod to a node.
5. That node's **kubelet** sees the pod, asks the CNI for an IP, the CSI
   for volumes, and containerd to pull and start containers. Then it runs
   probes.
6. Once the readiness probe passes, the **EndpointSlice controller** adds
   the pod IP to the Service, and **kube-proxy** programs the node's
   routing.

Every step is an independent controller watching the API. Nothing calls
anything directly. When something is stuck, ask **"which controller should
have acted, and what does the object's `status`/events say?"**

### 12.4 A production-like local cluster with kind

kind ("Kubernetes in Docker") runs each node as a container. Use a
multi-node cluster with zone labels so scheduling, anti-affinity and
topology spread behave like a real cloud.

```yaml
# kind-orbit.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: orbit
nodes:
  - role: control-plane
  - role: worker
    labels: { topology.kubernetes.io/zone: zone-a }
  - role: worker
    labels: { topology.kubernetes.io/zone: zone-b }
  - role: worker
    labels: { topology.kubernetes.io/zone: zone-c }
```

```bash
kind create cluster --config kind-orbit.yaml
kubectl get nodes -L topology.kubernetes.io/zone

# LoadBalancer Services on kind: run cloud-provider-kind in another terminal
go install sigs.k8s.io/cloud-provider-kind@latest
sudo cloud-provider-kind

# metrics-server (kubectl top, HPA). kind's kubelet certs are self-signed, hence the flag.
helm repo add metrics-server https://kubernetes-sigs.github.io/metrics-server/
helm upgrade --install metrics-server metrics-server/metrics-server -n kube-system \
  --set args='{--kubelet-insecure-tls}'

# load locally built images into the cluster nodes (no registry needed)
for s in catalog orders worker; do
  docker build --build-arg SERVICE=$s -t orbit-$s:dev . && kind load docker-image orbit-$s:dev --name orbit
done
```

Alternatives: **k3d** (k3s in Docker, very fast), **minikube** (many
drivers and add-ons), **OrbStack/Docker Desktop's built-in Kubernetes**
(single node). Use kind because it is what Kubernetes itself uses for CI
and supports multi-node.

> ⚠️ **Caveat (verify):** `cloud-provider-kind` flags and its macOS behavior change between releases. On macOS the LoadBalancer IP may not be reachable from the host; port-forward instead (§43.1), or use the Linux lab VM (§0), where it works natively. The metrics-server `--set args='{...}'` syntax varies by chart version; if it's rejected, put `args:` in a values file. Multi-node kind needs raised inotify limits on Linux hosts (`fs.inotify.max_user_watches=524288`, `max_user_instances=512`); the §0 bootstrap sets them.

### 12.5 kubectl, efficiently

```bash
# contexts and namespaces
kubectl config get-contexts
kubectx kind-orbit           # switch cluster
kubens orbit                 # switch default namespace

# discovery: the built-in documentation
kubectl api-resources                        # every type, short names, namespaced?
kubectl explain deployment.spec.strategy --recursive
kubectl explain pod.spec.containers.livenessProbe

# read
kubectl get pods -o wide                     # node + IP
kubectl get pods -l app.kubernetes.io/name=orders -w   # watch by label
kubectl get deploy orders -o yaml | yq '.status'
kubectl describe pod <pod>                   # EVENTS at the bottom: read them first
kubectl get events --sort-by=.lastTimestamp -A | tail -30
kubectl get pods -o custom-columns=NAME:.metadata.name,RESTARTS:.status.containerStatuses[0].restartCount,NODE:.spec.nodeName

# write: always declarative in real life
kubectl apply -f manifests/ --server-side
kubectl diff -f manifests/                   # what WOULD change (use before every apply)

# generate YAML instead of hand-writing from memory
kubectl create deployment demo --image=nginx --replicas=2 --dry-run=client -o yaml > demo.yaml

# interact
kubectl logs deploy/orders -f --all-containers --since=10m
stern orders -n orbit                        # tail all pods matching "orders", color-coded
kubectl exec -it deploy/orders -- sh         # (not on distroless; see §33.3)
kubectl port-forward svc/orders 8082:80
k9s                                          # terminal UI; learn it, you'll live in it
```

> **Best practice:** in production, humans don't `kubectl apply`. Git +
> Argo CD does (§36). Humans use `kubectl` read-only, plus break-glass
> actions that are audited. Practice that discipline on kind too.

---

## 13. Pods: the unit of everything

### 13.1 Anatomy

A Pod is one or more containers that share a **network namespace** (same
IP, talk via `localhost`), share **volumes**, and are scheduled together on
one node. Usually one app container plus optional helpers.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: hello
  labels: { app: hello }
spec:
  containers:
    - name: web
      image: nginx:1.27-alpine
      ports: [{ containerPort: 80, name: http }]
      resources:
        requests: { cpu: 50m, memory: 32Mi }
        limits:   { memory: 64Mi }
```

```bash
kubectl apply -f hello.yaml
kubectl get pod hello -o wide
kubectl port-forward pod/hello 8080:80
kubectl delete pod hello      # gone forever: nothing recreates a bare Pod
```

That last line is the point: **never run bare Pods in production.** Use a
controller (Deployment, StatefulSet, Job) that recreates them.

### 13.2 Pod lifecycle and statuses

Phases: `Pending` → `Running` → `Succeeded` / `Failed` (plus `Unknown` when
the node is unreachable). The interesting detail is in **container
states** and **conditions**:

| You see (`kubectl get pods`) | Meaning |
|---|---|
| `Pending` | not scheduled (resources, taints, affinity, PVC) or image still pulling |
| `ContainerCreating` | scheduled; pulling image, mounting volumes, setting up network |
| `Running` but `0/1 READY` | started, readiness probe failing |
| `CrashLoopBackOff` | container keeps exiting; kubelet waits 10s, 20s, 40s ... up to 5 min between restarts |
| `ImagePullBackOff` / `ErrImagePull` | wrong image name/tag, private registry without credentials, arch mismatch |
| `CreateContainerConfigError` | referenced ConfigMap/Secret/key missing |
| `OOMKilled` (last state) | exceeded memory limit |
| `Terminating` (stuck) | finalizers, or node unreachable |
| `Evicted` | node under pressure (memory/disk) evicted it |

Conditions (`kubectl get pod x -o yaml | yq .status.conditions`):
`PodScheduled`, `Initialized`, `ContainersReady`, `Ready`. **Only `Ready`
pods receive Service traffic.**

### 13.3 Init containers

Run to completion, in order, **before** app containers start. Use them for
one-time setup, not for "wait for the DB" loops: your app should retry
connections itself (§11).

```yaml
spec:
  initContainers:
    - name: fetch-config
      image: busybox:1.36
      command: ["sh", "-c", "wget -qO /work/feature-flags.json http://flags.internal/orbit"]
      volumeMounts: [{ name: work, mountPath: /work }]
  containers:
    - name: app
      image: orbit-orders:dev
      volumeMounts: [{ name: work, mountPath: /etc/orbit }]
  volumes:
    - name: work
      emptyDir: {}
```

### 13.4 Sidecars: native sidecar containers (1.29+, GA 1.33)

A sidecar is a helper that runs **alongside** the app for the pod's
lifetime: a log shipper, a proxy, a config reloader. Kubernetes now models
it explicitly: an **init container with `restartPolicy: Always`**. It
starts before the app, keeps running, is restarted if it crashes, and is
**stopped after the app** on shutdown. That fixes two long-standing bugs:
Jobs that never complete because the sidecar never exits, and proxies that
die before the app finishes draining.

```yaml
spec:
  initContainers:
    - name: cloud-sql-proxy          # starts first, stays up, stops last
      image: gcr.io/cloud-sql-connectors/cloud-sql-proxy:2.14.0
      restartPolicy: Always
      args: ["--port=5432", "acme:us-central1:orbit-db"]
      startupProbe:
        tcpSocket: { port: 5432 }
  containers:
    - name: app
      image: orbit-orders:dev
      env: [{ name: DATABASE_URL, value: "postgres://orbit@127.0.0.1:5432/orbit" }]
```

Common patterns: **sidecar** (extend the app), **ambassador** (proxy
outbound connections, like the SQL proxy above), **adapter** (normalize
output, e.g. a metrics exporter). From *Designing Distributed Systems*
(Burns) and *Kubernetes Patterns* (Ibryam & Huß).

### 13.5 The termination sequence (memorize this)

When a pod is deleted (rollout, scale-down, eviction, node drain):

```
t=0   Pod marked Terminating. Two things happen IN PARALLEL:
      (a) endpoint removal: EndpointSlice controller removes the pod IP ->
          kube-proxy/ingress/LB update over the next ~1-5+ seconds
      (b) on the kubelet: run preStop hook (if any), THEN send SIGTERM to PID 1
t=?   app must: stop accepting new work, finish in-flight work, close connections, exit 0
t=30s (terminationGracePeriodSeconds, default 30) -> SIGKILL, no matter what
```

The race between (a) and (b) is why rolling deploys drop requests: the app
receives SIGTERM and stops listening **before** every load balancer has
stopped sending it traffic. The fix is to **keep serving for a few seconds
after SIGTERM**:

```yaml
spec:
  terminationGracePeriodSeconds: 45     # > preStop + drain + shutdown timeout
  containers:
    - name: app
      lifecycle:
        preStop:
          sleep: { seconds: 5 }         # native sleep action (1.30+); no shell needed (works on distroless)
```

Our Go server *also* sleeps `DRAIN_DELAY` after SIGTERM while failing
readiness (§4.1). Use one or the other, or both (preStop runs first). The
budget must fit: `preStop (5s) + DRAIN_DELAY (5s) + SHUTDOWN_TIMEOUT (20s) <
terminationGracePeriodSeconds (45s)`.

> **War story:** a team saw a burst of 502s on *every* deploy, about 0.3% of
> requests, for months. The pods shut down cleanly in under a second. That
> was the problem: they shut down before the cloud load balancer
> deregistered them, which took up to 10 s. A 15 s `preStop` sleep made the
> errors disappear. Cloud LBs (AWS ALB in particular) deregister *slowly*;
> size the sleep to your LB.

> ⚠️ **Caveat (verify):** feature status by version: native sidecars (beta and on by default in 1.29, GA in 1.33) and the `preStop.sleep` action (beta and on by default since 1.30). Confirm on your cluster with `kubectl explain pod.spec.containers.lifecycle.preStop.sleep` and `kubectl explain pod.spec.initContainers.restartPolicy`. If `explain` doesn't know the field, the API server doesn't either.

---

## 14. Deployments, ReplicaSets, and rollouts

### 14.1 The Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: catalog
  labels: { app.kubernetes.io/name: catalog, app.kubernetes.io/part-of: orbit }
spec:
  replicas: 3
  revisionHistoryLimit: 5
  selector:
    matchLabels: { app.kubernetes.io/name: catalog }     # IMMUTABLE after creation
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 25%          # extra pods allowed above replicas during rollout
      maxUnavailable: 0      # never go below desired capacity
  minReadySeconds: 10        # a new pod must stay Ready 10s before it "counts"
  progressDeadlineSeconds: 300
  template:
    metadata:
      labels: { app.kubernetes.io/name: catalog, app.kubernetes.io/part-of: orbit }
    spec:
      containers:
        - name: catalog
          image: orbit-catalog:dev
          ports: [{ name: http, containerPort: 8080 }]
          readinessProbe: { httpGet: { path: /readyz, port: http }, periodSeconds: 5 }
```

`Deployment` → manages `ReplicaSet`s (one per pod-template version) →
manages `Pod`s. A rolling update scales the new ReplicaSet up and the old
one down within the `maxSurge`/`maxUnavailable` budget, **gated on
readiness**. That is why readiness probes are what make rollouts safe.

| Setting | Effect |
|---|---|
| `maxUnavailable: 0, maxSurge: 1` | slowest, safest; needs spare capacity for 1 extra pod |
| `maxUnavailable: 25%, maxSurge: 25%` | default; faster, temporarily reduced capacity |
| `strategy: Recreate` | kill all, then start new; for apps that can't run two versions at once (some DB-migration-coupled apps, single-writer apps) |
| `minReadySeconds` | catches "starts fine, crashes 5 s later" before the rollout proceeds |
| `progressDeadlineSeconds` | marks the rollout `Progressing=False` so CI/CD can detect a stuck deploy |

### 14.2 Rolling out and rolling back

```bash
kubectl set image deploy/catalog catalog=orbit-catalog:v2   # (demo only; GitOps in real life)
kubectl rollout status deploy/catalog --timeout=5m           # CI uses this exit code
kubectl rollout history deploy/catalog
kubectl rollout undo deploy/catalog                         # back to previous ReplicaSet
kubectl rollout undo deploy/catalog --to-revision=3
kubectl rollout restart deploy/catalog                      # new pods, same spec (e.g. reload a Secret)
kubectl rollout pause deploy/catalog                        # batch several changes into one rollout
```

Try a broken rollout to see the safety net:

```bash
kubectl set image deploy/catalog catalog=orbit-catalog:does-not-exist
kubectl get pods -l app.kubernetes.io/name=catalog   # old pods still serving; new one ImagePullBackOff
kubectl rollout status deploy/catalog                # eventually: "exceeded its progress deadline"
kubectl rollout undo deploy/catalog
```

With `maxUnavailable: 0`, a bad image **never** reduces capacity. The
rollout stalls with all old pods still serving.

> **Best practice:** `kubectl rollout undo` is a break-glass tool. In a
> GitOps setup the next sync reverts it. The durable rollback is
> **reverting the Git commit** (§36).

### 14.3 Why `selector` and labels matter

The selector is how a ReplicaSet finds "its" pods. If two Deployments'
selectors overlap, they fight over pods. Use the recommended labels
consistently:

```yaml
labels:
  app.kubernetes.io/name: orders          # what it is
  app.kubernetes.io/instance: orders-prod # which instance
  app.kubernetes.io/version: "1.4.2"      # NOT in the selector (it changes every release)
  app.kubernetes.io/component: api
  app.kubernetes.io/part-of: orbit
  app.kubernetes.io/managed-by: kustomize
```

---

## 15. Services, DNS, and the networking model

### 15.1 The Kubernetes network model (the rules every CNI implements)

1. Every pod gets its own IP.
2. Every pod can reach every other pod **without NAT**, across nodes.
3. Agents on a node can reach all pods on that node.

The **CNI plugin** implements this. Common ones: **Cilium** (eBPF, network
policy, observability; the default choice for new clusters), **Calico**,
cloud-native CNIs like the **AWS VPC CNI** (pods get real VPC IPs), and
kind's simple `kindnet`.

Pod IPs are ephemeral. **Services** give a stable virtual IP + DNS name in
front of a changing set of pods selected by label.

### 15.2 Service types

```yaml
apiVersion: v1
kind: Service
metadata:
  name: catalog
spec:
  type: ClusterIP                 # default
  selector: { app.kubernetes.io/name: catalog }
  ports:
    - name: http
      port: 80                    # the Service's port
      targetPort: http            # the container port NAME (rename-safe)
```

| Type | Reachable from | Use |
|---|---|---|
| `ClusterIP` | inside the cluster | service-to-service (99% of Services) |
| headless (`clusterIP: None`) | DNS returns pod IPs directly | StatefulSets, client-side load balancing, gRPC |
| `NodePort` | `<any-node-ip>:30000-32767` | rarely directly; building block for LBs |
| `LoadBalancer` | external cloud LB | exposing a gateway/ingress controller (one LB, not one per app) |
| `ExternalName` | DNS CNAME | aliasing an external DB host |

### 15.3 DNS

CoreDNS gives every Service a name:

```
catalog                                   # same namespace
catalog.orbit                             # other namespace
catalog.orbit.svc.cluster.local           # fully qualified
redis-0.redis.orbit.svc.cluster.local     # one StatefulSet pod via a headless Service
```

```bash
kubectl run -it --rm dnsutils --image=registry.k8s.io/e2e-test-images/agnhost:2.39 --restart=Never -- \
  nslookup catalog.orbit
```

> **War story:** pod `/etc/resolv.conf` has `ndots:5`, so the name
> `api.stripe.com` (2 dots) is first tried as
> `api.stripe.com.orbit.svc.cluster.local`, `...svc.cluster.local`,
> `...cluster.local`, plus the cloud's search domains, before the real
> name. That's 4–6 wasted queries per lookup (×2 for A and AAAA). At scale
> this overloaded CoreDNS and added latency to every external call. Fixes:
> use FQDNs with a trailing dot (`api.stripe.com.`), set
> `dnsConfig.options: [{name: ndots, value: "2"}]` on the pod, and run
> **NodeLocal DNSCache**.

### 15.4 How a Service actually routes

The **EndpointSlice** controller keeps a list of *ready* pod IPs for each
Service. **kube-proxy** (iptables/IPVS/nftables) or an **eBPF dataplane**
(Cilium) on every node rewrites packets addressed to the Service VIP to one
of those pod IPs. There is no proxy process in the data path. It is
connection-level load balancing.

Consequence: **long-lived connections don't rebalance**. HTTP/2 and gRPC
clients open one connection and pin to one pod. Scale from 3 to 10 pods
and the new 7 get no traffic from existing clients. Fixes: client-side
balancing with a headless Service, a service mesh / L7 proxy, or a max
connection age on the server (`MaxConnectionAge` in gRPC).

```bash
kubectl get endpointslices -l kubernetes.io/service-name=catalog -o wide
# If this is empty, your Service selects nothing or no pod is Ready (see §33).
```

---

## 16. Configuration: ConfigMaps, Secrets, and the Downward API

### 16.1 ConfigMaps

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: orders-config
data:
  LOG_LEVEL: info
  DB_MAX_CONNS: "10"
  DRAIN_DELAY: 5s
  features.yaml: |
    new_checkout: false
```

```yaml
# consume as env vars (read at start only) ...
envFrom:
  - configMapRef: { name: orders-config }
# ... or as files (updated in place ~1 min after the ConfigMap changes, unless mounted via subPath)
volumeMounts:
  - { name: config, mountPath: /etc/orbit, readOnly: true }
volumes:
  - name: config
    configMap: { name: orders-config, items: [{ key: features.yaml, path: features.yaml }] }
```

> **Best practice:** changing a ConfigMap does **not** restart pods. Make
> config changes roll out like code changes. Kustomize's
> `configMapGenerator` (and Helm's checksum-annotation pattern) put a hash
> in the name or annotation, so a config change produces a new pod
> template and a normal rolling update, with rollback (§26).

### 16.2 Secrets: what they are and are not

```bash
kubectl create secret generic orders-db --from-literal=DATABASE_URL='postgres://...' --dry-run=client -o yaml
```

A Secret is a ConfigMap with **base64 encoding (not encryption)**, tighter
default handling (tmpfs mounts, not printed by `describe`), and separate
RBAC. To make Secrets actually secret:

1. **Encryption at rest** in etcd (managed clusters: enable envelope
   encryption with your KMS key).
2. **RBAC**: `get/list secrets` is effectively read access to every
   credential in the namespace. `list` returns values too.
3. **Never commit Secret YAML to Git.** Use the External Secrets Operator
   (pulls from Vault/AWS Secrets Manager/GCP Secret Manager) or Sealed
   Secrets / SOPS (encrypted in Git). Covered in §34.4.
4. Prefer **files over env vars** for secrets. Env vars leak into crash
   dumps, `/proc/<pid>/environ`, and child processes, and get logged by
   "print config at startup" code.

### 16.3 Image pull secrets

```bash
kubectl create secret docker-registry ghcr --docker-server=ghcr.io \
  --docker-username=$USER --docker-password=$GHCR_TOKEN -n orbit
kubectl patch serviceaccount default -n orbit -p '{"imagePullSecrets":[{"name":"ghcr"}]}'
```

On EKS/GKE/AKS, prefer node or workload identity to the cloud registry: no
static credentials to rotate.

### 16.4 The Downward API: tell the pod about itself

```yaml
env:
  - name: POD_NAME
    valueFrom: { fieldRef: { fieldPath: metadata.name } }       # worker consumer name
  - name: POD_NAMESPACE
    valueFrom: { fieldRef: { fieldPath: metadata.namespace } }
  - name: NODE_NAME
    valueFrom: { fieldRef: { fieldPath: spec.nodeName } }
  - name: POD_IP
    valueFrom: { fieldRef: { fieldPath: status.podIP } }
  - name: MEMORY_LIMIT
    valueFrom: { resourceFieldRef: { resource: limits.memory } }
  - name: OTEL_RESOURCE_ATTRIBUTES          # telemetry knows which pod/node it came from
    value: "k8s.pod.name=$(POD_NAME),k8s.namespace.name=$(POD_NAMESPACE),k8s.node.name=$(NODE_NAME)"
```

`$(VAR)` references only work for variables defined **earlier** in the
same `env` list.

---

## 17. Storage: Volumes, PVs, PVCs, StorageClasses

### 17.1 The abstraction stack

```
Pod ──uses──► PersistentVolumeClaim ("I need 10Gi RWO, class fast-ssd")
                     │ bound to
                     ▼
              PersistentVolume (an actual disk: EBS vol-123, GCE PD, local path)
                     ▲ dynamically created by
              StorageClass (provisioner = CSI driver + parameters)
```

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast-ssd
provisioner: ebs.csi.aws.com            # kind uses rancher.io/local-path
parameters: { type: gp3, encrypted: "true" }
reclaimPolicy: Retain                   # keep the disk when the PVC is deleted (prod data!)
allowVolumeExpansion: true
volumeBindingMode: WaitForFirstConsumer # create the disk in the zone where the pod is scheduled
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata: { name: data }
spec:
  accessModes: [ReadWriteOnce]
  storageClassName: fast-ssd
  resources: { requests: { storage: 10Gi } }
```

| Access mode | Meaning | Typical backend |
|---|---|---|
| `ReadWriteOnce` (RWO) | one **node** mounts read-write | block storage: EBS, PD, Azure Disk |
| `ReadWriteOncePod` | one **pod** | same; stricter single-writer guarantee |
| `ReadOnlyMany` (ROX) | many nodes read-only | |
| `ReadWriteMany` (RWX) | many nodes read-write | NFS, EFS, Filestore, CephFS (slower, use sparingly) |

> **War story:** `volumeBindingMode: Immediate` (the old default in many
> classes) created a pod's EBS volume in `us-east-1a` before the pod was
> scheduled. The scheduler then put the pod in `1b`, where the volume
> couldn't attach, and the pod stayed `Pending` forever with
> "volume node affinity conflict". **Always use `WaitForFirstConsumer` for
> zonal disks.**

### 17.2 Ephemeral volumes

- `emptyDir`: scratch space for the pod's lifetime; `medium: Memory` makes
  it tmpfs (counts against the memory limit). Pair it with
  `readOnlyRootFilesystem` for `/tmp`.
- `configMap`, `secret`, `downwardAPI`, `projected`: config as files.
- **Never `hostPath`** in application workloads. It is a container-escape
  primitive and ties the pod to a node.

---

## 18. StatefulSets, DaemonSets, Jobs, and CronJobs

### 18.1 StatefulSet: stable identity + stable storage

Use a StatefulSet when each replica needs **its own** disk and a **stable
name**: databases, Kafka, ZooKeeper, Elasticsearch, Redis.

| | Deployment | StatefulSet |
|---|---|---|
| Pod names | random (`orders-7d9f-x2k`) | ordinal (`redis-0`, `redis-1`) |
| Storage | shared or none | `volumeClaimTemplates`: one PVC per pod, kept across reschedules |
| Network identity | via Service only | per-pod DNS through a headless Service |
| Order | parallel | ordered create/delete (default `OrderedReady`) |

Orbit's Redis (single instance, AOF persistence; Redis HA is a separate
topic: use a managed service or an operator in production):

```yaml
apiVersion: v1
kind: Service
metadata: { name: redis, labels: { app.kubernetes.io/name: redis } }
spec:
  clusterIP: None                 # headless: required for StatefulSet pod DNS
  selector: { app.kubernetes.io/name: redis }
  ports: [{ name: redis, port: 6379 }]
---
apiVersion: apps/v1
kind: StatefulSet
metadata: { name: redis }
spec:
  serviceName: redis
  replicas: 1
  selector: { matchLabels: { app.kubernetes.io/name: redis } }
  template:
    metadata: { labels: { app.kubernetes.io/name: redis } }
    spec:
      securityContext: { runAsUser: 999, runAsGroup: 999, fsGroup: 999, runAsNonRoot: true }
      containers:
        - name: redis
          image: redis:7.4-alpine
          args: ["--appendonly", "yes", "--maxmemory", "200mb", "--maxmemory-policy", "noeviction"]
          ports: [{ name: redis, containerPort: 6379 }]
          resources:
            requests: { cpu: 100m, memory: 256Mi }
            limits:   { memory: 256Mi }
          readinessProbe: { exec: { command: ["redis-cli", "ping"] }, periodSeconds: 5 }
          livenessProbe:  { tcpSocket: { port: redis }, periodSeconds: 10, failureThreshold: 6 }
          volumeMounts: [{ name: data, mountPath: /data }]
  volumeClaimTemplates:
    - metadata: { name: data }
      spec:
        accessModes: [ReadWriteOnce]
        resources: { requests: { storage: 1Gi } }
```

`noeviction` matters because Redis is our *queue* too. With an LRU policy
under memory pressure, Redis would silently drop unprocessed orders.

> **Best practice:** for databases, **use an operator** (CloudNativePG,
> Percona, Strimzi for Kafka) or a **managed service** (RDS, Cloud SQL).
> A StatefulSet gives you identity and disks. It does not give you
> failover, backups, replication, minor-version upgrades, or connection
> pooling. Operators encode a DBA's runbook as a controller. Orbit uses
> CloudNativePG (§20).

### 18.2 DaemonSet: one pod per node

For node agents: log collectors, the OTel Collector agent, CNI, CSI node
plugins, node-exporter, security agents.

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata: { name: node-agent, namespace: monitoring }
spec:
  selector: { matchLabels: { app: node-agent } }
  updateStrategy: { type: RollingUpdate, rollingUpdate: { maxUnavailable: 10% } }
  template:
    metadata: { labels: { app: node-agent } }
    spec:
      tolerations: [{ operator: Exists }]        # also run on tainted nodes (control plane, GPU, spot)
      priorityClassName: system-node-critical     # never evicted for app pods
      containers:
        - name: agent
          image: busybox:1.36
          command: ["sh", "-c", "while true; do date; sleep 60; done"]
          resources: { requests: { cpu: 10m, memory: 16Mi }, limits: { memory: 32Mi } }
```

You'll deploy a real one, the OTel Collector agent, in §30.

### 18.3 Jobs and CronJobs

```yaml
apiVersion: batch/v1
kind: Job
metadata: { name: reprice-catalog }
spec:
  backoffLimit: 3                 # retries before Failed
  activeDeadlineSeconds: 600      # hard timeout for the whole job
  ttlSecondsAfterFinished: 86400  # garbage-collect after a day
  completions: 10                 # 10 work items...
  parallelism: 3                  # ...3 at a time
  completionMode: Indexed         # each pod gets JOB_COMPLETION_INDEX 0..9 (shard the work)
  podFailurePolicy:               # don't retry bugs; do retry infra disruptions
    rules:
      - action: FailJob
        onExitCodes: { containerName: job, operator: In, values: [42] }
      - action: Ignore
        onPodConditions: [{ type: DisruptionTarget }]
  template:
    spec:
      restartPolicy: Never
      containers:
        - name: job
          image: orbit-catalog:dev
          command: ["/app", "reprice", "--shard=$(JOB_COMPLETION_INDEX)"]
---
apiVersion: batch/v1
kind: CronJob
metadata: { name: expire-pending-orders }
spec:
  schedule: "*/15 * * * *"
  timeZone: "Etc/UTC"
  concurrencyPolicy: Forbid       # don't start a run while the last is still going
  startingDeadlineSeconds: 300
  successfulJobsHistoryLimit: 3
  failedJobsHistoryLimit: 5
  jobTemplate:
    spec:
      backoffLimit: 2
      activeDeadlineSeconds: 300
      template:
        spec:
          restartPolicy: OnFailure
          containers:
            - name: expire
              image: postgres:17
              envFrom: [{ secretRef: { name: orbit-db-app } }]
              command: ["sh", "-c", "psql \"$uri\" -c \"UPDATE orders SET status='failed' WHERE status='pending' AND created_at < now() - interval '1 hour'\""]
```

> **Best practice:** CronJob runs are **at-least-once and may overlap or
> be skipped** (controller downtime, `startingDeadlineSeconds`). Make
> every job idempotent, set `concurrencyPolicy: Forbid` unless overlap is
> safe, and **alert on "last successful run too old"**, not just on
> failure. A job that silently stops being scheduled never fails.

---

## 19. Namespaces, labels, and organizing a cluster

### 19.1 Namespaces

Namespaces are the unit of **RBAC, quota, network policy, and naming**.
They are not a security boundary by themselves.

```
kube-system           # cluster components
monitoring            # Prometheus, Grafana, Loki, Tempo, OTel
argocd, argo-rollouts, keda, kyverno, cert-manager, envoy-gateway-system, cnpg-system, velero
orbit-dev / orbit-staging / orbit-prod   # (on one cluster for the lab; separate clusters for real prod)
```

> **Best practice:** **prod gets its own cluster** (and usually its own
> cloud account). Namespaces separate teams and apps within an
> environment; clusters separate environments and blast radius. "One
> namespace per env on a shared cluster" is fine for dev/staging, and
> it is how a bad CRD upgrade in staging takes down prod.

### 19.2 Labels vs annotations

- **Labels**: identifying, *selectable* (`-l`, Service selectors, policies,
  cost allocation). Keep them small and stable: `app.kubernetes.io/*`,
  `team`, `env`, `cost-center`.
- **Annotations**: non-identifying metadata for tools: `prometheus.io/*`,
  `argocd.argoproj.io/sync-wave`, change-cause, links to runbooks.

---

## 20. Part II capstone: Orbit from raw manifests

### 20.1 Postgres with CloudNativePG

```bash
helm repo add cnpg https://cloudnative-pg.github.io/charts
helm upgrade --install cnpg cnpg/cloudnative-pg -n cnpg-system --create-namespace --wait
kubectl create namespace orbit
```

```yaml
# deploy/base/postgres/cluster.yaml
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata: { name: orbit-db, namespace: orbit }
spec:
  instances: 3                    # 1 primary + 2 streaming replicas, automatic failover
  imageName: ghcr.io/cloudnative-pg/postgresql:17
  storage: { size: 2Gi }
  bootstrap:
    initdb: { database: orbit, owner: orbit }
  postgresql:
    parameters:
      max_connections: "200"
  affinity:
    topologyKey: topology.kubernetes.io/zone   # spread instances across zones
  resources:
    requests: { cpu: 250m, memory: 512Mi }
    limits:   { memory: 512Mi }
  monitoring:
    enablePodMonitor: true        # Prometheus scrapes it once kube-prometheus-stack exists (§30)
```

CloudNativePG creates:

- Services `orbit-db-rw` (primary), `orbit-db-ro` (replicas), and
  `orbit-db-r` (any instance).
- Secret `orbit-db-app` with `username`, `password`, `uri`, `host`, `port`
  and more.
- Automated failover: delete the primary pod and watch a replica get
  promoted in seconds.

```bash
kubectl apply -f deploy/base/postgres/cluster.yaml
kubectl get cluster -n orbit -w               # wait for "Cluster in healthy state"
kubectl get pods -n orbit -L role             # which one is primary
```

### 20.2 Migrations as a Job

```yaml
apiVersion: batch/v1
kind: Job
metadata: { name: migrate-0001, namespace: orbit }
spec:
  backoffLimit: 5
  ttlSecondsAfterFinished: 3600
  template:
    spec:
      restartPolicy: OnFailure
      containers:
        - name: goose
          image: ghcr.io/acme/orbit-migrations:sha-abc123   # FROM goose + COPY migrations/
          args: ["-dir", "/migrations", "postgres", "$(DATABASE_URL)", "up"]
          env:
            - name: DATABASE_URL
              valueFrom: { secretKeyRef: { name: orbit-db-app, key: uri } }
```

§38 covers how migrations are ordered relative to deploys in GitOps (sync
waves / pre-sync hooks) and the expand/contract pattern.

### 20.3 The application manifests

For each service: a Deployment and a Service. Here is `orders`. This is
the *minimal* version; Part III turns it into the production "golden
manifest".

```yaml
apiVersion: apps/v1
kind: Deployment
metadata: { name: orders, namespace: orbit, labels: { app.kubernetes.io/name: orders } }
spec:
  replicas: 2
  selector: { matchLabels: { app.kubernetes.io/name: orders } }
  template:
    metadata: { labels: { app.kubernetes.io/name: orders, app.kubernetes.io/part-of: orbit } }
    spec:
      containers:
        - name: orders
          image: orbit-orders:dev
          imagePullPolicy: IfNotPresent      # kind-loaded image
          ports: [{ name: http, containerPort: 8080 }]
          env:
            - name: DATABASE_URL
              valueFrom: { secretKeyRef: { name: orbit-db-app, key: uri } }
            - { name: REDIS_ADDR, value: "redis:6379" }
          readinessProbe: { httpGet: { path: /readyz, port: http } }
          livenessProbe:  { httpGet: { path: /livez,  port: http } }
---
apiVersion: v1
kind: Service
metadata: { name: orders, namespace: orbit }
spec:
  selector: { app.kubernetes.io/name: orders }
  ports: [{ name: http, port: 80, targetPort: http }]
```

```bash
kubectl apply -n orbit -f deploy/raw/
kubectl get pods -n orbit
kubectl port-forward -n orbit svc/orders 8082:80 &
curl -s -XPOST localhost:8082/orders -H 'Idempotency-Key: k8s-1' -H 'Content-Type: application/json' \
     -d '{"product_id":1,"quantity":1}'
```

Exercises (do them; they build intuition for Part V):

1. `kubectl delete pod -l app.kubernetes.io/name=orders`. Watch
   replacements appear. Did any request fail during it? (Probably yes;
   Part III fixes that.)
2. Delete the CNPG primary pod. Watch failover with
   `kubectl get pods -L role -w`. Do `orders` pods go not-ready, then ready
   again? Do they restart? (They shouldn't: liveness doesn't check the DB.)
3. Scale `worker` to 0, create 50 orders, scale back to 3. The backlog
   drains. This is the behavior KEDA automates in §27.
4. `kubectl drain <node> --ignore-daemonsets --delete-emptydir-data`, then
   `kubectl uncordon`. What happened to pods on it?

**Part II checkpoint:** you can read any manifest, explain which controller
acts on it, and trace a request from Service → EndpointSlice → pod.

---

# Part III — Production Workloads

## 21. Health checks: probes done right

Probes are how Kubernetes decides to **restart** a container and whether to
**send it traffic**. Get them wrong and Kubernetes' self-healing becomes
self-harm.

### 21.1 The three probes

| Probe | Question | On failure | Runs |
|---|---|---|---|
| **startupProbe** | "Has the app finished starting?" | keep waiting; after `failureThreshold` → restart | until first success, then never again |
| **livenessProbe** | "Is the process wedged beyond self-recovery?" | **restart the container** | after startup succeeds, forever |
| **readinessProbe** | "Can I handle traffic *right now*?" | remove from Service endpoints (no restart) | after startup succeeds, forever |

Mechanisms: `httpGet` (2xx/3xx = pass), `tcpSocket`, `grpc` (the standard
gRPC health protocol), `exec` (command exit code; costly, forks a process
every period).

### 21.2 Orbit's probes

```yaml
startupProbe:
  httpGet: { path: /livez, port: http }
  periodSeconds: 2
  failureThreshold: 30          # up to 60s to start; liveness is held off until then
livenessProbe:
  httpGet: { path: /livez, port: http }
  periodSeconds: 10
  timeoutSeconds: 2
  failureThreshold: 3           # 30s of consecutive failure before a restart
readinessProbe:
  httpGet: { path: /readyz, port: http }
  periodSeconds: 5
  timeoutSeconds: 2
  failureThreshold: 2           # out of rotation within ~10s
  successThreshold: 1
```

**Time to action = `periodSeconds × failureThreshold`** (+ timeouts). Pick
numbers by asking "how long can this be broken before acting is better than
waiting?"

### 21.3 Rules (each one is a past outage)

1. **Liveness must never check dependencies.** If `/livez` pings the
   database, a 30-second DB failover restarts *every pod of every service*
   at once. They all reconnect to the recovering DB simultaneously
   (thundering herd), it falls over again, and you have a self-inflicted
   outage loop.
2. **Liveness and readiness are different endpoints** with different
   semantics. Liveness: "the event loop/HTTP server responds". Readiness:
   "I can do useful work".
3. **Readiness and dependencies: be deliberate.** If *every* replica's
   readiness depends on a shared dependency, a blip of that dependency
   removes *all* endpoints. Clients then get connection errors instead of
   a fast, meaningful 503. Common compromise: readiness checks only
   **local** conditions plus **hard** dependencies without which every
   request fails. Use timeouts, circuit breakers and fallbacks (serve
   stale cache) for the rest. Orbit's `orders` includes Postgres (no
   orders without it). `catalog` would not include Redis, since a cache
   miss is just slower.
4. **Use a startupProbe for slow starters** (JVM warm-up, cache preload,
   ML model load) instead of a huge `initialDelaySeconds` on liveness.
5. **Probes must be cheap and fast.** No DB queries that scan tables, no
   synchronous calls to other services. A probe that times out under load
   restarts healthy-but-busy pods exactly when you need them most.
6. **Timeouts must tolerate GC pauses and CPU throttling.** A 1 s timeout
   on a CPU-throttled pod (§23.3) is a restart generator.
7. **Readiness goes false during shutdown** (our `SetDraining()`), so
   endpoints are removed even if they haven't been yet.

> **War story:** a payments company added a "deep health check" to
> liveness that verified Kafka, Postgres and two downstream APIs. During a
> partial Kafka incident, liveness failed fleet-wide and Kubernetes
> restarted ~400 pods in two minutes. Restart storms saturated the image
> registry and the database connection limits. The Kafka incident lasted 3
> minutes; the self-inflicted outage lasted 40. Liveness went back to
> "return 200".

### 21.4 gRPC services

```yaml
livenessProbe:
  grpc: { port: 9090, service: liveness }     # implement grpc.health.v1.Health
readinessProbe:
  grpc: { port: 9090, service: readiness }
```

---

## 22. Self-healing and auto-recovery

### 22.1 The layers of recovery

| Failure | Who recovers it | How fast | What you must configure |
|---|---|---|---|
| process crash / panic | kubelet restarts the container | seconds (then exponential backoff) | exit non-zero on fatal errors; don't swallow panics |
| process hung / deadlocked | kubelet via **liveness** | `period × threshold` | a liveness probe that can actually detect it |
| not ready (dependency, overload) | EndpointSlice removes it | `period × threshold` | readiness probe |
| pod deleted / evicted | ReplicaSet creates a replacement | seconds | use controllers, never bare pods |
| node dies | node lifecycle controller taints it `unreachable`; pods evicted after `tolerationSeconds` | ~40 s detect + **300 s default** | lower `tolerationSeconds` for stateless apps; replicas ≥ 2 on different nodes |
| zone outage | surviving replicas in other zones | immediate if spread | topology spread across zones (§24) + capacity headroom |
| bad deploy | rollout halts (readiness + `maxUnavailable: 0`), canary analysis rolls back (§37) | minutes | probes, `progressDeadlineSeconds`, Argo Rollouts |
| voluntary disruption (drain, upgrade, autoscaler) | **PodDisruptionBudget** limits concurrent evictions | n/a | a PDB per service |
| transient downstream errors | **your code**: timeouts, retries with backoff + jitter, circuit breaker | ms | client-side resilience |

### 22.2 Faster node-failure recovery for stateless pods

By default every pod tolerates `node.kubernetes.io/unreachable` and
`not-ready` for 300 s, so it waits 5 minutes before being replaced. For
stateless services with ≥ 2 replicas that's far too long:

```yaml
tolerations:
  - { key: node.kubernetes.io/unreachable, operator: Exists, effect: NoExecute, tolerationSeconds: 30 }
  - { key: node.kubernetes.io/not-ready,   operator: Exists, effect: NoExecute, tolerationSeconds: 30 }
```

Don't do this for StatefulSets with RWO volumes. A "dead" node might just
be partitioned and still writing to the disk, and that's how you get split
brain. Let the operator decide.

### 22.3 PodDisruptionBudgets

Voluntary disruptions (node drains for upgrades, Karpenter consolidation,
cluster autoscaler scale-down) go through the **Eviction API**, which
honors PDBs:

```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata: { name: orders }
spec:
  maxUnavailable: 1                                    # or minAvailable: 2 / "80%"
  selector: { matchLabels: { app.kubernetes.io/name: orders } }
  unhealthyPodEvictionPolicy: AlwaysAllow              # don't let already-broken pods block node drains
```

> **War story:** `minAvailable: 100%` (or `maxUnavailable: 0`) on a
> 1-replica Deployment means a node can **never** be drained. A cluster
> upgrade hangs for hours until someone deletes the PDB by hand. PDBs must
> allow at least one eviction, and single-replica workloads need either 2
> replicas or no PDB.

### 22.4 Retries, timeouts, and circuit breakers in the app

Kubernetes restarts things; it doesn't make individual requests succeed.
Every outbound call needs:

```go
// A reusable outbound HTTP client: always set timeouts; the zero-value http.Client has NONE.
var httpClient = &http.Client{
	Timeout: 3 * time.Second,
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 100,          // default 2 causes connection churn under load
		IdleConnTimeout:     90 * time.Second,
	},
}

// Retry only idempotent operations, with capped exponential backoff + full jitter,
// and stop when the caller's context deadline is exceeded.
func retry(ctx context.Context, attempts int, fn func(context.Context) error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = fn(ctx); err == nil || !isRetryable(err) {
			return err
		}
		backoff := time.Duration(rand.Int64N(int64(100*time.Millisecond) << i)) // full jitter
		select {
		case <-time.After(min(backoff, 2*time.Second)):
		case <-ctx.Done():
			return errors.Join(err, ctx.Err())
		}
	}
	return err
}
```

- **Timeouts propagate**: each hop's timeout must be smaller than its
  caller's. Use `context.Context` deadlines.
- **Retry budgets**: retries at every layer multiply. With 3 layers × 3
  retries, one user request can become 27 at the bottom. That is a retry
  storm. Retry at one layer, or cap retries as a percentage of traffic.
- **Circuit breakers** (e.g. `sony/gobreaker`) fail fast when a dependency
  is down, so you don't hold goroutines and connections waiting for it.
- **Load shedding**: return 503 quickly under overload instead of queueing
  forever. A fast "no" beats a slow timeout.

### 22.5 Priority and preemption

```yaml
apiVersion: scheduling.k8s.io/v1
kind: PriorityClass
metadata: { name: orbit-critical }
value: 100000
description: "User-facing request path (catalog, orders)"
---
apiVersion: scheduling.k8s.io/v1
kind: PriorityClass
metadata: { name: orbit-batch }
value: 1000
preemptionPolicy: Never          # waits for room instead of evicting others
```

When the cluster is full, higher-priority pods **preempt** lower-priority
ones. Put user-facing APIs above batch and workers, so a capacity crunch
degrades background work first.

---

## 23. Resources, QoS, and right-sizing

### 23.1 Requests vs limits

| | requests | limits |
|---|---|---|
| Used by | **the scheduler** (bin-packing) and HPA percentages | **the kernel** (cgroups) at runtime |
| CPU | guaranteed share (cgroup weight) | hard quota: **throttled** above it |
| Memory | used for scheduling and eviction ranking | hard cap: **OOMKilled** above it |

```yaml
resources:
  requests: { cpu: 250m, memory: 256Mi }    # 250 millicores = 1/4 core
  limits:   { memory: 256Mi }               # memory limit = request; NO cpu limit (see below)
```

### 23.2 QoS classes (who gets evicted first under node memory pressure)

| Class | Condition | Eviction order |
|---|---|---|
| `Guaranteed` | every container: requests == limits for CPU **and** memory | last |
| `Burstable` | some request set, not Guaranteed | middle (by usage above request) |
| `BestEffort` | no requests or limits at all | first |

### 23.3 The CPU-limit debate

A CPU limit is enforced by the CFS quota in 100 ms periods. A Go service
with 8 busy goroutines and `limits.cpu: 1` can burn its entire 100 ms quota
in 12.5 ms of wall time, then **freeze for 87.5 ms**. p99 latency spikes
appear even though average CPU usage looks low.

```promql
# fraction of CFS periods in which the container was throttled
sum by (pod) (rate(container_cpu_cfs_throttled_periods_total{namespace="orbit"}[5m]))
  / sum by (pod) (rate(container_cpu_cfs_periods_total{namespace="orbit"}[5m]))
```

Widely adopted guidance for latency-sensitive services:

- **Always set CPU requests** (scheduling + fair share under contention).
- **Omit CPU limits** for latency-sensitive services on clusters you
  control, or set them generously (2–4× request). Unused CPU on the node is
  then usable, and contention is resolved by request-proportional shares.
- **Always set memory limits, equal to the memory request.** Memory isn't
  compressible: overcommit leads to node-level OOM, which kills random
  pods.
- Multi-tenant clusters with untrusted neighbors may *require* CPU limits.
  Enforce that with LimitRange.

### 23.4 Make the runtime container-aware

```yaml
env:
  # Go 1.25+ sets GOMAXPROCS from the cgroup CPU limit automatically; without a limit, it uses
  # node cores. On older Go, use go.uber.org/automaxprocs.
  - name: GOMEMLIMIT          # soft memory target: GC works harder before the hard limit -> fewer OOMKills
    value: "230MiB"          # ~90% of a 256Mi limit
  # JVM:    -XX:MaxRAMPercentage=75   (container-aware since JDK 10)
  # Node:   --max-old-space-size=200  (in MB; Node does not size the heap from the cgroup on its own)
  # Python: watch out for per-worker memory x worker count (gunicorn --workers)
```

> **War story:** a Node.js service with a 512 Mi limit was OOMKilled every
> few hours. V8 sized its heap from the *node's* 64 GB RAM, so it saw no
> reason to collect garbage until the cgroup killed it. Setting
> `--max-old-space-size=400` fixed it. The JVM had the same problem before
> JDK 10.

> ⚠️ **Caveat (verify):** runtime behavior differs by version. Go 1.25+ derives GOMAXPROCS from the cgroup CPU limit (older Go needs `automaxprocs`). Newer Node.js releases take container memory into account when sizing the heap, so the war story applies mainly to older Node versions. Check yours: `node -e "console.log(v8.getHeapStatistics().heap_size_limit/1048576)"` inside the container.

### 23.5 Guardrails per namespace

```yaml
apiVersion: v1
kind: LimitRange
metadata: { name: defaults, namespace: orbit }
spec:
  limits:
    - type: Container
      defaultRequest: { cpu: 100m, memory: 128Mi }   # applied when a container omits requests
      default:        { memory: 256Mi }               # applied when it omits limits
      max:            { memory: 4Gi }
---
apiVersion: v1
kind: ResourceQuota
metadata: { name: orbit-quota, namespace: orbit }
spec:
  hard:
    requests.cpu: "20"
    requests.memory: 40Gi
    limits.memory: 40Gi
    pods: "200"
    services.loadbalancers: "0"    # force traffic through the shared Gateway
    persistentvolumeclaims: "20"
```

### 23.6 Right-sizing with data

Requests should come from **observed usage**, not guesses:

```promql
# p95 CPU usage per container over 7 days -> a candidate CPU request
quantile_over_time(0.95, (sum by (container) (rate(container_cpu_usage_seconds_total{namespace="orbit",container!=""}[5m])))[7d:5m])
# max working-set memory over 7 days -> memory request/limit plus ~20-30% headroom
max_over_time(max by (container) (container_memory_working_set_bytes{namespace="orbit",container!=""})[7d:])
```

VPA in recommendation mode (§27.5) does this calculation for you.

**In-place pod resize** (beta, on by default since 1.33) changes a running
pod's CPU and memory without a restart:
`kubectl patch pod orders-x --subresource resize -p '{"spec":{"containers":[{"name":"orders","resources":{"requests":{"cpu":"500m"}}}]}}'`.
It's useful for startup CPU boosts and for VPA's `InPlaceOrRecreate` mode.

> ⚠️ **Caveat (verify):** in-place pod resize was beta (on by default) in 1.33 and may be GA in your version; its API details have changed between releases. The `--subresource resize` patch needs a matching kubectl (1.32+). Check with `kubectl explain pod.spec.containers.resizePolicy`. The PromQL in §23.6 uses cAdvisor/kube-state-metrics names; check them in Prometheus before using them.

---

## 24. Scheduling: placement for resilience

### 24.1 Spread replicas across zones and nodes

```yaml
topologySpreadConstraints:
  - maxSkew: 1
    topologyKey: topology.kubernetes.io/zone
    whenUnsatisfiable: ScheduleAnyway        # prefer spread, don't block scheduling
    labelSelector: { matchLabels: { app.kubernetes.io/name: orders } }
    matchLabelKeys: [pod-template-hash]      # compute skew per rollout revision
  - maxSkew: 1
    topologyKey: kubernetes.io/hostname
    whenUnsatisfiable: ScheduleAnyway
    labelSelector: { matchLabels: { app.kubernetes.io/name: orders } }
    matchLabelKeys: [pod-template-hash]
```

Without this, the scheduler can put all 3 replicas on one node. That node
reboots and the service is down even though "replicas: 3".

### 24.2 Affinity, taints, and tolerations

```yaml
# node affinity: only arm64 nodes, prefer spot
affinity:
  nodeAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      nodeSelectorTerms:
        - matchExpressions: [{ key: kubernetes.io/arch, operator: In, values: [arm64] }]
    preferredDuringSchedulingIgnoredDuringExecution:
      - weight: 50
        preference: { matchExpressions: [{ key: karpenter.sh/capacity-type, operator: In, values: [spot] }] }
  # pod affinity: put the worker near Redis (same zone) to cut cross-zone traffic cost and latency
  podAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
      - weight: 100
        podAffinityTerm:
          topologyKey: topology.kubernetes.io/zone
          labelSelector: { matchLabels: { app.kubernetes.io/name: redis } }
```

**Taints repel, tolerations permit.** A dedicated GPU pool:

```bash
kubectl taint nodes gpu-node-1 nvidia.com/gpu=true:NoSchedule
```

```yaml
tolerations: [{ key: nvidia.com/gpu, operator: Exists, effect: NoSchedule }]
nodeSelector: { nvidia.com/gpu.present: "true" }   # toleration alone doesn't ATTRACT, it only allows
```

### 24.3 The golden Deployment (everything so far, in one manifest)

This is the reference manifest for a stateless HTTP service. Every field
is here for a reason covered above. Copy it as your template.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: orders
  labels: &labels
    app.kubernetes.io/name: orders
    app.kubernetes.io/part-of: orbit
    app.kubernetes.io/component: api
spec:
  replicas: 3                               # removed when an HPA/KEDA manages replicas (§27)
  revisionHistoryLimit: 5
  selector: { matchLabels: { app.kubernetes.io/name: orders } }
  strategy: { type: RollingUpdate, rollingUpdate: { maxSurge: 25%, maxUnavailable: 0 } }
  minReadySeconds: 5
  progressDeadlineSeconds: 300
  template:
    metadata:
      labels: *labels
    spec:
      serviceAccountName: orders            # its own identity (§34.1), never "default"
      automountServiceAccountToken: false   # it never calls the K8s API
      priorityClassName: orbit-critical
      terminationGracePeriodSeconds: 45
      securityContext:
        runAsNonRoot: true
        runAsUser: 65532
        runAsGroup: 65532
        fsGroup: 65532
        seccompProfile: { type: RuntimeDefault }
      topologySpreadConstraints:
        - { maxSkew: 1, topologyKey: topology.kubernetes.io/zone, whenUnsatisfiable: ScheduleAnyway,
            labelSelector: { matchLabels: { app.kubernetes.io/name: orders } }, matchLabelKeys: [pod-template-hash] }
        - { maxSkew: 1, topologyKey: kubernetes.io/hostname, whenUnsatisfiable: ScheduleAnyway,
            labelSelector: { matchLabels: { app.kubernetes.io/name: orders } }, matchLabelKeys: [pod-template-hash] }
      tolerations:
        - { key: node.kubernetes.io/unreachable, operator: Exists, effect: NoExecute, tolerationSeconds: 30 }
        - { key: node.kubernetes.io/not-ready,   operator: Exists, effect: NoExecute, tolerationSeconds: 30 }
      containers:
        - name: orders
          image: ghcr.io/acme/orbit-orders:sha-3f9c2a1   # immutable tag (GitOps pins the digest)
          imagePullPolicy: IfNotPresent
          ports: [{ name: http, containerPort: 8080 }]
          envFrom: [{ configMapRef: { name: orders-config } }]
          env:
            - name: DATABASE_URL
              valueFrom: { secretKeyRef: { name: orbit-db-app, key: uri } }
            - { name: GOMEMLIMIT, value: "230MiB" }
            - name: POD_NAME
              valueFrom: { fieldRef: { fieldPath: metadata.name } }
            - { name: OTEL_SERVICE_NAME, value: orders }
            - name: OTEL_EXPORTER_OTLP_ENDPOINT     # node-local collector agent (§30)
              value: "http://$(HOST_IP):4317"
          resources:
            requests: { cpu: 250m, memory: 256Mi }
            limits:   { memory: 256Mi }
          startupProbe:   { httpGet: { path: /livez,  port: http }, periodSeconds: 2,  failureThreshold: 30 }
          livenessProbe:  { httpGet: { path: /livez,  port: http }, periodSeconds: 10, timeoutSeconds: 2, failureThreshold: 3 }
          readinessProbe: { httpGet: { path: /readyz, port: http }, periodSeconds: 5,  timeoutSeconds: 2, failureThreshold: 2 }
          lifecycle:
            preStop: { sleep: { seconds: 5 } }
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            capabilities: { drop: ["ALL"] }
          volumeMounts: [{ name: tmp, mountPath: /tmp }]
      volumes:
        - { name: tmp, emptyDir: { sizeLimit: 64Mi } }
```

(`HOST_IP` comes from `fieldRef: status.hostIP`; it's omitted above for
brevity and included in the Kustomize base in §26.)

---

## 25. Traffic in: Gateway API, TLS, and cert-manager

### 25.1 Ingress → Gateway API

The `Ingress` API is frozen. Its common controller, **ingress-nginx, was
retired by Kubernetes SIG Network in March 2026** (no further releases or
security fixes). Annotations-for-everything also made Ingress unportable.
**Gateway API** is the successor: typed, role-oriented, and supported by
Envoy Gateway, Istio, Cilium, NGINX Gateway Fabric, Traefik, Kong, and the
cloud LBs (GKE, AWS LB Controller, Azure).

```
GatewayClass   (infra provider: "use Envoy Gateway")             owned by: platform team
   └─ Gateway  (listeners: :443 for *.orbit.example, TLS cert)    owned by: platform team
        └─ HTTPRoute (orbit.example/api/orders -> svc orders:80)  owned by: app team, in its namespace
```

> ⚠️ **Caveat (verify):** confirm the ingress-nginx retirement status and dates on the Kubernetes blog (kubernetes.io/blog) before citing them. The direction (move to Gateway API) is well established; the exact dates are what to check.

### 25.2 Install Envoy Gateway and cert-manager

```bash
helm upgrade --install eg oci://docker.io/envoyproxy/gateway-helm \
  -n envoy-gateway-system --create-namespace --wait
helm repo add jetstack https://charts.jetstack.io
helm upgrade --install cert-manager jetstack/cert-manager -n cert-manager --create-namespace \
  --set crds.enabled=true --set config.enableGatewayAPI=true --wait
```

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: GatewayClass
metadata: { name: eg }
spec: { controllerName: gateway.envoyproxy.io/gatewayclass-controller }
---
# Lab: a self-signed CA. Production: ACME/Let's Encrypt (below) or your corporate CA.
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata: { name: selfsigned }
spec: { selfSigned: {} }
---
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: public
  namespace: envoy-gateway-system
  annotations: { cert-manager.io/cluster-issuer: selfsigned }   # cert-manager creates the Secret
spec:
  gatewayClassName: eg
  listeners:
    - name: https
      protocol: HTTPS
      port: 443
      hostname: "orbit.localtest.me"          # *.localtest.me resolves to 127.0.0.1
      tls: { mode: Terminate, certificateRefs: [{ name: orbit-tls }] }
      allowedRoutes: { namespaces: { from: Selector, selector: { matchLabels: { gateway-access: "true" } } } }
    - name: http
      protocol: HTTP
      port: 80
      allowedRoutes: { namespaces: { from: Same } }
---
# HTTP -> HTTPS redirect
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata: { name: https-redirect, namespace: envoy-gateway-system }
spec:
  parentRefs: [{ name: public, sectionName: http }]
  rules: [{ filters: [{ type: RequestRedirect, requestRedirect: { scheme: https, statusCode: 301 } }] }]
```

Production issuer:

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata: { name: letsencrypt }
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: platform@acme.example
    privateKeySecretRef: { name: letsencrypt-account }
    solvers:
      - http01: { gatewayHTTPRoute: { parentRefs: [{ name: public, namespace: envoy-gateway-system, kind: Gateway }] } }
```

> ⚠️ **Caveat (verify):** pin `--version` for both charts. The `config.enableGatewayAPI` value and the Gateway-annotation certificate flow depend on the cert-manager version. Check its *Gateway API* docs page; on older releases it was a feature-gate flag. On kind, the Gateway gets an address only while `cloud-provider-kind` is running.

### 25.3 Routes owned by the app team

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata: { name: orbit-api, namespace: orbit }     # namespace labeled gateway-access=true
spec:
  parentRefs: [{ name: public, namespace: envoy-gateway-system, sectionName: https }]
  hostnames: ["orbit.localtest.me"]
  rules:
    - matches: [{ path: { type: PathPrefix, value: /api/products } }]
      filters:
        - type: URLRewrite
          urlRewrite: { path: { type: ReplacePrefixMatch, replacePrefixMatch: /products } }
      backendRefs: [{ name: catalog, port: 80 }]
    - matches: [{ path: { type: PathPrefix, value: /api/orders }, method: POST }]
      filters:
        - type: URLRewrite
          urlRewrite: { path: { type: ReplacePrefixMatch, replacePrefixMatch: /orders } }
      timeouts: { request: 5s }
      backendRefs: [{ name: orders, port: 80 }]
```

```bash
kubectl label namespace orbit gateway-access=true
kubectl get gateway -n envoy-gateway-system public        # PROGRAMMED=True, ADDRESS from cloud-provider-kind
GW=$(kubectl get gateway -n envoy-gateway-system public -o jsonpath='{.status.addresses[0].value}')
curl -sk --resolve orbit.localtest.me:443:$GW https://orbit.localtest.me/api/products | jq
```

Traffic splitting (the foundation for canaries, §37) is just weights:

```yaml
backendRefs:
  - { name: orders,        port: 80, weight: 90 }
  - { name: orders-canary, port: 80, weight: 10 }
```

### 25.4 Edge resilience policies (Envoy Gateway)

```yaml
apiVersion: gateway.envoyproxy.io/v1alpha1
kind: BackendTrafficPolicy
metadata: { name: orders-resilience, namespace: orbit }
spec:
  targetRefs: [{ group: gateway.networking.k8s.io, kind: HTTPRoute, name: orbit-api }]
  retry:
    numRetries: 2
    retryOn: { triggers: [connect-failure, retriable-status-codes], httpStatusCodes: [503] }
    perRetry: { backOff: { baseInterval: 100ms, maxInterval: 1s } }
  circuitBreaker: { maxPendingRequests: 1024, maxParallelRequests: 1024 }
  rateLimit:
    type: Local
    local: { rules: [{ limit: { requests: 200, unit: Second } }] }
```

Retrying `POST /orders` at the edge is safe **only because** the API
requires an `Idempotency-Key`. That's the payoff of designing for it in
§11.

> ⚠️ **Caveat (verify):** Envoy Gateway's `BackendTrafficPolicy` is `v1alpha1`, and its field names (retry, circuitBreaker, rateLimit) have moved between releases. Validate against the installed CRDs: `kubectl apply --dry-run=server -f policy.yaml`. Local rate limiting may need extra configuration depending on the version.

---

## 26. Packaging: Kustomize and Helm

Rule of thumb used by many platform teams: **Helm for software you
consume** (Prometheus, cert-manager, CNPG), **Kustomize for software you
own** (plain YAML you can read, patched per environment). Both are
first-class in Argo CD.

### 26.1 Kustomize: base + overlays

```
deploy/
├── base/
│   ├── kustomization.yaml
│   ├── orders/{deployment.yaml,service.yaml,pdb.yaml,serviceaccount.yaml,config.env}
│   ├── catalog/...
│   ├── worker/...
│   ├── redis/statefulset.yaml
│   ├── postgres/cluster.yaml
│   └── routes/httproute.yaml
└── overlays/
    ├── dev/kustomization.yaml
    └── prod/{kustomization.yaml,patches/*.yaml}
```

```yaml
# deploy/base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: orbit
labels:
  - pairs: { app.kubernetes.io/part-of: orbit }
    includeSelectors: false
resources:
  - orders/deployment.yaml
  - orders/service.yaml
  - orders/pdb.yaml
  - orders/serviceaccount.yaml
  - catalog/
  - worker/
  - redis/statefulset.yaml
  - postgres/cluster.yaml
  - routes/httproute.yaml
configMapGenerator:
  - name: orders-config
    envs: [orders/config.env]      # name gets a content hash -> config change = rolling update
images:
  - { name: orbit-orders,  newName: ghcr.io/acme/orbit-orders }
  - { name: orbit-catalog, newName: ghcr.io/acme/orbit-catalog }
  - { name: orbit-worker,  newName: ghcr.io/acme/orbit-worker }
```

```yaml
# deploy/overlays/prod/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: orbit-prod
resources: [../../base]
images:   # CI bumps these (§35): digest-pinned
  - { name: orbit-orders,  newName: ghcr.io/acme/orbit-orders,  digest: sha256:4f1c... }
  - { name: orbit-catalog, newName: ghcr.io/acme/orbit-catalog, digest: sha256:9ab0... }
  - { name: orbit-worker,  newName: ghcr.io/acme/orbit-worker,  digest: sha256:c77e... }
patches:
  - path: patches/orders-resources.yaml
  - target: { kind: Cluster, name: orbit-db }
    patch: |-
      - op: replace
        path: /spec/storage/size
        value: 100Gi
```

```bash
kustomize build deploy/overlays/prod | less               # render
kubectl diff -k deploy/overlays/prod                       # compare with the cluster
kustomize build deploy/overlays/prod | kubeconform -strict -summary -ignore-missing-schemas
```

### 26.2 Helm: authoring a chart

```bash
helm create orbit-service     # scaffold; then delete what you don't need
```

```
orbit-service/
├── Chart.yaml
├── values.yaml
├── values.schema.json        # validates values: catches typos like "replica: 3"
└── templates/
    ├── _helpers.tpl
    ├── deployment.yaml
    ├── service.yaml
    ├── pdb.yaml
    └── hpa.yaml
```

```yaml
# templates/deployment.yaml (excerpt)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "orbit-service.fullname" . }}
  labels: {{- include "orbit-service.labels" . | nindent 4 }}
spec:
  {{- if not .Values.autoscaling.enabled }}
  replicas: {{ .Values.replicaCount }}
  {{- end }}
  selector:
    matchLabels: {{- include "orbit-service.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      annotations:
        # restart pods when the rendered config changes
        checksum/config: {{ include (print $.Template.BasePath "/configmap.yaml") . | sha256sum }}
      labels: {{- include "orbit-service.selectorLabels" . | nindent 8 }}
    spec:
      containers:
        - name: {{ .Chart.Name }}
          image: "{{ .Values.image.repository }}{{ if .Values.image.digest }}@{{ .Values.image.digest }}{{ else }}:{{ .Values.image.tag | default .Chart.AppVersion }}{{ end }}"
          resources: {{- toYaml .Values.resources | nindent 12 }}
```

```bash
helm lint ./orbit-service
helm template orders ./orbit-service -f values-prod.yaml | kubeconform -strict
helm upgrade --install orders ./orbit-service -n orbit -f values-prod.yaml --atomic --timeout 5m
helm history orders -n orbit
helm rollback orders 3 -n orbit
helm plugin install https://github.com/databus23/helm-diff && helm diff upgrade orders ./orbit-service -f values-prod.yaml
```

`--atomic` rolls back automatically if the upgrade doesn't become healthy
within `--timeout`.

> **Best practice:** never `helm install` third-party charts with default
> values in prod. **Pin chart versions**, commit a values file, render it
> (`helm template`) in CI, and review the diff. A chart's minor version
> bump can change CRDs, RBAC or defaults.

### 26.3 Validate manifests before they reach a cluster

| Tool | Catches |
|---|---|
| `kubeconform` | schema errors (typos, wrong types), including CRDs with schemas |
| `kube-linter` / `kube-score` / Polaris | missing probes, limits, securityContext, `latest` tags |
| `conftest` (OPA/Rego) / Kyverno CLI | your org's policies, the same ones enforced at admission (§34.5) |
| `pluto` / `kubent` | deprecated/removed API versions before an upgrade |

---

## 27. Autoscaling: pods, sizes, and nodes

### 27.1 Three layers

```
 load ↑ ──► HPA / KEDA: more PODS ──► pods Pending (no room) ──► Cluster Autoscaler / Karpenter: more NODES
           VPA: right-size each pod's requests (vertical)
```

They depend on each other: **HPA percentages are relative to requests**,
and node autoscalers react to **Pending pods' requests**, not actual usage.
Wrong requests break all three layers.

### 27.2 HPA on CPU (catalog)

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata: { name: catalog }
spec:
  scaleTargetRef: { apiVersion: apps/v1, kind: Deployment, name: catalog }
  minReplicas: 3                    # >= number of zones: survive a zone loss
  maxReplicas: 30
  metrics:
    - type: Resource
      resource: { name: cpu, target: { type: Utilization, averageUtilization: 65 } }
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 0
      policies:
        - { type: Percent, value: 100, periodSeconds: 30 }   # at most double every 30s
        - { type: Pods,    value: 4,   periodSeconds: 30 }
      selectPolicy: Max
    scaleDown:
      stabilizationWindowSeconds: 300                         # wait 5 min of low load: no flapping
      policies: [{ type: Percent, value: 20, periodSeconds: 60 }]
```

The algorithm:
`desired = ceil(current × currentMetric / target)`. With 4 pods at 90% CPU
and a 65% target, `ceil(4 × 90/65) = 6`. The HPA syncs every 15 s and
ignores changes within a 10% tolerance.

> **Best practice:** target **60–70%**, not 90%. Scaling takes time (HPA
> sync + scheduling + image pull + startup + readiness ≈ 30–120 s). The
> headroom absorbs the spike while new pods arrive. Remove `replicas:` from
> the Deployment manifest when an HPA owns it, or every GitOps sync resets
> the count.

### 27.3 KEDA: scale on what actually matters

CPU is a lagging, indirect signal. A queue consumer should scale on
**backlog**, and an API on **request rate** or **latency**. KEDA
(CNCF-graduated) feeds external metrics to an HPA it manages, and can
**scale to zero**.

```bash
helm repo add kedacore https://kedacore.github.io/charts
helm upgrade --install keda kedacore/keda -n keda --create-namespace --wait
```

Worker scales on Redis Stream lag (0 → 20):

```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata: { name: worker, namespace: orbit }
spec:
  scaleTargetRef: { name: worker }
  minReplicaCount: 0              # scale to zero when there's nothing to do
  maxReplicaCount: 20
  pollingInterval: 10
  cooldownPeriod: 120             # seconds of zero activity before going to 0
  triggers:
    - type: redis-streams
      metadata:
        address: redis.orbit.svc.cluster.local:6379
        stream: orders.created
        consumerGroup: workers
        lagCount: "50"            # target: ~50 unprocessed entries per replica
        activationLagCount: "1"   # wake from 0 when anything arrives
```

Orders scales on request rate from Prometheus (set up in §31):

```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata: { name: orders, namespace: orbit }
spec:
  scaleTargetRef: { name: orders }
  minReplicaCount: 3
  maxReplicaCount: 30
  advanced:
    horizontalPodAutoscalerConfig:
      behavior:
        scaleDown: { stabilizationWindowSeconds: 300 }
  triggers:
    - type: prometheus
      metadata:
        serverAddress: http://kps-prometheus.monitoring.svc:9090
        query: sum(rate(http_server_request_duration_seconds_count{job="orbit/orders"}[1m]))
        threshold: "100"          # ~100 RPS per pod (measured by load test, §27.7)
    - type: cpu                   # safety net: CPU too
      metricType: Utilization
      metadata: { value: "70" }
```

KEDA has 70+ scalers: Kafka lag, SQS/PubSub/Service Bus depth, RabbitMQ,
NATS, Postgres queries, cron schedules (pre-scale before a known daily
peak), and more. `ScaledJob` runs one Job per message for long tasks such
as video transcoding.

> **War story:** a team autoscaled Kafka consumers on CPU. Consumers were
> I/O-bound (waiting on a slow downstream), so CPU stayed at 20% while lag
> grew to 6 hours. Scaling on consumer lag fixed it. Note that for Kafka,
> **replicas beyond the partition count do nothing**, so set
> `maxReplicaCount ≤ partitions`.

> ⚠️ **Caveat (verify):** the KEDA `redis-streams` trigger's `lagCount`/`activationLagCount` need a recent KEDA and Redis 7+. On older versions use `pendingEntriesCount`, which measures delivered-but-unacked entries, not backlog. Check the KEDA scaler docs for your version. The Prometheus trigger assumes the metric and label names in §30; check that the query returns data in Prometheus first.

### 27.4 Scaling on latency, carefully

Scaling on p95 latency sounds ideal, but latency rises for reasons that
more pods don't fix (a slow DB, lock contention). Adding pods then adds DB
connections and makes it worse. Prefer **throughput or saturation signals**
(RPS, in-flight requests, queue depth, CPU) for scaling, and use
**latency for alerting**.

Also check **connection math**: `maxReplicas × DB_MAX_CONNS` must stay
below Postgres `max_connections` (30 × 10 = 300 > 200 in our config!). Use
a pooler (PgBouncer, which CNPG provides via its `Pooler` resource) or
lower per-pod pools.

### 27.5 VPA: right-sizing recommendations

```yaml
apiVersion: autoscaling.k8s.io/v1
kind: VerticalPodAutoscaler
metadata: { name: orders }
spec:
  targetRef: { apiVersion: apps/v1, kind: Deployment, name: orders }
  updatePolicy: { updateMode: "Off" }     # recommend only; a human/PR applies it
```

```bash
kubectl describe vpa orders     # Target / Lower Bound / Upper Bound per container
```

Don't let VPA (in `Auto`/`Recreate` mode) and HPA act on the **same
metric** (CPU/memory); they fight. The common pattern is HPA on
RPS/queue, plus VPA in `Off` mode feeding periodic right-sizing PRs.

### 27.6 Node autoscaling: Cluster Autoscaler vs Karpenter

| | Cluster Autoscaler | Karpenter |
|---|---|---|
| Model | scales pre-defined node groups (ASGs/MIGs) | provisions individual instances to fit pending pods |
| Instance choice | fixed per group | picks from many types/sizes, spot/on-demand, arch |
| Consolidation | removes underutilized nodes | actively repacks pods onto cheaper/fewer nodes |
| Clouds | all | AWS (native), Azure (AKS NAP); others emerging |

```yaml
apiVersion: karpenter.sh/v1
kind: NodePool
metadata: { name: general }
spec:
  template:
    spec:
      nodeClassRef: { group: karpenter.k8s.aws, kind: EC2NodeClass, name: default }
      requirements:
        - { key: kubernetes.io/arch,             operator: In, values: [amd64, arm64] }
        - { key: karpenter.sh/capacity-type,     operator: In, values: [spot, on-demand] }
        - { key: karpenter.k8s.aws/instance-category, operator: In, values: [c, m, r] }
      expireAfter: 720h                # recycle nodes monthly: patched AMIs, no snowflakes
  limits: { cpu: "500" }               # cost ceiling
  disruption:
    consolidationPolicy: WhenEmptyOrUnderutilized
    consolidateAfter: 5m
    budgets: [{ nodes: "10%" }]        # disrupt at most 10% of nodes at once
```

**Overprovisioning** hides node spin-up time (60–120 s). Run low-priority
"balloon" pods that reserve spare capacity. Real pods preempt them
instantly, and the evicted balloons trigger a new node in the background.

```yaml
apiVersion: scheduling.k8s.io/v1
kind: PriorityClass
metadata: { name: overprovisioning }
value: -10
---
apiVersion: apps/v1
kind: Deployment
metadata: { name: overprovisioning, namespace: kube-system }
spec:
  replicas: 2
  selector: { matchLabels: { app: overprovisioning } }
  template:
    metadata: { labels: { app: overprovisioning } }
    spec:
      priorityClassName: overprovisioning
      containers:
        - name: pause
          image: registry.k8s.io/pause:3.10
          resources: { requests: { cpu: "1", memory: 2Gi } }
```

> ⚠️ **Caveat (verify):** Karpenter field names follow the `karpenter.sh/v1` API (v1.0+). The `instance-category` label and `EC2NodeClass` details are AWS-specific and change with releases. Validate with `kubectl apply --dry-run=server` on a cluster that has Karpenter installed.

### 27.7 Prove it: a load test that triggers scaling

```javascript
// load/orders.js
import http from 'k6/http';
import { check } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

export const options = {
  scenarios: {
    ramp: {
      executor: 'ramping-arrival-rate',        // open model: fixed RPS regardless of latency
      startRate: 20, timeUnit: '1s', preAllocatedVUs: 200, maxVUs: 1000,
      stages: [
        { target: 200, duration: '2m' },
        { target: 800, duration: '5m' },       // should force scale-out
        { target: 800, duration: '5m' },
        { target: 0,   duration: '2m' },
      ],
    },
  },
  thresholds: {
    http_req_failed:   ['rate<0.01'],          // < 1% errors
    http_req_duration: ['p(95)<300'],          // p95 < 300 ms
  },
  insecureSkipTLSVerify: true,
};

const BASE = __ENV.BASE || 'https://orbit.localtest.me';

export default function () {
  const res = http.post(`${BASE}/api/orders`,
    JSON.stringify({ product_id: 1 + Math.floor(Math.random() * 3), quantity: 1 }),
    { headers: { 'Content-Type': 'application/json', 'Idempotency-Key': uuidv4() } });
  check(res, { created: (r) => r.status === 201 });
}
```

```bash
k6 run load/orders.js &
watch -n2 'kubectl get hpa,scaledobject -n orbit; kubectl get pods -n orbit | grep -c Running'
kubectl get events -n orbit --field-selector reason=SuccessfulRescale -w
```

What you should observe: orders replicas climb as RPS crosses ~100/pod;
worker goes 0 → N as lag builds and back to 0 after the cooldown; p95
holds below threshold. If it doesn't, Part IV gives you the tools to see
why.

**Part III checkpoint:** zero-error rolling deploys (rerun the §20
exercise 1 during a k6 run: 0 failures), pods spread across zones,
autoscaling on meaningful signals, and one golden manifest you can
justify line by line.

---

# Part IV — Observability

## 28. Observability strategy

### 28.1 Signals and what each one answers

| Signal | Answers | Strength | Cost driver |
|---|---|---|---|
| **Metrics** | "Is something wrong? How much? Since when?" | cheap aggregates, alerting, trends | **cardinality** (unique label combinations) |
| **Logs** | "What exactly happened in this event?" | detail, context, audit | **volume** (bytes ingested) |
| **Traces** | "Where did the time go across services for this request?" | causality across hops | volume (sample!) |
| **Profiles** | "Which line of code burns CPU/allocates memory?" | continuous pprof (Pyroscope, Parca) | low |
| **K8s events** | "What did the platform do?" (OOMKill, eviction, scaling) | explains restarts | low, but they expire after 1 h |

The workflow they enable: **an alert fires on a metric → a dashboard
narrows it to a service and time → an exemplar or trace shows the slow
hop → that span's logs (same `trace_id`) show the error.** Correlation by
shared resource attributes (`service.name`, `k8s.pod.name`) and `trace_id`
is what turns three data stores into one investigation.

### 28.2 What to measure: Golden Signals, RED, USE

- **RED** (for every *service*): **R**ate, **E**rrors, **D**uration.
- **USE** (for every *resource*: CPU, memory, disk, connection pools,
  queues): **U**tilization, **S**aturation, **E**rrors.
- **Four Golden Signals** (Google SRE book): latency, traffic, errors,
  saturation, which is RED + saturation.

For Orbit:

| Component | Signals |
|---|---|
| catalog, orders | RED per route; DB pool in-use/wait; cache hit ratio |
| worker | messages processed/s, processing duration, failures, **stream lag** (saturation) |
| Postgres | connections vs max, replication lag, TPS, slow queries, disk |
| Redis | memory vs maxmemory, ops/s, stream length |
| Kubernetes | restarts, OOMKills, not-ready pods, HPA at max, pending pods, node pressure, PVC usage |

### 28.3 Why OpenTelemetry

OpenTelemetry (OTel) is the CNCF standard for generating and shipping
telemetry: **one SDK per language, one wire protocol (OTLP), one
Collector**. The backend becomes a swappable detail: Prometheus/Loki/Tempo
today, Datadog/Honeycomb/New Relic/Grafana Cloud tomorrow, by changing
Collector config, not application code.

```
 app (OTel SDK) ──OTLP──► Collector AGENT (DaemonSet, per node) ──OTLP──► Collector GATEWAY (Deployment)
 app stdout logs ──files──►   │  + k8s metadata, batching                  │ tail sampling, routing, redaction
                              │                                            ├──► Prometheus (metrics, OTLP receiver)
                              │                                            ├──► Loki       (logs, OTLP)
                              │                                            └──► Tempo      (traces, OTLP)
 kube-state-metrics, node-exporter, kubelet/cAdvisor ◄──scrape── Prometheus
```

> **Best practice:** apps send to a **local agent**, never directly to a
> SaaS endpoint. The agent buffers through backend outages, adds K8s
> metadata, and lets you change backends, sampling and redaction without
> redeploying 200 services.

---

## 29. Instrumenting Orbit with OpenTelemetry

```bash
go get go.opentelemetry.io/otel \
       go.opentelemetry.io/otel/sdk \
       go.opentelemetry.io/otel/sdk/metric \
       go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc \
       go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc \
       go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp \
       go.opentelemetry.io/contrib/instrumentation/runtime \
       github.com/exaring/otelpgx \
       github.com/redis/go-redis/extra/redisotel/v9
```

### 29.1 SDK setup, configured entirely by standard env vars

`internal/platform/telemetry.go`:

```go
package platform

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// SetupOTel installs global tracer/meter providers. Everything is driven by standard env vars:
//   OTEL_SERVICE_NAME, OTEL_RESOURCE_ATTRIBUTES, OTEL_EXPORTER_OTLP_ENDPOINT,
//   OTEL_TRACES_SAMPLER(_ARG), OTEL_SDK_DISABLED ...
// so the same binary works in Compose, kind and prod without code changes.
func SetupOTel(ctx context.Context) (func(context.Context) error, error) {
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
	)
	if err != nil {
		return nil, err
	}

	traceExp, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExp), // async, batched: never block a request on export
		sdktrace.WithResource(res),
	)

	metricExp, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		return nil, err
	}
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExp, sdkmetric.WithInterval(15*time.Second))),
		sdkmetric.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(mp)
	// W3C traceparent + baggage: what every other OTel-instrumented hop understands.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}))

	// Go runtime metrics: goroutines, GC, heap, scheduler latency.
	if err := runtime.Start(); err != nil {
		return nil, err
	}

	return func(ctx context.Context) error {
		return errors.Join(tp.Shutdown(ctx), mp.Shutdown(ctx)) // flush on exit
	}, nil
}
```

> ⚠️ **Caveat (verify):** the OTel Go packages move fast. Run `go build` against current versions and check whether `runtime.Start()` options or exporter constructors changed. All exporter settings come from the `OTEL_*` env vars; if spans don't arrive, set `OTEL_EXPORTER_OTLP_INSECURE=true` and confirm the endpoint scheme.

### 29.2 Logs that carry the trace

Applications log **structured JSON to stdout**. The node agent ships it
(§30). The only addition is putting the active `trace_id`/`span_id` on
every log line:

```go
package platform

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

type traceHandler struct{ slog.Handler }

func (h traceHandler) Handle(ctx context.Context, r slog.Record) error {
	if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
		r.AddAttrs(slog.String("trace_id", sc.TraceID().String()),
			slog.String("span_id", sc.SpanID().String()))
	}
	return h.Handler.Handle(ctx, r)
}
func (h traceHandler) WithAttrs(a []slog.Attr) slog.Handler { return traceHandler{h.Handler.WithAttrs(a)} }
func (h traceHandler) WithGroup(n string) slog.Handler       { return traceHandler{h.Handler.WithGroup(n)} }

// NewLogger: JSON to stdout, level from LOG_LEVEL, trace-correlated.
func NewLogger() *slog.Logger {
	var lvl slog.Level
	_ = lvl.UnmarshalText([]byte(Env("LOG_LEVEL", "info")))
	return slog.New(traceHandler{slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})})
}
```

Use `slog.InfoContext(ctx, ...)` / `slog.ErrorContext(ctx, ...)`, which pass
the context in. Plain `slog.Info` can't see the span.

Logging rules that keep log bills and incidents sane:

- One event per line, JSON, with stable keys (`order_id`, not
  `"order 123 failed"` prose).
- **Never log secrets or PII**: tokens, passwords, card numbers, full
  emails. Redact at the source; the Collector can redact as a second line
  of defense (§30.3).
- `INFO` for state changes, `ERROR` for things a human may need to act on.
  Don't log every successful request at INFO in a high-RPS service;
  metrics and traces cover that.

### 29.3 Auto-instrumentation: HTTP, Postgres, Redis

```go
// internal/platform/http.go
// Route registers a handler wrapped in an OTel span + http.server.* metrics named after the pattern.
func Route(mux *http.ServeMux, pattern string, h http.HandlerFunc) {
	mux.Handle(pattern, otelhttp.NewHandler(h, pattern))
}
```

`cmd/orders/main.go`, the changed lines only:

```go
func main() {
	slog.SetDefault(platform.NewLogger())
	ctx := context.Background()

	shutdownOTel, err := platform.SetupOTel(ctx)
	if err != nil {
		slog.Error("otel setup", "err", err) // telemetry failure must not stop the service
	} else {
		defer func() {
			c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = shutdownOTel(c)
		}()
	}

	pcfg, _ := pgxpool.ParseConfig(platform.MustEnv("DATABASE_URL"))
	pcfg.ConnConfig.Tracer = otelpgx.NewTracer()           // a span per query
	db, _ := pgxpool.NewWithConfig(ctx, pcfg)
	_ = otelpgx.RecordStats(db)                             // pool metrics: acquired, idle, wait time

	rdb := redis.NewClient(&redis.Options{Addr: platform.MustEnv("REDIS_ADDR")})
	_ = redisotel.InstrumentTracing(rdb)
	_ = redisotel.InstrumentMetrics(rdb)

	// ...
	platform.Route(mux, "POST /orders", h.Create)          // instead of mux.HandleFunc
	mux.HandleFunc("GET /livez", health.Livez)              // probes: deliberately NOT traced (noise)
	mux.HandleFunc("GET /readyz", health.Readyz)
	// ...
}
```

Outbound HTTP calls get the client-side half:
`&http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}`
injects `traceparent` into outgoing requests.

### 29.4 Business metrics and spans

```go
package orders

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	tracer = otel.Tracer("github.com/acme/orbit/orders")
	meter  = otel.Meter("github.com/acme/orbit/orders")

	ordersCreated, _ = meter.Int64Counter("orbit.orders.created",
		metric.WithDescription("Orders successfully created"), metric.WithUnit("{order}"))
	orderQuantity, _ = meter.Int64Histogram("orbit.orders.quantity",
		metric.WithDescription("Items per order"), metric.WithExplicitBucketBoundaries(1, 2, 5, 10, 50))
)

// inside Create, after a successful insert:
//   ordersCreated.Add(ctx, 1, metric.WithAttributes(attribute.Bool("idempotent_replay", !created)))
//   orderQuantity.Record(ctx, int64(req.Quantity))
//
// a custom span around a meaningful unit of work:
//   ctx, span := tracer.Start(ctx, "orders.insert")
//   defer span.End()
//   span.SetAttributes(attribute.Int64("orbit.product_id", req.ProductID))
//   if err != nil { span.RecordError(err); span.SetStatus(codes.Error, "insert failed") }
```

> **War story: cardinality.** Someone added `attribute.String("user_id",
> ...)` to a request counter. With 2 million users × 20 routes × 5 status
> codes, Prometheus went from 300k to 200M active series and ran out of
> memory, taking alerting down with it. **Metric attributes must be
> bounded** (route, method, status class, region). Unbounded values
> (user, order, request IDs) belong on **spans and logs**.

> ⚠️ **Caveat (verify):** **metric names depend on the otelhttp version.** Releases following the current HTTP semantic conventions emit `http.server.request.duration` in seconds (Prometheus: `http_server_request_duration_seconds_*`). Older releases emitted `http.server.duration` in **milliseconds**. Every PromQL query, alert, SLO, KEDA trigger and canary analysis in this guide assumes the first. Check what you actually get: `curl -s localhost:9090/api/v1/label/__name__/values | jq -r '.data[]' | grep http_server` (port-forward Prometheus first). Also check that `otelpgx.RecordStats` and `redisotel.InstrumentMetrics` exist in the versions you pin.

### 29.5 Propagating traces through the queue

HTTP propagation is automatic. A message queue is not: inject the context
into the message and extract it on the consumer side, so one trace covers
`POST /orders → XADD → worker → UPDATE`.

```go
// producer (orders): add the W3C trace context to the stream entry
carrier := propagation.MapCarrier{}
otel.GetTextMapPropagator().Inject(ctx, carrier)
values := map[string]any{"order_id": id}
for k, v := range carrier { // "traceparent", optionally "tracestate"/"baggage"
	values[k] = v
}
h.RDB.XAdd(ctx, &redis.XAddArgs{Stream: Stream, Values: values, MaxLen: 100_000, Approx: true})

// consumer (worker): continue the trace
carrier := propagation.MapCarrier{}
for k, v := range m.Values {
	if s, ok := v.(string); ok {
		carrier[k] = s
	}
}
ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
ctx, span := tracer.Start(ctx, "orders.created process",
	trace.WithSpanKind(trace.SpanKindConsumer),
	trace.WithAttributes(
		attribute.String("messaging.system", "redis"),
		attribute.String("messaging.destination.name", Stream),
		attribute.String("messaging.message.id", m.ID)))
defer span.End()
```

Worker metrics to add: `orbit.worker.processed` (counter, attribute
`outcome=paid|failed|retry`) and `orbit.worker.process.duration`
(histogram, seconds).

### 29.6 Configuration in Kubernetes

Added to every service's base manifest (and to the Compose file):

```yaml
env:
  - name: HOST_IP
    valueFrom: { fieldRef: { fieldPath: status.hostIP } }
  - name: POD_NAME
    valueFrom: { fieldRef: { fieldPath: metadata.name } }
  - { name: OTEL_SERVICE_NAME, value: orders }
  - { name: OTEL_EXPORTER_OTLP_ENDPOINT, value: "http://$(HOST_IP):4317" }     # node-local agent
  - name: OTEL_RESOURCE_ATTRIBUTES
    value: "service.namespace=orbit,service.version=$(APP_VERSION),deployment.environment.name=prod,k8s.pod.name=$(POD_NAME)"
```

(`APP_VERSION` is set by Kustomize from the image tag; `service.version`
is what lets dashboards compare canary vs stable in §37.)

### 29.7 Zero-code instrumentation (for services you don't own)

The **OpenTelemetry Operator** injects auto-instrumentation agents (Java,
Python, Node.js, .NET; Go via eBPF) using a pod annotation. That's useful
for third-party or legacy apps:

```yaml
apiVersion: opentelemetry.io/v1alpha1
kind: Instrumentation
metadata: { name: default, namespace: legacy }
spec:
  exporter: { endpoint: http://otel-agent.monitoring:4318 }
  propagators: [tracecontext, baggage]
  sampler: { type: parentbased_traceidratio, argument: "1" }
---
# on the Deployment's pod template:
#   annotations: { instrumentation.opentelemetry.io/inject-java: "true" }
```

Grafana **Beyla** / OTel eBPF instrumentation can produce RED metrics and
traces for *any* HTTP/gRPC service at the kernel level, with no code and
no restarts.

> ⚠️ **Caveat (verify):** the Operator `Instrumentation` exporter endpoint, ports and Go eBPF support vary by Operator version. Grafana Beyla is being folded into the OpenTelemetry eBPF instrumentation project, so check which one is current.

---

## 30. Deploying the observability stack

All components go into the `monitoring` namespace. Values files live in
`platform/` and are deployed by Argo CD in Part VII. Here we install them
by hand once to understand them.

```bash
kubectl create namespace monitoring
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts
helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
helm repo update
```

### 30.1 Prometheus, Alertmanager, Grafana: kube-prometheus-stack

One chart provides the Prometheus Operator, Prometheus, Alertmanager,
Grafana, node-exporter, kube-state-metrics, and a battle-tested set of
Kubernetes alerts and dashboards.

```yaml
# platform/kube-prometheus-stack/values.yaml
fullnameOverride: kps
prometheus:
  prometheusSpec:
    retention: 15d
    retentionSize: 40GB
    enableOTLPReceiver: true                 # accept OTLP metrics at /api/v1/otlp (Prometheus 3)
    otlp:
      promoteResourceAttributes:             # resource attrs -> labels (default: only in target_info)
        - service.version
        - deployment.environment.name
        - k8s.namespace.name
        - k8s.pod.name
        - k8s.deployment.name
    enableFeatures: [exemplar-storage]       # trace IDs attached to histogram samples
    # pick up ServiceMonitors/PodMonitors/PrometheusRules from ALL namespaces, not just this release's
    serviceMonitorSelectorNilUsesHelmValues: false
    podMonitorSelectorNilUsesHelmValues: false
    ruleSelectorNilUsesHelmValues: false
    resources:
      requests: { cpu: 500m, memory: 2Gi }
      limits:   { memory: 2Gi }
    storageSpec:
      volumeClaimTemplate:
        spec: { accessModes: [ReadWriteOnce], resources: { requests: { storage: 50Gi } } }
alertmanager:
  config:
    route:
      receiver: default
      group_by: [alertname, namespace, service]
      group_wait: 30s
      group_interval: 5m
      repeat_interval: 4h
      routes:
        - matchers: [severity="page"]
          receiver: pagerduty
        - matchers: [alertname="Watchdog"]   # dead-man's switch: must ALWAYS fire (see §31.4)
          receiver: deadmansswitch
          repeat_interval: 1m
    inhibit_rules:
      - source_matchers: [severity="page"]
        target_matchers: [severity="ticket"]
        equal: [namespace, service]
    receivers:
      - name: default
        slack_configs:
          - api_url_file: /etc/alertmanager/secrets/slack/url   # from a Secret, not inline
            channel: "#orbit-alerts"
            send_resolved: true
            title: '{{ .CommonLabels.alertname }} ({{ .Status }})'
            text: '{{ range .Alerts }}{{ .Annotations.summary }} - runbook: {{ .Annotations.runbook_url }}{{ "\n" }}{{ end }}'
      - name: pagerduty
        pagerduty_configs: [{ routing_key_file: /etc/alertmanager/secrets/pagerduty/key }]
      - name: deadmansswitch
        webhook_configs: [{ url: "https://hc-ping.com/<uuid>" }]
  alertmanagerSpec:
    secrets: [slack, pagerduty]
grafana:
  adminPassword: change-me                  # prod: admin.existingSecret + SSO
  sidecar:
    dashboards: { enabled: true, searchNamespace: ALL }   # dashboards as ConfigMaps (label grafana_dashboard=1)
  additionalDataSources:
    - name: Loki
      type: loki
      uid: loki
      url: http://loki-gateway.monitoring.svc
      jsonData:
        derivedFields:                       # log line -> trace
          - name: trace_id
            matcherType: label
            matcherRegex: trace_id
            datasourceUid: tempo
            url: "$${__value.raw}"
    - name: Tempo
      type: tempo
      uid: tempo
      url: http://tempo.monitoring.svc:3200
      jsonData:
        tracesToLogsV2:                      # trace -> its logs
          datasourceUid: loki
          filterByTraceID: true
          spanStartTimeShift: "-5m"
          spanEndTimeShift: "5m"
        serviceMap: { datasourceUid: prometheus }
        nodeGraph: { enabled: true }
```

```bash
kubectl create secret generic slack -n monitoring --from-literal=url=https://hooks.slack.com/services/...
kubectl create secret generic pagerduty -n monitoring --from-literal=key=...
helm upgrade --install kps prometheus-community/kube-prometheus-stack -n monitoring \
  -f platform/kube-prometheus-stack/values.yaml --wait
kubectl port-forward -n monitoring svc/kps-grafana 3000:80     # http://localhost:3000
```

> ⚠️ **Caveat (verify):** the least certain part of the stack. Check against the kube-prometheus-stack/Prometheus Operator version you install. (1) The field names `enableOTLPReceiver` and `otlp.promoteResourceAttributes` in `prometheusSpec`: run `kubectl explain prometheus.spec.otlp` after installing. If they're absent, use `enableFeatures: [otlp-write-receiver]` on older versions, or the Collector's `prometheusremotewrite` exporter with `enableRemoteWriteReceiver: true`. (2) The label mapping: OTLP `service.namespace`/`service.name` → `job="orbit/orders"`, and dots → underscores. (3) Grafana `derivedFields` with `matcherType: label` (structured metadata) needs a recent Grafana; otherwise use a regex matcher on `"trace_id":"(\w+)"`. (4) The `fullnameOverride: kps` → service `kps-prometheus` naming; confirm with `kubectl get svc -n monitoring`.

### 30.2 Loki (logs) and Tempo (traces)

```yaml
# platform/loki/values.yaml: single-binary mode for the lab; use "SimpleScalable"/"Distributed" + object storage in prod
deploymentMode: SingleBinary
loki:
  auth_enabled: false
  commonConfig: { replication_factor: 1 }
  storage: { type: filesystem }             # prod: s3/gcs/azure bucket
  schemaConfig:
    configs:
      - from: "2024-01-01"
        store: tsdb
        object_store: filesystem
        schema: v13
        index: { prefix: index_, period: 24h }
  limits_config:
    allow_structured_metadata: true          # required for OTLP ingestion
    retention_period: 168h
singleBinary: { replicas: 1, persistence: { size: 20Gi } }
read: { replicas: 0 }
write: { replicas: 0 }
backend: { replicas: 0 }
chunksCache: { enabled: false }
resultsCache: { enabled: false }
```

```yaml
# platform/tempo/values.yaml (single binary; tempo-distributed + object storage in prod)
tempo:
  retention: 72h
  receivers:
    otlp:
      protocols:
        grpc: { endpoint: 0.0.0.0:4317 }
        http: { endpoint: 0.0.0.0:4318 }
persistence: { enabled: true, size: 20Gi }
```

```bash
helm upgrade --install loki  grafana/loki  -n monitoring -f platform/loki/values.yaml --wait
helm upgrade --install tempo grafana/tempo -n monitoring -f platform/tempo/values.yaml --wait
```

> ⚠️ **Caveat (verify):** Loki and Tempo chart values change between chart majors, and these were written against Loki chart v6-style values. Grafana has been reorganizing the Tempo charts (single-binary vs `tempo-distributed`). Render with `helm template` and check `kubectl get svc -n monitoring` for the actual service names (`loki-gateway`, `tempo`) that §30.1 and §30.3 point to.

### 30.3 OpenTelemetry Collector: agent (DaemonSet) + gateway (Deployment)

**Agent:** receives OTLP from pods on its node (hostPort 4317), tails
container log files, enriches everything with Kubernetes metadata, and
forwards it to the gateway.

```yaml
# platform/otel-agent/values.yaml
mode: daemonset
fullnameOverride: otel-agent
image: { repository: otel/opentelemetry-collector-contrib }
presets:
  kubernetesAttributes: { enabled: true }   # k8sattributes processor + RBAC
  logsCollection: { enabled: true }         # filelog receiver on /var/log/pods
resources:
  requests: { cpu: 100m, memory: 256Mi }
  limits:   { memory: 512Mi }
config:
  receivers:
    otlp:
      protocols:
        grpc: { endpoint: ${env:MY_POD_IP}:4317 }
        http: { endpoint: ${env:MY_POD_IP}:4318 }
  processors:
    memory_limiter: { check_interval: 1s, limit_percentage: 80, spike_limit_percentage: 25 }
    batch: { send_batch_size: 8192, timeout: 2s }
    resourcedetection: { detectors: [env, system], override: false }
    transform/json-logs:                      # lift JSON fields + trace_id out of the log body
      error_mode: ignore
      log_statements:
        - context: log
          statements:
            - merge_maps(attributes, ParseJSON(body), "upsert") where IsMatch(body, "^\\{")
            - set(trace_id.string, attributes["trace_id"]) where attributes["trace_id"] != nil
            - set(span_id.string, attributes["span_id"]) where attributes["span_id"] != nil
            - set(severity_text, attributes["level"]) where attributes["level"] != nil
  exporters:
    otlp/gateway:
      endpoint: otel-gateway.monitoring.svc:4317
      tls: { insecure: true }                 # mTLS in prod, or rely on a mesh
      sending_queue: { enabled: true, queue_size: 5000 }
      retry_on_failure: { enabled: true, max_elapsed_time: 300s }
    loadbalancing:                            # traces: all spans of a trace -> same gateway pod (tail sampling)
      routing_key: traceID
      protocol: { otlp: { tls: { insecure: true } } }
      resolver:
        k8s: { service: otel-gateway-headless.monitoring }
  service:
    pipelines:
      traces:  { receivers: [otlp], processors: [memory_limiter, k8sattributes, resourcedetection, batch], exporters: [loadbalancing] }
      metrics: { receivers: [otlp], processors: [memory_limiter, k8sattributes, resourcedetection, batch], exporters: [otlp/gateway] }
      logs:    { receivers: [otlp, filelog], processors: [memory_limiter, k8sattributes, transform/json-logs, batch], exporters: [otlp/gateway] }
```

(The `loadbalancing` exporter with the k8s resolver needs RBAC to watch
EndpointSlices; add a small Role, or use the `dns` resolver against the
headless service.)

**Gateway:** central processing, with tail sampling, redaction and fan-out
to backends.

```yaml
# platform/otel-gateway/values.yaml
mode: deployment
fullnameOverride: otel-gateway
replicaCount: 2
image: { repository: otel/opentelemetry-collector-contrib }
resources:
  requests: { cpu: 250m, memory: 512Mi }
  limits:   { memory: 1Gi }
podDisruptionBudget: { enabled: true, minAvailable: 1 }
config:
  receivers:
    otlp: { protocols: { grpc: { endpoint: ${env:MY_POD_IP}:4317 } } }
  processors:
    memory_limiter: { check_interval: 1s, limit_percentage: 80, spike_limit_percentage: 25 }
    batch: {}
    tail_sampling:                            # decide AFTER seeing the whole trace
      decision_wait: 10s
      policies:
        - { name: errors,   type: status_code, status_code: { status_codes: [ERROR] } }
        - { name: slow,     type: latency,     latency: { threshold_ms: 500 } }
        - { name: baseline, type: probabilistic, probabilistic: { sampling_percentage: 10 } }
    attributes/redact:
      actions:
        - { key: http.request.header.authorization, action: delete }
        - { key: db.query.text, action: hash }     # or keep, if queries are parameterized (they should be)
  exporters:
    otlphttp/prometheus:
      endpoint: http://kps-prometheus.monitoring.svc:9090/api/v1/otlp
      tls: { insecure: true }
    otlphttp/loki:
      endpoint: http://loki-gateway.monitoring.svc/otlp
    otlp/tempo:
      endpoint: tempo.monitoring.svc:4317
      tls: { insecure: true }
  service:
    pipelines:
      traces:  { receivers: [otlp], processors: [memory_limiter, tail_sampling, attributes/redact, batch], exporters: [otlp/tempo] }
      metrics: { receivers: [otlp], processors: [memory_limiter, batch], exporters: [otlphttp/prometheus] }
      logs:    { receivers: [otlp], processors: [memory_limiter, attributes/redact, batch], exporters: [otlphttp/loki] }
```

```bash
helm upgrade --install otel-gateway open-telemetry/opentelemetry-collector -n monitoring -f platform/otel-gateway/values.yaml
helm upgrade --install otel-agent   open-telemetry/opentelemetry-collector -n monitoring -f platform/otel-agent/values.yaml
kubectl logs -n monitoring ds/otel-agent | grep -i error
```

Sampling strategy explained:

- **Head sampling** (decided at the first span, in the SDK) is cheap but
  blind: at 10% you drop 90% of your errors.
- **Tail sampling** (decided in the gateway after the trace completes)
  keeps **100% of errors and slow traces** plus a 10% baseline. That's the
  standard production choice. It costs gateway memory and requires that
  all spans of a trace reach the same gateway replica, hence the
  `loadbalancing` exporter.
- **Metrics are never sampled.** They come from the SDK before sampling,
  so RED dashboards and SLOs stay exact.

> **Best practice:** **monitor the monitoring.** Alert on Collector
> `otelcol_exporter_send_failed_*` and `otelcol_processor_refused_*`,
> Prometheus `up == 0` for important targets, and the Watchdog alert
> (§31.4). An observability pipeline that silently drops data fails in the
> worst possible moment: during the incident.

> ⚠️ **Caveat (verify):** (1) OTTL syntax changed in recent Collector releases (paths such as `log.body` with context inference). The `context: log` + `body` form shown here should still work but may print deprecation warnings. (2) Check the `loadbalancing` exporter's `resolver.k8s` RBAC needs and config shape in its README. (3) The chart's `presets` and the daemonset `hostPort: 4317` default depend on the chart version. Validate every config with `docker run --rm -v $PWD:/c otel/opentelemetry-collector-contrib validate --config=/c/config.yaml` (extract the `config:` block first), and pin the chart and image versions.

### 30.4 Scraping things that expose `/metrics`

Not everything speaks OTLP. Postgres (CNPG), Redis exporters and Envoy
expose Prometheus endpoints. Point the Prometheus Operator at them with
`ServiceMonitor`/`PodMonitor`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata: { name: redis, namespace: orbit }
spec:
  selector: { matchLabels: { app.kubernetes.io/name: redis-exporter } }
  endpoints: [{ port: metrics, interval: 30s }]
```

(Run `oliver006/redis_exporter` as a sidecar in the Redis StatefulSet, or
as a small Deployment, to get Redis metrics. CNPG's `enablePodMonitor:
true` from §20 is already scraped.)

---

## 31. Metrics, dashboards, and alerting

### 31.1 PromQL you will actually use

```promql
# --- RED for Orbit services (OTel http.server.request.duration, in seconds) ---
# Rate (RPS) per service
sum by (job) (rate(http_server_request_duration_seconds_count{job=~"orbit/.*"}[5m]))
# Error ratio (5xx / all)
sum by (job) (rate(http_server_request_duration_seconds_count{job=~"orbit/.*",http_response_status_code=~"5.."}[5m]))
  / sum by (job) (rate(http_server_request_duration_seconds_count{job=~"orbit/.*"}[5m]))
# p95 / p99 latency (ALWAYS aggregate buckets by le before histogram_quantile)
histogram_quantile(0.95, sum by (job, le) (rate(http_server_request_duration_seconds_bucket{job=~"orbit/.*"}[5m])))

# --- Business ---
sum(rate(orbit_orders_created_total[5m])) * 60                         # orders/minute
histogram_quantile(0.95, sum by (le) (rate(orbit_worker_process_duration_seconds_bucket[5m])))

# --- Kubernetes health (kube-state-metrics + cAdvisor) ---
sum by (namespace, pod) (increase(kube_pod_container_status_restarts_total{namespace="orbit"}[1h])) > 0
kube_pod_container_status_last_terminated_reason{reason="OOMKilled", namespace="orbit"}
sum by (namespace) (kube_pod_status_phase{phase="Pending"}) > 0
kube_deployment_status_replicas_available / kube_deployment_spec_replicas < 1
kube_horizontalpodautoscaler_status_current_replicas >= kube_horizontalpodautoscaler_spec_max_replicas   # HPA maxed out
# memory usage vs limit, per container
max by (pod, container) (container_memory_working_set_bytes{namespace="orbit",container!=""})
  / max by (pod, container) (kube_pod_container_resource_limits{namespace="orbit",resource="memory"})
# PVC will be full within 4 days?
predict_linear(kubelet_volume_stats_available_bytes{namespace="orbit"}[6h], 4 * 24 * 3600) < 0

# --- USE for nodes (node-exporter) ---
1 - avg by (instance) (rate(node_cpu_seconds_total{mode="idle"}[5m]))      # CPU utilization
node_pressure_cpu_waiting_seconds_total                                    # PSI saturation (rate it)
```

PromQL rules of thumb:

- `rate()` on counters, **then** `sum`. Never `sum` and then `rate`.
- Range ≥ 4× the scrape/export interval (`[1m]` for 15 s, `[5m]` is a safe
  default).
- `histogram_quantile` needs `le` in the `by` clause, and you can't
  average percentiles across pods. Aggregate buckets instead.

### 31.2 Recording rules: precompute what dashboards and alerts reuse

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata: { name: orbit-recording, namespace: monitoring }
spec:
  groups:
    - name: orbit.red
      interval: 30s
      rules:
        - record: job:http_requests:rate5m
          expr: sum by (job) (rate(http_server_request_duration_seconds_count{job=~"orbit/.*"}[5m]))
        - record: job:http_errors:ratio_rate5m
          expr: |
            sum by (job) (rate(http_server_request_duration_seconds_count{job=~"orbit/.*",http_response_status_code=~"5.."}[5m]))
            / sum by (job) (rate(http_server_request_duration_seconds_count{job=~"orbit/.*"}[5m]))
        - record: job:http_latency:p95_5m
          expr: histogram_quantile(0.95, sum by (job, le) (rate(http_server_request_duration_seconds_bucket{job=~"orbit/.*"}[5m])))
```

### 31.3 Alert rules: page on symptoms, ticket on causes

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata: { name: orbit-alerts, namespace: monitoring }
spec:
  groups:
    - name: orbit.alerts
      rules:
        - alert: OrbitWorkerBacklogGrowing
          expr: |
            sum(rate(orbit_orders_created_total[10m])) > 0
            and sum(rate(orbit_worker_processed_total{outcome="paid"}[10m])) < 0.8 * sum(rate(orbit_orders_created_total[10m]))
          for: 15m
          labels: { severity: ticket, service: worker }
          annotations:
            summary: "Worker processes < 80% of incoming orders for 15m"
            runbook_url: https://runbooks.acme.example/orbit/worker-backlog
        - alert: OrbitPodCrashLooping
          expr: increase(kube_pod_container_status_restarts_total{namespace=~"orbit.*"}[15m]) > 3
          for: 5m
          labels: { severity: ticket }
          annotations:
            summary: "{{ $labels.pod }} restarted {{ $value | humanize }} times in 15m"
            runbook_url: https://runbooks.acme.example/k8s/crashloop
        - alert: OrbitOOMKilled
          expr: increase(container_oom_events_total{namespace=~"orbit.*"}[30m]) > 0
          labels: { severity: ticket }
          annotations: { summary: "{{ $labels.pod }}/{{ $labels.container }} OOMKilled" }
        - alert: OrbitHPAMaxedOut
          expr: kube_horizontalpodautoscaler_status_current_replicas{namespace=~"orbit.*"} >= kube_horizontalpodautoscaler_spec_max_replicas{namespace=~"orbit.*"}
          for: 30m
          labels: { severity: ticket }
          annotations: { summary: "{{ $labels.horizontalpodautoscaler }} at maxReplicas for 30m: raise the max or find the bottleneck" }
        - alert: OrbitPostgresConnectionsHigh
          expr: sum(cnpg_backends_total{namespace=~"orbit.*"}) / max(cnpg_pg_settings_setting{name="max_connections",namespace=~"orbit.*"}) > 0.8
          for: 10m
          labels: { severity: ticket, service: postgres }
          annotations: { summary: "Postgres connections > 80% of max_connections" }
```

What makes a good alert:

- **Page only on user-visible symptoms** (SLO burn, §32). "CPU at 90%" is
  not a page. If users are fine, it can wait until morning.
- **Every alert has a runbook URL** with what it means, how to verify, and
  how to mitigate.
- **`for:` durations** suppress flapping on transient blips.
- **Every page must be actionable.** If the on-call's response is "ack and
  ignore", delete the alert or downgrade it to a ticket. Alert fatigue is
  how real pages get missed.

> ⚠️ **Caveat (verify):** every query here assumes the metric names from the §29 caveat and the label promotion from the §30.1 caveat. Exporter metric names (`cnpg_backends_total`, `cnpg_pg_settings_setting`, `container_oom_events_total`) also vary by exporter version. Check each in Prometheus's Explore view and run `promtool check rules` on the extracted rule groups before deploying.

### 31.4 The dead man's switch

kube-prometheus-stack ships an always-firing `Watchdog` alert. We route it
to an external heartbeat service (healthchecks.io, Cronitor, PagerDuty's
heartbeat). If the heartbeats **stop**, the external service pages you.
That's the only way to learn your Prometheus/Alertmanager pipeline itself
is dead.

### 31.5 Dashboards as code

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: orbit-red-dashboard
  namespace: monitoring
  labels: { grafana_dashboard: "1" }   # the Grafana sidecar loads it automatically
data:
  orbit-red.json: |
    { "title": "Orbit / RED", "uid": "orbit-red", "schemaVersion": 39,
      "templating": { "list": [ { "name": "job", "type": "query", "datasource": "prometheus",
          "query": "label_values(http_server_request_duration_seconds_count{job=~\"orbit/.*\"}, job)" } ] },
      "panels": [
        { "type": "timeseries", "title": "Requests/s", "gridPos": {"x":0,"y":0,"w":8,"h":8},
          "targets": [ { "expr": "sum by (job) (rate(http_server_request_duration_seconds_count{job=~\"$job\"}[5m]))" } ] },
        { "type": "timeseries", "title": "Error ratio", "gridPos": {"x":8,"y":0,"w":8,"h":8},
          "targets": [ { "expr": "job:http_errors:ratio_rate5m{job=~\"$job\"}" } ] },
        { "type": "timeseries", "title": "p95 latency", "gridPos": {"x":16,"y":0,"w":8,"h":8},
          "targets": [ { "expr": "histogram_quantile(0.95, sum by (job, le) (rate(http_server_request_duration_seconds_bucket{job=~\"$job\"}[5m])))", "exemplar": true } ] }
      ] }
```

Dashboard design that works during incidents:

1. **Top row = the SLO/RED summary.** Is it broken, and since when?
2. Next rows follow the request path: gateway → service → dependencies
   (DB, Redis, queue).
3. Then saturation: CPU throttling, memory vs limit, pool usage, HPA
   replicas vs max.
4. Add deploy annotations (Argo CD sync events) so "it started at the
   deploy" is visible at a glance.

Import community dashboards as a starting point: Node Exporter Full
(1860), the CloudNativePG dashboard, Envoy Gateway's, and Go runtime
metrics.

---

## 32. SLOs and burn-rate alerting

### 32.1 Define the SLOs

| SLI | Definition (good events / valid events) | SLO (30d) |
|---|---|---|
| Orders availability | non-5xx responses to `POST /orders` / all | 99.9% |
| Orders latency | responses < 250 ms / all | 99% |
| Order fulfilment | orders `paid` within 60 s of creation / orders created | 99.5% |

**Error budget** = 1 − SLO. At 99.9% over 30 days that's 43.2 minutes of
full outage, or 0.1% of requests failing.

### 32.2 Multi-window, multi-burn-rate alerts (Google SRE Workbook)

**Burn rate** = how fast you're spending the budget relative to "exactly
use it all in 30 days". A burn rate of 14.4 over 1 hour spends 2% of the
monthly budget in that hour.

| Severity | Long window | Short window | Burn rate | Budget consumed when it fires |
|---|---|---|---|---|
| page | 1h | 5m | 14.4 | 2% |
| page | 6h | 30m | 6 | 5% |
| ticket | 1d | 2h | 3 | 10% |
| ticket | 3d | 6h | 1 | 10% |

The short window makes the alert **reset quickly** once the problem is
fixed. The long window keeps it from firing on a 30-second blip.

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata: { name: orbit-slo-orders, namespace: monitoring }
spec:
  groups:
    - name: slo.orders.availability.recording
      rules:
        - record: slo:orders_errors:ratio_rate5m
          expr: |
            sum(rate(http_server_request_duration_seconds_count{job="orbit/orders",http_response_status_code=~"5.."}[5m]))
            / sum(rate(http_server_request_duration_seconds_count{job="orbit/orders"}[5m]))
        # ... identical rules for 30m, 1h, 2h, 6h, 1d, 3d (generate them; see below)
    - name: slo.orders.availability.alerts
      rules:
        - alert: OrdersErrorBudgetBurnFast
          expr: |
            (slo:orders_errors:ratio_rate1h > (14.4 * 0.001) and slo:orders_errors:ratio_rate5m > (14.4 * 0.001))
            or
            (slo:orders_errors:ratio_rate6h > (6 * 0.001) and slo:orders_errors:ratio_rate30m > (6 * 0.001))
          labels: { severity: page, service: orders }
          annotations:
            summary: "orders is burning its 99.9% availability error budget fast"
            dashboard: https://grafana.acme.example/d/orbit-red?var-job=orbit/orders
            runbook_url: https://runbooks.acme.example/orbit/orders-availability
        - alert: OrdersErrorBudgetBurnSlow
          expr: |
            (slo:orders_errors:ratio_rate1d > (3 * 0.001) and slo:orders_errors:ratio_rate2h > (3 * 0.001))
            or
            (slo:orders_errors:ratio_rate3d > 0.001 and slo:orders_errors:ratio_rate6h > 0.001)
          labels: { severity: ticket, service: orders }
          annotations: { summary: "orders error budget burning steadily" }
```

Latency SLO uses the histogram bucket at the threshold (250 ms = the
`le="0.25"` bucket boundary in the OTel default buckets):

```promql
1 - (
  sum(rate(http_server_request_duration_seconds_bucket{job="orbit/orders",le="0.25"}[1h]))
  / sum(rate(http_server_request_duration_seconds_count{job="orbit/orders"}[1h]))
)   # = fraction of SLOW requests: compare against (1 - 0.99) * burn rate
```

> **Best practice:** don't hand-write these. **Sloth** or **Pyrra** take a
> 10-line SLO spec and generate all recording rules, alerts and a
> dashboard. OpenSLO is the vendor-neutral spec format.

```yaml
# Sloth spec (generates ~30 rules)
version: prometheus/v1
service: orders
slos:
  - name: availability
    objective: 99.9
    sli:
      events:
        error_query: sum(rate(http_server_request_duration_seconds_count{job="orbit/orders",http_response_status_code=~"5.."}[{{.window}}]))
        total_query: sum(rate(http_server_request_duration_seconds_count{job="orbit/orders"}[{{.window}}]))
    alerting:
      page_alert: { labels: { severity: page } }
      ticket_alert: { labels: { severity: ticket } }
```

> ⚠️ **Caveat (verify):** the latency SLI assumes a histogram bucket boundary at exactly 0.25 s (the OTel default HTTP buckets). Confirm with `count by (le) (http_server_request_duration_seconds_bucket)`. Check the Sloth spec against the Sloth version you run.

### 32.3 Error budget policy: the part that changes behavior

An SLO only matters if it drives decisions. Write it down and agree on it
with product owners:

- Budget remaining > 50%: ship freely; run chaos experiments in prod.
- Budget < 25%: releases need extra review; prioritize reliability work.
- Budget exhausted: **feature freeze** for the service except reliability
  fixes, until the 30-day window recovers.

**Part IV checkpoint:** run the k6 test, then break something (scale
Postgres to 1 instance and kill it, or add a `time.Sleep` to orders). Go
from the burn-rate alert → RED dashboard → an exemplar → the Tempo trace
→ that trace's logs in Loki, without typing a pod name.

---

# Part V — Debugging

## 33. The Kubernetes debugging playbook

### 33.1 The universal triage sequence

Whatever the symptom, run these in order. Ninety percent of incidents are
diagnosed by step 3.

```bash
# 1. What state is it in?
kubectl get pods -n orbit -o wide
# 2. What did Kubernetes try, and what went wrong? (read EVENTS at the bottom)
kubectl describe pod <pod> -n orbit
# 3. What did the app say, including the PREVIOUS (crashed) container?
kubectl logs <pod> -n orbit --previous --tail=200
# 4. What does the controller think? (spec vs status, conditions)
kubectl get deploy orders -n orbit -o yaml | yq '.status'
# 5. What happened around it? (events expire after ~1h by default; ship them to Loki for history)
kubectl get events -n orbit --sort-by=.lastTimestamp | tail -40
# 6. What do the metrics say? (Grafana: restarts, memory vs limit, throttling, RED)
kubectl top pods -n orbit --containers
```

> **Best practice:** **mitigate first, debug second.** If a deploy just
> went out and errors rose, roll back (§36.5) *then* investigate. Users
> don't care about your root cause while it's broken. Capture evidence
> before you mitigate if it's cheap: `kubectl logs --previous > file`, a
> heap profile, `describe` output.

### 33.2 Symptom → cause → fix

**`Pending`**

```bash
kubectl describe pod <pod> | sed -n '/Events/,$p'
#   0/4 nodes are available: 3 Insufficient cpu      -> requests too big or cluster full (autoscaler? quota?)
#   ... had untolerated taint {node-role...}         -> needs a toleration / wrong nodeSelector
#   ... didn't match pod affinity/anti-affinity      -> required anti-affinity with too few nodes/zones
#   ... pod has unbound immediate PersistentVolumeClaims -> PVC Pending: check StorageClass/CSI
#   exceeded quota                                   -> ResourceQuota (shows as a ReplicaSet event, not a pod event!)
kubectl get events -n orbit --field-selector reason=FailedCreate     # quota/admission rejections live here
```

**`ImagePullBackOff` / `ErrImagePull`**

```bash
kubectl describe pod <pod> | grep -A3 -i 'failed to pull'
#   not found                    -> typo or tag not pushed (did CI push? which registry?)
#   unauthorized / 403           -> imagePullSecret missing/wrong namespace, or node IAM role
#   no match for platform        -> arm64-only image on amd64 nodes (build multi-arch)
#   toomanyrequests              -> Docker Hub rate limit: mirror images to your own registry
```

**`CrashLoopBackOff`**

```bash
kubectl logs <pod> --previous                                  # the crash reason, nearly always
kubectl get pod <pod> -o jsonpath='{.status.containerStatuses[0].lastState.terminated}' | jq
#   exitCode 1  + config error in logs  -> missing env/secret/bad flag: fix config
#   exitCode 137 + reason OOMKilled      -> raise memory limit / fix leak / set GOMEMLIMIT (§23.4)
#   exitCode 137 without OOMKilled       -> killed by liveness probe failing: check probe events
#   exitCode 0                           -> process exited "successfully": wrong command for a server
#   no logs at all                       -> exec format error (arch), missing binary, or crash before logging
```

Liveness kills look like this in the events: `Liveness probe failed: ...
Container orders failed liveness probe, will be restarted`.

**`CreateContainerConfigError`**: a referenced Secret/ConfigMap or key does
not exist (`describe` names it). Common in GitOps when the ExternalSecret
hasn't synced yet.

**`Running` but `0/1 READY`**

```bash
kubectl describe pod <pod> | grep -i readiness
kubectl port-forward pod/<pod> 8080:8080 & curl -i localhost:8080/readyz   # what does the app say?
```

Our `/readyz` returns JSON listing the failing dependency.

**Service returns connection refused / no response**

```bash
kubectl get endpointslices -n orbit -l kubernetes.io/service-name=orders
#   no endpoints -> selector doesn't match pod labels, OR no pod is Ready
kubectl get svc orders -o yaml | yq '.spec.selector'; kubectl get pods -l app.kubernetes.io/name=orders --show-labels
#   endpoints exist -> port mismatch? targetPort must match containerPort (name or number)
#   still failing -> a NetworkPolicy may be dropping it (§34.3)
```

**DNS failures** (`no such host`, intermittent 5 s delays)

```bash
kubectl run -it --rm dns --image=nicolaka/netshoot --restart=Never -- bash
#   dig orders.orbit.svc.cluster.local     (works?)  dig +search orders
#   cat /etc/resolv.conf                    (ndots, search domains)
kubectl -n kube-system logs -l k8s-app=kube-dns --tail=50
kubectl -n kube-system top pods -l k8s-app=kube-dns       # CoreDNS saturated? scale it / NodeLocal DNSCache
```

**Stuck `Terminating`**

```bash
kubectl get pod <pod> -o jsonpath='{.metadata.finalizers}'      # a controller must remove these
kubectl get node <node>                                         # node NotReady? kubelet can't confirm the kill
# last resort for a pod on a DEAD node (understand the consequences for StatefulSets first):
kubectl delete pod <pod> --grace-period=0 --force
```

Namespaces stuck `Terminating` are almost always a CRD-backed resource
whose controller was uninstalled first, leaving finalizers that nothing
will ever remove. **Uninstall in reverse order:** CRs, then the operator.

**Latency spikes with low average CPU** → CPU throttling (§23.3). Check
the throttled-periods ratio; remove or raise the CPU limit.

**Errors only during deploys** → the termination race (§13.5): add a
preStop sleep, confirm readiness flips to false on SIGTERM, and set
`maxUnavailable: 0`.

**Node `NotReady`**

```bash
kubectl describe node <node> | sed -n '/Conditions/,/Addresses/p'   # MemoryPressure/DiskPressure/PIDPressure
kubectl debug node/<node> -it --image=ubuntu                       # shell on the node (host fs at /host)
#   chroot /host journalctl -u kubelet --since "30 min ago" | tail -100
#   chroot /host crictl ps -a ; crictl logs <id>                   # containers seen by the runtime
#   df -h /host/var/lib/containerd                                 # disk full from images/logs?
```

### 33.3 Ephemeral debug containers (distroless-friendly)

You can't `exec` into a distroless pod: there's no shell. Attach a
temporary toolbox container **into the running pod**, sharing its process
namespace with `--target`:

```bash
kubectl debug -it orders-7d9f-x2k -n orbit --image=nicolaka/netshoot --target=orders
#   ps aux                       -> you see the app's processes
#   ls -l /proc/1/root/          -> the app container's filesystem
#   ss -tnp                      -> its connections (to DB? stuck in CLOSE_WAIT?)
#   curl -s localhost:8080/readyz
#   tcpdump -i any -nn port 5432 -c 50

# Profiles add capabilities: --profile=netadmin (tcpdump/iptables), sysadmin, general, restricted
kubectl debug -it orders-7d9f-x2k --image=nicolaka/netshoot --target=orders --profile=netadmin

# Copy a crash-looping pod with a different command, to poke at it without the app dying:
kubectl debug orders-7d9f-x2k -it --copy-to=orders-debug --container=orders \
  --image=gcr.io/distroless/static-debian12:debug-nonroot -- sh
kubectl delete pod orders-debug    # clean up
```

> **Real-world example:** "We need a shell in the image for debugging" is
> the most common pushback against distroless. Ephemeral containers remove
> it: tools are attached on demand, by people with RBAC permission for the
> `pods/ephemeralcontainers` subresource, and the attach is audited.

### 33.4 Profiling a live pod

Expose `pprof` on a **separate, non-public port** bound to localhost, or
behind auth:

```go
// cmd/orders/main.go
import _ "net/http/pprof" // registers on http.DefaultServeMux

go func() { _ = http.ListenAndServe("127.0.0.1:6060", nil) }() // NOT routed by the Service
```

```bash
kubectl port-forward pod/orders-7d9f-x2k 6060:6060 -n orbit
go tool pprof -http=:9999 'http://localhost:6060/debug/pprof/profile?seconds=30'   # CPU
go tool pprof -http=:9999 http://localhost:6060/debug/pprof/heap                   # memory
curl -s 'http://localhost:6060/debug/pprof/goroutine?debug=2' | head -100          # goroutine leak? deadlock?
```

(`kubectl port-forward` connects via the pod's network namespace, so it
reaches `127.0.0.1:6060` inside the pod.)

**Continuous profiling** (Grafana Pyroscope, Parca, Datadog) samples all
pods all the time at ~1–2% overhead. You can compare flame graphs from
"before the deploy" and "after the deploy" without reproducing anything.

### 33.5 Incident workflow and postmortems

1. **Declare** an incident early, with roles (incident commander, comms,
   operator). It's cheap to downgrade.
2. **Mitigate:** roll back, scale up, fail over, shed load, disable a
   feature flag.
3. **Communicate** on a cadence: status page, stakeholders.
4. **Resolve**, then write a **blameless postmortem** within a week:
   timeline, impact (in SLO terms), contributing factors (plural: there is
   never one root cause), what went well, action items with owners.
   Action items should change the system, not "be more careful".

> **Reference:** the public **Kubernetes failure stories** collection
> (k8s.af) is a catalogue of real postmortems: DNS, OOM, CPU limits,
> etcd, admission webhooks, ingress. Reading ten of them teaches you more
> about production Kubernetes than any certification.

---

# Part VI — Security

## 34. Securing the cluster and the workloads

Defense in depth, in the order an attacker meets it: **supply chain**
(§8, §34.6) → **admission** (what is allowed to run) → **runtime
privileges** (what a compromised pod can do) → **network** (what it can
reach) → **identity** (what it can access) → **detection** (Falco, audit
logs).

### 34.1 ServiceAccounts and RBAC

Every workload gets its **own ServiceAccount**. Mount a token only if it
calls the Kubernetes API.

```yaml
apiVersion: v1
kind: ServiceAccount
metadata: { name: orders, namespace: orbit }
automountServiceAccountToken: false
```

RBAC: Role/ClusterRole (verbs on resources) + RoleBinding/ClusterRoleBinding
(to users, groups, ServiceAccounts). Example: a read-only on-call role for
the `orbit` namespace.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata: { name: oncall-readonly, namespace: orbit }
rules:
  - apiGroups: ["", "apps", "batch", "autoscaling", "keda.sh", "argoproj.io"]
    resources: ["pods", "pods/log", "services", "endpoints", "configmaps", "events",
                "deployments", "replicasets", "statefulsets", "jobs", "cronjobs",
                "horizontalpodautoscalers", "scaledobjects", "rollouts"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["pods/portforward", "pods/ephemeralcontainers"]   # debugging, audited
    verbs: ["create", "patch"]
  # deliberately NO secrets
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata: { name: oncall-readonly, namespace: orbit }
subjects: [{ kind: Group, name: "orbit-oncall", apiGroup: rbac.authorization.k8s.io }]   # from your SSO/OIDC
roleRef: { kind: Role, name: oncall-readonly, apiGroup: rbac.authorization.k8s.io }
```

```bash
kubectl auth can-i list secrets -n orbit --as-group=orbit-oncall --as=jane
kubectl auth can-i --list -n orbit --as=system:serviceaccount:orbit:orders
```

RBAC escalation traps: `create pods` in a namespace means access to every
Secret and ServiceAccount there (mount it into a new pod). `escalate`,
`bind` and `impersonate` verbs, wildcards (`*`), and `nodes/proxy` are
effectively cluster-admin. Audit them with tools like `rbac-tool` or
`kubectl-who-can`.

### 34.2 Pod Security Admission

Built in (replaced PodSecurityPolicy). Label namespaces with a level:
`privileged`, `baseline`, or `restricted`.

```bash
kubectl label namespace orbit \
  pod-security.kubernetes.io/enforce=restricted \
  pod-security.kubernetes.io/enforce-version=latest \
  pod-security.kubernetes.io/warn=restricted \
  pod-security.kubernetes.io/audit=restricted
```

`restricted` requires exactly the securityContext in our golden manifest:
non-root, no privilege escalation, `drop: ALL`, a seccomp profile. Roll it
out with `warn`/`audit` first, read the warnings, then `enforce`.

### 34.3 NetworkPolicy: default deny, then allow what's needed

By default every pod can talk to every pod in every namespace. A
compromised `catalog` could reach Postgres, the Kubernetes API, and the
cloud metadata endpoint. NetworkPolicies need a CNI that enforces them
(Cilium, Calico; recent kind's kindnet does).

```yaml
# 1. default deny everything in the namespace (ingress AND egress)
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata: { name: default-deny, namespace: orbit }
spec:
  podSelector: {}
  policyTypes: [Ingress, Egress]
---
# 2. everyone may resolve DNS
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata: { name: allow-dns, namespace: orbit }
spec:
  podSelector: {}
  policyTypes: [Egress]
  egress:
    - to:
        - namespaceSelector: { matchLabels: { kubernetes.io/metadata.name: kube-system } }
          podSelector: { matchLabels: { k8s-app: kube-dns } }
      ports: [{ protocol: UDP, port: 53 }, { protocol: TCP, port: 53 }]
---
# 3. the gateway may reach catalog and orders
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata: { name: api-from-gateway, namespace: orbit }
spec:
  podSelector:
    matchExpressions: [{ key: app.kubernetes.io/name, operator: In, values: [catalog, orders] }]
  policyTypes: [Ingress]
  ingress:
    - from: [{ namespaceSelector: { matchLabels: { kubernetes.io/metadata.name: envoy-gateway-system } } }]
      ports: [{ port: 8080 }]
---
# 4. Postgres accepts only orders, catalog, worker (+ CNPG operator and its own replicas)
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata: { name: postgres-ingress, namespace: orbit }
spec:
  podSelector: { matchLabels: { cnpg.io/cluster: orbit-db } }
  policyTypes: [Ingress]
  ingress:
    - from:
        - podSelector:
            matchExpressions: [{ key: app.kubernetes.io/name, operator: In, values: [catalog, orders, worker] }]
        - podSelector: { matchLabels: { cnpg.io/cluster: orbit-db } }                   # replication
        - namespaceSelector: { matchLabels: { kubernetes.io/metadata.name: cnpg-system } }  # operator
        - namespaceSelector: { matchLabels: { kubernetes.io/metadata.name: monitoring } }   # metrics :9187
      ports: [{ port: 5432 }, { port: 8000 }, { port: 9187 }]
---
# 5. app egress: Postgres, Redis, and the node-local OTel agent
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata: { name: app-egress, namespace: orbit }
spec:
  podSelector:
    matchExpressions: [{ key: app.kubernetes.io/name, operator: In, values: [catalog, orders, worker] }]
  policyTypes: [Egress]
  egress:
    - to: [{ podSelector: { matchLabels: { cnpg.io/cluster: orbit-db } } }]
      ports: [{ port: 5432 }]
    - to: [{ podSelector: { matchLabels: { app.kubernetes.io/name: redis } } }]
      ports: [{ port: 6379 }]
    - to: [{ ipBlock: { cidr: 172.18.0.0/16 } }]   # node IPs (kind's docker network): OTel agent hostPort
      ports: [{ port: 4317 }]
```

Write the policies (and the remaining Redis ingress / monitoring policies)
**while** you build the service, not afterwards. Then test them:

```bash
kubectl run -n orbit -it --rm probe --image=nicolaka/netshoot --restart=Never -- \
  nc -zv -w2 orbit-db-rw 5432      # expect: timeout (probe pod isn't allowed)
```

Cilium adds L7 (HTTP method/path) and FQDN egress policies
(`toFQDNs: api.stripe.com`). It also adds Hubble, which shows every
dropped flow and why, and turns NetworkPolicy debugging from guesswork
into a query.

> **War story:** pods on cloud VMs can reach the instance metadata
> endpoint (`169.254.169.254`) by default and fetch **the node's IAM
> credentials**. Several cloud breaches started with an SSRF in a web app
> that fetched that URL. Block it with an egress policy (allow
> `0.0.0.0/0` except `169.254.169.254/32`), enforce IMDSv2 with hop limit 1
> on AWS, and use workload identity (§42.3) so pods never need node
> credentials.

### 34.4 Secrets management: External Secrets Operator

The source of truth lives in a real secret manager (AWS Secrets Manager,
GCP Secret Manager, Vault, Azure Key Vault). ESO syncs it into Kubernetes
Secrets and refreshes on rotation. Git holds only *references*.

```bash
helm repo add external-secrets https://charts.external-secrets.io
helm upgrade --install external-secrets external-secrets/external-secrets -n external-secrets --create-namespace
```

```yaml
apiVersion: external-secrets.io/v1
kind: ClusterSecretStore
metadata: { name: aws-secrets-manager }
spec:
  provider:
    aws:
      service: SecretsManager
      region: us-east-1
      auth:   # the ESO pod's own IAM identity (Pod Identity / IRSA): no static keys
        jwt: { serviceAccountRef: { name: external-secrets, namespace: external-secrets } }
---
apiVersion: external-secrets.io/v1
kind: ExternalSecret
metadata: { name: payment-provider, namespace: orbit }
spec:
  refreshInterval: 1h
  secretStoreRef: { kind: ClusterSecretStore, name: aws-secrets-manager }
  target: { name: payment-provider, creationPolicy: Owner }
  data:
    - secretKey: API_KEY
      remoteRef: { key: prod/orbit/payment-provider, property: api_key }
```

Alternatives: **Sealed Secrets** (encrypt with the cluster's public key,
commit the ciphertext; simple, cluster-bound), and **SOPS** + age/KMS
(encrypted YAML in Git, decrypted by Argo CD/Flux plugins). For
database credentials, CNPG already generates and stores them.

> ⚠️ **Caveat (verify):** (1) NetworkPolicy enforcement on kind depends on the kind/kindnet version; recent releases enforce it, older ones silently ignore policies. Test with the `nc` probe *before* trusting the policies, or install Cilium. (2) `172.18.0.0/16` is kind's usual Docker network; confirm with `docker network inspect kind | jq '.[0].IPAM'`. (3) Check CNPG's instance ports (5432, status 8000, metrics 9187) in its docs. (4) External Secrets uses the `external-secrets.io/v1` API on recent releases (older ones use `v1beta1`), and the AWS `auth.jwt` shape differs from Pod Identity setups; check the provider docs.

### 34.5 Admission policy: guardrails as code

**ValidatingAdmissionPolicy** (built in, CEL expressions, GA since 1.30)
needs no extra components:

```yaml
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingAdmissionPolicy
metadata: { name: orbit-workload-standards }
spec:
  failurePolicy: Fail
  matchConstraints:
    resourceRules:
      - { apiGroups: ["apps"], apiVersions: ["v1"], operations: ["CREATE", "UPDATE"], resources: ["deployments"] }
  validations:
    - expression: "object.spec.template.spec.containers.all(c, !c.image.endsWith(':latest') && c.image.contains(':') || c.image.contains('@sha256:'))"
      message: "images must be pinned to a tag or digest, never :latest"
    - expression: "object.spec.template.spec.containers.all(c, has(c.resources.requests) && has(c.resources.limits) && has(c.resources.limits.memory))"
      message: "every container needs requests and a memory limit"
    - expression: "object.spec.template.spec.containers.all(c, has(c.readinessProbe))"
      message: "every container needs a readinessProbe"
---
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingAdmissionPolicyBinding
metadata: { name: orbit-workload-standards }
spec:
  policyName: orbit-workload-standards
  validationActions: [Deny]           # start with [Warn, Audit]
  matchResources:
    namespaceSelector: { matchLabels: { orbit.acme.example/enforce: "true" } }
```

**Kyverno** (YAML policies; validate, *mutate*, *generate*, verify
images) or **OPA Gatekeeper** (Rego) cover what CEL can't, plus a policy
library and reporting.

> **War story:** an admission webhook with `failurePolicy: Fail` whose
> pods were themselves down blocked *every* pod creation in the cluster,
> including its own replacement pods and CoreDNS. Exclude `kube-system`
> and the policy engine's own namespace, run ≥ 3 replicas with a PDB, and
> set webhook timeouts low.

### 34.6 Only run signed images (Kyverno verifyImages)

CI signs images keylessly (§35). This policy rejects any `orbit` image not
signed by *our* workflow on *our* repo:

```yaml
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata: { name: verify-orbit-images }
spec:
  validationFailureAction: Enforce
  webhookTimeoutSeconds: 15
  rules:
    - name: verify-signature
      match: { any: [{ resources: { kinds: [Pod], namespaces: ["orbit*"] } }] }
      verifyImages:
        - imageReferences: ["ghcr.io/acme/orbit-*"]
          mutateDigest: true                 # rewrite tag -> digest: what was verified is what runs
          attestors:
            - entries:
                - keyless:
                    issuer: https://token.actions.githubusercontent.com
                    subject: "https://github.com/acme/orbit/.github/workflows/ci.yaml@refs/heads/main"
                    rekor: { url: https://rekor.sigstore.dev }
```

> ⚠️ **Caveat (verify):** Kyverno is moving image verification to newer policy types (`ImageValidatingPolicy`), and spec-level `validationFailureAction` is deprecated in favor of per-rule `failureAction`. Check this ClusterPolicy against your Kyverno version with `kyverno apply` / `kubectl apply --dry-run=server`. The keyless `subject` must match your workflow path and ref exactly; check it with `cosign verify --certificate-identity ... --certificate-oidc-issuer ...` first.

### 34.7 Runtime detection and continuous audit

- **Falco** (eBPF) alerts on suspicious runtime behavior: a shell spawned
  in a container, writes to `/etc`, unexpected outbound connections,
  reading `/proc/*/environ`.
- **Kubernetes audit logs**: who did what via the API. Enable them on
  managed clusters, ship them to your SIEM, and alert on `exec` into prod
  pods, Secret reads by humans, and RBAC changes.
- **trivy-operator / Kubescape**: continuously scan running workloads
  for CVEs, misconfigurations, and exposed secrets.
- **kube-bench**: CIS Kubernetes Benchmark checks for nodes and the
  control plane.

---

# Part VII — Delivery

## 35. CI: build, test, scan, sign, publish

### 35.1 Pipeline shape

```
PR:     lint → unit tests (-race) → build image → scan → integration test (compose) → manifests validate
main:   ... same ... → push (sha tag) → sign + SBOM attest → bump image digest in deploy/overlays/dev
tag v*: promote: open PRs bumping staging/prod overlays to the SAME digest (build once, promote the artifact)
```

> **Best practice:** **build once, promote the same artifact.** Never
> rebuild for prod. The digest that passed staging is the digest that
> ships. Rebuilding produces a different artifact that has never been
> tested.

### 35.2 GitHub Actions workflow

```yaml
# .github/workflows/ci.yaml
name: ci
on:
  pull_request:
  push: { branches: [main] }

permissions: { contents: read }

concurrency: { group: "${{ github.workflow }}-${{ github.ref }}", cancel-in-progress: true }

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version-file: go.mod }
      - uses: golangci/golangci-lint-action@v8
      - run: go test -race -count=1 ./...
      - name: integration test against the real stack
        run: |
          docker compose up -d --build --wait
          go test -tags=integration ./test/integration/...
          docker compose down -v
      - name: validate manifests
        run: |
          for o in deploy/overlays/*; do
            kustomize build "$o" | kubeconform -strict -summary -ignore-missing-schemas
          done

  image:
    needs: test
    runs-on: ubuntu-latest
    strategy:
      matrix: { service: [catalog, orders, worker] }
    permissions:
      contents: read
      packages: write          # push to GHCR
      id-token: write          # keyless signing (OIDC)
      attestations: write
    outputs:
      digest-catalog: ${{ steps.out.outputs.digest-catalog }}
      digest-orders:  ${{ steps.out.outputs.digest-orders }}
      digest-worker:  ${{ steps.out.outputs.digest-worker }}
    steps:
      - uses: actions/checkout@v4
      - uses: docker/setup-buildx-action@v3
      - uses: docker/login-action@v3
        if: github.event_name == 'push'
        with: { registry: ghcr.io, username: "${{ github.actor }}", password: "${{ secrets.GITHUB_TOKEN }}" }
      - id: meta
        uses: docker/metadata-action@v5
        with:
          images: ghcr.io/acme/orbit-${{ matrix.service }}
          tags: type=sha,format=short,prefix=sha-
      - id: build
        uses: docker/build-push-action@v6
        with:
          context: .
          platforms: linux/amd64,linux/arm64
          build-args: |
            SERVICE=${{ matrix.service }}
            VERSION=${{ github.sha }}
          push: ${{ github.event_name == 'push' }}
          load: false
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          sbom: true
          provenance: mode=max
          cache-from: type=gha,scope=${{ matrix.service }}
          cache-to: type=gha,mode=max,scope=${{ matrix.service }}
      - name: scan (fail on fixable CRITICAL)
        if: github.event_name == 'push'
        uses: aquasecurity/trivy-action@<release-sha>   # pin to the current release by commit SHA
        with:
          image-ref: ghcr.io/acme/orbit-${{ matrix.service }}@${{ steps.build.outputs.digest }}
          severity: CRITICAL
          ignore-unfixed: true
          exit-code: "1"
      - uses: sigstore/cosign-installer@v3
        if: github.event_name == 'push'
      - name: sign (keyless)
        if: github.event_name == 'push'
        run: cosign sign --yes ghcr.io/acme/orbit-${{ matrix.service }}@${{ steps.build.outputs.digest }}
      - id: out
        run: echo "digest-${{ matrix.service }}=${{ steps.build.outputs.digest }}" >> "$GITHUB_OUTPUT"

  deploy-dev:
    if: github.event_name == 'push'
    needs: image
    runs-on: ubuntu-latest
    permissions: { contents: write }
    steps:
      - uses: actions/checkout@v4
      - name: bump dev overlay to the new digests
        run: |
          cd deploy/overlays/dev
          kustomize edit set image \
            orbit-catalog=ghcr.io/acme/orbit-catalog@${{ needs.image.outputs.digest-catalog }} \
            orbit-orders=ghcr.io/acme/orbit-orders@${{ needs.image.outputs.digest-orders }} \
            orbit-worker=ghcr.io/acme/orbit-worker@${{ needs.image.outputs.digest-worker }}
          git config user.name "orbit-ci" && git config user.email "ci@acme.example"
          git commit -am "deploy(dev): ${GITHUB_SHA::7}" && git push
```

Notes:

- `permissions:` are least-privilege per job. `id-token: write` only where
  signing happens.
- Pin third-party actions by **commit SHA** in real repos (tags are
  mutable; compromised action tags have been used in real supply-chain
  attacks, e.g. `tj-actions/changed-files` in 2025). Let Dependabot or
  Renovate bump them.
- Many teams keep deploy manifests in a **separate config repo**, so CI
  needs only write access to that repo, and app-repo commits don't
  trigger deploy loops. Either works; be consistent.
- Add an **ephemeral-cluster e2e job** with `helm/kind-action` for
  platform-level changes (new CRDs, policies).

> ⚠️ **Caveat (verify):** Action major versions (`checkout@v4`, `setup-go@v5`, `build-push-action@v6`, `golangci-lint-action@v8`) were current when this was written. Check each action's releases page and pin to commit SHAs. The trivy-action placeholder must be filled in. `kustomize edit set image name=repo@digest` syntax depends on the kustomize version; run it locally once.

---

## 36. GitOps with Argo CD

### 36.1 Principles (OpenGitOps)

1. **Declarative**: the whole system is described as data.
2. **Versioned and immutable**: Git is the source of truth with full
   history.
3. **Pulled automatically**: an in-cluster agent pulls; CI never holds
   cluster credentials.
4. **Continuously reconciled**: drift is detected and corrected.

Practical wins: every change is a reviewed PR, every rollback is
`git revert`, an audit trail comes for free, and disaster recovery is
"point a new cluster at the repo".

### 36.2 Install and bootstrap

```bash
helm repo add argo https://argoproj.github.io/argo-helm
helm upgrade --install argocd argo/argo-cd -n argocd --create-namespace \
  --set configs.params."server\.insecure"=true --wait
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | base64 -d; echo
kubectl -n argocd port-forward svc/argocd-server 8088:80     # http://localhost:8088
```

### 36.3 App of apps: the platform layer

```
platform/
├── root.yaml                    # the only thing applied by hand
├── apps/
│   ├── cert-manager.yaml        # Argo CD Application -> Helm chart + values
│   ├── envoy-gateway.yaml
│   ├── cnpg.yaml
│   ├── kube-prometheus-stack.yaml
│   ├── loki.yaml  tempo.yaml  otel-agent.yaml  otel-gateway.yaml
│   ├── keda.yaml  argo-rollouts.yaml  kyverno.yaml  external-secrets.yaml  velero.yaml
│   └── orbit.yaml               # ApplicationSet: one Application per environment
└── <component>/values.yaml
```

```yaml
# platform/apps/kube-prometheus-stack.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: kube-prometheus-stack
  namespace: argocd
  annotations: { argocd.argoproj.io/sync-wave: "-5" }     # platform before apps
  finalizers: [resources-finalizer.argocd.argoproj.io]
spec:
  project: platform
  sources:
    - repoURL: https://prometheus-community.github.io/helm-charts
      chart: kube-prometheus-stack
      targetRevision: "77.*"          # pin the current major; Renovate bumps it via PR
      helm: { valueFiles: [$values/platform/kube-prometheus-stack/values.yaml] }
    - repoURL: https://github.com/acme/orbit
      targetRevision: main
      ref: values
  destination: { server: https://kubernetes.default.svc, namespace: monitoring }
  syncPolicy:
    automated: { prune: true, selfHeal: true }
    syncOptions: [CreateNamespace=true, ServerSideApply=true]   # SSA: large CRDs exceed the client-side annotation limit
```

```yaml
# platform/apps/orbit.yaml: one Application per environment overlay
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata: { name: orbit, namespace: argocd }
spec:
  goTemplate: true
  generators:
    - list:
        elements:
          - { env: dev,     autoSync: true }
          - { env: staging, autoSync: true }
          - { env: prod,    autoSync: true }   # still gated: changes reach prod only via a reviewed PR
  template:
    metadata: { name: "orbit-{{.env}}" }
    spec:
      project: orbit
      source:
        repoURL: https://github.com/acme/orbit
        targetRevision: main
        path: "deploy/overlays/{{.env}}"
      destination: { server: https://kubernetes.default.svc, namespace: "orbit-{{.env}}" }
      syncPolicy:
        automated: { prune: true, selfHeal: true }
        syncOptions: [CreateNamespace=true]
        retry: { limit: 5, backoff: { duration: 10s, factor: 2, maxDuration: 3m } }
```

```bash
kubectl apply -n argocd -f platform/root.yaml     # root Application pointing at platform/apps/
argocd app list
argocd app diff orbit-prod
argocd app history orbit-prod
```

(For real multi-cluster setups, use a **cluster generator** instead of a
list, so each environment's cluster is registered in Argo CD and its
Applications target it.)

### 36.4 Sync waves and hooks: ordering

```yaml
metadata:
  annotations:
    argocd.argoproj.io/sync-wave: "-1"       # lower waves first: CRDs/namespaces < DB < migrations < apps
    argocd.argoproj.io/hook: PreSync          # run before the main sync (migrations, §38)
    argocd.argoproj.io/hook-delete-policy: BeforeHookCreation
```

### 36.5 Promotion and rollback

- **Promotion** = a PR that copies the digest from the `staging` overlay
  to the `prod` overlay. Automate opening the PR; keep merging human (or
  gated on the staging SLO). Tools like **Kargo** or Argo CD Image Updater
  automate multi-stage promotion.
- **Rollback** = `git revert` the promotion commit. Argo CD syncs the
  previous digest. With canaries (§37), a bad release usually aborts
  itself before anyone has to.
- **Drift**: `selfHeal` reverts manual `kubectl edit`s. For a break-glass
  manual change, disable auto-sync on that Application first, and record
  why.

> **War story:** with `prune: true`, removing a file from Git deletes the
> resource. Someone refactored Kustomize paths, a PVC-owning StatefulSet
> briefly vanished from the rendered output, and Argo CD pruned it. Protect
> stateful resources with
> `argocd.argoproj.io/sync-options: Prune=false,Delete=false`, use
> `reclaimPolicy: Retain` on the StorageClass, and review Argo CD diffs in
> PRs (the `argocd app diff` output, or a rendered-manifests diff check in
> CI).

> ⚠️ **Caveat (verify):** the kube-prometheus-stack major in `targetRevision` is a placeholder for whatever is current. Look it up with `helm search repo prometheus-community/kube-prometheus-stack`. Check the ApplicationSet template fields against your Argo CD version.

---

## 37. Progressive delivery: canary with automatic rollback

Rolling updates are gated only on readiness, which says the pod works,
not that the new code is correct. A **canary** sends a small share of real
traffic to the new version, measures it against the SLO, and proceeds or
aborts automatically.

### 37.1 Install Argo Rollouts with the Gateway API plugin

```bash
helm upgrade --install argo-rollouts argo/argo-rollouts -n argo-rollouts --create-namespace \
  --set dashboard.enabled=true
```

```yaml
# ConfigMap argo-rollouts-config in argo-rollouts: enable the Gateway API traffic router plugin
apiVersion: v1
kind: ConfigMap
metadata: { name: argo-rollouts-config, namespace: argo-rollouts }
data:
  trafficRouterPlugins: |
    - name: "argoproj-labs/gatewayAPI"
      location: "https://github.com/argoproj-labs/rollouts-plugin-trafficrouter-gatewayapi/releases/download/<latest-version>/gatewayapi-plugin-linux-amd64"   # pin a release
```

(The controller also needs RBAC to patch HTTPRoutes; the plugin's README
has the ClusterRole.)

### 37.2 Convert `orders` to a Rollout

Replace the `Deployment` with a `Rollout`. The pod template is identical,
so the golden manifest carries over unchanged. Add a second Service for
the canary and reference both from the HTTPRoute (`orders` weight 100,
`orders-canary` weight 0).

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Rollout
metadata: { name: orders, namespace: orbit }
spec:
  revisionHistoryLimit: 5
  selector: { matchLabels: { app.kubernetes.io/name: orders } }
  template:
    # ... exactly the pod template from §24.3, plus:
    #   env: APP_VERSION from fieldRef metadata.labels['app.kubernetes.io/version']
  strategy:
    canary:
      stableService: orders
      canaryService: orders-canary
      trafficRouting:
        plugins:
          argoproj-labs/gatewayAPI:
            httpRoute: orbit-api
            namespace: orbit
      analysis:                                   # background analysis during the whole rollout
        templates: [{ templateName: orders-slo }]
        startingStep: 1
        args:
          - name: canary-version
            valueFrom: { podTemplateHashValue: Latest }
      steps:
        - setWeight: 5
        - pause: { duration: 2m }
        - setWeight: 25
        - pause: { duration: 5m }
        - setWeight: 50
        - pause: { duration: 5m }
        # 100% after the last step
```

### 37.3 The analysis: abort if the canary breaks the SLO

```yaml
apiVersion: argoproj.io/v1alpha1
kind: AnalysisTemplate
metadata: { name: orders-slo, namespace: orbit }
spec:
  args: [{ name: canary-version }]
  metrics:
    - name: error-ratio
      interval: 1m
      failureLimit: 1
      successCondition: "len(result) == 0 || isNaN(result[0]) || result[0] < 0.01"   # no data yet = not a failure
      provider:
        prometheus:
          address: http://kps-prometheus.monitoring.svc:9090
          query: |
            sum(rate(http_server_request_duration_seconds_count{job="orbit/orders",k8s_pod_name=~"orders-{{args.canary-version}}-.*",http_response_status_code=~"5.."}[2m]))
            /
            sum(rate(http_server_request_duration_seconds_count{job="orbit/orders",k8s_pod_name=~"orders-{{args.canary-version}}-.*"}[2m]))
    - name: p95-latency
      interval: 1m
      failureLimit: 1
      successCondition: "len(result) == 0 || isNaN(result[0]) || result[0] < 0.3"
      provider:
        prometheus:
          address: http://kps-prometheus.monitoring.svc:9090
          query: |
            histogram_quantile(0.95, sum by (le) (rate(http_server_request_duration_seconds_bucket{job="orbit/orders",k8s_pod_name=~"orders-{{args.canary-version}}-.*"}[2m])))
```

(Canary pods are named `orders-<pod-template-hash>-xxxxx`, and
`k8s.pod.name` was promoted to a label in §30.1. That's how the query
isolates the canary from the stable version.)

```bash
kubectl argo rollouts get rollout orders -n orbit --watch
kubectl argo rollouts promote orders -n orbit        # skip the current pause
kubectl argo rollouts abort orders -n orbit          # manual abort -> traffic back to stable
```

Try it: release a version of `orders` that returns 500 for 5% of requests
while k6 runs. The analysis fails at the 5% or 25% step, traffic snaps
back to stable, and **at most a few percent of a few minutes of traffic**
saw the bug. Compare that with a rolling update, which would have reached
100% in under a minute.

> ⚠️ **Caveat (verify):** (1) Fill in the Gateway API plugin version and add the ClusterRole from the plugin README; without it the controller can't patch HTTPRoutes. On arm64, use the `-linux-arm64` binary. (2) The analysis queries depend on `k8s.pod.name` being promoted to the `k8s_pod_name` label (§30.1 caveat) and on canary pod names containing the pod-template hash. Check the actual names with `kubectl get pods`. (3) `isNaN` / `len(result)` come from the expr language Argo Rollouts uses; test the template with a dry run on a non-critical service first.

### 37.4 Other strategies

| Strategy | How | Good for | Cost |
|---|---|---|---|
| Rolling update | replace pods gradually | most changes | no traffic control |
| Canary | % of traffic, metric-gated | user-facing APIs | needs good metrics + a traffic router |
| Blue/green | full new stack, switch Service selector at once, instant switch back | big-bang changes, quick rollback | 2× capacity during deploy |
| Feature flags | ship code dark, enable per user/segment at runtime (OpenFeature, LaunchDarkly, Unleash) | decoupling deploy from release | flag debt: remove old flags |
| Shadow/mirror | copy live traffic to the new version, discard responses | validating rewrites | side effects must be disabled |

---

## 38. Database migrations in a continuously deployed system

During a rolling update or canary, **old and new code run against the same
database at the same time**. And a rollback means old code runs against
the *new* schema. So every migration must be compatible with both
versions: the **expand/contract** (parallel change) pattern.

Renaming `orders.status` to `orders.state`, done safely across releases:

| Release | Schema change | Code |
|---|---|---|
| 1 (expand) | add `state` (nullable) | write both columns; read `status` |
| 2 | backfill `state` in batches (a Job) | read `state`, fall back to `status` |
| 3 | add NOT NULL / constraints | read and write `state` only |
| 4 (contract) | drop `status`, **only after release 3 is proven and can't be rolled back past** | — |

Running migrations in GitOps:

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  generateName: orbit-migrate-
  annotations:
    argocd.argoproj.io/hook: PreSync                     # before the new app version rolls out
    argocd.argoproj.io/hook-delete-policy: BeforeHookCreation
spec:
  backoffLimit: 2
  activeDeadlineSeconds: 900
  template:
    spec:
      restartPolicy: Never
      containers:
        - name: migrate
          image: ghcr.io/acme/orbit-migrations@sha256:...   # versioned alongside the app images
          args: ["-dir", "/migrations", "postgres", "$(DATABASE_URL)", "up"]
          env:
            - { name: DATABASE_URL, valueFrom: { secretKeyRef: { name: orbit-db-app, key: uri } } }
```

Migration rules for production Postgres:

- `CREATE INDEX CONCURRENTLY` (no table lock; it can't run inside a
  transaction, so mark it `-- +goose NO TRANSACTION`).
- Add columns as nullable or with a constant default (instant in PG 11+).
  Add `NOT NULL` via a `CHECK ... NOT VALID` constraint, then `VALIDATE`.
- Set `lock_timeout` (e.g. `SET lock_timeout = '5s'`) so a migration
  waiting on a lock fails fast instead of queueing every query behind it.
- Large backfills run as batched Jobs, never inside the migration.
- **Never run migrations from app startup** with N replicas. They race;
  goose and Flyway use locks, but a slow migration then blocks every pod's
  startup and readiness.

> **War story:** `ALTER TABLE orders ADD COLUMN ... DEFAULT now()` on a
> 400M-row table (a *volatile* default forces a table rewrite even on
> modern Postgres) took an `ACCESS EXCLUSIVE` lock for 40 minutes. Every
> order request queued behind it, connection pools filled, and readiness
> probes failed fleet-wide. `lock_timeout` plus a review checklist (or a
> linter like `squawk`) in CI prevents this class of outage.

---

# Part VIII — Reliability Operations

## 39. Backup and disaster recovery

### 39.1 What actually needs backing up

| State | Where it lives | How it's recovered |
|---|---|---|
| Cluster configuration (all manifests, Helm values, policies) | **Git** | re-bootstrap Argo CD; it rebuilds everything |
| Application data | Postgres (CNPG), PVs | **database-native backups + WAL archiving** (PITR) |
| Secrets | cloud secret manager (ESO) | the secret manager's own replication; ESO re-syncs |
| Cluster objects not in Git (CR status, generated certs) | etcd | Velero (namespace-level), or regenerate |
| etcd itself | control plane | managed: provider's job; self-managed: `etcdctl snapshot` |
| Container images | registry | registry replication / immutable tags |

> **Best practice:** if your cluster can be rebuilt from Git, then
> **cluster DR = data DR**. Treat clusters as cattle: "restore the
> cluster" becomes "create a new cluster, point Argo CD at the repo,
> restore the databases".

Define targets per service before choosing tools. **RPO**: how much data
can we lose? **RTO**: how long can we be down? Orbit's targets: RPO 5 min,
RTO 1 h.

### 39.2 Postgres: continuous backup and point-in-time recovery (CNPG)

CloudNativePG ships base backups plus continuous **WAL archiving** to
object storage through the Barman Cloud plugin. The RPO is the WAL archive
interval, which is seconds to minutes.

```yaml
apiVersion: barmancloud.cnpg.io/v1
kind: ObjectStore
metadata: { name: orbit-backups, namespace: orbit }
spec:
  retentionPolicy: "30d"
  configuration:
    destinationPath: s3://acme-orbit-pg-backups/       # lab: a MinIO bucket
    s3Credentials: { inheritFromIAMRole: true }         # Pod Identity / IRSA in the cloud
    wal: { compression: gzip }
    data: { compression: gzip }
---
# added to the Cluster spec (§20.1)
#   plugins:
#     - name: barman-cloud.cloudnative-pg.io
#       isWALArchiver: true
#       parameters: { barmanObjectName: orbit-backups }
---
apiVersion: postgresql.cnpg.io/v1
kind: ScheduledBackup
metadata: { name: orbit-db-daily, namespace: orbit }
spec:
  schedule: "0 0 2 * * *"           # six fields: sec min hour dom mon dow
  backupOwnerReference: self
  cluster: { name: orbit-db }
  method: plugin
  pluginConfiguration: { name: barman-cloud.cloudnative-pg.io }
```

Point-in-time restore into a **new** cluster (never over the broken one):

```yaml
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata: { name: orbit-db-restore, namespace: orbit }
spec:
  instances: 3
  storage: { size: 2Gi }
  bootstrap:
    recovery:
      source: orbit-db
      recoveryTarget: { targetTime: "2026-10-04 09:41:00+00" }   # just before the bad DELETE
  externalClusters:
    - name: orbit-db
      plugin:
        name: barman-cloud.cloudnative-pg.io
        parameters: { barmanObjectName: orbit-backups, serverName: orbit-db }
```

> **War story:** GitLab's 2017 outage. An engineer ran `rm -rf` on the
> primary's data directory during a replication incident. Then the team
> found that **none of their five backup mechanisms worked**: pg_dump
> had been failing silently (version mismatch), and other mechanisms had
> never been tested. They restored from a 6-hour-old staging copy. Lesson
> burned into every SRE team since: **a backup you haven't restored is
> not a backup.** Automate restore tests (§39.4).

> ⚠️ **Caveat (verify):** the Barman Cloud plugin is a **separate install** (it needs cert-manager). Its API group, ObjectStore fields, and the Cluster `plugins:` / recovery `externalClusters.plugin` shapes follow the plugin docs for recent CNPG (1.26+). Older CNPG versions use the in-tree `spec.backup.barmanObjectStore` instead. Check against the CNPG and plugin docs for your versions, and use a MinIO bucket with explicit credentials in the lab (no IAM role on kind).

### 39.3 Velero: Kubernetes objects and volumes

```bash
# lab: MinIO as the S3 target
helm repo add vmware-tanzu https://vmware-tanzu.github.io/helm-charts
helm upgrade --install velero vmware-tanzu/velero -n velero --create-namespace -f platform/velero/values.yaml

velero schedule create orbit-hourly --schedule="@every 1h" --include-namespaces orbit --ttl 168h
velero backup create orbit-manual --include-namespaces orbit --wait
velero backup describe orbit-manual --details
velero restore create --from-backup orbit-manual --namespace-mappings orbit:orbit-restore-test
```

For databases, prefer the DB-native backup (consistent, PITR-capable) and
use Velero for the surrounding objects and for simple PVs (file-system
backup via Kopia, or CSI snapshots).

> ⚠️ **Caveat (verify):** the Velero values file isn't shown. It must configure the AWS object-store plugin (`velero/velero-plugin-for-aws`), the MinIO URL with `s3ForcePathStyle`, and credentials. Follow the Velero chart README for your chart version.

### 39.4 Prove it: an automated restore test

A weekly CronJob (or CI job) that restores the latest backup into a
scratch namespace, runs sanity queries, and reports success as a metric:

```bash
#!/usr/bin/env bash
# restore-test.sh: runs weekly; alert if orbit_restore_test_last_success_timestamp is older than 8 days
set -euo pipefail
ns=restore-test-$(date +%s)
kubectl create ns "$ns"
sed "s/namespace: orbit/namespace: $ns/" restore-cluster.yaml | kubectl apply -f -
kubectl wait -n "$ns" cluster/orbit-db-restore --for=condition=Ready --timeout=30m
kubectl exec -n "$ns" orbit-db-restore-1 -- psql -U postgres -d orbit -tAc \
  "select count(*) > 0 from orders where created_at > now() - interval '1 day'" | grep -q t
curl -s --data-binary "orbit_restore_test_last_success_timestamp $(date +%s)" \
  http://pushgateway.monitoring:9091/metrics/job/restore_test
kubectl delete ns "$ns"
```

### 39.5 DR strategies and the full-cluster drill

| Strategy | RTO | RPO | Cost | Description |
|---|---|---|---|---|
| Backup & restore | hours | minutes (WAL) | $ | rebuild from Git + backups in another region |
| Pilot light | ~1 h | minutes | $$ | DB replica running in region B, cluster scaled to ~0 |
| Warm standby | minutes | seconds | $$$ | smaller full stack in B, scaled up on failover |
| Active-active | ~0 | ~0 | $$$$ | both regions serve; needs multi-region data design |

The drill (do it on kind; it maps 1:1 to the cloud):

1. `kind delete cluster --name orbit`. The "region" is gone.
2. `kind create cluster --config kind-orbit.yaml`, install Argo CD, apply
   `platform/root.yaml`.
3. Argo CD rebuilds the platform and apps. The DB Cluster bootstraps from
   backup (a `recovery` bootstrap in a DR overlay).
4. Verify with the k6 smoke test plus data checks. **Write down the
   wall-clock time: that is your real RTO.**

---

## 40. Chaos engineering

Chaos engineering is **experiments** that test a hypothesis about
resilience. It isn't random breakage. Format: *"Given steady state X
(SLO met), when Y fails, we expect Z, with users seeing at most W."*

### 40.1 Chaos Mesh

```bash
helm repo add chaos-mesh https://charts.chaos-mesh.org
helm upgrade --install chaos-mesh chaos-mesh/chaos-mesh -n chaos-mesh --create-namespace \
  --set chaosDaemon.runtime=containerd --set chaosDaemon.socketPath=/run/containerd/containerd.sock
```

```yaml
# H1: killing 1/3 of orders pods every 2 minutes causes no failed requests (PDB, readiness, retries)
apiVersion: chaos-mesh.org/v1alpha1
kind: Schedule
metadata: { name: orders-pod-kill, namespace: orbit }
spec:
  schedule: "@every 2m"
  type: PodChaos
  historyLimit: 5
  concurrencyPolicy: Forbid
  podChaos:
    action: pod-kill
    mode: fixed-percent
    value: "33"
    selector: { namespaces: [orbit], labelSelectors: { app.kubernetes.io/name: orders } }
---
# H2: 200ms extra latency to Postgres keeps orders p95 under 300 ms? (tests timeouts/pool sizing)
apiVersion: chaos-mesh.org/v1alpha1
kind: NetworkChaos
metadata: { name: db-latency, namespace: orbit }
spec:
  action: delay
  mode: all
  selector: { namespaces: [orbit], labelSelectors: { app.kubernetes.io/name: orders } }
  direction: to
  target:
    mode: all
    selector: { namespaces: [orbit], labelSelectors: { cnpg.io/cluster: orbit-db } }
  delay: { latency: "200ms", jitter: "50ms" }
  duration: "10m"
---
# H3: worker memory pressure -> OOMKill -> unacked messages are reclaimed, nothing lost
apiVersion: chaos-mesh.org/v1alpha1
kind: StressChaos
metadata: { name: worker-mem, namespace: orbit }
spec:
  mode: one
  selector: { namespaces: [orbit], labelSelectors: { app.kubernetes.io/name: worker } }
  stressors: { memory: { workers: 1, size: "300MB" } }
  duration: "5m"
```

> ⚠️ **Caveat (verify):** the Chaos Mesh daemon's containerd socket path is correct for kind nodes. Other distros differ (k3s: `/run/k3s/containerd/containerd.sock`). Check field names against the installed CRDs with `kubectl explain podchaos.spec`.

### 40.2 Game day template

| Step | Example |
|---|---|
| Hypothesis | "Losing zone-b causes < 1 min of elevated errors and no data loss" |
| Steady state | k6 at 200 RPS, SLO dashboards green |
| Blast radius / abort criteria | staging only; abort if error ratio > 5% for 2 min |
| Inject | `kubectl cordon` + delete all pods on zone-b nodes (or `docker stop` the kind node) |
| Observe | alerts fired? which? dashboards tell the story? runbook accurate? |
| Results | what broke that we didn't expect (there is always something) |
| Actions | tickets with owners; re-run the experiment after fixes |

> **Real-world example:** Netflix's Chaos Monkey (2011) randomly
> terminated production instances during business hours, so engineers had
> to build services that tolerated it. The lasting lesson isn't "break
> prod". It's that **failure handling which is never exercised doesn't
> work** when you need it. Start in staging, during work hours, with an
> abort button.

---

## 41. Cost and capacity

Kubernetes cost is mostly **requested** capacity (you pay for nodes sized
to requests, not to usage) plus data transfer and storage.

| Lever | Typical saving | How |
|---|---|---|
| Right-size requests | 30–50% | VPA recommendations, p95 usage dashboards (§23.6) |
| Spot/preemptible for stateless + workers | 60–90% on those nodes | Karpenter `capacity-type: spot`, PDBs, ≥ 2 replicas, graceful shutdown |
| arm64 nodes | 20–40% | multi-arch images (§10.2) |
| Bin-packing / consolidation | 10–30% | Karpenter consolidation, fewer node pools |
| Scale to zero | ~100% off-hours for dev | KEDA `cron` trigger, or sleep schedules for dev namespaces |
| Cross-AZ traffic | often a surprise line item | topology-aware routing (`trafficDistribution: PreferClose` on Services), co-locate chatty pairs |
| Observability bills | often 10–30% of infra | drop unused metrics (`metric_relabel_configs`), sample traces, don't log every request |
| NAT gateway data | per-GB charge | VPC endpoints for S3/ECR/STS, pull-through cache |

Make cost visible per team and namespace with **OpenCost / Kubecost**,
using `team` and `cost-center` labels (§19.2).

---

## 42. Going to the cloud: Orbit on EKS

Everything so far runs on kind. In the cloud, only the **bottom layer**
changes: cluster, network, load balancers, storage classes, identity. The
GitOps repo, add-ons and app manifests stay the same, apart from an
overlay.

### 42.1 Infrastructure as code with Terraform

```
infra/
├── modules/                   # (or use the community modules directly)
├── envs/
│   ├── staging/main.tf
│   └── prod/main.tf           # separate state + separate AWS account per env
└── bootstrap/                 # state bucket, lock table, OIDC role for CI
```

```hcl
# infra/envs/prod/main.tf  (sketch: check each module's docs for current input names)
terraform {
  required_version = ">= 1.10"
  backend "s3" {
    bucket       = "acme-tfstate-prod"
    key          = "orbit/eks.tfstate"
    region       = "us-east-1"
    use_lockfile = true
  }
  required_providers { aws = { source = "hashicorp/aws", version = "~> 6.0" } }
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 6.0"
  name    = "orbit-prod"
  cidr    = "10.40.0.0/16"
  azs             = ["us-east-1a", "us-east-1b", "us-east-1c"]
  private_subnets = ["10.40.0.0/19", "10.40.32.0/19", "10.40.64.0/19"]
  public_subnets  = ["10.40.96.0/22", "10.40.100.0/22", "10.40.104.0/22"]
  enable_nat_gateway     = true
  one_nat_gateway_per_az = true                    # a single NAT is a single point of failure
  private_subnet_tags = { "karpenter.sh/discovery" = "orbit-prod", "kubernetes.io/role/internal-elb" = 1 }
  public_subnet_tags  = { "kubernetes.io/role/elb" = 1 }
}

module "eks" {
  source             = "terraform-aws-modules/eks/aws"
  version            = "~> 21.0"
  name               = "orbit-prod"
  kubernetes_version = "1.34"
  vpc_id             = module.vpc.vpc_id
  subnet_ids         = module.vpc.private_subnets
  endpoint_public_access = false                   # API reachable via VPN/bastion only
  enable_cluster_creator_admin_permissions = true

  addons = {
    coredns                = {}
    kube-proxy             = {}
    vpc-cni                = { before_compute = true }
    eks-pod-identity-agent = { before_compute = true }
    aws-ebs-csi-driver     = {}
  }

  # small, stable on-demand group for system add-ons and Karpenter itself; Karpenter does the rest
  eks_managed_node_groups = {
    system = {
      instance_types = ["m7g.large"]
      ami_type       = "AL2023_ARM_64_STANDARD"
      min_size = 3, max_size = 4, desired_size = 3
      labels = { "node-role" = "system" }
      taints = { system = { key = "CriticalAddonsOnly", value = "true", effect = "NO_SCHEDULE" } }
    }
  }
}
```

```bash
cd infra/envs/prod && terraform init && terraform plan -out tf.plan && terraform apply tf.plan
aws eks update-kubeconfig --name orbit-prod --region us-east-1
```

> ⚠️ **Caveat (verify):** the Terraform module majors (`vpc ~> 6`, `eks ~> 21`, AWS provider `~> 6`) and v21 input names (`name`, `kubernetes_version`, `addons`, `endpoint_public_access`) were written from memory of the v21 changes. Run `terraform init && terraform validate` and compare with each module's README and UPGRADE notes. Check that `kubernetes_version` is a version EKS currently offers (`aws eks describe-cluster-versions`). The S3 backend's `use_lockfile` needs Terraform 1.10+.

### 42.2 What changes in the GitOps repo for EKS

| Concern | kind | EKS |
|---|---|---|
| LoadBalancer | cloud-provider-kind | AWS Load Balancer Controller (NLB in front of Envoy Gateway) |
| StorageClass | local-path | `gp3` via EBS CSI, `WaitForFirstConsumer`, encrypted |
| Node scaling | fixed 3 workers | Karpenter NodePools (§27.6) |
| Postgres | CNPG in-cluster | **RDS/Aurora** (managed) *or* CNPG on EBS with S3 backups. Choose by team skill and needs |
| Redis | StatefulSet | ElastiCache / MemoryDB |
| Secrets | K8s Secrets | AWS Secrets Manager + ESO |
| Pod → AWS access | n/a | **EKS Pod Identity** (or IRSA) |
| DNS + certs | localtest.me + self-signed | external-dns (Route 53) + cert-manager with Let's Encrypt DNS-01 |

### 42.3 Workload identity: no static cloud credentials, ever

```hcl
# IAM role the orders ServiceAccount assumes, scoped to one bucket
resource "aws_iam_role" "orders" {
  name               = "orbit-prod-orders"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{ Effect = "Allow", Principal = { Service = "pods.eks.amazonaws.com" },
                   Action = ["sts:AssumeRole", "sts:TagSession"] }]
  })
}
resource "aws_iam_role_policy" "orders_s3" {
  role   = aws_iam_role.orders.id
  policy = jsonencode({ Version = "2012-10-17", Statement = [{ Effect = "Allow",
    Action = ["s3:PutObject"], Resource = "arn:aws:s3:::acme-orbit-invoices/*" }] })
}
resource "aws_eks_pod_identity_association" "orders" {
  cluster_name    = module.eks.cluster_name
  namespace       = "orbit-prod"
  service_account = "orders"
  role_arn        = aws_iam_role.orders.arn
}
```

The AWS SDK in the pod picks up short-lived, auto-rotated credentials. No
access keys in Secrets, nothing to leak. The equivalents are **GKE
Workload Identity Federation** and **AKS Workload Identity**.

### 42.4 Cluster upgrades without drama

Kubernetes ships 3 minors a year, and managed providers support each for
roughly 14 months (then extended support costs extra). Upgrade at least
twice a year:

1. Read the changelog's deprecations. Run `pluto detect-all-in-cluster`
   / `kubent` for removed APIs.
2. Upgrade add-ons to versions that support the target minor (CNI, CSI,
   CoreDNS, Karpenter, ingress/gateway, operators).
3. Upgrade the **control plane** one minor at a time (1.33 → 1.34).
4. Roll the **nodes**: Karpenter drift replaces them automatically and
   respects PDBs. This is where bad PDBs and single replicas bite (§22.3).
5. For high-risk jumps, use **blue/green clusters**: build a new cluster
   at the new version, sync it from Git, shift traffic via DNS/LB weights,
   then delete the old one. GitOps is what makes this cheap.

> ⚠️ **Caveat (verify):** provider support windows and pricing for extended support change; check the EKS/GKE/AKS version calendars.

### 42.5 GKE and AKS equivalents

| Concept | EKS | GKE | AKS |
|---|---|---|---|
| Node autoscaling | Karpenter / CA | node auto-provisioning, **Autopilot** (no nodes to manage) | Node Auto Provisioning (Karpenter-based) / CA |
| Workload identity | Pod Identity / IRSA | Workload Identity Federation | Entra Workload ID |
| Managed Postgres | RDS / Aurora | Cloud SQL / AlloyDB | Azure Database for PostgreSQL |
| Gateway API | AWS LB Controller, or Envoy/Cilium | built-in GKE Gateway controller | Application Gateway for Containers |
| Managed Prometheus | Amazon Managed Prometheus | Google Managed Prometheus | Azure Monitor managed Prometheus |
| Registry | ECR | Artifact Registry | ACR |

---

# Part IX — Capstone and Reference

## 43. Capstone: assemble and verify the whole system

### 43.1 Bootstrap order (kind)

```makefile
# Makefile (excerpt)
cluster:        ## 1. cluster + LB
	kind create cluster --config kind-orbit.yaml
images:         ## 2. build + load images (CI pushes to GHCR in the real flow)
	for s in catalog orders worker; do \
	  docker build --build-arg SERVICE=$$s -t orbit-$$s:dev . && kind load docker-image orbit-$$s:dev --name orbit; \
	done
gitops:         ## 3. Argo CD + root app: everything else is pulled from Git
	helm upgrade --install argocd argo/argo-cd -n argocd --create-namespace --wait
	kubectl apply -n argocd -f platform/root.yaml
verify:         ## 4. smoke + load + chaos
	k6 run -e BASE=https://orbit.localtest.me load/smoke.js
	k6 run load/orders.js
```

Sync waves (from §36.4) give the platform its order: CRDs and operators
(cert-manager, CNPG, KEDA, Rollouts, Kyverno, ESO) → gateways, monitoring
and Collectors → policies → the `orbit-*` ApplicationSet (DB → migration
hook → services → routes → ScaledObjects).

> **macOS note:** if the Gateway's LoadBalancer address isn't reachable
> from the host, port-forward the Envoy service instead:
> `kubectl -n envoy-gateway-system port-forward svc/$(kubectl -n envoy-gateway-system get svc -l gateway.envoyproxy.io/owning-gateway-name=public -o name | cut -d/ -f2) 8443:443`
> and use `https://orbit.localtest.me:8443`.
>
> **Lab VM note (§0):** inside the VM the LoadBalancer IP works directly
> (`cloud-provider-kind` runs as a systemd service). To use the browser on
> your Mac, run the same port-forward *inside the VM*: Lima forwards
> `8443` to your Mac's `localhost` automatically (§0.5).

### 43.2 Acceptance tests: the system is "done" when all of these pass

| # | Test | Proves |
|---|---|---|
| 1 | `git push` to main → new digest running in dev in < 10 min, no human steps | CI + GitOps |
| 2 | Unsigned image applied by hand → rejected by admission | supply chain policy |
| 3 | `kubectl run` with root user in `orbit` → rejected | Pod Security |
| 4 | k6 800 RPS: p95 < 300 ms, errors < 1%; orders and worker replicas scale out, then back in (worker to 0) | autoscaling |
| 5 | Rolling restart of all services during k6 → 0 failed requests | probes, preStop, PDB, `maxUnavailable: 0` |
| 6 | Kill the Postgres primary → failover < 30 s; orders not-ready, then ready; **no pod restarts** | liveness/readiness design, CNPG |
| 7 | Bad canary (5% 500s) → automatic abort at ≤ 25% weight | progressive delivery |
| 8 | Burn-rate alert fires for #7 if you force-promote; Slack message has a runbook link | SLO alerting |
| 9 | From that alert: dashboard → exemplar → trace → logs in < 5 min | observability correlation |
| 10 | Drain a node / stop a kind "zone" during k6 → SLO holds | topology spread, tolerations, PDB |
| 11 | Delete the `orbit` namespace → restore DB (PITR) + Argo CD resync; measure RTO; verify orders exist | backup & DR |
| 12 | Catalog pod cannot open a connection to Postgres admin port / metadata IP | NetworkPolicy |
| 13 | NetworkChaos 200 ms to DB → p95 degrades but error rate stays < 1%, alerts are tickets, not pages | timeouts, alert design |
| 14 | A poison message (unknown `order_id` that fails forever) → moved to `orders.dlq` after 5 attempts; alert on DLQ length > 0 | queue hygiene |

### 43.3 The dead-letter queue (test #14)

```go
// in worker.handle, on failure: check delivery count from XPENDING; after maxAttempts, park it
const maxAttempts = 5

func (w *Worker) deadLetterIfExhausted(ctx context.Context, m redis.XMessage, cause error) bool {
	p, err := w.RDB.XPendingExt(ctx, &redis.XPendingExtArgs{
		Stream: w.Stream, Group: w.Group, Start: m.ID, End: m.ID, Count: 1,
	}).Result()
	if err != nil || len(p) == 0 || p[0].RetryCount < maxAttempts {
		return false
	}
	values := map[string]any{"original_id": m.ID, "error": cause.Error(), "attempts": p[0].RetryCount}
	for k, v := range m.Values {
		values[k] = v
	}
	pipe := w.RDB.TxPipeline()                       // MULTI: add to DLQ + ACK atomically
	pipe.XAdd(ctx, &redis.XAddArgs{Stream: w.Stream + ".dlq", Values: values})
	pipe.XAck(ctx, w.Stream, w.Group, m.ID)
	_, err = pipe.Exec(ctx)
	slog.ErrorContext(ctx, "message dead-lettered", "order_id", m.Values["order_id"], "attempts", p[0].RetryCount)
	return err == nil
}
```

Alert on `XLEN orders.created.dlq > 0` (via redis_exporter's
`--check-streams`). A DLQ nobody watches just loses messages slowly.

> ⚠️ **Caveat (verify):** check the `XPendingExt` result field (`RetryCount`) against your go-redis version, and check whether redis_exporter's `--check-streams` flag exists in yours (`redis_exporter --help`).

### 43.4 What you built

```
Git (app + deploy + platform)
  └─► CI: test → build multi-arch → scan → sign → SBOM → bump digest
        └─► Argo CD: platform add-ons + orbit (dev → staging → prod via PR)
              └─► Kubernetes: Envoy Gateway (TLS) → catalog / orders (Rollout canary) → CNPG Postgres, Redis
                    worker (KEDA 0..N on stream lag), CronJobs, PDBs, spread, NetworkPolicies, PSA restricted
                    Kyverno (signed images), ESO (secrets), Velero + CNPG PITR (backups)
                    OTel agent/gateway → Prometheus / Loki / Tempo → Grafana, SLO burn alerts → Alertmanager
```

---

## 44. Production-readiness checklist

Use it as a PR template for a new service, or as a launch review.

**Image**
- [ ] Multi-stage build; minimal base (distroless/Chainguard/slim); pinned by digest
- [ ] Runs as numeric non-root UID; no secrets in layers, ARGs or ENV
- [ ] Exec-form ENTRYPOINT; handles SIGTERM; graceful shutdown within the grace period
- [ ] Multi-arch; SBOM + provenance attached; signed; CRITICAL CVEs gated in CI

**Workload**
- [ ] Deployment/StatefulSet (never bare pods); ≥ 2 replicas (≥ 3 across zones for critical paths)
- [ ] CPU + memory requests from measured data; memory limit = request; runtime memory-aware (GOMEMLIMIT/JVM flags)
- [ ] startup/liveness/readiness probes, with liveness free of dependency checks
- [ ] preStop sleep + `terminationGracePeriodSeconds` budget adds up
- [ ] PDB allows ≥ 1 eviction; topology spread across zones/nodes
- [ ] PriorityClass set; tolerations tuned for stateless fast failover
- [ ] HPA/KEDA on a meaningful signal; `replicas` removed from the manifest; DB connection math checked
- [ ] Config via ConfigMap with hash-based rollout; secrets via ESO, mounted as files where possible

**Security**
- [ ] Own ServiceAccount; token not mounted unless needed; least-privilege RBAC
- [ ] PSA `restricted` passes: non-root, drop ALL, no privilege escalation, read-only root FS, seccomp RuntimeDefault
- [ ] NetworkPolicies: default deny + explicit ingress/egress; metadata endpoint blocked
- [ ] Workload identity for cloud APIs; no static credentials anywhere

**Observability**
- [ ] OTel traces with context propagation (HTTP + messaging); RED metrics; bounded-cardinality attributes
- [ ] Structured JSON logs with trace_id; no secrets/PII in logs
- [ ] Dashboard (RED + saturation + dependencies) as code
- [ ] SLOs defined; burn-rate alerts; every alert has a runbook; on-call rotation knows the service

**Delivery and data**
- [ ] Deployed only via GitOps; digest-pinned; promotion = PR; rollback = revert, rehearsed
- [ ] Canary with automated analysis for user-facing services
- [ ] Migrations are backward compatible (expand/contract); run as a hook, not at app start
- [ ] Backups automated *and restore-tested*; RPO/RTO written down and measured
- [ ] Load-tested to 2× expected peak; chaos experiments for pod, node and dependency failure

---

## 45. Cheat sheets

### 45.1 Docker

```bash
docker build -t img:tag --build-arg K=V --target stage .     # build (a specific stage)
docker buildx build --platform linux/amd64,linux/arm64 -t reg/img:tag --push .
docker run -d --name n -p 127.0.0.1:8080:8080 -e K=V -v vol:/data img:tag
docker logs -f --tail 100 n ; docker exec -it n sh ; docker inspect n | jq '.[0].State'
docker stats --no-stream ; docker system df ; docker system prune -a --volumes   # (careful)
docker compose up -d --build --wait ; docker compose logs -f svc ; docker compose down -v
docker run --rm -it --network container:n --pid container:n nicolaka/netshoot  # debug a distroless container
trivy image --severity HIGH,CRITICAL --ignore-unfixed img:tag
dive img:tag                                                   # explore layers
```

### 45.2 kubectl

```bash
kubectl get pods -A -o wide --field-selector=status.phase!=Running   # everything not running
kubectl get pods -o custom-columns=NAME:.metadata.name,STATUS:.status.phase,RESTARTS:.status.containerStatuses[*].restartCount,NODE:.spec.nodeName
kubectl describe pod P ; kubectl logs P -c C --previous ; kubectl logs -l app.kubernetes.io/name=X --prefix -f
kubectl get events -A --sort-by=.lastTimestamp | tail -30
kubectl top pods --containers ; kubectl top nodes
kubectl rollout status|history|undo|restart deploy/D
kubectl scale deploy/D --replicas=N
kubectl debug -it P --image=nicolaka/netshoot --target=C ; kubectl debug node/N -it --image=ubuntu
kubectl port-forward svc/S 8080:80
kubectl auth can-i VERB RESOURCE --as=system:serviceaccount:NS:SA
kubectl get endpointslices -l kubernetes.io/service-name=S
kubectl explain RESOURCE.spec.FIELD --recursive
kubectl diff -k overlay/ ; kubectl apply --server-side -k overlay/
kubectl get pod P -o jsonpath='{.status.containerStatuses[0].lastState.terminated.reason}'
kubectl cordon N ; kubectl drain N --ignore-daemonsets --delete-emptydir-data ; kubectl uncordon N
```

### 45.3 Helm, Argo CD, Rollouts, KEDA

```bash
helm upgrade --install R chart -n NS -f values.yaml --atomic --timeout 5m ; helm history R ; helm rollback R REV
helm template R chart -f values.yaml | kubeconform -strict ; helm get values R -n NS
argocd app list ; argocd app diff A ; argocd app sync A ; argocd app history A ; argocd app rollback A ID
kubectl argo rollouts get rollout R -w ; kubectl argo rollouts promote|abort|retry rollout R
kubectl get scaledobject,hpa -n NS ; kubectl describe scaledobject X     # KEDA status + trigger errors
```

---

## 46. Pitfalls that cause real outages

| Pitfall | Symptom | Fix (section) |
|---|---|---|
| Shell-form `CMD` / no SIGTERM handling | slow stops, dropped requests on deploy | exec form, signal handling (§3.3) |
| No preStop / instant shutdown | 502s on every rollout | preStop sleep + drain (§13.5) |
| Liveness checks dependencies | fleet-wide restart storm during a DB blip | liveness = process health only (§21.3) |
| CPU limits on latency-sensitive apps | p99 spikes at low average CPU | requests only, or generous limits (§23.3) |
| Runtime not container-aware | OOMKilled loops | GOMEMLIMIT / JVM / Node flags (§23.4) |
| No requests set | noisy neighbors, BestEffort evictions, broken HPA | measured requests (§23.6) |
| All replicas on one node/zone | outage on a single node reboot | topology spread (§24.1) |
| PDB that allows 0 evictions | stuck cluster upgrades | allow ≥ 1 (§22.3) |
| `:latest` / mutable tags | can't tell what's running; failed rollbacks | sha tags + digests (§10.1) |
| `ndots:5` + external calls | DNS latency, CoreDNS overload | FQDN dots, ndots 2, NodeLocal DNS (§15.3) |
| Long-lived gRPC/HTTP2 connections | new pods idle after scale-out | client LB / max connection age (§15.4) |
| HPA + `replicas` in Git | GitOps resets replicas every sync | remove `replicas` (§27.2) |
| Pods × pool size > DB max_connections | "too many connections" at peak scale | pooler / pool math (§27.4) |
| Unbounded metric labels | Prometheus OOM, monitoring outage | bounded attributes (§29.4) |
| Untested backups | data loss when you need them | automated restore tests (§39.4) |
| Locking migration at peak | full outage behind one ALTER | expand/contract, lock_timeout (§38) |
| Admission webhook down + `Fail` | nothing can be created | exclusions, HA, timeouts (§34.5) |
| Prune without protection | stateful resource deleted by GitOps | Prune=false on stateful (§36.5) |
| Pods can reach instance metadata | cloud credential theft via SSRF | egress policy, IMDSv2, workload identity (§34.3) |

---

## 47. Further reading: the real-life references

**Official documentation (the primary sources)**
- Kubernetes docs: kubernetes.io/docs. Especially *Concepts → Workloads*,
  *Configure Liveness, Readiness and Startup Probes*, *Pod Lifecycle*,
  *Production environment*, and *Security checklist*.
- Docker docs: docs.docker.com. *Building best practices*, *Multi-stage
  builds*, *Build secrets*, *Compose file reference*.
- OpenTelemetry docs: opentelemetry.io/docs. *Go getting started*, the
  *Collector* configuration, and *Semantic Conventions* for HTTP and
  messaging.
- Prometheus (prometheus.io/docs) and *Prometheus Operator* docs; the
  Grafana Loki and Tempo docs.
- Gateway API (gateway-api.sigs.k8s.io), Argo CD and Argo Rollouts
  (argo-cd.readthedocs.io, argoproj.github.io/rollouts), KEDA (keda.sh),
  CloudNativePG (cloudnative-pg.io), Karpenter (karpenter.sh), Kyverno
  (kyverno.io).

**Books**
- *Kubernetes in Action*, 2nd ed. (Marko Lukša): the deepest
  practitioner's explanation of how Kubernetes works.
- *Production Kubernetes* (Rosso, Lander, Brand, Harris): building
  platforms on Kubernetes; real trade-offs from VMware field engineers.
- *Kubernetes Patterns*, 2nd ed. (Ibryam & Huß): sidecar, ambassador,
  init container, operator and more, as reusable patterns.
- *Docker Deep Dive* (Nigel Poulton): short, current, practical.
- *Site Reliability Engineering* and *The Site Reliability Workbook*
  (Google, free at sre.google/books): SLOs, burn-rate alerting (Workbook
  ch. 5), postmortems.
- *Observability Engineering* (Majors, Fong-Jones, Miranda).
- *Learning OpenTelemetry* (Young & Parker).
- *Designing Distributed Systems*, 2nd ed. (Brendan Burns).
- *Container Security* (Liz Rice), and *Hacking Kubernetes* (Martin &
  Hausenblas) for the attacker's view.

**Field guides and real-world material**
- **Kubernetes Failure Stories**, k8s.af: public postmortems; read them.
- learnk8s.io: excellent deep dives (graceful shutdown, CPU limits,
  architecture, the "troubleshooting deployments" flowchart).
- Google Cloud *Kubernetes best practices*, AWS *EKS Best Practices
  Guide* (aws.github.io/aws-eks-best-practices), Azure AKS baseline
  architecture. These are opinionated, battle-tested checklists.
- CNCF *Cloud Native Security Whitepaper* and the NSA/CISA *Kubernetes
  Hardening Guide*.
- OpenGitOps principles (opengitops.dev), SLSA (slsa.dev), Sigstore
  (sigstore.dev).
- Brendan Gregg's *USE Method* and Tom Wilkie's *RED Method* write-ups.
- Henning Jacobs' *Kubernetes on AWS: CPU limits* talks and the "Stop
  using CPU limits" discussion; read both sides before deciding (§23.3).

**Practice**
- killercoda.com and killer.sh: browser-based Kubernetes labs (and CKA,
  CKAD and CKS exam simulators). CKA → CKAD → CKS is a sensible
  certification order once you've finished this guide.
- Re-run §43's acceptance tests after every significant change to the
  platform. That's what keeps this system production-grade.

