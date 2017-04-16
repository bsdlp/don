# don

Manages your environment variables to set credentials according to profiles
configured in your `$HOME/.aws/credentials` file.

## Usage

Prints the shell commands exporting this profile's aws credentials to the
environment.

Usage:

```
> eval $(don personal)
```

Example:

```
> don personal
export AWS_ACCESS_KEY_ID="AKIDWEOIFJWOPEJF"
export AWS_SECRET_ACCESS_KEY=""weijpiojwfpaoewifjawfoipj
```
