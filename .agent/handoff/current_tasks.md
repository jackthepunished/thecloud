# 🚀 Mini AWS - Master Task Breakdown

## Overview
Granular task checklist for Phases 5-7, organized into 2-week sprints.  
Each implementation task has a corresponding **🧪 Test** task.

---

## 📋 Phase 5: The Console (Sprints 1-3)

### Sprint 1: Backend Foundation
- [x] **Architect**: Define `domain/dashboard.go` (ResourceSummary, MetricPoint)
- [x] **Architect**: Define `domain/ws_event.go` (WSEventType, WSEvent struct)
- [x] **Architect**: Create ADR-006: Real-time Communication Strategy
- [x] **Backend**: Implement `services/dashboard.go` (GetSummary, CountResources)
- [x] 🧪 **QA**: Unit tests for DashboardService (≥90% coverage)
- [x] **Backend**: Add `GET /api/dashboard/summary` endpoint
- [x] 🧪 **QA**: API tests for dashboard/summary endpoint
- [x] **Infra**: Create `migrations/009_metrics_history.sql`
- [ ] 🧪 **QA**: Migration rollback test

### Sprint 2: Real-time Streaming
- [x] **Infra**: Implement `docker/stats.go` (real-time container metrics)
- [x] 🧪 **QA**: Integration test for Docker stats adapter
- [x] **Backend**: Implement SSE endpoint `GET /api/dashboard/stream`
- [x] 🧪 **QA**: SSE connection and event delivery test
- [x] **Backend**: Create WebSocket hub `handlers/ws/hub.go`
- [x] 🧪 **QA**: WebSocket lifecycle test (connect/message/disconnect)
- [x] **Security**: Add WS handshake authentication middleware
- [x] 🧪 **QA**: Auth rejection test for invalid API keys
- [ ] **Security**: Configure CORS for frontend origin
- [ ] **Platform**: Add `mini_aws_ws_connections_active` gauge
- [ ] 🧪 **QA**: Metrics endpoint validation

### Sprint 3: Frontend Dashboard
- [ ] **Frontend**: Initialize Next.js 14 project in `/frontend`
- [ ] **Frontend**: Create dashboard layout with sidebar navigation
- [ ] 🧪 **QA**: Component snapshot tests
- [ ] **Frontend**: Build ResourceCard components (Instances, Volumes, VPCs)
- [ ] 🧪 **QA**: Unit tests for ResourceCard props
- [ ] **Frontend**: Implement real-time metrics charts (CPU/Memory)
- [ ] 🧪 **QA**: Chart rendering with mock data
- [ ] **Frontend**: Build Activity Feed (audit logs stream)
- [ ] **CLI**: Add `cloud dashboard open` command
- [ ] 🧪 **QA**: CLI command execution test
- [ ] **Docs**: Create `docs/guides/console.md`
- [ ] 🧪 **QA**: Playwright E2E tests for full dashboard flow

---

## 📋 Phase 6: The Elastic Cloud (Sprints 4-6)

### Sprint 4: Load Balancer Core
- [ ] **Architect**: Define `domain/loadbalancer.go` (LB, Target, Algorithm)
- [ ] **Architect**: Define `domain/healthcheck.go` (HealthCheckConfig)
- [ ] **Architect**: Create ADR-007: Load Balancer Implementation Strategy
- [ ] **Backend**: Implement `services/loadbalancer.go` (Create, AddTarget)
- [ ] 🧪 **QA**: Unit tests for LBService CRUD operations
- [ ] **Backend**: Implement LB state machine (CREATING → ACTIVE → DRAINING)
- [ ] 🧪 **QA**: State transition tests (valid + invalid paths)
- [ ] **Infra**: Create `migrations/006_load_balancers.sql`
- [ ] 🧪 **QA**: Migration up/down test

### Sprint 5: Docker Proxying & Health Checks
- [ ] **Infra**: Implement HAProxy config generator
- [ ] 🧪 **QA**: Config generation unit test
- [ ] **Infra**: Create `docker/loadbalancer.go` (deploy proxy container)
- [ ] 🧪 **QA**: Integration test for proxy container lifecycle
- [ ] **Backend**: Implement health check scheduler (background goroutine)
- [ ] 🧪 **QA**: Health check timing and result recording test
- [ ] **Backend**: Add `/api/lb` endpoints (CRUD + targets)
- [ ] 🧪 **QA**: API integration tests for LB endpoints
- [ ] **Security**: Ensure LB only routes within same VPC
- [ ] 🧪 **QA**: Cross-VPC routing rejection test
- [ ] **Platform**: Add `mini_aws_lb_requests_total` counter

