# GitHub Authentication

vulhub-cli downloads environment files from GitHub. To avoid API rate limits, you can authenticate with your GitHub account.

## Rate Limits

| Authentication Status | Rate Limit |
|----------------------|------------|
| Unauthenticated | 60 requests/hour |
| Authenticated | 5,000 requests/hour |

For typical usage, 60 requests/hour is often sufficient. However, if you're:
- Downloading many environments
- Running frequent syncs
- Using vulhub-cli in CI/CD pipelines

You may want to authenticate to avoid hitting rate limits.

## OAuth Device Flow

vulhub-cli uses GitHub's OAuth Device Flow for authentication. This is a secure method that:
- Does not require you to enter your password in the CLI
- Does not store your GitHub password
- Only requests minimal permissions (public repository access)
- Can be revoked at any time from GitHub settings

## Authentication Process

### Step 1: Start Authentication

```bash
vulhub github-auth
```

### Step 2: Authorization

The CLI will display:
```
GitHub Authentication
Authenticate with GitHub to increase API rate limit from 60 to 5,000 requests/hour.

ℹ Requesting authorization code...

Please visit the URL below and enter the code:
  URL:  https://github.com/login/device
  Code: ABCD-1234

ℹ Browser opened automatically
⠋ Waiting for authorization...
```

### Step 3: Enter Code on GitHub

1. The browser opens automatically (or manually visit the URL)
2. Log in to GitHub if not already logged in
3. Enter the code displayed in the CLI
4. Click "Authorize"

### Step 4: Completion

Once authorized, the CLI displays:
```
✓ Authentication successful!
API rate limit increased to 5,000 requests/hour
```

The token is automatically saved to your configuration file.

## Managing Authentication

### Check Status

View your current authentication status:

```bash
vulhub github-auth --status
```

Output when authenticated:
```
● Authenticated
  Token: ghp_****····****
ℹ API rate limit: 5,000 requests/hour
```

Output when not authenticated:
```
○ Not authenticated
ℹ API rate limit: 60 requests/hour (unauthenticated)

ℹ Run 'vulhub github-auth' to authenticate.
```

### Remove Authentication

To remove your saved authentication:

```bash
vulhub github-auth --remove
```

You will be prompted to confirm:
```
Remove GitHub authentication?
You will need to re-authenticate to avoid rate limits.
  [Yes, remove]  [Cancel]
```

## Automatic Authentication Prompt

When vulhub-cli encounters a rate limit error, it automatically prompts you to authenticate:

```bash
vulhub start log4j
# If rate limited:
# ⚠ GitHub API rate limit exceeded
# Authenticate to increase limit from 60 to 5,000 requests/hour
#   Authenticate with GitHub now? [Yes] [Later]
```

If you choose "Yes":
1. The OAuth flow starts immediately
2. After successful authentication, the original operation automatically retries
3. No need to re-run the command

## Using Environment Variable

You can also provide a GitHub token via environment variable:

```bash
export GITHUB_TOKEN=ghp_your_personal_access_token
vulhub start log4j
```

The environment variable takes precedence over the saved token in config.toml.

### Creating a Personal Access Token

If you prefer to use a Personal Access Token instead of OAuth:

1. Go to GitHub Settings → Developer settings → Personal access tokens
2. Click "Generate new token (classic)"
3. Select scope: `public_repo` (or no scopes for public repositories)
4. Generate and copy the token
5. Set the environment variable or add to config.toml

## Token Storage

The OAuth token is stored in `~/.vulhub/config.toml`:

```toml
[github]
token = "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
```

### Security Notes

- The token file is created with user-only permissions (0600)
- The token only has access to public repositories
- You can revoke the token at any time from GitHub settings
- The token is masked when displayed (`ghp_****····****`)

## Revoking Access

To revoke vulhub-cli's access to your GitHub account:

1. Go to GitHub Settings → Applications → Authorized OAuth Apps
2. Find "Vulhub CLI"
3. Click "Revoke"

Or from the CLI:
```bash
vulhub github-auth --remove
```

## Troubleshooting

### "Authorization timed out"

The authorization code expires after 15 minutes. If you see this error:
1. Run `vulhub github-auth` again
2. Complete the authorization within 15 minutes

### "Authorization was denied"

If you clicked "Cancel" or "Deny" on the GitHub authorization page:
1. Run `vulhub github-auth` again
2. Click "Authorize" when prompted

### Rate limit still exceeded after authentication

1. Check if the token is saved:
   ```bash
   vulhub github-auth --status
   ```

2. Verify the token works:
   ```bash
   curl -H "Authorization: token $(grep token ~/.vulhub/config.toml | cut -d'"' -f2)" \
        https://api.github.com/rate_limit
   ```

3. Re-authenticate:
   ```bash
   vulhub github-auth --remove
   vulhub github-auth
   ```

### Environment variable not recognized

Ensure the variable is exported:
```bash
# Wrong
GITHUB_TOKEN=xxx vulhub start log4j

# Correct
export GITHUB_TOKEN=xxx
vulhub start log4j

# Or in one line
GITHUB_TOKEN=xxx vulhub start log4j  # Works in most shells
```
