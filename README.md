# What is it?

A great util to format you git commit message!

It's inspired by this gif:

![git-cz](assets/commitizen.gif)

# Usage?

## Build

You need to stay in `leopard-zlc/day03/commitizen-go` first, and then:

```bash
unix> make v1
unix> sudo mv git-cz $GOPATH/bin/
```

### version 2

A new plugin was added in version, it comes from [here](https://github.com/c-bata/go-prompt)  
In this version, you need to `go mod init` first, and then :

```bash
unix> make v2
unix> sudo mv git-cz $GOPATH/bin/
```

## And then

After `git add .`, all you need to do is:  

```bash
git cz
```

That's all.
