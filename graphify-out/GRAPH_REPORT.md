# Graph Report - .  (2026-06-25)

## Corpus Check
- Corpus is ~9,284 words - fits in a single context window. You may not need a graph.

## Summary
- 398 nodes · 591 edges · 46 communities (36 shown, 10 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 39 edges (avg confidence: 0.82)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Server Bootstrap|Server Bootstrap]]
- [[_COMMUNITY_Error & Response Middleware|Error & Response Middleware]]
- [[_COMMUNITY_Account Repository|Account Repository]]
- [[_COMMUNITY_DI Container|DI Container]]
- [[_COMMUNITY_Pagination Table Logic|Pagination Table Logic]]
- [[_COMMUNITY_User Credential Repository|User Credential Repository]]
- [[_COMMUNITY_User Credential Service|User Credential Service]]
- [[_COMMUNITY_Authentication Service|Authentication Service]]
- [[_COMMUNITY_User Entity & DTOs|User Entity & DTOs]]
- [[_COMMUNITY_Account Service|Account Service]]
- [[_COMMUNITY_Organization Domain|Organization Domain]]
- [[_COMMUNITY_Pagination Query Builder|Pagination Query Builder]]
- [[_COMMUNITY_SQLGORM Connection|SQL/GORM Connection]]
- [[_COMMUNITY_User Repository|User Repository]]
- [[_COMMUNITY_User Service|User Service]]
- [[_COMMUNITY_Architecture Docs|Architecture Docs]]
- [[_COMMUNITY_Common Model Base|Common Model Base]]
- [[_COMMUNITY_Account HTTP Handler|Account HTTP Handler]]
- [[_COMMUNITY_Request Context|Request Context]]
- [[_COMMUNITY_Dashboard Entities|Dashboard Entities]]
- [[_COMMUNITY_App Status Codes|App Status Codes]]
- [[_COMMUNITY_Organization Repository|Organization Repository]]
- [[_COMMUNITY_Organization Signup|Organization Signup]]
- [[_COMMUNITY_DB View Service|DB View Service]]
- [[_COMMUNITY_External Service Config|External Service Config]]
- [[_COMMUNITY_Pagination DTO|Pagination DTO]]
- [[_COMMUNITY_Transaction Manager|Transaction Manager]]
- [[_COMMUNITY_Auth HTTP Handler|Auth HTTP Handler]]
- [[_COMMUNITY_Organization HTTP Handler|Organization HTTP Handler]]
- [[_COMMUNITY_User DB Model|User DB Model]]
- [[_COMMUNITY_Database Config Loader|Database Config Loader]]
- [[_COMMUNITY_AI Provider Configs|AI Provider Configs]]
- [[_COMMUNITY_Account DB Model|Account DB Model]]
- [[_COMMUNITY_Architecture Rationale|Architecture Rationale]]
- [[_COMMUNITY_Main DB YAML Pairs|Main DB YAML Pairs]]
- [[_COMMUNITY_Redis YAML Pairs|Redis YAML Pairs]]
- [[_COMMUNITY_REST Server YAML Pairs|REST Server YAML Pairs]]
- [[_COMMUNITY_Organization DB Model|Organization DB Model]]
- [[_COMMUNITY_Generic Response Helpers|Generic Response Helpers]]
- [[_COMMUNITY_Structured Logging Rationale|Structured Logging Rationale]]
- [[_COMMUNITY_DTO View Tag|DTO View Tag]]
- [[_COMMUNITY_Module Root|Module Root]]
- [[_COMMUNITY_Transaction Repo|Transaction Repo]]

## God Nodes (most connected - your core abstractions)
1. `Repository` - 13 edges
2. `Service` - 13 edges
3. `GetListRequest` - 12 edges
4. `Container` - 12 edges
5. `Service` - 11 edges
6. `Table` - 11 edges
7. `GetConfig()` - 10 edges
8. `Repository` - 10 edges
9. `Response` - 10 edges
10. `Context` - 9 edges

## Surprising Connections (you probably didn't know these)
- `Connect()` --calls--> `GetDatabaseDriverMigration()`  [INFERRED]
  cmd/db/migrate/db_migrate.go → pkg/infrastructure/database/sql_connection.go
