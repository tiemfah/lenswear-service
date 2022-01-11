package protocol

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/tiemfah/lenswear-service/configs"
	"github.com/tiemfah/lenswear-service/internal/core/services/apparelsrv"
	"github.com/tiemfah/lenswear-service/internal/core/services/authsrv"
	"github.com/tiemfah/lenswear-service/internal/core/services/usersrv"
	"github.com/tiemfah/lenswear-service/internal/handlers/apparelhdl"
	"github.com/tiemfah/lenswear-service/internal/handlers/authhdl"
	"github.com/tiemfah/lenswear-service/internal/handlers/userhdl"
	"github.com/tiemfah/lenswear-service/internal/repositories/apparelrepo"
	"github.com/tiemfah/lenswear-service/internal/repositories/authrepo"
	"github.com/tiemfah/lenswear-service/internal/repositories/userrepo"
	"github.com/tiemfah/lenswear-service/pkg/bucket"
	"github.com/tiemfah/lenswear-service/pkg/hash"
	"github.com/tiemfah/lenswear-service/pkg/middlewares"
	ttjwt "github.com/tiemfah/lenswear-service/pkg/token/jwt"
	"github.com/tiemfah/lenswear-service/pkg/token/rsa"
	"github.com/tiemfah/lenswear-service/pkg/uidgen"
	"github.com/tiemfah/lenswear-service/protocol/adminroutes"
	"github.com/tiemfah/lenswear-service/protocol/userroutes"
	"google.golang.org/api/option"
)

func ServeHTTP() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var cfg struct{ Env string }
	flag.StringVar(&cfg.Env, "env", "", "the environment and config that used")
	flag.Parse()
	configs.InitViper("./configs", cfg.Env)

	dsClient, err := datastore.NewClient(ctx, configs.GetViper().GCP.ProjectID, option.WithCredentialsFile(configs.GetViper().GCP.ServiceAccount))
	if err != nil {
		return fmt.Errorf("cannot connect to the datastore client: '%s'", err)
	}
	defer dsClient.Close()

	csClient, err := storage.NewClient(ctx, option.WithCredentialsFile(configs.GetViper().GCP.ServiceAccount))
	if err != nil {
		return fmt.Errorf("cannot connect to the cloud storage client: '%s'", err)
	}
	defer csClient.Close()
	gcpBucket := bucket.NewGCPBucket(ctx, csClient, configs.GetViper().GCP.ProjectID)

	uidgen := uidgen.New()
	hash := hash.New()

	privateKey, publicKey := rsa.GenerateRSA(configs.GetViper().KeyPath.PublicKey, configs.GetViper().KeyPath.PrivateKey)

	authRepository := authrepo.NewAuthenticationRepository(ttjwt.New(privateKey, publicKey))
	userRepository := userrepo.NewUserRepository(dsClient, hash)
	apparelRepository := apparelrepo.NewApprelRepository(dsClient, gcpBucket)

	authService := authsrv.NewAuthenticationService(authRepository, userRepository)
	userService := usersrv.NewUserService(userRepository, uidgen, hash)
	apparelService := apparelsrv.NewApparelService(apparelRepository, uidgen)

	authHandler := authhdl.NewHTTPHandler(authService)
	userHandler := userhdl.NewHTTPHandler(userService)
	apparelHandler := apparelhdl.NewHTTPHandler(apparelService)

	app := fiber.New()
	api := app.Group("api", middlewares.CORSMiddleware(), middlewares.LoggerMiddleware())
	{
		adminGroup := api.Group("a")
		{
			authAPI := adminGroup.Group("auth")
			adminroutes.AuthEndPoint(authAPI, authHandler)
			userAPI := adminGroup.Group("user", middlewares.AuthMiddleware(publicKey))
			adminroutes.UserEndPoint(userAPI, userHandler)
			apparelAPI := adminGroup.Group("apparel", middlewares.AuthMiddleware(publicKey))
			adminroutes.ApparelEndPoint(apparelAPI, apparelHandler)
		}
		userGroup := api.Group("u")
		{
			authAPI := userGroup.Group("auth")
			userroutes.AuthEndPoint(authAPI, authHandler)
			userAPI := userGroup.Group("user", middlewares.AuthMiddleware(publicKey))
			userroutes.UserEndPoint(userAPI, userHandler)
			apparelAPI := userGroup.Group("apparel")
			userroutes.ApparelEndPoint(apparelAPI, apparelHandler)
		}
	}
	app.Get("health", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	if err := app.Listen(":" + configs.GetViper().HTTPPort); err != nil {
		return err
	}
	return nil
}
