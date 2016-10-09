# templenv
templenv reads golang template files and outputs the result to stdout. From there you can pipe it wherever you want.

there are some (only two atm) useful functions you can use:
- `env "<NAME>"`: returns the value of the environment variable with the name <NAME>
- `exec "<CMD>"`: executes the cmd in a bash shell, retrieves the stdout and strips the final EOL character

# examples

test.templ
```
there's no place like home => {{ env "HOME" }}.
today is {{ exec "cal | head -n 1" }}
```

`> templenv test.templ`
```
there's no place like home => /Users/ksophocleous.
today is     October 2016
```

# why?
I needed an easy way to generate config files from env variables and command outputs, searched for something small and not bash based and found nothing, so I created this.

essentially I wanted to generate ansible inventories that contain ip addresses from terraform outputs, like so

```
[db]
{{ exec "terraform output db.node1.ipv4.public" }} var1=somevalue
```
