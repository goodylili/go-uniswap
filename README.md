# go-uniswap
Dummy Uniswap v3 template with Go. Nothing too serious here

It's a template, not a package!

Check out [this blog](https://blog.logrocket.com/using-go-generate-reduce-boilerplate-code/) I wrote on how to use Templates with Go via `go generate`.

Currently busy with other projects, so I won't be able to document this indepth. 

You'll need abigen to generate the Uniswap v3 router and factory ABIs. 

Install it via:

```shell
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
```

I've included the Factory, ERC20 and Router ABIs in the `contract` folder.


Then generate the ABIs:

```shell
abigen --abi=contract/FactoryABI.json --pkg=uniswap --out=internal/uniswap/factory.go
```

Check the `README.md` in configurations for customizing the networks and adding more on Uniswap v3.

You'll find swap functionality in the `internal/uniswap/model.go` file. That's pretty much what this repo is about (easy to understand v3 swap functionality).

Feel free to fork and modify it to your needs.

Reach out to me on [Contacts & Socials](https://goodnessuc.github.io/about) if you have any questions.


Have fun, run wild!