- `Connect()` --calls--> `GetConfig()`  [INFERRED]
  cmd/db/migrate/db_migrate.go → config/loader.go
- `MigrateUpAll()` --calls--> `GetConfig()`  [INFERRED]
  cmd/db/migrate/db_migrate.go → config/loader.go
- `startRest()` --calls--> `GetConfig()`  [INFERRED]
  cmd/server/rest.go → config/loader.go
- `startRest()` --calls--> `NewRest()`  [INFERRED]
  cmd/server/rest.go → pkg/infrastructure/server/rest.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Layered Request Flow** — _github__copilot_instructions_handler_layer, _github__copilot_instructions_service_layer, _github__copilot_instructions_repository_layer, _github__copilot_instructions_database_layer [EXTRACTED 1.00]
- **Database Configuration Variants** — example_database_main_database, example_database_redis, resources_database_main_database, resources_database_redis [INFERRED 0.95]
- **Model Provider Configuration** — example_server_groq, resources_server_groq, resources_server_minimax [INFERRED 0.85]

## Communities (46 total, 10 thin omitted)

### Community 0 - "Server Bootstrap"
Cohesion: 0.08
Nodes (28): Command, Config, GetConfig(), DatabaseList, Groq, jwtError(), Protected(), SuccessHandler() (+20 more)

### Community 1 - "Error & Response Middleware"
Cohesion: 0.12
Nodes (20): Error, DefaultErrorHandler(), DefaultResponseHandler(), Ctx, AppStatus, T, Time, ErrorResponse (+12 more)

### Community 2 - "Account Repository"
Cohesion: 0.14
Nodes (15): Account, Repository, NewRepository(), AccountType, NewAccountDtoFromModel(), AccountDto, AccountFilterDto, AccountRequest (+7 more)

### Community 3 - "DI Container"
Cohesion: 0.17
Nodes (10): Container, NewContainer(), Context, DB, Manager, App, Manager, contextKey (+2 more)

### Community 4 - "Pagination Table Logic"
Cohesion: 0.27
Nodes (17): Table, calculateTotalPages(), getMappingField(), getOperator(), NewTable(), sortValidation(), Table[T, U, V], validateAllowedFields() (+9 more)

### Community 5 - "User Credential Repository"
Cohesion: 0.24
Nodes (8): Context, DB, PaginationRequestDto, ResultPagination, UserCredential, UserCredentialDto, Repository, NewRepository()

### Community 6 - "User Credential Service"
Cohesion: 0.23
Nodes (8): Context, PaginationRequestDto, Repository, ResultPagination, UserCredential, UserCredentialDto, Service, NewService()

### Community 7 - "Authentication Service"
Cohesion: 0.15
Nodes (12): Service, NewUserSessionDtoFromClaims(), LoginDto, LoginRequestDto, SignUpRequestDto, SignUpResponseDto, UserSessionDto, LoginDto (+4 more)

### Community 8 - "User Entity & DTOs"
Cohesion: 0.24
Nodes (10): NewUserCredentialDtoFromModel(), NewUserDtoFromModel(), UserCredentialDto, UserDto, UserFilterDto, GetListRequest, UserCredential, UserCredentialDto (+2 more)

### Community 9 - "Account Service"
Cohesion: 0.30
Nodes (8): Service, NewService(), AccountDto, Context, Manager, PaginationRequestDto, Repository, ResultPagination

### Community 10 - "Organization Domain"
Cohesion: 0.18
Nodes (12): NewOrganizationDtoFromModel(), OrganizationDto, Config, configDefault(), getOrganizationByOrigin(), New(), OrganizationDto, CommonWithIDs (+4 more)

### Community 11 - "Pagination Query Builder"
Cohesion: 0.19
Nodes (3): FilterItem, GetListRequest, View

### Community 12 - "SQL/GORM Connection"
Cohesion: 0.36
Nodes (10): connect(), GetDatabaseDriverMigration(), getGormDialect(), getSqlDB(), NewSQLConnection(), SQLConnection, Dialector, Driver (+2 more)

### Community 13 - "User Repository"
Cohesion: 0.30
Nodes (7): Context, DB, PaginationRequestDto, ResultPagination, UserDto, Repository, NewRepository()

