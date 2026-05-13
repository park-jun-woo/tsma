# TestMaster

Extract tests from legacy code with LLM agents.

## Install

```bash
make install
```

## Usage

```
tsma next                          # show next incomplete endpoint + call chain
tsma submit <endpoint> <test-file> # submit a test file for verification
tsma list [--page N] [--size S]    # list all endpoints with status
tsma status [--endpoint E]         # show progress summary or endpoint detail
tsma reset [<endpoint> | --all]    # reset endpoint or delete session
```

## Workflow

1. `tsma next` -- get the next endpoint and its function call chain
2. Read the chain functions, write a test covering all branches
3. `tsma submit <endpoint> <test-file>` -- submit for validation
4. If PARTIAL, check uncovered branches and improve the test
5. Repeat until all endpoints are DONE
