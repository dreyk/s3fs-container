package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
)

go_binary(
    name = "server",
    library = ":go_default_library",
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    tags = ["automanaged"],
    deps = [
        "//examples/apiserver:go_default_library",
        "//pkg/util/flag:go_default_library",
        "//pkg/util/wait:go_default_library",
        "//vendor:github.com/golang/glog",
        "//vendor:github.com/spf13/pflag",
    ],
)
