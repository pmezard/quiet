# What is it?

crontab sends emails when a job prints anything on standard output. I find that
annoying at best and much prefer to receive emails about failing jobs, and
possibly with their output tailed. That is what `quiet` is about:
```
quiet sh -c 'echo something; false'
```
executes `sh -c 'echo something; false'` and because the command fails, prints
`echo` output and forward the exit status code. It would have printed nothing
if `false` where to be removed or replaced by `true`.

If you know a simple standard way to achieve this, please tell me!