### Sprint 6: Auto-Scaling 🆕
- [ ] **Architect**: Define `domain/autoscaling.go` (ScalingPolicy, ScalingGroup)
- [ ] **Architect**: Create ADR-008: Auto-Scaling Strategy
- [ ] **Backend**: Implement `services/autoscaling.go` (Evaluate, ScaleOut, ScaleIn)
- [ ] 🧪 **QA**: Unit tests for scaling decision logic
- [ ] **Backend**: Implement metric-based scaling triggers (CPU > 80%)
- [ ] 🧪 **QA**: Trigger threshold boundary tests
- [ ] **Backend**: Add cooldown period logic (prevent thrashing)
- [ ] 🧪 **QA**: Cooldown enforcement test
- [ ] **Infra**: Integrate with Instance service for spawn/terminate
- [ ] 🧪 **QA**: Integration test for instance scaling
- [ ] **CLI**: Add `cloud autoscaling create|list|describe` commands
- [ ] 🧪 **QA**: CLI command output validation
- [ ] **Docs**: Create `docs/guides/autoscaling.md`

---

## 📋 Phase 7: The Managed Cloud (Sprints 7-9)

### Sprint 7: RDS-lite Core
- [ ] **Architect**: Define `domain/rds.go` (ManagedDatabase, DBEngine)
- [ ] **Architect**: Define `domain/snapshot.go` (VolumeSnapshot)
- [ ] **Architect**: Create ADR-009: Managed Database Architecture
- [ ] **Backend**: Implement `services/rds.go` (Create, Delete, GetEndpoint)
- [ ] 🧪 **QA**: Unit tests for RDS service
- [ ] **Infra**: Create preconfigured Docker images (postgres:15, mysql:8)
- [ ] 🧪 **QA**: Container health check validation
- [ ] **Infra**: Create `migrations/007_managed_databases.sql`
- [ ] **Security**: Implement credential generation (random password)
- [ ] 🧪 **QA**: Password strength validation test
- [ ] **Security**: Encrypt DB credentials at rest
- [ ] 🧪 **QA**: Encryption/decryption round-trip test

### Sprint 8: Snapshots & Backups
- [ ] **Infra**: Implement `docker/snapshot.go` (volume backup to tar)
- [ ] 🧪 **QA**: Snapshot creation and file integrity test
- [ ] **Backend**: Implement `services/snapshot.go` (Create, Restore, List)
- [ ] 🧪 **QA**: Full backup/restore cycle test
- [ ] **Backend**: Add scheduled backup support (cron-like)
- [ ] 🧪 **QA**: Scheduler timing test
- [ ] **CLI**: Add `cloud rds create|list|connect` commands
- [ ] 🧪 **QA**: RDS CLI integration test
- [ ] **CLI**: Add `cloud snapshot create|list|restore` commands
- [ ] 🧪 **QA**: Snapshot CLI integration test

### Sprint 9: Additional Managed Services 🆕
- [ ] **Architect**: Define `domain/cache.go` (ManagedCache - CloudCache)
- [ ] **Backend**: Implement `services/cache.go` (CloudCache service)
- [ ] 🧪 **QA**: CloudCache CRUD and connection test
- [ ] **Infra**: Redis cluster mode support (master + replica)
- [ ] 🧪 **QA**: Replication failover test
- [ ] **Architect**: Define `domain/queue.go` (MessageQueue - CloudQueue)
- [ ] **Backend**: Implement `services/queue.go` (SendMessage, ReceiveMessage)
- [ ] 🧪 **QA**: Message ordering and visibility test
- [ ] **CLI**: Add `cloud cache create|list` commands
- [ ] 🧪 **QA**: Cache CLI test
- [ ] **CLI**: Add `cloud queue create|send|receive` commands
- [ ] 🧪 **QA**: Queue CLI test
- [ ] **Docs**: Update all guides for new services

---

## 📋 Phase Final: Marketplace (Sprint 10)

### Sprint 10: Templates & 1-Click Deploy
- [ ] **Architect**: Define `domain/template.go` (CloudFormation-lite schema)
- [ ] **Backend**: Implement template parser and validator
- [ ] 🧪 **QA**: Template validation test (valid + malformed)
- [ ] **Backend**: Implement stack provisioning orchestrator
- [ ] 🧪 **QA**: Full stack deployment E2E test
- [ ] **Frontend**: Build template gallery UI
- [ ] 🧪 **QA**: Gallery rendering and selection test
- [ ] **Docs**: Create sample templates (WordPress, MERN, Django)

---

## 🆕 Additional Managed Services (Mini AWS Branding)

| Service | Mini AWS Name | AWS Equivalent | Priority |
|---------|---------------|----------------|----------|
| **Cache** | **CloudCache** | ElastiCache | P1 (Sprint 9) |
| **Queue** | **CloudQueue** | SQS | P2 (Sprint 9) |
| **DNS** | **Route25** | Route53 | P3 (Future) |
| **Functions** | **CloudFunctions** | Lambda | P3 (Future) |
| **Secrets** | **Vault** | Secrets Manager | P2 (Sprint 8) |

---

## 📊 Progress Summary

| Phase | Sprints | Tasks | Tests | Status |
|-------|---------|-------|-------|--------|
| Phase 5 | 1-3 | 24 | 16 | ⬜ Not Started |
| Phase 6 | 4-6 | 27 | 18 | ⬜ Not Started |
| Phase 7 | 7-9 | 28 | 19 | ⬜ Not Started |
| Marketplace | 10 | 8 | 4 | ⬜ Not Started |
| **Total** | **10** | **87** | **57** | - |
