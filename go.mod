module github.com/couchbase/query

go 1.13

replace github.com/couchbase/cbauth => ../cbauth

replace github.com/couchbase/cbft => ../../../../../cbft

replace github.com/couchbase/cbftx => ../../../../../cbftx

replace github.com/couchbase/cbgt => ../../../../../cbgt

replace github.com/couchbase/eventing-ee => ../eventing-ee

replace github.com/couchbase/go-couchbase => ../go-couchbase

replace github.com/couchbase/go_json => ../go_json

replace github.com/couchbase/gomemcached => ../gomemcached

replace github.com/couchbase/indexing => ../indexing

replace github.com/couchbase/n1fty => ../n1fty

replace github.com/couchbase/plasma => ../plasma

replace github.com/couchbase/query => ./empty

replace github.com/couchbase/query-ee => ../query-ee

require (
	github.com/couchbase/cbauth v0.0.0-20201026062450-0eaf917092a2
	github.com/couchbase/clog v0.0.0-20190523192451-b8e6d5d421bc
	github.com/couchbase/eventing-ee v0.0.0-20210326071855-9c6d3f17f3a3
	github.com/couchbase/go-couchbase v0.0.0-20210330201927-1d32284da76d
	github.com/couchbase/go_json v0.0.0-20210413120532-1d2f74e0ecac
	github.com/couchbase/gocbcore-transactions v0.0.0-20210325232044-4715b797c11c
	github.com/couchbase/gocbcore/v9 v9.1.4-0.20210325182448-577aecce6dc6
	github.com/couchbase/godbc v0.0.0-20201207142944-d43b329cdf71
	github.com/couchbase/gomemcached v0.1.2
	github.com/couchbase/gometa v0.0.0-20200717102231-b0e38b71d711 // indirect
	github.com/couchbase/goutils v0.0.0-20210118111533-e33d3ffb5401
	github.com/couchbase/indexing v0.0.0-20210412045640-c428bf2b7b8e
	github.com/couchbase/n1fty v0.0.0-20210414101030-ec9f1b69eb28
	github.com/couchbase/query-ee v0.0.0-20210325205437-af39a170ba58
	github.com/couchbase/retriever v0.0.0-20150311081435-e3419088e4d3
	github.com/couchbasedeps/go-curl v0.0.0-20190830233031-f0b2afc926ec
	github.com/gorilla/mux v1.8.0
	github.com/natefinch/npipe v0.0.0-20160621034901-c1b8fa8bdcce // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/peterh/liner v1.2.0
	github.com/russross/blackfriday v1.5.2
	github.com/samuel/go-zookeeper v0.0.0-20200724154423-2164a8ac840e
	github.com/sbinet/liner v0.0.0-20150202172121-d9335eee40a4
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1
	gopkg.in/couchbase/gocb.v1 v1.6.7
	gopkg.in/couchbase/gocbcore.v7 v7.1.18 // indirect
	gopkg.in/couchbaselabs/gocbconnstr.v1 v1.0.4 // indirect
	gopkg.in/couchbaselabs/gojcbmock.v1 v1.0.4 // indirect
	gopkg.in/couchbaselabs/jsonx.v1 v1.0.0 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
)
