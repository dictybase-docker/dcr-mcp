---
allowed-tools: Bash(git worktree:*), Bash(git status:*), Bash(basename:*), Bash(openssl:*), Bash(echo:*), Bash(git rev-parse:*), Bash(git branch:*), Bash(ls:*), Bash(pwd:*), Bash(tr:*), Bash(cd:*), TodoWrite
argument-hint: [branch-name]
description: Create a git worktree from current or specified branch
model: claude-3-5-haiku@20241022 
---

I'll help you create a git worktree from the current branch. Let me first create
a TodoList to track this task.

Steps to create the worktree:
1. Parse current repository name and branch information
2. Format the worktree path according to: {repo-name}-{branch-name}-{random-string}-claude
3. Generate a random alphanumeric string using openssl
4. Check if the target branch is already checked out
5. Create the git worktree with appropriate strategy
6. Verify the worktree creation

The worktree will be created one directory up from the current repository.

**Important Git Worktree Behavior:**
- If the target branch is already checked out in the main repository, Git will
not allow creating a worktree from the same branch
- In this case, we must create a new branch from the current branch using: `git
worktree add -b <new-branch-name> <path>`
- If the branch is not checked out elsewhere, we can use: `git worktree add
<path> <branch-name>`

**Command Strategy:**
1. First try: `git worktree add ../repo-branch-RANDOM-claude branch-name`
2. If that fails with "already checked out" error, use: `git worktree add -b branch-name-worktree ../repo-branch-RANDOM-claude`

Format details:
- Repository name: basename of current directory
- Branch name: current branch (or $1 if provided) with forward slashes converted to hyphens
- Random string: 6-character uppercase hex string
- Final format: repo-branch-RANDOM-claude
- New branch name (if needed): original-branch-name-worktree
