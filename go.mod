module github.con/inverse-inc/go-utils

go 1.15

require (
	github.com/cevaris/ordered_map v0.0.0-20171019141434-01ce2b16ad4f
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/google/uuid v1.1.1
	github.com/inconshreveable/log15 v0.0.0-20171019012758-0decfc6c20d9
	github.com/inverse-inc/go-utils v0.0.0-00010101000000-000000000000
	github.com/kr/pretty v0.2.1
	github.com/mattn/go-colorable v0.1.8 // indirect
	golang.org/x/net v0.0.0-20210415231046-e915ea6b2b7d
	gopkg.in/alexcesaro/statsd.v2 v2.0.0-20160320182110-7fea3f0d2fab
)

replace github.com/inverse-inc => ../

replace github.com/inverse-inc/go-utils => ./
