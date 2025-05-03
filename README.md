# PTC
A easy to use patch CLI built for the [ThreadMC](https://github.com/threadmc) open-source projects

## Commands

- `ptc create <file> <patch>`: Create a patch from a file.
- `ptc apply <patch> [--no-hash-check]`: Apply a patch to the target file. Use `--no-hash-check` to skip hash verification.
- `ptc list <directory>`: List all patch files in a directory.
- `ptc diff <patch>`: Show the diff between the current file and the patch content.
- `ptc backup <file>`: Create a backup of a file.
- `ptc restore <file>`: Restore a file from its backup.

## Example

```sh
ptc create foo.txt foo.ptc
ptc apply foo.ptc
ptc diff foo.ptc
ptc backup foo.txt
ptc restore foo.txt
```
