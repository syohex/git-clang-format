# git-clang-format

[clang-format](https://clang.llvm.org/docs/ClangFormat.html) utility which applies C family files(C, C++, Objective-C, Objective-C++)
in git repository

## Usage

```bash
# apply clang-format for C family files under current directory
% git-clang-format

# apply clang-format for C family files under some_dir
% git-clang-format some_dir

# apply all C family files in repository
% git-clang-format -all

# apply only modified files C family files in repository
% git-clang-format -modified

# apply only staged files C family files in repository
% git-clang-format -staged
```
