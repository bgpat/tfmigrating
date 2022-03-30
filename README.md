# tfmigrating

Generate a migration file for [tfmigrate](https://github.com/minamijoyo/tfmigrate) from tfplan.

## Features

- migration block types
  - [x] `state`
  - [ ] `multi_state`
- migration actions
  - [x] `mv`
  - [ ] `rm`
  - [ ] `import`
- migration options
  - [ ] set custom migration name
  - [ ] set `dir`
  - [ ] set `workspace`
  - [ ] set `force`
- others
  - [ ] write migration with default filename
  - [ ] prettify migration output

## Install

```bash
go install github.com/bgpat/tfmigrating@latest
```

## Usage

```bash
terraform plan -out=tfplan
terraform show -json tfplan | tfmigrating > migration.hcl
```

