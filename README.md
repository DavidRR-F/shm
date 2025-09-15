# SHM (Simple Home Manager)

Simple dotfile manager and package installer for linux/mac

## Synopsis

`Stow` is great, but it doesn't help you manager your user/system packages and `Nix Home Manager` does but, if you are not using `NixOS`,
the management plan becomes overly complex. I just wanted a simple tool that will link my dotfiles and install any package I want from
any package manager I want with a single command so that setting up my development environment on new devices/vm is one command regardless
of the platform I am using. Enter `shm`, a simple cil tool that uses yaml configurations to initialize your workstations.

## Installation

```bash
curl -sSfL https://raw.githubusercontent.com/DavidRR-F/shm/main/install.sh | bash
```

## Usage

### Links

`.shm/shm.yml`

```yaml
links:
  - src: ~/.dotfiles/.bashrc 
    dest: ~/.bashrc
  - src: ~/.dotfiles/myscript.sh 
    dest: ~/.local/bin/myscript
    exe: true
```

```bash
shm ~/dotfiles/path
```

### Packages Managers

[Manager Configuration Examples](https://github.com/DavidRR-F/shm/tree/main/.shm/managers)

`.shm/managers/nix-env`

```yaml
# package manager name
name: nix-env
command: ["nix-env"]
# arguments for install command
args: ["-iA"]
# if will be invoked if package manager not already installed
install:
  # package manager installation command
  pre: ["sh", "<(curl -L https://nixos.org/nix/install)", "--daemon"]
  # commands to run before installing packages
  post: [".", "/etc/profile.d/nix.sh"]
# list of packages
packages:
  - nixpkgs.starship
```

`.shm/base.yml`

```yaml
links:
  ...

# list of package managers
managers:
  - nix-env
```

```bash
shm ~/dotfiles/path --install-packages
```

### Profiles

`.shm/mac.yml`

```yaml
links:
  - src: ~/.dotfiles/lazygit
    dest: ~/Applications/Library Support/lazygit

managers:
  - brew
```

```bash
shm ~/dotfiles/path --profile mac
```

### Development

copy the examples config to as a yhm in the project base directory

```bash
cp -R examples/config .shm
```

run the following to invoke the cli from source

```bash
go run ./cmd/yhm/main.go . --dry-run
```
