# Go Wallet Watcher
A service built with Golang to watch crypto balance from public wallets (has private key shared publicly) and automate transfers to another wallet, also expose an API endpoint to query wallet info.

## Table of contents
* [Getting Started](#getting-started)
    * [Stack used](#stack-used)
* [Deploying the App](#deploying-the-app)
* [Running Locally](#running-locally)

## Getting Started
To run this service you must provide <b>ETH</b> network (<i>BSC/POLYGON will be implemented soon</i>) wallets address with their respective secret keys.

Also, if you want to record the automated transactions result into BetterStack Log management service you must provide a source token.  

### Stack used
* Golang with concurrency
* Gin web framework
* Go-Ethereum lib

## Deploying the App
You can easily deploy this service to Digital Ocean using their App Platform just clicking in this button.

[![Deploy to DigitalOcean](https://www.deploytodo.com/do-btn-white.svg)](https://cloud.digitalocean.com/apps/new?repo=https://github.com/ronilsonalves/go-wallet-watcher/tree/main&refcode=128ab6cf920e&utm_campaign=Referral_Invite&utm_medium=Referral_Program&utm_source=badge)

Using this button disables the ability to automatically re-deploy your app when pushing to a branch or tag in your repository as you are using this repo directly.

If you want to automatically re-deploy your app, fork the GitHub repository to your account so that you have a copy of it stored to the cloud. Click the Fork button in the GitHub repository and follow the on-screen instructions.

Also, you can build a Docker image and push it to Digital Ocean Container Registry.

<b>Note: Deploy will fail, after you have to configure environment variables to make this service run correctly, check the env var required at .env.example </b>

## Running Locally
If you want to run this service locally, just rename .env.example and fill with your data, open terminal in root project and run:

```cmd
go run cmd/server/main.go
```

## Support
[![DigitalOcean Referral Badge](https://web-platforms.sfo2.digitaloceanspaces.com/WWW/Badge%202.svg)](https://www.digitalocean.com/?refcode=128ab6cf920e&utm_campaign=Referral_Invite&utm_medium=Referral_Program&utm_source=badge)
[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/gbraad)