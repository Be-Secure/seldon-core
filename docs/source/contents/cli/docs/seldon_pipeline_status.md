## seldon pipeline status

status of a pipeline

### Synopsis

status of a pipeline

```
seldon pipeline status [flags]
```

### Options

```
  -h, --help                    help for status
  -p, --pipeline-name string    pipeline name for status
      --scheduler-host string   seldon scheduler host (default "0.0.0.0")
      --scheduler-port int      seldon scheduler port (default 9004)
  -w, --wait string             pipeline wait condition
```

### Options inherited from parent commands

```
  -r, --show-request    show request
  -o, --show-response   show response (default true)
```

### SEE ALSO

* [seldon pipeline](seldon_pipeline.md)	 - manage pipelines

###### Auto generated by spf13/cobra on 16-Apr-2022