# 📦 SHM (Simple Home Manager)

**A simple dotfile manager and package manager wrapper for Linux/macOS**

Define symbolic links for multiple os profiles, and install packages from multiple package managers like

| Package Manager | Icon |
|-----------------|:----:|
| **nix-env**     | 🧪   |
| **brew**        | 🍺   |
| **apt**         | 🐧   |
| **dnf**         | 🔧   |
| **flatpak**     | 📦   |
| **snap**        | 📦   |

> [!WARN] still testing package manager wrapper configurations may not work will all managers
---

## 📌 Synopsis

`stow` is great for symlinking dotfiles, but it doesn't help manage your system or user packages.  
Meanwhile, `Nix Home Manager` does this — but if you're not on NixOS, things can get overly complex.

I just wanted a **simple tool** to:
- ✅ Link dotfiles
- ✅ Install packages from **any** package manager

So I made **`shm`** — a lightweight CLI tool that uses **YAML** configuration to initialize your dev environment in seconds.
---

## ⚙️ Installation

Run the following installation script

```bash
curl -sSfL https://raw.githubusercontent.com/DavidRR-F/shm/main/install.sh | bash
```

> [!HINT] This script installs the executable in your `~/.local/bin` by default to remain user scoped

## 💡Usage

For examples, you can view my personal [shm configurations](https://github.com/DavidRR-F/dotfiles/tree/main/.shm)

### 🗂️File Structure

Create your shm configuration in your dotfiles directory

```bash
.shm/
├── base.yml
├── <profile>.yml
└── managers/
    └── <manager>.yml
```

| File/Directory | Description |
|:---------------|:------------|
| **.shm/base.yml** | base shm configuration for common configurations between os's or single configuration (Can be empty) |
| **.shm/<profile>.yml** | additional profiles configurations that can be added to base configuration |
| **.shm/managers/<manager>.yml** | manager object configuration defines for package managers |

### 🔗 Dotfile Linking

You can define symbolic link definitions in your `base.yml` or `<profile>.yml`

```yaml
links:
  - src: ~/.dotfiles/.bashrc 
    dest: ~/.bashrc
  - src: ~/.dotfiles/myscript.sh 
    dest: ~/.local/bin/myscript
    exe: true # makes file execute
```

Simply run shm on your dotfile base directory to apply the links in your `base.yml`

```bash
shm ~/dotfiles/path
```

### 📦 Packages Managers

In your managers directory you can define manager yaml object

```yaml
# package manager name (must match file name and be executable)
name: nix-env
command: ["nix-env"]
# arguments for install command
args: ["-iA"]
# Nested commands will be invoked if package manager not already installed
install:
  # package manager installation command
  pre: ["sh", "<(curl -L https://nixos.org/nix/install)", "--daemon"]
  # commands to run after package manager installation
  post: [".", "/etc/profile.d/nix.sh"]
# list of packages to install
packages:
  - nixpkgs.starship
  - nixpkgs.lazygit
```

You can then add the manager(s) references to your `base.yml` or `<profile>.yml` configurations

```yaml
managers:
  - nix-env
```

By default, only links are applied when running shm, add the `install-packages` flag to additionally defined install packages

```bash
shm ~/dotfiles/path --install-packages
```

### 🧩 Profiles

You are able do define profile files in your `.shm` configuration

Example: `.shm/mac.yml`

```yaml
links:
  - src: ~/.dotfiles/lazygit
    dest: ~/Applications/Library Support/lazygit

managers:
  - brew
```

When running *shm* you can specify a profile and those configuration will be added to your *base* configuration

```bash
shm ~/dotfiles/path --profile mac --install-packages
```

### 🧪 Dry Run

You can preview what shm will do without making changes using `--dry-run`.

```bash
shm ~/dotfiles/path -i -p mac --dry-run
```