### Community 14 - "User Service"
Cohesion: 0.30
Nodes (7): Context, PaginationRequestDto, Repository, ResultPagination, UserDto, Service, NewService()

### Community 15 - "Architecture Docs"
Cohesion: 0.27
Nodes (10): Clean Architecture, Database Layer, Handler Layer, Repository Layer, Service Layer, docs/api-standards.md, docs/architecture.md, docs/database-standards.md (+2 more)

### Community 16 - "Common Model Base"
Cohesion: 0.31
Nodes (5): CommonWithID, CommonWithIDs, DeletedAt, DB, Time

### Community 17 - "Account HTTP Handler"
Cohesion: 0.50
Nodes (8): CreateAccount(), DeleteAccountByID(), FindAccountByCustomerID(), FindAccountByID(), FindAllAccount(), UpdateAccountByID(), Handler, Service

### Community 18 - "Request Context"
Cohesion: 0.43
Nodes (7): GetOrganization(), GetOrigin(), GetOriginTypeKey(), GetUserCredential(), OrganizationData, UserCredentialData, Context

### Community 19 - "Dashboard Entities"
Cohesion: 0.39
Nodes (7): CashFlowData, DailySales, DashboardFilterDto, DashboardSummaryDto, ExpenseByCategory, ProfitLossData, Time

### Community 20 - "App Status Codes"
Cohesion: 0.43
Nodes (7): AppStatus, ClientErrorCase, ServerErrorCase, NewClientErrorAppStatus(), NewServerErrorAppStatus(), NewSuccessAppStatus(), SuccessCase

### Community 21 - "Organization Repository"
Cohesion: 0.48
Nodes (5): Repository, NewRepository(), Context, DB, Organization

### Community 22 - "Organization Signup"
Cohesion: 0.38
Nodes (6): NewService(), Context, Service, Repository, SignUpRequestDto, SignUpResponseDto

### Community 23 - "DB View Service"
Cohesion: 0.52
Nodes (5): T, View, Service, NewViewService(), service[T]

### Community 24 - "External Service Config"
Cohesion: 0.40
Nodes (6): Groq, MessageBroker, Minimax, Server, ServerList, Server

### Community 25 - "Pagination DTO"
Cohesion: 0.50
Nodes (5): FilterItem, Pagination, PaginationRequestDto, ResultPagination, T

### Community 26 - "Transaction Manager"
Cohesion: 0.40
Nodes (3): Context, DB, AppRepository

### Community 27 - "Auth HTTP Handler"
Cohesion: 0.50
Nodes (3): SignIn(), Handler, Service

### Community 28 - "Organization HTTP Handler"
Cohesion: 0.50
Nodes (3): SignUp(), Handler, Service

### Community 30 - "User DB Model"
Cohesion: 0.67
Nodes (4): User, UserType, CommonWithIDs, UserCredential

### Community 31 - "Database Config Loader"
Cohesion: 1.00
Nodes (3): Database, Database, DatabaseList

### Community 32 - "AI Provider Configs"
Cohesion: 0.67
Nodes (3): Example Groq Model Config, Resource Groq Model Config, Resource MiniMax Model Config

### Community 33 - "Account DB Model"
Cohesion: 1.00
Nodes (3): Account, AccountType, CommonWithIDs

## Knowledge Gaps
- **93 isolated node(s):** `Database`, `ServerList`, `DatabaseList`, `Groq`, `Minimax` (+88 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **10 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `GetConfig()` connect `Server Bootstrap` to `Authentication Service`?**
  _High betweenness centrality (0.040) - this node is a cross-community bridge._
- **Why does `NewContainer()` connect `DI Container` to `Server Bootstrap`?**
  _High betweenness centrality (0.022) - this node is a cross-community bridge._
- **What connects `Database`, `ServerList`, `DatabaseList` to the rest of the system?**
  _95 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Server Bootstrap` be split into smaller, more focused modules?**
  _Cohesion score 0.08412698412698413 - nodes in this community are weakly interconnected._
- **Should `Error & Response Middleware` be split into smaller, more focused modules?**
  _Cohesion score 0.11666666666666667 - nodes in this community are weakly interconnected._
- **Should `Account Repository` be split into smaller, more focused modules?**
  _Cohesion score 0.14130434782608695 - nodes in this community are weakly interconnected._