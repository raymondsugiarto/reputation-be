package di

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/raymondsugiarto/reputation-be/pkg/infrastructure/database"
	"github.com/raymondsugiarto/reputation-be/pkg/module/admin"
	"github.com/raymondsugiarto/reputation-be/pkg/module/authentication"
	"github.com/raymondsugiarto/reputation-be/pkg/module/customer"
	"github.com/raymondsugiarto/reputation-be/pkg/module/organization"
	"github.com/raymondsugiarto/reputation-be/pkg/module/user"
	usercredential "github.com/raymondsugiarto/reputation-be/pkg/module/user-credential"
	txDatabase "github.com/raymondsugiarto/reputation-be/pkg/shared/database/transaction"
)

type Container interface {
	RegisterServices()
}

type container struct {
	app *fiber.App
}

func NewContainer(app *fiber.App) Container {
	return &container{
		app: app,
	}
}

func (c *container) RegisterServices() {
	c.gormManager()

	c.userCredentialService()
	c.userService()
	c.organizationService()

	// customerService is registered BEFORE authenticationService so the
	// auth module can use it as a customerLookup closure (for
	// approval-status gating on sign-in).
	c.customerService()

	c.authenticationService()

	// Internal-admin platform service. Wired after the customer service
	// because the admin router re-uses customer.Service.FindAll for its
	// platform-wide customer listing.
	c.adminService()

	// Financial tracking services
	// c.customerService()
	// c.merchantService()
	// c.categoryService()
	// c.accountService()
	// c.transactionService()
	// c.dashboardService()
	// c.groqService()
	// c.minimaxService()
	// c.pdfService()
}

func (c *container) add(name string, service any) {
	if c.app.State().Has(name) {
		panic("Service already registered: " + name)
	}
	c.app.State().Set(name, service)
}

func (c *container) gormManager() {
	gormTxManager := txDatabase.NewGormManager(database.DBConn)
	c.add(txDatabase.GormServiceName, gormTxManager)
}

func (c *container) gormManagerWithTx() txDatabase.Manager {
	return fiber.MustGetState[txDatabase.Manager](c.app.State(), txDatabase.GormServiceName)
}

func (c *container) authenticationService() {
	userCredentialService := fiber.MustGetState[usercredential.Service](c.app.State(), usercredential.ServiceName)
	customerService := fiber.MustGetState[customer.Service](c.app.State(), customer.ServiceName)

	// Wire the customer-approval gate into authentication.SignIn so
	// PENDING_APPROVAL / REJECTED customers cannot log in. Defined as
	// a closure so authentication doesn't have to import the customer
	// module directly (which would create a cycle).
	customerLookup := func(ctx context.Context, userID string) (authentication.CustomerStatusInfo, error) {
		cust, err := customerService.FindByUserID(ctx, userID)
		if err != nil || cust == nil {
			return authentication.CustomerStatusInfo{}, err
		}
		return authentication.CustomerStatusInfo{Status: cust.Status}, nil
	}

	authenticationService := authentication.NewServiceWithCustomerLookup(
		userCredentialService,
		customerLookup,
	)
	c.add(authentication.ServiceName, authenticationService)
}

func (c *container) userCredentialService() {
	userCredentialRepository := usercredential.NewRepository(database.DBConn)
	userCredentialService := usercredential.NewService(userCredentialRepository)
	c.add(usercredential.ServiceName, userCredentialService)
}

func (c *container) userService() {
	userCredentialService := fiber.MustGetState[usercredential.Service](c.app.State(), usercredential.ServiceName)
	userRepository := user.NewRepository(database.DBConn)
	userService := user.NewService(userRepository, userCredentialService)
	c.add(user.ServiceName, userService)
}

// func (c *container) funderService() {
// 	userService := fiber.MustGetState[user.Service](c.app.State(), user.ServiceName)
// 	userCredentialService := fiber.MustGetState[usercredential.Service](c.app.State(), usercredential.ServiceName)
// 	funderRepository := funder.NewRepository(database.DBConn)
// 	funderService := funder.NewService(c.gormManagerWithTx(), funderRepository, userService, userCredentialService)
// 	c.add(funder.ServiceName, funderService)
// }

