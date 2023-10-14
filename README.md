# daily-update

small script/tool I use to produce a daily standup update based on a list of tasks I keep stored in Notion. This repo is likely not useful to you.

## add to path

```sh
#!/bin/sh

cd ${PATH_TO_CLONED_REPO}
NOTION_TOKEN=${NOTION_TOKEN_VAL} NOTION_TASK_DB_ID=${DB_ID_VAL} go run main.go
```

> DB ID can be found in notion page's **URI** (note this is not the `v` query parameter).

## Run it

```sh
daily-update | pbcopy
```
