# Workflow Tracker

A system tray application that monitors GitHub workflow status for your open pull requests and sends desktop notifications when workflows complete.

## Overview

Workflow Tracker runs as a background process in your system tray, periodically checking the status of GitHub Actions workflows for all pull requests you have authored. When a workflow completes (success or failure), you receive a desktop notification so you don't need to constantly check GitHub.

## How It Works

### Startup Flow

1. **Configuration Setup**: On first run, the app creates a config file in `~/.wft/config.json` to store your GitHub token and username.

2. **Authentication**: If no token is found, the app initiates GitHub OAuth authentication through your browser. Once authenticated, the token is saved to the config file for future use.

3. **User Information**: The app fetches your GitHub username using the authenticated token and saves it to the config.

4. **State Management**: A state file (`~/.wft/state.json`) is created to track PR details and workflow statuses across sessions.

### Tracking Process

1. **PR Discovery**: Every 2 minutes, the tracker searches for all open pull requests where you are the author using the GitHub Search API.

2. **Workflow Monitoring**: For each PR found:
   - Fetches current PR details using conditional requests (ETags) to minimize API calls
   - Checks if the HEAD commit SHA has changed (new commits)
   - Queries GitHub Actions API for workflow runs associated with the current commit
   - Tracks workflow status: pending, in_progress, completed
   - Records workflow conclusion: success, failure, cancelled, etc.

3. **Change Detection**: The app compares the current state with previously saved state to detect:
   - New commits pushed to PRs
   - Workflow status changes
   - Workflow completions

4. **Notifications**: When a workflow transitions from any state to "completed":
   - Success: Displays a notification with a success icon
   - Failure: Displays a notification with a failure icon

### State Persistence

All PR states are persisted to disk, including:
- PR number and repository
- Current HEAD commit SHA
- Last checked commit SHA
- ETag for conditional requests
- Workflow status and conclusion
- Workflow run ID and URL
- Last update timestamp

This allows the app to resume monitoring seamlessly after restarts without losing context.

### System Tray Interface

The app runs in your system tray with minimal UI:
- Shows "Workflow notifier" title
- Displays a startup notification when launched
- Provides a "Quit" menu option to exit the application

No idea if this works on linux right now

## Running the app

For Linux:
```bash
make linux
```

For Windows:
```bash
make run-win
```

The app will continue running in the background until you quit it from the system tray menu.
