# Architecture Documentation

## System Design for 1 Billion Users

This document outlines the architecture decisions and scaling strategies for supporting up to 1 billion users.

## High-Level Architecture

```
┌─────────────┐
│   Clients   │ (Mobile/Web)
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────┐
│   CDN & Edge (CloudFront/CF)        │
└──────┬──────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────┐
│   API Gateway (Envoy/Kong) + WAF    │
└──────┬──────────────────────────────┘
       │
       ├──► Auth Service ──► Auth DB
       ├──► User Service ──► User DB
       ├──► Product Service ──► Product DB
       ├──► Task Service ──► Task DB
       └──► Media Service ──► Media DB
       │
       ▼
┌─────────────────────────────────────┐
│   Message Bus (Kafka)                │
└──────┬──────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────┐
│   Analytics (ClickHouse)             │
└──────────────────────────────────────┘

┌─────────────────────────────────────┐
│   Cache Layer (Redis Cluster)       │
└──────────────────────────────────────┘

┌─────────────────────────────────────┐
│   Object Storage (S3/MinIO)         │
└──────────────────────────────────────┘
```

## Component Details

### 1. API Services

**Technology**: Go 1.21+ with Fiber framework

**Why Go?**
- Compiled language with low latency
- Excellent concurrency model (goroutines)
- Small memory footprint
- Fast startup time

**Scaling Strategy**:
- Stateless services → horizontal scaling
- Each service runs in Kubernetes pods
- Auto-scaling based on CPU/RPS metrics
- Load balancing via Kubernetes Service

**Service Ports**:
- Auth Service: 3001
- User Service: 3002
- Product Service: 3003
- Task Service: 3004
- Media Service: 3005

### 2. Database Layer

**Current**: PostgreSQL 15 (single instance for dev)

**Production Recommendations**:

**Option A: CockroachDB** (Recommended for multi-region)
- Distributed SQL database
- Automatic sharding and replication
- ACID transactions across regions
- Built-in geo-partitioning

**Option B: Vitess** (MySQL-based)
- Horizontal sharding for MySQL
- Battle-tested at YouTube scale
- Good for read-heavy workloads

**Option C: Citus** (PostgreSQL extension)
- PostgreSQL-compatible
- Automatic sharding
- Good for analytics workloads

**Sharding Strategy**:
- Shard by user_id (hash-based or range-based)
- Each shard handles ~10-100M users
- Cross-shard queries minimized

**Read Replicas**:
- Primary for writes, multiple replicas for reads
- Route reads to nearest replica (geo-aware)
- Eventual consistency acceptable for profile reads

### 3. Caching Layer

**Technology**: Redis Cluster

**Use Cases**:
- Session storage (refresh tokens)
- Hot user profiles
- Rate limiting counters
- Frequently accessed data

**Strategy**:
- L1: CDN edge cache (static content)
- L2: Redis cluster (hot data, 1-24 hour TTL)
- L3: In-memory cache per pod (very hot, 1-5 min TTL)

**Cache Invalidation**:
- Publish invalidation events via Kafka
- TTL-based expiration
- Cache-aside pattern

### 4. Message Queue

**Technology**: Apache Kafka

**Topics**:
- `user.created`, `user.updated`
- `product.created`, `product.updated`
- `task.created`, `task.updated`
- `media.uploaded`

**Consumers**:
- Analytics pipeline (ClickHouse)
- Search indexer (Elasticsearch)
- Notification service
- Audit logging

**Scaling**:
- Partition topics by user_id for parallel processing
- Multiple consumer groups for different use cases
- Replication factor 3+ for durability

### 5. Object Storage

**Technology**: S3-compatible (AWS S3, MinIO, etc.)

**Use Cases**:
- User avatars
- Product images
- Media files
- Backup storage

**Strategy**:
- Presigned URLs for direct upload/download
- CDN in front (CloudFront/Cloudflare)
- Lifecycle policies for old data
- Multi-region replication

### 6. Authentication & Authorization

**Technology**: JWT (JSON Web Tokens)

**Token Strategy**:
- Access tokens: Short-lived (15 minutes)
- Refresh tokens: Long-lived (7 days), stored in DB
- Token rotation on refresh
- Revocation via refresh token deletion

**Security**:
- HTTPS/TLS everywhere
- mTLS for service-to-service communication
- Rate limiting per IP/user
- WAF for DDoS protection

## Scaling Strategies

### Horizontal Scaling

All services are stateless and can scale horizontally:

```yaml
# Kubernetes HPA example
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: user-service-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: user-service
  minReplicas: 3
  maxReplicas: 100
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

### Database Scaling

**Phase 1** (0-10M users):
- Single PostgreSQL instance with read replicas
- Connection pooling (PgBouncer)
- Query optimization and indexing

**Phase 2** (10M-100M users):
- Shard database by user_id
- Use Citus or Vitess for automatic sharding
- Cross-shard queries minimized

**Phase 3** (100M-1B users):
- Multi-region deployment
- CockroachDB or Spanner for global distribution
- Geo-partitioning (data near users)

### Caching Strategy

**Cache Hierarchy**:
1. CDN (edge locations) - static content, images
2. Redis Cluster (regional) - hot data, sessions
3. In-memory (per pod) - very hot data

**Cache Patterns**:
- Cache-aside: App checks cache, then DB
- Write-through: Write to cache and DB
- Write-behind: Write to cache, async to DB

### Rate Limiting

**Layers**:
1. Edge (Cloudflare/AWS WAF) - IP-based
2. API Gateway - User-based, token-based
3. Service level - Per-endpoint limits

**Implementation**:
- Redis for distributed rate limiting
- Token bucket algorithm
- Sliding window for accuracy

## Multi-Region Deployment

### Architecture

```
Region 1 (US-East)          Region 2 (EU-West)          Region 3 (APAC)
┌──────────────┐           ┌──────────────┐           ┌──────────────┐
│ API Services │           │ API Services │           │ API Services │
│   (Active)   │           │   (Active)   │           │   (Active)   │
└──────┬───────┘           └──────┬───────┘           └──────┬───────┘
       │                         │                         │
       ▼                         ▼                         ▼
┌──────────────┐           ┌──────────────┐           ┌──────────────┐
│   Database   │◄──────────►│   Database   │◄──────────►│   Database   │
│  (Primary)   │  Replication│  (Replica)   │  Replication│  (Replica)   │
└──────────────┘           └──────────────┘           └──────────────┘
```

### Traffic Routing

- **DNS**: Route53 latency-based routing
- **API Gateway**: Geo-aware routing
- **Database**: Read from nearest replica, write to primary

### Data Consistency

- **Strong Consistency**: User auth, payments (synchronous replication)
- **Eventual Consistency**: User profiles, product listings (async replication)
- **Conflict Resolution**: Last-write-wins or application-level merging

## Monitoring & Observability

### Metrics (Prometheus)

- Request rate (RPS)
- Latency (p50, p95, p99)
- Error rate
- Database connection pool usage
- Cache hit/miss ratio
- Kafka lag

### Logging (Loki/ELK)

- Structured JSON logs
- Centralized log aggregation
- Log retention: 30 days hot, 1 year cold

### Tracing (Jaeger)

- Distributed tracing across services
- Request flow visualization
- Performance bottleneck identification

### Alerting

- SLO targets: 99.9% uptime, <200ms p95 latency
- Alert on error rate > 1%
- Alert on latency p95 > 500ms
- Alert on database connection pool exhaustion

## Security

### Network Security

- VPC isolation
- Security groups (least privilege)
- mTLS for service-to-service
- WAF for DDoS protection

### Data Security

- Encryption at rest (database, S3)
- Encryption in transit (TLS 1.3)
- Secrets management (Vault/AWS Secrets Manager)
- Regular security audits

### Compliance

- GDPR: Data deletion workflows
- Data residency per region
- Audit logs for admin operations
- PII encryption

## Cost Optimization

### Compute

- Use spot instances for non-critical workloads
- Reserved instances for baseline capacity
- Auto-scaling to reduce idle costs

### Storage

- S3 lifecycle policies (move to Glacier after 90 days)
- Database archiving (move old data to cold storage)
- Compression for logs and backups

### Network

- CDN for static content (reduces origin traffic)
- Regional data transfer optimization
- Use private networking where possible

## Disaster Recovery

### Backup Strategy

- Database: Daily snapshots + continuous WAL archiving
- S3: Versioning enabled, cross-region replication
- Configuration: Git-based, versioned

### RTO/RPO Targets

- **Critical Services** (Auth, Payments): RTO < 1 hour, RPO < 5 minutes
- **Standard Services** (User, Product): RTO < 4 hours, RPO < 1 hour
- **Analytics**: RTO < 24 hours, RPO < 6 hours

### Failover Procedures

1. Detect failure (health checks, monitoring)
2. Route traffic away from failed region
3. Promote replica to primary (if database)
4. Scale up services in healthy regions
5. Verify data consistency
6. Document incident

## Performance Targets

### Latency

- API response time: p95 < 200ms
- Database query: p95 < 50ms
- Cache hit: < 5ms

### Throughput

- Peak: 500k requests/second
- Sustained: 100k requests/second
- Per service: 50k requests/second

### Availability

- Target: 99.9% uptime (8.76 hours downtime/year)
- Critical paths: 99.99% (52.56 minutes/year)

## Next Steps

1. **Phase 1** (Current): Single-region, basic services
2. **Phase 2**: Add caching, optimize queries, add monitoring
3. **Phase 3**: Multi-region, database sharding, advanced caching
4. **Phase 4**: Full observability, chaos engineering, cost optimization

