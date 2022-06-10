## on_machine

A helper CLI tool to determine which computer you're currently on. Often in ones dotfiles or personal scripts, you do something like:

```bash
case "$(uname -s)" in
Linux*) command_on_linux ;;
Darwin*) command_on_mac ;;
*) echo 'unknown...' ;;
esac
```

...to run particular commands on different machines/installs.

This is fine if you're always going to have two installs, but as you add more it starts to get more and more complicated. For example, I have:

- multiple linux installs (for example one on `Ubuntu` and another on `Arch`)
- [termux](https://termux.com/) on my phone
- [`wsl`](https://docs.microsoft.com/en-us/windows/wsl/install) will also return `Linux`, when you likely want to do something custom on windows

So, `on_machine` generates a unique-enough fingerprint of your system (which you can tune to as simple or complicated as you want), so you can do:

```bash
case "$(on_machine)" in
linux_arch_*) command_on_arch ;;
linux_ubuntu_*) command_on_ubuntu ;;
android_termux_*) command_on_termux ;;
windows_*) command_on_wsl ;;
mac_*) command_on_mac ;;
esac
```

This borrows a lot of ideas from tools like [`neofetch`](https://github.com/dylanaraps/neofetch) to figure out what operating system/distribution/window manager one is using

### Install

Using `go install` to put it on your `$GOBIN`:

`go install github.com/seanbreckenridge/on_machine/cmd/on_machine@latest`

I recommend you have both `uname` and `lsb-release` installed if possible on `linux`, that makes distrobution detection much nicer. Otherwise, this defaults to the `golang` `runtime` module defaults

To manually build:

```bash
git clone https://github.com/seanbreckenridge/on_machine
cd ./on_machine
go build ./cmd/on_machine
# copy binary somewhere on your $PATH
sudo cp ./on_machine /usr/local/bin
```

## Usage

```
usage: on_machine [-h] [PATTERN]

Tool to determine which operating system/machine you're on.

PATTERN is a printf-styled format string, supporting the following sequences:

%o - Operating System (using uname)
%d - Distro (using lsb_release)
%h - Hostname (name of the computer)
%a - Arch (detected by golang)
%O - Golang OS (unmodified golang detected operating system)

By default, this uses '%o_%d_%h'
```

## Contributing

I test this on all my machines, but it gets increasingly difficult to test on systems I don't have access to.

If you're able to test this on some operating system not listed above, or use one this doesn't support, am happy to accept PRs for new operating systems, or possible strategies to detect new systems based on some command output/metadata file that exists
