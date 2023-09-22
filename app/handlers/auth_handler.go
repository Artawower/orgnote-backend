package handlers

import (
	"bytes"
	"encoding/gob"
	"net/url"
	"orgnote/app/configs"
	"orgnote/app/models"
	"orgnote/app/services"
	"orgnote/app/tools"

	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/rs/zerolog/log"
	"github.com/shareed2k/goth_fiber"
)

type OAuthRedirectData struct {
	RedirectURL string `json:"redirectUrl"`
}

func mapToUser(user goth.User) *models.User {
	return &models.User{
		Provider:            user.Provider,
		Email:               user.Email,
		Name:                user.Name,
		NickName:            user.NickName,
		AvatarURL:           user.AvatarURL,
		ExternalID:          user.UserID,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		Token:               user.AccessToken,
		RefreshToken:        &user.RefreshToken,
		TokenExpirationDate: user.ExpiresAt,
		ProfileURL:          user.RawData["html_url"].(string),
		Notes:               []models.Note{},
		APITokens:           []models.APIToken{},
		SpaceLimit:          tools.ConvertMegabyte2Bytes(100),
		UsedSpace:           0,
	}
}

type AuthHandler struct {
	userService    *services.UserService
	config         configs.Config
	authMiddleware fiber.Handler
}

