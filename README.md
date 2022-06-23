## on_machine

A helper CLI tool to determine which computer you're currently on. Often in ones dotfiles or personal scripts, you do something like:

```bash
case "$(uname -s)" in
  Linux) command_on_linux ;;
  Darwin) command_on_mac ;;
  *) echo 'unknown...' ;;
esac
```

...to run particular commands on different machines/installs.

This is fine if you're always going to have two installs, but as you add more it starts to get more and more complicated. For example, I have:

- multiple linux installs (for example one on `Ubuntu` and another on `Arch`)
- [termux](https://termux.com/) on my phone
- [`wsl`](https://docs.microsoft.com/en-us/windows/wsl/install) will also return `Linux`, when you likely want to do something custom on windows

So, `on_machine` generates a unique-enough fingerprint of your system (which you can tune to be as simple or complicated as you want), so you can do:

```bash
case "$(on_machine)" in
  linux_arch_*) command_on_arch ;;
  linux_ubuntu_*) command_on_ubuntu ;;
  android_termux_*) command_on_termux ;;
  windows_*) command_on_wsl ;;
  mac_*) command_on_mac ;;
esac
```

This borrows a lot of ideas from tools like [`neofetch`](https://github.com/dylanaraps/neofetch) to figure out what operating system/distribution one is using

### Install

Using `go install` to put it on your `$GOBIN`:

`go install github.com/seanbreckenridge/on_machine/cmd/on_machine@latest`

I recommend you have both `uname` and `lsb_release` installed if possible on `linux`, that makes distribution detection much nicer. Otherwise, this defaults to the `golang` `runtime` module defaults

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
%d - Distro (using lsb_release, falling back to uname if lsb_release isnt available)
%h - Hostname (name of the computer)
%a - Arch (detected by golang)
%O - Golang OS (unmodified golang detected operating system)

By default, this uses '%o_%d_%h'
```

### `-match-paths`

This can be thought of as an alternative to the above, its a different way to figure out what code to run on different machines, by placing scripts in a particular directory structure

You can manually do case/regex statements in bash (and often that is enough), but in some cases that can become complicated. For example, I use this for my [background processes](https://github.com/seanbreckenridge/bgproc), as I organize each job as a small bash script. To figure out which bash scripts to run, I organize my scripts into a directory like:

```
matching_examples/dir_based
├── all
│   └── this_dir_is_matched_everytime
├── android
│   └── on_android
├── linux
│   ├── arch
│   │   ├── home
│   │   │   └── on_arch_when_hostname_is_home
│   │   ├── only_on_arch
│   │   └── work
│   │       └── on_arch_when_hostname_is_work
│   ├── matched_on_any_linux_install
│   └── ubuntu
│       └── only_on_ubuntu
└── mac
    └── on_mac
```

Then, say my hostname is `home`, and I'm on my `arch` machine. `-match-paths` computes the following:

```bash
$ on_machine -match-paths ./matching_examples/dir_based '%o/%d/%h'
/home/sean/Repos/on_machine/matching_examples/dir_based/all
/home/sean/Repos/on_machine/matching_examples/dir_based/linux
/home/sean/Repos/on_machine/matching_examples/dir_based/linux/arch
/home/sean/Repos/on_machine/matching_examples/dir_based/linux/arch/home
```

Note: `all` is like a `*`, its always matched -- so that would be where I store shared jobs.

Then, I could use `find` to loop over each directory, searching with `-maxdepth 1` to find all my jobs that match my current OS/distro:

TODO: ADD LINK TO bgproc_on_machine

(For a real example, see [here](https://github.com/seanbreckenridge/HPI-personal/tree/master/jobs))

```
TODO: ADD EXAMPLE USING -PRINT0 FLAG TO LOOP OVER RESULTS
```

If the pattern includes an extension, this extracts that and tries to match at each level going down. For example, I use this in my `~/.zshrc` setup. I want some code that runs everywhere (`all.zsh`), some that runs on `android`, some that runs on `linux`, and then additional code that runs on `linux` and `arch`. So, given:

```
./matching_examples/with_extensions
├── all.zsh
├── android.zsh
├── linux
│   ├── arch.zsh
│   └── ubuntu.zsh
├── linux.zsh
└── mac.zsh
```

On my arch machine, using the pattern `%o/%d.zsh`, this matches:

```
$ on_machine -match-paths ./matching_examples/with_extensions '%o/%d.zsh
/home/sean/Repos/on_machine/matching_examples/with_extensions/all.zsh
/home/sean/Repos/on_machine/matching_examples/with_extensions/linux
/home/sean/Repos/on_machine/matching_examples/with_extensions/linux/arch.zsh
/home/sean/Repos/on_machine/matching_examples/with_extensions/linux.zsh
```

That does (purposefully, incase you want it) include the `linux` directory, but that can just be filtered later in bash with a quick `[[ -f "$pth" ]]` check before sourcing each `.zsh` file.

Essentially, this lets me pick what scripts to run organized as a directory, instead of ever-growing `case` statement in a `bash` script somewhere.

## Contributing

I test this on all my machines, but it gets increasingly difficult to test on systems I don't have access to.

If you're able to test this on some operating system not listed above, or use one this doesn't support, am happy to accept PRs for new operating systems, or possible strategies to detect new systems based on some command output/metadata file that exists
