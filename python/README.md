# Python Example
## Interacting with the Github API

### Usage

- `export GHTOKEN=<your github token>` 
- `./githubApiInteraction.py`

#### Redirect output to a file

- `githubApiInteraction.py &> authors.txt`

#### Redirect output to a file and view results

- `githubApiInteraction.py 2>&1 | tee [-a] authors.txt`