// Login godoc
// @Summary      OAuth Login
// @Description  Entrypoint for login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        provider path string true "provider"
// @Success      200  {object}  handlers.HttpResponse[OAuthRedirectData, any]
// @Failure      400  {object}  handlers.HttpError[any]
// @Failure      404  {object}  handlers.HttpError[any]
// @Failure      500  {object}  handlers.HttpError[any]
// @Router       /auth/{provider}/login  [get]
func (a *AuthHandler) Login(c *fiber.Ctx) error {
	log.Info().Msgf("Fiber context: %v", c.BaseURL())
	url, err := goth_fiber.GetAuthURL(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	// return c.Redirect(url, fiber.StatusTemporaryRedirect)
	log.Info().Msgf("Redirecting to %s", url)
	data := NewHttpResponse[OAuthRedirectData, any](OAuthRedirectData{
		RedirectURL: url,
	}, nil)
	return c.JSON(data)
}

// LoginCallback godoc
// @Summary      Callback for OAuth
// @Description
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        provider path string true "provider"
// @Success      200  {object}  any
// @Failure      400  {object}  handlers.HttpError[any]
// @Failure      404  {object}  handlers.HttpError[any]
// @Failure      500  {object}  handlers.HttpError[any]
// @Router       /auth/{provider}/callback  [get]
func (a *AuthHandler) LoginCallback(c *fiber.Ctx) error {
	user, err := goth_fiber.CompleteUserAuth(c)
	if err != nil {
		log.Error().Err(err).Msgf("auth handlers: github auth handler: complete user auth")
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}
	var userBytes bytes.Buffer
	enc := gob.NewEncoder(&userBytes)
	err = enc.Encode(user)
	if err != nil {
		log.Error().Err(err).Msgf("auth handlers: github auth handler: encode user: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}
	u, err := a.userService.Login(*mapToUser(user))
	if err != nil {
		log.Error().Err(err).Msgf("auth handlers: github auth handler: login user %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}
	// TODO: master client url for redirect. Read from env
	redirectURL := a.config.ClientAddress + "/#/auth/login/"
	parsedURL, err := url.Parse(redirectURL)
	if err != nil {
		log.Error().Err(err).Msgf("auth handlers: github auth handler: parse redirect url %v", err)
	}
	q := parsedURL.Query()
	q.Set("token", u.Token)
	q.Set("id", u.ID.Hex())
	q.Set("username", u.NickName)
	q.Set("avatarUrl", u.AvatarURL)
	q.Set("email", u.Email)
	q.Set("profileUrl", u.ProfileURL)
	parsedURL.RawQuery = q.Encode()

	return c.Redirect(redirectURL + "?" + parsedURL.RawQuery)

}

// Logout godoc
// @Summary      Logout
// @Description
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  any
// @Failure      500  {object}  handlers.HttpError[any]
// @Router       /auth/logout  [get]
func (a *AuthHandler) Logout(c *fiber.Ctx) error {
	if err := goth_fiber.Logout(c); err != nil {
		log.Error().Err(err).Msgf("auth handlers: github auth handler: logout")
		return c.Status(500).SendString("Internal server error")
	}
	// TODO: master delete user token here
	c.SendString("logout")
	return c.Status(200).JSON(struct{}{})
}

// CreateApiToken godoc
// @Summary      Create API token
// @Description  Create API token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  HttpResponse[models.APIToken, any]
// @Failure      500  {object}  handlers.HttpError[any]
// @Router       /auth/token  [post]
func (a *AuthHandler) CreateToken(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	token, err := a.userService.CreateToken(user)
	if err != nil {
		log.Error().Err(err).Msgf("auth handlers: github auth handler: create token")
		return c.Status(500).SendString("Internal server error")
	}
	return c.Status(200).JSON(NewHttpResponse[*models.APIToken, any](token, nil))
}

type ParamsDeleteToken struct {
	TokenID string `json:"tokenId"`
}

// DeleteToken   godoc
// @Summary      Delete API token
// @Description  Delete API token
// @Tags         auth
// @Param 		   tokenId path string true "token id"
// @Accept       json
// @Produce      json
// @Success      200  {object}  any
// @Failure      500  {object}  handlers.HttpError[any]
// @Router       /auth/token/{tokenId}  [delete]
func (a *AuthHandler) DeleteToken(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	b := new(ParamsDeleteToken)
	if err := c.ParamsParser(b); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(NewHttpError[any]("Token doesn't provided", nil))
	}

	log.Info().Msgf("Delete token %s", b.TokenID)

	err := a.userService.DeleteToken(user, b.TokenID)
	if err != nil {
		log.Error().Err(err).Msgf("auth handlers: github auth handler: delete token")
		return c.Status(500).SendString("Internal server error")
	}
	return c.Status(200).JSON(NewHttpResponse[any, any](nil, nil))
}

// VerifyUser godoc
// @Summary      Verify user
// @Description  Return found user by provided token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  handlers.HttpResponse[models.PublicUser, any]
// @Failure      403  {object}  handlers.HttpError[any]
// @Failure      500  {object}  handlers.HttpError[any]
// @Router       /auth/verify  [get]
func (a *AuthHandler) VerifyUser(c *fiber.Ctx) error {
	token := tools.ExtractBearerTokenFromCtx(c)
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(NewHttpError[any](ErrTokenNotProvided, nil))
	}
	user, err := a.userService.FindUser(token)
	if err != nil {
		log.Info().Err(err).Msgf("auth handlers: github auth handler: find user")
		return c.Status(fiber.StatusBadRequest).SendString(ErrInvalidToken)
	}
	return c.Status(fiber.StatusOK).JSON(NewHttpResponse[*models.PublicUser, any](user, nil))
}

// GetApiTokens godoc
// @Summary      Get API tokens
// @Description  Return all available API tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  handlers.HttpResponse[[]models.APIToken, any]
// @Failure      500  {object}  handlers.HttpError[any]
// @Failure      400  {object}  handlers.HttpError[any]
// @Router       /auth/api-tokens  [get]
func (a *AuthHandler) GetAPITokens(c *fiber.Ctx) error {
	ctxUser := c.Locals("user")
	if ctxUser == (*models.User)(nil) {
		return c.Status(fiber.StatusBadRequest).SendString("Could not find api tokens for current user")
	}
	user := c.Locals("user").(*models.User)
	tokens, err := a.userService.GetAPITokens(user.ID.Hex())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Could not find api tokens for current user")
	}
	return c.Status(fiber.StatusOK).JSON(NewHttpResponse[[]models.APIToken, any](tokens, nil))
}

// TODO: master refactor this code.
func RegisterAuthHandler(app fiber.Router, userService *services.UserService, config configs.Config, authMiddleware fiber.Handler) {
	goth.UseProviders(
		github.New(config.GithubID, config.GithubSecret, config.BackendHost()+"/auth/github/callback"),
	)

	authHandler := &AuthHandler{
		userService:    userService,
		config:         config,
		authMiddleware: authMiddleware,
	}

	app.Get("/auth/:provider/login", authHandler.Login)
	app.Get("/auth/:provider/callback", authHandler.LoginCallback)
	app.Get("/auth/logout", authHandler.Logout)
	app.Post("/auth/token", authMiddleware, authHandler.CreateToken)
	app.Delete("/auth/token/:tokenId", authMiddleware, authHandler.DeleteToken)
	app.Get("/auth/verify", authHandler.VerifyUser)
	app.Get("/auth/api-tokens", authHandler.GetAPITokens)
}
