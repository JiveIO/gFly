package gfly

import (
	"app/core/oauth"
	"app/core/oauth/providers/facebook"
	"app/core/oauth/providers/github"
	"app/core/oauth/providers/google"
	"app/core/utils"
	"os"
)

func setupOAuth() {
	loadOAuthFacebook()
	loadOAuthGoogle()
	loadOAuthGithub()
}

func loadOAuthFacebook() {
	key := utils.Getenv("OAUTH_FACEBOOK_KEY", "")
	secret := utils.Getenv("OAUTH_FACEBOOK_SECRET", "")
	callback := utils.Getenv("OAUTH_FACEBOOK_CALLBACK", "")

	if key != "" && secret != "" && callback != "" {
		oauth.UseProviders(
			facebook.New(os.Getenv("OAUTH_FACEBOOK_KEY"), os.Getenv("OAUTH_FACEBOOK_SECRET"), os.Getenv("OAUTH_FACEBOOK_CALLBACK")),
		)
	}
}

func loadOAuthGoogle() {
	key := utils.Getenv("OAUTH_GOOGLE_KEY", "")
	secret := utils.Getenv("OAUTH_GOOGLE_SECRET", "")
	callback := utils.Getenv("OAUTH_GOOGLE_CALLBACK", "")

	if key != "" && secret != "" && callback != "" {
		oauth.UseProviders(
			google.New(os.Getenv("OAUTH_GOOGLE_KEY"), os.Getenv("OAUTH_GOOGLE_SECRET"), os.Getenv("OAUTH_GOOGLE_CALLBACK")),
		)
	}
}

func loadOAuthGithub() {
	key := utils.Getenv("OAUTH_GITHUB_KEY", "")
	secret := utils.Getenv("OAUTH_GITHUB_SECRET", "")
	callback := utils.Getenv("OAUTH_GITHUB_CALLBACK", "")

	if key != "" && secret != "" && callback != "" {
		oauth.UseProviders(
			github.New(os.Getenv("OAUTH_GITHUB_KEY"), os.Getenv("OAUTH_GITHUB_SECRET"), os.Getenv("OAUTH_GITHUB_CALLBACK")),
		)
	}
}