// func (c *container) contractService() {
// 	contractRepository := contract.NewRepository(database.DBConn)
// 	contractService := contract.NewService(c.gormManagerWithTx(), contractRepository)
// 	c.add(contract.ServiceName, contractService)
// }

// func (c *container) contractPaymentService() {
// 	contractPaymentRepository := contractpayment.NewRepository(database.DBConn)
// 	contractService := fiber.MustGetState[contract.Service](c.app.State(), contract.ServiceName)
// 	contractPaymentService := contractpayment.NewService(c.gormManagerWithTx(), contractPaymentRepository, contractService)
// 	c.add(contractpayment.ServiceName, contractPaymentService)
// }

// func (c *container) customerService() {
// 	customerRepository := customer.NewRepository(database.DBConn)
// 	customerService := customer.NewService(c.gormManagerWithTx(), customerRepository)
//

func (c *container) customerService() {
	customerRepository := customer.NewRepository(database.DBConn)
	userService := fiber.MustGetState[user.Service](c.app.State(), user.ServiceName)
	userCredentialService := fiber.MustGetState[usercredential.Service](c.app.State(), usercredential.ServiceName)
	customerService := customer.NewService(c.gormManagerWithTx(), customerRepository, userService, userCredentialService)
	c.add(customer.ServiceName, customerService)
}

// adminService registers the platform-internal-admin module. It depends
// only on the database — it does not need user / user-credential services
// because admin authentication flows through the existing JWT middleware
// (authentication.SuccessHandler), and admin identification happens in
// middleware.AdminOnly via the user_type check on UserSession.
func (c *container) adminService() {
	adminRepository := admin.NewRepository(database.DBConn)
	adminService := admin.NewService(adminRepository)
	c.add(admin.ServiceName, adminService)
}

// func (c *container) merchantService() {
// 	merchantRepository := merchant.NewRepository(database.DBConn)
// 	merchantService := merchant.NewService(c.gormManagerWithTx(), merchantRepository)
// 	c.add(merchant.ServiceName, merchantService)
// }

// func (c *container) categoryService() {
// 	categoryRepository := category.NewRepository(database.DBConn)
// 	categoryService := category.NewService(c.gormManagerWithTx(), categoryRepository)
// 	c.add(category.ServiceName, categoryService)
// }

// func (c *container) accountService() {
// 	accountRepository := account.NewRepository(database.DBConn)
// 	accountService := account.NewService(c.gormManagerWithTx(), accountRepository)
// 	c.add(account.ServiceName, accountService)
// }

// func (c *container) transactionService() {
// 	transactionRepository := transaction.NewRepository(database.DBConn)
// 	transactionService := transaction.NewService(c.gormManagerWithTx(), transactionRepository)
// 	c.add(transaction.ServiceName, transactionService)
// }

// func (c *container) dashboardService() {
// 	transactionRepository := transaction.NewRepository(database.DBConn)
// 	dashboardService := dashboard.NewService(transactionRepository)
// 	c.add(dashboard.ServiceName, dashboardService)
// }

func (c *container) organizationService() {
	organizationRepository := organization.NewRepository(database.DBConn)
	userService := fiber.MustGetState[user.Service](c.app.State(), user.ServiceName)
	userCredentialService := fiber.MustGetState[usercredential.Service](c.app.State(), usercredential.ServiceName)
	organizationService := organization.NewService(organizationRepository, userService, userCredentialService)
	c.add(organization.ServiceName, organizationService)
}

// func (c *container) groqService() {
// 	groqService := groq.NewService(database.DBConn, fiber.MustGetState[category.Service](c.app.State(), category.ServiceName))
// 	c.add(groq.ServiceName, groqService)
// }

// func (c *container) minimaxService() {
// 	minimaxService := minimax.NewService(database.DBConn, fiber.MustGetState[category.Service](c.app.State(), category.ServiceName))
// 	c.add(minimax.ServiceName, minimaxService)
// }

// func (c *container) pdfService() {
// 	pdfService := pdf.NewService(database.DBConn)
// 	c.add(pdf.ServiceName, pdfService)
// }
